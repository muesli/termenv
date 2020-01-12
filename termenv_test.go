package termenv

import (
	"testing"
)

func TestTermEnv(t *testing.T) {
	p := SupportedColorProfile()
	if p != TrueColor {
		t.Errorf("Expected %d, got %d", TrueColor, p)
	}

	fg := DefaultForegroundColor()
	fgexp := "37;1m"
	if fg.Sequence(false) != fgexp {
		t.Errorf("Expected %s, got %s", fgexp, fg.Sequence(false))
	}

	bg := DefaultBackgroundColor()
	bgexp := "48;2;0;0;0m"
	if bg.Sequence(true) != bgexp {
		t.Errorf("Expected %s, got %s", bgexp, bg.Sequence(true))
	}
}

func TestANSIProfile(t *testing.T) {
	p := ANSI

	c := p.Color("#abcdef")
	exp := "37m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSIColor); !ok {
		t.Errorf("Expected type termenv.ANSIColor, got %T", c)
	}

	c = p.Color("139")
	exp = "30;1m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSIColor); !ok {
		t.Errorf("Expected type termenv.ANSIColor, got %T", c)
	}

	c = p.Color("2")
	exp = "32m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSIColor); !ok {
		t.Errorf("Expected type termenv.ANSIColor, got %T", c)
	}
}

func TestANSI256Profile(t *testing.T) {
	p := ANSI256

	c := p.Color("#abcdef")
	exp := "38;5;153m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSI256Color); !ok {
		t.Errorf("Expected type termenv.ANSI256Color, got %T", c)
	}

	c = p.Color("139")
	exp = "38;5;139m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSI256Color); !ok {
		t.Errorf("Expected type termenv.ANSI256Color, got %T", c)
	}

	c = p.Color("2")
	exp = "32m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSIColor); !ok {
		t.Errorf("Expected type termenv.ANSIColor, got %T", c)
	}
}

func TestTrueColorProfile(t *testing.T) {
	p := TrueColor

	c := p.Color("#abcdef")
	exp := "38;2;171;205;239m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(RGBColor); !ok {
		t.Errorf("Expected type termenv.HexColor, got %T", c)
	}

	c = p.Color("139")
	exp = "38;5;139m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSI256Color); !ok {
		t.Errorf("Expected type termenv.ANSI256Color, got %T", c)
	}

	c = p.Color("2")
	exp = "32m"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSIColor); !ok {
		t.Errorf("Expected type termenv.ANSIColor, got %T", c)
	}
}

func TestStyles(t *testing.T) {
	s := String("foobar").Foreground(TrueColor.Color("2"))

	exp := "\x1b[32mfoobar\x1b[0m"
	if s.String() != exp {
		t.Errorf("Expected %s, got %s", exp, s.String())
	}
}
