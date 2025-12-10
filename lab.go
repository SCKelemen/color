package color

import "math"

// LAB represents a color in the CIE LAB color space.
// L is lightness [0, 100], A and B are color-opponent dimensions.
type LAB struct {
	L, A, B, A_ float64 // A_ is alpha to avoid conflict with method
}

// NewLAB creates a new LAB color.
func NewLAB(l, a, b, alpha float64) *LAB {
	return &LAB{
		L:  clamp(l, 0, 100),
		A:  a,
		B:  b,
		A_: clamp01(alpha),
	}
}

// RGBA converts LAB to RGBA via XYZ.
func (c *LAB) RGBA() (r, g, b, a float64) {
	xyz := c.toXYZ()
	return xyz.RGBA()
}

// Alpha implements Color.
func (c *LAB) Alpha() float64 {
	return c.A_
}

// WithAlpha implements Color.
func (c *LAB) WithAlpha(alpha float64) Color {
	return &LAB{L: c.L, A: c.A, B: c.B, A_: clamp01(alpha)}
}

// toXYZ converts LAB to XYZ.
func (c *LAB) toXYZ() *XYZ {
	// D65 white point
	const (
		xn = 0.95047
		yn = 1.00000
		zn = 1.08883
	)

	// Convert LAB to XYZ
	fy := (c.L + 16) / 116
	fx := c.A/500 + fy
	fz := fy - c.B/200

	// Calculate x, y, z
	var x, y, z float64
	if fx3 := fx * fx * fx; fx3 > 0.008856 {
		x = xn * fx3
	} else {
		x = (fx - 16.0/116.0) * 3 * 0.008856 * xn
	}

	if c.L > 8 {
		y = yn * math.Pow((c.L+16)/116, 3)
	} else {
		y = c.L / 903.3 * yn
	}

	if fz3 := fz * fz * fz; fz3 > 0.008856 {
		z = zn * fz3
	} else {
		z = (fz - 16.0/116.0) * 3 * 0.008856 * zn
	}

	return &XYZ{X: x, Y: y, Z: z, A: c.A_}
}

// ToLAB converts an RGBA color to LAB.
func ToLAB(c Color) *LAB {
	xyz := ToXYZ(c)
	return xyz.toLAB()
}

// toLAB converts XYZ to LAB.
func (c *XYZ) toLAB() *LAB {
	// D65 white point
	const (
		xn = 0.95047
		yn = 1.00000
		zn = 1.08883
	)

	// Normalize by white point
	x := c.X / xn
	y := c.Y / yn
	z := c.Z / zn

	// Convert to LAB
	fx := labF(x)
	fy := labF(y)
	fz := labF(z)

	l := 116*fy - 16
	a := 500 * (fx - fy)
	b := 200 * (fy - fz)

	return &LAB{L: l, A: a, B: b, A_: c.A}
}

// labF is the helper function for LAB conversion.
func labF(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3.0)
	}
	return (7.787*t + 16.0/116.0)
}

