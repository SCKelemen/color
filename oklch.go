package color

import "math"

// OKLCH represents a color in the OKLCH color space (polar representation of OKLAB).
// L is perceived lightness [0, 1], C is chroma [0, ~0.4], H is hue [0, 360).
type OKLCH struct {
	L, C, H, A_ float64 // A_ is alpha to avoid conflict with method
}

// NewOKLCH creates a new OKLCH color.
func NewOKLCH(l, c, h, alpha float64) *OKLCH {
	return &OKLCH{
		L:  clamp01(l),
		C:  math.Max(0, c),
		H:  normalizeHue(h),
		A_: clamp01(alpha),
	}
}

// RGBA converts OKLCH to RGBA via OKLAB.
func (c *OKLCH) RGBA() (r, g, b, a float64) {
	oklab := c.toOKLAB()
	return oklab.RGBA()
}

// Alpha implements Color.
func (c *OKLCH) Alpha() float64 {
	return c.A_
}

// WithAlpha implements Color.
func (c *OKLCH) WithAlpha(alpha float64) Color {
	return &OKLCH{L: c.L, C: c.C, H: c.H, A_: clamp01(alpha)}
}

// toOKLAB converts OKLCH to OKLAB.
func (c *OKLCH) toOKLAB() *OKLAB {
	rad := c.H * math.Pi / 180
	a := c.C * math.Cos(rad)
	b := c.C * math.Sin(rad)
	return &OKLAB{L: c.L, A: a, B: b, A_: c.A_}
}

// ToOKLCH converts an RGBA color to OKLCH.
func ToOKLCH(c Color) *OKLCH {
	oklab := ToOKLAB(c)
	return oklab.toOKLCH()
}

// toOKLCH converts OKLAB to OKLCH.
func (c *OKLAB) toOKLCH() *OKLCH {
	c_ := math.Sqrt(c.A*c.A + c.B*c.B)
	h := math.Atan2(c.B, c.A) * 180 / math.Pi
	h = normalizeHue(h)
	return &OKLCH{L: c.L, C: c_, H: h, A_: c.A_}
}

