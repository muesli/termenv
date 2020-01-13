package termenv

type Profile int

const (
	Monochrome = Profile(iota)
	ANSI
	ANSI256
	TrueColor
)
