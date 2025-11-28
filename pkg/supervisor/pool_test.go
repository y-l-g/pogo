package supervisor

import (
	"os"
	"testing"
	"time"
)

// TestObserver captures events for verification in tests.
type TestObserver struct {
	WorkerStarted chan int
	WorkerExited  chan int
}

func NewTestObserver() *TestObserver {
	return &TestObserver{
		WorkerStarted: make(chan int, 100),
		WorkerExited:  make(chan int, 100),
	}
}

func (o *TestObserver) OnWorkerStart(workerID int) {
	select {
	case o.WorkerStarted <- workerID:
	default:
	}
}

func (o *TestObserver) OnWorkerExit(workerID int) {
	select {
	case o.WorkerExited <- workerID:
	default:
	}
}

func setEnvOrFatal(t *testing.T, key, value string) {
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set env var %s: %v", key, err)
	}
}

func unsetEnvOrFatal(t *testing.T, key string) {
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Failed to unset env var %s: %v", key, err)
	}
}

func TestCrashResilience(t *testing.T) {
	// Locate the test binary itself to use as the worker
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("Failed to get executable: %v", err)
	}

	setEnvOrFatal(t, "POGO_TEST_PHP_BINARY", exe)
	setEnvOrFatal(t, "GO_WANT_HELPER_PROCESS", "1")
	setEnvOrFatal(t, "POGO_MOCK_WORKER_MODE", "crash_immediate")

	defer func() {
		unsetEnvOrFatal(t, "POGO_TEST_PHP_BINARY")
		unsetEnvOrFatal(t, "GO_WANT_HELPER_PROCESS")
		unsetEnvOrFatal(t, "POGO_MOCK_WORKER_MODE")
	}()

	p := NewPool(999)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CRITICAL FAILURE: Supervisor Panicked: %v", r)
		}
		p.Shutdown()
	}()

	observer := NewTestObserver()

	cfg := PoolConfig{
		ShmSize:      1024 * 1024,
		IpcTimeout:   100 * time.Millisecond,
		ScaleLatency: 50,
		Observer:     observer,
	}

	t.Log("Starting Pool with Mock Worker (Crash Mode)...")
	p.Start("-test.run=TestHelperProcess", 1, 1, 0, cfg)

	// Wait for the first worker to die
	select {
	case id := <-observer.WorkerExited:
		t.Logf("Worker #%d died as expected.", id)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for initial worker crash")
	}

	// The supervisor should try to respawn.
	// Wait for the second worker to die (proving the loop is active)
	select {
	case id := <-observer.WorkerExited:
		t.Logf("Replacement Worker #%d died as expected.", id)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for replacement worker crash")
	}

	stats := p.GetStats()
	t.Logf("Stats: %+v", stats)
	t.Log("Supervisor survived crash cycle without hanging.")
}

func TestNormalOperation(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("Failed to get executable: %v", err)
	}

	setEnvOrFatal(t, "POGO_TEST_PHP_BINARY", exe)
	setEnvOrFatal(t, "GO_WANT_HELPER_PROCESS", "1")
	setEnvOrFatal(t, "POGO_MOCK_WORKER_MODE", "normal")

	defer func() {
		unsetEnvOrFatal(t, "POGO_TEST_PHP_BINARY")
		unsetEnvOrFatal(t, "GO_WANT_HELPER_PROCESS")
		unsetEnvOrFatal(t, "POGO_MOCK_WORKER_MODE")
	}()

	p := NewPool(1000)
	defer p.Shutdown()

	observer := NewTestObserver()

	cfg := PoolConfig{
		ShmSize:      1024 * 1024,
		IpcTimeout:   2000 * time.Millisecond,
		ScaleLatency: 50,
		Observer:     observer,
	}

	p.Start("-test.run=TestHelperProcess", 1, 1, 0, cfg)

	// Wait for boot via observer instead of sleep
	select {
	case id := <-observer.WorkerStarted:
		t.Logf("Worker #%d started successfully.", id)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for worker start")
	}

	stats := p.GetStats()
	if stats["total_workers"] != 1 {
		t.Errorf("Expected 1 worker, got %v", stats["total_workers"])
	}
}
