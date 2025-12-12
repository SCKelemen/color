package color

import stdcolor "image/color"

// ToStdColor converts a color from this library to the standard library's color.Color interface.
// The standard library uses uint32 values in the range [0, 65535] (16-bit per channel).
//
// Example:
//
//	import (
//	    "image/color"
//	    "github.com/SCKelemen/color"
//	)
//
//	myColor := color.RGB(1.0, 0.5, 0.0)
//	stdColor := color.ToStdColor(myColor)
//	// Now you can use stdColor with image processing libraries
func ToStdColor(c Color) stdcolor.Color {
	r, g, b, a := c.RGBA()
	return &stdRGBA{
		R: uint32(r * 65535),
		G: uint32(g * 65535),
		B: uint32(b * 65535),
		A: uint32(a * 65535),
	}
}

// FromStdColor converts a standard library color.Color to this library's Color interface.
// This allows you to use colors from image processing libraries with this library.
//
// Example:
//
//	import (
//	    "image/color"
//	    "github.com/SCKelemen/color"
//	)
//
//	stdColor := color.RGBA{R: 255, G: 128, B: 0, A: 255}
//	myColor := color.FromStdColor(stdColor)
//	// Now you can use myColor with this library's functions
//	oklch := color.ToOKLCH(myColor)
func FromStdColor(c stdcolor.Color) Color {
	r, g, b, a := c.RGBA()
	// Standard library returns alpha-premultiplied values in range [0, 65535]
	// We need to convert to non-premultiplied [0, 1] range
	alpha := float64(a) / 65535.0
	if alpha == 0 {
		return NewRGBA(0, 0, 0, 0)
	}
	return NewRGBA(
		float64(r)/65535.0/alpha,
		float64(g)/65535.0/alpha,
		float64(b)/65535.0/alpha,
		alpha,
	)
}

// stdRGBA is a wrapper that implements image/color.Color interface.
type stdRGBA struct {
	R, G, B, A uint32
}

// RGBA implements image/color.Color interface.
// Returns alpha-premultiplied values in range [0, 65535].
func (c *stdRGBA) RGBA() (r, g, b, a uint32) {
	// Alpha-premultiply the color components
	alpha := float64(c.A) / 65535.0
	r = uint32(float64(c.R) * alpha)
	g = uint32(float64(c.G) * alpha)
	b = uint32(float64(c.B) * alpha)
	a = c.A
	return
}

