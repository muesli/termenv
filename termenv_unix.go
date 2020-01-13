// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package termenv

import (
	"bufio"
	"errors"
	"fmt"
	"log"
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
)

func stdoutWithTimeout(d time.Duration) *os.File {
	initStdout.Do(func() {
		if err := syscall.SetNonblock(1, true); err != nil {
			panic(err)
		}

		nbStdout = os.NewFile(1, "stdout")
	})

	if err := nbStdout.SetReadDeadline(time.Now().Add(d)); err != nil {
		panic(err)
	}
	return nbStdout
}

func SupportedColorProfile() ColorProfile {
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

func DefaultForegroundColor() ColorSequencer {
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

func DefaultBackgroundColor() ColorSequencer {
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
		log.Fatal(err)
	}
	defer t.Set(os.Stdout)

	noecho := t
	noecho.Lflag = noecho.Lflag &^ term.ECHO
	noecho.Lflag = noecho.Lflag &^ term.ICANON
	if err := noecho.Set(os.Stdout); err != nil {
		log.Fatal(err)
	}

	f := stdoutWithTimeout(100 * time.Millisecond)
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
	if !ok {
		return "", errors.New("unable to retrieve status report")
	}
	return s, nil
}
