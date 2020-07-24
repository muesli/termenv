// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package termenv

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func colorProfile() Profile {
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
	var readfds syscall.FdSet
	fd := f.Fd()
	readfds.Bits[fd/64] |= 1 << (fd % 64)

	for {
		// Use select to attempt to read from os.Stdout for 100 ms
		err := sysSelect(int(fd)+1,
			&readfds,
			&syscall.Timeval{Usec: 100000})
		if err == nil {
			break
		}
		// On MacOS we can see EINTR here if the user
		// pressed ^Z. Similar to issue https://github.com/golang/go/issues/22838
		if runtime.GOOS == "darwin" && err == syscall.EINTR {
			continue
		}
		// log.Printf("select(read error): %v", err)
		return "", false
	}

	if readfds.Bits[fd/64]&(1<<(fd%64)) == 0 {
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
		if b[0] == '\a' || (b[0] == '\\' && len(data) > 2) {
			break
		}
	}

	// fmt.Printf("read %d bytes from stdout: %s\n", n, data)
	return string(data), true
}
