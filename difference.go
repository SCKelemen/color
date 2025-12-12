package color

import "math"

// DeltaEOK calculates the perceptual difference between two colors using OKLAB space.
// This is a simpler and faster alternative to DeltaE2000, using Euclidean distance in OKLAB.
// Returns a value where 0 means identical colors, and larger values mean more different.
// Values < 0.02 are typically imperceptible, < 0.05 are barely perceptible.
func DeltaEOK(c1, c2 Color) float64 {
	oklab1 := ToOKLAB(c1)
	oklab2 := ToOKLAB(c2)

	dl := oklab1.L - oklab2.L
	da := oklab1.A - oklab2.A
	db := oklab1.B - oklab2.B

	return math.Sqrt(dl*dl + da*da + db*db)
}

// DeltaE76 calculates the perceptual difference using CIE76 formula (Euclidean distance in LAB).
// This is the original DeltaE formula, simpler but less accurate than DeltaE2000.
// Returns a value where 0 means identical colors, and larger values mean more different.
// Values < 1.0 are barely perceptible, 1-2 are small differences, > 2 are noticeable.
func DeltaE76(c1, c2 Color) float64 {
	lab1 := ToLAB(c1)
	lab2 := ToLAB(c2)

	dl := lab1.L - lab2.L
	da := lab1.A - lab2.A
	db := lab1.B - lab2.B

	return math.Sqrt(dl*dl + da*da + db*db)
}

// DeltaE2000 calculates the perceptual difference using the CIEDE2000 formula.
// This is the industry standard for color difference measurement.
// Returns a value where 0 means identical colors, and larger values mean more different.
// Values < 1.0 are barely perceptible, 1-2 are small differences, > 2 are noticeable.
func DeltaE2000(c1, c2 Color) float64 {
	lab1 := ToLAB(c1)
	lab2 := ToLAB(c2)

	// Convert to LCH for intermediate calculations
	l1, a1, b1 := lab1.L, lab1.A, lab1.B
	l2, a2, b2 := lab2.L, lab2.A, lab2.B

	// Calculate C1, C2 (chroma)
	c1_ := math.Sqrt(a1*a1 + b1*b1)
	c2_ := math.Sqrt(a2*a2 + b2*b2)

	// Calculate average C
	cAvg := (c1_ + c2_) / 2.0

	// Calculate G (adjustment factor)
	cAvg7 := cAvg * cAvg * cAvg * cAvg * cAvg * cAvg * cAvg
	g := 0.5 * (1 - math.Sqrt(cAvg7/(cAvg7+6103515625.0))) // 6103515625 = 25^7

	// Calculate adjusted a values
	a1Prime := a1 * (1 + g)
	a2Prime := a2 * (1 + g)

	// Calculate adjusted C' and h'
	c1Prime := math.Sqrt(a1Prime*a1Prime + b1*b1)
	c2Prime := math.Sqrt(a2Prime*a2Prime + b2*b2)

	h1Prime := 0.0
	if b1 != 0 || a1Prime != 0 {
		h1Prime = math.Atan2(b1, a1Prime) * 180.0 / math.Pi
		if h1Prime < 0 {
			h1Prime += 360
		}
	}

	h2Prime := 0.0
	if b2 != 0 || a2Prime != 0 {
		h2Prime = math.Atan2(b2, a2Prime) * 180.0 / math.Pi
		if h2Prime < 0 {
			h2Prime += 360
		}
	}

	// Calculate deltas
	deltaLPrime := l2 - l1
	deltaCPrime := c2Prime - c1Prime

	// Calculate deltaHPrime
	deltaHPrime := 0.0
	if c1Prime*c2Prime != 0 {
		deltaHPrime = h2Prime - h1Prime
		if deltaHPrime > 180 {
			deltaHPrime -= 360
		} else if deltaHPrime < -180 {
			deltaHPrime += 360
		}
	}
	deltaHPrime = 2 * math.Sqrt(c1Prime*c2Prime) * math.Sin(deltaHPrime*math.Pi/360.0)

	// Calculate averages for weighting factors
	lPrimeAvg := (l1 + l2) / 2.0
	cPrimeAvg := (c1Prime + c2Prime) / 2.0

	hPrimeAvg := 0.0
	if c1Prime*c2Prime != 0 {
		hPrimeAvg = (h1Prime + h2Prime) / 2.0
		if math.Abs(h1Prime-h2Prime) > 180 {
			if h1Prime+h2Prime < 360 {
				hPrimeAvg += 180
			} else {
				hPrimeAvg -= 180
			}
		}
	}

	// Calculate T (function of hue)
	t := 1 - 0.17*math.Cos((hPrimeAvg-30)*math.Pi/180.0) +
		0.24*math.Cos(2*hPrimeAvg*math.Pi/180.0) +
		0.32*math.Cos((3*hPrimeAvg+6)*math.Pi/180.0) -
		0.20*math.Cos((4*hPrimeAvg-63)*math.Pi/180.0)

	// Calculate rotation term
	cPrimeAvg7 := cPrimeAvg * cPrimeAvg * cPrimeAvg * cPrimeAvg * cPrimeAvg * cPrimeAvg * cPrimeAvg
	rc := 2 * math.Sqrt(cPrimeAvg7/(cPrimeAvg7+6103515625.0))

	deltaTheta := 30 * math.Exp(-((hPrimeAvg-275)/25)*((hPrimeAvg-275)/25))
	rt := -rc * math.Sin(2*deltaTheta*math.Pi/180.0)

	// Calculate weighting factors
	sl := 1 + (0.015*(lPrimeAvg-50)*(lPrimeAvg-50))/math.Sqrt(20+(lPrimeAvg-50)*(lPrimeAvg-50))
	sc := 1 + 0.045*cPrimeAvg
	sh := 1 + 0.015*cPrimeAvg*t

	// Calculate final DeltaE2000
	kL, kC, kH := 1.0, 1.0, 1.0 // Standard weighting factors

	deltaE := math.Sqrt(
		(deltaLPrime/(kL*sl))*(deltaLPrime/(kL*sl)) +
			(deltaCPrime/(kC*sc))*(deltaCPrime/(kC*sc)) +
			(deltaHPrime/(kH*sh))*(deltaHPrime/(kH*sh)) +
			rt*(deltaCPrime/(kC*sc))*(deltaHPrime/(kH*sh)))

	return deltaE
}
