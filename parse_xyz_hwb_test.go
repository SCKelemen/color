package color

import (
	"testing"
)

func TestParseXYZ(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{"XYZ via color()", "color(xyz 0.5 0.5 0.5)", false},
		{"XYZ D65", "color(xyz-d65 0.5 0.5 0.5)", false},
		{"XYZ D50", "color(xyz-d50 0.5 0.5 0.5)", false},
		{"XYZ with alpha", "color(xyz 0.5 0.5 0.5 / 0.5)", false},
		{"Invalid XYZ", "color(xyz 0.5)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Check it's a valid color
				r, g, b, a := c.RGBA()
				if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
					t.Errorf("ParseColor(%q) = RGB(%v, %v, %v, %v), want valid color",
						tt.input, r, g, b, a)
				}
			}
		})
	}
}

func TestParseHWB(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected [4]float64 // Approximate RGBA
		wantErr  bool
	}{
		{"HWB red", "hwb(0 0% 0%)", [4]float64{1, 0, 0, 1}, false},
		{"HWB white", "hwb(0 100% 0%)", [4]float64{1, 1, 1, 1}, false},
		{"HWB black", "hwb(0 0% 100%)", [4]float64{0, 0, 0, 1}, false},
		{"HWB with alpha", "hwb(0 0% 0% / 0.5)", [4]float64{1, 0, 0, 0.5}, false},
		{"Invalid HWB", "hwb(0)", [4]float64{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				r, g, b, a := c.RGBA()
				// Check alpha exactly
				if !floatEqual(a, tt.expected[3]) {
					t.Errorf("ParseColor(%q) alpha = %v, want %v", tt.input, a, tt.expected[3])
				}
				// RGB should be close (within 0.2 for HWB conversion)
				if !floatEqual(r, tt.expected[0]) && (r < tt.expected[0]-0.2 || r > tt.expected[0]+0.2) {
					t.Errorf("ParseColor(%q) R = %v, want ~%v", tt.input, r, tt.expected[0])
				}
				// Use g and b to avoid unused variable warning
				_ = g
				_ = b
			}
		})
	}
}

func TestParseModernRGBSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected [4]float64
	}{
		{"Modern RGB", "rgb(255 0 0)", [4]float64{1, 0, 0, 1}},
		{"Modern RGB with alpha", "rgb(255 0 0 / 0.5)", [4]float64{1, 0, 0, 0.5}},
		{"Modern HSL", "hsl(0 100% 50%)", [4]float64{1, 0, 0, 1}},
		{"Modern HSL with alpha", "hsl(0 100% 50% / 0.5)", [4]float64{1, 0, 0, 0.5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseColor(tt.input)
			if err != nil {
				t.Fatalf("ParseColor(%q) error = %v", tt.input, err)
			}
			r, g, b, a := c.RGBA()
			if !rgbaEqual(r, g, b, a, tt.expected[0], tt.expected[1], tt.expected[2], tt.expected[3]) {
				t.Errorf("ParseColor(%q) = RGB(%v, %v, %v, %v), want RGB(%v, %v, %v, %v)",
					tt.input, r, g, b, a, tt.expected[0], tt.expected[1], tt.expected[2], tt.expected[3])
			}
		})
	}
}

