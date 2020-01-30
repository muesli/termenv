package termenv

import "errors"

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

// HasDarkBackground returns whether terminal uses a dark-ish background
func HasDarkBackground() bool {
	c := ConvertToRGB(BackgroundColor())
	_, _, l := c.Hsl()
	return l < 0.5
}
