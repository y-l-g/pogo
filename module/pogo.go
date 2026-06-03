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
	pools   map[string]*pool
	tasks   map[uint64]*task
	tasksMu sync.Mutex
	nextID  atomic.Uint64
	closed  atomic.Bool
}

type pool struct {
	name    string
	workers frankenphp.Workers
	maxWait time.Duration
	closed  atomic.Bool
}

type task struct {
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
		pools: pools,
		tasks: make(map[uint64]*task),
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

	m.tasksMu.Lock()
	tasks := m.tasks
	m.tasks = make(map[uint64]*task)
	m.tasksMu.Unlock()

	for _, task := range tasks {
		task.cancelIfSet()
		task.finish(result{err: "Pogo pool is shutting down"})
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

func (m *manager) cancel(taskID uint64) {
	if m == nil {
		return
	}

	m.tasksMu.Lock()
	task := m.tasks[taskID]
	delete(m.tasks, taskID)
	m.tasksMu.Unlock()

	if task != nil {
		task.cancelIfSet()
	}
}

func (m *manager) take(taskID uint64) (*task, bool) {
	m.tasksMu.Lock()
	task, ok := m.tasks[taskID]
	if ok {
		delete(m.tasks, taskID)
	}
	m.tasksMu.Unlock()

	return task, ok
}

func (m *manager) has(taskID uint64) bool {
	m.tasksMu.Lock()
	_, ok := m.tasks[taskID]
	m.tasksMu.Unlock()

	return ok
}

func (m *manager) store(taskID uint64, t *task) error {
	m.tasksMu.Lock()
	if m.closed.Load() {
		m.tasksMu.Unlock()
		t.cancelIfSet()
		return errors.New("Pogo pool is shutting down")
	}
	m.tasks[taskID] = t
	m.tasksMu.Unlock()
	return nil
}

func (m *manager) nextTaskID() uint64 {
	for {
		id := m.nextID.Add(1)
		if id != 0 {
			return id
		}
	}
}

func (t *task) cancelIfSet() {
	if t.cancel != nil {
		t.cancel()
	}
}

func (t *task) finish(result result) {
	if !t.done.CompareAndSwap(false, true) {
		return
	}

	select {
	case t.result <- result:
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

func awaitResult(t *task, timeout time.Duration) (string, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case result := <-t.result:
		if result.err != "" {
			return "", errors.New(result.err)
		}
		if !json.Valid([]byte(result.valueJSON)) {
			return "", errors.New("Pogo worker returned a non-JSON-compatible result")
		}
		return result.valueJSON, nil
	case <-timer.C:
		t.cancelIfSet()
		return "", fmt.Errorf("Pogo await timed out after %s", timeout)
	}
}

func (m *manager) spawn(poolName string, className string, argsJSON string) (uint64, error) {
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

	taskID := m.nextTaskID()
	ctx, cancel := context.WithTimeout(context.Background(), p.maxWait)
	t := &task{
		pool:   p.name,
		result: make(chan result, 1),
		cancel: cancel,
	}
	if err := m.store(taskID, t); err != nil {
		return 0, err
	}

	go p.run(ctx, t, payloadJSON)

	return taskID, nil
}

func (p *pool) run(ctx context.Context, t *task, payload string) {
	defer t.cancelIfSet()
	response, err := p.workers.SendMessage(ctx, payload, nil)
	if err != nil {
		t.finish(result{err: err.Error()})
		return
	}

	normalized, err := normalizeWorkerResponse(response)
	if err != nil {
		t.finish(result{err: err.Error()})
		return
	}

	if !normalized.OK {
		if normalized.Error == "" {
			normalized.Error = "Pogo worker returned an error"
		}
		t.finish(result{err: normalized.Error})
		return
	}

	if normalized.Result == nil {
		normalized.Result = json.RawMessage("null")
	}

	t.finish(result{valueJSON: string(normalized.Result)})
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

func (m *manager) await(taskID uint64, timeout time.Duration) (string, error) {
	if m == nil {
		return "", errors.New("Pogo is not configured")
	}

	task, ok := m.take(taskID)
	if !ok {
		return "", errors.New("unknown or already awaited Pogo task")
	}

	return awaitResult(task, timeout)
}

//export go_pogo_spawn
func go_pogo_spawn(poolName *C.char, poolNameLen C.size_t, className *C.char, classNameLen C.size_t, argsJSON *C.char, argsJSONLen C.size_t, errOut **C.char) C.uint64_t {
	pool := C.GoStringN(poolName, C.int(poolNameLen))
	class := C.GoStringN(className, C.int(classNameLen))
	args := C.GoStringN(argsJSON, C.int(argsJSONLen))

	taskID, err := currentManager().spawn(pool, class, args)
	if err != nil {
		*errOut = C.CString(err.Error())
		return 0
	}

	return C.uint64_t(taskID)
}

//export go_pogo_await
func go_pogo_await(taskID C.uint64_t, timeoutSeconds C.double, errOut **C.char) *C.char {
	if timeoutSeconds < 0 {
		*errOut = C.CString("Pogo await timeout must be greater than or equal to zero")
		return nil
	}

	result, err := currentManager().await(uint64(taskID), time.Duration(float64(timeoutSeconds)*float64(time.Second)))
	if err != nil {
		*errOut = C.CString(err.Error())
		return nil
	}

	return C.CString(result)
}

//export go_pogo_cancel
func go_pogo_cancel(taskID C.uint64_t) {
	m := currentManager()
	if m == nil {
		return
	}

	m.cancel(uint64(taskID))
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
