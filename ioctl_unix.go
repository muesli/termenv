//go:build (darwin || dragonfly || freebsd || netbsd || openbsd) && !linux && !solaris && !illumos
// +build darwin dragonfly freebsd netbsd openbsd
// +build !linux
// +build !solaris
// +build !illumos

package termenv

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	// nolint:revive
	_FIONREAD = 0x4004667f
)

func tcFlush(fd int, selector uintptr) error {
	return unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, int(selector))
}

func waitForData(fd uintptr) error {
	var avail int
	var err syscall.Errno

	for i := 1; i < 10; i++ {
		_, _, err = syscall.Syscall(syscall.SYS_IOCTL, fd, _FIONREAD, uintptr(unsafe.Pointer(&avail)))
		if err != 0 || avail > 0 {
			break
		}

		time.Sleep(time.Duration(i*i) * time.Millisecond)
	}

	if avail == 0 || err != 0 {
		return fmt.Errorf("timeout")
	}

	return nil
}
