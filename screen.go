package termenv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Sequence definitions.
const (
	// Cursor positioning.
	CursorUpSeq              = "%dA"
	CursorDownSeq            = "%dB"
	CursorForwardSeq         = "%dC"
	CursorBackSeq            = "%dD"
	CursorNextLineSeq        = "%dE"
	CursorPreviousLineSeq    = "%dF"
	CursorHorizontalSeq      = "%dG"
	CursorPositionSeq        = "%d;%dH"
	GetCursorPositionSeq     = "%c[6n"
	EraseDisplaySeq          = "%dJ"
	EraseLineSeq             = "%dK"
	ScrollUpSeq              = "%dS"
	ScrollDownSeq            = "%dT"
	SaveCursorPositionSeq    = "s"
	RestoreCursorPositionSeq = "u"
	ChangeScrollingRegionSeq = "%d;%dr"
	InsertLineSeq            = "%dL"
	DeleteLineSeq            = "%dM"

	// Explicit values for EraseLineSeq.
	EraseLineRightSeq  = "0K"
	EraseLineLeftSeq   = "1K"
	EraseEntireLineSeq = "2K"

	// Mouse.
	EnableMousePressSeq       = "?9h" // press only (X10)
	DisableMousePressSeq      = "?9l"
	EnableMouseSeq            = "?1000h" // press, release, wheel
	DisableMouseSeq           = "?1000l"
	EnableMouseHiliteSeq      = "?1001h" // highlight
	DisableMouseHiliteSeq     = "?1001l"
	EnableMouseCellMotionSeq  = "?1002h" // press, release, move on pressed, wheel
	DisableMouseCellMotionSeq = "?1002l"
	EnableMouseAllMotionSeq   = "?1003h" // press, release, move, wheel
	DisableMouseAllMotionSeq  = "?1003l"

	// Screen.
	RestoreScreenSeq = "?47l"
	SaveScreenSeq    = "?47h"
	AltScreenSeq     = "?1049h"
	ExitAltScreenSeq = "?1049l"

	// Session.
	SetWindowTitleSeq     = "2;%s\007"
	SetForegroundColorSeq = "10;%s\007"
	SetBackgroundColorSeq = "11;%s\007"
	SetCursorColorSeq     = "12;%s\007"
	ShowCursorSeq         = "?25h"
	HideCursorSeq         = "?25l"
)

// Reset the terminal to its default style, removing any active styles.
func (o Output) Reset() {
	fmt.Fprint(o.tty, CSI+ResetSeq+"m")
}

// SetForegroundColor sets the default foreground color.
func (o Output) SetForegroundColor(color Color) {
	fmt.Fprintf(o.tty, OSC+SetForegroundColorSeq, color)
}

// SetBackgroundColor sets the default background color.
func (o Output) SetBackgroundColor(color Color) {
	fmt.Fprintf(o.tty, OSC+SetBackgroundColorSeq, color)
}

// SetCursorColor sets the cursor color.
func (o Output) SetCursorColor(color Color) {
	fmt.Fprintf(o.tty, OSC+SetCursorColorSeq, color)
}

// RestoreScreen restores a previously saved screen state.
func (o Output) RestoreScreen() {
	fmt.Fprint(o.tty, CSI+RestoreScreenSeq)
}

// SaveScreen saves the screen state.
func (o Output) SaveScreen() {
	fmt.Fprint(o.tty, CSI+SaveScreenSeq)
}

// AltScreen switches to the alternate screen buffer. The former view can be
// restored with ExitAltScreen().
func (o Output) AltScreen() {
	fmt.Fprint(o.tty, CSI+AltScreenSeq)
}

// ExitAltScreen exits the alternate screen buffer and returns to the former
// terminal view.
func (o Output) ExitAltScreen() {
	fmt.Fprint(o.tty, CSI+ExitAltScreenSeq)
}

// ClearScreen clears the visible portion of the terminal.
func (o Output) ClearScreen() {
	fmt.Fprintf(o.tty, CSI+EraseDisplaySeq, 2)
	o.MoveCursor(1, 1)
}

// MoveCursor moves the cursor to a given position.
func (o Output) MoveCursor(row int, column int) {
	fmt.Fprintf(o.tty, CSI+CursorPositionSeq, row, column)
}

// HideCursor hides the cursor.
func (o Output) HideCursor() {
	fmt.Fprintf(o.tty, CSI+HideCursorSeq)
}

// ShowCursor shows the cursor.
func (o Output) ShowCursor() {
	fmt.Fprintf(o.tty, CSI+ShowCursorSeq)
}

// SaveCursorPosition saves the cursor position.
func (o Output) SaveCursorPosition() {
	fmt.Fprint(o.tty, CSI+SaveCursorPositionSeq)
}

// RestoreCursorPosition restores a saved cursor position.
func (o Output) RestoreCursorPosition() {
	fmt.Fprint(o.tty, CSI+RestoreCursorPositionSeq)
}

// CursorUp moves the cursor up a given number of lines.
func (o Output) CursorUp(n int) {
	fmt.Fprintf(o.tty, CSI+CursorUpSeq, n)
}

// CursorDown moves the cursor down a given number of lines.
func (o Output) CursorDown(n int) {
	fmt.Fprintf(o.tty, CSI+CursorDownSeq, n)
}

// CursorForward moves the cursor up a given number of lines.
func (o Output) CursorForward(n int) {
	fmt.Fprintf(o.tty, CSI+CursorForwardSeq, n)
}

// CursorBack moves the cursor backwards a given number of cells.
func (o Output) CursorBack(n int) {
	fmt.Fprintf(o.tty, CSI+CursorBackSeq, n)
}

// CursorNextLine moves the cursor down a given number of lines and places it at
// the beginning of the line.
func (o Output) CursorNextLine(n int) {
	fmt.Fprintf(o.tty, CSI+CursorNextLineSeq, n)
}

// CursorPrevLine moves the cursor up a given number of lines and places it at
// the beginning of the line.
func (o Output) CursorPrevLine(n int) {
	fmt.Fprintf(o.tty, CSI+CursorPreviousLineSeq, n)
}

// ClearLine clears the current line.
func (o Output) ClearLine() {
	fmt.Fprint(o.tty, CSI+EraseEntireLineSeq)
}

// ClearLineLeft clears the line to the left of the cursor.
func (o Output) ClearLineLeft() {
	fmt.Fprint(o.tty, CSI+EraseLineLeftSeq)
}

// ClearLineRight clears the line to the right of the cursor.
func (o Output) ClearLineRight() {
	fmt.Fprint(o.tty, CSI+EraseLineRightSeq)
}

// ClearLines clears a given number of lines.
func (o Output) ClearLines(n int) {
	clearLine := fmt.Sprintf(CSI+EraseLineSeq, 2)
	cursorUp := fmt.Sprintf(CSI+CursorUpSeq, 1)
	fmt.Fprint(o.tty, clearLine+strings.Repeat(cursorUp+clearLine, n))
}

// ChangeScrollingRegion sets the scrolling region of the terminal.
func (o Output) ChangeScrollingRegion(top, bottom int) {
	fmt.Fprintf(o.tty, CSI+ChangeScrollingRegionSeq, top, bottom)
}

// InsertLines inserts the given number of lines at the top of the scrollable
// region, pushing lines below down.
func (o Output) InsertLines(n int) {
	fmt.Fprintf(o.tty, CSI+InsertLineSeq, n)
}

// DeleteLines deletes the given number of lines, pulling any lines in
// the scrollable region below up.
func (o Output) DeleteLines(n int) {
	fmt.Fprintf(o.tty, CSI+DeleteLineSeq, n)
}

// EnableMousePress enables X10 mouse mode. Button press events are sent only.
func (o Output) EnableMousePress() {
	fmt.Fprint(o.tty, CSI+EnableMousePressSeq)
}

// DisableMousePress disables X10 mouse mode.
func (o Output) DisableMousePress() {
	fmt.Fprint(o.tty, CSI+DisableMousePressSeq)
}

// EnableMouse enables Mouse Tracking mode.
func (o Output) EnableMouse() {
	fmt.Fprint(o.tty, CSI+EnableMouseSeq)
}

// DisableMouse disables Mouse Tracking mode.
func (o Output) DisableMouse() {
	fmt.Fprint(o.tty, CSI+DisableMouseSeq)
}

// EnableMouseHilite enables Hilite Mouse Tracking mode.
func (o Output) EnableMouseHilite() {
	fmt.Fprint(o.tty, CSI+EnableMouseHiliteSeq)
}

// DisableMouseHilite disables Hilite Mouse Tracking mode.
func (o Output) DisableMouseHilite() {
	fmt.Fprint(o.tty, CSI+DisableMouseHiliteSeq)
}

// EnableMouseCellMotion enables Cell Motion Mouse Tracking mode.
func (o Output) EnableMouseCellMotion() {
	fmt.Fprint(o.tty, CSI+EnableMouseCellMotionSeq)
}

// DisableMouseCellMotion disables Cell Motion Mouse Tracking mode.
func (o Output) DisableMouseCellMotion() {
	fmt.Fprint(o.tty, CSI+DisableMouseCellMotionSeq)
}

// EnableMouseAllMotion enables All Motion Mouse mode.
func (o Output) EnableMouseAllMotion() {
	fmt.Fprint(o.tty, CSI+EnableMouseAllMotionSeq)
}

// DisableMouseAllMotion disables All Motion Mouse mode.
func (o Output) DisableMouseAllMotion() {
	fmt.Fprint(o.tty, CSI+DisableMouseAllMotionSeq)
}

// SetWindowTitle sets the terminal window title.
func (o Output) SetWindowTitle(title string) {
	fmt.Fprintf(o.tty, OSC+SetWindowTitleSeq, title)
}

// Legacy functions.

// Reset the terminal to its default style, removing any active styles.
func Reset() {
	NewOutputWithProfile(os.Stdout, ANSI).Reset()
}

// SetForegroundColor sets the default foreground color.
func SetForegroundColor(color Color) {
	NewOutputWithProfile(os.Stdout, ANSI).SetForegroundColor(color)
}

// SetBackgroundColor sets the default background color.
func SetBackgroundColor(color Color) {
	NewOutputWithProfile(os.Stdout, ANSI).SetBackgroundColor(color)
}

// SetCursorColor sets the cursor color.
func SetCursorColor(color Color) {
	NewOutputWithProfile(os.Stdout, ANSI).SetCursorColor(color)
}

// RestoreScreen restores a previously saved screen state.
func RestoreScreen() {
	NewOutputWithProfile(os.Stdout, ANSI).RestoreScreen()
}

// SaveScreen saves the screen state.
func SaveScreen() {
	NewOutputWithProfile(os.Stdout, ANSI).SaveScreen()
}

// AltScreen switches to the alternate screen buffer. The former view can be
// restored with ExitAltScreen().
func AltScreen() {
	NewOutputWithProfile(os.Stdout, ANSI).AltScreen()
}

// ExitAltScreen exits the alternate screen buffer and returns to the former
// terminal view.
func ExitAltScreen() {
	NewOutputWithProfile(os.Stdout, ANSI).ExitAltScreen()
}

// ClearScreen clears the visible portion of the terminal.
func ClearScreen() {
	NewOutputWithProfile(os.Stdout, ANSI).ClearScreen()
}

// MoveCursor moves the cursor to a given position.
func MoveCursor(row int, column int) {
	NewOutputWithProfile(os.Stdout, ANSI).MoveCursor(row, column)
}

// HideCursor hides the cursor.
func HideCursor() {
	NewOutputWithProfile(os.Stdout, ANSI).HideCursor()
}

// ShowCursor shows the cursor.
func ShowCursor() {
	NewOutputWithProfile(os.Stdout, ANSI).ShowCursor()
}

// SaveCursorPosition saves the cursor position.
func SaveCursorPosition() {
	NewOutputWithProfile(os.Stdout, ANSI).SaveCursorPosition()
}

// RestoreCursorPosition restores a saved cursor position.
func RestoreCursorPosition() {
	NewOutputWithProfile(os.Stdout, ANSI).RestoreCursorPosition()
}

// CursorUp moves the cursor up a given number of lines.
func CursorUp(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).CursorUp(n)
}

// CursorDown moves the cursor down a given number of lines.
func CursorDown(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).CursorDown(n)
}

// CursorForward moves the cursor up a given number of lines.
func CursorForward(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).CursorForward(n)
}

// CursorBack moves the cursor backwards a given number of cells.
func CursorBack(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).CursorBack(n)
}

// CursorNextLine moves the cursor down a given number of lines and places it at
// the beginning of the line.
func CursorNextLine(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).CursorNextLine(n)
}

// CursorPrevLine moves the cursor up a given number of lines and places it at
// the beginning of the line.
func CursorPrevLine(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).CursorPrevLine(n)
}

// ClearLine clears the current line.
func ClearLine() {
	NewOutputWithProfile(os.Stdout, ANSI).ClearLine()
}

// ClearLineLeft clears the line to the left of the cursor.
func ClearLineLeft() {
	NewOutputWithProfile(os.Stdout, ANSI).ClearLineLeft()
}

// ClearLineRight clears the line to the right of the cursor.
func ClearLineRight() {
	NewOutputWithProfile(os.Stdout, ANSI).ClearLineRight()
}

// ClearLines clears a given number of lines.
func ClearLines(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).ClearLines(n)
}

// ChangeScrollingRegion sets the scrolling region of the terminal.
func ChangeScrollingRegion(top, bottom int) {
	NewOutputWithProfile(os.Stdout, ANSI).ChangeScrollingRegion(top, bottom)
}

// InsertLines inserts the given number of lines at the top of the scrollable
// region, pushing lines below down.
func InsertLines(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).InsertLines(n)
}

// DeleteLines deletes the given number of lines, pulling any lines in
// the scrollable region below up.
func DeleteLines(n int) {
	NewOutputWithProfile(os.Stdout, ANSI).DeleteLines(n)
}

// EnableMousePress enables X10 mouse mode. Button press events are sent only.
func EnableMousePress() {
	NewOutputWithProfile(os.Stdout, ANSI).EnableMousePress()
}

// DisableMousePress disables X10 mouse mode.
func DisableMousePress() {
	NewOutputWithProfile(os.Stdout, ANSI).DisableMousePress()
}

// EnableMouse enables Mouse Tracking mode.
func EnableMouse() {
	NewOutputWithProfile(os.Stdout, ANSI).EnableMouse()
}

// DisableMouse disables Mouse Tracking mode.
func DisableMouse() {
	NewOutputWithProfile(os.Stdout, ANSI).DisableMouse()
}

// EnableMouseHilite enables Hilite Mouse Tracking mode.
func EnableMouseHilite() {
	NewOutputWithProfile(os.Stdout, ANSI).EnableMouseHilite()
}

// DisableMouseHilite disables Hilite Mouse Tracking mode.
func DisableMouseHilite() {
	NewOutputWithProfile(os.Stdout, ANSI).DisableMouseHilite()
}

// EnableMouseCellMotion enables Cell Motion Mouse Tracking mode.
func EnableMouseCellMotion() {
	NewOutputWithProfile(os.Stdout, ANSI).EnableMouseCellMotion()
}

// DisableMouseCellMotion disables Cell Motion Mouse Tracking mode.
func DisableMouseCellMotion() {
	NewOutputWithProfile(os.Stdout, ANSI).DisableMouseCellMotion()
}

// EnableMouseAllMotion enables All Motion Mouse mode.
func EnableMouseAllMotion() {
	NewOutputWithProfile(os.Stdout, ANSI).EnableMouseAllMotion()
}

// DisableMouseAllMotion disables All Motion Mouse mode.
func DisableMouseAllMotion() {
	NewOutputWithProfile(os.Stdout, ANSI).DisableMouseAllMotion()
}

// SetWindowTitle sets the terminal window title.
func SetWindowTitle(title string) {
	NewOutputWithProfile(os.Stdout, ANSI).SetWindowTitle(title)
}

// GetCursorPosition return the current position of the cursor on a terminal window in (row, column) format.
func GetCurosrPosition() (int, int, error) {
	/* The method this function uses is a bit out of ordinary. Essentially, it changes
	the command line to 'raw' mode, then prints an ANSI special character
	"\033[6n" in the terminal. The terminal then prints the cursor's position
	in this format: row;columnR . This function then parses this output to get row and
	column numbers. Finally, turns back the terminal mode from 'raw' to 'normal'.
	*/

	var row int
	var col int
	// Set the terminal to raw mode (to be undone with `-raw`)
	rawMode := exec.Command("/bin/stty", "raw")
	rawMode.Stdin = os.Stdin
	_ = rawMode.Run()
	err := rawMode.Wait()
	if err != nil {
		return -1, -1, err
	}

	// Running command $ echo -e "\033[6n" | read -dR
	cmd := exec.Command("echo", fmt.Sprintf(GetCursorPositionSeq, 27))
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	_ = cmd.Start()

	// capture keyboard output from echo command
	reader := bufio.NewReader(os.Stdin)
	err = cmd.Wait()
	if err != nil {
		return -1, -1, err
	}

	// by printing the command output, we are triggering input
	fmt.Print(randomBytes)
	text, _ := reader.ReadSlice('R') // The output ends with 'R'

	// check for the desired output
	if strings.Contains(string(text), ";") {
		re := regexp.MustCompile(`\d+;\d+`)
		line := re.FindString(string(text))
		delimiters := strings.Split(line, ";")
		// converting row and col strings to int
		row, _ = strconv.Atoi(delimiters[0])
		col, _ = strconv.Atoi(delimiters[1])
	} else {
		return -1, -1, errors.New("Could not parse cursor position output.")
	}
	// Set the terminal back from raw mode to 'normal'
	rawModeOff := exec.Command("/bin/stty", "-raw")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	err = rawModeOff.Wait()
	if err != nil {
		return -1, -1, err
	}
	return row, col, err
}
