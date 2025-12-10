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

