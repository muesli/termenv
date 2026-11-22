//go:build haiku

package termenv

import "os"

func isTerminal(fd uintptr) bool {
	fi, err := os.NewFile(fd, "").Stat()
	return err == nil && (fi.Mode()&os.ModeCharDevice) != 0
}
