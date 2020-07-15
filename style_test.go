package termenv

import (
	"testing"
)

func TestStyleLen(t *testing.T) {
	s := String("Hello World")
	if s.Len() != 11 {
		t.Errorf("Expected len of 11, got %d", s.Len())
	}

	s = s.Bold()
	if s.Len() != 11 {
		t.Errorf("Expected len of 11, got %d", s.Len())
	}

	s = s.Italic()
	if s.Len() != 11 {
		t.Errorf("Expected len of 11, got %d", s.Len())
	}

	s = s.Foreground(TrueColor.Color("#abcdef"))
	s = s.Background(TrueColor.Color("69"))
	if s.Len() != 11 {
		t.Errorf("Expected len of 11, got %d", s.Len())
	}
}
