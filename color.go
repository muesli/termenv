package termenv

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

var (
	// ErrInvalidColor gets returned when a color is invalid.
	ErrInvalidColor = errors.New("invalid color")
)

const (
	// Foreground sequence code.
	ForegroudSeq = "38"
	// Background sequence code.
	BackgroundSeq = "48"
	// Underline color sequence code.
	UnderColorSeq = "58"
)

// Foreground and Background sequence codes
const (
	// Foreground sequence code.
	//
	// Deprecated: use ForegroudSeq
	Foreground = ForegroudSeq
	// Background sequence code.
	//
	// Deprecated: use BackgroundSeq
	Background = BackgroundSeq
)

// Color is an interface implemented by all colors that can be converted to an
// ANSI sequence.
type Color interface {
	// Sequence returns the ANSI Sequence for the color.
	Sequence() string
}

// NoColor is a nop for terminals that don't support colors.
type NoColor struct{}

func (c NoColor) String() string {
	return ""
}

// ANSIColor is a color (0-15) as defined by the ANSI Standard.
type ANSIColor struct {
	Color uint8
	bg    bool
}

func (c ANSIColor) String() string {
	if c.Color >= 16 {
		return ""
	}
	return ansiHex[c.Color]
}

// ANSI256Color is a color (16-255) as defined by the ANSI Standard.
type ANSI256Color uint8

func (c ANSI256Color) String() string {
	return ansiHex[c]
}

// RGBColor is a hex-encoded color, e.g. "#abcdef".
type RGBColor string

// ConvertToRGB converts a Color to a colorful.Color.
func ConvertToRGB(c Color) colorful.Color {
	var hex string
	switch v := c.(type) {
	case RGBColor:
		hex = string(v)
	case ANSIColor:
		hex = ansiHex[v.Color]
	case ANSI256Color:
		hex = ansiHex[v]
	}

	ch, _ := colorful.Hex(hex)
	return ch
}

// Sequence returns the ANSI Sequence for the color.
func (c NoColor) Sequence() string {
	return ""
}

// Sequence returns the ANSI Sequence for the color.
func (c ANSIColor) Sequence() string {
	col := int(c.Color)
	bgMod := func(col int) int {
		if c.bg {
			return col + 10
		}
		return col
	}

	if col < 8 {
		return fmt.Sprintf("%d", bgMod(col)+30)
	}
	return fmt.Sprintf("%d", bgMod(col-8)+90)
}

// Sequence returns the ANSI Sequence for the color.
func (c ANSI256Color) Sequence() string {
	return fmt.Sprintf("5;%d", c)
}

// Sequence returns the ANSI Sequence for the color.
func (c RGBColor) Sequence() string {
	f, err := colorful.Hex(string(c))
	if err != nil {
		return ""
	}

	return fmt.Sprintf("2;%d;%d;%d", uint8(f.R*255), uint8(f.G*255), uint8(f.B*255))
}

func xTermColor(s string) (RGBColor, error) {
	if len(s) < 24 || len(s) > 25 {
		return RGBColor(""), ErrInvalidColor
	}

	switch {
	case strings.HasSuffix(s, string(BEL)):
		s = strings.TrimSuffix(s, string(BEL))
	case strings.HasSuffix(s, string(ESC)):
		s = strings.TrimSuffix(s, string(ESC))
	case strings.HasSuffix(s, ST):
		s = strings.TrimSuffix(s, ST)
	default:
		return RGBColor(""), ErrInvalidColor
	}

	s = s[4:]

	prefix := ";rgb:"
	if !strings.HasPrefix(s, prefix) {
		return RGBColor(""), ErrInvalidColor
	}
	s = strings.TrimPrefix(s, prefix)

	h := strings.Split(s, "/")
	hex := fmt.Sprintf("#%s%s%s", h[0][:2], h[1][:2], h[2][:2])
	return RGBColor(hex), nil
}

func ansi256ToANSIColor(c ANSI256Color) ANSIColor {
	var r uint8
	md := math.MaxFloat64

	h, _ := colorful.Hex(ansiHex[c])
	for i := 0; i <= 15; i++ {
		hb, _ := colorful.Hex(ansiHex[i])
		d := h.DistanceHSLuv(hb)

		if d < md {
			md = d
			r = uint8(i)
		}
	}

	return ANSIColor{Color: r}
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
	colorDist := c.DistanceHSLuv(c2)
	grayDist := c.DistanceHSLuv(g2)

	if colorDist <= grayDist {
		return ANSI256Color(16 + ci)
	}
	return ANSI256Color(232 + grayIdx)
}
