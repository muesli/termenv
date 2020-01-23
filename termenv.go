package termenv

type Profile int

const (
	Monochrome = Profile(iota)
	ANSI
	ANSI256
	TrueColor
)

func HasDarkBackground() bool {
	c := ConvertToRGB(BackgroundColor())
	_, _, l := c.Hsl()
	return l < 0.5
}
