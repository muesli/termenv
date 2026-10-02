package termenv

import (
	"fmt"
	"testing"
)

func TestXTermColor(t *testing.T) {
	var tests = []struct {
		input string
		color RGBColor
		valid bool
	}{
		{
			"\033]11;rgb:fafa/fafa/fafa\033",
			RGBColor("#fafafa"),
			true,
		},
		{
			"\033]11;rgb:fafa/fafa/fafa\033\\",
			RGBColor("#fafafa"),
			true,
		},
		{
			"\033]11;rgb:1212/3434/5656\a",
			RGBColor("#123456"),
			true,
		},
		{
			"\033]11;foo:fafa/fafa/fafaZZ",
			"",
			false,
		},
		{
			"\033]11;rgb:fafa/fafa",
			"",
			false,
		},
		{
			"\033]11;rgb:fafa/fafa/fafaY",
			"",
			false,
		},
		{
			"\033]11;rgb:fafa/fafa/fafaZZ",
			"",
			false,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			color, err := xTermColor(test.input)
			if err != nil && test.valid {
				t.Fatalf("unexpected error for input %q: %v", test.input, err)
			}

			if err == nil && !test.valid {
				t.Fatalf("expected error for input %v not found", test.input)
			}

			if color != test.color {
				t.Fatalf("wrong color returned, want %v, got %v", test.color, color)
			}
		})
	}
}

func TestRGBColorSequence(t *testing.T) {
	tests := []struct {
		color    RGBColor
		bg       bool
		expected string
	}{
		{RGBColor("#657084"), false, "38;2;101;112;132"},
		{RGBColor("#eceff4"), true, "48;2;236;239;244"},
		{RGBColor("#4c566a"), false, "38;2;76;86;106"},
		{RGBColor("#2d6271"), false, "38;2;45;98;113"},
		{RGBColor("invalid"), false, ""},
	}

	for _, tc := range tests {
		t.Run(string(tc.color), func(t *testing.T) {
			actual := tc.color.Sequence(tc.bg)
			if actual != tc.expected {
				t.Fatalf("wrong sequence returned, want %q, got %q", tc.expected, actual)
			}
		})
	}

	// Test all 256 channel values to ensure no floating-point truncation occurs
	for i := 0; i < 256; i++ {
		c := RGBColor(fmt.Sprintf("#%02x%02x%02x", i, i, i))
		expected := fmt.Sprintf("38;2;%d;%d;%d", i, i, i)
		if actual := c.Sequence(false); actual != expected {
			t.Fatalf("channel %d truncated for %s: want %q, got %q", i, c, expected, actual)
		}
	}
}
