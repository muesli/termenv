# termenv

[![Latest Release](https://img.shields.io/github/release/muesli/termenv.svg)](https://github.com/muesli/termenv/releases) [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/termenv) [![Build Status](https://github.com/muesli/termenv/workflows/build/badge.svg)](https://github.com/muesli/termenv/actions) [![Go ReportCard](http://goreportcard.com/badge/muesli/termenv)](http://goreportcard.com/report/muesli/termenv)

`termenv` gathers information about the terminal environment in terms of its
ANSI & color support. You can then use its convenient methods to colorize and
style your text output with ANSI escape sequences.

## Query Terminal Status

```go
// returns supported color profile: Monochrome, ANSI, ANSI256, or TrueColor
p := termenv.ColorProfile()

// returns default foreground color
fg := termenv.ForegroundColor()

// returns default background color
bg := termenv.BackgroundColor()

// returns whether terminal uses a dark-ish background
dark := termenv.HasDarkBackground()
```

## Apply Colors

`termenv` will automatically degrade colors to the closest available color
in the current color profile: `TrueColor` => `ANSI 256 Colors` =>
`ANSI 16 Colors` => `Monochrome`.

```go
p := termenv.ColorProfile()
out := termenv.String("Hello World")

// supports hex colors
out = out.Foreground(p.Color("#abcdef"))
// but also supports ANSI colors (0-255)
out = out.Background(p.Color("69"))

fmt.Println(out)
```

## Styles

```go
out := termenv.String("foobar")

// text styles
out.Bold()
out.Faint()
out.Italic()
out.Underline()
out.Overline()

// reverse swaps current fore- & background colors
out.Reverse()

// blinking text
out.Blink()
```
