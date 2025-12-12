package color

// DisplayP3Space represents Display P3 color space
var DisplayP3Space Space = &rgbSpace{
	name: "display-p3",
	xyzToRGBMatrix: [9]float64{
		2.493496911941425, -0.9313836179191239, -0.40271078445071684,
		-0.8294889695615747, 1.7626640603183463, 0.023624685841943577,
		0.03584583024378447, -0.07617238926804182, 0.9568845240076872,
	},
	rgbToXYZMatrix: [9]float64{
		0.4865709486482162, 0.26566769316909306, 0.1982172852343625,
		0.2289745640697488, 0.6917385218365064, 0.079286914093745,
		0.000000000000000, 0.04511338185890264, 1.043944368900976,
	},
	transferFunc:        sRGBTransfer, // Display P3 uses sRGB transfer function
	inverseTransferFunc: sRGBInverseTransfer,
	whitePoint:          WhiteD65,
}

// A98RGBSpace represents Adobe RGB 1998 color space
var A98RGBSpace Space = &rgbSpace{
	name: "a98-rgb",
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
	transferFunc:        gammaTransferFunc(2.2),
	inverseTransferFunc: gammaInverseTransferFunc(2.2),
	whitePoint:          WhiteD65,
}

// ProPhotoRGBSpace represents ProPhoto RGB color space (D50 white point)
var ProPhotoRGBSpace Space = &rgbSpace{
	name: "prophoto-rgb",
	xyzToRGBMatrix: [9]float64{
		1.3459433, -0.2556075, -0.0511118,
		-0.5445989, 1.5081673, 0.0205351,
		0.0000000, 0.0000000, 1.2118128,
	},
	rgbToXYZMatrix: [9]float64{
		0.7976749, 0.1351917, 0.0313534,
		0.2880402, 0.7118741, 0.0000857,
		0.0000000, 0.0000000, 0.8252100,
	},
	transferFunc:        gammaTransferFunc(1.8),
	inverseTransferFunc: gammaInverseTransferFunc(1.8),
	whitePoint:          WhiteD50, // ProPhoto RGB uses D50 white point
}

// Rec2020Space represents Rec. 2020 color space (UHDTV)
var Rec2020Space Space = &rgbSpace{
	name: "rec2020",
	xyzToRGBMatrix: [9]float64{
		1.7166511, -0.3556708, -0.2533663,
		-0.6666844, 1.6164812, 0.0157685,
		0.0176399, -0.0427706, 0.9421031,
	},
	rgbToXYZMatrix: [9]float64{
		0.6369580483012914, 0.14461690358620832, 0.1688809751641721,
		0.262704531669281, 0.6779980715188708, 0.05930171646986196,
		0.000000000000000, 0.028072693049087428, 1.060985057710791,
	},
	transferFunc:        rec2020Transfer,
	inverseTransferFunc: rec2020InverseTransfer,
	whitePoint:          WhiteD65,
}

// Rec709Space represents Rec. 709 color space (HDTV)
// Very similar to sRGB but technically has different primaries
var Rec709Space Space = &rgbSpace{
	name: "rec709",
	xyzToRGBMatrix: [9]float64{
		3.2404542, -1.5371385, -0.4985314,
		-0.9692660, 1.8760108, 0.0415560,
		0.0556434, -0.2040259, 1.0572252,
	},
	rgbToXYZMatrix: [9]float64{
		0.4124564, 0.3575761, 0.1804375,
		0.2126729, 0.7151522, 0.0721750,
		0.0193339, 0.1191920, 0.9503041,
	},
	transferFunc:        rec709Transfer,
	inverseTransferFunc: rec709InverseTransfer,
	whitePoint:          WhiteD65,
}

// DCIP3Space represents DCI-P3 color space (Digital Cinema)
// Similar to Display P3 but uses gamma 2.6 and different white point
var DCIP3Space Space = &rgbSpace{
	name: "dci-p3",
	xyzToRGBMatrix: [9]float64{
		2.493496911941425, -0.9313836179191239, -0.40271078445071684,
		-0.8294889695615747, 1.7626640603183463, 0.023624685841943577,
		0.03584583024378447, -0.07617238926804182, 0.9568845240076872,
	},
	rgbToXYZMatrix: [9]float64{
		0.4865709486482162, 0.26566769316909306, 0.1982172852343625,
		0.2289745640697488, 0.6917385218365064, 0.079286914093745,
		0.000000000000000, 0.04511338185890264, 1.043944368900976,
	},
	transferFunc:        gammaTransferFunc(2.6), // DCI-P3 uses gamma 2.6
	inverseTransferFunc: gammaInverseTransferFunc(2.6),
	whitePoint:          WhiteD65,
}

