package input

import (
	"io"
	"os"
	"strings"
	"syscall"
	"unicode/utf8"
	"unsafe"
)

type Keyboard struct {
	orig syscall.Termios
	out  chan rune
	err  chan error
	done chan struct{}
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		out:  make(chan rune, 32),
		err:  make(chan error, 1),
		done: make(chan struct{}),
	}
}

func (k *Keyboard) Start() (<-chan rune, <-chan error, error) {
	if err := makeRaw(int(os.Stdin.Fd()), &k.orig); err != nil {
		return nil, nil, err
	}
	go k.readLoop()
	return k.out, k.err, nil
}

func (k *Keyboard) Stop() {
	close(k.done)
	_ = restore(int(os.Stdin.Fd()), &k.orig)
}

func (k *Keyboard) readLoop() {
	buf := make([]byte, 8)
	for {
		select {
		case <-k.done:
			close(k.out)
			return
		default:
		}

		n, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				close(k.out)
				return
			}
			k.err <- err
			continue
		}
		if n == 0 {
			continue
		}
		raw := string(buf[:n])
		r, _ := utf8.DecodeRuneInString(raw)
		if r == utf8.RuneError {
			continue
		}
		lower := []rune(strings.ToLower(string(r)))[0]
		k.out <- lower
	}
}

func makeRaw(fd int, orig *syscall.Termios) error {
	termios, err := getTermios(fd)
	if err != nil {
		return err
	}
	*orig = *termios

	termios.Iflag &^= syscall.ICRNL | syscall.INLCR | syscall.IGNCR | syscall.IXON
	termios.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	termios.Cflag |= syscall.CS8
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0

	return setTermios(fd, termios)
}

func restore(fd int, orig *syscall.Termios) error {
	return setTermios(fd, orig)
}

func getTermios(fd int) (*syscall.Termios, error) {
	termios := &syscall.Termios{}
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	return termios, nil
}

func setTermios(fd int, termios *syscall.Termios) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(termios)), 0, 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
