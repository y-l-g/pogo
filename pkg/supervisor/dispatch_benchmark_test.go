package supervisor

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/ugorji/go/codec"
)

// BenchmarkSerialization compares JSON vs MsgPack for a typical job payload.
// This validates the decision to use MsgPack for binary transport optimization.
func BenchmarkSerialization(b *testing.B) {
	// A typical payload: class name, nested data, mixed types
	payload := map[string]any{
		"job_class": "App\\Jobs\\ProcessOrder",
		"payload": map[string]any{
			"order_id": 12345,
			"user_id":  9876,
			"items": []any{
				map[string]any{"id": 1, "qty": 2, "price": 10.50},
				map[string]any{"id": 55, "qty": 1, "price": 199.99},
			},
			"options": map[string]any{
				"send_email": true,
				"retry":      false,
			},
		},
		"meta": "some metadata string to pad the size slightly",
	}

	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[string]any(nil))

	b.Run("JSON", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if _, err := json.Marshal(payload); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("MsgPack", func(b *testing.B) {
		var buf []byte
		// Mimic the logic in pool.go: executeOnWorker
		for i := 0; i < b.N; i++ {
			buf = buf[:0]
			enc := codec.NewEncoderBytes(&buf, &mh)
			if err := enc.Encode(payload); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkHandleValidation measures the cost of the security scan
// which runs on every dispatch to prevent handle hijacking.
func BenchmarkHandleValidation(b *testing.B) {
	p := NewPool(0)
	defer p.Shutdown()

	// Nested payload to force recursive traversal
	payload := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": map[string]any{
					"data": "noop",
					// integer mimicking a handle
					"maybe_handle": uint64(123456),
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := p.ValidateHandles(payload); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkInternalBus measures the latency of the goWorkerLoop processing a task.
// This isolates the Go Runtime overhead from the PHP Process overhead.
func BenchmarkInternalBus(b *testing.B) {
	p := NewPool(0)
	defer p.Shutdown()

	done := make(chan struct{})

	// Inject a dummy internal worker
	// Accessing p.registry is allowed because we are in package supervisor
	p.registry["bench.noop"] = func(payload map[string]any) {
		done <- struct{}{}
	}

	task := GoTask{
		Name:       "bench.noop",
		Payload:    map[string]any{},
		EnqueuedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// We must verify the waitgroup counter logic matches dispatch_to_pool_wrapper
		p.Wg().Add(1)
		p.Tasks() <- task
		<-done
	}
}
