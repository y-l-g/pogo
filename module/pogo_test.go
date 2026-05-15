package pogo

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

type fakeWorkers struct {
	delay    time.Duration
	response any
	err      error
	threads  int
	calls    atomic.Int64
	canceled chan struct{}
}

func (f *fakeWorkers) SendRequest(http.ResponseWriter, *http.Request) error {
	return nil
}

func (f *fakeWorkers) SendMessage(ctx context.Context, message any, _ http.ResponseWriter) (any, error) {
	f.calls.Add(1)
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			if f.canceled != nil {
				close(f.canceled)
			}
			return nil, ctx.Err()
		}
	}
	if f.err != nil {
		return nil, f.err
	}
	return f.response, nil
}

func (f *fakeWorkers) NumThreads() int {
	return f.threads
}

func testManager(pools map[string]*pool) *manager {
	return newManager(pools)
}

func TestDispatchUsesNamedPool(t *testing.T) {
	defaultWorkers := &fakeWorkers{response: `{"ok":true,"result":"default"}`}
	apiWorkers := &fakeWorkers{response: `{"ok":true,"result":"api"}`}
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, defaultWorkers, time.Second),
		"external_api":  newPool("external_api", apiWorkers, time.Second),
	})

	handle, err := m.dispatch("external_api", "App\\Job", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	result, err := m.await(handle, time.Second)
	if err != nil {
		t.Fatalf("await failed: %v", err)
	}

	if result != `"api"` {
		t.Fatalf("unexpected result: %s", result)
	}
	if apiWorkers.calls.Load() != 1 {
		t.Fatalf("expected external_api pool to be used")
	}
	if defaultWorkers.calls.Load() != 0 {
		t.Fatalf("default pool should not have been used")
	}
}

func TestUnknownPoolDispatchFails(t *testing.T) {
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, &fakeWorkers{}, time.Second),
	})

	if _, err := m.dispatch("missing", "App\\Job", `[]`); err == nil {
		t.Fatal("expected unknown pool error")
	}
}

func TestHandlesAreGloballyUniqueAcrossPools(t *testing.T) {
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, &fakeWorkers{response: `{"ok":true,"result":1}`}, time.Second),
		"cpu":           newPool("cpu", &fakeWorkers{response: `{"ok":true,"result":2}`}, time.Second),
	})

	a, err := m.dispatch(defaultPoolName, "App\\Job", `[]`)
	if err != nil {
		t.Fatalf("dispatch default failed: %v", err)
	}
	b, err := m.dispatch("cpu", "App\\Job", `[]`)
	if err != nil {
		t.Fatalf("dispatch cpu failed: %v", err)
	}

	if a == b {
		t.Fatal("handles from different pools must be globally unique")
	}
}

func TestAwaitTimeoutDeletesHandleAndLateResultDoesNotBlock(t *testing.T) {
	workers := &fakeWorkers{
		delay:    50 * time.Millisecond,
		response: `{"ok":true,"result":"late"}`,
		threads:  1,
	}
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, workers, time.Second),
	})

	handle, err := m.dispatch(defaultPoolName, "App\\SlowJob", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	if _, err := m.await(handle, time.Millisecond); err == nil {
		t.Fatal("expected await timeout")
	}

	if m.has(handle) {
		t.Fatal("timed out handle was not deleted")
	}

	time.Sleep(100 * time.Millisecond)
}

func TestCancelDeletesHandleAndCancelsWorkerContext(t *testing.T) {
	canceled := make(chan struct{})
	workers := &fakeWorkers{
		delay:    time.Second,
		response: `{"ok":true,"result":"late"}`,
		canceled: canceled,
	}
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, workers, time.Second),
	})

	handle, err := m.dispatch(defaultPoolName, "App\\SlowJob", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	m.cancel(handle)

	if m.has(handle) {
		t.Fatal("canceled handle was not deleted")
	}

	select {
	case <-canceled:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("worker context was not canceled")
	}
}

func TestAwaitUnknownHandleFails(t *testing.T) {
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, &fakeWorkers{}, time.Second),
	})

	if _, err := m.await(42, time.Millisecond); err == nil {
		t.Fatal("expected unknown handle error")
	}
}

func TestAwaitSameHandleTwiceFails(t *testing.T) {
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, &fakeWorkers{response: `{"ok":true,"result":123}`}, time.Second),
	})

	handle, err := m.dispatch(defaultPoolName, "App\\Job", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	if result, err := m.await(handle, time.Second); err != nil || result != "123" {
		t.Fatalf("unexpected first await result=%q err=%v", result, err)
	}

	if _, err := m.await(handle, time.Second); err == nil {
		t.Fatal("expected second await to fail")
	}
}

func TestWorkerErrorPropagates(t *testing.T) {
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, &fakeWorkers{err: errors.New("worker failed")}, time.Second),
	})

	handle, err := m.dispatch(defaultPoolName, "App\\Job", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	if _, err := m.await(handle, time.Second); err == nil || err.Error() != "worker failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPoolSizeUsesSelectedPool(t *testing.T) {
	m := testManager(map[string]*pool{
		defaultPoolName: newPool(defaultPoolName, &fakeWorkers{threads: 2}, time.Second),
		"external_api":  newPool("external_api", &fakeWorkers{threads: 7}, time.Second),
	})

	p, err := m.pool("external_api")
	if err != nil {
		t.Fatalf("pool lookup failed: %v", err)
	}

	if p.workers.NumThreads() != 7 {
		t.Fatalf("unexpected pool size: %d", p.workers.NumThreads())
	}
}
