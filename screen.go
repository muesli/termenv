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
	GetCursorPosition        = "%c[6n"
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
func Reset() {
	fmt.Print(CSI + ResetSeq + "m")
}

// SetForegroundColor sets the default foreground color.
func SetForegroundColor(color Color) {
	fmt.Printf(OSC+SetForegroundColorSeq, color)
}

// SetBackgroundColor sets the default background color.
func SetBackgroundColor(color Color) {
	fmt.Printf(OSC+SetBackgroundColorSeq, color)
}

// SetCursorColor sets the cursor color.
func SetCursorColor(color Color) {
	fmt.Printf(OSC+SetCursorColorSeq, color)
}

// RestoreScreen restores a previously saved screen state.
func RestoreScreen() {
	fmt.Print(CSI + RestoreScreenSeq)
}

// SaveScreen saves the screen state.
func SaveScreen() {
	fmt.Print(CSI + SaveScreenSeq)
}

// AltScreen switches to the alternate screen buffer. The former view can be
// restored with ExitAltScreen().
func AltScreen() {
	fmt.Print(CSI + AltScreenSeq)
}

// ExitAltScreen exits the alternate screen buffer and returns to the former
// terminal view.
func ExitAltScreen() {
	fmt.Print(CSI + ExitAltScreenSeq)
}

// ClearScreen clears the visible portion of the terminal.
func ClearScreen() {
	fmt.Printf(CSI+EraseDisplaySeq, 2)
	MoveCursor(1, 1)
}

// MoveCursor moves the cursor to a given position.
func MoveCursor(row int, column int) {
	fmt.Printf(CSI+CursorPositionSeq, row, column)
}

// HideCursor hides the cursor.
func HideCursor() {
	fmt.Printf(CSI + HideCursorSeq)
}

// ShowCursor shows the cursor.
func ShowCursor() {
	fmt.Printf(CSI + ShowCursorSeq)
}

// SaveCursorPosition saves the cursor position.
func SaveCursorPosition() {
	fmt.Print(CSI + SaveCursorPositionSeq)
}

// RestoreCursorPosition restores a saved cursor position.
func RestoreCursorPosition() {
	fmt.Print(CSI + RestoreCursorPositionSeq)
}

// CursorUp moves the cursor up a given number of lines.
func CursorUp(n int) {
	fmt.Printf(CSI+CursorUpSeq, n)
}

// CursorDown moves the cursor down a given number of lines.
func CursorDown(n int) {
	fmt.Printf(CSI+CursorDownSeq, n)
}

// CursorForward moves the cursor up a given number of lines.
func CursorForward(n int) {
	fmt.Printf(CSI+CursorForwardSeq, n)
}

// CursorBack moves the cursor backwards a given number of cells.
func CursorBack(n int) {
	fmt.Printf(CSI+CursorBackSeq, n)
}

// CursorNextLine moves the cursor down a given number of lines and places it at
// the beginning of the line.
func CursorNextLine(n int) {
	fmt.Printf(CSI+CursorNextLineSeq, n)
}

// CursorPrevLine moves the cursor up a given number of lines and places it at
// the beginning of the line.
func CursorPrevLine(n int) {
	fmt.Printf(CSI+CursorPreviousLineSeq, n)
}

// ClearLine clears the current line.
func ClearLine() {
	fmt.Print(CSI + EraseEntireLineSeq)
}

// ClearLineLeft clears the line to the left of the cursor.
func ClearLineLeft() {
	fmt.Print(CSI + EraseLineLeftSeq)
}

// ClearLineRight clears the line to the right of the cursor.
func ClearLineRight() {
	fmt.Print(CSI + EraseLineRightSeq)
}

// ClearLines clears a given number of lines.
func ClearLines(n int) {
	clearLine := fmt.Sprintf(CSI+EraseLineSeq, 2)
	cursorUp := fmt.Sprintf(CSI+CursorUpSeq, 1)
	fmt.Print(clearLine + strings.Repeat(cursorUp+clearLine, n))
}

// ChangeScrollingRegion sets the scrolling region of the terminal.
func ChangeScrollingRegion(top, bottom int) {
	fmt.Printf(CSI+ChangeScrollingRegionSeq, top, bottom)
}

// InsertLines inserts the given number of lines at the top of the scrollable
// region, pushing lines below down.
func InsertLines(n int) {
	fmt.Printf(CSI+InsertLineSeq, n)
}

// DeleteLines deletes the given number of lines, pulling any lines in
// the scrollable region below up.
func DeleteLines(n int) {
	fmt.Printf(CSI+DeleteLineSeq, n)
}

// EnableMousePress enables X10 mouse mode. Button press events are sent only.
func EnableMousePress() {
	fmt.Print(CSI + EnableMousePressSeq)
}

// DisableMousePress disables X10 mouse mode.
func DisableMousePress() {
	fmt.Print(CSI + DisableMousePressSeq)
}

// EnableMouse enables Mouse Tracking mode.
func EnableMouse() {
	fmt.Print(CSI + EnableMouseSeq)
}

// DisableMouse disables Mouse Tracking mode.
func DisableMouse() {
	fmt.Print(CSI + DisableMouseSeq)
}

// EnableMouseHilite enables Hilite Mouse Tracking mode.
func EnableMouseHilite() {
	fmt.Print(CSI + EnableMouseHiliteSeq)
}

// DisableMouseHilite disables Hilite Mouse Tracking mode.
func DisableMouseHilite() {
	fmt.Print(CSI + DisableMouseHiliteSeq)
}

// EnableMouseCellMotion enables Cell Motion Mouse Tracking mode.
func EnableMouseCellMotion() {
	fmt.Print(CSI + EnableMouseCellMotionSeq)
}

// DisableMouseCellMotion disables Cell Motion Mouse Tracking mode.
func DisableMouseCellMotion() {
	fmt.Print(CSI + DisableMouseCellMotionSeq)
}

// EnableMouseAllMotion enables All Motion Mouse mode.
func EnableMouseAllMotion() {
	fmt.Print(CSI + EnableMouseAllMotionSeq)
}

// DisableMouseAllMotion disables All Motion Mouse mode.
func DisableMouseAllMotion() {
	fmt.Print(CSI + DisableMouseAllMotionSeq)
}

// SetWindowTitle sets the terminal window title.
func SetWindowTitle(title string) {
	fmt.Printf(OSC+SetWindowTitleSeq, title)
}

// GetPosition return the current position of the cursor on a terminal window in (row, column) format.
func GetPosition() (int, int, error) {
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
	cmd := exec.Command("echo", fmt.Sprintf(GetCursorPosition, 27))
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
