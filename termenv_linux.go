// +build dragonfly freebsd linux netbsd openbsd solaris

package termenv

import (
	"fmt"
	"os"

	"github.com/google/goterm/term"
)

func termStatusReport(sequence int) (string, error) {
	t, err := term.Attr(os.Stdout)
	if err != nil {
		return "", ErrStatusReport
	}
	defer t.Set(os.Stdout)

	noecho := t
	noecho.Lflag = noecho.Lflag &^ term.ECHO
	noecho.Lflag = noecho.Lflag &^ term.ICANON
	if err := noecho.Set(os.Stdout); err != nil {
		return "", ErrStatusReport
	}

	fmt.Printf("\033]%d;?\007", sequence)
	s, ok := readWithTimeout(os.Stdout)
	if !ok {
		return "", ErrStatusReport
	}
	// fmt.Println("Rcvd", s[1:])
	return s, nil
}
