package termenv

import "os"

// Output is a terminal output.
type Output struct {
	Profile
	tty *os.File
}

// NewOutput returns a new Output for the given file descriptor.
func NewOutput(tty *os.File) *Output {
	p := Ascii
	if isTTY(tty.Fd()) {
		p = envColorProfile()
	}

	return NewOutputWithProfile(tty, p)
}

// NewOutputWithProfile returns a new Output for the given file descriptor and
// profile.
func NewOutputWithProfile(tty *os.File, profile Profile) *Output {
	return &Output{
		Profile: profile,
		tty:     tty,
	}
}

// NewOutputWithProfileEnv returns a new Output for the given file descriptor and
// autodetects the profile based on provided TERM and COLORTERM variables.
func NewOutputWithProfileEnv(tty *os.File, term, colorTerm string) *Output {
	return &Output{
		Profile: colorProfile(term, colorTerm),
		tty:     tty,
	}
}

// ForegroundColor returns the terminal's default foreground color.
func (o Output) ForegroundColor() Color {
	if !isTTY(o.tty.Fd()) {
		return NoColor{}
	}

	return o.foregroundColor()
}

// BackgroundColor returns the terminal's default background color.
func (o Output) BackgroundColor() Color {
	if !isTTY(o.tty.Fd()) {
		return NoColor{}
	}

	return o.backgroundColor()
}

// HasDarkBackground returns whether terminal uses a dark-ish background.
func (o Output) HasDarkBackground() bool {
	c := ConvertToRGB(o.BackgroundColor())
	_, _, l := c.Hsl()
	return l < 0.5
}
