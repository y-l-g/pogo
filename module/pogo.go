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

var (
	globalPool   *pool
	globalPoolMu sync.RWMutex
)

type pool struct {
	workers   frankenphp.Workers
	handles   map[uint64]*pending
	handlesMu sync.Mutex
	nextID    atomic.Uint64
	maxWait   time.Duration
	closed    atomic.Bool
}

type pending struct {
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

func newPool(workers frankenphp.Workers, maxWait time.Duration) *pool {
	return &pool{
		workers: workers,
		handles: make(map[uint64]*pending),
		maxWait: maxWait,
	}
}

func currentPool() *pool {
	globalPoolMu.RLock()
	p := globalPool
	globalPoolMu.RUnlock()
	return p
}

func (p *pool) close() {
	p.closed.Store(true)

	p.handlesMu.Lock()
	handles := p.handles
	p.handles = make(map[uint64]*pending)
	p.handlesMu.Unlock()

	for _, pending := range handles {
		pending.cancelIfSet()
		pending.finish(result{err: "Pogo pool is shutting down"})
	}
}

func (p *pool) cancel(handle uint64) {
	if p == nil {
		return
	}

	p.handlesMu.Lock()
	pending := p.handles[handle]
	delete(p.handles, handle)
	p.handlesMu.Unlock()

	if pending != nil {
		pending.cancelIfSet()
	}
}

func (p *pool) take(handle uint64) (*pending, bool) {
	p.handlesMu.Lock()
	pending, ok := p.handles[handle]
	if ok {
		delete(p.handles, handle)
	}
	p.handlesMu.Unlock()

	return pending, ok
}

func (p *pool) has(handle uint64) bool {
	p.handlesMu.Lock()
	_, ok := p.handles[handle]
	p.handlesMu.Unlock()

	return ok
}

func (p *pool) store(handle uint64, pending *pending) error {
	p.handlesMu.Lock()
	if p.closed.Load() {
		p.handlesMu.Unlock()
		pending.cancelIfSet()
		return errors.New("Pogo pool is shutting down")
	}
	p.handles[handle] = pending
	p.handlesMu.Unlock()
	return nil
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

func (p *pool) awaitResult(pending *pending, timeout time.Duration) (string, error) {
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

func (p *pool) dispatch(className string, argsJSON string) (uint64, error) {
	if p == nil || p.workers == nil || p.closed.Load() {
		return 0, errors.New("Pogo is not configured")
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

	handle := p.nextHandle()
	ctx, cancel := context.WithTimeout(context.Background(), p.maxWait)
	pending := &pending{
		result: make(chan result, 1),
		cancel: cancel,
	}
	if err := p.store(handle, pending); err != nil {
		return 0, err
	}

	go p.run(ctx, pending, payloadJSON)

	return handle, nil
}

func (p *pool) nextHandle() uint64 {
	for {
		id := p.nextID.Add(1)
		if id != 0 {
			return id
		}
	}
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

func (p *pool) await(handle uint64, timeout time.Duration) (string, error) {
	if p == nil {
		return "", errors.New("Pogo is not configured")
	}

	pending, ok := p.take(handle)
	if !ok {
		return "", errors.New("unknown or already awaited Pogo handle")
	}

	return p.awaitResult(pending, timeout)
}

//export go_pogo_dispatch
func go_pogo_dispatch(className *C.char, classNameLen C.size_t, argsJSON *C.char, argsJSONLen C.size_t, errOut **C.char) C.uint64_t {
	class := C.GoStringN(className, C.int(classNameLen))
	args := C.GoStringN(argsJSON, C.int(argsJSONLen))

	handle, err := currentPool().dispatch(class, args)
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

	result, err := currentPool().await(uint64(handle), time.Duration(float64(timeoutSeconds)*float64(time.Second)))
	if err != nil {
		*errOut = C.CString(err.Error())
		return nil
	}

	return C.CString(result)
}

//export go_pogo_cancel
func go_pogo_cancel(handle C.uint64_t) {
	p := currentPool()
	if p == nil {
		return
	}

	p.cancel(uint64(handle))
}

//export go_pogo_pool_size
func go_pogo_pool_size() C.int {
	p := currentPool()
	if p == nil || p.workers == nil {
		return 0
	}

	return C.int(p.workers.NumThreads())
}

var _ = unsafe.Pointer(nil)
