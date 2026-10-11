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
		name     string
		color    RGBColor
		bg       bool
		expected string
	}{
		{
			name:     "hex #657084 foreground preserves channel precision",
			color:    RGBColor("#657084"),
			bg:       false,
			expected: "38;2;101;112;132",
		},
		{
			name:     "hex #eceff4 background preserves channel precision",
			color:    RGBColor("#eceff4"),
			bg:       true,
			expected: "48;2;236;239;244",
		},
		{
			name:     "hex #4c566a foreground preserves channel precision",
			color:    RGBColor("#4c566a"),
			bg:       false,
			expected: "38;2;76;86;106",
		},
		{
			name:     "hex #2d6271 background preserves channel precision",
			color:    RGBColor("#2d6271"),
			bg:       true,
			expected: "48;2;45;98;113",
		},
		{
			name:     "black foreground boundary",
			color:    RGBColor("#000000"),
			bg:       false,
			expected: "38;2;0;0;0",
		},
		{
			name:     "white background boundary",
			color:    RGBColor("#ffffff"),
			bg:       true,
			expected: "48;2;255;255;255",
		},
		{
			name:     "invalid hex returns empty string",
			color:    RGBColor("invalid"),
			bg:       false,
			expected: "",
		},
		{
			name:     "invalid hex background returns empty string",
			color:    RGBColor("not-a-color"),
			bg:       true,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.color.Sequence(tc.bg)
			if got != tc.expected {
				t.Fatalf("RGBColor(%q).Sequence(%t) = %q, want %q", tc.color, tc.bg, got, tc.expected)
			}
		})
	}

	t.Run("exhaustive channel range 0-255", func(t *testing.T) {
		for i := 0; i <= 255; i++ {
			r := uint8(i)
			g := uint8((i * 3) % 256)
			b := uint8(255 - i)
			hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
			color := RGBColor(hex)

			wantFg := fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
			if gotFg := color.Sequence(false); gotFg != wantFg {
				t.Fatalf("foreground truncation for %s: got %q, want %q", hex, gotFg, wantFg)
			}

			wantBg := fmt.Sprintf("48;2;%d;%d;%d", r, g, b)
			if gotBg := color.Sequence(true); gotBg != wantBg {
				t.Fatalf("background truncation for %s: got %q, want %q", hex, gotBg, wantBg)
			}
		}
	})
}

