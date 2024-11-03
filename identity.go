package termenv

// TerminalIdentity is a guesstimate of a vague category of terminals that this program is running in. It should
// not be considered definitive, some categories overlap, and terminals also lie about this stuff all the time.
//
// It's used to get a general sense of certain things.  For instance, dumb terminals absolutely won't respond to cursor
// queries, and GNU Screen will sometimes say that it supports truecolor when it doesn't.
type TerminalIdentity int

const (
	// GNUScreen is a terminal multiplexer released in 1987
	GNUScreen TerminalIdentity = iota
	// TMux is a terminal multiplexer released in 2007
	TMux
	// DumbTerminal indicates a terminal with no capabilities- some embedded terminals identify themselves this way
	DumbTerminal
	// GoogleCloudShell is a browser-based CLI for GCP
	GoogleCloudShell
	// WindowsTerminalHosted indicates a terminal operating as a Windows Terminal tab.
	WindowsTerminalHosted
	// XTermCompatible means that the terminal has identified itself as being xterm, which may or may not be true
	XTermCompatible
	// OtherWindows means PowerShell or Windows' CMD, running standalone
	OtherWindows
	// OtherTerminal means any terminal that does not belong to one of the above categories
	OtherTerminal
)

func (i TerminalIdentity) String() string {
	switch i {
	case GNUScreen:
		return "GNUScreen"
	case TMux:
		return "TMux"
	case DumbTerminal:
		return "DumbTerminal"
	case GoogleCloudShell:
		return "GoogleCloudShell"
	case WindowsTerminalHosted:
		return "WindowsTerminalHosted"
	case XTermCompatible:
		return "XTermCompatible"
	case OtherWindows:
		return "OtherWindows"
	default:
		return "OtherTerminal"
	}
}
