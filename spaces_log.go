package color

import "math"

// LOG Color Space Transfer Functions
// These implement logarithmic transfer functions used in professional video/cinema cameras

// Canon C-Log Transfer Functions
// C-Log is Canon's logarithmic color space for cinema cameras
// Reference: Canon C-Log specification

func cLogTransfer(linear float64) float64 {
	// Canon C-Log transfer function
	// Simplified formula for C-Log (Cinema Gamut / C-Log)
	if linear < 0 {
		return 0
	}

	// C-Log is a pure logarithmic curve
	// Encoded value can be negative for very dark values (this is correct)
	return 0.529136*math.Log10(10.1596*linear+1) + 0.0730597
}

func cLogInverseTransfer(encoded float64) float64 {
	// Canon C-Log inverse transfer function
	linear := (math.Pow(10, (encoded-0.0730597)/0.529136) - 1) / 10.1596

	// Clamp to non-negative
	if linear < 0 {
		return 0
	}
	return linear
}

// Sony S-Log3 Transfer Functions
// S-Log3 is Sony's logarithmic color space for cinema cameras (newest version)
// Reference: Sony S-Log3 specification

func sLog3Transfer(linear float64) float64 {
	// S-Log3 transfer function: linear to S-Log3
	if linear < 0 {
		linear = 0
	}

	// S-Log3 parameters
	const a = 0.01125000
	const b = 0.42188671
	const c = 0.42188671
	const k = 261.5

	if linear >= a {
		return (420.0 + math.Log10((linear+0.01)/0.18)*(c*k)) / 1023.0
	}

	return (linear * (171.2102946929 - 95.0) / 0.01125000 + 95.0) / 1023.0
}

func sLog3InverseTransfer(encoded float64) float64 {
	// S-Log3 inverse: S-Log3 to linear
	encoded = encoded * 1023.0

	const a = 0.01125000
	const b = 0.42188671
	const c = 0.42188671
	const k = 261.5

	if encoded >= 171.2102946929 {
		return (math.Pow(10, ((encoded-420.0)/(c*k))) * 0.18) - 0.01
	}

	return ((encoded - 95.0) / (171.2102946929 - 95.0) * 0.01125000)
}

// Panasonic V-Log Transfer Functions
// V-Log is Panasonic's logarithmic color space for cinema cameras
// Reference: Panasonic V-Log specification

func vLogTransfer(linear float64) float64 {
	// V-Log transfer function: linear to V-Log
	if linear < 0 {
		return 0
	}

	// V-Log parameters
	const cut1 = 0.01
	const cut2 = 0.181
	const b = 0.00873
	const c = 0.241514
	const d = 0.598206

	if linear < cut1 {
		return 5.6 * linear + 0.125
	}

	return c*math.Log10(linear+b) + d
}

func vLogInverseTransfer(encoded float64) float64 {
	// V-Log inverse: V-Log to linear
	const cut1 = 0.01
	const cut2 = 0.181
	const b = 0.00873
	const c = 0.241514
	const d = 0.598206

	if encoded < cut2 {
		return (encoded - 0.125) / 5.6
	}

	return math.Pow(10, (encoded-d)/c) - b
}

// Arri LogC Transfer Functions (V3, EI 800)
// LogC is Arri's logarithmic color space for Alexa cameras
// Reference: Arri LogC specification (Version 3, Exposure Index 800)

func arriLogCTransfer(linear float64) float64 {
	// Arri LogC (V3, EI 800) transfer function: linear to LogC
	if linear < 0 {
		return 0
	}

	// LogC V3 EI 800 parameters
	const cut = 0.010591
	const a = 5.555556
	const b = 0.052272
	const c = 0.247190
	const d = 0.385537
	const e = 5.367655
	const f = 0.092809

	if linear > cut {
		return c*math.Log10(a*linear+b) + d
	}

	return e*linear + f
}

func arriLogCInverseTransfer(encoded float64) float64 {
	// Arri LogC inverse: LogC to linear
	const cut = 0.010591
	const e = 5.367655
	const f = 0.092809
	const cut2 = e*cut + f
	const a = 5.555556
	const b = 0.052272
	const c = 0.247190
	const d = 0.385537

	if encoded > cut2 {
		return (math.Pow(10, (encoded-d)/c) - b) / a
	}

	return (encoded - f) / e
}

// Red Log3G10 Transfer Functions
// Log3G10 is Red's logarithmic color space for Red cameras
// Reference: Red Log3G10 specification

func redLog3G10Transfer(linear float64) float64 {
	// Red Log3G10 transfer function: linear to Log3G10
	if linear < 0 {
		return 0
	}

	// Log3G10 parameters
	const a = 0.224282
	const b = 155.975327
	const c = 0.01

	if linear < 0 {
		return 0
	}

	return a*math.Log10(linear*b+1) + c
}

func redLog3G10InverseTransfer(encoded float64) float64 {
	// Red Log3G10 inverse: Log3G10 to linear
	const a = 0.224282
	const b = 155.975327
	const c = 0.01

	return (math.Pow(10, (encoded-c)/a) - 1) / b
}

// Blackmagic Film Transfer Functions
// BMDFilm is Blackmagic's logarithmic color space
// Reference: Blackmagic Film specification

func bmdFilmTransfer(linear float64) float64 {
	// BMDFilm transfer function: linear to BMDFilm
	if linear < 0.0 {
		return 0.0
	}

	// Simplified curve based on log base 2
	const linearRange = 0.005
	const logOffset = 0.075

	if linear < linearRange {
		return linear * 8.0
	}

	return math.Log2(linear+linearRange)*0.07 + logOffset + 0.5
}

func bmdFilmInverseTransfer(encoded float64) float64 {
	// BMDFilm inverse: BMDFilm to linear
	const linearRange = 0.005
	const logOffset = 0.075
	const linearBreak = 0.04

	if encoded < linearBreak {
		return encoded / 8.0
	}

	return math.Pow(2, (encoded-logOffset-0.5)/0.07) - linearRange
}

// LOG Color Space Definitions
// These use the same primaries as their base RGB spaces but with LOG transfer functions

var (
	// CLogSpace represents Canon C-Log color space
	// Uses Canon Cinema Gamut primaries with C-Log transfer
	CLogSpace Space = &rgbSpace{
		name: "c-log",
		// Canon Cinema Gamut primaries (approximately Rec. 2020)
		xyzToRGBMatrix: [9]float64{
			1.9624274, -0.6105343, -0.3413404,
			-0.9787684, 1.9161415, 0.0334540,
			0.0286869, -0.1406752, 1.3487655,
		},
		rgbToXYZMatrix: [9]float64{
			0.5766690429101305, 0.1855582379065463, 0.1882286462349947,
			0.29734497525053605, 0.6273635662554661, 0.07529145849399788,
			0.02703136138641234, 0.07068885253582723, 0.9913375368376388,
		},
		transferFunc:        cLogTransfer,
		inverseTransferFunc: cLogInverseTransfer,
		whitePoint:          WhiteD65,
	}

	// SLog3Space represents Sony S-Log3 color space
	// Uses Sony S-Gamut3 primaries with S-Log3 transfer
	SLog3Space Space = &rgbSpace{
		name: "s-log3",
		// Sony S-Gamut3 primaries (wide gamut)
		xyzToRGBMatrix: [9]float64{
			1.9624274, -0.6105343, -0.3413404,
			-0.9787684, 1.9161415, 0.0334540,
			0.0286869, -0.1406752, 1.3487655,
		},
		rgbToXYZMatrix: [9]float64{
			0.7064, 0.1288, 0.1213,
			0.2709, 0.7869, -0.0578,
			-0.0096, 0.0045, 1.1156,
		},
		transferFunc:        sLog3Transfer,
		inverseTransferFunc: sLog3InverseTransfer,
		whitePoint:          WhiteD65,
	}

	// VLogSpace represents Panasonic V-Log color space
	// Uses V-Gamut primaries with V-Log transfer
	VLogSpace Space = &rgbSpace{
		name: "v-log",
		// V-Gamut primaries (wide gamut)
		xyzToRGBMatrix: [9]float64{
			1.5890, -0.3130, -0.1802,
			-0.5340, 1.3960, 0.0950,
			-0.0110, -0.0640, 1.1570,
		},
		rgbToXYZMatrix: [9]float64{
			0.7790, 0.0780, 0.1020,
			0.3270, 0.7230, -0.0500,
			0.0010, 0.0510, 0.8390,
		},
		transferFunc:        vLogTransfer,
		inverseTransferFunc: vLogInverseTransfer,
		whitePoint:          WhiteD65,
	}

	// ArriLogCSpace represents Arri LogC (V3, EI 800) color space
	// Uses Arri Wide Gamut primaries with LogC transfer
	ArriLogCSpace Space = &rgbSpace{
		name: "arri-logc",
		// Arri Wide Gamut primaries
		xyzToRGBMatrix: [9]float64{
			1.789066, -0.482534, -0.200076,
			-0.639849, 1.396400, 0.194432,
			-0.041532, 0.082335, 1.040648,
		},
		rgbToXYZMatrix: [9]float64{
			0.6380, 0.2140, 0.0970,
			0.2910, 0.8240, -0.1150,
			0.0020, -0.0380, 1.0910,
		},
		transferFunc:        arriLogCTransfer,
		inverseTransferFunc: arriLogCInverseTransfer,
		whitePoint:          WhiteD65,
	}

	// RedLog3G10Space represents Red Log3G10 color space
	// Uses RedWideGamutRGB primaries with Log3G10 transfer
	RedLog3G10Space Space = &rgbSpace{
		name: "red-log3g10",
		// Red Wide Gamut RGB primaries
		xyzToRGBMatrix: [9]float64{
			1.7827, -0.4969, -0.1768,
			-0.6702, 1.4985, 0.1323,
			-0.0530, 0.0403, 1.0672,
		},
		rgbToXYZMatrix: [9]float64{
			0.7347, 0.1596, 0.0366,
			0.2653, 0.8404, -0.1057,
			0.0000, -0.0000, 1.0690,
		},
		transferFunc:        redLog3G10Transfer,
		inverseTransferFunc: redLog3G10InverseTransfer,
		whitePoint:          WhiteD65,
	}

	// BMDFilmSpace represents Blackmagic Film color space
	// Uses Blackmagic Wide Gamut primaries with BMDFilm transfer
	BMDFilmSpace Space = &rgbSpace{
		name: "bmd-film",
		// Blackmagic Wide Gamut (approximately Rec. 2020)
		xyzToRGBMatrix: [9]float64{
			1.9624274, -0.6105343, -0.3413404,
			-0.9787684, 1.9161415, 0.0334540,
			0.0286869, -0.1406752, 1.3487655,
		},
		rgbToXYZMatrix: [9]float64{
			0.5766690429101305, 0.1855582379065463, 0.1882286462349947,
			0.29734497525053605, 0.6273635662554661, 0.07529145849399788,
			0.02703136138641234, 0.07068885253582723, 0.9913375368376388,
		},
		transferFunc:        bmdFilmTransfer,
		inverseTransferFunc: bmdFilmInverseTransfer,
		whitePoint:          WhiteD65,
	}
)
