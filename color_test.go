package termenv

import "testing"

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

func TestRGBColorSequenceRoundsChannels(t *testing.T) {
	cases := []struct {
		hex string
		fg  string
		bg  string
	}{
		{"#657084", "38;2;101;112;132", "48;2;101;112;132"},
		{"#eceff4", "38;2;236;239;244", "48;2;236;239;244"},
		{"#4c566a", "38;2;76;86;106", "48;2;76;86;106"},
		{"#2d6271", "38;2;45;98;113", "48;2;45;98;113"},
		{"#fafafa", "38;2;250;250;250", "48;2;250;250;250"},
	}
	for _, tc := range cases {
		c := RGBColor(tc.hex)
		if got := c.Sequence(false); got != tc.fg {
			t.Errorf("%s fg: got %q want %q", tc.hex, got, tc.fg)
		}
		if got := c.Sequence(true); got != tc.bg {
			t.Errorf("%s bg: got %q want %q", tc.hex, got, tc.bg)
		}
	}
}
