//go:build (darwin || dragonfly || freebsd || netbsd || openbsd) && !linux && !solaris && !illumos
// +build darwin dragonfly freebsd netbsd openbsd
// +build !linux
// +build !solaris
// +build !illumos

package termenv

import (
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const (
	// nolint:revive
	_FIONREAD = 0x4004667f
)

func tcFlush(fd int, selector uintptr) error {
	return unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, int(selector))
}

func waitForData(fd uintptr, timeout time.Duration) error {
	tv := syscall.NsecToTimeval(int64(timeout))

	var fds syscall.FdSet
	fds.Bits[0] = 1 << uint(fd)

	return syscall.Select(int(fd)+1, &fds, nil, nil, &tv)
}
