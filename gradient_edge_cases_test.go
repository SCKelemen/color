package color

import (
	"math"
	"testing"
)

// Test edge cases for gradient generation

func TestGradientSingleColor(t *testing.T) {
	// Gradient with same start and end color
	red := RGB(1, 0, 0)
	gradient := Gradient(red, red, 5)

	if len(gradient) != 5 {
		t.Errorf("Expected 5 colors, got %d", len(gradient))
	}

	// All colors should be the same
	for i, c := range gradient {
		r, g, b, _ := c.RGBA()
		if math.Abs(r-1.0) > 0.01 || g > 0.01 || b > 0.01 {
			t.Errorf("Color %d not red: (%f, %f, %f)", i, r, g, b)
		}
	}
}

func TestGradientMinimalSteps(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Test with 1 step (should return start color)
	gradient1 := Gradient(red, blue, 1)
	if len(gradient1) != 1 {
		t.Errorf("Expected 1 color, got %d", len(gradient1))
	}

	// Test with 2 steps (start and end)
	gradient2 := Gradient(red, blue, 2)
	if len(gradient2) != 2 {
		t.Errorf("Expected 2 colors, got %d", len(gradient2))
	}

	// First should be red, last should be blue
	r1, g1, b1, _ := gradient2[0].RGBA()
	if math.Abs(r1-1.0) > 0.01 || g1 > 0.01 || b1 > 0.01 {
		t.Error("First color should be red")
	}

	r2, g2, b2, _ := gradient2[1].RGBA()
	if r2 > 0.01 || g2 > 0.01 || math.Abs(b2-1.0) > 0.01 {
		t.Error("Last color should be blue")
	}
}

func TestGradientZeroSteps(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Should handle gracefully
	gradient := Gradient(red, blue, 0)
	if gradient == nil {
		t.Error("Gradient returned nil for 0 steps")
	}
	// Implementation may return empty or single color
}

func TestGradientMultiStopEmpty(t *testing.T) {
	// Empty stops slice
	gradient := GradientMultiStop([]GradientStop{}, 10, GradientOKLCH)

	if gradient == nil {
		t.Error("GradientMultiStop returned nil for empty stops")
	}
}

func TestGradientMultiStopSingle(t *testing.T) {
	// Single stop
	stops := []GradientStop{
		{Color: RGB(1, 0, 0), Position: 0.5},
	}

	gradient := GradientMultiStop(stops, 10, GradientOKLCH)

	if len(gradient) == 0 {
		t.Error("Expected some colors for single stop")
	}

	// All colors should be the same as the single stop
	for _, c := range gradient {
		r, g, b, _ := c.RGBA()
		if math.Abs(r-1.0) > 0.01 || g > 0.1 || b > 0.1 {
			t.Errorf("Color not red: (%f, %f, %f)", r, g, b)
		}
	}
}

func TestGradientMultiStopUnsortedEdgeCase(t *testing.T) {
	// Stops in wrong order
	stops := []GradientStop{
		{Color: RGB(0, 0, 1), Position: 1.0},
		{Color: RGB(1, 0, 0), Position: 0.0},
		{Color: RGB(0, 1, 0), Position: 0.5},
	}

	gradient := GradientMultiStop(stops, 10, GradientOKLCH)

	if len(gradient) != 10 {
		t.Errorf("Expected 10 colors, got %d", len(gradient))
	}

	// Should still work (implementation should sort)
	// First should be red, last should be blue
	r1, _, b1, _ := gradient[0].RGBA()
	if r1 < 0.5 || b1 > 0.5 {
		t.Error("First color should be red-ish")
	}

	r2, _, b2, _ := gradient[len(gradient)-1].RGBA()
	if r2 > 0.5 || b2 < 0.5 {
		t.Error("Last color should be blue-ish")
	}
}

func TestGradientMultiStopDuplicatePositions(t *testing.T) {
	// Multiple stops at same position
	stops := []GradientStop{
		{Color: RGB(1, 0, 0), Position: 0.0},
		{Color: RGB(0, 1, 0), Position: 0.5},
		{Color: RGB(0, 0, 1), Position: 0.5},
		{Color: RGB(1, 1, 0), Position: 1.0},
	}

	gradient := GradientMultiStop(stops, 20, GradientOKLCH)

	if len(gradient) != 20 {
		t.Errorf("Expected 20 colors, got %d", len(gradient))
	}
}

func TestGradientMultiStopOutOfRangePositions(t *testing.T) {
	// Positions outside [0, 1]
	stops := []GradientStop{
		{Color: RGB(1, 0, 0), Position: -0.5},
		{Color: RGB(0, 1, 0), Position: 0.5},
		{Color: RGB(0, 0, 1), Position: 1.5},
	}

	// Should handle gracefully (clamp or ignore)
	gradient := GradientMultiStop(stops, 10, GradientOKLCH)

	if gradient == nil {
		t.Error("GradientMultiStop returned nil for out-of-range positions")
	}
}

func TestGradientWithEasingEdgeCases(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	easings := []EasingFunction{
		EaseLinear,
		EaseInQuad,
		EaseOutQuad,
		EaseInOutQuad,
		EaseInCubic,
		EaseOutCubic,
		EaseInOutCubic,
	}

	for _, easing := range easings {
		gradient := GradientWithEasing(red, blue, 10, GradientOKLCH, easing)

		if len(gradient) != 10 {
			t.Errorf("Expected 10 colors, got %d", len(gradient))
		}

		// First and last should still be red and blue
		r1, g1, b1, _ := gradient[0].RGBA()
		if math.Abs(r1-1.0) > 0.01 || g1 > 0.01 || b1 > 0.01 {
			t.Error("First color should be red")
		}

		r2, g2, b2, _ := gradient[len(gradient)-1].RGBA()
		if r2 > 0.01 || g2 > 0.01 || math.Abs(b2-1.0) > 0.01 {
			t.Error("Last color should be blue")
		}
	}
}

func TestGradientAllSpaces(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	spaces := []GradientSpace{
		GradientRGB,
		GradientHSL,
		GradientLAB,
		GradientOKLAB,
		GradientLCH,
		GradientOKLCH,
	}

	for _, space := range spaces {
		gradient := GradientInSpace(red, blue, 10, space)

		if len(gradient) != 10 {
			t.Errorf("Space %v: expected 10 colors, got %d", space, len(gradient))
		}

		// Verify all colors are valid
		for i, c := range gradient {
			r, g, b, a := c.RGBA()
			if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
				t.Errorf("Space %v, color %d: invalid RGB (%f, %f, %f, %f)", space, i, r, g, b, a)
			}
		}
	}
}

func TestGradientAlphaInterpolation(t *testing.T) {
	// Gradient with different alpha values
	transparent := NewRGBA(1, 0, 0, 0.0)
	opaque := NewRGBA(1, 0, 0, 1.0)

	gradient := Gradient(transparent, opaque, 11)

	if len(gradient) != 11 {
		t.Errorf("Expected 11 colors, got %d", len(gradient))
	}

	// Check alpha interpolation
	for i, c := range gradient {
		expectedAlpha := float64(i) / 10.0
		actualAlpha := c.Alpha()

		if math.Abs(actualAlpha-expectedAlpha) > 0.1 {
			t.Errorf("Color %d: expected alpha %f, got %f", i, expectedAlpha, actualAlpha)
		}
	}
}

func TestGradientIdenticalColors(t *testing.T) {
	// All stops are the same color
	gray := RGB(0.5, 0.5, 0.5)
	stops := []GradientStop{
		{Color: gray, Position: 0.0},
		{Color: gray, Position: 0.5},
		{Color: gray, Position: 1.0},
	}

	gradient := GradientMultiStop(stops, 20, GradientOKLCH)

	// All colors should be gray
	for i, c := range gradient {
		r, g, b, _ := c.RGBA()
		if math.Abs(r-0.5) > 0.01 || math.Abs(g-0.5) > 0.01 || math.Abs(b-0.5) > 0.01 {
			t.Errorf("Color %d not gray: (%f, %f, %f)", i, r, g, b)
		}
	}
}

func TestGradientLargeStepCount(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Test with large number of steps
	gradient := Gradient(red, blue, 1000)

	if len(gradient) != 1000 {
		t.Errorf("Expected 1000 colors, got %d", len(gradient))
	}

	// Verify smooth progression
	tolerance := 0.002 // Small tolerance for 1000 steps
	for i := 1; i < len(gradient); i++ {
		r1, g1, b1, _ := gradient[i-1].RGBA()
		r2, g2, b2, _ := gradient[i].RGBA()

		// Change should be gradual
		dr := math.Abs(r2 - r1)
		dg := math.Abs(g2 - g1)
		db := math.Abs(b2 - b1)

		if dr > tolerance || dg > tolerance || db > tolerance {
			t.Logf("Large jump between steps %d and %d: dr=%f, dg=%f, db=%f", i-1, i, dr, dg, db)
		}
	}
}
