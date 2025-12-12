package color

import (
	"math"
)

// WCAG (Web Content Accessibility Guidelines) provides functions for
// checking color contrast and accessibility compliance.

// ContrastRatio calculates the WCAG 2.1 contrast ratio between two colors.
// The ratio ranges from 1:1 (no contrast) to 21:1 (maximum contrast).
//
// WCAG Requirements:
//   - Normal text (AA): 4.5:1
//   - Normal text (AAA): 7:1
//   - Large text (AA): 3:1
//   - Large text (AAA): 4.5:1
//   - UI components: 3:1
//
// Reference: https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum.html
func ContrastRatio(c1, c2 Color) float64 {
	l1 := relativeLuminance(c1)
	l2 := relativeLuminance(c2)

	// Ensure l1 is the lighter color
	if l2 > l1 {
		l1, l2 = l2, l1
	}

	return (l1 + 0.05) / (l2 + 0.05)
}

// relativeLuminance calculates the relative luminance according to WCAG 2.1.
// This is different from perceptual lightness (L* in LAB/OKLCH).
func relativeLuminance(c Color) float64 {
	r, g, b, _ := c.RGBA()

	// Convert to linear RGB (remove gamma correction)
	r = sRGBToLinear(r)
	g = sRGBToLinear(g)
	b = sRGBToLinear(b)

	// Calculate relative luminance using ITU-R BT.709 coefficients
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// sRGBToLinear converts sRGB component to linear RGB for luminance calculation.
func sRGBToLinear(component float64) float64 {
	if component <= 0.04045 {
		return component / 12.92
	}
	return math.Pow((component+0.055)/1.055, 2.4)
}

// WCAGLevel represents WCAG conformance levels.
type WCAGLevel int

const (
	// WCAGAA is the minimum WCAG 2.1 Level AA compliance
	WCAGAA WCAGLevel = iota
	// WCAGAAA is the enhanced WCAG 2.1 Level AAA compliance
	WCAGAAA
)

// TextSize represents the size category of text for WCAG guidelines.
type TextSize int

const (
	// NormalText is regular-sized text (< 18pt or < 14pt bold)
	NormalText TextSize = iota
	// LargeText is large text (≥ 18pt or ≥ 14pt bold)
	LargeText
)

// ContrastCompliance checks if two colors meet WCAG contrast requirements.
type ContrastCompliance struct {
	Ratio         float64
	AANormal      bool // AA level for normal text (4.5:1)
	AALarge       bool // AA level for large text (3:1)
	AAANormal     bool // AAA level for normal text (7:1)
	AAALarge      bool // AAA level for large text (4.5:1)
	UIComponents  bool // UI components and graphics (3:1)
}

// CheckContrast returns detailed WCAG compliance information for two colors.
func CheckContrast(foreground, background Color) ContrastCompliance {
	ratio := ContrastRatio(foreground, background)

	return ContrastCompliance{
		Ratio:         ratio,
		AANormal:      ratio >= 4.5,
		AALarge:       ratio >= 3.0,
		AAANormal:     ratio >= 7.0,
		AAALarge:      ratio >= 4.5,
		UIComponents:  ratio >= 3.0,
	}
}

// IsAccessible checks if a color pair meets WCAG requirements for given level and text size.
func IsAccessible(foreground, background Color, level WCAGLevel, textSize TextSize) bool {
	ratio := ContrastRatio(foreground, background)

	switch level {
	case WCAGAA:
		if textSize == NormalText {
			return ratio >= 4.5
		}
		return ratio >= 3.0
	case WCAGAAA:
		if textSize == NormalText {
			return ratio >= 7.0
		}
		return ratio >= 4.5
	default:
		return false
	}
}

// SuggestAccessibleForeground finds an accessible foreground color for a given background.
// It adjusts the lightness of the base color to meet WCAG requirements.
// Returns the suggested color and whether it meets the requirements.
func SuggestAccessibleForeground(base, background Color, level WCAGLevel, textSize TextSize) (Color, bool) {
	baseOKLCH := ToOKLCH(base)
	bgLuminance := relativeLuminance(background)

	// Determine target contrast ratio
	var targetRatio float64
	switch level {
	case WCAGAA:
		if textSize == NormalText {
			targetRatio = 4.5
		} else {
			targetRatio = 3.0
		}
	case WCAGAAA:
		if textSize == NormalText {
			targetRatio = 7.0
		} else {
			targetRatio = 4.5
		}
	}

	// Try different lightness values to find one that works
	// First, determine if we need a lighter or darker foreground
	needLighter := bgLuminance < 0.5

	bestColor := base
	bestRatio := 0.0

	// Try lightness values from 0.1 to 0.9
	for l := 0.1; l <= 0.9; l += 0.05 {
		if needLighter && l < 0.5 {
			continue // Skip dark colors if we need lighter
		}
		if !needLighter && l > 0.5 {
			continue // Skip light colors if we need darker
		}

		testColor := NewOKLCH(l, baseOKLCH.C, baseOKLCH.H, baseOKLCH.A_)
		ratio := ContrastRatio(testColor, background)

		if ratio >= targetRatio {
			return testColor, true
		}

		if ratio > bestRatio {
			bestRatio = ratio
			bestColor = testColor
		}
	}

	// If we couldn't meet the target, return the best we found
	return bestColor, bestRatio >= targetRatio
}

// SuggestAccessibleBackground finds an accessible background color for a given foreground.
// It adjusts the lightness of the base color to meet WCAG requirements.
func SuggestAccessibleBackground(foreground, baseBackground Color, level WCAGLevel, textSize TextSize) (Color, bool) {
	// This is symmetric to SuggestAccessibleForeground
	return SuggestAccessibleForeground(baseBackground, foreground, level, textSize)
}

// ColorBlindnessType represents different types of color vision deficiency.
type ColorBlindnessType int

const (
	// Protanopia is red-blind (missing L-cone, ~1% males)
	Protanopia ColorBlindnessType = iota
	// Protanomaly is red-weak (anomalous L-cone, ~1% males)
	Protanomaly
	// Deuteranopia is green-blind (missing M-cone, ~1% males)
	Deuteranopia
	// Deuteranomaly is green-weak (anomalous M-cone, ~5% males)
	Deuteranomaly
	// Tritanopia is blue-blind (missing S-cone, very rare)
	Tritanopia
	// Tritanomaly is blue-weak (anomalous S-cone, very rare)
	Tritanomaly
	// Achromatopsia is complete color blindness (very rare)
	Achromatopsia
	// Achromatomaly is partial color blindness (very rare)
	Achromatomaly
)

// SimulateColorBlindness simulates how a color appears to someone with color vision deficiency.
// Uses the Brettel et al. (1997) and Viénot et al. (1999) simulation matrices.
func SimulateColorBlindness(c Color, cvdType ColorBlindnessType) Color {
	r, g, b, a := c.RGBA()

	// Convert to linear RGB for accurate simulation
	r = sRGBToLinear(r)
	g = sRGBToLinear(g)
	b = sRGBToLinear(b)

	var rOut, gOut, bOut float64

	switch cvdType {
	case Protanopia:
		// Red-blind: confuse red and green
		rOut = 0.56667*r + 0.43333*g + 0.00000*b
		gOut = 0.55833*r + 0.44167*g + 0.00000*b
		bOut = 0.00000*r + 0.24167*g + 0.75833*b

	case Protanomaly:
		// Red-weak: partial red confusion
		rOut = 0.81667*r + 0.18333*g + 0.00000*b
		gOut = 0.33333*r + 0.66667*g + 0.00000*b
		bOut = 0.00000*r + 0.12500*g + 0.87500*b

	case Deuteranopia:
		// Green-blind: confuse green and red
		rOut = 0.625*r + 0.375*g + 0.000*b
		gOut = 0.700*r + 0.300*g + 0.000*b
		bOut = 0.000*r + 0.300*g + 0.700*b

	case Deuteranomaly:
		// Green-weak: partial green confusion
		rOut = 0.80000*r + 0.20000*g + 0.00000*b
		gOut = 0.25833*r + 0.74167*g + 0.00000*b
		bOut = 0.00000*r + 0.14167*g + 0.85833*b

	case Tritanopia:
		// Blue-blind: confuse blue and yellow
		rOut = 0.95000*r + 0.05000*g + 0.00000*b
		gOut = 0.00000*r + 0.43333*g + 0.56667*b
		bOut = 0.00000*r + 0.47500*g + 0.52500*b

	case Tritanomaly:
		// Blue-weak: partial blue confusion
		rOut = 0.96667*r + 0.03333*g + 0.00000*b
		gOut = 0.00000*r + 0.73333*g + 0.26667*b
		bOut = 0.00000*r + 0.18333*g + 0.81667*b

	case Achromatopsia:
		// Complete color blindness: see only luminance
		lum := 0.2126*r + 0.7152*g + 0.0722*b
		rOut, gOut, bOut = lum, lum, lum

	case Achromatomaly:
		// Partial color blindness: reduced color perception
		lum := 0.2126*r + 0.7152*g + 0.0722*b
		rOut = 0.618*lum + 0.382*r
		gOut = 0.618*lum + 0.382*g
		bOut = 0.618*lum + 0.382*b

	default:
		rOut, gOut, bOut = r, g, b
	}

	// Convert back to sRGB
	rOut = linearToSRGB(rOut)
	gOut = linearToSRGB(gOut)
	bOut = linearToSRGB(bOut)

	return NewRGBA(rOut, gOut, bOut, a)
}

// linearToSRGB converts linear RGB to sRGB (applies gamma correction).
func linearToSRGB(component float64) float64 {
	if component <= 0.0031308 {
		return component * 12.92
	}
	return 1.055*math.Pow(component, 1/2.4) - 0.055
}

// IsColorBlindSafe checks if two colors are distinguishable for people with color blindness.
// It simulates the colors for the given CVD type and checks if they have sufficient contrast.
func IsColorBlindSafe(c1, c2 Color, cvdType ColorBlindnessType, minContrast float64) bool {
	sim1 := SimulateColorBlindness(c1, cvdType)
	sim2 := SimulateColorBlindness(c2, cvdType)

	ratio := ContrastRatio(sim1, sim2)
	return ratio >= minContrast
}

// CheckColorBlindSafety checks if two colors are safe for all common types of color blindness.
// Returns a map showing which types pass the minimum contrast requirement.
func CheckColorBlindSafety(c1, c2 Color, minContrast float64) map[ColorBlindnessType]bool {
	types := []ColorBlindnessType{
		Protanopia,
		Protanomaly,
		Deuteranopia,
		Deuteranomaly,
		Tritanopia,
		Tritanomaly,
	}

	results := make(map[ColorBlindnessType]bool)
	for _, cvdType := range types {
		results[cvdType] = IsColorBlindSafe(c1, c2, cvdType, minContrast)
	}

	return results
}
