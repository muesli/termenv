package main

import (
	"fmt"

	"github.com/muesli/termenv"
)

func main() {
	p := termenv.ColorProfile()

	fmt.Printf("\n\t%s %s %s %s %s",
		termenv.String("bold").Bold(),
		termenv.String("faint").Faint(),
		termenv.String("italic").Italic(),
		termenv.String("underline").Underline(),
		termenv.String("crossout").CrossOut(),
	)

	fmt.Printf("\n\t%s %s %s %s %s %s %s",
		termenv.String("red").Foreground(p.Color("#E88388")),
		termenv.String("green").Foreground(p.Color("#A8CC8C")),
		termenv.String("yellow").Foreground(p.Color("#DBAB79")),
		termenv.String("blue").Foreground(p.Color("#71BEF2")),
		termenv.String("magenta").Foreground(p.Color("#D290E4")),
		termenv.String("cyan").Foreground(p.Color("#66C2CD")),
		termenv.String("gray").Foreground(p.Color("#B9BFCA")),
	)

	fmt.Printf("\n\t%s %s %s %s %s %s %s\n\n",
		termenv.String("red").Foreground(p.Color("0")).Background(p.Color("#E88388")),
		termenv.String("green").Foreground(p.Color("0")).Background(p.Color("#A8CC8C")),
		termenv.String("yellow").Foreground(p.Color("0")).Background(p.Color("#DBAB79")),
		termenv.String("blue").Foreground(p.Color("0")).Background(p.Color("#71BEF2")),
		termenv.String("magenta").Foreground(p.Color("0")).Background(p.Color("#D290E4")),
		termenv.String("cyan").Foreground(p.Color("0")).Background(p.Color("#66C2CD")),
		termenv.String("gray").Foreground(p.Color("0")).Background(p.Color("#B9BFCA")),
	)

	fmt.Printf("\n\t%s %t\n", termenv.String("Has dark background?").Bold(), termenv.HasDarkBackground())
}
