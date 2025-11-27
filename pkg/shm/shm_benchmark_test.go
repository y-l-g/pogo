package shm

import (
	"crypto/rand"
	"testing"
)

// BenchmarkAllocate measures the cost of sequential allocations.
// This tests the efficiency of the ring buffer logic and map metadata insertions.
func BenchmarkAllocate(b *testing.B) {
	// Setup: Create a 64MB SHM
	shm, err := NewSharedMemory(64 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Failed to init SHM: %v", err)
	}
	defer shm.Close()

	payloadSize := 1024 // 1KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		offset, err := shm.Allocate(payloadSize)
		if err != nil {
			// If buffer fills up (unlikely with just 1KB in benchmarks unless N is huge without free),
			// we free everything to reset.
			// However, standard Benchmark loop usually doesn't free.
			// For a Ring Buffer, we must free to prevent OOM errors in long benchmarks.
			b.StopTimer()
			// Reset head/tail/queue manual cleanup or just fail if logic assumes infinite space
			// Simple approach: Free immediately to simulate throughput
			b.Fatalf("Allocation failed at i=%d: %v", i, err)
		}
		shm.Free(offset)
	}
}

// BenchmarkAllocateParallel measures contention on the Mutex.
func BenchmarkAllocateParallel(b *testing.B) {
	shm, err := NewSharedMemory(64 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Failed to init SHM: %v", err)
	}
	defer shm.Close()

	payloadSize := 512 // 512 Bytes

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			offset, err := shm.Allocate(payloadSize)
			if err != nil {
				// In parallel, we can't easily fatal without stopping others,
				// but let's assume 64MB is enough for the benchmark window.
				// We definitely need to free to sustain the ring.
				continue
			}
			shm.Free(offset)
		}
	})
}

// BenchmarkWriteAt measures the cost of copying data into the mmap region.
func BenchmarkWriteAt(b *testing.B) {
	shm, err := NewSharedMemory(10 * 1024 * 1024)
	if err != nil {
		b.Fatalf("Failed to init SHM: %v", err)
	}
	defer shm.Close()

	data := make([]byte, 4096) // 4KB
	_, _ = rand.Read(data)
	offset, _ := shm.Allocate(len(data))

	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := shm.WriteAt(offset, data); err != nil {
			b.Fatalf("WriteAt failed: %v", err)
		}
	}
}
