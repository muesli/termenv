// +build windows

package termenv

func colorProfile() Profile {
	return TrueColor
}

func foregroundColor() Color {
	// default gray
	return ANSIColor(7)
}

func backgroundColor() Color {
	// default black
	return ANSIColor(0)
}
