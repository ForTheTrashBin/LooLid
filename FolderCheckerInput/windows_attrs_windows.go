//go:build windows

package checker

import "syscall"

func hasHiddenOrSystemOnWindows(path string) bool {

	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}

	attr, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return false
	}

	const (
		HIDDEN = 0x2
		SYSTEM = 0x4
	)

	return ((attr&HIDDEN != 0) || (attr&SYSTEM != 0))
}
