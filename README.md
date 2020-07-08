<p align="center">
    <img src="https://stuff.charm.sh/termenv.png" width="480" alt="termenv Logo">
    <br />
    <a href="https://github.com/muesli/termenv/releases"><img src="https://img.shields.io/github/release/muesli/termenv.svg" alt="Latest Release"></a>
    <a href="https://godoc.org/github.com/muesli/termenv"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/muesli/termenv/actions"><img src="https://github.com/muesli/termenv/workflows/build/badge.svg" alt="Build Status"></a>
    <a href="https://coveralls.io/github/muesli/termenv?branch=master"><img src="https://coveralls.io/repos/github/muesli/termenv/badge.svg?branch=master" alt="Coverage Status"></a>
    <a href="http://goreportcard.com/report/muesli/termenv"><img src="http://goreportcard.com/badge/muesli/termenv" alt="Go ReportCard"></a>
</p>

`termenv` lets you safely use advanced styling options on the terminal. It
gathers information about the terminal environment in terms of its ANSI & color
support and offers you convenient methods to colorize and style your output,
without you having to deal with all kinds of weird ANSI escape sequences and
color conversions.

![Example output](https://github.com/muesli/termenv/raw/master/examples/hello-world/hello-world.png)

## Installation

```bash
go get github.com/muesli/termenv
```

## Query Terminal Status

```go
// returns supported color profile: Ascii, ANSI, ANSI256, or TrueColor
termenv.ColorProfile()

// returns default foreground color
termenv.ForegroundColor()

// returns default background color
termenv.BackgroundColor()

// returns whether terminal uses a dark-ish background
termenv.HasDarkBackground()
```

## Colors

`termenv` supports multiple color profiles: ANSI (16 colors), ANSI Extended
(256 colors), and TrueColor (24-bit RGB). Colors will automatically be degraded
to the best matching available color in the desired profile:

`TrueColor` => `ANSI 256 Colors` => `ANSI 16 Colors` => `Ascii`

```go
out := termenv.String("Hello World")

// retrieve color profile supported by terminal
p := termenv.ColorProfile()

// supports hex values
// will automatically degrade colors on terminals not supporting RGB
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
out.CrossOut()
out.Underline()
out.Overline()

// reverse swaps current fore- & background colors
out.Reverse()

// blinking text
out.Blink()

// combine multiple options
out.Bold().Underline()
```

## Template Helpers

```go
// load template helpers
f := termenv.TemplateFuncs(termenv.ColorProfile())
tpl := template.New("tpl").Funcs(f)

// apply bold style in a template
bold := `{{ Bold "Hello World" }}`

// examples for colorized templates
col := `{{ Color "#ff0000" "#0000ff" "Red on Blue" }}`
fg := `{{ Foreground "#ff0000" "Red Foreground" }}`
bg := `{{ Background "#0000ff" "Blue Background" }}`

// wrap styles
wrap := `{{ Bold (Underline "Hello World") }}`

// parse and render
tpl = tpl.Parse(bold)

var buf bytes.Buffer
tpl.Execute(&buf, nil)
fmt.Println(buf)
```

Other available helper functions are: `Faint`, `Italic`, `CrossOut`,
`Underline`, `Overline`, `Reverse`, and `Blink`.

## Screen

```go
// Reset the terminal to its default style, removing any active styles
termenv.Reset()

// Switch to the altscreen. The former view can be restored with ExitAltScreen()
termenv.AltScreen()

// Exit the altscreen and return to the former terminal view
termenv.ExitAltScreen()

// Clear the visible portion of the terminal
termenv.ClearScreen()

// Move the cursor to a given position
termenv.MoveCursor(row, column)

// Hide the cursor
termenv.HideCursor()

// Show the cursor
termenv.ShowCursor()

// Move the cursor up a given number of lines
termenv.CursorUp(n)

// Move the cursor down a given number of lines
termenv.CursorDown(n)

// Move the cursor down a given number of lines and place it at the beginning
// of the line
termenv.CursorNextLine(n)

// Move the cursor up a given number of lines and place it at the beginning of
// the line
termenv.CursorPrevLine(n)

// Clear the current line
termenv.ClearLine()

// Clear a given number of lines
termenv.ClearLines(n)
```

## Color Chart

![ANSI color chart](https://github.com/muesli/termenv/raw/master/examples/color-chart/color-chart.png)

You can find the source code used to create this chart in `termenv`'s examples.

## Related Projects

Check out [Glow](https://github.com/charmbracelet/glow), a markdown renderer for
the command-line, which uses `termenv`.

## License

[MIT](https://github.com/muesli/termenv/raw/master/LICENSE)
