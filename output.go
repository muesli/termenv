package termenv

import (
	"io"
	"os"
)

var (
	// output is the default global output.
	output = NewOutputWithProfile(os.Stdout, &osEnviron{}, ANSI)
)

// File represents a file descriptor.
type File interface {
	io.ReadWriter
	Fd() uintptr
}

// Output is a terminal output.
type Output struct {
	Profile
	tty     File
	environ Environ
}

// Environ is an interface for getting environment variables.
type Environ interface {
	Environ() []string
	Getenv(string) string
}

type osEnviron struct{}

func (oe *osEnviron) Environ() []string {
	return os.Environ()
}

func (oe *osEnviron) Getenv(key string) string {
	return os.Getenv(key)
}

// DefaultOutput returns the default global output.
func DefaultOutput() *Output {
	return output
}

// NewOutput returns a new Output for the given file descriptor.
func NewOutput(tty File, environ Environ) *Output {
	o := NewOutputWithProfile(tty, environ, Ascii)
	if o.isTTY() {
		o.Profile = o.EnvColorProfile()
	}
	return o
}

// NewOutputWithProfile returns a new Output for the given file descriptor and
// profile.
func NewOutputWithProfile(tty File, environ Environ, profile Profile) *Output {
	return &Output{
		Profile: profile,
		tty:     tty,
		environ: environ,
	}
}

// ForegroundColor returns the terminal's default foreground color.
func (o Output) ForegroundColor() Color {
	if !o.isTTY() {
		return NoColor{}
	}

	return o.foregroundColor()
}

// BackgroundColor returns the terminal's default background color.
func (o Output) BackgroundColor() Color {
	if !o.isTTY() {
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
