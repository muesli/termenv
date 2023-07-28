package termenv

import (
	"io"
	"testing"
)

func TestOutputRace(t *testing.T) {
	o := NewOutput(io.Discard)
	for i := 0; i < 100; i++ {
		t.Run("Test race", func(t *testing.T) {
			t.Parallel()
			o.Write([]byte("test"))
			o.SetColorProfile(ANSI)
			o.ColorProfile()
			o.HasDarkBackground()
			o.TTY()
			o.ForegroundColor()
			o.BackgroundColor()
		})
	}
}
