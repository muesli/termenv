// +build windows

package termenv

func SupportedColorProfile() ColorProfile {
	return TrueColor
}

func DefaultForegroundColor() ColorSequencer {
	// default gray
	return ANSIColor(7)
}

func DefaultBackgroundColor() ColorSequencer {
	// default black
	return ANSIColor(0)
}
