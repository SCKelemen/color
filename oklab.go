package color

import "math"

// OKLAB represents a color in the OKLAB color space.
// L is perceived lightness [0, 1], A and B are color-opponent dimensions.
type OKLAB struct {
	L, A, B, A_ float64 // A_ is alpha to avoid conflict with method
}

// NewOKLAB creates a new OKLAB color.
func NewOKLAB(l, a, b, alpha float64) *OKLAB {
	return &OKLAB{
		L:  clamp01(l),
		A:  a,
		B:  b,
		A_: clamp01(alpha),
	}
}

// RGBA converts OKLAB to RGBA via linear RGB.
func (c *OKLAB) RGBA() (r, g, b, a float64) {
	// Convert OKLAB to linear RGB
	l := c.L + 0.3963377774*c.A + 0.2158037573*c.B
	m := c.L - 0.1055613458*c.A - 0.0638541728*c.B
	s := c.L - 0.0894841775*c.A - 1.2914855480*c.B

	l3 := l * l * l
	m3 := m * m * m
	s3 := s * s * s

	linearR := +4.0767416621*l3 - 3.3077115913*m3 + 0.2309699292*s3
	linearG := -1.2684380046*l3 + 2.6097574011*m3 - 0.3413193965*s3
	linearB := -0.0041960863*l3 - 0.7034186147*m3 + 1.7076147010*s3

	// Apply gamma correction to convert linear RGB to sRGB
	r = gammaCorrection(linearR)
	g = gammaCorrection(linearG)
	b = gammaCorrection(linearB)
	a = c.A_

	return clamp01(r), clamp01(g), clamp01(b), clamp01(a)
}

// Alpha implements Color.
func (c *OKLAB) Alpha() float64 {
	return c.A_
}

// WithAlpha implements Color.
func (c *OKLAB) WithAlpha(alpha float64) Color {
	return &OKLAB{L: c.L, A: c.A, B: c.B, A_: clamp01(alpha)}
}

// ToOKLAB converts an RGBA color to OKLAB.
func ToOKLAB(c Color) *OKLAB {
	r, g, b, _ := c.RGBA()

	// Convert sRGB to linear RGB
	linearR := inverseGammaCorrection(r)
	linearG := inverseGammaCorrection(g)
	linearB := inverseGammaCorrection(b)

	// Convert linear RGB to OKLAB
	l := 0.4122214708*linearR + 0.5363325363*linearG + 0.0514459929*linearB
	m := 0.2119034982*linearR + 0.6806995451*linearG + 0.1073969566*linearB
	s := 0.0883024619*linearR + 0.2817188376*linearG + 0.6299787005*linearB

	l3 := math.Cbrt(l)
	m3 := math.Cbrt(m)
	s3 := math.Cbrt(s)

	okl := 0.2104542553*l3 + 0.7936177850*m3 - 0.0040720468*s3
	oka := 1.9779984951*l3 - 2.4285922050*m3 + 0.4505937099*s3
	okb := 0.0259040371*l3 + 0.7827717662*m3 - 0.8086757660*s3

	_, _, _, a := c.RGBA()
	return &OKLAB{L: okl, A: oka, B: okb, A_: a}
}

