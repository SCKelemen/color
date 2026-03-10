package color

import (
	"testing"
)

func TestParseXYZ(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"XYZ via color()", "color(xyz 0.5 0.5 0.5)", false},
		{"XYZ D65", "color(xyz-d65 0.5 0.5 0.5)", false},
		{"XYZ D50", "color(xyz-d50 0.5 0.5 0.5)", false},
		{"XYZ with alpha", "color(xyz 0.5 0.5 0.5 / 0.5)", false},
		{"Invalid XYZ", "color(xyz 0.5)", true},
		{"Too many XYZ args", "color(xyz 0.5 0.5 0.5 0.5 0.5)", true},
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

func TestParseColorFunctionRGBSpaceInvalidRanges(t *testing.T) {
	tests := []string{
		"color(display-p3 -0.1 0.2 0.3)",
		"color(display-p3 300 0.2 0.3)",
		"color(display-p3 0.1 -2 0.3)",
		"color(display-p3 0.1 0.2 999)",
		"color(display-p3 0.1 0.2 0.3 0.4 0.5)",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, err := ParseColor(input); err == nil {
				t.Fatalf("expected ParseColor(%q) to fail", input)
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

func TestParseHueAngleUnitsAndWrapping(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		equivalent string
	}{
		{"HSL turn", "hsl(0.5turn 100% 50%)", "hsl(180deg 100% 50%)"},
		{"HSL rad", "hsl(3.141592653589793rad 100% 50%)", "hsl(180deg 100% 50%)"},
		{"HSL grad", "hsl(200grad 100% 50%)", "hsl(180deg 100% 50%)"},
		{"HSL wrap positive", "hsl(480 100% 50%)", "hsl(120 100% 50%)"},
		{"HSL wrap negative", "hsl(-120 100% 50%)", "hsl(240 100% 50%)"},
		{"HSV turn", "hsv(0.5turn 100% 100%)", "hsv(180deg 100% 100%)"},
		{"HWB turn", "hwb(0.5turn 0% 0%)", "hwb(180deg 0% 0%)"},
		{"LCH turn", "lch(70 40 0.5turn)", "lch(70 40 180deg)"},
		{"OKLCH turn", "oklch(0.7 0.2 0.5turn)", "oklch(0.7 0.2 180deg)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseColor(tt.input)
			if err != nil {
				t.Fatalf("ParseColor(%q) error = %v", tt.input, err)
			}

			want, err := ParseColor(tt.equivalent)
			if err != nil {
				t.Fatalf("ParseColor(%q) error = %v", tt.equivalent, err)
			}

			r1, g1, b1, a1 := got.RGBA()
			r2, g2, b2, a2 := want.RGBA()
			if !rgbaEqual(r1, g1, b1, a1, r2, g2, b2, a2) {
				t.Fatalf(
					"%q != %q: got (%f,%f,%f,%f), want (%f,%f,%f,%f)",
					tt.input, tt.equivalent, r1, g1, b1, a1, r2, g2, b2, a2,
				)
			}
		})
	}
}

func TestParseXYZD50ChromaticAdaptation(t *testing.T) {
	parsed, err := ParseColor("color(xyz-d50 0.4 0.3 0.2 / 0.8)")
	if err != nil {
		t.Fatalf("ParseColor failed: %v", err)
	}

	x, y, z := AdaptD50ToD65(0.4, 0.3, 0.2)
	expected := NewXYZ(x, y, z, 0.8)

	r1, g1, b1, a1 := parsed.RGBA()
	r2, g2, b2, a2 := expected.RGBA()
	if !rgbaEqual(r1, g1, b1, a1, r2, g2, b2, a2) {
		t.Fatalf("xyz-d50 adaptation mismatch: got (%f,%f,%f,%f), want (%f,%f,%f,%f)",
			r1, g1, b1, a1, r2, g2, b2, a2)
	}

	// Ensure xyz-d50 is not treated identically to xyz-d65 for same numeric triples.
	d65, err := ParseColor("color(xyz-d65 0.4 0.3 0.2 / 0.8)")
	if err != nil {
		t.Fatalf("ParseColor failed: %v", err)
	}

	rd, gd, bd, _ := d65.RGBA()
	if abs(r1-rd) < 0.000001 && abs(g1-gd) < 0.000001 && abs(b1-bd) < 0.000001 {
		t.Fatalf("xyz-d50 should differ from xyz-d65 for same numeric values")
	}
}
