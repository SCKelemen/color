package color

import "math"

// SRGBSpace represents the sRGB color space (D65 white point, sRGB primaries, sRGB transfer function)
var SRGBSpace Space = &rgbSpace{
	name:           "sRGB",
	xyzToRGBMatrix: [9]float64{3.2404542, -1.5371385, -0.4985314, -0.9692660, 1.8760108, 0.0415560, 0.0556434, -0.2040259, 1.0572252},
	rgbToXYZMatrix: [9]float64{0.4124564, 0.3575761, 0.1804375, 0.2126729, 0.7151522, 0.0721750, 0.0193339, 0.1191920, 0.9503041},
	transferFunc:   sRGBTransfer,
	inverseTransferFunc: sRGBInverseTransfer,
}

// SRGBLinearSpace represents linear sRGB (no gamma encoding)
var SRGBLinearSpace Space = &rgbSpace{
	name:           "sRGB-linear",
	xyzToRGBMatrix: [9]float64{3.2404542, -1.5371385, -0.4985314, -0.9692660, 1.8760108, 0.0415560, 0.0556434, -0.2040259, 1.0572252},
	rgbToXYZMatrix: [9]float64{0.4124564, 0.3575761, 0.1804375, 0.2126729, 0.7151522, 0.0721750, 0.0193339, 0.1191920, 0.9503041},
	transferFunc:   linearTransfer,
	inverseTransferFunc: linearInverseTransfer,
}

// rgbSpace implements Space for RGB color spaces
type rgbSpace struct {
	name               string
	xyzToRGBMatrix     [9]float64 // Matrix to convert XYZ to linear RGB
	rgbToXYZMatrix     [9]float64 // Matrix to convert linear RGB to XYZ
	transferFunc       func(float64) float64
	inverseTransferFunc func(float64) float64
}

func (s *rgbSpace) Name() string {
	return s.name
}

func (s *rgbSpace) Channels() int {
	return 3
}

func (s *rgbSpace) ChannelNames() []string {
	return []string{"R", "G", "B"}
}

func (s *rgbSpace) ToXYZ(channels []float64) (x, y, z float64) {
	if len(channels) != 3 {
		panic("RGB space requires 3 channels")
	}
	
	// Apply inverse transfer function to get linear RGB
	r := s.inverseTransferFunc(channels[0])
	g := s.inverseTransferFunc(channels[1])
	b := s.inverseTransferFunc(channels[2])
	
	// Convert linear RGB to XYZ using matrix
	m := s.rgbToXYZMatrix
	x = m[0]*r + m[1]*g + m[2]*b
	y = m[3]*r + m[4]*g + m[5]*b
	z = m[6]*r + m[7]*g + m[8]*b
	
	return x, y, z
}

func (s *rgbSpace) FromXYZ(x, y, z float64) []float64 {
	// Convert XYZ to linear RGB using matrix
	m := s.xyzToRGBMatrix
	linearR := m[0]*x + m[1]*y + m[2]*z
	linearG := m[3]*x + m[4]*y + m[5]*z
	linearB := m[6]*x + m[7]*y + m[8]*z
	
	// Apply transfer function
	r := s.transferFunc(linearR)
	g := s.transferFunc(linearG)
	b := s.transferFunc(linearB)
	
	return []float64{r, g, b}
}

// sRGBTransfer applies sRGB gamma correction
func sRGBTransfer(linear float64) float64 {
	if linear <= 0.0031308 {
		return 12.92 * linear
	}
	return 1.055*math.Pow(linear, 1.0/2.4) - 0.055
}

// sRGBInverseTransfer reverses sRGB gamma correction
func sRGBInverseTransfer(encoded float64) float64 {
	if encoded <= 0.04045 {
		return encoded / 12.92
	}
	return math.Pow((encoded+0.055)/1.055, 2.4)
}

// linearTransfer is the identity function (no encoding)
func linearTransfer(linear float64) float64 {
	return linear
}

// linearInverseTransfer is the identity function (no decoding)
func linearInverseTransfer(encoded float64) float64 {
	return encoded
}

