package color

import "math"

// XYZ represents a color in the CIE XYZ color space.
// This is an intermediate color space used for conversions.
type XYZ struct {
	X, Y, Z, A float64
}

// NewXYZ creates a new XYZ color.
func NewXYZ(x, y, z, a float64) *XYZ {
	return &XYZ{X: x, Y: y, Z: z, A: clamp01(a)}
}

// RGBA converts XYZ to RGBA (sRGB).
func (c *XYZ) RGBA() (r, g, b, a float64) {
	// Convert XYZ to linear RGB
	linearR := c.X*3.2404542 - c.Y*1.5371385 - c.Z*0.4985314
	linearG := -c.X*0.9692660 + c.Y*1.8760108 + c.Z*0.0415560
	linearB := c.X*0.0556434 - c.Y*0.2040259 + c.Z*1.0572252

	// Apply gamma correction to convert linear RGB to sRGB
	r = gammaCorrection(linearR)
	g = gammaCorrection(linearG)
	b = gammaCorrection(linearB)
	a = c.A

	return clamp01(r), clamp01(g), clamp01(b), clamp01(a)
}

// Alpha implements Color.
func (c *XYZ) Alpha() float64 {
	return c.A
}

// WithAlpha implements Color.
func (c *XYZ) WithAlpha(alpha float64) Color {
	return &XYZ{X: c.X, Y: c.Y, Z: c.Z, A: clamp01(alpha)}
}

// ToXYZ converts an RGBA color to XYZ.
func ToXYZ(c Color) *XYZ {
	r, g, b, a := c.RGBA()

	// Convert sRGB to linear RGB
	linearR := inverseGammaCorrection(r)
	linearG := inverseGammaCorrection(g)
	linearB := inverseGammaCorrection(b)

	// Convert linear RGB to XYZ
	x := linearR*0.4124564 + linearG*0.3575761 + linearB*0.1804375
	y := linearR*0.2126729 + linearG*0.7151522 + linearB*0.0721750
	z := linearR*0.0193339 + linearG*0.1191920 + linearB*0.9503041

	return &XYZ{X: x, Y: y, Z: z, A: a}
}

// gammaCorrection applies sRGB gamma correction.
func gammaCorrection(linear float64) float64 {
	if linear <= 0.0031308 {
		return 12.92 * linear
	}
	return 1.055*math.Pow(linear, 1.0/2.4) - 0.055
}

// inverseGammaCorrection reverses sRGB gamma correction.
func inverseGammaCorrection(srgb float64) float64 {
	if srgb <= 0.04045 {
		return srgb / 12.92
	}
	return math.Pow((srgb+0.055)/1.055, 2.4)
}

