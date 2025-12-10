package color

import (
	"testing"
)

func TestGradient(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Generate 5-step gradient
	gradient := Gradient(red, blue, 5)

	if len(gradient) != 5 {
		t.Errorf("Gradient length = %d, want 5", len(gradient))
	}

	// First should be red
	r1, g1, b1, _ := gradient[0].RGBA()
	if !floatEqual(r1, 1) || !floatEqual(g1, 0) || !floatEqual(b1, 0) {
		t.Errorf("First color = RGB(%v, %v, %v), want RGB(1, 0, 0)", r1, g1, b1)
	}

	// Last should be blue
	r2, g2, b2, _ := gradient[4].RGBA()
	if !floatEqual(r2, 0) || !floatEqual(g2, 0) || !floatEqual(b2, 1) {
		t.Errorf("Last color = RGB(%v, %v, %v), want RGB(0, 0, 1)", r2, g2, b2)
	}

	// Middle should be a mix (purple-ish)
	r3, g3, b3, _ := gradient[2].RGBA()
	// Should have both red and blue components (purple-ish)
	if r3 < 0.2 || b3 < 0.2 {
		t.Errorf("Middle color = RGB(%v, %v, %v), expected purple-ish (both red and blue)", r3, g3, b3)
	}
}

func TestGradientInSpace(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Test different color spaces
	spaces := []GradientSpace{
		GradientRGB,
		GradientHSL,
		GradientLAB,
		GradientOKLAB,
		GradientLCH,
		GradientOKLCH,
	}

	for _, space := range spaces {
		t.Run(space.String(), func(t *testing.T) {
			gradient := GradientInSpace(red, blue, 3, space)
			if len(gradient) != 3 {
				t.Errorf("Gradient length = %d, want 3", len(gradient))
			}

			// All colors should be valid
			for i, c := range gradient {
				r, g, b, a := c.RGBA()
				if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
					t.Errorf("Gradient[%d] = RGB(%v, %v, %v, %v), want valid color", i, r, g, b, a)
				}
			}
		})
	}
}

func TestGradientMultiStop(t *testing.T) {
	red := RGB(1, 0, 0)
	yellow := RGB(1, 1, 0)
	blue := RGB(0, 0, 1)

	stops := []GradientStop{
		{Color: red, Position: 0.0},
		{Color: yellow, Position: 0.5},
		{Color: blue, Position: 1.0},
	}

	gradient := GradientMultiStop(stops, 5, GradientOKLCH)

	if len(gradient) != 5 {
		t.Errorf("Gradient length = %d, want 5", len(gradient))
	}

	// First should be red
	r1, g1, b1, _ := gradient[0].RGBA()
	if !floatEqual(r1, 1) || !floatEqual(g1, 0) || !floatEqual(b1, 0) {
		t.Errorf("First color = RGB(%v, %v, %v), want RGB(1, 0, 0)", r1, g1, b1)
	}

	// Last should be blue
	r2, g2, b2, _ := gradient[4].RGBA()
	if !floatEqual(r2, 0) || !floatEqual(g2, 0) || !floatEqual(b2, 1) {
		t.Errorf("Last color = RGB(%v, %v, %v), want RGB(0, 0, 1)", r2, g2, b2)
	}

	// Middle should be yellow-ish
	r3, g3, b3, _ := gradient[2].RGBA()
	// Should have high red and green components
	if r3 < 0.5 || g3 < 0.5 {
		t.Errorf("Middle color = RGB(%v, %v, %v), expected yellow-ish", r3, g3, b3)
	}
}

func TestGradientMultiStopUnsorted(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Stops in wrong order - should be sorted automatically
	stops := []GradientStop{
		{Color: blue, Position: 1.0},
		{Color: red, Position: 0.0},
	}

	gradient := GradientMultiStop(stops, 3, GradientOKLCH)

	if len(gradient) != 3 {
		t.Errorf("Gradient length = %d, want 3", len(gradient))
	}

	// Should still start with red and end with blue
	r1, _, _, _ := gradient[0].RGBA()
	r2, _, _, _ := gradient[2].RGBA()
	if !floatEqual(r1, 1) {
		t.Errorf("First color should be red, got R=%v", r1)
	}
	if !floatEqual(r2, 0) {
		t.Errorf("Last color should be blue, got R=%v", r2)
	}
}

func TestGradientWithEasing(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Test with ease-in
	gradient := GradientWithEasing(red, blue, 5, GradientOKLCH, EaseInQuad)

	if len(gradient) != 5 {
		t.Errorf("Gradient length = %d, want 5", len(gradient))
	}

	// First should be red
	r1, g1, b1, _ := gradient[0].RGBA()
	if !floatEqual(r1, 1) || !floatEqual(g1, 0) || !floatEqual(b1, 0) {
		t.Errorf("First color = RGB(%v, %v, %v), want RGB(1, 0, 0)", r1, g1, b1)
	}

	// Last should be blue
	r2, g2, b2, _ := gradient[4].RGBA()
	if !floatEqual(r2, 0) || !floatEqual(g2, 0) || !floatEqual(b2, 1) {
		t.Errorf("Last color = RGB(%v, %v, %v), want RGB(0, 0, 1)", r2, g2, b2)
	}

	// With ease-in, the first steps should change more slowly
	// (more red in early steps compared to linear)
	r3, _, _, _ := gradient[1].RGBA()
	r4, _, _, _ := gradient[2].RGBA()
	if r3 < r4 {
		t.Errorf("With ease-in, early steps should change slower (r3=%v should be > r4=%v)", r3, r4)
	}
}

func TestGradientMultiStopWithEasing(t *testing.T) {
	red := RGB(1, 0, 0)
	yellow := RGB(1, 1, 0)
	blue := RGB(0, 0, 1)

	stops := []GradientStop{
		{Color: red, Position: 0.0},
		{Color: yellow, Position: 0.5},
		{Color: blue, Position: 1.0},
	}

	gradient := GradientMultiStopWithEasing(stops, 5, GradientOKLCH, EaseInOutQuad)

	if len(gradient) != 5 {
		t.Errorf("Gradient length = %d, want 5", len(gradient))
	}

	// First should be red
	r1, g1, b1, _ := gradient[0].RGBA()
	if !floatEqual(r1, 1) || !floatEqual(g1, 0) || !floatEqual(b1, 0) {
		t.Errorf("First color = RGB(%v, %v, %v), want RGB(1, 0, 0)", r1, g1, b1)
	}

	// Last should be blue
	r2, g2, b2, _ := gradient[4].RGBA()
	if !floatEqual(r2, 0) || !floatEqual(g2, 0) || !floatEqual(b2, 1) {
		t.Errorf("Last color = RGB(%v, %v, %v), want RGB(0, 0, 1)", r2, g2, b2)
	}
}

func TestEasingFunctions(t *testing.T) {
	// Test that easing functions map [0, 1] to [0, 1]
	easingFuncs := []EasingFunction{
		EaseLinear,
		EaseInQuad,
		EaseOutQuad,
		EaseInOutQuad,
		EaseInCubic,
		EaseOutCubic,
		EaseInOutCubic,
		EaseInSine,
		EaseOutSine,
		EaseInOutSine,
	}

	for _, easing := range easingFuncs {
		// Test boundaries
		result0 := easing(0)
		result1 := easing(1)
		if !floatEqual(result0, 0) {
			t.Errorf("Easing function should map 0 to 0, got %v", result0)
		}
		if !floatEqual(result1, 1) {
			t.Errorf("Easing function should map 1 to 1, got %v", result1)
		}

		// Test that result is in [0, 1] for t in [0, 1]
		for testT := 0.0; testT <= 1.0; testT += 0.1 {
			result := easing(testT)
			if result < 0 || result > 1 {
				t.Errorf("Easing function should map %v to [0, 1], got %v", testT, result)
			}
		}
	}
}

func (gs GradientSpace) String() string {
	switch gs {
	case GradientRGB:
		return "RGB"
	case GradientHSL:
		return "HSL"
	case GradientLAB:
		return "LAB"
	case GradientOKLAB:
		return "OKLAB"
	case GradientLCH:
		return "LCH"
	case GradientOKLCH:
		return "OKLCH"
	default:
		return "Unknown"
	}
}

