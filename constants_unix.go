// +build darwin dragonfly freebsd netbsd openbsd solaris

package termenv

import "golang.org/x/sys/unix"

const (
	tcgetattr = unix.TIOCGETA
	tcsetattr = unix.TIOCSETA
)
