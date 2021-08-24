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

// modifier is a function that is applied before a style is rendered. This
// provides the opportunity to apply post-processing based on the style's
// properties, even if the rendered style is a modified copy of the style to
// which the modifier was added.
type modifier func(*Style) string

// sequenceModifier creates a simple modifier that applies an ANSI sequence.
func sequenceModifier(sequence string) modifier {
	return func(*Style) string {
		if sequence == "" {
			return sequence
		}

		return sequence + ";"
	}
}

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

	processed := t // apply modifiers on a copy
	for _, mod := range t.modifier {
		builder.WriteString(mod(&processed))
	}

	if !isNoColor(processed.fgColor) {
		builder.WriteString(processed.fgColor.Sequence(false) + ";")
	}

	if !isNoColor(processed.bgColor) {
		builder.WriteString(processed.bgColor.Sequence(true) + ";")
	}

	return strings.TrimSuffix(builder.String(), ";") + "m" + s + CSI + ResetSeq + "m"
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
	t.modifier = append(t.modifier, sequenceModifier(BoldSeq))
	return t
}

// Faint enables faint rendering using the ANSI faint/dim sequence. Not all
// terminals render faint text appropriately, especially when a light color
// scheme is used. See ForcedFaint and AdaptiveFaint for alternative solutions.
func (t Style) Faint() Style {
	t.modifier = append(t.modifier, sequenceModifier(FaintSeq))
	return t
}

// ForceFaint produces a consistent faint effect by changing the foreground
// color to a blend of the foreground and background color of the Style. This
// foreground or background color was set, the terminal's default colors are
// used. If this fails, ForceFaint will produce a grey foreground color as
// fallback.
func (t Style) ForceFaint() Style {
	t.modifier = append(t.modifier, func(s *Style) string {
		bgColor := s.backgroundColor()
		fgColor := s.foregroundColor()

		*s = s.Foreground(blend(fgColor, bgColor))

		return ""
	})

	return t
}

// AdaptiveFaint produces a faint effect using Faint for terminals with a dark
// color scheme and ForceFaint for terminals with a light color scheme (see
// HasDarkColorScheme). If the color scheme cannot be detected, it uses Faint.
// This behaviour remedies the fact that many terminals do produce an
// appropriate faint effect for light color schemes.
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

// foregroundColor returns the foreground color when the style is applied
// meaning the style's foreground color or the terminal's current foreground
// color if the style does not have a foreground color set.
func (t *Style) foregroundColor() Color {
	if t.fgColor != nil {
		return t.fgColor
	}

	return ForegroundColor()
}

// backgroundColor returns the background color when the style is applied
// meaning the style's background color or the terminal's current background
// color if the style does not have a background color set.
func (t *Style) backgroundColor() Color {
	if t.bgColor != nil {
		return t.bgColor
	}

	return BackgroundColor()
}

// blend produces a blend between two colors. If one of the arguments is
// NoColor{} a grey color is returned.
func blend(c1 Color, c2 Color) Color {
	profile := ColorProfile()

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
	t.modifier = append(t.modifier, sequenceModifier(ItalicSeq))
	return t
}

// Underline enables underline rendering.
func (t Style) Underline() Style {
	t.modifier = append(t.modifier, sequenceModifier(UnderlineSeq))
	return t
}

// Overline enables overline rendering.
func (t Style) Overline() Style {
	t.modifier = append(t.modifier, sequenceModifier(OverlineSeq))
	return t
}

// Blink enables blink mode.
func (t Style) Blink() Style {
	t.modifier = append(t.modifier, sequenceModifier(BlinkSeq))
	return t
}

// Reverse enables reverse color mode.
func (t Style) Reverse() Style {
	t.modifier = append(t.modifier, sequenceModifier(ReverseSeq))
	return t
}

// CrossOut enables crossed-out rendering.
func (t Style) CrossOut() Style {
	t.modifier = append(t.modifier, sequenceModifier(CrossOutSeq))
	return t
}

// Width returns the width required to print all runes in Style.
func (t Style) Width() int {
	return runewidth.StringWidth(t.string)
}

func isNoColor(c Color) bool {
	return c == nil || c == NoColor{}
}
