// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package termenv

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/goterm/term"
)

func ColorProfile() Profile {
	colorTerm := os.Getenv("COLORTERM")
	if colorTerm == "truecolor" {
		return TrueColor
	}

	term := os.Getenv("TERM")
	if strings.Contains(term, "256color") {
		return ANSI256
	}
	if strings.Contains(term, "color") {
		return ANSI
	}

	return Monochrome
}

func ForegroundColor() Color {
	s, err := termStatusReport(10)
	if err == nil {
		c, err := xTermColor(s)
		if err == nil {
			return c
		}
	}

	colorFGBG := os.Getenv("COLORFGBG")
	if strings.Contains(colorFGBG, ";") {
		c := strings.Split(colorFGBG, ";")
		i, err := strconv.Atoi(c[0])
		if err == nil {
			return ANSIColor(i)
		}
	}

	// default gray
	return ANSIColor(7)
}

func BackgroundColor() Color {
	s, err := termStatusReport(11)
	if err == nil {
		c, err := xTermColor(s)
		if err == nil {
			return c
		}
	}

	colorFGBG := os.Getenv("COLORFGBG")
	if strings.Contains(colorFGBG, ";") {
		c := strings.Split(colorFGBG, ";")
		i, err := strconv.Atoi(c[1])
		if err == nil {
			return ANSIColor(i)
		}
	}

	// default black
	return ANSIColor(0)
}

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

func readWithTimeout(f *os.File) (string, bool) {
	var readfds syscall.FdSet
	fd := f.Fd()
	readfds.Bits[fd/64] |= 1 << (fd % 64)

	// Use select to attempt to read from os.Stdout for 100 ms
	n, err := syscall.Select(int(fd)+1,
		&readfds, nil, nil,
		&syscall.Timeval{Usec: 100000})

	if err != nil {
		// log.Printf("select(read stdout): %v", err)
		return "", false
	}
	if n == 0 {
		// log.Printf("select(read stdout): timed out")
		return "", false
	}

	// n > 0 => is readable
	data := make([]byte, 24)
	n, err = f.Read(data)
	if err != nil {
		// log.Printf("read(stdout): %v", err)
		return "", false
	}

	// fmt.Printf("read %d bytes from stdout: %s\n", n, data)
	return string(data), true
}
