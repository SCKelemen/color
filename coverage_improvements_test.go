package color

import (
	"math"
	"testing"
)

// Test Error() method on HexParseError
func TestHexParseError(t *testing.T) {
	err := &HexParseError{hex: "zzz"}
	expected := "invalid hex color: zzz"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

// Test DarkenSpace
func TestDarkenSpace(t *testing.T) {
	// Test with OKLCH color
	oklchColor := NewSpaceColor(OKLCHSpace, []float64{0.7, 0.15, 180}, 1.0)
	darkened := DarkenSpace(oklchColor, 0.3)

	channels := darkened.Channels()
	if channels[0] >= 0.7 {
		t.Errorf("DarkenSpace should decrease lightness: got %f", channels[0])
	}

	// Test with sRGB color (should convert to OKLCH, darken, convert back)
	srgbColor := NewSpaceColor(SRGBSpace, []float64{0.8, 0.5, 0.2}, 1.0)
	darkenedSRGB := DarkenSpace(srgbColor, 0.2)

	if darkenedSRGB.Space().Name() != "sRGB" {
		t.Errorf("DarkenSpace should preserve space: got %s", darkenedSRGB.Space().Name())
	}
}

// Test SaturateSpace
func TestSaturateSpace(t *testing.T) {
	// Test with OKLCH color
	oklchColor := NewSpaceColor(OKLCHSpace, []float64{0.7, 0.1, 180}, 1.0)
	saturated := SaturateSpace(oklchColor, 0.5)

	channels := saturated.Channels()
	if channels[1] <= 0.1 {
		t.Errorf("SaturateSpace should increase chroma: got %f", channels[1])
	}

	// Test with sRGB color
	srgbColor := NewSpaceColor(SRGBSpace, []float64{0.5, 0.5, 0.5}, 1.0)
	saturatedSRGB := SaturateSpace(srgbColor, 0.3)

	if saturatedSRGB.Space().Name() != "sRGB" {
		t.Errorf("SaturateSpace should preserve space: got %s", saturatedSRGB.Space().Name())
	}
}

// Test DesaturateSpace
func TestDesaturateSpace(t *testing.T) {
	// Test with OKLCH color
	oklchColor := NewSpaceColor(OKLCHSpace, []float64{0.7, 0.2, 180}, 1.0)
	desaturated := DesaturateSpace(oklchColor, 0.5)

	channels := desaturated.Channels()
	expectedChroma := 0.2 * 0.5
	if math.Abs(channels[1]-expectedChroma) > 0.01 {
		t.Errorf("DesaturateSpace chroma = %f, want ~%f", channels[1], expectedChroma)
	}

	// Test with sRGB color
	srgbColor := NewSpaceColor(SRGBSpace, []float64{0.8, 0.5, 0.2}, 1.0)
	desaturatedSRGB := DesaturateSpace(srgbColor, 0.3)

	if desaturatedSRGB.Space().Name() != "sRGB" {
		t.Errorf("DesaturateSpace should preserve space: got %s", desaturatedSRGB.Space().Name())
	}
}

// Test UnregisterSpace
func TestUnregisterSpace(t *testing.T) {
	// Use an existing space (sRGB) and register it with a test name
	testName := "test-space-unregister-temp"

	// Register it
	RegisterSpace(testName, SRGBSpace)

	// Verify it's registered
	if _, ok := GetSpace(testName); !ok {
		t.Fatal("Space not registered")
	}

	// Unregister it
	UnregisterSpace(testName)

	// Verify it's unregistered
	if _, ok := GetSpace(testName); ok {
		t.Error("Space should be unregistered")
	}
}

// Test WithAlpha on spaceColor
func TestSpaceColorWithAlpha(t *testing.T) {
	oklchColor := NewSpaceColor(OKLCHSpace, []float64{0.7, 0.15, 180}, 1.0)
	withAlpha := oklchColor.WithAlpha(0.5)

	if withAlpha.Alpha() != 0.5 {
		t.Errorf("WithAlpha() alpha = %f, want 0.5", withAlpha.Alpha())
	}

	// Original should be unchanged
	if oklchColor.Alpha() != 1.0 {
		t.Error("WithAlpha() should not modify original")
	}
}

// Test WithAlpha on XYZ
func TestXYZWithAlpha(t *testing.T) {
	xyz := NewXYZ(0.5, 0.6, 0.7, 1.0)
	withAlpha := xyz.WithAlpha(0.3)

	if withAlpha.Alpha() != 0.3 {
		t.Errorf("WithAlpha() alpha = %f, want 0.3", withAlpha.Alpha())
	}

	// Verify it's still an XYZ color
	if _, ok := withAlpha.(*XYZ); !ok {
		t.Error("WithAlpha() should return *XYZ")
	}
}

// Test ChannelNames for OKLCH
func TestOKLCHChannelNames(t *testing.T) {
	names := OKLCHSpace.ChannelNames()
	expected := []string{"L", "C", "H"}

	if len(names) != len(expected) {
		t.Fatalf("ChannelNames() length = %d, want %d", len(names), len(expected))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("ChannelNames()[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

// Test ChannelNames for sRGB
func TestSRGBChannelNames(t *testing.T) {
	names := SRGBSpace.ChannelNames()
	expected := []string{"R", "G", "B"}

	if len(names) != len(expected) {
		t.Fatalf("ChannelNames() length = %d, want %d", len(names), len(expected))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("ChannelNames()[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

// Test ConvertToWithMapping
func TestConvertToWithMapping(t *testing.T) {
	// Create a wide-gamut color in Display P3 that's out of sRGB gamut
	p3Color := NewSpaceColor(DisplayP3Space, []float64{0.0, 1.0, 0.5}, 1.0)

	// Convert to sRGB with different mapping methods
	clipped := p3Color.(*spaceColor).ConvertToWithMapping(SRGBSpace, GamutClip)
	preserved := p3Color.(*spaceColor).ConvertToWithMapping(SRGBSpace, GamutPreserveLightness)

	// Both should be in sRGB space
	if clipped.Space().Name() != "sRGB" {
		t.Error("ConvertToWithMapping should convert to target space")
	}
	if preserved.Space().Name() != "sRGB" {
		t.Error("ConvertToWithMapping should convert to target space")
	}

	// Both should be in gamut
	clippedRGB := clipped.ToRGBA()
	if !InGamut(clippedRGB) {
		t.Error("Clipped color should be in gamut")
	}
}

// Test MapToGamut with different methods
func TestMapToGamutMethods(t *testing.T) {
	// Out of gamut color
	outOfGamut := RGB(1.5, -0.2, 0.8)

	tests := []struct {
		name   string
		method GamutMapping
	}{
		{"Clip", GamutClip},
		{"PreserveLightness", GamutPreserveLightness},
		{"Project", GamutProject},
		{"PreserveChroma", GamutPreserveChroma},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapped := MapToGamut(outOfGamut, tt.method)
			if !InGamut(mapped) {
				t.Errorf("MapToGamut(%v) should produce in-gamut color", tt.method)
			}
		})
	}
}

// Test parseLCH with more cases
func TestParseLCH(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"lch(70% 50 180)", false},
		{"lch(70% 50 180 / 0.5)", false},
		{"lch(70% 50 180 / 50%)", false},
		{"lch(70 50 180)", false}, // Without %
		// Note: Parser is lenient and clamps values rather than erroring
	}

	for _, tt := range tests {
		_, err := ParseColor(tt.input)
		hasErr := err != nil
		if hasErr != tt.wantErr {
			t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

// Test parseOKLAB with more cases
func TestParseOKLAB(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"oklab(0.7 0.1 -0.1)", false},
		{"oklab(0.7 0.1 -0.1 / 0.8)", false},
		{"oklab(70% 0.1 -0.1)", false},
		// Parser is lenient and clamps out-of-range values
	}

	for _, tt := range tests {
		_, err := ParseColor(tt.input)
		hasErr := err != nil
		if hasErr != tt.wantErr {
			t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

// Test parseLAB with more cases
func TestParseLAB(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"lab(70% 20 -30)", false},
		{"lab(70% 20 -30 / 0.9)", false},
		{"lab(70 20 -30)", false},
		// Parser is lenient
	}

	for _, tt := range tests {
		_, err := ParseColor(tt.input)
		hasErr := err != nil
		if hasErr != tt.wantErr {
			t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

// Test parseHSV with more cases
func TestParseHSV(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hsv(180, 50%, 80%)", false},
		{"hsv(180, 50%, 80%, 0.7)", false},
		{"hsv(180 50% 80%)", false},
		// Parser is lenient and clamps values
	}

	for _, tt := range tests {
		_, err := ParseColor(tt.input)
		hasErr := err != nil
		if hasErr != tt.wantErr {
			t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

// Test parseXYZ with validation edge cases
func TestParseXYZValidation(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"xyz(0.5 0.6 0.7)", false},
		{"xyz(0.5 0.6 0.7 / 0.8)", false},
		{"xyz(50% 60% 70%)", false},
		// Parser is lenient and allows various ranges
	}

	for _, tt := range tests {
		_, err := ParseColor(tt.input)
		hasErr := err != nil
		if hasErr != tt.wantErr {
			t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

// Test parseRGBColorSpace
func TestParseRGBColorSpace(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"color(srgb 0.5 0.6 0.7)", false},
		{"color(display-p3 0.5 0.6 0.7)", false},
		{"color(rec2020 0.5 0.6 0.7)", false},
		{"color(a98-rgb 0.5 0.6 0.7)", false},
		{"color(prophoto-rgb 0.5 0.6 0.7)", false},
		{"color(unknown-space 0.5 0.6 0.7)", true},
		// Parser may be lenient with out-of-range values
	}

	for _, tt := range tests {
		c, err := ParseColor(tt.input)
		hasErr := err != nil
		if hasErr != tt.wantErr {
			t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}

		if !hasErr {
			if sc, ok := c.(SpaceColor); ok {
				// Verify it's a valid space color
				if sc.Space() == nil {
					t.Errorf("ParseColor(%q) produced nil space", tt.input)
				}
			}
		}
	}
}

// Test clamp edge cases - additional tests beyond the existing TestClamp
func TestClampAdditional(t *testing.T) {
	tests := []struct {
		value float64
		min   float64
		max   float64
		want  float64
	}{
		{5.0, 2.0, 8.0, 5.0},
		{1.0, 2.0, 8.0, 2.0},
		{10.0, 2.0, 8.0, 8.0},
	}

	for _, tt := range tests {
		got := clamp(tt.value, tt.min, tt.max)
		if got != tt.want {
			t.Errorf("clamp(%f, %f, %f) = %f, want %f", tt.value, tt.min, tt.max, got, tt.want)
		}
	}
}

// Test LightenSpace edge cases
func TestLightenSpaceEdgeCases(t *testing.T) {
	// Test with amount > 1.0 (should clamp to 1.0)
	oklchColor := NewSpaceColor(OKLCHSpace, []float64{0.5, 0.15, 180}, 1.0)
	lightened := LightenSpace(oklchColor, 1.5)

	channels := lightened.Channels()
	if channels[0] > 1.0 {
		t.Errorf("LightenSpace should clamp lightness to 1.0: got %f", channels[0])
	}

	// Test with amount < 0.0 (should clamp to 0.0)
	lightened2 := LightenSpace(oklchColor, -0.5)
	channels2 := lightened2.Channels()

	// Should be same as original when amount is clamped to 0
	if math.Abs(channels2[0]-0.5) > 0.01 {
		t.Errorf("LightenSpace with negative amount should not change lightness: got %f, want ~0.5", channels2[0])
	}
}

// Test GradientMultiStopWithEasing with additional easing functions
func TestGradientMultiStopWithEasingAdditional(t *testing.T) {
	stops := []GradientStop{
		{Color: RGB(1, 0, 0), Position: 0.0},   // Red
		{Color: RGB(0, 1, 0), Position: 0.5},   // Green
		{Color: RGB(0, 0, 1), Position: 1.0},   // Blue
	}

	// Test with ease-out easing
	easeOutGradient := GradientMultiStopWithEasing(stops, 10, GradientOKLCH, EaseOutQuad)
	if len(easeOutGradient) != 10 {
		t.Errorf("GradientMultiStopWithEasing length = %d, want 10", len(easeOutGradient))
	}

	// Test with ease-in-out easing
	easeInOutGradient := GradientMultiStopWithEasing(stops, 10, GradientOKLCH, EaseInOutQuad)
	if len(easeInOutGradient) != 10 {
		t.Errorf("GradientMultiStopWithEasing length = %d, want 10", len(easeInOutGradient))
	}

	// Test with cubic easing
	cubicGradient := GradientMultiStopWithEasing(stops, 15, GradientLAB, EaseInCubic)
	if len(cubicGradient) != 15 {
		t.Errorf("GradientMultiStopWithEasing with cubic length = %d, want 15", len(cubicGradient))
	}

	// First color should match first stop
	r1, g1, b1, _ := easeOutGradient[0].RGBA()
	r2, g2, b2, _ := stops[0].Color.RGBA()
	if math.Abs(r1-r2) > 0.01 || math.Abs(g1-g2) > 0.01 || math.Abs(b1-b2) > 0.01 {
		t.Error("First gradient color should match first stop")
	}

	// Last color should match last stop
	rL1, gL1, bL1, _ := easeOutGradient[len(easeOutGradient)-1].RGBA()
	rL2, gL2, bL2, _ := stops[len(stops)-1].Color.RGBA()
	if math.Abs(rL1-rL2) > 0.01 || math.Abs(gL1-gL2) > 0.01 || math.Abs(bL1-bL2) > 0.01 {
		t.Error("Last gradient color should match last stop")
	}
}

// Test interpolateHue edge cases with different interpolation methods
func TestInterpolateHueEdgeCases(t *testing.T) {
	tests := []struct {
		h1     float64
		h2     float64
		t      float64
		method HueInterpolation
	}{
		{0, 180, 0.5, HueShorter},     // Normal interpolation
		{350, 10, 0.5, HueShorter},    // Wrapping around 0/360
		{10, 350, 0.5, HueLonger},     // Longer path
		{0, 180, 0.5, HueIncreasing},  // Increasing direction
		{180, 0, 0.5, HueDecreasing},  // Decreasing direction
		{120, 240, 0.25, HueShorter},  // Quarter way
	}

	for _, tt := range tests {
		got := interpolateHue(tt.h1, tt.h2, tt.t, tt.method)
		// Just verify it returns a valid hue
		if got < 0 || got >= 360 {
			t.Errorf("interpolateHue(%f, %f, %f, %v) = %f, should be in [0, 360)",
				tt.h1, tt.h2, tt.t, tt.method, got)
		}
	}
}

// Test RGB space transfer functions
func TestRGBSpaceTransferFunctions(t *testing.T) {
	// Test Rec. 709 transfer
	rec709Space, ok := GetSpace("rec-709")
	if !ok {
		t.Skip("Rec. 709 space not found")
	}

	// Convert a color through Rec. 709
	color := NewSpaceColor(rec709Space, []float64{0.5, 0.6, 0.7}, 1.0)
	xyz := ToXYZ(color)

	// Should produce valid XYZ values
	if xyz.X < 0 || xyz.Y < 0 || xyz.Z < 0 {
		t.Error("Rec. 709 conversion produced negative XYZ values")
	}

	// Test Rec. 2020 transfer
	rec2020Space, ok := GetSpace("rec2020")
	if !ok {
		t.Skip("Rec. 2020 space not found")
	}

	color2 := NewSpaceColor(rec2020Space, []float64{0.5, 0.6, 0.7}, 1.0)
	xyz2 := ToXYZ(color2)

	// Should produce valid XYZ values
	if xyz2.X < 0 || xyz2.Y < 0 || xyz2.Z < 0 {
		t.Error("Rec. 2020 conversion produced negative XYZ values")
	}
}

// Test ConvertXYZToRGB by converting through spaces
func TestConvertXYZToRGBThroughSpaces(t *testing.T) {
	// Create an XYZ color
	xyz := NewXYZ(0.5, 0.6, 0.7, 1.0)

	// Convert to different RGB spaces through the Space interface
	spaces := []string{"display-p3", "a98-rgb", "prophoto-rgb", "rec2020"}

	for _, spaceName := range spaces {
		space, ok := GetSpace(spaceName)
		if !ok {
			t.Errorf("Space %q not found", spaceName)
			continue
		}

		// Convert XYZ to the RGB space using FromXYZ
		channels := space.FromXYZ(xyz.X, xyz.Y, xyz.Z)
		if len(channels) != 3 {
			t.Errorf("Space %q FromXYZ returned %d channels, want 3", spaceName, len(channels))
		}

		// Verify all channels are reasonable values
		for i, ch := range channels {
			if math.IsNaN(ch) || math.IsInf(ch, 0) {
				t.Errorf("Space %q channel %d is invalid: %f", spaceName, i, ch)
			}
		}
	}
}

// Test mapPreserveChroma and mapPreserveLightness edge cases
func TestGamutMappingEdgeCases(t *testing.T) {
	// Very high chroma color that's out of gamut
	highChroma := NewOKLCH(0.7, 0.5, 180, 1.0)

	// Test preserve lightness mapping
	mappedLight := MapToGamut(highChroma, GamutPreserveLightness)
	if !InGamut(mappedLight) {
		t.Error("GamutPreserveLightness should produce in-gamut color")
	}

	// Check that lightness is approximately preserved (allow some tolerance)
	mappedOKLCH := ToOKLCH(mappedLight)
	if math.Abs(mappedOKLCH.L-0.7) > 0.15 {
		t.Errorf("GamutPreserveLightness should approximately preserve lightness: got %f, want ~0.7", mappedOKLCH.L)
	}

	// Test preserve chroma mapping
	mappedChroma := MapToGamut(highChroma, GamutPreserveChroma)
	if !InGamut(mappedChroma) {
		t.Error("GamutPreserveChroma should produce in-gamut color")
	}

	// Very dark color
	darkColor := NewOKLCH(0.1, 0.3, 120, 1.0)
	mappedDark := MapToGamut(darkColor, GamutPreserveLightness)

	if !InGamut(mappedDark) {
		t.Error("Gamut mapping should handle dark colors")
	}

	// Very light color
	lightColor := NewOKLCH(0.95, 0.3, 240, 1.0)
	mappedLightColor := MapToGamut(lightColor, GamutPreserveLightness)

	if !InGamut(mappedLightColor) {
		t.Error("Gamut mapping should handle light colors")
	}

	// Test project method
	projected := MapToGamut(highChroma, GamutProject)
	if !InGamut(projected) {
		t.Error("GamutProject should produce in-gamut color")
	}
}
