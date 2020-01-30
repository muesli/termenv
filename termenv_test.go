package termenv

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"
)

func TestTermEnv(t *testing.T) {
	p := ColorProfile()
	if p != TrueColor && p != Ascii {
		t.Errorf("Expected %d, got %d", TrueColor, p)
	}

	fg := ForegroundColor()
	fgexp := "97"
	if fg.Sequence(false) != fgexp && fg.Sequence(false) != "" {
		t.Errorf("Expected %s, got %s", fgexp, fg.Sequence(false))
	}

	bg := BackgroundColor()
	bgexp := "48;2;0;0;0"
	if bg.Sequence(true) != bgexp && bg.Sequence(true) != "" {
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

	exp = "foobar"
	mono := String(exp)
	mono = mono.Foreground(Ascii.Color("#abcdef"))
	if mono.String() != exp {
		t.Errorf("Ascii profile should not apply color styles")
	}
}

func TestColorConversion(t *testing.T) {
	a := ANSI.Color("7")
	c := ConvertToRGB(a)

	exp := "#c0c0c0"
	if c.Hex() != exp {
		t.Errorf("Expected %s, got %s", exp, c.Hex())
	}

	a256 := ANSI256.Color("91")
	c = ConvertToRGB(a256)

	exp = "#8700af"
	if c.Hex() != exp {
		t.Errorf("Expected %s, got %s", exp, c.Hex())
	}

	hex := "#abcdef"
	argb := TrueColor.Color(hex)
	c = ConvertToRGB(argb)

	if c.Hex() != hex {
		t.Errorf("Expected %s, got %s", exp, c.Hex())
	}
}

func TestAscii(t *testing.T) {
	c := Ascii.Color("#abcdef")
	if c.Sequence(false) != "" {
		t.Errorf("Expected empty sequence, got %s", c.Sequence(false))
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
	exp = "90"
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

func TestTemplateHelpers(t *testing.T) {
	p := TrueColor

	exp := String("Hello World")
	basetpl := `{{ %s "Hello World" }}`
	wraptpl := `{{ %s (%s "Hello World") }}`

	tt := []struct {
		Template string
		Expected string
	}{
		{
			Template: fmt.Sprintf(basetpl, "Bold"),
			Expected: exp.Bold().String(),
		},
		{
			Template: fmt.Sprintf(basetpl, "Faint"),
			Expected: exp.Faint().String(),
		},
		{
			Template: fmt.Sprintf(basetpl, "Italic"),
			Expected: exp.Italic().String(),
		},
		{
			Template: fmt.Sprintf(basetpl, "Underline"),
			Expected: exp.Underline().String(),
		},
		{
			Template: fmt.Sprintf(basetpl, "Overline"),
			Expected: exp.Overline().String(),
		},
		{
			Template: fmt.Sprintf(basetpl, "Blink"),
			Expected: exp.Blink().String(),
		},
		{
			Template: fmt.Sprintf(basetpl, "Reverse"),
			Expected: exp.Reverse().String(),
		},
		{
			Template: fmt.Sprintf(basetpl, "CrossOut"),
			Expected: exp.CrossOut().String(),
		},
		{
			Template: fmt.Sprintf(wraptpl, "Underline", "Bold"),
			Expected: String(exp.Bold().String()).Underline().String(),
		},
		{
			Template: `{{ Color "#ff0000" "foobar" }}`,
			Expected: String("foobar").Foreground(p.Color("#ff0000")).String(),
		},
		{
			Template: `{{ Color "#ff0000" "#0000ff" "foobar" }}`,
			Expected: String("foobar").
				Foreground(p.Color("#ff0000")).
				Background(p.Color("#0000ff")).
				String(),
		},
		{
			Template: `{{ Foreground "#ff0000" "foobar" }}`,
			Expected: String("foobar").Foreground(p.Color("#ff0000")).String(),
		},
		{
			Template: `{{ Background "#ff0000" "foobar" }}`,
			Expected: String("foobar").Background(p.Color("#ff0000")).String(),
		},
	}

	for i, v := range tt {
		tpl, err := template.New(fmt.Sprintf("test_%d", i)).Funcs(TemplateFuncs(p)).Parse(v.Template)
		if err != nil {
			t.Error(err)
		}

		var buf bytes.Buffer
		err = tpl.Execute(&buf, nil)
		if err != nil {
			t.Error(err)
		}

		if buf.String() != v.Expected {
			v1 := strings.Replace(v.Expected, "\x1b", "", -1)
			v2 := strings.Replace(buf.String(), "\x1b", "", -1)
			t.Errorf("Expected %s, got %s", v1, v2)
		}
	}
}
