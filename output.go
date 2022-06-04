package termenv

import (
	"io"
	"os"
	"sync"
)

var (
	// output is the default global output.
	output = NewOutput(os.Stdout, WithProfile(ANSI))
)

// File represents a file descriptor.
type File interface {
	io.ReadWriter
	Fd() uintptr
}

// Output is a terminal output.
type Output struct {
	Profile
	tty     io.Writer
	environ Environ

	fgSync  *sync.Once
	fgColor Color
	bgSync  *sync.Once
	bgColor Color
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
func NewOutput(tty io.Writer, opts ...func(*Output)) *Output {
	o := &Output{
		tty:     tty,
		environ: &osEnviron{},
		Profile: Ascii,
		fgSync:  &sync.Once{},
		fgColor: NoColor{},
		bgSync:  &sync.Once{},
		bgColor: NoColor{},
	}
	o.Profile = o.EnvColorProfile()

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithProfile returns a new Output for the given environment. The profile gets
// detected from this environment.
func WithEnvironment(environ Environ) func(*Output) {
	return func(o *Output) {
		o.environ = environ
		o.Profile = o.EnvColorProfile()
	}
}

// WithProfile returns a new Output for the given file descriptor and profile.
func WithProfile(profile Profile) func(*Output) {
	return func(o *Output) {
		o.Profile = profile
	}
}

// ForegroundColor returns the terminal's default foreground color.
func (o Output) ForegroundColor() Color {
	o.fgSync.Do(func() {
		if !o.isTTY() {
			return
		}

		o.fgColor = o.foregroundColor()
	})

	return o.fgColor
}

// BackgroundColor returns the terminal's default background color.
func (o Output) BackgroundColor() Color {
	o.bgSync.Do(func() {
		if !o.isTTY() {
			return
		}

		o.bgColor = o.backgroundColor()
	})

	return o.bgColor
}

// HasDarkBackground returns whether terminal uses a dark-ish background.
func (o Output) HasDarkBackground() bool {
	c := ConvertToRGB(o.BackgroundColor())
	_, _, l := c.Hsl()
	return l < 0.5
}

// TTY returns the terminal's file descriptor. This may be nil if the output is
// not a terminal.
func (o Output) TTY() File {
	if f, ok := o.tty.(File); ok {
		return f
	}
	return nil
}

func (o Output) Write(p []byte) (int, error) {
	return o.tty.Write(p)
}
