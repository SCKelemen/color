package color

import "math"

// HSV represents a color in the HSV (HSB) color space.
// H is hue [0, 360), S is saturation [0, 1], V is value/brightness [0, 1].
type HSV struct {
	H, S, V, A float64
}

// NewHSV creates a new HSV color.
func NewHSV(h, s, v, a float64) *HSV {
	return &HSV{
		H: normalizeHue(h),
		S: clamp01(s),
		V: clamp01(v),
		A: clamp01(a),
	}
}

// RGBA converts HSV to RGBA.
func (c *HSV) RGBA() (r, g, b, a float64) {
	h := c.H / 60.0
	s := c.S
	v := c.V

	c_ := v * s
	x := c_ * (1 - math.Abs(math.Mod(h, 2)-1))
	m := v - c_

	var r1, g1, b1 float64
	switch {
	case h < 1:
		r1, g1, b1 = c_, x, 0
	case h < 2:
		r1, g1, b1 = x, c_, 0
	case h < 3:
		r1, g1, b1 = 0, c_, x
	case h < 4:
		r1, g1, b1 = 0, x, c_
	case h < 5:
		r1, g1, b1 = x, 0, c_
	default:
		r1, g1, b1 = c_, 0, x
	}

	r = r1 + m
	g = g1 + m
	b = b1 + m

	return clamp01(r), clamp01(g), clamp01(b), c.A
}

// Alpha implements Color.
func (c *HSV) Alpha() float64 {
	return c.A
}

// WithAlpha implements Color.
func (c *HSV) WithAlpha(alpha float64) Color {
	return &HSV{H: c.H, S: c.S, V: c.V, A: clamp01(alpha)}
}

// ToHSV converts an RGBA color to HSV.
func ToHSV(c Color) *HSV {
	r, g, b, a := c.RGBA()

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min

	// Value
	v := max

	// Saturation
	var s float64
	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	// Hue
	var h float64
	if delta == 0 {
		h = 0
	} else if max == r {
		h = 60 * (((g - b) / delta))
		if h < 0 {
			h += 360
		}
	} else if max == g {
		h = 60*((b-r)/delta) + 120
	} else {
		h = 60*((r-g)/delta) + 240
	}

	return &HSV{H: normalizeHue(h), S: s, V: v, A: a}
}

