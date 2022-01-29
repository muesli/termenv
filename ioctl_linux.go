//go:build linux
// +build linux

package termenv

import (
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

func tcFlush(fd int, selector uintptr) error {
	return unix.IoctlSetInt(fd, unix.TCFLSH, int(selector))
}

func waitForData(fd uintptr, timeout time.Duration) error {
	tv := syscall.NsecToTimeval(int64(timeout))

	var fds syscall.FdSet
	fds.Bits[0] = 1 << uint(fd)

	_, err := syscall.Select(int(fd)+1, &fds, nil, nil, &tv)
	return err
}
