package color

import (
	"math"
	"strings"
)

// RGBColorSpace represents a wide-gamut RGB color space with its own
// conversion matrix and transfer function.
type RGBColorSpace struct {
	// Name of the color space
	Name string
	
	// Matrix to convert from XYZ to linear RGB
	// Format: [rX, rY, rZ, gX, gY, gZ, bX, bY, bZ]
	XYZToRGBMatrix [9]float64
	
	// Matrix to convert from linear RGB to XYZ
	// Format: [Xr, Xg, Xb, Yr, Yg, Yb, Zr, Zg, Zb]
	RGBToXYZMatrix [9]float64
	
	// Transfer function: converts linear RGB to encoded RGB
	TransferFunc func(float64) float64
	
	// Inverse transfer function: converts encoded RGB to linear RGB
	InverseTransferFunc func(float64) float64
}

// ConvertXYZToRGB converts XYZ to this RGB color space.
func (cs *RGBColorSpace) ConvertXYZToRGB(xyz *XYZ) *RGBA {
	// Apply matrix: linear RGB = matrix * XYZ
	m := cs.XYZToRGBMatrix
	linearR := m[0]*xyz.X + m[1]*xyz.Y + m[2]*xyz.Z
	linearG := m[3]*xyz.X + m[4]*xyz.Y + m[5]*xyz.Z
	linearB := m[6]*xyz.X + m[7]*xyz.Y + m[8]*xyz.Z
	
	// Apply transfer function
	r := cs.TransferFunc(linearR)
	g := cs.TransferFunc(linearG)
	b := cs.TransferFunc(linearB)
	
	return NewRGBA(clamp01(r), clamp01(g), clamp01(b), xyz.A)
}

// ConvertRGBToXYZ converts RGB in this color space to XYZ.
func (cs *RGBColorSpace) ConvertRGBToXYZ(r, g, b, a float64) *XYZ {
	// Apply inverse transfer function to get linear RGB
	linearR := cs.InverseTransferFunc(r)
	linearG := cs.InverseTransferFunc(g)
	linearB := cs.InverseTransferFunc(b)
	
	// Apply matrix: XYZ = matrix * linear RGB
	m := cs.RGBToXYZMatrix
	x := m[0]*linearR + m[1]*linearG + m[2]*linearB
	y := m[3]*linearR + m[4]*linearG + m[5]*linearB
	z := m[6]*linearR + m[7]*linearG + m[8]*linearB
	
	return NewXYZ(x, y, z, a)
}

// sRGB transfer function (same as gammaCorrection)
func sRGBTransfer(linear float64) float64 {
	if linear <= 0.0031308 {
		return 12.92 * linear
	}
	return 1.055*math.Pow(linear, 1.0/2.4) - 0.055
}

// sRGB inverse transfer function (same as inverseGammaCorrection)
func sRGBInverseTransfer(encoded float64) float64 {
	if encoded <= 0.04045 {
		return encoded / 12.92
	}
	return math.Pow((encoded+0.055)/1.055, 2.4)
}

// Linear transfer function (no encoding)
func linearTransfer(linear float64) float64 {
	return linear
}

// Linear inverse transfer function (no decoding)
func linearInverseTransfer(encoded float64) float64 {
	return encoded
}

// Gamma transfer function with specified gamma value
func gammaTransferFunc(gamma float64) func(float64) float64 {
	return func(linear float64) float64 {
		if linear < 0 {
			return 0
		}
		return math.Pow(linear, 1.0/gamma)
	}
}

// Gamma inverse transfer function
func gammaInverseTransferFunc(gamma float64) func(float64) float64 {
	return func(encoded float64) float64 {
		if encoded < 0 {
			return 0
		}
		return math.Pow(encoded, gamma)
	}
}

// Define color spaces
var (
	// sRGB color space (already supported, but defined here for consistency)
	sRGBSpace = &RGBColorSpace{
		Name: "srgb",
		// D65 white point, sRGB primaries
		XYZToRGBMatrix: [9]float64{
			3.2404542, -1.5371385, -0.4985314,
			-0.9692660, 1.8760108, 0.0415560,
			0.0556434, -0.2040259, 1.0572252,
		},
		RGBToXYZMatrix: [9]float64{
			0.4124564, 0.3575761, 0.1804375,
			0.2126729, 0.7151522, 0.0721750,
			0.0193339, 0.1191920, 0.9503041,
		},
		TransferFunc:        sRGBTransfer,
		InverseTransferFunc: sRGBInverseTransfer,
	}

	// sRGB-linear: linear sRGB (no gamma encoding)
	sRGBLinearSpace = &RGBColorSpace{
		Name: "srgb-linear",
		// Same matrices as sRGB, but linear transfer
		XYZToRGBMatrix: [9]float64{
			3.2404542, -1.5371385, -0.4985314,
			-0.9692660, 1.8760108, 0.0415560,
			0.0556434, -0.2040259, 1.0572252,
		},
		RGBToXYZMatrix: [9]float64{
			0.4124564, 0.3575761, 0.1804375,
			0.2126729, 0.7151522, 0.0721750,
			0.0193339, 0.1191920, 0.9503041,
		},
		TransferFunc:        linearTransfer,
		InverseTransferFunc: linearInverseTransfer,
	}

	// Display P3: Wide gamut RGB (D65 white point)
	displayP3Space = &RGBColorSpace{
		Name: "display-p3",
		// D65 white point, Display P3 primaries
		XYZToRGBMatrix: [9]float64{
			2.493496911941425, -0.9313836179191239, -0.40271078445071684,
			-0.8294889695615747, 1.7626640603183463, 0.023624685841943577,
			0.03584583024378447, -0.07617238926804182, 0.9568845240076872,
		},
		RGBToXYZMatrix: [9]float64{
			0.4865709486482162, 0.26566769316909306, 0.1982172852343625,
			0.2289745640697488, 0.6917385218365064, 0.079286914093745,
			0.000000000000000, 0.04511338185890264, 1.043944368900976,
		},
		TransferFunc:        sRGBTransfer, // Display P3 uses sRGB transfer function
		InverseTransferFunc: sRGBInverseTransfer,
	}

	// Adobe RGB 1998: D65 white point
	a98RGBSpace = &RGBColorSpace{
		Name: "a98-rgb",
		// D65 white point, Adobe RGB 1998 primaries
		XYZToRGBMatrix: [9]float64{
			1.9624274, -0.6105343, -0.3413404,
			-0.9787684, 1.9161415, 0.0334540,
			0.0286869, -0.1406752, 1.3487655,
		},
		RGBToXYZMatrix: [9]float64{
			0.5766690429101305, 0.1855582379065463, 0.1882286462349947,
			0.29734497525053605, 0.6273635662554661, 0.07529145849399788,
			0.02703136138641234, 0.07068885253582723, 0.9913375368376388,
		},
		TransferFunc:        gammaTransferFunc(2.2), // Adobe RGB uses gamma 2.2
		InverseTransferFunc: gammaInverseTransferFunc(2.2),
	}

	// ProPhoto RGB: D50 white point
	proPhotoRGBSpace = &RGBColorSpace{
		Name: "prophoto-rgb",
		// D50 white point, ProPhoto RGB primaries
		XYZToRGBMatrix: [9]float64{
			1.3459433, -0.2556075, -0.0511118,
			-0.5445989, 1.5081673, 0.0205351,
			0.0000000, 0.0000000, 1.2118128,
		},
		RGBToXYZMatrix: [9]float64{
			0.7976749, 0.1351917, 0.0313534,
			0.2880402, 0.7118741, 0.0000857,
			0.0000000, 0.0000000, 0.8252100,
		},
		TransferFunc:        gammaTransferFunc(1.8), // ProPhoto RGB uses gamma 1.8
		InverseTransferFunc: gammaInverseTransferFunc(1.8),
	}

	// Rec. 2020: Wide gamut for UHDTV (D65 white point)
	rec2020Space = &RGBColorSpace{
		Name: "rec2020",
		// D65 white point, Rec. 2020 primaries
		XYZToRGBMatrix: [9]float64{
			1.7166511, -0.3556708, -0.2533663,
			-0.6666844, 1.6164812, 0.0157685,
			0.0176399, -0.0427706, 0.9421031,
		},
		RGBToXYZMatrix: [9]float64{
			0.6369580483012914, 0.14461690358620832, 0.1688809751641721,
			0.262704531669281, 0.6779980715188708, 0.05930171646986196,
			0.000000000000000, 0.028072693049087428, 1.060985057710791,
		},
		TransferFunc:        rec2020Transfer, // Rec. 2020 uses a different transfer function
		InverseTransferFunc: rec2020InverseTransfer,
	}
)

// Rec. 2020 transfer function (PQ-like, but simplified to gamma 2.4 for compatibility)
func rec2020Transfer(linear float64) float64 {
	// Rec. 2020 actually uses a more complex transfer function, but for compatibility
	// with CSS, we use a simplified version. Full implementation would use:
	// if linear < beta: return alpha * linear
	// else: return (1 + delta) * linear^(1/gamma) - delta
	// For now, use gamma 2.4 as approximation
	if linear < 0 {
		return 0
	}
	return math.Pow(linear, 1.0/2.4)
}

// Rec. 2020 inverse transfer function
func rec2020InverseTransfer(encoded float64) float64 {
	if encoded < 0 {
		return 0
	}
	return math.Pow(encoded, 2.4)
}

// Rec. 709 transfer function (similar to sRGB but with different constants)
func rec709Transfer(linear float64) float64 {
	// Rec. 709 uses essentially the same transfer function as sRGB
	// The differences are negligible for practical purposes
	if linear <= 0.0031308 {
		return 12.92 * linear
	}
	return 1.055*math.Pow(linear, 1.0/2.4) - 0.055
}

// Rec. 709 inverse transfer function
func rec709InverseTransfer(encoded float64) float64 {
	if encoded <= 0.04045 {
		return encoded / 12.92
	}
	return math.Pow((encoded+0.055)/1.055, 2.4)
}

// getRGBColorSpace returns the RGBColorSpace for a given name.
func getRGBColorSpace(name string) *RGBColorSpace {
	switch strings.ToLower(name) {
	case "srgb":
		return sRGBSpace
	case "srgb-linear":
		return sRGBLinearSpace
	case "display-p3", "display-p3-d65":
		return displayP3Space
	case "a98-rgb", "a98rgb":
		return a98RGBSpace
	case "prophoto-rgb", "prophoto":
		return proPhotoRGBSpace
	case "rec2020", "rec-2020":
		return rec2020Space
	default:
		return nil
	}
}

