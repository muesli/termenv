package main

import (
	"fmt"

	"github.com/muesli/termenv"
)

func main() {
	var s termenv.Style

	compare(s)

	for i := 1; i < 255; i++ {
		if i < 16 {
			compare(s.Foreground(termenv.ANSIColor(i)))
		} else {
			compare(s.Foreground(termenv.ANSI256Color(i)))
		}
	}
}

func compare(s termenv.Style) {
	fmt.Print(s.Styled("Regular") + " | ")
	fmt.Print(s.Faint().Styled("Faint") + " | ")
	fmt.Print(s.ForceFaint().Styled("Forced") + " | ")
	fmt.Print(s.AdaptiveFaint().Styled("Adaptive"))
	fmt.Println()
}
