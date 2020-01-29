// +build linux

package termenv

import (
	"syscall"
)

func sysSelect(nfd int, r *syscall.FdSet, timeout *syscall.Timeval) error {
	_, err := syscall.Select(nfd, r, nil, nil, timeout)
	return err
}
