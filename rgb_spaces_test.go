package color

import (
	"testing"
)

func TestParseWideGamutRGB(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Display P3", "color(display-p3 1 0 0)", false},
		{"Display P3 with alpha", "color(display-p3 1 0 0 / 0.5)", false},
		{"sRGB linear", "color(srgb-linear 1 0 0)", false},
		{"Adobe RGB", "color(a98-rgb 1 0 0)", false},
		{"ProPhoto RGB", "color(prophoto-rgb 1 0 0)", false},
		{"Rec. 2020", "color(rec2020 1 0 0)", false},
		{"Invalid color space", "color(invalid 1 0 0)", true},
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

func TestRGBColorSpaceConversion(t *testing.T) {
	// Test that wide-gamut colors can be converted
	displayP3Red, err := ParseColor("color(display-p3 1 0 0)")
	if err != nil {
		t.Fatalf("ParseColor failed: %v", err)
	}

	r, g, b, a := displayP3Red.RGBA()
	// Display P3 red should convert to a valid sRGB color
	// (may be slightly different due to gamut mapping)
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
		t.Errorf("Display P3 red = RGB(%v, %v, %v, %v), want valid color", r, g, b, a)
	}

	// Red should still be mostly red
	if r < 0.9 {
		t.Errorf("Display P3 red R = %v, expected > 0.9", r)
	}
}

func TestParseWideGamutPreservesSpaceColor(t *testing.T) {
	c, err := ParseColor("color(display-p3 0.9 0.2 0.8)")
	if err != nil {
		t.Fatalf("ParseColor failed: %v", err)
	}

	sc, ok := c.(SpaceColor)
	if !ok {
		t.Fatalf("expected ParseColor(color(display-p3 ...)) to return SpaceColor")
	}

	if sc.Space() == nil || sc.Space().Name() != "display-p3" {
		t.Fatalf("expected display-p3 space, got %v", sc.Space())
	}

	channels := sc.Channels()
	if len(channels) < 3 {
		t.Fatalf("expected at least 3 channels, got %d", len(channels))
	}
	if abs(channels[0]-0.9) > 0.01 || abs(channels[1]-0.2) > 0.01 || abs(channels[2]-0.8) > 0.01 {
		t.Fatalf("expected preserved channels [0.9 0.2 0.8], got %v", channels[:3])
	}
}

func TestSRGBLinear(t *testing.T) {
	// sRGB-linear should handle values differently than sRGB
	srgbLinear, err := ParseColor("color(srgb-linear 0.5 0.5 0.5)")
	if err != nil {
		t.Fatalf("ParseColor failed: %v", err)
	}

	r, g, b, a := srgbLinear.RGBA()
	// Should be a valid color
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
		t.Errorf("sRGB-linear = RGB(%v, %v, %v, %v), want valid color", r, g, b, a)
	}
}
