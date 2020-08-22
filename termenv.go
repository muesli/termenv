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
	CSI = "\x1b["

	Ascii = Profile(iota)
	ANSI
	ANSI256
	TrueColor
)

// ColorProfile returns the supported color profile:
// Ascii, ANSI, ANSI256, or TrueColor.
func ColorProfile() Profile {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return Ascii
	}

	return colorProfile()
}

// ForegroundColor returns the terminal's default foreground color.
func ForegroundColor() Color {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return NoColor{}
	}

	return foregroundColor()
}

// BackgroundColor returns the terminal's default background color.
func BackgroundColor() Color {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return NoColor{}
	}

	return backgroundColor()
}

// HasDarkBackground returns whether terminal uses a dark-ish background.
func HasDarkBackground() bool {
	c := ConvertToRGB(BackgroundColor())
	_, _, l := c.Hsl()
	return l < 0.5
}

// EnvNoColor returns true if the environment variables explicitly disable color output
// by setting NO_COLOR (https://no-color.org/)
// or CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If NO_COLOR is set, this will return true, ignoring CLICOLOR/CLICOLOR_FORCE
// If CLICOLOR=="0", it will be true only if CLICOLOR_FORCE is also "0" or is unset.
func EnvNoColor() bool {
	return os.Getenv("NO_COLOR") != "" || (os.Getenv("CLICOLOR") == "0" && !cliColorForced())
}

// EnvColorProfile returns the color profile based on environment variables set
// Supports NO_COLOR (https://no-color.org/)
// and CLICOLOR/CLICOLOR_FORCE (https://bixense.com/clicolors/)
// If none of these environment variables are set, this behaves the same as ColorProfile()
// It will return the Ascii color profile if EnvNoColor() returns true
// If the terminal does not support any colors, but CLICOLOR_FORCE is set and not "0"
// then the ANSI color profile will be returned.
func EnvColorProfile() Profile {
	if EnvNoColor() {
		return Ascii
	}
	p := ColorProfile()
	if cliColorForced() && p == Ascii {
		return ANSI
	}
	return p
}

func cliColorForced() bool {
	if forced := os.Getenv("CLICOLOR_FORCE"); forced != "" {
		return forced != "0"
	}
	return false
}
