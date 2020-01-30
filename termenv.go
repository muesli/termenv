package termenv

import (
	"errors"
	"os"

	"github.com/mattn/go-isatty"
)

var (
	ErrStatusReport = errors.New("unable to retrieve status report")
)

type Profile int

const (
	Ascii = Profile(iota)
	ANSI
	ANSI256
	TrueColor
)

// ColorProfile returns the supported color profile:
// Ascii, ANSI, ANSI256, or TrueColor
func ColorProfile() Profile {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return Ascii
	}

	return colorProfile()
}

// ForegroundColor returns the terminal's default foreground color
func ForegroundColor() Color {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return NoColor{}
	}

	return foregroundColor()
}

// BackgroundColor returns the terminal's default background color
func BackgroundColor() Color {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return NoColor{}
	}

	return backgroundColor()
}

// HasDarkBackground returns whether terminal uses a dark-ish background
func HasDarkBackground() bool {
	c := ConvertToRGB(BackgroundColor())
	_, _, l := c.Hsl()
	return l < 0.5
}
