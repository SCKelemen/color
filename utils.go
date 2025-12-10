package color

import "math"

// Lighten increases the lightness of a color by the specified amount.
// Amount should be in the range [0, 1], where 0 is no change and 1 is maximum lightening.
func Lighten(c Color, amount float64) Color {
	amount = clamp01(amount)
	
	// Convert to OKLCH for perceptually uniform lightening
	oklch := ToOKLCH(c)
	oklch.L = clamp01(oklch.L + amount*(1-oklch.L))
	return oklch
}

// Darken decreases the lightness of a color by the specified amount.
// Amount should be in the range [0, 1], where 0 is no change and 1 is maximum darkening.
func Darken(c Color, amount float64) Color {
	amount = clamp01(amount)
	
	// Convert to OKLCH for perceptually uniform darkening
	oklch := ToOKLCH(c)
	oklch.L = clamp01(oklch.L * (1 - amount))
	return oklch
}

// Saturate increases the saturation of a color by the specified amount.
// Amount should be in the range [0, 1], where 0 is no change and 1 is maximum saturation.
func Saturate(c Color, amount float64) Color {
	amount = clamp01(amount)
	
	oklch := ToOKLCH(c)
	maxC := estimateMaxChroma(oklch.L, oklch.H)
	oklch.C = clamp(oklch.C+amount*(maxC-oklch.C), 0, maxC)
	return oklch
}

// Desaturate decreases the saturation of a color by the specified amount.
// Amount should be in the range [0, 1], where 0 is no change and 1 is complete desaturation.
func Desaturate(c Color, amount float64) Color {
	amount = clamp01(amount)
	
	oklch := ToOKLCH(c)
	oklch.C = oklch.C * (1 - amount)
	return oklch
}

// Mix blends two colors together.
// Weight should be in the range [0, 1], where 0 returns c1 and 1 returns c2.
func Mix(c1, c2 Color, weight float64) Color {
	weight = clamp01(weight)
	
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	
	// Mix RGB values
	r := r1*(1-weight) + r2*weight
	g := g1*(1-weight) + g2*weight
	b := b1*(1-weight) + b2*weight
	
	// Mix alpha values
	a := a1*(1-weight) + a2*weight
	
	return NewRGBA(r, g, b, a)
}

// MixOKLCH blends two colors in OKLCH space for perceptually uniform mixing.
func MixOKLCH(c1, c2 Color, weight float64) Color {
	weight = clamp01(weight)
	
	oklch1 := ToOKLCH(c1)
	oklch2 := ToOKLCH(c2)
	
	// Interpolate in OKLCH space
	l := oklch1.L*(1-weight) + oklch2.L*weight
	c := oklch1.C*(1-weight) + oklch2.C*weight
	
	// Handle hue interpolation (shortest path)
	h1, h2 := oklch1.H, oklch2.H
	dh := h2 - h1
	if math.Abs(dh) > 180 {
		if dh > 0 {
			dh -= 360
		} else {
			dh += 360
		}
	}
	h := normalizeHue(h1 + dh*weight)
	
	// Mix alpha
	a := oklch1.Alpha()*(1-weight) + oklch2.Alpha()*weight
	
	return NewOKLCH(l, c, h, a)
}

// AdjustHue shifts the hue of a color by the specified degrees.
func AdjustHue(c Color, degrees float64) Color {
	oklch := ToOKLCH(c)
	oklch.H = normalizeHue(oklch.H + degrees)
	return oklch
}

// Invert inverts the RGB values of a color.
func Invert(c Color) Color {
	r, g, b, a := c.RGBA()
	return NewRGBA(1-r, 1-g, 1-b, a)
}

// Grayscale converts a color to grayscale.
func Grayscale(c Color) Color {
	r, g, b, a := c.RGBA()
	// Use luminance formula
	gray := 0.299*r + 0.587*g + 0.114*b
	return NewRGBA(gray, gray, gray, a)
}

// Complement returns the complementary color (hue shifted by 180 degrees).
func Complement(c Color) Color {
	return AdjustHue(c, 180)
}

// estimateMaxChroma estimates the maximum chroma for a given lightness and hue in OKLCH.
// This is a simplified approximation.
func estimateMaxChroma(l, h float64) float64 {
	// Rough approximation - in practice, this would require gamut mapping
	// For now, use a conservative estimate
	return 0.4 * (1 - math.Abs(l-0.5)*2)
}

// Opacity returns a color with the specified opacity (alpha).
// Opacity should be in the range [0, 1], where 0 is transparent and 1 is opaque.
func Opacity(c Color, opacity float64) Color {
	return c.WithAlpha(clamp01(opacity))
}

// FadeOut decreases the opacity of a color by the specified amount.
func FadeOut(c Color, amount float64) Color {
	amount = clamp01(amount)
	return c.WithAlpha(c.Alpha() * (1 - amount))
}

// FadeIn increases the opacity of a color by the specified amount.
func FadeIn(c Color, amount float64) Color {
	amount = clamp01(amount)
	return c.WithAlpha(clamp01(c.Alpha() + amount*(1-c.Alpha())))
}

