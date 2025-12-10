package color

import (
	"math"
	"testing"
)

const epsilon = 1e-5 // Tolerance for floating point comparisons

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

func rgbaEqual(r1, g1, b1, a1, r2, g2, b2, a2 float64) bool {
	return floatEqual(r1, r2) && floatEqual(g1, g2) && floatEqual(b1, b2) && floatEqual(a1, a2)
}

func TestRGB(t *testing.T) {
	tests := []struct {
		name     string
		r, g, b   float64
		expected [4]float64
	}{
		{"Black", 0, 0, 0, [4]float64{0, 0, 0, 1}},
		{"White", 1, 1, 1, [4]float64{1, 1, 1, 1}},
		{"Red", 1, 0, 0, [4]float64{1, 0, 0, 1}},
		{"Green", 0, 1, 0, [4]float64{0, 1, 0, 1}},
		{"Blue", 0, 0, 1, [4]float64{0, 0, 1, 1}},
		{"Cyan", 0, 1, 1, [4]float64{0, 1, 1, 1}},
		{"Magenta", 1, 0, 1, [4]float64{1, 0, 1, 1}},
		{"Yellow", 1, 1, 0, [4]float64{1, 1, 0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := RGB(tt.r, tt.g, tt.b)
			r, g, b, a := c.RGBA()
			if !rgbaEqual(r, g, b, a, tt.expected[0], tt.expected[1], tt.expected[2], tt.expected[3]) {
				t.Errorf("RGB(%v, %v, %v) = (%v, %v, %v, %v), want (%v, %v, %v, %v)",
					tt.r, tt.g, tt.b, r, g, b, a, tt.expected[0], tt.expected[1], tt.expected[2], tt.expected[3])
			}
		})
	}
}

func TestNewRGBA(t *testing.T) {
	c := NewRGBA(0.5, 0.3, 0.7, 0.8)
	r, g, b, a := c.RGBA()
	if !rgbaEqual(r, g, b, a, 0.5, 0.3, 0.7, 0.8) {
		t.Errorf("NewRGBA(0.5, 0.3, 0.7, 0.8) = (%v, %v, %v, %v), want (0.5, 0.3, 0.7, 0.8)", r, g, b, a)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"Normal", 0.5, 0.5},
		{"Below zero", -0.5, 0},
		{"Above one", 1.5, 1},
		{"Zero", 0, 0},
		{"One", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clamp01(tt.value)
			if !floatEqual(result, tt.expected) {
				t.Errorf("clamp01(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestAlpha(t *testing.T) {
	c := NewRGBA(1, 0, 0, 0.5)
	if !floatEqual(c.Alpha(), 0.5) {
		t.Errorf("Alpha() = %v, want 0.5", c.Alpha())
	}
}

func TestWithAlpha(t *testing.T) {
	c := RGB(1, 0, 0)
	newC := c.WithAlpha(0.7)
	if !floatEqual(newC.Alpha(), 0.7) {
		t.Errorf("WithAlpha(0.7).Alpha() = %v, want 0.7", newC.Alpha())
	}
	// Original should be unchanged
	if !floatEqual(c.Alpha(), 1.0) {
		t.Errorf("Original Alpha() = %v, want 1.0", c.Alpha())
	}
}

