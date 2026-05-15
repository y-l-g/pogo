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

func TestAwaitTimeoutDeletesHandleAndLateResultDoesNotBlock(t *testing.T) {
	workers := &fakeWorkers{
		delay:    50 * time.Millisecond,
		response: `{"ok":true,"result":"late"}`,
		threads:  1,
	}
	p := newPool(workers, time.Second)

	handle, err := p.dispatch("App\\SlowJob", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	if _, err := p.await(handle, time.Millisecond); err == nil {
		t.Fatal("expected await timeout")
	}

	if p.has(handle) {
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
	p := newPool(workers, time.Second)

	handle, err := p.dispatch("App\\SlowJob", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	p.cancel(handle)

	if p.has(handle) {
		t.Fatal("canceled handle was not deleted")
	}

	select {
	case <-canceled:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("worker context was not canceled")
	}
}

func TestAwaitUnknownHandleFails(t *testing.T) {
	p := newPool(&fakeWorkers{}, time.Second)

	if _, err := p.await(42, time.Millisecond); err == nil {
		t.Fatal("expected unknown handle error")
	}
}

func TestAwaitSameHandleTwiceFails(t *testing.T) {
	p := newPool(&fakeWorkers{response: `{"ok":true,"result":123}`}, time.Second)

	handle, err := p.dispatch("App\\Job", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	if result, err := p.await(handle, time.Second); err != nil || result != "123" {
		t.Fatalf("unexpected first await result=%q err=%v", result, err)
	}

	if _, err := p.await(handle, time.Second); err == nil {
		t.Fatal("expected second await to fail")
	}
}

func TestWorkerErrorPropagates(t *testing.T) {
	p := newPool(&fakeWorkers{err: errors.New("worker failed")}, time.Second)

	handle, err := p.dispatch("App\\Job", `[]`)
	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	if _, err := p.await(handle, time.Second); err == nil || err.Error() != "worker failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}
