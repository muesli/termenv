package termenv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type testEnv map[string]string

func (te testEnv) Environ() []string {
	var output []string
	for key, value := range te {
		output = append(output, fmt.Sprintf("%s=%s", key, value))
	}

	return output
}

func (te testEnv) Getenv(key string) string {
	value, ok := te[key]
	if !ok {
		return ""
	}

	return value
}

func tempOutput(t *testing.T) *Output {
	t.Helper()

	f, err := os.CreateTemp("", "termenv")
	if err != nil {
		t.Fatal(err)
	}

	return NewOutput(f, WithEnvironment(testEnv{
		"TERM": "xterm-256color",
	}), WithProfile(TrueColor))
}

func verify(t *testing.T, o *Output, exp string) {
	t.Helper()
	tty := o.w.(*os.File)

	if _, err := tty.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	b, err := io.ReadAll(tty)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != exp {
		b = bytes.ReplaceAll(b, []byte("\x1b"), []byte("\\x1b"))
		exp = strings.ReplaceAll(exp, "\x1b", "\\x1b")
		t.Errorf("output does not match, expected %s, got %s", exp, string(b))
	}

	// remove temp file
	os.Remove(tty.Name())
}

func TestReset(t *testing.T) {
	o := tempOutput(t)
	o.Reset()
	verify(t, o, "\x1b[0m")
}

func TestSetForegroundColor(t *testing.T) {
	o := tempOutput(t)
	o.SetForegroundColor(ANSI.Color("0"))
	verify(t, o, "\x1b]10;#000000\a")
}

func TestSetBackgroundColor(t *testing.T) {
	o := tempOutput(t)
	o.SetBackgroundColor(ANSI.Color("0"))
	verify(t, o, "\x1b]11;#000000\a")
}

func TestSetCursorColor(t *testing.T) {
	o := tempOutput(t)
	o.SetCursorColor(ANSI.Color("0"))
	verify(t, o, "\x1b]12;#000000\a")
}

func TestRestoreScreen(t *testing.T) {
	o := tempOutput(t)
	o.RestoreScreen()
	verify(t, o, "\x1b[?47l")
}

func TestSaveScreen(t *testing.T) {
	o := tempOutput(t)
	o.SaveScreen()
	verify(t, o, "\x1b[?47h")
}

func TestAltScreen(t *testing.T) {
	o := tempOutput(t)
	o.AltScreen()
	verify(t, o, "\x1b[?1049h")
}

func TestExitAltScreen(t *testing.T) {
	o := tempOutput(t)
	o.ExitAltScreen()
	verify(t, o, "\x1b[?1049l")
}

func TestClearScreen(t *testing.T) {
	o := tempOutput(t)
	o.ClearScreen()
	verify(t, o, "\x1b[2J\x1b[1;1H")
}

func TestMoveCursor(t *testing.T) {
	o := tempOutput(t)
	o.MoveCursor(16, 8)
	verify(t, o, "\x1b[16;8H")
}

func TestHideCursor(t *testing.T) {
	o := tempOutput(t)
	o.HideCursor()
	verify(t, o, "\x1b[?25l")
}

func TestShowCursor(t *testing.T) {
	o := tempOutput(t)
	o.ShowCursor()
	verify(t, o, "\x1b[?25h")
}

func TestSaveCursorPosition(t *testing.T) {
	o := tempOutput(t)
	o.SaveCursorPosition()
	verify(t, o, "\x1b[s")
}

func TestRestoreCursorPosition(t *testing.T) {
	o := tempOutput(t)
	o.RestoreCursorPosition()
	verify(t, o, "\x1b[u")
}

func TestCursorUp(t *testing.T) {
	o := tempOutput(t)
	o.CursorUp(8)
	verify(t, o, "\x1b[8A")
}

func TestCursorDown(t *testing.T) {
	o := tempOutput(t)
	o.CursorDown(8)
	verify(t, o, "\x1b[8B")
}

func TestCursorForward(t *testing.T) {
	o := tempOutput(t)
	o.CursorForward(8)
	verify(t, o, "\x1b[8C")
}

func TestCursorBack(t *testing.T) {
	o := tempOutput(t)
	o.CursorBack(8)
	verify(t, o, "\x1b[8D")
}

func TestCursorNextLine(t *testing.T) {
	o := tempOutput(t)
	o.CursorNextLine(8)
	verify(t, o, "\x1b[8E")
}

func TestCursorPrevLine(t *testing.T) {
	o := tempOutput(t)
	o.CursorPrevLine(8)
	verify(t, o, "\x1b[8F")
}

func TestClearLine(t *testing.T) {
	o := tempOutput(t)
	o.ClearLine()
	verify(t, o, "\x1b[2K")
}

func TestClearLineLeft(t *testing.T) {
	o := tempOutput(t)
	o.ClearLineLeft()
	verify(t, o, "\x1b[1K")
}

func TestClearLineRight(t *testing.T) {
	o := tempOutput(t)
	o.ClearLineRight()
	verify(t, o, "\x1b[0K")
}

func TestClearLines(t *testing.T) {
	o := tempOutput(t)
	o.ClearLines(8)
	verify(t, o, "\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[1A\x1b[2K")
}

func TestChangeScrollingRegion(t *testing.T) {
	o := tempOutput(t)
	o.ChangeScrollingRegion(16, 8)
	verify(t, o, "\x1b[16;8r")
}

func TestInsertLines(t *testing.T) {
	o := tempOutput(t)
	o.InsertLines(8)
	verify(t, o, "\x1b[8L")
}

func TestDeleteLines(t *testing.T) {
	o := tempOutput(t)
	o.DeleteLines(8)
	verify(t, o, "\x1b[8M")
}

func TestEnableMousePress(t *testing.T) {
	o := tempOutput(t)
	o.EnableMousePress()
	verify(t, o, "\x1b[?9h")
}

func TestDisableMousePress(t *testing.T) {
	o := tempOutput(t)
	o.DisableMousePress()
	verify(t, o, "\x1b[?9l")
}

func TestEnableMouse(t *testing.T) {
	o := tempOutput(t)
	o.EnableMouse()
	verify(t, o, "\x1b[?1000h")
}

func TestDisableMouse(t *testing.T) {
	o := tempOutput(t)
	o.DisableMouse()
	verify(t, o, "\x1b[?1000l")
}

func TestEnableMouseHilite(t *testing.T) {
	o := tempOutput(t)
	o.EnableMouseHilite()
	verify(t, o, "\x1b[?1001h")
}

func TestDisableMouseHilite(t *testing.T) {
	o := tempOutput(t)
	o.DisableMouseHilite()
	verify(t, o, "\x1b[?1001l")
}

func TestEnableMouseCellMotion(t *testing.T) {
	o := tempOutput(t)
	o.EnableMouseCellMotion()
	verify(t, o, "\x1b[?1002h")
}

func TestDisableMouseCellMotion(t *testing.T) {
	o := tempOutput(t)
	o.DisableMouseCellMotion()
	verify(t, o, "\x1b[?1002l")
}

func TestEnableMouseAllMotion(t *testing.T) {
	o := tempOutput(t)
	o.EnableMouseAllMotion()
	verify(t, o, "\x1b[?1003h")
}

func TestDisableMouseAllMotion(t *testing.T) {
	o := tempOutput(t)
	o.DisableMouseAllMotion()
	verify(t, o, "\x1b[?1003l")
}

func TestEnableMouseExtendedMode(t *testing.T) {
	o := tempOutput(t)
	o.EnableMouseExtendedMode()
	verify(t, o, "\x1b[?1006h")
}

func TestDisableMouseExtendedMode(t *testing.T) {
	o := tempOutput(t)
	o.DisableMouseExtendedMode()
	verify(t, o, "\x1b[?1006l")
}

func TestEnableMousePixelsMode(t *testing.T) {
	o := tempOutput(t)
	o.EnableMousePixelsMode()
	verify(t, o, "\x1b[?1016h")
}

func TestDisableMousePixelsMode(t *testing.T) {
	o := tempOutput(t)
	o.DisableMousePixelsMode()
	verify(t, o, "\x1b[?1016l")
}

func TestSetWindowTitle(t *testing.T) {
	o := tempOutput(t)
	o.SetWindowTitle("test")
	verify(t, o, "\x1b]2;test\a")
}

func TestCopyClipboard(t *testing.T) {
	o := tempOutput(t)
	o.Copy("hello")
	verify(t, o, "\x1b]52;c;aGVsbG8=\a")
}

func TestCopyPrimary(t *testing.T) {
	o := tempOutput(t)
	o.CopyPrimary("hello")
	verify(t, o, "\x1b]52;p;aGVsbG8=\a")
}

func TestHyperlink(t *testing.T) {
	o := tempOutput(t)
	o.WriteString(o.Hyperlink("http://example.com", "example"))
	verify(t, o, "\x1b]8;;http://example.com\x1b\\example\x1b]8;;\x1b\\")
}
