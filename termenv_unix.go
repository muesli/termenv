//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package termenv

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

const (
	// timeout for OSC queries.
	OSCTimeout = 5 * time.Second
)

// ColorProfile returns the supported color profile:
// Ascii, ANSI, ANSI256, or TrueColor.
func (o *Output) ColorProfile() Profile {
	if !o.isTTY() {
		return Ascii
	}

	if o.environ.Getenv("GOOGLE_CLOUD_SHELL") == "true" {
		return TrueColor
	}

	term := o.environ.Getenv("TERM")
	colorTerm := o.environ.Getenv("COLORTERM")

	switch strings.ToLower(colorTerm) {
	case "24bit":
		fallthrough
	case "truecolor":
		if strings.HasPrefix(term, "screen") {
			// tmux supports TrueColor, screen only ANSI256
			if o.environ.Getenv("TERM_PROGRAM") != "tmux" {
				return ANSI256
			}
		}
		return TrueColor
	case "yes":
		fallthrough
	case "true":
		return ANSI256
	}

	switch term {
	case "xterm-kitty", "wezterm":
		return TrueColor
	case "linux":
		return ANSI
	}

	if strings.Contains(term, "256color") {
		return ANSI256
	}
	if strings.Contains(term, "color") {
		return ANSI
	}
	if strings.Contains(term, "ansi") {
		return ANSI
	}

	return Ascii
}

func (o Output) foregroundColor() Color {
	s, err := o.termStatusReport(10)
	if err == nil {
		c, err := xTermColor(s)
		if err == nil {
			return c
		}
	}

	colorFGBG := o.environ.Getenv("COLORFGBG")
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

func (o Output) backgroundColor() Color {
	s, err := o.termStatusReport(11)
	if err == nil {
		c, err := xTermColor(s)
		if err == nil {
			return c
		}
	}

	colorFGBG := o.environ.Getenv("COLORFGBG")
	if strings.Contains(colorFGBG, ";") {
		c := strings.Split(colorFGBG, ";")
		i, err := strconv.Atoi(c[len(c)-1])
		if err == nil {
			return ANSIColor(i)
		}
	}

	// default black
	return ANSIColor(0)
}

func (o Output) kittyKeyboardProtocolSupport() byte {
	// screen/tmux can't support OSC, because they can be connected to multiple
	// terminals concurrently.
	term := o.environ.Getenv("TERM")
	if strings.HasPrefix(term, "screen") || strings.HasPrefix(term, "tmux") {
		return 0
	}

	tty := o.TTY()
	if tty == nil {
		return 0
	}

	if !o.unsafe {
		fd := int(tty.Fd())
		// if in background, we can't control the terminal
		if !isForeground(fd) {
			return 0
		}

		t, err := unix.IoctlGetTermios(fd, tcgetattr)
		if err != nil {
			return 0
		}
		defer unix.IoctlSetTermios(fd, tcsetattr, t) //nolint:errcheck

		noecho := *t
		noecho.Lflag = noecho.Lflag &^ unix.ECHO
		noecho.Lflag = noecho.Lflag &^ unix.ICANON
		if err := unix.IoctlSetTermios(fd, tcsetattr, &noecho); err != nil {
			return 0
		}
	}

	// first, send CSI query to see whether this terminal supports the
	// kitty keyboard protocol
	fmt.Fprintf(tty, CSI+"?u")

	// then, query primary device data, should be supported by all terminals
	// if we receive a response for the primary device data befor the kitty keyboard
	// protocol response, this terminal does not support kitty keyboard protocol.
	fmt.Fprintf(tty, CSI+"c")

	response, isAttrs, err := o.readNextResponseKittyKeyboardProtocol()

	// we queried for the kitty keyboard protocol current progressive enhancements
	// but received the primary device attributes response, therefore this terminal
	// does not support the kitty keyboard protocol.
	if err != nil || isAttrs {
		return 0
	}

	// read the primary attrs response and ignore it.
	_, _, err = o.readNextResponseKittyKeyboardProtocol()
	if err != nil {
		return 0
	}

	// we receive a valid response to the kitty keyboard protocol query, this
	// terminal supports the protocol.
	//
	// parse the response and return the flags supported.
	//
	//   0    1 2 3 4
	//   \x1b [ ? 1 u
	//
	if len(response) <= 3 {
		return 0
	}

	return response[3]
}

func (o *Output) waitForData(timeout time.Duration) error {
	fd := o.TTY().Fd()
	tv := unix.NsecToTimeval(int64(timeout))
	var readfds unix.FdSet
	readfds.Set(int(fd))

	for {
		n, err := unix.Select(int(fd)+1, &readfds, nil, nil, &tv)
		if err == unix.EINTR {
			continue
		}
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("timeout")
		}

		break
	}

	return nil
}

func (o *Output) readNextByte() (byte, error) {
	if !o.unsafe {
		if err := o.waitForData(OSCTimeout); err != nil {
			return 0, err
		}
	}

	var b [1]byte
	n, err := o.TTY().Read(b[:])
	if err != nil {
		return 0, err
	}

	if n == 0 {
		panic("read returned no data")
	}

	return b[0], nil
}

// readNextResponseKittyKeyboardProtocol reads either a CSI response to the current
// progressive enhancement status or primary device attributes response.
//   - CSI response: "\x1b]?31u"
//   - primary device attributes response: "\x1b]?64;1;2;7;8;9;15;18;21;44;45;46c"
func (o *Output) readNextResponseKittyKeyboardProtocol() (response string, isAttrs bool, err error) {
	start, err := o.readNextByte()
	if err != nil {
		return "", false, ErrStatusReport
	}

	// first byte must be ESC
	for start != ESC {
		start, err = o.readNextByte()
		if err != nil {
			return "", false, ErrStatusReport
		}
	}

	response += string(start)

	// next byte is [
	tpe, err := o.readNextByte()
	if err != nil {
		return "", false, ErrStatusReport
	}
	response += string(tpe)

	if tpe != '[' {
		return "", false, ErrStatusReport
	}

	for {
		b, err := o.readNextByte()
		if err != nil {
			return "", false, ErrStatusReport
		}
		response += string(b)

		switch b {
		case 'u':
			// kitty keyboard protocol response
			return response, false, nil
		case 'c':
			// primary device attributes response
			return response, true, nil
		}

		// both responses have less than 38 bytes, so if we read more, that's an error
		if len(response) > 38 {
			break
		}
	}

	return response, isAttrs, nil
}

// readNextResponse reads either an OSC response or a cursor position response:
//   - OSC response: "\x1b]11;rgb:1111/1111/1111\x1b\\"
//   - cursor position response: "\x1b[42;1R"
func (o *Output) readNextResponse() (response string, isOSC bool, err error) {
	start, err := o.readNextByte()
	if err != nil {
		return "", false, err
	}

	// first byte must be ESC
	for start != ESC {
		start, err = o.readNextByte()
		if err != nil {
			return "", false, err
		}
	}

	response += string(start)

	// next byte is either '[' (cursor position response) or ']' (OSC response)
	tpe, err := o.readNextByte()
	if err != nil {
		return "", false, err
	}

	response += string(tpe)

	var oscResponse bool
	switch tpe {
	case '[':
		oscResponse = false
	case ']':
		oscResponse = true
	default:
		return "", false, ErrStatusReport
	}

	for {
		b, err := o.readNextByte()
		if err != nil {
			return "", false, err
		}

		response += string(b)

		if oscResponse {
			// OSC can be terminated by BEL (\a) or ST (ESC)
			if b == BEL || strings.HasSuffix(response, string(ESC)) {
				return response, true, nil
			}
		} else {
			// cursor position response is terminated by 'R'
			if b == 'R' {
				return response, false, nil
			}
		}

		// both responses have less than 25 bytes, so if we read more, that's an error
		if len(response) > 25 {
			break
		}
	}

	return "", false, ErrStatusReport
}

func (o Output) termStatusReport(sequence int) (string, error) {
	// screen/tmux can't support OSC, because they can be connected to multiple
	// terminals concurrently.
	term := o.environ.Getenv("TERM")
	if strings.HasPrefix(term, "screen") || strings.HasPrefix(term, "tmux") {
		return "", ErrStatusReport
	}

	tty := o.TTY()
	if tty == nil {
		return "", ErrStatusReport
	}

	if !o.unsafe {
		fd := int(tty.Fd())
		// if in background, we can't control the terminal
		if !isForeground(fd) {
			return "", ErrStatusReport
		}

		t, err := unix.IoctlGetTermios(fd, tcgetattr)
		if err != nil {
			return "", fmt.Errorf("%s: %s", ErrStatusReport, err)
		}
		defer unix.IoctlSetTermios(fd, tcsetattr, t) //nolint:errcheck

		noecho := *t
		noecho.Lflag = noecho.Lflag &^ unix.ECHO
		noecho.Lflag = noecho.Lflag &^ unix.ICANON
		if err := unix.IoctlSetTermios(fd, tcsetattr, &noecho); err != nil {
			return "", fmt.Errorf("%s: %s", ErrStatusReport, err)
		}
	}

	// first, send OSC query, which is ignored by terminal which do not support it
	fmt.Fprintf(tty, OSC+"%d;?"+ST, sequence)

	// then, query cursor position, should be supported by all terminals
	fmt.Fprintf(tty, CSI+"6n")

	// read the next response
	res, isOSC, err := o.readNextResponse()
	if err != nil {
		return "", fmt.Errorf("%s: %s", ErrStatusReport, err)
	}

	// if this is not OSC response, then the terminal does not support it
	if !isOSC {
		return "", ErrStatusReport
	}

	// read the cursor query response next and discard the result
	_, _, err = o.readNextResponse()
	if err != nil {
		return "", err
	}

	// fmt.Println("Rcvd", res[1:])
	return res, nil
}

// EnableVirtualTerminalProcessing enables virtual terminal processing on
// Windows for w and returns a function that restores w to its previous state.
// On non-Windows platforms, or if w does not refer to a terminal, then it
// returns a non-nil no-op function and no error.
func EnableVirtualTerminalProcessing(_ io.Writer) (func() error, error) {
	return func() error { return nil }, nil
}
