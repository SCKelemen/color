package color

// Harmony provides functions for generating harmonious color schemes
// based on color theory principles. All functions work in the OKLCH color
// space for perceptually uniform results.

// Complementary returns the complementary color (opposite on the color wheel).
// This creates maximum contrast while maintaining harmony.
func Complementary(c Color) Color {
	oklch := ToOKLCH(c)
	return NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+180), oklch.A_)
}

// Triadic returns a triadic color scheme (3 colors evenly spaced on the color wheel).
// Returns the input color plus two harmonious colors at 120° intervals.
func Triadic(c Color) []Color {
	oklch := ToOKLCH(c)
	return []Color{
		c,
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+120), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+240), oklch.A_),
	}
}

// Tetradic returns a tetradic (square) color scheme (4 colors evenly spaced).
// Returns the input color plus three harmonious colors at 90° intervals.
func Tetradic(c Color) []Color {
	oklch := ToOKLCH(c)
	return []Color{
		c,
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+90), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+180), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+270), oklch.A_),
	}
}

// Square is an alias for Tetradic (4 colors at 90° intervals).
func Square(c Color) []Color {
	return Tetradic(c)
}

// Analogous returns an analogous color scheme (colors adjacent on the wheel).
// By default, returns 3 colors: the input color and two neighbors at ±30°.
// Use AnalogousN for custom angles and count.
func Analogous(c Color) []Color {
	return AnalogousN(c, 3, 30)
}

// AnalogousN returns n analogous colors with custom angle spacing.
// The colors are distributed evenly within the angle range on both sides.
//
// Example:
//   AnalogousN(color, 5, 30) returns 5 colors from -30° to +30°
func AnalogousN(c Color, n int, angle float64) []Color {
	if n <= 0 {
		return []Color{}
	}
	if n == 1 {
		return []Color{c}
	}

	oklch := ToOKLCH(c)
	colors := make([]Color, n)

	// Distribute colors evenly from -angle to +angle
	for i := 0; i < n; i++ {
		offset := -angle + (float64(i) / float64(n-1) * 2 * angle)
		colors[i] = NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+offset), oklch.A_)
	}

	return colors
}

// SplitComplementary returns a split-complementary color scheme.
// This uses colors on either side of the complement (opposite + ±30°).
// Returns 3 colors: the input color and two colors flanking the complement.
func SplitComplementary(c Color) []Color {
	return SplitComplementaryN(c, 30)
}

// SplitComplementaryN returns a split-complementary scheme with custom angle.
// The angle determines how far from the complement the flanking colors are.
func SplitComplementaryN(c Color, angle float64) []Color {
	oklch := ToOKLCH(c)
	complement := oklch.H + 180

	return []Color{
		c,
		NewOKLCH(oklch.L, oklch.C, normalizeHue(complement-angle), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(complement+angle), oklch.A_),
	}
}

// Monochromatic returns a monochromatic color scheme (same hue, varying lightness).
// Returns n colors with the same hue but different lightness values.
// The colors range from dark to light, preserving the original chroma.
func Monochromatic(c Color, n int) []Color {
	if n <= 0 {
		return []Color{}
	}
	if n == 1 {
		return []Color{c}
	}

	oklch := ToOKLCH(c)
	colors := make([]Color, n)

	// Distribute lightness from 0.2 to 0.9
	minL := 0.2
	maxL := 0.9

	for i := 0; i < n; i++ {
		l := minL + (float64(i) / float64(n-1) * (maxL - minL))
		// Adjust chroma for very light/dark colors to keep them in gamut
		adjustedC := oklch.C
		if l < 0.3 || l > 0.85 {
			adjustedC = oklch.C * 0.7
		}
		colors[i] = NewOKLCH(l, adjustedC, oklch.H, oklch.A_)
	}

	return colors
}

// MonochromaticCentered returns a monochromatic scheme centered on the input color.
// The input color will be in the middle of the returned palette.
func MonochromaticCentered(c Color, n int) []Color {
	if n <= 0 {
		return []Color{}
	}
	if n == 1 {
		return []Color{c}
	}

	oklch := ToOKLCH(c)
	colors := make([]Color, n)

	// Center the palette around the input lightness
	centerL := oklch.L
	range_ := 0.35 // Total range of ±0.35

	for i := 0; i < n; i++ {
		offset := -range_ + (float64(i) / float64(n-1) * 2 * range_)
		l := clamp01(centerL + offset)

		// Adjust chroma for extreme lightness
		adjustedC := oklch.C
		if l < 0.3 || l > 0.85 {
			adjustedC = oklch.C * 0.7
		}

		colors[i] = NewOKLCH(l, adjustedC, oklch.H, oklch.A_)
	}

	return colors
}

// Shades returns darker variations of the color (reduced lightness).
// Returns n colors from the original to nearly black, preserving hue and chroma.
func Shades(c Color, n int) []Color {
	if n <= 0 {
		return []Color{}
	}
	if n == 1 {
		return []Color{c}
	}

	oklch := ToOKLCH(c)
	colors := make([]Color, n)

	minL := 0.1

	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		l := oklch.L * (1 - t) + minL * t
		// Reduce chroma as we get darker
		adjustedC := oklch.C * (1 - t*0.3)
		colors[i] = NewOKLCH(l, adjustedC, oklch.H, oklch.A_)
	}

	return colors
}

// Tints returns lighter variations of the color (increased lightness).
// Returns n colors from the original to nearly white, preserving hue.
func Tints(c Color, n int) []Color {
	if n <= 0 {
		return []Color{}
	}
	if n == 1 {
		return []Color{c}
	}

	oklch := ToOKLCH(c)
	colors := make([]Color, n)

	maxL := 0.95

	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		l := oklch.L * (1 - t) + maxL * t
		// Reduce chroma as we get lighter
		adjustedC := oklch.C * (1 - t*0.5)
		colors[i] = NewOKLCH(l, adjustedC, oklch.H, oklch.A_)
	}

	return colors
}

// Tones returns variations with reduced chroma (desaturated, more gray).
// Returns n colors from the original to a neutral gray, preserving hue and lightness.
func Tones(c Color, n int) []Color {
	if n <= 0 {
		return []Color{}
	}
	if n == 1 {
		return []Color{c}
	}

	oklch := ToOKLCH(c)
	colors := make([]Color, n)

	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		// Reduce chroma to 0 (gray)
		adjustedC := oklch.C * (1 - t)
		colors[i] = NewOKLCH(oklch.L, adjustedC, oklch.H, oklch.A_)
	}

	return colors
}

// Rectangle returns a rectangular (double split-complementary) color scheme.
// This uses two pairs of complementary colors, creating 4 colors total.
// The angle parameter controls the spacing (typically 30-60°).
func Rectangle(c Color, angle float64) []Color {
	oklch := ToOKLCH(c)

	return []Color{
		c,
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+angle), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+180), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+180+angle), oklch.A_),
	}
}

// DoubleSplitComplementary returns a double split-complementary scheme (6 colors).
// This combines two split-complementary schemes from opposite sides of the wheel.
func DoubleSplitComplementary(c Color) []Color {
	oklch := ToOKLCH(c)
	angle := 30.0
	complement := oklch.H + 180

	return []Color{
		c,
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H-angle), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(oklch.H+angle), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(complement), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(complement-angle), oklch.A_),
		NewOKLCH(oklch.L, oklch.C, normalizeHue(complement+angle), oklch.A_),
	}
}
