package color

import (
	"testing"
)

func TestInGamut(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		expected bool
	}{
		{
			name:     "Pure red in gamut",
			color:    RGB(1, 0, 0),
			expected: true,
		},
		{
			name:     "Mid gray in gamut",
			color:    RGB(0.5, 0.5, 0.5),
			expected: true,
		},
		{
			name:     "Black in gamut",
			color:    RGB(0, 0, 0),
			expected: true,
		},
		{
			name:     "White in gamut",
			color:    RGB(1, 1, 1),
			expected: true,
		},
		{
			name:     "Vivid OKLCH (out of sRGB gamut)",
			color:    NewOKLCH(0.7, 0.35, 150, 1.0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InGamut(tt.color)
			if result != tt.expected {
				r, g, b, _ := tt.color.RGBA()
				t.Errorf("InGamut(%f, %f, %f) = %v, want %v", r, g, b, result, tt.expected)
			}
		})
	}
}

func TestMapToGamutClip(t *testing.T) {
	// Test that clipping brings out-of-gamut colors into gamut
	outOfGamut := NewRGBA(1.2, 0.5, -0.1, 1)
	mapped := MapToGamut(outOfGamut, GamutClip)

	if !InGamut(mapped) {
		r, g, b, _ := mapped.RGBA()
		t.Errorf("MapToGamut(GamutClip) produced out-of-gamut color: (%f, %f, %f)", r, g, b)
	}

	// Check that clamping was applied
	r, g, b, _ := mapped.RGBA()
	if r > 1 || g > 1 || b > 1 {
		t.Errorf("Clipped color has components > 1: (%f, %f, %f)", r, g, b)
	}
	if r < 0 || g < 0 || b < 0 {
		t.Errorf("Clipped color has components < 0: (%f, %f, %f)", r, g, b)
	}
}

func TestMapToGamutPreserveLightness(t *testing.T) {
	// Create a vivid out-of-gamut color
	vivid := NewOKLCH(0.7, 0.35, 150, 1.0) // High chroma teal

	mapped := MapToGamut(vivid, GamutPreserveLightness)

	if !InGamut(mapped) {
		t.Error("MapToGamut(GamutPreserveLightness) produced out-of-gamut color")
	}

	// Check that lightness is approximately preserved
	originalOKLCH := ToOKLCH(vivid)
	mappedOKLCH := ToOKLCH(mapped)

	if abs(originalOKLCH.L-mappedOKLCH.L) > 0.05 {
		t.Errorf("Lightness not preserved: original %f, mapped %f", originalOKLCH.L, mappedOKLCH.L)
	}
}

func TestMapToGamutPreserveChroma(t *testing.T) {
	// Create a vivid out-of-gamut color
	vivid := NewOKLCH(0.7, 0.35, 150, 1.0)

	mapped := MapToGamut(vivid, GamutPreserveChroma)

	if !InGamut(mapped) {
		t.Error("MapToGamut(GamutPreserveChroma) produced out-of-gamut color")
	}

	// Chroma should be approximately preserved
	originalOKLCH := ToOKLCH(vivid)
	mappedOKLCH := ToOKLCH(mapped)

	if abs(originalOKLCH.C-mappedOKLCH.C) > 0.05 {
		t.Errorf("Chroma not preserved: original %f, mapped %f", originalOKLCH.C, mappedOKLCH.C)
	}
}

func TestMapToGamutProject(t *testing.T) {
	// Create out-of-gamut color
	vivid := NewOKLCH(0.7, 0.35, 150, 1.0)

	mapped := MapToGamut(vivid, GamutProject)

	if !InGamut(mapped) {
		t.Error("MapToGamut(GamutProject) produced out-of-gamut color")
	}

	// Project should be a good balance
	originalOKLCH := ToOKLCH(vivid)
	mappedOKLCH := ToOKLCH(mapped)

	// Hue should always be preserved
	hueDiff := abs(originalOKLCH.H - mappedOKLCH.H)
	if hueDiff > 360 {
		hueDiff = 720 - hueDiff
	}
	if hueDiff > 5 {
		t.Errorf("Hue not preserved: original %f, mapped %f", originalOKLCH.H, mappedOKLCH.H)
	}
}

func TestMapToGamutAlreadyInGamut(t *testing.T) {
	// Colors already in gamut should not change significantly
	inGamut := RGB(0.5, 0.3, 0.8)

	strategies := []GamutMapping{
		GamutClip,
		GamutPreserveLightness,
		GamutPreserveChroma,
		GamutProject,
	}

	for _, strategy := range strategies {
		mapped := MapToGamut(inGamut, strategy)
		r1, g1, b1, _ := inGamut.RGBA()
		r2, g2, b2, _ := mapped.RGBA()

		diff := abs(r1-r2) + abs(g1-g2) + abs(b1-b2)
		if diff > 0.01 {
			t.Errorf("Strategy %d changed in-gamut color too much: diff = %f", strategy, diff)
		}
	}
}

func TestGamutMappingPreservesAlpha(t *testing.T) {
	// All gamut mapping should preserve alpha
	color := NewRGBA(1.2, 0.5, 0.3, 0.7) // Out of gamut with alpha

	strategies := []GamutMapping{
		GamutClip,
		GamutPreserveLightness,
		GamutPreserveChroma,
		GamutProject,
	}

	for _, strategy := range strategies {
		mapped := MapToGamut(color, strategy)
		alpha := mapped.Alpha()
		if abs(alpha-0.7) > 0.01 {
			t.Errorf("Strategy %d didn't preserve alpha: got %f, want 0.7", strategy, alpha)
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func BenchmarkInGamut(b *testing.B) {
	c := RGB(0.5, 0.6, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InGamut(c)
	}
}

func BenchmarkMapToGamutClip(b *testing.B) {
	c := NewRGBA(1.2, 0.5, 0.3, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToGamut(c, GamutClip)
	}
}

func BenchmarkMapToGamutPreserveLightness(b *testing.B) {
	c := NewOKLCH(0.7, 0.35, 150, 1.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToGamut(c, GamutPreserveLightness)
	}
}

func BenchmarkMapToGamutProject(b *testing.B) {
	c := NewOKLCH(0.7, 0.35, 150, 1.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MapToGamut(c, GamutProject)
	}
}
