package color

import "math"

// LCH represents a color in the LCH color space (polar representation of LAB).
// L is lightness [0, 100], C is chroma [0, ~132], H is hue [0, 360).
type LCH struct {
	L, C, H, A_ float64 // A_ is alpha to avoid conflict with method
}

// NewLCH creates a new LCH color.
func NewLCH(l, c, h, alpha float64) *LCH {
	return &LCH{
		L:  clamp(l, 0, 100),
		C:  math.Max(0, c),
		H:  normalizeHue(h),
		A_: clamp01(alpha),
	}
}

// RGBA converts LCH to RGBA via LAB.
func (c *LCH) RGBA() (r, g, b, a float64) {
	lab := c.toLAB()
	return lab.RGBA()
}

// Alpha implements Color.
func (c *LCH) Alpha() float64 {
	return c.A_
}

// WithAlpha implements Color.
func (c *LCH) WithAlpha(alpha float64) Color {
	return &LCH{L: c.L, C: c.C, H: c.H, A_: clamp01(alpha)}
}

// toLAB converts LCH to LAB.
func (c *LCH) toLAB() *LAB {
	rad := c.H * math.Pi / 180
	a := c.C * math.Cos(rad)
	b := c.C * math.Sin(rad)
	return &LAB{L: c.L, A: a, B: b, A_: c.A_}
}

// ToLCH converts an RGBA color to LCH.
func ToLCH(c Color) *LCH {
	lab := ToLAB(c)
	return lab.toLCH()
}

// toLCH converts LAB to LCH.
func (c *LAB) toLCH() *LCH {
	c_ := math.Sqrt(c.A*c.A + c.B*c.B)
	h := math.Atan2(c.B, c.A) * 180 / math.Pi
	h = normalizeHue(h)
	return &LCH{L: c.L, C: c_, H: h, A_: c.A_}
}

// normalizeHue normalizes hue to [0, 360).
func normalizeHue(h float64) float64 {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	return h
}

