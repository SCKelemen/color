package color

import "math"

// LUV represents a color in the CIE LUV color space.
// L is lightness [0, 100], U and V are color-opponent dimensions.
// LUV is an alternative to LAB that's better suited for emissive displays.
type LUV struct {
	L, U, V, A_ float64 // A_ is alpha to avoid conflict with method
}

// NewLUV creates a new LUV color.
func NewLUV(l, u, v, alpha float64) *LUV {
	return &LUV{
		L:  clamp(l, 0, 100),
		U:  u,
		V:  v,
		A_: clamp01(alpha),
	}
}

// RGBA converts LUV to RGBA via XYZ.
func (c *LUV) RGBA() (r, g, b, a float64) {
	xyz := c.toXYZ()
	return xyz.RGBA()
}

// Alpha implements Color.
func (c *LUV) Alpha() float64 {
	return c.A_
}

// WithAlpha implements Color.
func (c *LUV) WithAlpha(alpha float64) Color {
	return &LUV{L: c.L, U: c.U, V: c.V, A_: clamp01(alpha)}
}

// toXYZ converts LUV to XYZ.
func (c *LUV) toXYZ() *XYZ {
	// D65 white point
	const (
		xn = 0.95047
		yn = 1.00000
		zn = 1.08883
	)

	// Reference white u', v'
	unPrime := (4 * xn) / (xn + 15*yn + 3*zn)
	vnPrime := (9 * yn) / (xn + 15*yn + 3*zn)

	// Calculate Y
	var y float64
	if c.L > 8 {
		y = yn * math.Pow((c.L+16)/116, 3)
	} else {
		y = yn * c.L / 903.3
	}

	// Calculate u', v'
	var uPrime, vPrime float64
	if c.L == 0 {
		uPrime = unPrime
		vPrime = vnPrime
	} else {
		uPrime = c.U/(13*c.L) + unPrime
		vPrime = c.V/(13*c.L) + vnPrime
	}

	// Convert to XYZ
	x := y * (9 * uPrime) / (4 * vPrime)
	z := y * (12 - 3*uPrime - 20*vPrime) / (4 * vPrime)

	return &XYZ{X: x, Y: y, Z: z, A: c.A_}
}

// ToLUV converts a Color to LUV.
func ToLUV(c Color) *LUV {
	xyz := ToXYZ(c)
	return xyz.toLUV()
}

// toLUV converts XYZ to LUV.
func (c *XYZ) toLUV() *LUV {
	// D65 white point
	const (
		xn = 0.95047
		yn = 1.00000
		zn = 1.08883
	)

	// Calculate L*
	yr := c.Y / yn
	var l float64
	if yr > 0.008856 {
		l = 116*math.Pow(yr, 1.0/3.0) - 16
	} else {
		l = 903.3 * yr
	}

	// Calculate u', v'
	denominator := c.X + 15*c.Y + 3*c.Z
	var uPrime, vPrime float64
	if denominator != 0 {
		uPrime = (4 * c.X) / denominator
		vPrime = (9 * c.Y) / denominator
	}

	// Reference white u', v'
	unPrime := (4 * xn) / (xn + 15*yn + 3*zn)
	vnPrime := (9 * yn) / (xn + 15*yn + 3*zn)

	// Calculate u*, v*
	u := 13 * l * (uPrime - unPrime)
	v := 13 * l * (vPrime - vnPrime)

	return &LUV{L: l, U: u, V: v, A_: c.A}
}

// LCHuv represents a color in the CIE LCHuv color space (cylindrical LUV).
// L is lightness [0, 100], C is chroma, H is hue [0, 360).
type LCHuv struct {
	L, C, H, A_ float64 // A_ is alpha to avoid conflict with method
}

// NewLCHuv creates a new LCHuv color.
func NewLCHuv(l, c, h, alpha float64) *LCHuv {
	return &LCHuv{
		L:  clamp(l, 0, 100),
		C:  math.Max(0, c),
		H:  NormalizeHue(h),
		A_: clamp01(alpha),
	}
}

// RGBA converts LCHuv to RGBA via LUV.
func (c *LCHuv) RGBA() (r, g, b, a float64) {
	luv := c.toLUV()
	return luv.RGBA()
}

// Alpha implements Color.
func (c *LCHuv) Alpha() float64 {
	return c.A_
}

// WithAlpha implements Color.
func (c *LCHuv) WithAlpha(alpha float64) Color {
	return &LCHuv{L: c.L, C: c.C, H: c.H, A_: clamp01(alpha)}
}

// toLUV converts LCHuv to LUV.
func (c *LCHuv) toLUV() *LUV {
	rad := c.H * math.Pi / 180
	u := c.C * math.Cos(rad)
	v := c.C * math.Sin(rad)
	return &LUV{L: c.L, U: u, V: v, A_: c.A_}
}

// ToLCHuv converts a Color to LCHuv.
func ToLCHuv(c Color) *LCHuv {
	luv := ToLUV(c)
	return luv.toLCHuv()
}

// toLCHuv converts LUV to LCHuv.
func (c *LUV) toLCHuv() *LCHuv {
	chroma := math.Sqrt(c.U*c.U + c.V*c.V)
	h := math.Atan2(c.V, c.U) * 180 / math.Pi
	h = NormalizeHue(h)
	return &LCHuv{L: c.L, C: chroma, H: h, A_: c.A_}
}
