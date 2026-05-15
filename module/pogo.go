package pogo

/*
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/dunglas/frankenphp"
)

const defaultPoolName = "default"

var (
	globalManager   *manager
	globalManagerMu sync.RWMutex
)

type manager struct {
	pools     map[string]*pool
	handles   map[uint64]*pending
	handlesMu sync.Mutex
	nextID    atomic.Uint64
	closed    atomic.Bool
}

type pool struct {
	name    string
	workers frankenphp.Workers
	maxWait time.Duration
	closed  atomic.Bool
}

type pending struct {
	pool   string
	result chan result
	cancel context.CancelFunc
	done   atomic.Bool
}

type result struct {
	valueJSON string
	err       string
}

type workerEnvelope struct {
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

func newManager(pools map[string]*pool) *manager {
	return &manager{
		pools:   pools,
		handles: make(map[uint64]*pending),
	}
}

func newPool(name string, workers frankenphp.Workers, maxWait time.Duration) *pool {
	return &pool{
		name:    name,
		workers: workers,
		maxWait: maxWait,
	}
}

func currentManager() *manager {
	globalManagerMu.RLock()
	m := globalManager
	globalManagerMu.RUnlock()
	return m
}

func (m *manager) close() {
	m.closed.Store(true)
	for _, p := range m.pools {
		p.closed.Store(true)
	}

	m.handlesMu.Lock()
	handles := m.handles
	m.handles = make(map[uint64]*pending)
	m.handlesMu.Unlock()

	for _, pending := range handles {
		pending.cancelIfSet()
		pending.finish(result{err: "Pogo pool is shutting down"})
	}
}

func (m *manager) pool(name string) (*pool, error) {
	if m == nil || m.closed.Load() {
		return nil, errors.New("Pogo is not configured")
	}

	if name == "" {
		name = defaultPoolName
	}

	p := m.pools[name]
	if p == nil || p.workers == nil || p.closed.Load() {
		return nil, fmt.Errorf("unknown Pogo pool %q", name)
	}

	return p, nil
}

func (m *manager) cancel(handle uint64) {
	if m == nil {
		return
	}

	m.handlesMu.Lock()
	pending := m.handles[handle]
	delete(m.handles, handle)
	m.handlesMu.Unlock()

	if pending != nil {
		pending.cancelIfSet()
	}
}

func (m *manager) take(handle uint64) (*pending, bool) {
	m.handlesMu.Lock()
	pending, ok := m.handles[handle]
	if ok {
		delete(m.handles, handle)
	}
	m.handlesMu.Unlock()

	return pending, ok
}

func (m *manager) has(handle uint64) bool {
	m.handlesMu.Lock()
	_, ok := m.handles[handle]
	m.handlesMu.Unlock()

	return ok
}

func (m *manager) store(handle uint64, pending *pending) error {
	m.handlesMu.Lock()
	if m.closed.Load() {
		m.handlesMu.Unlock()
		pending.cancelIfSet()
		return errors.New("Pogo pool is shutting down")
	}
	m.handles[handle] = pending
	m.handlesMu.Unlock()
	return nil
}

func (m *manager) nextHandle() uint64 {
	for {
		id := m.nextID.Add(1)
		if id != 0 {
			return id
		}
	}
}

func (p *pending) cancelIfSet() {
	if p.cancel != nil {
		p.cancel()
	}
}

func (p *pending) finish(result result) {
	if !p.done.CompareAndSwap(false, true) {
		return
	}

	select {
	case p.result <- result:
	default:
	}
}

func buildPayloadJSON(className string, argsJSON string) (string, error) {
	classJSON, err := json.Marshal(className)
	if err != nil {
		return "", fmt.Errorf("failed to encode Pogo job class: %w", err)
	}

	return fmt.Sprintf(`{"class":%s,"args":%s}`, classJSON, argsJSON), nil
}

func awaitResult(pending *pending, timeout time.Duration) (string, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case result := <-pending.result:
		if result.err != "" {
			return "", errors.New(result.err)
		}
		if !json.Valid([]byte(result.valueJSON)) {
			return "", errors.New("Pogo worker returned a non-JSON-compatible result")
		}
		return result.valueJSON, nil
	case <-timer.C:
		pending.cancelIfSet()
		return "", fmt.Errorf("Pogo await timed out after %s", timeout)
	}
}

func (m *manager) dispatch(poolName string, className string, argsJSON string) (uint64, error) {
	p, err := m.pool(poolName)
	if err != nil {
		return 0, err
	}

	if className == "" {
		return 0, errors.New("Pogo job class must not be empty")
	}

	if !json.Valid([]byte(argsJSON)) {
		return 0, errors.New("Pogo job args must be JSON-compatible")
	}

	payloadJSON, err := buildPayloadJSON(className, argsJSON)
	if err != nil {
		return 0, err
	}

	handle := m.nextHandle()
	ctx, cancel := context.WithTimeout(context.Background(), p.maxWait)
	pending := &pending{
		pool:   p.name,
		result: make(chan result, 1),
		cancel: cancel,
	}
	if err := m.store(handle, pending); err != nil {
		return 0, err
	}

	go p.run(ctx, pending, payloadJSON)

	return handle, nil
}

func (p *pool) run(ctx context.Context, pending *pending, payload string) {
	defer pending.cancelIfSet()
	response, err := p.workers.SendMessage(ctx, payload, nil)
	if err != nil {
		pending.finish(result{err: err.Error()})
		return
	}

	normalized, err := normalizeWorkerResponse(response)
	if err != nil {
		pending.finish(result{err: err.Error()})
		return
	}

	if !normalized.OK {
		if normalized.Error == "" {
			normalized.Error = "Pogo worker returned an error"
		}
		pending.finish(result{err: normalized.Error})
		return
	}

	if normalized.Result == nil {
		normalized.Result = json.RawMessage("null")
	}

	pending.finish(result{valueJSON: string(normalized.Result)})
}

func normalizeWorkerResponse(response any) (workerEnvelope, error) {
	switch value := response.(type) {
	case string:
		return decodeEnvelope([]byte(value))
	case []byte:
		return decodeEnvelope(value)
	case map[string]any:
		data, err := json.Marshal(value)
		if err != nil {
			return workerEnvelope{}, fmt.Errorf("Pogo worker returned a non-JSON-compatible response: %w", err)
		}
		return decodeEnvelope(data)
	case frankenphp.AssociativeArray[any]:
		data, err := json.Marshal(value.Map)
		if err != nil {
			return workerEnvelope{}, fmt.Errorf("Pogo worker returned a non-JSON-compatible response: %w", err)
		}
		return decodeEnvelope(data)
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return workerEnvelope{}, fmt.Errorf("Pogo worker returned unsupported response type %T", response)
		}
		return decodeEnvelope(data)
	}
}

func decodeEnvelope(data []byte) (workerEnvelope, error) {
	var envelope workerEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return workerEnvelope{}, fmt.Errorf("Pogo worker returned an invalid response envelope: %w", err)
	}
	return envelope, nil
}

func (m *manager) await(handle uint64, timeout time.Duration) (string, error) {
	if m == nil {
		return "", errors.New("Pogo is not configured")
	}

	pending, ok := m.take(handle)
	if !ok {
		return "", errors.New("unknown or already awaited Pogo handle")
	}

	return awaitResult(pending, timeout)
}

//export go_pogo_dispatch
func go_pogo_dispatch(poolName *C.char, poolNameLen C.size_t, className *C.char, classNameLen C.size_t, argsJSON *C.char, argsJSONLen C.size_t, errOut **C.char) C.uint64_t {
	pool := C.GoStringN(poolName, C.int(poolNameLen))
	class := C.GoStringN(className, C.int(classNameLen))
	args := C.GoStringN(argsJSON, C.int(argsJSONLen))

	handle, err := currentManager().dispatch(pool, class, args)
	if err != nil {
		*errOut = C.CString(err.Error())
		return 0
	}

	return C.uint64_t(handle)
}

//export go_pogo_await
func go_pogo_await(handle C.uint64_t, timeoutSeconds C.double, errOut **C.char) *C.char {
	if timeoutSeconds < 0 {
		*errOut = C.CString("Pogo await timeout must be greater than or equal to zero")
		return nil
	}

	result, err := currentManager().await(uint64(handle), time.Duration(float64(timeoutSeconds)*float64(time.Second)))
	if err != nil {
		*errOut = C.CString(err.Error())
		return nil
	}

	return C.CString(result)
}

//export go_pogo_cancel
func go_pogo_cancel(handle C.uint64_t) {
	m := currentManager()
	if m == nil {
		return
	}

	m.cancel(uint64(handle))
}

//export go_pogo_pool_size
func go_pogo_pool_size(poolName *C.char, poolNameLen C.size_t) C.int {
	pool := C.GoStringN(poolName, C.int(poolNameLen))
	p, err := currentManager().pool(pool)
	if err != nil {
		return 0
	}

	return C.int(p.workers.NumThreads())
}

var _ = unsafe.Pointer(nil)
