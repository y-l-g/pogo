//go:build !windows

package shm

import (
	"os"
	"syscall"
)

func mapFile(f *os.File, size int64) ([]byte, error) {
	return syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
}

func unmapFile(data []byte) error {
	return syscall.Munmap(data)
}
