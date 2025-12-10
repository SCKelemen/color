package color

import "math"

// HSL represents a color in the HSL color space.
// H is hue [0, 360), S is saturation [0, 1], L is lightness [0, 1].
type HSL struct {
	H, S, L, A float64
}

// NewHSL creates a new HSL color.
func NewHSL(h, s, l, a float64) *HSL {
	return &HSL{
		H: normalizeHue(h),
		S: clamp01(s),
		L: clamp01(l),
		A: clamp01(a),
	}
}

// RGBA converts HSL to RGBA.
func (c *HSL) RGBA() (r, g, b, a float64) {
	h := c.H / 360.0
	s := c.S
	l := c.L

	if s == 0 {
		// Achromatic (gray)
		return l, l, l, c.A
	}

	var q, p float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p = 2*l - q

	r = hueToRGB(p, q, h+1.0/3.0)
	g = hueToRGB(p, q, h)
	b = hueToRGB(p, q, h-1.0/3.0)

	return clamp01(r), clamp01(g), clamp01(b), c.A
}

// Alpha implements Color.
func (c *HSL) Alpha() float64 {
	return c.A
}

// WithAlpha implements Color.
func (c *HSL) WithAlpha(alpha float64) Color {
	return &HSL{H: c.H, S: c.S, L: c.L, A: clamp01(alpha)}
}

// hueToRGB is a helper function for HSL to RGB conversion.
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// ToHSL converts an RGBA color to HSL.
func ToHSL(c Color) *HSL {
	r, g, b, a := c.RGBA()

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min

	// Lightness
	l := (max + min) / 2.0

	// Saturation
	var s float64
	if delta == 0 {
		s = 0
	} else {
		if l < 0.5 {
			s = delta / (max + min)
		} else {
			s = delta / (2 - max - min)
		}
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

	return &HSL{H: normalizeHue(h), S: s, L: l, A: a}
}

