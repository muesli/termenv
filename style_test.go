package termenv

import (
	"testing"
)

func TestStyleWidth(t *testing.T) {
	s := String("Hello World")
	if s.Width() != 11 {
		t.Errorf("Expected width of 11, got %d", s.Width())
	}

	s = s.Bold()
	if s.Width() != 11 {
		t.Errorf("Expected width of 11, got %d", s.Width())
	}

	s = s.Italic()
	if s.Width() != 11 {
		t.Errorf("Expected width of 11, got %d", s.Width())
	}

	s = s.Foreground(TrueColor.Color("#abcdef"))
	s = s.Background(TrueColor.Color("69"))
	if s.Width() != 11 {
		t.Errorf("Expected width of 11, got %d", s.Width())
	}
}

func TestForceFaint(t *testing.T) {
	s := String("Hello World").Foreground(TrueColor.Color("#40ff00")).Background(TrueColor.Color("#605e10")).ForceFaint()

	exp := "\x1b[38;2;80;175;8;48;2;96;94;16mHello World\x1b[0m"
	if s.String() != exp {
		t.Errorf("Expected %s (%q), got %s (%q)", exp, exp, s.String(), s.String())
	}
}
