package shm

import (
	"fmt"
	"os"
	"sync"
)

// AllocationMeta tracks the lifecycle of a shared memory region.
type AllocationMeta struct {
	Offset    int64 // Logical offset
	Size      int64
	Freed     bool
	IsPadding bool
	WorkerID  int // ID of the worker owning this allocation
}

// SharedMemory implements a Ring Buffer backed by a Memory-Mapped File.
type SharedMemory struct {
	file *os.File
	data []byte
	Size int64

	mu        sync.Mutex
	head      int64             // Oldest occupied byte (logical)
	tail      int64             // Next free byte (logical)
	wasted    int64             // Bytes currently consumed by padding
	queue     []*AllocationMeta // FIFO queue of all allocations
	queueHead int               // Index of the first item in queue
	lookup    map[int64]*AllocationMeta
}

func NewSharedMemory(size int64) (*SharedMemory, error) {
	// Create a temporary file
	f, err := os.CreateTemp("", "frankenphp_shm_")
	if err != nil {
		return nil, err
	}

	// Resize it
	if err := f.Truncate(size); err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return nil, err
	}

	// Map it (OS-specific implementation)
	data, err := mapFile(f, size)
	if err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return nil, err
	}

	// Unlink immediately
	_ = os.Remove(f.Name())

	s := &SharedMemory{
		file:   f,
		data:   data,
		Size:   size,
		queue:  make([]*AllocationMeta, 0, 128),
		lookup: make(map[int64]*AllocationMeta),
		head:   8, // Reserve signature
		tail:   8,
	}

	if len(s.data) >= 8 {
		s.data[0] = 0x02
		copy(s.data[1:], []byte("GOSHM"))
	}

	return s, nil
}

func (s *SharedMemory) Close() error {
	if s.data != nil {
		_ = unmapFile(s.data)
	}
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}

func (s *SharedMemory) File() *os.File {
	return s.file
}

// compress advances the Head if the oldest allocations are Freed.
func (s *SharedMemory) compress() {
	// 1. Advance logical head
	for s.queueHead < len(s.queue) {
		meta := s.queue[s.queueHead]
		if meta.Freed {
			s.head += meta.Size
			if meta.IsPadding {
				s.wasted -= meta.Size
			}
			// Clear pointer to help GC
			s.queue[s.queueHead] = nil
			s.queueHead++
		} else {
			break
		}
	}

	// 2. Compaction (Prevent slice from growing indefinitely)
	if s.queueHead > 1024 && s.queueHead > len(s.queue)/2 {
		active := len(s.queue) - s.queueHead
		copy(s.queue, s.queue[s.queueHead:])
		s.queue = s.queue[:active]
		s.queueHead = 0
	}
}

// Allocate reserves a contiguous block in the ring buffer.
func (s *SharedMemory) Allocate(length int, workerID int) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.compress()

	len64 := int64(length)

	// 1. Check logical capacity
	used := s.tail - s.head
	free := s.Size - used

	if len64 > free {
		return 0, fmt.Errorf("shm full (capacity)")
	}

	// 2. Determine Physical Offset
	physHead := s.tail % s.Size

	// 3. Check for Wrap-Around Fragmentation
	if physHead+len64 > s.Size {
		pad := s.Size - physHead

		if len64+pad > free {
			return 0, fmt.Errorf("shm full (fragmentation)")
		}

		// Insert Padding Meta (Auto-Freed)
		padMeta := &AllocationMeta{
			Offset:    s.tail,
			Size:      pad,
			Freed:     true,
			IsPadding: true,
			WorkerID:  -1, // Padding belongs to no one
		}
		s.queue = append(s.queue, padMeta)
		s.tail += pad
		s.wasted += pad
	}

	// 4. Final Check
	used = s.tail - s.head
	free = s.Size - used
	if len64 > free {
		return 0, fmt.Errorf("shm full (post-pad)")
	}

	// 5. Commit
	offset := s.tail
	meta := &AllocationMeta{
		Offset:   offset,
		Size:     len64,
		Freed:    false,
		WorkerID: workerID,
	}

	s.queue = append(s.queue, meta)
	s.lookup[offset%s.Size] = meta
	s.tail += len64

	return offset % s.Size, nil
}

func (s *SharedMemory) Free(physOffset int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if meta, ok := s.lookup[physOffset]; ok {
		meta.Freed = true
		delete(s.lookup, physOffset)
		s.compress()
	}
}

// FreeByWorkerID reclaims all allocations owned by a specific worker.
// This is O(N) where N is the number of active allocations.
// It is intended for cleanup during worker crash/shutdown.
func (s *SharedMemory) FreeByWorkerID(workerID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for physOffset, meta := range s.lookup {
		if meta.WorkerID == workerID {
			meta.Freed = true
			delete(s.lookup, physOffset)
		}
	}
	s.compress()
}

func (s *SharedMemory) WriteAt(offset int64, data []byte) error {
	if offset < 0 || offset+int64(len(data)) > s.Size {
		return fmt.Errorf("write out of bounds")
	}
	copy(s.data[offset:], data)
	return nil
}

// Stats structure
type ShmStats struct {
	TotalBytes  int64
	UsedBytes   int64
	FreeBytes   int64
	WastedBytes int64
}

func (s *SharedMemory) GetStats() ShmStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	used := s.tail - s.head
	return ShmStats{
		TotalBytes:  s.Size,
		UsedBytes:   used,
		FreeBytes:   s.Size - used,
		WastedBytes: s.wasted,
	}
}
