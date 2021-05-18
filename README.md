<p align="center">
    <img src="https://stuff.charm.sh/termenv.png" width="480" alt="termenv Logo">
    <br />
    <a href="https://github.com/muesli/termenv/releases"><img src="https://img.shields.io/github/release/muesli/termenv.svg" alt="Latest Release"></a>
    <a href="https://godoc.org/github.com/muesli/termenv"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/muesli/termenv/actions"><img src="https://github.com/muesli/termenv/workflows/build/badge.svg" alt="Build Status"></a>
    <a href="https://coveralls.io/github/muesli/termenv?branch=master"><img src="https://coveralls.io/repos/github/muesli/termenv/badge.svg?branch=master" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/muesli/termenv"><img src="https://goreportcard.com/badge/muesli/termenv" alt="Go ReportCard"></a>
</p>

`termenv` lets you safely use advanced styling options on the terminal. It
gathers information about the terminal environment in terms of its ANSI & color
support and offers you convenient methods to colorize and style your output,
without you having to deal with all kinds of weird ANSI escape sequences and
color conversions.

![Example output](https://github.com/muesli/termenv/raw/master/examples/hello-world/hello-world.png)

## Features

- RGB/TrueColor support
- Detects the supported color range of your terminal
- Automatically converts colors to the best matching, available colors
- Terminal theme (light/dark) detection
- Chainable syntax
- Nested styles

## Installation

```bash
go get github.com/muesli/termenv
```

## Query Terminal Support

`termenv` can query the terminal it is running in, so you can safely use
advanced features, like RGB colors. `ColorProfile` returns the color profile
supported by the terminal:

```go
profile := termenv.ColorProfile()
```

This returns one of the supported color profiles:

- `termenv.Ascii` - no ANSI support detected, ASCII only
- `termenv.ANSI` - 16 color ANSI support
- `termenv.ANSI256` - Extended 256 color ANSI support
- `termenv.TrueColor` - RGB/TrueColor support

You can also query the terminal for its color scheme, so you know whether your
app is running in a light- or dark-themed environment:

```go
// Returns terminal's foreground color
color := termenv.ForegroundColor()

// Returns terminal's background color
color := termenv.BackgroundColor()

// Returns whether terminal uses a dark-ish background
darkTheme := termenv.HasDarkBackground()
```

## Colors

`termenv` supports multiple color profiles: ANSI (16 colors), ANSI Extended
(256 colors), and TrueColor (24-bit RGB). Colors will automatically be degraded
to the best matching available color in the desired profile:

`TrueColor` => `ANSI 256 Colors` => `ANSI 16 Colors` => `Ascii`

```go
s := termenv.String("Hello World")

// Retrieve color profile supported by terminal
p := termenv.ColorProfile()

// Supports hex values
// Will automatically degrade colors on terminals not supporting RGB
s.Foreground(p.Color("#abcdef"))
// but also supports ANSI colors (0-255)
s.Background(p.Color("69"))
// ...or the color.Color interface
s.Foreground(p.FromColor(color.RGBA{255, 128, 0, 255}))

// Combine fore- & background colors
s.Foreground(p.Color("#ffffff")).Background(p.Color("#0000ff"))

// Supports the fmt.Stringer interface
fmt.Println(s)
```

## Styles

You can use a chainable syntax to compose your own styles:

```go
s := termenv.String("foobar")

// Text styles
s.Bold()
s.Faint()
s.Italic()
s.CrossOut()
s.Underline()
s.Overline()

// Reverse swaps current fore- & background colors
s.Reverse()

// Blinking text
s.Blink()

// Combine multiple options
s.Bold().Underline()
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
tpl, err = tpl.Parse(bold)

var buf bytes.Buffer
tpl.Execute(&buf, nil)
fmt.Println(&buf)
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

// Save the cursor position
termenv.SaveCursorPosition()

// Restore a saved cursor position
termenv.RestoreCursorPosition()

// Move the cursor up a given number of lines
termenv.CursorUp(n)

// Move the cursor down a given number of lines
termenv.CursorDown(n)

// Move the cursor up a given number of lines
termenv.CursorForward(n)

// Move the cursor backwards a given number of cells
termenv.CursorBack(n)

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

// Set the scrolling region of the terminal
termenv.ChangeScrollingRegion(top, bottom)

// Insert the given number of lines at the top of the scrollable region, pushing
// lines below down
termenv.InsertLines(n)

// Delete the given number of lines, pulling any lines in the scrollable region
// below up
termenv.DeleteLines(n)
```

## Mouse

```go
// Enable X10 mouse mode, only button press events are sent
termenv.EnableMousePress()

// Disable X10 mouse mode
termenv.DisableMousePress()

// Enable Mouse Tracking mode
termenv.EnableMouse()

// Disable Mouse Tracking mode
termenv.DisableMouse()

// Enable Hilite Mouse Tracking mode
termenv.EnableMouseHilite()

// Disable Hilite Mouse Tracking mode
termenv.DisableMouseHilite()

// Enable Cell Motion Mouse Tracking mode
termenv.EnableMouseCellMotion()

// Disable Cell Motion Mouse Tracking mode
termenv.DisableMouseCellMotion()

// Enable All Motion Mouse mode
termenv.EnableMouseAllMotion()

// Disable All Motion Mouse mode
termenv.DisableMouseAllMotion()
```

## Color Chart

![ANSI color chart](https://github.com/muesli/termenv/raw/master/examples/color-chart/color-chart.png)

You can find the source code used to create this chart in `termenv`'s examples.

## Related Projects

- [reflow](https://github.com/muesli/reflow) - ANSI-aware text operations
- [Glow](https://github.com/charmbracelet/glow) - a markdown renderer for
the command-line, which uses `termenv`

## License

[MIT](https://github.com/muesli/termenv/raw/master/LICENSE)
