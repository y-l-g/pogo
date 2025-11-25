package supervisor

import (
	"os"
	"testing"
	"time"
)

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

	// Set env vars for the pool to pick up
	setEnvOrFatal(t, "POGO_TEST_PHP_BINARY", exe)
	// Set env vars for the worker process to behave correctly
	setEnvOrFatal(t, "GO_WANT_HELPER_PROCESS", "1")
	setEnvOrFatal(t, "POGO_MOCK_WORKER_MODE", "crash_immediate")

	defer func() {
		unsetEnvOrFatal(t, "POGO_TEST_PHP_BINARY")
		unsetEnvOrFatal(t, "GO_WANT_HELPER_PROCESS")
		unsetEnvOrFatal(t, "POGO_MOCK_WORKER_MODE")
	}()

	// 2. Initialize Pool
	p := NewPool(999)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CRITICAL FAILURE: Supervisor Panicked: %v", r)
		}
		p.Shutdown()
	}()

	cfg := PoolConfig{
		ShmSize:      1024 * 1024,
		IpcTimeout:   100 * time.Millisecond,
		ScaleLatency: 50,
	}

	t.Log("Starting Pool with Mock Worker (Crash Mode)...")
	// The entrypoint becomes the flag passed to the binary
	p.Start("-test.run=TestHelperProcess", 1, 1, 0, cfg)

	// 3. Wait for the crash cycle to stabilize
	// The backoff is 500ms, so 1.5s allows ~2 restart attempts.
	time.Sleep(1500 * time.Millisecond)

	stats := p.GetStats()
	// We expect workers to be crashing, so active might be 0 or 1 depending on timing,
	// but the Supervisor must still be alive.
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
	setEnvOrFatal(t, "POGO_MOCK_WORKER_MODE", "normal") // Normal echo mode

	defer func() {
		unsetEnvOrFatal(t, "POGO_TEST_PHP_BINARY")
		unsetEnvOrFatal(t, "GO_WANT_HELPER_PROCESS")
		unsetEnvOrFatal(t, "POGO_MOCK_WORKER_MODE")
	}()

	p := NewPool(1000)
	defer p.Shutdown()

	cfg := PoolConfig{
		ShmSize:      1024 * 1024,
		IpcTimeout:   2000 * time.Millisecond,
		ScaleLatency: 50,
	}

	p.Start("-test.run=TestHelperProcess", 1, 1, 0, cfg)

	// Wait for boot
	time.Sleep(500 * time.Millisecond)

	stats := p.GetStats()
	if stats["total_workers"] != 1 {
		t.Errorf("Expected 1 worker, got %v", stats["total_workers"])
	}
}
