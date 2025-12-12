package color

import "math"

// HWB represents a color in the HWB color space.
// H is hue [0, 360), W is whiteness [0, 1], B is blackness [0, 1].
// HWB is more intuitive than HSL - whiteness adds white, blackness adds black.
// This color space is part of CSS Color Level 4.
type HWB struct {
	H, W, B, A float64
}

// NewHWB creates a new HWB color.
// H is in [0, 360), W and B are in [0, 1].
// If W+B > 1, they will be normalized to sum to 1.
func NewHWB(h, w, b, a float64) *HWB {
	h = NormalizeHue(h)
	w = clamp01(w)
	b = clamp01(b)

	// Normalize whiteness and blackness if their sum exceeds 1
	sum := w + b
	if sum > 1 {
		w = w / sum
		b = b / sum
	}

	return &HWB{
		H: h,
		W: w,
		B: b,
		A: clamp01(a),
	}
}

// RGBA converts HWB to RGBA.
func (c *HWB) RGBA() (r, g, b, a float64) {
	// If whiteness + blackness = 1, return gray
	sum := c.W + c.B
	if sum >= 1 {
		gray := c.W / sum
		return gray, gray, gray, c.A
	}

	// Convert hue to base RGB (pure hue)
	h := c.H / 60.0
	x := 1 - math.Abs(math.Mod(h, 2)-1)

	var r1, g1, b1 float64
	switch int(h) {
	case 0:
		r1, g1, b1 = 1, x, 0
	case 1:
		r1, g1, b1 = x, 1, 0
	case 2:
		r1, g1, b1 = 0, 1, x
	case 3:
		r1, g1, b1 = 0, x, 1
	case 4:
		r1, g1, b1 = x, 0, 1
	case 5:
		r1, g1, b1 = 1, 0, x
	default:
		r1, g1, b1 = 1, 0, 0
	}

	// Apply whiteness and blackness
	// Formula: RGB = (RGB_pure * (1 - W - B)) + W
	r = r1*(1-c.W-c.B) + c.W
	g = g1*(1-c.W-c.B) + c.W
	b = b1*(1-c.W-c.B) + c.W
	a = c.A

	return clamp01(r), clamp01(g), clamp01(b), clamp01(a)
}

// Alpha implements Color.
func (c *HWB) Alpha() float64 {
	return c.A
}

// WithAlpha implements Color.
func (c *HWB) WithAlpha(alpha float64) Color {
	return &HWB{H: c.H, W: c.W, B: c.B, A: clamp01(alpha)}
}

// ToHWB converts a Color to HWB.
func ToHWB(c Color) *HWB {
	r, g, b, a := c.RGBA()

	// Find min and max of RGB
	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	chroma := max - min

	// Calculate hue (same as HSL/HSV)
	h := 0.0
	if chroma != 0 {
		switch max {
		case r:
			h = 60 * math.Mod((g-b)/chroma, 6)
		case g:
			h = 60 * ((b-r)/chroma + 2)
		case b:
			h = 60 * ((r-g)/chroma + 4)
		}
	}
	h = NormalizeHue(h)

	// Whiteness is the minimum RGB value
	w := min

	// Blackness is 1 - maximum RGB value
	blackness := 1 - max

	return &HWB{
		H: h,
		W: w,
		B: blackness,
		A: a,
	}
}
