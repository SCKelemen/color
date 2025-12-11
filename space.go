package color

// Space defines a color space with conversion to/from XYZ (the reference space).
// All inter-space conversions go through XYZ to ensure accuracy.
type Space interface {
	// Name returns the name of the color space (e.g., "sRGB", "Display P3", "OKLCH")
	Name() string
	
	// ToXYZ converts color values from this space to XYZ (linear, D65 white point).
	// channels: color values in this space's native format
	// Returns: x, y, z in CIE XYZ space
	ToXYZ(channels []float64) (x, y, z float64)
	
	// FromXYZ converts XYZ values to this color space.
	// x, y, z: CIE XYZ values (linear, D65 white point)
	// Returns: color values in this space's native format
	FromXYZ(x, y, z float64) []float64
	
	// Channels returns the number of color channels (3 for RGB, 4 for CMYK, etc.)
	Channels() int
	
	// ChannelNames returns the names of the channels (e.g., ["R", "G", "B"] or ["L", "C", "H"])
	ChannelNames() []string
}

// SpaceColor represents a color in a specific color space.
// This preserves the color space information and allows lossless operations.
type SpaceColor interface {
	Color // Still implements the base Color interface for compatibility
	
	// Space returns the color space this color is in
	Space() Space
	
	// Channels returns the color values in the native space
	Channels() []float64
	
	// ConvertTo converts this color to a different color space.
	// This is where data loss may occur (gamut clipping, etc.)
	ConvertTo(space Space) SpaceColor
	
	// ToRGBA converts to sRGB RGBA (explicit conversion, may lose data for wide-gamut colors)
	ToRGBA() *RGBA
}

// spaceColor is the concrete implementation of SpaceColor
type spaceColor struct {
	space  Space
	values []float64
	alpha  float64
}

// NewSpaceColor creates a new color in a specific color space.
func NewSpaceColor(space Space, channels []float64, alpha float64) SpaceColor {
	if len(channels) != space.Channels() {
		panic("channel count mismatch")
	}
	
	// Copy channels to avoid external mutation
	values := make([]float64, len(channels))
	copy(values, channels)
	
	return &spaceColor{
		space:  space,
		values: values,
		alpha:  clamp01(alpha),
	}
}

// Space implements SpaceColor
func (c *spaceColor) Space() Space {
	return c.space
}

// Channels implements SpaceColor
func (c *spaceColor) Channels() []float64 {
	// Return a copy to prevent mutation
	result := make([]float64, len(c.values))
	copy(result, c.values)
	return result
}

// Alpha implements Color
func (c *spaceColor) Alpha() float64 {
	return c.alpha
}

// WithAlpha implements Color
func (c *spaceColor) WithAlpha(alpha float64) Color {
	return &spaceColor{
		space:  c.space,
		values: c.values,
		alpha:  clamp01(alpha),
	}
}

// ConvertTo implements SpaceColor
func (c *spaceColor) ConvertTo(target Space) SpaceColor {
	// Convert through XYZ (the reference space)
	x, y, z := c.space.ToXYZ(c.values)
	targetChannels := target.FromXYZ(x, y, z)
	
	return NewSpaceColor(target, targetChannels, c.alpha)
}

// RGBA implements Color (converts to sRGB RGBA)
func (c *spaceColor) RGBA() (r, g, b, a float64) {
	rgba := c.ToRGBA()
	return rgba.R, rgba.G, rgba.B, rgba.A
}

// ToRGBA implements SpaceColor (explicit conversion to sRGB)
func (c *spaceColor) ToRGBA() *RGBA {
	// Convert through XYZ to sRGB
	x, y, z := c.space.ToXYZ(c.values)
	
	// Convert XYZ to linear sRGB
	linearR := x*3.2404542 - y*1.5371385 - z*0.4985314
	linearG := -x*0.9692660 + y*1.8760108 + z*0.0415560
	linearB := x*0.0556434 - y*0.2040259 + z*1.0572252
	
	// Apply sRGB gamma correction (using function from xyz.go)
	r := sRGBTransfer(linearR)
	g := sRGBTransfer(linearG)
	b := sRGBTransfer(linearB)
	
	return NewRGBA(
		clamp01(r),
		clamp01(g),
		clamp01(b),
		c.alpha,
	)
}

