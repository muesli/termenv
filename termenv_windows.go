// +build windows

package termenv

func SupportedColorProfile() ColorProfile {
	return TrueColor
}

func DefaultForegroundColor() ColorSequencer {
	// default white
	return 15
}

func DefaultBackgroundColor() ColorSequencer {
	// default black
	return 0
}
