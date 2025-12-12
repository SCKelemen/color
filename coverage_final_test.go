package color

import (
	stdcolor "image/color"
	"math"
	"testing"
)

// Test linear transfer functions
func TestLinearTransfer(t *testing.T) {
	tests := []float64{0.0, 0.5, 1.0}
	for _, v := range tests {
		result := linearTransfer(v)
		if result != v {
			t.Errorf("linearTransfer(%f) = %f, want %f", v, result, v)
		}

		inv := linearInverseTransfer(v)
		if inv != v {
			t.Errorf("linearInverseTransfer(%f) = %f, want %f", v, inv, v)
		}
	}
}

// Test rec709 transfer functions
func TestRec709Transfer(t *testing.T) {
	tests := []float64{0.0, 0.04, 0.5, 1.0}
	for _, v := range tests {
		encoded := rec709Transfer(v)
		decoded := rec709InverseTransfer(encoded)

		// Round-trip should be close
		if math.Abs(decoded-v) > 0.001 {
			t.Errorf("Rec709 round-trip failed: %f -> %f -> %f", v, encoded, decoded)
		}
	}
}

// Test rec2020 transfer functions
func TestRec2020Transfer(t *testing.T) {
	tests := []float64{0.0, 0.04, 0.5, 1.0}
	for _, v := range tests {
		encoded := rec2020Transfer(v)
		decoded := rec2020InverseTransfer(encoded)

		// Round-trip should be close
		if math.Abs(decoded-v) > 0.001 {
			t.Errorf("Rec2020 round-trip failed: %f -> %f -> %f", v, encoded, decoded)
		}
	}

	// Test negative value handling
	if rec2020Transfer(-0.1) < 0 {
		t.Error("rec2020Transfer should handle negative values")
	}
}

// Test gamma transfer functions
func TestGammaTransferFunctions(t *testing.T) {
	gamma := 2.2
	transfer := gammaTransferFunc(gamma)
	inverse := gammaInverseTransferFunc(gamma)

	tests := []float64{0.0, 0.25, 0.5, 0.75, 1.0}
	for _, v := range tests {
		encoded := transfer(v)
		decoded := inverse(encoded)

		// Round-trip should be close
		if math.Abs(decoded-v) > 0.001 {
			t.Errorf("Gamma round-trip failed: %f -> %f -> %f", v, encoded, decoded)
		}
	}

	// Test negative value handling
	if transfer(-0.1) != 0 {
		t.Error("gammaTransferFunc should return 0 for negative values")
	}
	if inverse(-0.1) != 0 {
		t.Error("gammaInverseTransferFunc should return 0 for negative values")
	}
}

// Test ConvertXYZToRGB directly
func TestConvertXYZToRGBDirect(t *testing.T) {
	xyz := NewXYZ(0.5, 0.6, 0.7, 1.0)

	// Get Display P3 as an RGBColorSpace
	space, ok := GetSpace("display-p3")
	if !ok {
		t.Skip("Display P3 not available")
	}

	// Convert through the space's FromXYZ method
	channels := space.FromXYZ(xyz.X, xyz.Y, xyz.Z)
	if len(channels) != 3 {
		t.Errorf("FromXYZ returned %d channels, want 3", len(channels))
	}

	// Test with ProPhoto RGB (wider gamut)
	prophoto, ok := GetSpace("prophoto-rgb")
	if !ok {
		t.Skip("ProPhoto RGB not available")
	}

	channelsP := prophoto.FromXYZ(xyz.X, xyz.Y, xyz.Z)
	if len(channelsP) != 3 {
		t.Errorf("FromXYZ returned %d channels, want 3", len(channelsP))
	}
}

// Test isHexString with various inputs
func TestIsHexString(t *testing.T) {
	validHex := []string{
		"#fff",
		"#ffffff",
		"#ffff",
		"#ffffffff",
		"fff",
		"ffffff",
	}

	for _, hex := range validHex {
		if _, err := HexToRGB(hex); err != nil {
			t.Errorf("HexToRGB(%q) should succeed, got error: %v", hex, err)
		}
	}

	invalidHex := []string{
		"#ff",     // Too short
		"#fffff",  // Wrong length
		"#ggg",    // Invalid characters
		"zzzzzz",  // Invalid characters
	}

	for _, hex := range invalidHex {
		if _, err := HexToRGB(hex); err == nil {
			t.Errorf("HexToRGB(%q) should fail, but succeeded", hex)
		}
	}
}

// Test MapToGamut default case
func TestMapToGamutInvalidMethod(t *testing.T) {
	outOfGamut := RGB(1.5, -0.2, 0.8)

	// Use an invalid method value (should default to clip)
	mapped := MapToGamut(outOfGamut, GamutMapping(999))
	if !InGamut(mapped) {
		t.Error("MapToGamut with invalid method should still produce in-gamut color")
	}
}

// Test LightenSpace with different spaces
func TestLightenSpaceNonOKLCH(t *testing.T) {
	// Test with sRGB space
	srgbColor := NewSpaceColor(SRGBSpace, []float64{0.5, 0.3, 0.2}, 1.0)
	lightened := LightenSpace(srgbColor, 0.2)

	if lightened.Space().Name() != "sRGB" {
		t.Error("LightenSpace should preserve space")
	}

	// Test with Display P3
	p3Color := NewSpaceColor(DisplayP3Space, []float64{0.4, 0.3, 0.2}, 1.0)
	lightenedP3 := LightenSpace(p3Color, 0.3)

	if lightenedP3.Space().Name() != "display-p3" {
		t.Errorf("LightenSpace should preserve Display P3 space, got %s", lightenedP3.Space().Name())
	}
}

// Test DarkenSpace with non-OKLCH spaces
func TestDarkenSpaceNonOKLCH(t *testing.T) {
	// Test with Display P3
	p3Color := NewSpaceColor(DisplayP3Space, []float64{0.8, 0.7, 0.6}, 1.0)
	darkened := DarkenSpace(p3Color, 0.3)

	if darkened.Space().Name() != "display-p3" {
		t.Errorf("DarkenSpace should preserve Display P3 space, got %s", darkened.Space().Name())
	}
}

// Test SaturateSpace with non-OKLCH spaces
func TestSaturateSpaceNonOKLCH(t *testing.T) {
	// Test with Display P3
	p3Color := NewSpaceColor(DisplayP3Space, []float64{0.5, 0.5, 0.5}, 1.0)
	saturated := SaturateSpace(p3Color, 0.4)

	if saturated.Space().Name() != "display-p3" {
		t.Errorf("SaturateSpace should preserve Display P3 space, got %s", saturated.Space().Name())
	}
}

// Test DesaturateSpace with non-OKLCH spaces
func TestDesaturateSpaceNonOKLCH(t *testing.T) {
	// Test with Display P3
	p3Color := NewSpaceColor(DisplayP3Space, []float64{0.8, 0.3, 0.2}, 1.0)
	desaturated := DesaturateSpace(p3Color, 0.5)

	if desaturated.Space().Name() != "display-p3" {
		t.Errorf("DesaturateSpace should preserve Display P3 space, got %s", desaturated.Space().Name())
	}
}

// Test interpolateHue with all methods
func TestInterpolateHueAllMethods(t *testing.T) {
	methods := []HueInterpolation{
		HueShorter,
		HueLonger,
		HueIncreasing,
		HueDecreasing,
	}

	for _, method := range methods {
		result := interpolateHue(10, 350, 0.5, method)
		if result < 0 || result >= 360 {
			t.Errorf("interpolateHue with method %v returned invalid hue: %f", method, result)
		}

		// Test wrap-around cases
		result2 := interpolateHue(350, 10, 0.5, method)
		if result2 < 0 || result2 >= 360 {
			t.Errorf("interpolateHue wrap-around with method %v returned invalid hue: %f", method, result2)
		}
	}
}

// Test mapPreserveChroma with different scenarios
func TestMapPreserveChromaScenarios(t *testing.T) {
	// Test with a color that needs to be lightened
	darkOutOfGamut := NewOKLCH(0.2, 0.4, 120, 1.0)
	mapped := MapToGamut(darkOutOfGamut, GamutPreserveChroma)

	if !InGamut(mapped) {
		t.Error("mapPreserveChroma should produce in-gamut color for dark colors")
	}

	// Test with mid-range color
	midOutOfGamut := NewOKLCH(0.5, 0.5, 240, 1.0)
	mappedMid := MapToGamut(midOutOfGamut, GamutPreserveChroma)

	if !InGamut(mappedMid) {
		t.Error("mapPreserveChroma should produce in-gamut color for mid-range colors")
	}
}

// Test GradientWithEasing
func TestGradientWithEasingFunction(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Test with different easing functions
	easingFuncs := []EasingFunction{
		EaseLinear,
		EaseInQuad,
		EaseOutQuad,
		EaseInOutQuad,
		EaseInCubic,
		EaseOutCubic,
		EaseInOutCubic,
	}

	for _, easing := range easingFuncs {
		gradient := GradientWithEasing(red, blue, 10, GradientOKLCH, easing)
		if len(gradient) != 10 {
			t.Errorf("GradientWithEasing length = %d, want 10", len(gradient))
		}

		// First should be red, last should be blue
		r1, _, _, _ := gradient[0].RGBA()
		if math.Abs(r1-1.0) > 0.01 {
			t.Error("First color should be red")
		}

		_, _, b2, _ := gradient[len(gradient)-1].RGBA()
		if math.Abs(b2-1.0) > 0.01 {
			t.Error("Last color should be blue")
		}
	}
}

// Test HWB RGBA method
func TestHWBRGBA(t *testing.T) {
	// Create an HWB color
	hwb := NewHWB(180, 0.3, 0.2, 1.0)
	r, g, b, a := hwb.RGBA()

	// Verify it's in valid range
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
		t.Errorf("RGBA values out of range: (%f, %f, %f, %f)", r, g, b, a)
	}

	// Test with different hue
	hwb2 := NewHWB(240, 0.1, 0.1, 0.8)
	r2, g2, b2, a2 := hwb2.RGBA()

	if r2 < 0 || r2 > 1 || g2 < 0 || g2 > 1 || b2 < 0 || b2 > 1 {
		t.Errorf("RGBA values out of range: (%f, %f, %f)", r2, g2, b2)
	}

	if a2 != 0.8 {
		t.Errorf("Alpha = %f, want 0.8", a2)
	}
}

// Test FromStdColor with paletted image color
func TestFromStdColorPaletted(t *testing.T) {
	// Create a standard library color
	stdColor := stdcolor.RGBA{R: 128, G: 64, B: 192, A: 255}

	// Convert to our Color type
	c := FromStdColor(stdColor)

	r, g, b, a := c.RGBA()

	// Should be approximately 128/255, 64/255, 192/255, 1.0
	expectedR := 128.0 / 255.0
	expectedG := 64.0 / 255.0
	expectedB := 192.0 / 255.0

	if math.Abs(r-expectedR) > 0.01 || math.Abs(g-expectedG) > 0.01 || math.Abs(b-expectedB) > 0.01 {
		t.Errorf("FromStdColor conversion incorrect: got (%f, %f, %f), want (%f, %f, %f)",
			r, g, b, expectedR, expectedG, expectedB)
	}

	if a != 1.0 {
		t.Errorf("Alpha = %f, want 1.0", a)
	}
}

// Test NewSpaceColor with invalid channel count (should panic)
func TestNewSpaceColorPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewSpaceColor should panic with wrong channel count")
		}
	}()

	// This should panic because we're providing 2 channels instead of 3
	NewSpaceColor(SRGBSpace, []float64{0.5, 0.6}, 1.0)
}

// Test ConvertFromRGBSpace with different spaces
func TestConvertFromRGBSpaceVariety(t *testing.T) {
	// Convert from different RGB spaces
	spaces := []string{"srgb", "display-p3", "a98-rgb", "prophoto-rgb"}

	for _, spaceName := range spaces {
		// Convert from the space (assumes the color 0.5, 0.6, 0.7 was in that space)
		converted, err := ConvertFromRGBSpace(0.5, 0.6, 0.7, 1.0, spaceName)
		if err != nil {
			t.Errorf("ConvertFromRGBSpace(%q) error: %v", spaceName, err)
			continue
		}

		// Should be valid
		r, g, b, _ := converted.RGBA()
		if math.IsNaN(r) || math.IsNaN(g) || math.IsNaN(b) {
			t.Errorf("ConvertFromRGBSpace(%q) produced NaN", spaceName)
		}
	}

	// Test with invalid space name (should error)
	_, err := ConvertFromRGBSpace(0.5, 0.6, 0.7, 1.0, "invalid-space")
	if err == nil {
		t.Error("ConvertFromRGBSpace with invalid space should return error")
	}
}
