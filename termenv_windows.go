// +build windows

package termenv

func SupportedColorProfile() ColorProfile {
	return TrueColor
}

func ForegroundColor() Color {
	// default gray
	return ANSIColor(7)
}

func BackgroundColor() Color {
	// default black
	return ANSIColor(0)
}
