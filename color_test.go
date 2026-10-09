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

func TestRGBColorSequenceChannelRounding(t *testing.T) {
	// colorful.Hex divides the 8-bit channels by 255, and the float can land
	// just below the original integer. Truncating with uint8 then drops a step.
	for _, tt := range []struct {
		color  RGBColor
		fg, bg string
	}{
		{"#657084", "38;2;101;112;132", "48;2;101;112;132"},
		{"#eceff4", "38;2;236;239;244", "48;2;236;239;244"},
		{"#4c566a", "38;2;76;86;106", "48;2;76;86;106"},
		{"#2d6271", "38;2;45;98;113", "48;2;45;98;113"},
		{"#000000", "38;2;0;0;0", "48;2;0;0;0"},
		{"#ffffff", "38;2;255;255;255", "48;2;255;255;255"},
	} {
		if got := tt.color.Sequence(false); got != tt.fg {
			t.Errorf("Sequence(false) for %q: want %q, got %q", tt.color, tt.fg, got)
		}
		if got := tt.color.Sequence(true); got != tt.bg {
			t.Errorf("Sequence(true) for %q: want %q, got %q", tt.color, tt.bg, got)
		}
	}
}

func TestRGBColorSequenceAllChannelValues(t *testing.T) {
	// every 8-bit value must survive the round-trip through colorful's floats.
	for i := 0; i < 256; i++ {
		c := RGBColor(fmt.Sprintf("#%02x%02x%02x", i, i, i))
		want := fmt.Sprintf("38;2;%d;%d;%d", i, i, i)
		if got := c.Sequence(false); got != want {
			t.Fatalf("Sequence(false) for %q: want %q, got %q", c, want, got)
		}
	}
}
