package color

import (
	"testing"
)

// Test round-trip conversions RGB <-> HSL
func TestRGBHSLRoundTrip(t *testing.T) {
	colors := []struct {
		name string
		rgb  *RGBA
	}{
		{"Red", RGB(1, 0, 0)},
		{"Green", RGB(0, 1, 0)},
		{"Blue", RGB(0, 0, 1)},
		{"White", RGB(1, 1, 1)},
		{"Black", RGB(0, 0, 0)},
		{"Cyan", RGB(0, 1, 1)},
		{"Magenta", RGB(1, 0, 1)},
		{"Yellow", RGB(1, 1, 0)},
		{"Orange", RGB(1, 0.5, 0)},
		{"Purple", RGB(0.5, 0, 0.5)},
		{"Gray", RGB(0.5, 0.5, 0.5)},
	}

	for _, tt := range colors {
		t.Run(tt.name, func(t *testing.T) {
			// RGB -> HSL -> RGB
			hsl := ToHSL(tt.rgb)
			r2, g2, b2, a2 := hsl.RGBA()

			r1, g1, b1, a1 := tt.rgb.RGBA()

			// Allow some tolerance for round-trip errors
			if !rgbaEqual(r1, g1, b1, a1, r2, g2, b2, a2) {
				t.Errorf("RGB->HSL->RGB round-trip failed:\n"+
					"  Original: RGB(%v, %v, %v, %v)\n"+
					"  HSL: H=%v, S=%v, L=%v\n"+
					"  Result: RGB(%v, %v, %v, %v)",
					r1, g1, b1, a1, hsl.H, hsl.S, hsl.L, r2, g2, b2, a2)
			}
		})
	}
}

// Test known HSL values
func TestHSLKnownValues(t *testing.T) {
	tests := []struct {
		name     string
		h, s, l  float64
		expected [3]float64 // RGB (no alpha)
	}{
		// Pure colors at full saturation and 50% lightness
		{"Red HSL", 0, 1, 0.5, [3]float64{1, 0, 0}},
		{"Green HSL", 120, 1, 0.5, [3]float64{0, 1, 0}},
		{"Blue HSL", 240, 1, 0.5, [3]float64{0, 0, 1}},
		// White (any hue, 0 saturation, 100% lightness)
		{"White HSL", 0, 0, 1, [3]float64{1, 1, 1}},
		// Black (any hue, any saturation, 0 lightness)
		{"Black HSL", 0, 1, 0, [3]float64{0, 0, 0}},
		// Gray (0 saturation)
		{"Gray HSL", 0, 0, 0.5, [3]float64{0.5, 0.5, 0.5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsl := NewHSL(tt.h, tt.s, tt.l, 1.0)
			r, g, b, _ := hsl.RGBA()
			if !floatEqual(r, tt.expected[0]) || !floatEqual(g, tt.expected[1]) || !floatEqual(b, tt.expected[2]) {
				t.Errorf("HSL(%v, %v, %v) = RGB(%v, %v, %v), want RGB(%v, %v, %v)",
					tt.h, tt.s, tt.l, r, g, b, tt.expected[0], tt.expected[1], tt.expected[2])
			}
		})
	}
}

// Test RGB to HSV conversion
func TestRGBHSVRoundTrip(t *testing.T) {
	colors := []struct {
		name string
		rgb  *RGBA
	}{
		{"Red", RGB(1, 0, 0)},
		{"Green", RGB(0, 1, 0)},
		{"Blue", RGB(0, 0, 1)},
		{"White", RGB(1, 1, 1)},
		{"Black", RGB(0, 0, 0)},
		{"Cyan", RGB(0, 1, 1)},
		{"Magenta", RGB(1, 0, 1)},
		{"Yellow", RGB(1, 1, 0)},
		{"Orange", RGB(1, 0.5, 0)},
	}

	for _, tt := range colors {
		t.Run(tt.name, func(t *testing.T) {
			// RGB -> HSV -> RGB
			hsv := ToHSV(tt.rgb)
			r2, g2, b2, a2 := hsv.RGBA()

			r1, g1, b1, a1 := tt.rgb.RGBA()

			if !rgbaEqual(r1, g1, b1, a1, r2, g2, b2, a2) {
				t.Errorf("RGB->HSV->RGB round-trip failed:\n"+
					"  Original: RGB(%v, %v, %v, %v)\n"+
					"  HSV: H=%v, S=%v, V=%v\n"+
					"  Result: RGB(%v, %v, %v, %v)",
					r1, g1, b1, a1, hsv.H, hsv.S, hsv.V, r2, g2, b2, a2)
			}
		})
	}
}

// Test known HSV values
func TestHSVKnownValues(t *testing.T) {
	tests := []struct {
		name     string
		h, s, v  float64
		expected [3]float64 // RGB
	}{
		{"Red HSV", 0, 1, 1, [3]float64{1, 0, 0}},
		{"Green HSV", 120, 1, 1, [3]float64{0, 1, 0}},
		{"Blue HSV", 240, 1, 1, [3]float64{0, 0, 1}},
		{"White HSV", 0, 0, 1, [3]float64{1, 1, 1}},
		{"Black HSV", 0, 0, 0, [3]float64{0, 0, 0}},
		{"Gray HSV", 0, 0, 0.5, [3]float64{0.5, 0.5, 0.5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsv := NewHSV(tt.h, tt.s, tt.v, 1.0)
			r, g, b, _ := hsv.RGBA()
			if !floatEqual(r, tt.expected[0]) || !floatEqual(g, tt.expected[1]) || !floatEqual(b, tt.expected[2]) {
				t.Errorf("HSV(%v, %v, %v) = RGB(%v, %v, %v), want RGB(%v, %v, %v)",
					tt.h, tt.s, tt.v, r, g, b, tt.expected[0], tt.expected[1], tt.expected[2])
			}
		})
	}
}

