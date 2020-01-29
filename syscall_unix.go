// +build !linux,!windows,!plan9,!solaris

package termenv

import (
	"syscall"
)

func sysSelect(nfd int, r *syscall.FdSet, timeout *syscall.Timeval) error {
	return syscall.Select(nfd, r, nil, nil, timeout)
}
