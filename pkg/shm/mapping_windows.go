//go:build windows

package shm

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

func mapFile(f *os.File, size int64) ([]byte, error) {
	// 1. Create File Mapping
	h, err := windows.CreateFileMapping(windows.Handle(f.Fd()), nil, windows.PAGE_READWRITE, 0, uint32(size), nil)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(h)

	// 2. Map View
	addr, err := windows.MapViewOfFile(h, windows.FILE_MAP_WRITE, 0, 0, uintptr(size))
	if err != nil {
		return nil, err
	}

	// 3. Convert uintptr to []byte
	// Note: In Go 1.20+ we can use unsafe.Slice
	data := unsafe.Slice((*byte)(unsafe.Pointer(addr)), int(size))
	return data, nil
}

func unmapFile(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	addr := uintptr(unsafe.Pointer(&data[0]))
	return windows.UnmapViewOfFile(addr)
}
