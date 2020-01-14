// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package termenv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/goterm/term"
)

var (
	nbStdout   *os.File
	initStdout sync.Once

	ErrStatusReport = errors.New("unable to retrieve status report")
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

	f := stdoutWithTimeout(100 * time.Millisecond)
	if f == nil {
		return "", ErrStatusReport
	}

	ch := make(chan string)
	fmt.Printf("\033]%d;?\007", sequence)

	go func() {
		r := bufio.NewReader(f)
		out, err := r.ReadBytes('\a')
		if err != nil {
			// timeout
			close(ch)
			return
		}
		ch <- string(out)
	}()

	s, ok := <-ch
	if err = syscall.SetNonblock(int(os.Stdout.Fd()), false); err != nil {
		return "", err
	}

	if !ok {
		return "", ErrStatusReport
	}
	return s, nil
}

func stdoutWithTimeout(d time.Duration) *os.File {
	var err error
	initStdout.Do(func() {
		if err = syscall.SetNonblock(int(os.Stdout.Fd()), true); err != nil {
			return
		}

		nbStdout = os.NewFile(os.Stdout.Fd(), "stdout")
	})
	if err != nil {
		return nil
	}

	if err := nbStdout.SetReadDeadline(time.Now().Add(d)); err != nil {
		return nil
	}
	return nbStdout
}
