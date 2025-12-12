package color

import (
	"testing"
)

// Test uncovered String() methods
func TestStringMethods(t *testing.T) {
	// Test RGBA String
	rgba := NewRGBA(1, 0.5, 0, 1.0)
	str := rgba.String()
	if str == "" {
		t.Error("RGBA String() returned empty string")
	}

	// Test OKLCH String
	oklch := NewOKLCH(0.7, 0.2, 180, 1.0)
	str2 := oklch.String()
	if str2 == "" {
		t.Error("OKLCH String() returned empty string")
	}

	// Test OKLCH String with alpha
	oklchAlpha := NewOKLCH(0.7, 0.2, 180, 0.5)
	str3 := oklchAlpha.String()
	if str3 == "" {
		t.Error("OKLCH String() with alpha returned empty string")
	}
}

// Test WithAlpha methods
func TestWithAlphaMethods(t *testing.T) {
	// HSL WithAlpha
	hsl := NewHSL(180, 0.5, 0.5, 1.0)
	hslWithAlpha := hsl.WithAlpha(0.5)
	if hslWithAlpha.Alpha() != 0.5 {
		t.Errorf("HSL WithAlpha failed: got %f, want 0.5", hslWithAlpha.Alpha())
	}

	// HSV WithAlpha and Alpha
	hsv := NewHSV(180, 0.5, 0.5, 1.0)
	if hsv.Alpha() != 1.0 {
		t.Errorf("HSV Alpha failed: got %f, want 1.0", hsv.Alpha())
	}
	hsvWithAlpha := hsv.WithAlpha(0.7)
	if hsvWithAlpha.Alpha() != 0.7 {
		t.Errorf("HSV WithAlpha failed: got %f, want 0.7", hsvWithAlpha.Alpha())
	}

	// HWB WithAlpha and Alpha
	hwb := NewHWB(180, 0.2, 0.2, 1.0)
	if hwb.Alpha() != 1.0 {
		t.Errorf("HWB Alpha failed: got %f, want 1.0", hwb.Alpha())
	}
	hwbWithAlpha := hwb.WithAlpha(0.3)
	if hwbWithAlpha.Alpha() != 0.3 {
		t.Errorf("HWB WithAlpha failed: got %f, want 0.3", hwbWithAlpha.Alpha())
	}

	// LAB WithAlpha
	lab := NewLAB(50, 20, 30, 1.0)
	labWithAlpha := lab.WithAlpha(0.8)
	if labWithAlpha.Alpha() != 0.8 {
		t.Errorf("LAB WithAlpha failed: got %f, want 0.8", labWithAlpha.Alpha())
	}

	// LCH WithAlpha
	lch := NewLCH(70, 50, 180, 1.0)
	lchWithAlpha := lch.WithAlpha(0.6)
	if lchWithAlpha.Alpha() != 0.6 {
		t.Errorf("LCH WithAlpha failed: got %f, want 0.6", lchWithAlpha.Alpha())
	}

	// OKLAB WithAlpha
	oklab := NewOKLAB(0.7, 0.1, -0.1, 1.0)
	oklabWithAlpha := oklab.WithAlpha(0.4)
	if oklabWithAlpha.Alpha() != 0.4 {
		t.Errorf("OKLAB WithAlpha failed: got %f, want 0.4", oklabWithAlpha.Alpha())
	}

	// OKLCH WithAlpha
	oklch := NewOKLCH(0.7, 0.2, 180, 1.0)
	oklchWithAlpha := oklch.WithAlpha(0.9)
	if oklchWithAlpha.Alpha() != 0.9 {
		t.Errorf("OKLCH WithAlpha failed: got %f, want 0.9", oklchWithAlpha.Alpha())
	}
}

// Test ToHWB conversion
func TestToHWB(t *testing.T) {
	colors := []Color{
		RGB(1, 0, 0),     // Red
		RGB(0, 1, 0),     // Green
		RGB(0, 0, 1),     // Blue
		RGB(0.5, 0.5, 0.5), // Gray
		RGB(1, 1, 1),     // White
		RGB(0, 0, 0),     // Black
	}

	for _, c := range colors {
		hwb := ToHWB(c)
		if hwb == nil {
			t.Error("ToHWB returned nil")
		}

		// Convert back to RGB and check it's similar
		r1, g1, b1, _ := c.RGBA()
		r2, g2, b2, _ := hwb.RGBA()

		tolerance := 0.01
		if abs(r1-r2) > tolerance || abs(g1-g2) > tolerance || abs(b1-b2) > tolerance {
			t.Errorf("ToHWB round trip failed: (%f,%f,%f) -> (%f,%f,%f)",
				r1, g1, b1, r2, g2, b2)
		}
	}
}

// Test ClipToGamut function
func TestClipToGamut(t *testing.T) {
	// Out of gamut colors
	outOfGamut := []Color{
		NewRGBA(1.5, 0.5, 0.3, 1.0),
		NewRGBA(0.5, -0.2, 0.8, 1.0),
		NewRGBA(0.2, 0.8, 1.2, 1.0),
	}

	for _, c := range outOfGamut {
		clipped := ClipToGamut(c)
		if !InGamut(clipped) {
			t.Error("ClipToGamut produced out-of-gamut color")
		}

		r, g, b, a := clipped.RGBA()
		if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
			t.Errorf("ClipToGamut produced invalid RGB: (%f, %f, %f, %f)", r, g, b, a)
		}
	}

	// In gamut colors should not change
	inGamut := RGB(0.5, 0.3, 0.8)
	clipped := ClipToGamut(inGamut)
	r1, g1, b1, _ := inGamut.RGBA()
	r2, g2, b2, _ := clipped.RGBA()

	if abs(r1-r2) > 0.001 || abs(g1-g2) > 0.001 || abs(b1-b2) > 0.001 {
		t.Error("ClipToGamut changed in-gamut color")
	}
}

// Test internal gamut mapping functions
func TestGamutMappingInternals(t *testing.T) {
	// Create out-of-gamut color
	vivid := NewOKLCH(0.7, 0.35, 150, 1.0)

	// Test mapClip
	clipped := mapClip(vivid)
	if !InGamut(clipped) {
		t.Error("mapClip produced out-of-gamut color")
	}

	// Test mapPreserveLightness
	lightness := mapPreserveLightness(vivid)
	if !InGamut(lightness) {
		t.Error("mapPreserveLightness produced out-of-gamut color")
	}

	// Test mapPreserveChroma
	chroma := mapPreserveChroma(vivid)
	if !InGamut(chroma) {
		t.Error("mapPreserveChroma produced out-of-gamut color")
	}

	// Test mapProject
	project := mapProject(vivid)
	if !InGamut(project) {
		t.Error("mapProject produced out-of-gamut color")
	}
}

// Test Error() method
func TestParseError(t *testing.T) {
	err := &ParseError{input: "invalid", reason: "test reason"}
	errStr := err.Error()
	if errStr == "" {
		t.Error("ParseError.Error() returned empty string")
	}
	if errStr != `cannot parse color "invalid": test reason` {
		t.Errorf("ParseError.Error() returned unexpected string: %q", errStr)
	}
}

// Test formatString helper
func TestFormatString(t *testing.T) {
	result := formatString("test %d %s", 42, "hello")
	expected := "test 42 hello"
	if result != expected {
		t.Errorf("formatString failed: got %q, want %q", result, expected)
	}
}

// Test Lipgloss integration
func TestLipglossIntegration(t *testing.T) {
	c := RGB(1, 0.5, 0)

	// Test ToLipglossColor
	lg := ToLipglossColor(c)
	if lg == "" {
		t.Error("ToLipglossColor returned empty string")
	}

	// Test ToLipglossColorWithAlpha
	cAlpha := NewRGBA(1, 0.5, 0, 0.5)
	lgAlpha := ToLipglossColorWithAlpha(cAlpha)
	if lgAlpha == "" {
		t.Error("ToLipglossColorWithAlpha returned empty string")
	}
}
