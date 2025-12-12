package color

import (
	"math"
	"sort"
)

// Gradient generates a gradient between two colors in a perceptually uniform color space.
// Steps is the number of colors to generate (including start and end).
// The gradient is computed in OKLCH space for perceptually uniform results.
func Gradient(start, end Color, steps int) []Color {
	if steps <= 0 {
		return []Color{}
	}
	if steps == 1 {
		return []Color{start}
	}

	result := make([]Color, steps)
	for i := 0; i < steps; i++ {
		weight := float64(i) / float64(steps-1)
		result[i] = MixOKLCH(start, end, weight)
	}

	return result
}

// GradientInSpace generates a gradient in a specific color space.
// This allows you to control which color space is used for interpolation.
func GradientInSpace(start, end Color, steps int, space GradientSpace) []Color {
	if steps <= 0 {
		return []Color{}
	}
	if steps == 1 {
		return []Color{start}
	}

	result := make([]Color, steps)
	for i := 0; i < steps; i++ {
		weight := float64(i) / float64(steps-1)
		result[i] = MixInSpace(start, end, weight, space)
	}

	return result
}

// GradientSpace specifies which color space to use for gradient interpolation.
type GradientSpace int

const (
	// GradientRGB interpolates in RGB space (fast but not perceptually uniform)
	GradientRGB GradientSpace = iota
	// GradientHSL interpolates in HSL space
	GradientHSL
	// GradientLAB interpolates in CIE LAB space
	GradientLAB
	// GradientOKLAB interpolates in OKLAB space (perceptually uniform)
	GradientOKLAB
	// GradientLCH interpolates in CIE LCH space
	GradientLCH
	// GradientOKLCH interpolates in OKLCH space (perceptually uniform, recommended)
	GradientOKLCH
)

// HueInterpolation specifies how to interpolate hue values in cylindrical color spaces.
type HueInterpolation int

const (
	// HueShorter interpolates hue using the shortest path around the color wheel (default)
	HueShorter HueInterpolation = iota
	// HueLonger interpolates hue using the longer path around the color wheel
	HueLonger
	// HueIncreasing interpolates hue in the direction of increasing values
	HueIncreasing
	// HueDecreasing interpolates hue in the direction of decreasing values
	HueDecreasing
)

// MixInSpace mixes two colors in the specified color space.
func MixInSpace(c1, c2 Color, weight float64, space GradientSpace) Color {
	weight = clamp01(weight)

	switch space {
	case GradientRGB:
		return Mix(c1, c2, weight)
	case GradientHSL:
		return mixHSL(c1, c2, weight)
	case GradientLAB:
		return mixLAB(c1, c2, weight)
	case GradientOKLAB:
		return mixOKLAB(c1, c2, weight)
	case GradientLCH:
		return mixLCH(c1, c2, weight)
	case GradientOKLCH:
		return MixOKLCH(c1, c2, weight)
	default:
		return MixOKLCH(c1, c2, weight) // Default to OKLCH
	}
}

// mixHSL mixes colors in HSL space.
func mixHSL(c1, c2 Color, weight float64) Color {
	hsl1 := ToHSL(c1)
	hsl2 := ToHSL(c2)

	// Interpolate HSL components
	h := interpolateHue(hsl1.H, hsl2.H, weight, HueShorter)
	s := hsl1.S*(1-weight) + hsl2.S*weight
	l := hsl1.L*(1-weight) + hsl2.L*weight
	a := hsl1.A*(1-weight) + hsl2.A*weight

	return NewHSL(h, s, l, a)
}

// mixLAB mixes colors in CIE LAB space.
func mixLAB(c1, c2 Color, weight float64) Color {
	lab1 := ToLAB(c1)
	lab2 := ToLAB(c2)

	l := lab1.L*(1-weight) + lab2.L*weight
	a := lab1.A*(1-weight) + lab2.A*weight
	b := lab1.B*(1-weight) + lab2.B*weight
	alpha := lab1.Alpha()*(1-weight) + lab2.Alpha()*weight

	return NewLAB(l, a, b, alpha)
}

// mixOKLAB mixes colors in OKLAB space.
func mixOKLAB(c1, c2 Color, weight float64) Color {
	oklab1 := ToOKLAB(c1)
	oklab2 := ToOKLAB(c2)

	l := oklab1.L*(1-weight) + oklab2.L*weight
	a := oklab1.A*(1-weight) + oklab2.A*weight
	b := oklab1.B*(1-weight) + oklab2.B*weight
	alpha := oklab1.Alpha()*(1-weight) + oklab2.Alpha()*weight

	return NewOKLAB(l, a, b, alpha)
}

// mixLCH mixes colors in CIE LCH space.
func mixLCH(c1, c2 Color, weight float64) Color {
	lch1 := ToLCH(c1)
	lch2 := ToLCH(c2)

	l := lch1.L*(1-weight) + lch2.L*weight
	c := lch1.C*(1-weight) + lch2.C*weight
	h := interpolateHue(lch1.H, lch2.H, weight, HueShorter)
	alpha := lch1.Alpha()*(1-weight) + lch2.Alpha()*weight

	return NewLCH(l, c, h, alpha)
}

// GradientStop represents a color stop in a multistop gradient.
// Position should be in the range [0, 1], where 0 is the start and 1 is the end.
type GradientStop struct {
	Color    Color
	Position float64
}

// GradientMultiStop generates a gradient with multiple color stops.
// Stops should be sorted by position (0 to 1). If not sorted, they will be sorted automatically.
// Steps is the total number of colors to generate.
//
// Example:
//
//	stops := []GradientStop{
//	    {Color: color.RGB(1, 0, 0), Position: 0.0},   // Red at start
//	    {Color: color.RGB(1, 1, 0), Position: 0.5}, // Yellow in middle
//	    {Color: color.RGB(0, 0, 1), Position: 1.0}, // Blue at end
//	}
//	gradient := color.GradientMultiStop(stops, 20, color.GradientOKLCH)
func GradientMultiStop(stops []GradientStop, steps int, space GradientSpace) []Color {
	if len(stops) == 0 {
		return []Color{}
	}
	if len(stops) == 1 {
		result := make([]Color, steps)
		for i := range result {
			result[i] = stops[0].Color
		}
		return result
	}

	// Sort stops by position
	sortedStops := make([]GradientStop, len(stops))
	copy(sortedStops, stops)
	sort.Slice(sortedStops, func(i, j int) bool {
		return sortedStops[i].Position < sortedStops[j].Position
	})

	// Ensure first stop is at 0 and last is at 1
	if sortedStops[0].Position > 0 {
		sortedStops = append([]GradientStop{{Color: sortedStops[0].Color, Position: 0}}, sortedStops...)
	}
	if sortedStops[len(sortedStops)-1].Position < 1 {
		lastColor := sortedStops[len(sortedStops)-1].Color
		sortedStops = append(sortedStops, GradientStop{Color: lastColor, Position: 1})
	}

	result := make([]Color, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		result[i] = interpolateMultiStop(sortedStops, t, space)
	}

	return result
}

// interpolateMultiStop interpolates a color at position t using the given stops.
func interpolateMultiStop(stops []GradientStop, t float64, space GradientSpace) Color {
	// Clamp t to [0, 1]
	t = clamp01(t)

	// Find the two stops to interpolate between
	for i := 0; i < len(stops)-1; i++ {
		if t >= stops[i].Position && t <= stops[i+1].Position {
			// Normalize t to [0, 1] between these two stops
			segmentStart := stops[i].Position
			segmentEnd := stops[i+1].Position
			if segmentEnd == segmentStart {
				return stops[i].Color
			}
			localT := (t - segmentStart) / (segmentEnd - segmentStart)
			return MixInSpace(stops[i].Color, stops[i+1].Color, localT, space)
		}
	}

	// Fallback (shouldn't happen)
	return stops[len(stops)-1].Color
}

// EasingFunction defines a function that maps a linear progress [0, 1] to an eased progress [0, 1].
// Users can create custom easing functions by implementing this interface.
//
// Requirements:
//   - Input t should be in range [0, 1]
//   - Output should be in range [0, 1]
//   - Should map 0 to 0 and 1 to 1
//
// Example of creating a custom easing function:
//
//	// Custom exponential easing
//	myEasing := color.EasingFunction(func(t float64) float64 {
//	    return 1 - math.Pow(1-t, 3) // Ease-out cubic
//	})
//
//	// Use it in a gradient
//	gradient := color.GradientWithEasing(red, blue, 20, color.GradientOKLCH, myEasing)
type EasingFunction func(t float64) float64

// Easing functions for non-linear gradients
var (
	// EaseLinear is the default linear easing (no change).
	EaseLinear EasingFunction = func(t float64) float64 { return t }

	// EaseInQuad provides a quadratic ease-in curve.
	EaseInQuad EasingFunction = func(t float64) float64 { return t * t }

	// EaseOutQuad provides a quadratic ease-out curve.
	EaseOutQuad EasingFunction = func(t float64) float64 { return t * (2 - t) }

	// EaseInOutQuad provides a quadratic ease-in-out curve.
	EaseInOutQuad EasingFunction = func(t float64) float64 {
		if t < 0.5 {
			return 2 * t * t
		}
		return -1 + (4-2*t)*t
	}

	// EaseInCubic provides a cubic ease-in curve.
	EaseInCubic EasingFunction = func(t float64) float64 { return t * t * t }

	// EaseOutCubic provides a cubic ease-out curve.
	EaseOutCubic EasingFunction = func(t float64) float64 {
		t--
		return t*t*t + 1
	}

	// EaseInOutCubic provides a cubic ease-in-out curve.
	EaseInOutCubic EasingFunction = func(t float64) float64 {
		if t < 0.5 {
			return 4 * t * t * t
		}
		t = 2*t - 2
		return t*t*t/2 + 1
	}

	// EaseInSine provides a sinusoidal ease-in curve.
	EaseInSine EasingFunction = func(t float64) float64 {
		return 1 - math.Cos(t*math.Pi/2)
	}

	// EaseOutSine provides a sinusoidal ease-out curve.
	EaseOutSine EasingFunction = func(t float64) float64 {
		return math.Sin(t * math.Pi / 2)
	}

	// EaseInOutSine provides a sinusoidal ease-in-out curve.
	EaseInOutSine EasingFunction = func(t float64) float64 {
		return -(math.Cos(math.Pi*t) - 1) / 2
	}
)

// GradientWithEasing generates a gradient with an easing function applied.
// The easing function transforms the linear interpolation progress for non-linear gradients.
//
// Example:
//
//	red := color.RGB(1, 0, 0)
//	blue := color.RGB(0, 0, 1)
//	// Create a gradient that starts slow and speeds up (ease-in)
//	gradient := color.GradientWithEasing(red, blue, 20, color.GradientOKLCH, color.EaseInQuad)
func GradientWithEasing(start, end Color, steps int, space GradientSpace, easing EasingFunction) []Color {
	if steps <= 0 {
		return []Color{}
	}
	if steps == 1 {
		return []Color{start}
	}

	if easing == nil {
		easing = EaseLinear
	}

	result := make([]Color, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		easedT := easing(t)
		result[i] = MixInSpace(start, end, easedT, space)
	}

	return result
}

// GradientMultiStopWithEasing generates a multistop gradient with easing applied.
// The easing function is applied to the overall gradient progress.
func GradientMultiStopWithEasing(stops []GradientStop, steps int, space GradientSpace, easing EasingFunction) []Color {
	if len(stops) == 0 {
		return []Color{}
	}
	if len(stops) == 1 {
		result := make([]Color, steps)
		for i := range result {
			result[i] = stops[0].Color
		}
		return result
	}

	if easing == nil {
		easing = EaseLinear
	}

	// Sort stops by position
	sortedStops := make([]GradientStop, len(stops))
	copy(sortedStops, stops)
	sort.Slice(sortedStops, func(i, j int) bool {
		return sortedStops[i].Position < sortedStops[j].Position
	})

	// Ensure first stop is at 0 and last is at 1
	if sortedStops[0].Position > 0 {
		sortedStops = append([]GradientStop{{Color: sortedStops[0].Color, Position: 0}}, sortedStops...)
	}
	if sortedStops[len(sortedStops)-1].Position < 1 {
		lastColor := sortedStops[len(sortedStops)-1].Color
		sortedStops = append(sortedStops, GradientStop{Color: lastColor, Position: 1})
	}

	result := make([]Color, steps)
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		easedT := easing(t)
		result[i] = interpolateMultiStop(sortedStops, easedT, space)
	}

	return result
}

// interpolateHue interpolates hue using the specified interpolation method.
func interpolateHue(h1, h2, weight float64, method HueInterpolation) float64 {
	dh := h2 - h1

	switch method {
	case HueShorter:
		// Take the shortest path around the color wheel
		if math.Abs(dh) > 180 {
			if dh > 0 {
				dh -= 360
			} else {
				dh += 360
			}
		}
	case HueLonger:
		// Take the longer path around the color wheel
		if math.Abs(dh) <= 180 {
			if dh > 0 {
				dh -= 360
			} else {
				dh += 360
			}
		}
	case HueIncreasing:
		// Always go in increasing direction
		if dh < 0 {
			dh += 360
		}
	case HueDecreasing:
		// Always go in decreasing direction
		if dh > 0 {
			dh -= 360
		}
	}

	return normalizeHue(h1 + dh*weight)
}

