//go:build linux
// +build linux

package termenv

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

func tcFlush(fd int, selector uintptr) error {
	return unix.IoctlSetInt(fd, unix.TCFLSH, int(selector))
}

func waitForData(fd uintptr) error {
	var avail int
	var err error

	for i := 1; i < 10; i++ {
		avail, err = unix.IoctlGetInt(int(fd), unix.TIOCINQ)
		if err != nil || avail > 0 {
			break
		}

		time.Sleep(time.Duration(i*i) * time.Millisecond)
	}

	if avail == 0 || err != nil {
		return fmt.Errorf("timeout")
	}

	return nil
}
