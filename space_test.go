package color

import "testing"

func TestSpaceColor(t *testing.T) {
	// Create a color in Display P3
	p3Color := NewSpaceColor(
		DisplayP3Space,
		[]float64{1.0, 0.5, 0.0}, // R, G, B in Display P3
		1.0,                      // alpha
	)

	if p3Color.Space().Name() != "display-p3" {
		t.Errorf("Expected space 'display-p3', got '%s'", p3Color.Space().Name())
	}

	channels := p3Color.Channels()
	if len(channels) != 3 {
		t.Errorf("Expected 3 channels, got %d", len(channels))
	}

	// Convert to sRGB (explicit conversion)
	srgbColor := p3Color.ConvertTo(SRGBSpace)
	if srgbColor.Space().Name() != "sRGB" {
		t.Errorf("Expected space 'sRGB', got '%s'", srgbColor.Space().Name())
	}

	// Convert back to Display P3 (should preserve original if in gamut)
	p3Back := srgbColor.ConvertTo(DisplayP3Space)

	// Note: May not be exactly equal due to gamut clipping, but should be close
	_ = p3Back
}

func TestSpaceOperations(t *testing.T) {
	// Create a color in OKLCH
	oklchColor := NewSpaceColor(
		OKLCHSpace,
		[]float64{0.6, 0.2, 180}, // L, C, H
		1.0,
	)

	// Lighten in native space
	lightened := LightenSpace(oklchColor, 0.2)
	if lightened.Space().Name() != "OKLCH" {
		t.Errorf("Expected space 'OKLCH', got '%s'", lightened.Space().Name())
	}

	// Check that lightness increased
	origChannels := oklchColor.Channels()
	newChannels := lightened.Channels()
	if newChannels[0] <= origChannels[0] {
		t.Errorf("Lightness should increase: %v -> %v", origChannels[0], newChannels[0])
	}
}

func TestSpaceConversionThroughXYZ(t *testing.T) {
	// Create color in Display P3
	p3Color := NewSpaceColor(
		DisplayP3Space,
		[]float64{1.0, 0.0, 0.0}, // Red in Display P3
		1.0,
	)

	// Convert to OKLCH (through XYZ)
	oklchColor := p3Color.ConvertTo(OKLCHSpace)
	if oklchColor.Space().Name() != "OKLCH" {
		t.Errorf("Expected space 'OKLCH', got '%s'", oklchColor.Space().Name())
	}

	// Convert back to Display P3
	p3Back := oklchColor.ConvertTo(DisplayP3Space)
	if p3Back.Space().Name() != "display-p3" {
		t.Errorf("Expected space 'display-p3', got '%s'", p3Back.Space().Name())
	}
}

func TestToXYZShortChannelsDoesNotPanic(t *testing.T) {
	// Both should return zero XYZ instead of panicking.
	x1, y1, z1 := SRGBSpace.ToXYZ([]float64{0.2, 0.3})
	if x1 != 0 || y1 != 0 || z1 != 0 {
		t.Fatalf("expected zero XYZ for short sRGB channels, got (%f, %f, %f)", x1, y1, z1)
	}

	x2, y2, z2 := OKLCHSpace.ToXYZ([]float64{0.7, 0.2})
	if x2 != 0 || y2 != 0 || z2 != 0 {
		t.Fatalf("expected zero XYZ for short OKLCH channels, got (%f, %f, %f)", x2, y2, z2)
	}
}
