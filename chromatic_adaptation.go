package color

// Chromatic adaptation transforms for converting between different white points.
// This is essential for accurate color conversions when source and destination
// color spaces use different white points (e.g., D50 vs D65).

// Common white points in XYZ
var (
	// D65 white point (standard for most RGB spaces, sRGB, Display P3, Rec.2020)
	// Corresponds to  6500K daylight
	whiteD65 = [3]float64{0.95047, 1.00000, 1.08883}

	// D50 white point (used by ProPhoto RGB, ICC LAB)
	// Corresponds to 5000K daylight (horizon light)
	whiteD50 = [3]float64{0.96422, 1.00000, 0.82521}
)

// Bradford chromatic adaptation matrix
// This matrix is used to transform XYZ values for chromatic adaptation
var bradfordMatrix = [9]float64{
	0.8951000, 0.2664000, -0.1614000,
	-0.7502000, 1.7135000, 0.0367000,
	0.0389000, -0.0685000, 1.0296000,
}

// Inverse Bradford matrix
var bradfordMatrixInv = [9]float64{
	0.9869929, -0.1470543, 0.1599627,
	0.4323053, 0.5183603, 0.0492912,
	-0.0085287, 0.0400428, 0.9684867,
}

// AdaptD65ToD50 adapts XYZ values from D65 white point to D50 white point.
// This is used when converting to color spaces that use D50 (like ProPhoto RGB).
func AdaptD65ToD50(x, y, z float64) (float64, float64, float64) {
	return adaptWhitePoint(x, y, z, whiteD65, whiteD50)
}

// AdaptD50ToD65 adapts XYZ values from D50 white point to D65 white point.
// This is used when converting from color spaces that use D50 (like ProPhoto RGB).
func AdaptD50ToD65(x, y, z float64) (float64, float64, float64) {
	return adaptWhitePoint(x, y, z, whiteD50, whiteD65)
}

// adaptWhitePoint performs chromatic adaptation using the Bradford transform.
// This adapts XYZ values from one white point to another.
func adaptWhitePoint(x, y, z float64, sourceWhite, destWhite [3]float64) (float64, float64, float64) {
	// Convert XYZ to LMS (cone response) using Bradford matrix
	m := bradfordMatrix
	lSource := m[0]*x + m[1]*y + m[2]*z
	mSource := m[3]*x + m[4]*y + m[5]*z
	sSource := m[6]*x + m[7]*y + m[8]*z

	// Get LMS for source and destination white points
	lWhiteSrc := m[0]*sourceWhite[0] + m[1]*sourceWhite[1] + m[2]*sourceWhite[2]
	mWhiteSrc := m[3]*sourceWhite[0] + m[4]*sourceWhite[1] + m[5]*sourceWhite[2]
	sWhiteSrc := m[6]*sourceWhite[0] + m[7]*sourceWhite[1] + m[8]*sourceWhite[2]

	lWhiteDst := m[0]*destWhite[0] + m[1]*destWhite[1] + m[2]*destWhite[2]
	mWhiteDst := m[3]*destWhite[0] + m[4]*destWhite[1] + m[5]*destWhite[2]
	sWhiteDst := m[6]*destWhite[0] + m[7]*destWhite[1] + m[8]*destWhite[2]

	// Scale the LMS values
	lAdapted := lSource * (lWhiteDst / lWhiteSrc)
	mAdapted := mSource * (mWhiteDst / mWhiteSrc)
	sAdapted := sSource * (sWhiteDst / sWhiteSrc)

	// Convert back to XYZ using inverse Bradford matrix
	mi := bradfordMatrixInv
	xAdapted := mi[0]*lAdapted + mi[1]*mAdapted + mi[2]*sAdapted
	yAdapted := mi[3]*lAdapted + mi[4]*mAdapted + mi[5]*sAdapted
	zAdapted := mi[6]*lAdapted + mi[7]*mAdapted + mi[8]*sAdapted

	return xAdapted, yAdapted, zAdapted
}
