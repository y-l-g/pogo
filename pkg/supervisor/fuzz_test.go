package supervisor

import (
	"io"
	"os"
	"testing"
	"time"
)

// MockWorkerRW simulates the pipe connection
type MockWorkerRW struct {
	Reader io.Reader
	Writer io.Writer
}

func (m *MockWorkerRW) Read(p []byte) (n int, err error) {
	return m.Reader.Read(p)
}
func (m *MockWorkerRW) Write(p []byte) (n int, err error) {
	return m.Writer.Write(p)
}
func (m *MockWorkerRW) Close() error { return nil }

func FuzzPacketReading(f *testing.F) {
	// Seed some valid packets
	// 1. Header: Len=2, Type=Data (0x00). Body: "{}"
	seed1 := []byte{0, 0, 0, 2, 0, '{', '}'}
	f.Add(seed1)

	// 2. Malformed Large Packet
	seed2 := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0}
	f.Add(seed2)

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 5 {
			return
		}

		// Setup Pipes
		r, w, _ := os.Pipe()
		defer func() { _ = r.Close() }()
		defer func() { _ = w.Close() }()

		// Write fuzz data to pipe asynchronously to prevent blocking
		go func() {
			if _, err := w.Write(data); err != nil {
				return
			}
			// Don't close immediately, let the reader try to read
			time.Sleep(10 * time.Millisecond)
			_ = w.Close()
		}()

		// Setup minimal Worker with Transport
		worker := &phpWorker{
			id:        1,
			transport: NewTransport(r, w),
			// Process is nil, which is handled by killWorker
		}

		// Setup Minimal Pool
		p := &Pool{
			config: PoolConfig{
				IpcTimeout: 10 * time.Millisecond, // Fast timeout for fuzzing
				JobTimeout: 10 * time.Millisecond,
			},
			// Initialize the workers channel to prevent blocking on return
			workers: make(chan *phpWorker, 1),
		}

		p.workersList = make(map[int]*phpWorker)
		p.workersList[1] = worker

		r_local, w_remote, _ := os.Pipe()
		r_remote, w_local, _ := os.Pipe()

		worker.transport = NewTransport(r_local, w_local)

		go func() {
			if _, err := w_remote.Write(data); err != nil {
				_ = w_remote.Close()
				_ = r_remote.Close()
				return
			}
			_ = w_remote.Close()
			_, _ = io.Copy(io.Discard, r_remote)
			_ = r_remote.Close()
		}()

		payload := map[string]any{"fuzz": true}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Panic detected on input %v: %v", data, r)
			}
		}()

		p.executeOnWorker(worker, payload, nil)

		_ = worker.transport.Close()
	})
}
