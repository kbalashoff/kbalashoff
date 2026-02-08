//go:build darwin

package input

import "syscall"

const (
	ioctlReadTermios  = uintptr(syscall.TIOCGETA)
	ioctlWriteTermios = uintptr(syscall.TIOCSETA)
)
