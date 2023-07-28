package termenv

import (
	"io"
	"os"
	"sync"
)

var (
	// output is the default global output.
	output = NewOutput(os.Stdout)
)

// File represents a file descriptor.
type File interface {
	io.ReadWriter
	Fd() uintptr
}

// OutputOption sets an option on Output.
type OutputOption = func(*Output)

// Output is a terminal output.
type Output struct {
	profile Profile
	tty     io.Writer
	environ Environ

	assumeTTY bool
	unsafe    bool
	cache     bool
	fgSync    *sync.Once
	fgColor   Color
	bgSync    *sync.Once
	bgColor   Color

	mtx sync.RWMutex
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

// SetDefaultOutput sets the default global output.
func SetDefaultOutput(o *Output) {
	output = o
}

// NewOutput returns a new Output for the given file descriptor.
func NewOutput(tty io.Writer, opts ...OutputOption) *Output {
	o := &Output{
		tty:     tty,
		environ: &osEnviron{},
		profile: -1,
		fgSync:  &sync.Once{},
		fgColor: NoColor{},
		bgSync:  &sync.Once{},
		bgColor: NoColor{},
	}

	if o.tty == nil {
		o.tty = os.Stdout
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.profile < 0 {
		o.profile = o.EnvColorProfile()
	}

	return o
}

// WithEnvironment returns a new OutputOption for the given environment.
func WithEnvironment(environ Environ) OutputOption {
	return func(o *Output) {
		o.environ = environ
	}
}

// WithProfile returns a new OutputOption for the given profile.
func WithProfile(profile Profile) OutputOption {
	return func(o *Output) {
		o.profile = profile
	}
}

// WithColorCache returns a new OutputOption with fore- and background color values
// pre-fetched and cached.
func WithColorCache(v bool) OutputOption {
	return func(o *Output) {
		o.cache = v

		// cache the values now
		_ = o.ForegroundColor()
		_ = o.BackgroundColor()
	}
}

// WithTTY returns a new OutputOption to assume whether or not the output is a TTY.
// This is useful when mocking console output.
func WithTTY(v bool) OutputOption {
	return func(o *Output) {
		o.assumeTTY = v
	}
}

// WithUnsafe returns a new OutputOption with unsafe mode enabled. Unsafe mode doesn't
// check whether or not the terminal is a TTY.
//
// This option supersedes WithTTY.
//
// This is useful when mocking console output and enforcing ANSI escape output
// e.g. on SSH sessions.
func WithUnsafe() OutputOption {
	return func(o *Output) {
		o.unsafe = true
	}
}

// ColorProfile returns the supported color profile:
func (o *Output) ColorProfile() Profile {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.profile
}

// SetColorProfile sets the color profile.
func (o *Output) SetColorProfile(profile Profile) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.profile = profile
}

// ForegroundColor returns the terminal's default foreground color.
func (o *Output) ForegroundColor() Color {
	f := func() {
		if !o.isTTY() {
			return
		}

		o.mtx.Lock()
		o.fgColor = o.foregroundColor()
		o.mtx.Unlock()
	}

	if o.cache {
		o.fgSync.Do(f)
	} else {
		f()
	}

	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.fgColor
}

// BackgroundColor returns the terminal's default background color.
func (o *Output) BackgroundColor() Color {
	f := func() {
		if !o.isTTY() {
			return
		}

		o.mtx.Lock()
		o.bgColor = o.backgroundColor()
		o.mtx.Unlock()
	}

	if o.cache {
		o.bgSync.Do(f)
	} else {
		f()
	}

	o.mtx.RLock()
	defer o.mtx.RUnlock()
	return o.bgColor
}

// HasDarkBackground returns whether terminal uses a dark-ish background.
func (o *Output) HasDarkBackground() bool {
	c := ConvertToRGB(o.BackgroundColor())
	_, _, l := c.Hsl()
	return l < 0.5
}

// TTY returns the terminal's file descriptor. This may be nil if the output is
// not a terminal.
func (o *Output) TTY() File {
	if f, ok := o.tty.(File); ok {
		return f
	}
	return nil
}

func (o *Output) Write(p []byte) (int, error) {
	return o.tty.Write(p)
}

// WriteString writes the given string to the output.
func (o *Output) WriteString(s string) (int, error) {
	return o.Write([]byte(s))
}
