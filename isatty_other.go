//go:build !haiku

package termenv

import "github.com/mattn/go-isatty"

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd)
}
