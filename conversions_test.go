package color

import (
	"testing"
)

func TestConvertToRGBSpace(t *testing.T) {
	red := RGB(1, 0, 0)

	// Convert to display-p3
	displayP3, err := ConvertToRGBSpace(red, "display-p3")
	if err != nil {
		t.Fatalf("ConvertToRGBSpace failed: %v", err)
	}

	r, g, b, a := displayP3.RGBA()
	// Should be a valid color
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
		t.Errorf("ConvertToRGBSpace = RGB(%v, %v, %v, %v), want valid color", r, g, b, a)
	}

	// Should still be mostly red
	if r < 0.9 {
		t.Errorf("ConvertToRGBSpace red R = %v, expected > 0.9", r)
	}
}

func TestConvertFromRGBSpace(t *testing.T) {
	// Convert from display-p3 RGB values
	displayP3Color, err := ConvertFromRGBSpace(1.0, 0.0, 0.0, 1.0, "display-p3")
	if err != nil {
		t.Fatalf("ConvertFromRGBSpace failed: %v", err)
	}

	sc, ok := displayP3Color.(SpaceColor)
	if !ok {
		t.Fatalf("ConvertFromRGBSpace should preserve source as SpaceColor")
	}
	if sc.Space() == nil || sc.Space().Name() != "display-p3" {
		t.Fatalf("expected display-p3 source space, got %v", sc.Space())
	}
	channels := sc.Channels()
	if len(channels) < 3 {
		t.Fatalf("expected 3 channels, got %d", len(channels))
	}
	if !floatEqual(channels[0], 1.0) || !floatEqual(channels[1], 0.0) || !floatEqual(channels[2], 0.0) {
		t.Fatalf("expected preserved display-p3 channels [1 0 0], got %v", channels[:3])
	}

	// Should be convertible to other color spaces
	oklch := ToOKLCH(displayP3Color)
	if oklch == nil {
		t.Error("ConvertFromRGBSpace result should be convertible to OKLCH")
	}

	// Should be a valid color
	r, g, b, a := displayP3Color.RGBA()
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
		t.Errorf("ConvertFromRGBSpace = RGB(%v, %v, %v, %v), want valid color", r, g, b, a)
	}
}

func TestConvertBetweenAllSpaces(t *testing.T) {
	// Test that we can convert between all color spaces
	red := RGB(1, 0, 0)

	// Convert to all spaces
	hsl := ToHSL(red)
	hsv := ToHSV(red)
	lab := ToLAB(red)
	oklab := ToOKLAB(red)
	lch := ToLCH(red)
	oklch := ToOKLCH(red)
	xyz := ToXYZ(red)

	// Convert back to RGB
	spaces := []Color{hsl, hsv, lab, oklab, lch, oklch, xyz}
	for i, c := range spaces {
		r, g, b, a := c.RGBA()
		if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
			t.Errorf("Space %d produced invalid color: RGB(%v, %v, %v, %v)", i, r, g, b, a)
		}
	}
}

func TestConvertToRGBSpaceFromSpaceColorPreservesGamutPath(t *testing.T) {
	source := NewSpaceColor(DisplayP3Space, []float64{0.9, 0.2, 0.8}, 1.0)
	converted, err := ConvertToRGBSpace(source, "display-p3")
	if err != nil {
		t.Fatalf("ConvertToRGBSpace failed: %v", err)
	}

	got := converted.Channels()
	if len(got) < 3 {
		t.Fatalf("expected 3 channels, got %d", len(got))
	}

	// Same-space conversion should stay close to original channels.
	if abs(got[0]-0.9) > 0.03 || abs(got[1]-0.2) > 0.03 || abs(got[2]-0.8) > 0.03 {
		t.Fatalf("expected near [0.9 0.2 0.8], got %v", got[:3])
	}
}
