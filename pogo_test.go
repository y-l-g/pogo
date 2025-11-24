package pogo

import (
	"testing"
)

func TestSharedMemory_RingBufferStrategy(t *testing.T) {
	// 1. Initialize SHM with small size for testing wrapping
	// Size = 100 bytes. Head=8. Free=92.
	size := int64(100)
	shm, err := NewSharedMemory(size)
	if err != nil {
		t.Fatalf("NewSharedMemory failed: %v", err)
	}
	defer func() { _ = shm.Close() }()

	if shm.head != 8 || shm.tail != 8 {
		t.Errorf("Expected initial head/tail at 8, got H:%d T:%d", shm.head, shm.tail)
	}

	// 2. Alloc 1: 30 bytes.
	// Tail: 8 -> 38.
	off1, err := shm.Allocate(30)
	if err != nil {
		t.Fatalf("Alloc 1 failed: %v", err)
	}
	if off1 != 8 {
		t.Errorf("Expected off1 at 8, got %d", off1)
	}
	if shm.tail != 38 {
		t.Errorf("Expected tail at 38, got %d", shm.tail)
	}

	// 3. Alloc 2: 50 bytes.
	// Tail: 38 -> 88.
	off2, err := shm.Allocate(50)
	if err != nil {
		t.Fatalf("Alloc 2 failed: %v", err)
	}
	if off2 != 38 {
		t.Errorf("Expected off2 at 38, got %d", off2)
	}

	// 4. Alloc 3: 20 bytes.
	// Space check: Used = 88-8 = 80. Free = 20.
	// Phys check: Tail=88. Size=100. 88+20 = 108 > 100. Wrap needed.
	// Pad needed: 100-88 = 12 bytes.
	// Total needed: 20 + 12 = 32 bytes.
	// Available: 20 bytes. -> Fail.
	_, err = shm.Allocate(20)
	if err == nil {
		t.Fatal("Expected Alloc 3 to fail due to fragmentation/full")
	}

	// 5. Free Alloc 1 (30 bytes).
	// Queue: [Alloc1(Free), Alloc2(Active)]
	// Compress: Head moves past Alloc1.
	// Head: 8 -> 38.
	shm.Free(off1)
	if shm.head != 38 {
		t.Errorf("Expected head at 38 after freeing 1, got %d", shm.head)
	}

	// 6. Alloc 3 Again (20 bytes).
	// Used: 88-38 = 50. Free = 50.
	// Wrap: Tail=88. Pad=12. New Tail=100 (Phys 0).
	// Alloc: 20 bytes at 0. New Tail=120.
	// Total consumed: 12(pad) + 20(data) = 32.
	// Free space was 50. 32 <= 50. OK.
	off3, err := shm.Allocate(20)
	if err != nil {
		t.Fatalf("Alloc 3 failed after free: %v", err)
	}
	if off3 != 0 {
		t.Errorf("Expected off3 at 0 (wrapped), got %d", off3)
	}

	// 7. Verify FIFO Blocking
	// Free Alloc 3 (20 bytes).
	// Queue: [Alloc2(Active), Pad(Free), Alloc3(Free)]
	// Compress: Head cannot move past Alloc2.
	shm.Free(off3)
	if shm.head != 38 {
		t.Errorf("Expected head to stay at 38 (Alloc2 blocking), got %d", shm.head)
	}

	// Free Alloc 2.
	// Queue: [Alloc2(Free), Pad(Free), Alloc3(Free)]
	// Compress:
	// - Alloc2 (50b) -> Head=88
	// - Pad (12b) -> Head=100
	// - Alloc3 (20b) -> Head=120
	shm.Free(off2)
	if shm.head != 120 {
		t.Errorf("Expected head to clear all (120), got %d", shm.head)
	}
}

func TestSharedMemory_BoundsCheck(t *testing.T) {
	shm, _ := NewSharedMemory(128)
	defer func() { _ = shm.Close() }()

	_, err := shm.Allocate(200)
	if err == nil {
		t.Error("Expected error for allocation larger than SHM size")
	}
}
