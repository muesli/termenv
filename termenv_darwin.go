// +build darwin

package termenv

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/google/goterm/term"
)

type Termios struct {
	Iflag  uintptr
	Oflag  uintptr
	Cflag  uintptr
	Lflag  uintptr
	Cc     [20]byte
	Ispeed uintptr
	Ospeed uintptr
}

func Attr(file *os.File) (Termios, error) {
	var t Termios
	fd := file.Fd()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&t)))
	if errno != 0 {
		return t, errno
	}
	return t, nil
}

func Set(t Termios, file *os.File) error {
	fd := file.Fd()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&t)))
	if errno != 0 {
		return errno
	}
	return nil
}

func termStatusReport(sequence int) (string, error) {
	t, err := Attr(os.Stdout)
	if err != nil {
		return "", ErrStatusReport
	}
	defer Set(t, os.Stdout)

	noecho := t
	noecho.Lflag &= term.ECHO &^ term.ICANON
	if err := Set(noecho, os.Stdout); err != nil {
		return "", ErrStatusReport
	}
	fmt.Printf("\033]%d;?\007", sequence)
	s, ok := readWithTimeout(os.Stdout)
	if !ok {
		return "", ErrStatusReport
	}
	// fmt.Println("Rcvd", s[1:])
	return s, nil
}
