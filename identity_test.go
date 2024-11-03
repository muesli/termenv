package termenv

import (
	"bufio"
	"bytes"
	"testing"
)

func TestTerminalIdentity(t *testing.T) {
	testCases := map[string]struct {
		Env              testEnv
		ExpectedIdentity TerminalIdentity
	}{
		"Google Cloud Shell": {
			testEnv{
				"GOOGLE_CLOUD_SHELL": "true",
			},
			GoogleCloudShell,
		},
		"Other": {
			testEnv{},
			OtherTerminal,
		},
		"Windows Terminal": {
			testEnv{
				"WT_SESSION": "aec55a8a-d1c8-48a9-93f9-65da8c5e468d",
				"GOOS":       "windows",
			},
			WindowsTerminalHosted,
		},
		"Windows Terminal Hosting XTerm Compatible": {
			testEnv{
				"WT_SESSION": "aec55a8a-d1c8-48a9-93f9-65da8c5e468d",
				"GOOS":       "windows",
				"TERM":       "xterm-256color",
			},
			WindowsTerminalHosted,
		},
		"Windows Terminal Hosting WSL": {
			testEnv{
				"WT_SESSION": "aec55a8a-d1c8-48a9-93f9-65da8c5e468d",
				"GOOS":       "linux",
				"TERM":       "xterm-256color",
			},
			WindowsTerminalHosted,
		},
		"GNU Screen": {
			testEnv{
				"TERM": "screen-256color",
			},
			GNUScreen,
		},
		"TMux Pretending To Be Screen": {
			testEnv{
				"TERM":         "screen-256color",
				"TERM_PROGRAM": "tmux",
			},
			TMux,
		},
		"TMux Being Honest": {
			testEnv{
				"TERM": "tmux-256color",
			},
			TMux,
		},
		"Dumb Terminal": {
			testEnv{
				"TERM": "dumb-256color",
			},
			DumbTerminal,
		},
		"XTerm Compatible": {
			testEnv{
				"TERM": "xterm-256color",
			},
			XTermCompatible,
		},
		"XTerm Compatible In Windows": {
			testEnv{
				"TERM": "xterm-256color",
				"GOOS": "windows",
			},
			XTermCompatible,
		},
		"Windows Native Terminal": {
			testEnv{
				"GOOS": "windows",
			},
			OtherWindows,
		},
		"Wezterm": {
			testEnv{
				"TERM": "wezterm",
			},
			OtherTerminal,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			output := NewOutput(bufio.NewWriter(bytes.NewBuffer(nil)), WithTTY(true), WithEnvironment(testCase.Env))

			if output.TerminalIdentity != testCase.ExpectedIdentity {
				t.Errorf("output does not match, expected %s, got %s", testCase.ExpectedIdentity, output.TerminalIdentity)
			}
		})
	}
}
