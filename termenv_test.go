package termenv

import (
	"testing"
)

func TestTermEnv(t *testing.T) {
	p := ColorProfile()
	if p != TrueColor && p != Monochrome {
		t.Errorf("Expected %d, got %d", TrueColor, p)
	}

	fg := ForegroundColor()
	fgexp := "37;1"
	if fg.Sequence(false) != fgexp && fg.Sequence(false) != "37" {
		t.Errorf("Expected %s, got %s", fgexp, fg.Sequence(false))
	}

	bg := BackgroundColor()
	bgexp := "48;2;0;0;0"
	if bg.Sequence(true) != bgexp && bg.Sequence(true) != "40" {
		t.Errorf("Expected %s, got %s", bgexp, bg.Sequence(true))
	}

	_ = HasDarkBackground()
}

func TestRendering(t *testing.T) {
	out := String("foobar")
	if out.String() != "foobar" {
		t.Errorf("Unstyled strings should be returned as plain text")
	}

	out = out.Foreground(TrueColor.Color("#abcdef"))
	out = out.Background(TrueColor.Color("69"))
	out = out.Bold()
	out = out.Italic()
	out = out.Faint()
	out = out.Underline()
	out = out.Blink()

	exp := "\x1b[38;2;171;205;239;48;5;69;1;3;2;4;5mfoobar\x1b[0m"
	if out.String() != exp {
		t.Errorf("Expected %s, got %s", exp, out.String())
	}
}

func TestColorConversion(t *testing.T) {
	a := ANSI.Color("7")
	c := convertToRGB(a)

	exp := "#c0c0c0"
	if c.Hex() != exp {
		t.Errorf("Expected %s, got %s", exp, c.Hex())
	}

	a256 := ANSI256.Color("91")
	c = convertToRGB(a256)

	exp = "#8700af"
	if c.Hex() != exp {
		t.Errorf("Expected %s, got %s", exp, c.Hex())
	}

	hex := "#abcdef"
	argb := TrueColor.Color(hex)
	c = convertToRGB(argb)

	if c.Hex() != hex {
		t.Errorf("Expected %s, got %s", exp, c.Hex())
	}
}

func TestANSIProfile(t *testing.T) {
	p := ANSI

	c := p.Color("#abcdef")
	exp := "37"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSIColor); !ok {
		t.Errorf("Expected type termenv.ANSIColor, got %T", c)
	}

	c = p.Color("139")
	exp = "30;1"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSIColor); !ok {
		t.Errorf("Expected type termenv.ANSIColor, got %T", c)
	}

	c = p.Color("2")
	exp = "32"
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
	exp := "38;5;153"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSI256Color); !ok {
		t.Errorf("Expected type termenv.ANSI256Color, got %T", c)
	}

	c = p.Color("139")
	exp = "38;5;139"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSI256Color); !ok {
		t.Errorf("Expected type termenv.ANSI256Color, got %T", c)
	}

	c = p.Color("2")
	exp = "32"
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
	exp := "38;2;171;205;239"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(RGBColor); !ok {
		t.Errorf("Expected type termenv.HexColor, got %T", c)
	}

	c = p.Color("139")
	exp = "38;5;139"
	if c.Sequence(false) != exp {
		t.Errorf("Expected %s, got %s", exp, c.Sequence(false))
	}
	if _, ok := c.(ANSI256Color); !ok {
		t.Errorf("Expected type termenv.ANSI256Color, got %T", c)
	}

	c = p.Color("2")
	exp = "32"
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
