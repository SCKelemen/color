package color

// Color represents a color in any color space with an alpha channel.
// All color spaces must implement conversion to and from RGBA.
type Color interface {
	// RGBA returns the red, green, blue, and alpha values
	// in the range [0, 1].
	RGBA() (r, g, b, a float64)
	
	// Alpha returns the alpha channel value in the range [0, 1].
	Alpha() float64
	
	// WithAlpha returns a new color with the specified alpha value.
	WithAlpha(alpha float64) Color
}

// RGBA represents a color in the RGB color space with alpha.
type RGBA struct {
	R, G, B, A float64
}

// RGB creates a new RGBA color with full opacity.
func RGB(r, g, b float64) *RGBA {
	return &RGBA{R: clamp01(r), G: clamp01(g), B: clamp01(b), A: 1.0}
}

// RGBA creates a new RGBA color with the specified alpha.
func NewRGBA(r, g, b, a float64) *RGBA {
	return &RGBA{R: clamp01(r), G: clamp01(g), B: clamp01(b), A: clamp01(a)}
}

// RGBA implements Color.
func (c *RGBA) RGBA() (r, g, b, a float64) {
	return c.R, c.G, c.B, c.A
}

// Alpha implements Color.
func (c *RGBA) Alpha() float64 {
	return c.A
}

// WithAlpha implements Color.
func (c *RGBA) WithAlpha(alpha float64) Color {
	return &RGBA{R: c.R, G: c.G, B: c.B, A: clamp01(alpha)}
}

// String returns a string representation of the color.
func (c *RGBA) String() string {
	return RGBToHex(c)
}

// clamp01 clamps a value to the range [0, 1].
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// clamp clamps a value to the range [min, max].
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

