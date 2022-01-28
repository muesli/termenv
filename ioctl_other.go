//go:build solaris || illumos
// +build solaris illumos

package termenv

import (
	"golang.org/x/sys/unix"
)

func tcFlush(fd int, selector uintptr) error {
	return unix.IoctlSetInt(fd, unix.TCFLSH, int(selector))
}

func waitForData(fd uintptr) error {
	return nil
}
