package termenv

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

var (
	ErrInvalidColor = errors.New("invalid color")
)

const (
	Foreground = "38"
	Background = "48"
)

type Color interface {
	Sequence(bg bool) string
}

type NoColor struct{}

// ANSIColor is a color (0-15) as defined by the ANSI Standard
type ANSIColor int

// ANSI256Color is a color (16-255) as defined by the ANSI Standard
type ANSI256Color int

// RGBColor is a hex-encoded color, e.g. "#abcdef"
type RGBColor string

func ConvertToRGB(c Color) colorful.Color {
	var hex string
	switch v := c.(type) {
	case RGBColor:
		hex = string(v)
	case ANSIColor:
		hex = ansiHex[v]
	case ANSI256Color:
		hex = ansiHex[v]
	}

	ch, _ := colorful.Hex(hex)
	return ch
}

func (p Profile) Convert(c Color) Color {
	if p == Ascii {
		return NoColor{}
	}

	switch v := c.(type) {
	case ANSIColor:
		return v

	case ANSI256Color:
		if p == ANSI {
			return ansi256ToANSIColor(v)
		}
		return v

	case RGBColor:
		h, err := colorful.Hex(string(v))
		if err != nil {
			return nil
		}
		if p < TrueColor {
			ac := hexToANSI256Color(h)
			if p == ANSI {
				return ansi256ToANSIColor(ac)
			}
			return ac
		}
		return v
	}

	return c
}

func (p Profile) Color(s string) Color {
	if len(s) == 0 {
		return nil
	}

	var c Color
	if strings.HasPrefix(s, "#") {
		c = RGBColor(s)
	} else {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}

		if i < 16 {
			c = ANSIColor(i)
		} else {
			c = ANSI256Color(i)
		}
	}

	return p.Convert(c)
}

func (c NoColor) Sequence(bg bool) string {
	return ""
}

func (c ANSIColor) Sequence(bg bool) string {
	col := int(c)
	bgMod := func(c int) int {
		if bg {
			return c + 10
		}
		return c
	}

	if col < 8 {
		return fmt.Sprintf("%d", bgMod(col)+30)
	}
	return fmt.Sprintf("%d", bgMod(col-8)+90)
}

func (c ANSI256Color) Sequence(bg bool) string {
	prefix := Foreground
	if bg {
		prefix = Background
	}
	return fmt.Sprintf("%s;5;%d", prefix, c)
}

func (c RGBColor) Sequence(bg bool) string {
	f, err := colorful.Hex(string(c))
	if err != nil {
		return ""
	}

	prefix := Foreground
	if bg {
		prefix = Background
	}
	return fmt.Sprintf("%s;2;%d;%d;%d", prefix, uint8(f.R*255), uint8(f.G*255), uint8(f.B*255))
}

func xTermColor(s string) (RGBColor, error) {
	if len(s) != 24 {
		return RGBColor(""), ErrInvalidColor
	}
	s = s[4:]

	prefix := ";rgb:"
	if !strings.HasPrefix(s, prefix) {
		return RGBColor(""), ErrInvalidColor
	}
	s = strings.TrimPrefix(s, prefix)
	s = strings.TrimSuffix(s, "\a")

	h := strings.Split(s, "/")
	hex := fmt.Sprintf("#%s%s%s", h[0][:2], h[1][:2], h[2][:2])
	return RGBColor(hex), nil
}

func ansi256ToANSIColor(c ANSI256Color) ANSIColor {
	var r int
	md := math.MaxFloat64

	h, _ := colorful.Hex(ansiHex[c])
	for i := 0; i <= 15; i++ {
		hb, _ := colorful.Hex(ansiHex[i])
		d := h.DistanceLab(hb)

		if d < md {
			md = d
			r = i
		}
	}

	return ANSIColor(r)
}

func hexToANSI256Color(c colorful.Color) ANSI256Color {
	v2ci := func(v float64) int {
		if v < 48 {
			return 0
		}
		if v < 115 {
			return 1
		}
		return int((v - 35) / 40)
	}

	// Calculate the nearest 0-based color index at 16..231
	r := v2ci(c.R * 255.0) // 0..5 each
	g := v2ci(c.G * 255.0)
	b := v2ci(c.B * 255.0)
	ci := 36*r + 6*g + b /* 0..215 */

	// Calculate the represented colors back from the index
	i2cv := [6]int{0, 0x5f, 0x87, 0xaf, 0xd7, 0xff}
	cr := i2cv[r] // r/g/b, 0..255 each
	cg := i2cv[g]
	cb := i2cv[b]

	// Calculate the nearest 0-based gray index at 232..255
	var grayIdx int
	average := (r + g + b) / 3
	if average > 238 {
		grayIdx = 23
	} else {
		grayIdx = (average - 3) / 10 // 0..23
	}
	gv := 8 + 10*grayIdx // same value for r/g/b, 0..255

	// Return the one which is nearer to the original input rgb value
	c2 := colorful.Color{R: float64(cr) / 255.0, G: float64(cg) / 255.0, B: float64(cb) / 255.0}
	g2 := colorful.Color{R: float64(gv) / 255.0, G: float64(gv) / 255.0, B: float64(gv) / 255.0}
	colorDist := c.DistanceLab(c2)
	grayDist := c.DistanceLab(g2)

	if colorDist <= grayDist {
		return ANSI256Color(16 + ci)
	}
	return ANSI256Color(232 + grayIdx)
}
