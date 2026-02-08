//go:build linux

package input

import "syscall"

const (
	ioctlReadTermios  = uintptr(syscall.TCGETS)
	ioctlWriteTermios = uintptr(syscall.TCSETS)
)
