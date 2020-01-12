package termenv

type ColorProfile int

const (
	Monochrome = ColorProfile(iota)
	ANSI
	ANSI256
	TrueColor
)
