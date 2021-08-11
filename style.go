package termenv

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

const (
	ResetSeq     = "0"
	BoldSeq      = "1"
	FaintSeq     = "2"
	ItalicSeq    = "3"
	UnderlineSeq = "4"
	BlinkSeq     = "5"
	ReverseSeq   = "7"
	CrossOutSeq  = "9"
	OverlineSeq  = "53"
)

// Style is a string that various rendering styles can be applied to.
type Style struct {
	string
	fgColor  Color
	bgColor  Color
	modifier []modifier
}

// String returns a new Style.
func String(s ...string) Style {
	return Style{
		string: strings.Join(s, " "),
	}
}

func (t Style) String() string {
	return t.Styled(t.string)
}

// Styled renders s with all applied styles.
func (t Style) Styled(s string) string {
	if len(t.modifier) == 0 && isNoColor(t.bgColor) && isNoColor(t.fgColor) {
		return s
	}

	var builder strings.Builder

	builder.WriteString(CSI)

	for _, mod := range t.modifier {
		builder.WriteString(mod(&t))
	}

	if !isNoColor(t.fgColor) {
		builder.WriteString(t.fgColor.Sequence(false) + ";")
	}

	if !isNoColor(t.bgColor) {
		builder.WriteString(t.bgColor.Sequence(true) + ";")
	}

	builder.WriteString("m" + s + CSI + ResetSeq + "m")

	return builder.String()
}

// Foreground sets a foreground color.
func (t Style) Foreground(c Color) Style {
	t.fgColor = c

	return t
}

// Background sets a background color.
func (t Style) Background(c Color) Style {
	t.bgColor = c

	return t
}

// Bold enables bold rendering.
func (t Style) Bold() Style {
	t.modifier = append(t.modifier, newModifier(BoldSeq))
	return t
}

// Faint enables faint rendering.
func (t Style) Faint() Style {
	t.modifier = append(t.modifier, newModifier(FaintSeq))
	return t
}

func (t Style) ForcedFaint() Style {
	t.modifier = append(t.modifier, func(s *Style) string {
		bgColor := s.backgroundColor()
		fgColor := s.foregroundColor()

		*s = s.Foreground(blend(fgColor, bgColor))

		return ""
	})

	return t
}

func (t Style) AdaptiveFaint() Style {
	t.modifier = append(t.modifier, func(s *Style) string {
		bgColor := s.backgroundColor()
		fgColor := s.foregroundColor()

		if isDarker(bgColor, fgColor) || (bgColor == NoColor{}) || (fgColor == NoColor{}) {
			return FaintSeq + ";"
		}

		*s = s.Foreground(blend(fgColor, bgColor))

		return ""
	})

	return t
}

func (t *Style) foregroundColor() Color {
	if t.fgColor != nil {
		return t.fgColor
	}

	return ForegroundColor()
}

func (t *Style) backgroundColor() Color {
	if t.bgColor != nil {
		return t.bgColor
	}

	return BackgroundColor()
}

func blend(c1 Color, c2 Color) Color {
	profile := colorProfile()

	if (c1 == NoColor{}) || (c2 == NoColor{}) {
		if profile != Ascii {
			return ANSIColor(8)
		}

		return NoColor{}
	}

	c1Rgb := ConvertToRGB(c1)
	c2Rgb := ConvertToRGB(c2)

	return profile.FromColor(c1Rgb.BlendRgb(c2Rgb, 0.5))
}

// Italic enables italic rendering.
func (t Style) Italic() Style {
	t.modifier = append(t.modifier, newModifier(ItalicSeq))
	return t
}

// Underline enables underline rendering.
func (t Style) Underline() Style {
	t.modifier = append(t.modifier, newModifier(UnderlineSeq))
	return t
}

// Overline enables overline rendering.
func (t Style) Overline() Style {
	t.modifier = append(t.modifier, newModifier(OverlineSeq))
	return t
}

// Blink enables blink mode.
func (t Style) Blink() Style {
	t.modifier = append(t.modifier, newModifier(BlinkSeq))
	return t
}

// Reverse enables reverse color mode.
func (t Style) Reverse() Style {
	t.modifier = append(t.modifier, newModifier(ReverseSeq))
	return t
}

// CrossOut enables crossed-out rendering.
func (t Style) CrossOut() Style {
	t.modifier = append(t.modifier, newModifier(CrossOutSeq))
	return t
}

// Width returns the width required to print all runes in Style.
func (t Style) Width() int {
	return runewidth.StringWidth(t.string)
}

type modifier func(*Style) string

func newModifier(sequence string) modifier {
	return func(*Style) string {
		if sequence == "" {
			return sequence
		}

		return sequence + ";"
	}
}

func isNoColor(c Color) bool {
	return c == nil || c == NoColor{}
}
