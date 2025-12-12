package color

import "math"

// GamutMapping specifies the algorithm to use when a color is out of the target gamut.
type GamutMapping int

const (
	// GamutClip simply clips RGB values to [0, 1]. This is fast but may shift hue.
	GamutClip GamutMapping = iota

	// GamutPreserveChroma reduces lightness while preserving chroma until the color is in-gamut.
	// This maintains saturation but may significantly alter brightness.
	GamutPreserveChroma

	// GamutPreserveLightness reduces chroma while preserving lightness until the color is in-gamut.
	// This maintains brightness but reduces saturation (recommended for most uses).
	GamutPreserveLightness

	// GamutProject projects the color onto the gamut boundary along the shortest perceptual path.
	// This balances lightness and chroma adjustments for the most perceptually accurate result.
	GamutProject
)

// InGamut checks if a color is within the sRGB gamut.
// A color is in-gamut if all RGB components are in [0, 1] when converted to sRGB.
func InGamut(c Color) bool {
	r, g, b, _ := c.RGBA()
	return r >= 0 && r <= 1 && g >= 0 && g <= 1 && b >= 0 && b <= 1
}

// MapToGamut maps a color into the sRGB gamut using the specified mapping algorithm.
// If the color is already in-gamut, it is returned unchanged.
func MapToGamut(c Color, mapping GamutMapping) Color {
	if InGamut(c) {
		return c
	}

	switch mapping {
	case GamutClip:
		return mapClip(c)
	case GamutPreserveChroma:
		return mapPreserveChroma(c)
	case GamutPreserveLightness:
		return mapPreserveLightness(c)
	case GamutProject:
		return mapProject(c)
	default:
		return mapClip(c)
	}
}

// mapClip clips RGB values to [0, 1]. Fast but may shift hue.
func mapClip(c Color) Color {
	r, g, b, a := c.RGBA()
	return NewRGBA(clamp01(r), clamp01(g), clamp01(b), a)
}

// mapPreserveChroma reduces lightness while keeping chroma constant.
func mapPreserveChroma(c Color) Color {
	oklch := ToOKLCH(c)

	// Binary search for the maximum lightness that keeps the color in-gamut
	lMin := 0.0
	lMax := oklch.L

	// If current color is darker, search upward
	testColor := NewOKLCH(lMax, oklch.C, oklch.H, oklch.Alpha())
	if !InGamut(testColor) {
		// Need to darken
		for i := 0; i < 20; i++ { // 20 iterations gives ~0.0001% precision
			lMid := (lMin + lMax) / 2
			testColor = NewOKLCH(lMid, oklch.C, oklch.H, oklch.Alpha())
			if InGamut(testColor) {
				lMin = lMid
			} else {
				lMax = lMid
			}
		}
		return NewOKLCH(lMin, oklch.C, oklch.H, oklch.Alpha())
	}

	// If we can't keep the chroma, reduce lightness toward 0
	return NewOKLCH(0, oklch.C, oklch.H, oklch.Alpha())
}

// mapPreserveLightness reduces chroma while keeping lightness constant.
func mapPreserveLightness(c Color) Color {
	oklch := ToOKLCH(c)

	// Binary search for the maximum chroma that keeps the color in-gamut
	cMin := 0.0
	cMax := oklch.C

	for i := 0; i < 20; i++ { // 20 iterations gives ~0.0001% precision
		cMid := (cMin + cMax) / 2
		testColor := NewOKLCH(oklch.L, cMid, oklch.H, oklch.Alpha())
		if InGamut(testColor) {
			cMin = cMid
		} else {
			cMax = cMid
		}
	}

	return NewOKLCH(oklch.L, cMin, oklch.H, oklch.Alpha())
}

// mapProject projects the color onto the gamut boundary.
// This uses a perceptual projection that balances lightness and chroma.
func mapProject(c Color) Color {
	oklch := ToOKLCH(c)

	// Start with the current color
	bestL := oklch.L
	bestC := oklch.C
	bestDistance := math.MaxFloat64

	// Search for the closest in-gamut color using a grid search
	// This is a simplified version; a full implementation would use more sophisticated methods

	// Try reducing chroma
	for ratio := 0.0; ratio <= 1.0; ratio += 0.05 {
		testC := oklch.C * (1 - ratio)
		testL := oklch.L
		testColor := NewOKLCH(testL, testC, oklch.H, oklch.Alpha())

		if InGamut(testColor) {
			// Calculate perceptual distance
			distance := math.Sqrt((testL-oklch.L)*(testL-oklch.L) + (testC-oklch.C)*(testC-oklch.C))
			if distance < bestDistance {
				bestDistance = distance
				bestL = testL
				bestC = testC
			}
			break // Found an in-gamut color, stop searching in this direction
		}
	}

	// Try reducing lightness
	for ratio := 0.0; ratio <= 1.0; ratio += 0.05 {
		testL := oklch.L * (1 - ratio)
		testC := oklch.C
		testColor := NewOKLCH(testL, testC, oklch.H, oklch.Alpha())

		if InGamut(testColor) {
			distance := math.Sqrt((testL-oklch.L)*(testL-oklch.L) + (testC-oklch.C)*(testC-oklch.C))
			if distance < bestDistance {
				bestDistance = distance
				bestL = testL
				bestC = testC
			}
			break
		}
	}

	// Try adjusting both
	for ratioC := 0.0; ratioC <= 1.0; ratioC += 0.1 {
		for ratioL := 0.0; ratioL <= 1.0; ratioL += 0.1 {
			testL := oklch.L * (1 - ratioL)
			testC := oklch.C * (1 - ratioC)
			testColor := NewOKLCH(testL, testC, oklch.H, oklch.Alpha())

			if InGamut(testColor) {
				distance := math.Sqrt((testL-oklch.L)*(testL-oklch.L) + (testC-oklch.C)*(testC-oklch.C))
				if distance < bestDistance {
					bestDistance = distance
					bestL = testL
					bestC = testC
				}
			}
		}
	}

	// If still not found, fallback to preserve lightness
	if bestDistance == math.MaxFloat64 {
		return mapPreserveLightness(c)
	}

	return NewOKLCH(bestL, bestC, oklch.H, oklch.Alpha())
}

// ClipToGamut is a convenience function that clips a color to the sRGB gamut.
// This is equivalent to MapToGamut(c, GamutClip).
func ClipToGamut(c Color) Color {
	return MapToGamut(c, GamutClip)
}
