package termenv

import "golang.org/x/sys/unix"

const (
	tcgetattr = unix.TCGETA
	tcsetattr = unix.TCSETA
)
