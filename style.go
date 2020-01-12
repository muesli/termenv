package termenv

import (
	"fmt"
	"strings"
)

const (
	CSI          = "\x1b["
	ResetSeq     = "0m"
	BoldSeq      = "1m"
	FaintSeq     = "2m"
	ItalicSeq    = "3m"
	UnderlineSeq = "4m"
	BlinkSeq     = "5m"
)

// Style is a string that various rendering styles can be applied to.
type Style struct {
	string
	styles []styleFunc
}

type styleFunc func(string) string

// String returns a new Style
func String(s ...string) Style {
	return Style{
		string: strings.Join(s, " "),
	}
}

func (t Style) String() string {
	return t.Styled(t.string)
}

// Styled renders s with all applied styles
func (t Style) Styled(s string) string {
	for i := len(t.styles) - 1; i >= 0; i-- {
		s = t.styles[i](s)
	}
	return fmt.Sprintf("%s%s", s, CSI+ResetSeq)
}

// Foreground sets a foreground color
func (t Style) Foreground(c ColorSequencer) Style {
	t.styles = append(t.styles, func(s string) string {
		return wrapSequence(c.Sequence(false), s)
	})
	return t
}

// Background sets a background color
func (t Style) Background(c ColorSequencer) Style {
	t.styles = append(t.styles, func(s string) string {
		return wrapSequence(c.Sequence(true), s)
	})
	return t
}

// Bold enables bold rendering
func (t Style) Bold() Style {
	t.styles = append(t.styles, func(s string) string {
		return wrapSequence(BoldSeq, s)
	})
	return t
}

// Faint enables faint rendering
func (t Style) Faint() Style {
	t.styles = append(t.styles, func(s string) string {
		return wrapSequence(FaintSeq, s)
	})
	return t
}

// Italic enables italic rendering
func (t Style) Italic() Style {
	t.styles = append(t.styles, func(s string) string {
		return wrapSequence(ItalicSeq, s)
	})
	return t
}

// Underline enables underline rendering
func (t Style) Underline() Style {
	t.styles = append(t.styles, func(s string) string {
		return wrapSequence(UnderlineSeq, s)
	})
	return t
}

// Blink enables blink mode
func (t Style) Blink() Style {
	t.styles = append(t.styles, func(s string) string {
		return wrapSequence(BlinkSeq, s)
	})
	return t
}

func wrapSequence(seq string, s string) string {
	return fmt.Sprintf("%s%s", CSI+seq, s)
}
