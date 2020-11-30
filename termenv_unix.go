// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package termenv

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func colorProfile() Profile {
	term := os.Getenv("TERM")
	colorTerm := os.Getenv("COLORTERM")

	switch strings.ToLower(colorTerm) {
	case "24bit":
		fallthrough
	case "truecolor":
		if term == "screen" || !strings.HasPrefix(term, "screen") {
			// enable TrueColor in tmux, but not for old-school screen
			return TrueColor
		}
	case "yes":
		fallthrough
	case "true":
		return ANSI256
	}

	if strings.Contains(term, "256color") {
		return ANSI256
	}
	if strings.Contains(term, "color") {
		return ANSI
	}

	return Ascii
}

func foregroundColor() Color {
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

func backgroundColor() Color {
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

func readWithTimeout(f *os.File) (string, bool) {
	var readfds unix.FdSet
	fd := int(f.Fd())
	readfds.Set(fd)

	for {
		// Use select to attempt to read from os.Stdout for 100 ms
		_, err := unix.Select(fd+1, &readfds, nil, nil, &unix.Timeval{Usec: 100000})
		if err == nil {
			break
		}
		// On MacOS we can see EINTR here if the user
		// pressed ^Z. Similar to issue https://github.com/golang/go/issues/22838
		if runtime.GOOS == "darwin" && err == unix.EINTR {
			continue
		}
		// log.Printf("select(read error): %v", err)
		return "", false
	}

	if !readfds.IsSet(fd) {
		// log.Print("select(read timeout)")
		return "", false
	}

	// n > 0 => is readable
	var data []byte
	b := make([]byte, 1)
	for {
		_, err := f.Read(b)
		if err != nil {
			// log.Printf("read(%d): %v %d", fd, err, n)
			return "", false
		}
		// log.Printf("read %d bytes from stdout: %s %d\n", n, data, data[len(data)-1])

		data = append(data, b[0])

		// data sent by terminal is either terminated by BEL (\a) or ST (ESC \)
		if bytes.HasSuffix(data, []byte("\a")) || bytes.HasSuffix(data, []byte("\033\\")) {
			break
		}
	}

	// fmt.Printf("read %d bytes from stdout: %s\n", n, data)
	return string(data), true
}

func termStatusReport(sequence int) (string, error) {
	term := os.Getenv("TERM")
	if strings.HasPrefix(term, "screen") {
		return "", ErrStatusReport
	}

	t, err := unix.IoctlGetTermios(unix.Stdout, tcgetattr)
	if err != nil {
		return "", ErrStatusReport
	}
	defer unix.IoctlSetTermios(unix.Stdout, tcsetattr, t)

	noecho := *t
	noecho.Lflag = noecho.Lflag &^ unix.ECHO
	noecho.Lflag = noecho.Lflag &^ unix.ICANON
	if err := unix.IoctlSetTermios(unix.Stdout, tcsetattr, &noecho); err != nil {
		return "", ErrStatusReport
	}

	fmt.Printf("\033]%d;?\033\\", sequence)
	s, ok := readWithTimeout(os.Stdout)
	if !ok {
		return "", ErrStatusReport
	}
	// fmt.Println("Rcvd", s[1:])
	return s, nil
}
