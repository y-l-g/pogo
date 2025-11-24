package pogo

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

// AllocationMeta tracks the lifecycle of a shared memory region.
type AllocationMeta struct {
	Offset int64 // Logical offset (always increases)
	Size   int64
	Freed  bool
}

// SharedMemory implements a Ring Buffer backed by a Memory-Mapped File.
type SharedMemory struct {
	file *os.File
	data []byte
	Size int64

	mu          sync.Mutex
	head        int64            // Oldest occupied byte (logical)
	tail        int64            // Next free byte (logical)
	allocations []AllocationMeta // FIFO queue of active allocations
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

	// Map it
	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return nil, err
	}

	// Unlink immediately (file stays open, but name is gone)
	_ = os.Remove(f.Name())

	s := &SharedMemory{
		file:        f,
		data:        data,
		Size:        size,
		allocations: make([]AllocationMeta, 0, 128),
		// Reserve first 8 bytes for Signature "GOSHM"
		head: 8,
		tail: 8,
	}

	// Initialize signature
	if len(s.data) >= 8 {
		s.data[0] = 0x02
		copy(s.data[1:], []byte("GOSHM"))
	}

	return s, nil
}

func (s *SharedMemory) Close() error {
	if s.data != nil {
		_ = syscall.Munmap(s.data)
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
	consumedCount := 0
	for i := range s.allocations {
		if s.allocations[i].Freed {
			s.head += s.allocations[i].Size
			consumedCount++
		} else {
			break
		}
	}

	if consumedCount > 0 {
		// Drop freed items
		// Optimization: Re-slice. Underlying array will be reused when appended to.
		s.allocations = s.allocations[consumedCount:]
	}
}

// Allocate reserves a contiguous block in the ring buffer.
func (s *SharedMemory) Allocate(length int) (int64, error) {
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
		// We must wrap to 0.
		pad := s.Size - physHead

		// Ensure space for Padding + Data
		if len64+pad > free {
			return 0, fmt.Errorf("shm full (fragmentation)")
		}

		// Register Dummy Allocation for Padding (Auto-freed)
		// It blocks 'Head' until it reaches the front, then instantly vanishes.
		s.allocations = append(s.allocations, AllocationMeta{
			Offset: s.tail,
			Size:   pad,
			Freed:  true,
		})

		s.tail += pad
	}

	// 4. Final Check
	used = s.tail - s.head
	free = s.Size - used
	if len64 > free {
		return 0, fmt.Errorf("shm full (post-pad)")
	}

	// 5. Commit
	offset := s.tail
	s.allocations = append(s.allocations, AllocationMeta{
		Offset: offset,
		Size:   len64,
		Freed:  false,
	})

	s.tail += len64

	return offset % s.Size, nil
}

func (s *SharedMemory) Free(physOffset int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the active allocation with this physical start offset.
	// Since overlapping active allocations are impossible, physOffset is unique among active items.
	for i := range s.allocations {
		if !s.allocations[i].Freed {
			if s.allocations[i].Offset%s.Size == physOffset {
				s.allocations[i].Freed = true
				break
			}
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
