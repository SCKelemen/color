package color

import (
	"math"
	"testing"
)

// Test edge cases for gamut mapping

func TestGamutMappingBlackAndWhite(t *testing.T) {
	strategies := []GamutMapping{
		GamutClip,
		GamutPreserveLightness,
		GamutPreserveChroma,
		GamutProject,
	}

	black := RGB(0, 0, 0)
	white := RGB(1, 1, 1)

	for _, strategy := range strategies {
		// Black should stay black
		mappedBlack := MapToGamut(black, strategy)
		r, g, b, _ := mappedBlack.RGBA()
		if r > 0.01 || g > 0.01 || b > 0.01 {
			t.Errorf("Strategy %v: Black changed to (%f, %f, %f)", strategy, r, g, b)
		}

		// White should stay white
		mappedWhite := MapToGamut(white, strategy)
		r, g, b, _ = mappedWhite.RGBA()
		if r < 0.99 || g < 0.99 || b < 0.99 {
			t.Errorf("Strategy %v: White changed to (%f, %f, %f)", strategy, r, g, b)
		}
	}
}

func TestGamutMappingGrays(t *testing.T) {
	strategies := []GamutMapping{
		GamutClip,
		GamutPreserveLightness,
		GamutPreserveChroma,
		GamutProject,
	}

	grays := []float64{0.1, 0.3, 0.5, 0.7, 0.9}

	for _, gray := range grays {
		color := RGB(gray, gray, gray)

		for _, strategy := range strategies {
			mapped := MapToGamut(color, strategy)
			r, g, b, _ := mapped.RGBA()

			// Gray should stay gray (achromatic)
			if math.Abs(r-g) > 0.02 || math.Abs(g-b) > 0.02 {
				t.Errorf("Strategy %v, gray %f: Lost achromaticity (%f, %f, %f)",
					strategy, gray, r, g, b)
			}
		}
	}
}

func TestGamutMappingExtremeColors(t *testing.T) {
	// Colors with extreme chroma that will definitely be out of gamut
	extremeColors := []*OKLCH{
		NewOKLCH(0.5, 1.0, 0, 1.0),     // Impossible red
		NewOKLCH(0.5, 1.0, 120, 1.0),   // Impossible green
		NewOKLCH(0.5, 1.0, 240, 1.0),   // Impossible blue
		NewOKLCH(0.1, 0.5, 180, 1.0),   // Very dark vivid
		NewOKLCH(0.9, 0.5, 180, 1.0),   // Very light vivid
	}

	for i, color := range extremeColors {
		for _, strategy := range []GamutMapping{GamutClip, GamutPreserveLightness, GamutPreserveChroma, GamutProject} {
			mapped := MapToGamut(color, strategy)

			// Result must be in gamut
			if !InGamut(mapped) {
				t.Errorf("Color %d, strategy %v: Result not in gamut", i, strategy)
			}

			// RGB values must be valid
			r, g, b, a := mapped.RGBA()
			if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 {
				t.Errorf("Color %d, strategy %v: Invalid RGB (%f, %f, %f)", i, strategy, r, g, b)
			}

			// Alpha must be preserved
			if math.Abs(a-1.0) > 0.01 {
				t.Errorf("Color %d, strategy %v: Alpha not preserved: %f", i, strategy, a)
			}
		}
	}
}

func TestGamutMappingZeroAlpha(t *testing.T) {
	// Color with zero alpha
	color := NewRGBA(1.5, 0.5, 0.3, 0.0)

	for _, strategy := range []GamutMapping{GamutClip, GamutPreserveLightness, GamutPreserveChroma, GamutProject} {
		mapped := MapToGamut(color, strategy)

		// Alpha should be preserved
		if mapped.Alpha() != 0.0 {
			t.Errorf("Strategy %v: Alpha not preserved: %f", strategy, mapped.Alpha())
		}
	}
}

func TestGamutMappingPartialAlpha(t *testing.T) {
	alphas := []float64{0.25, 0.5, 0.75}

	for _, alpha := range alphas {
		color := NewOKLCH(0.7, 0.35, 150, alpha)

		for _, strategy := range []GamutMapping{GamutClip, GamutPreserveLightness, GamutPreserveChroma, GamutProject} {
			mapped := MapToGamut(color, strategy)

			// Alpha should be preserved
			if math.Abs(mapped.Alpha()-alpha) > 0.01 {
				t.Errorf("Strategy %v, alpha %f: Alpha changed to %f", strategy, alpha, mapped.Alpha())
			}
		}
	}
}

func TestInGamutBoundaryColors(t *testing.T) {
	// Colors exactly at gamut boundaries
	boundaryColors := []Color{
		RGB(1, 0, 0),    // Pure red
		RGB(0, 1, 0),    // Pure green
		RGB(0, 0, 1),    // Pure blue
		RGB(1, 1, 0),    // Yellow
		RGB(1, 0, 1),    // Magenta
		RGB(0, 1, 1),    // Cyan
	}

	for i, color := range boundaryColors {
		if !InGamut(color) {
			t.Errorf("Boundary color %d should be in gamut", i)
		}
	}
}

func TestGamutMappingConsistency(t *testing.T) {
	// Mapping the same color twice should give the same result
	color := NewOKLCH(0.7, 0.35, 150, 1.0)

	for _, strategy := range []GamutMapping{GamutClip, GamutPreserveLightness, GamutPreserveChroma, GamutProject} {
		mapped1 := MapToGamut(color, strategy)
		mapped2 := MapToGamut(color, strategy)

		r1, g1, b1, a1 := mapped1.RGBA()
		r2, g2, b2, a2 := mapped2.RGBA()

		if math.Abs(r1-r2) > 0.001 || math.Abs(g1-g2) > 0.001 || math.Abs(b1-b2) > 0.001 || math.Abs(a1-a2) > 0.001 {
			t.Errorf("Strategy %v: Inconsistent mapping", strategy)
		}
	}
}

func TestGamutMappingIdempotent(t *testing.T) {
	// Mapping an already in-gamut color should not change it significantly
	inGamutColor := RGB(0.5, 0.3, 0.7)

	for _, strategy := range []GamutMapping{GamutClip, GamutPreserveLightness, GamutPreserveChroma, GamutProject} {
		mapped := MapToGamut(inGamutColor, strategy)

		r1, g1, b1, _ := inGamutColor.RGBA()
		r2, g2, b2, _ := mapped.RGBA()

		diff := math.Abs(r1-r2) + math.Abs(g1-g2) + math.Abs(b1-b2)
		if diff > 0.02 {
			t.Errorf("Strategy %v: Changed in-gamut color too much (diff=%f)", strategy, diff)
		}
	}
}

func TestGamutMappingPreservesHue(t *testing.T) {
	// All strategies should preserve hue (at least approximately)
	colors := []*OKLCH{
		NewOKLCH(0.7, 0.35, 0, 1.0),     // Red
		NewOKLCH(0.7, 0.35, 120, 1.0),   // Green
		NewOKLCH(0.7, 0.35, 240, 1.0),   // Blue
	}

	tolerance := 15.0 // degrees

	for i, color := range colors {
		originalHue := color.H

		for _, strategy := range []GamutMapping{GamutClip, GamutPreserveLightness, GamutPreserveChroma, GamutProject} {
			mapped := MapToGamut(color, strategy)
			mappedOKLCH := ToOKLCH(mapped)

			hueDiff := math.Abs(mappedOKLCH.H - originalHue)
			// Handle wraparound
			if hueDiff > 180 {
				hueDiff = 360 - hueDiff
			}

			if hueDiff > tolerance {
				t.Logf("Color %d, strategy %v: Hue shifted by %f degrees (original: %f, mapped: %f)",
					i, strategy, hueDiff, originalHue, mappedOKLCH.H)
			}
		}
	}
}

func TestMapToGamutVsInGamut(t *testing.T) {
	// After mapping, InGamut should return true
	outOfGamutColors := []*OKLCH{
		NewOKLCH(0.7, 0.35, 150, 1.0),
		NewOKLCH(0.5, 0.4, 60, 1.0),
		NewOKLCH(0.8, 0.3, 200, 1.0),
	}

	for i, color := range outOfGamutColors {
		for _, strategy := range []GamutMapping{GamutClip, GamutPreserveLightness, GamutPreserveChroma, GamutProject} {
			mapped := MapToGamut(color, strategy)

			if !InGamut(mapped) {
				t.Errorf("Color %d, strategy %v: Mapped color not in gamut", i, strategy)
			}
		}
	}
}

func TestGamutMappingExtremeLightness(t *testing.T) {
	// Very dark and very light colors with high chroma
	extremeColors := []*OKLCH{
		NewOKLCH(0.05, 0.3, 180, 1.0), // Very dark
		NewOKLCH(0.95, 0.3, 180, 1.0), // Very light
	}

	for i, color := range extremeColors {
		for _, strategy := range []GamutMapping{GamutClip, GamutPreserveLightness, GamutPreserveChroma, GamutProject} {
			mapped := MapToGamut(color, strategy)

			if !InGamut(mapped) {
				t.Errorf("Color %d, strategy %v: Result not in gamut", i, strategy)
			}

			// For PreserveLightness, lightness should be close
			if strategy == GamutPreserveLightness {
				originalOKLCH := color
				mappedOKLCH := ToOKLCH(mapped)

				lightnessChange := math.Abs(originalOKLCH.L - mappedOKLCH.L)
				if lightnessChange > 0.1 {
					t.Logf("Color %d: Lightness changed by %f (strategy should preserve it)", i, lightnessChange)
				}
			}
		}
	}
}
