// +build windows

package termenv

// ColorProfile returns the supported color profile:
// Monochrome, ANSI, ANSI256, or TrueColor
func ColorProfile() Profile {
	return TrueColor
}

// ForegroundColor returns the terminal's default foreground color
func ForegroundColor() Color {
	// default gray
	return ANSIColor(7)
}

// BackgroundColor returns the terminal's default background color
func BackgroundColor() Color {
	// default black
	return ANSIColor(0)
}
