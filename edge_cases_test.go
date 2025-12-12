package color

import (
	"math"
	"testing"
)

// Test edge cases for color space conversions

func TestRGBAEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		r, g, b, a float64
		expectClamped bool
	}{
		{"Negative values", -0.5, -0.3, -0.2, -0.1, true},
		{"Values over 1", 1.5, 1.3, 1.2, 1.1, true},
		{"Mixed out of range", -0.1, 0.5, 1.2, 0.8, true},
		{"Zero values", 0, 0, 0, 0, false},
		{"All ones", 1, 1, 1, 1, false},
		{"Very small values", 0.001, 0.001, 0.001, 0.001, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewRGBA(tt.r, tt.g, tt.b, tt.a)
			r, g, b, a := c.RGBA()

			// Check clamping
			if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
				t.Errorf("Values not clamped: RGBA(%f, %f, %f, %f)", r, g, b, a)
			}
		})
	}
}

func TestOKLCHEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		l, c, h, a float64
	}{
		{"Negative lightness", -0.1, 0.2, 180, 1.0},
		{"Lightness over 1", 1.5, 0.2, 180, 1.0},
		{"Negative chroma", 0.5, -0.1, 180, 1.0},
		{"Extreme hue values", 0.5, 0.2, 720, 1.0},
		{"Negative hue", 0.5, 0.2, -90, 1.0},
		{"Zero values", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oklch := NewOKLCH(tt.l, tt.c, tt.h, tt.a)

			// Check lightness is clamped [0, 1]
			if oklch.L < 0 || oklch.L > 1 {
				t.Errorf("Lightness not clamped: %f", oklch.L)
			}

			// Check chroma is non-negative
			if oklch.C < 0 {
				t.Errorf("Chroma is negative: %f", oklch.C)
			}

			// Check hue is normalized [0, 360)
			if oklch.H < 0 || oklch.H >= 360 {
				t.Errorf("Hue not normalized: %f", oklch.H)
			}

			// Check alpha is clamped
			if oklch.A_ < 0 || oklch.A_ > 1 {
				t.Errorf("Alpha not clamped: %f", oklch.A_)
			}
		})
	}
}

func TestHSLEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		h, s, l, a float64
	}{
		{"Saturation 0 (gray)", 180, 0, 0.5, 1.0},
		{"Lightness 0 (black)", 180, 1, 0, 1.0},
		{"Lightness 1 (white)", 180, 1, 1, 1.0},
		{"Extreme hue", 720, 0.5, 0.5, 1.0},
		{"Negative hue", -180, 0.5, 0.5, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hsl := NewHSL(tt.h, tt.s, tt.l, tt.a)

			// Convert to RGB and back
			r, g, b, a := hsl.RGBA()

			// Check all values are in range
			if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 {
				t.Errorf("RGB values out of range: (%f, %f, %f)", r, g, b)
			}

			if a < 0 || a > 1 {
				t.Errorf("Alpha out of range: %f", a)
			}
		})
	}
}

func TestHWBEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		h, w, b, a float64
	}{
		{"W+B > 1 (normalization)", 180, 0.7, 0.5, 1.0},
		{"W+B = 1 (gray)", 180, 0.5, 0.5, 1.0},
		{"W=1, B=0 (white)", 180, 1, 0, 1.0},
		{"W=0, B=1 (black)", 180, 0, 1, 1.0},
		{"All zero", 0, 0, 0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hwb := NewHWB(tt.h, tt.w, tt.b, tt.a)

			// Check normalization when W+B > 1
			if hwb.W + hwb.B > 1.001 { // Small tolerance for float precision
				t.Errorf("W+B not normalized: W=%f, B=%f, sum=%f", hwb.W, hwb.B, hwb.W+hwb.B)
			}

			// Convert to RGB
			r, g, b, a := hwb.RGBA()
			if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
				t.Errorf("RGB values out of range: (%f, %f, %f, %f)", r, g, b, a)
			}
		})
	}
}

func TestColorConversionRoundTrips(t *testing.T) {
	testColors := []Color{
		RGB(0.5, 0.5, 0.5),  // Gray
		RGB(1, 0, 0),        // Red
		RGB(0, 1, 0),        // Green
		RGB(0, 0, 1),        // Blue
		RGB(0, 0, 0),        // Black
		RGB(1, 1, 1),        // White
		RGB(0.3, 0.7, 0.2),  // Random color
	}

	tolerance := 0.01

	for _, original := range testColors {
		t.Run(RGBToHex(original), func(t *testing.T) {
			// RGB -> OKLCH -> RGB
			oklch := ToOKLCH(original)
			r1, g1, b1, _ := oklch.RGBA()
			r0, g0, b0, _ := original.RGBA()

			if math.Abs(r1-r0) > tolerance || math.Abs(g1-g0) > tolerance || math.Abs(b1-b0) > tolerance {
				t.Errorf("OKLCH round trip failed: (%f,%f,%f) -> (%f,%f,%f)", r0, g0, b0, r1, g1, b1)
			}

			// RGB -> HSL -> RGB
			hsl := ToHSL(original)
			r2, g2, b2, _ := hsl.RGBA()

			if math.Abs(r2-r0) > tolerance || math.Abs(g2-g0) > tolerance || math.Abs(b2-b0) > tolerance {
				t.Errorf("HSL round trip failed: (%f,%f,%f) -> (%f,%f,%f)", r0, g0, b0, r2, g2, b2)
			}
		})
	}
}

func TestXYZConversionsEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		color Color
	}{
		{"Black", RGB(0, 0, 0)},
		{"White", RGB(1, 1, 1)},
		{"Very dark", RGB(0.001, 0.001, 0.001)},
		{"Very bright", RGB(0.999, 0.999, 0.999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xyz := ToXYZ(tt.color)

			// XYZ values should be non-negative
			if xyz.X < 0 || xyz.Y < 0 || xyz.Z < 0 {
				t.Errorf("XYZ has negative values: X=%f, Y=%f, Z=%f", xyz.X, xyz.Y, xyz.Z)
			}

			// Convert back to RGB
			r, g, b, _ := xyz.RGBA()
			r0, g0, b0, _ := tt.color.RGBA()

			tolerance := 0.01
			if math.Abs(r-r0) > tolerance || math.Abs(g-g0) > tolerance || math.Abs(b-b0) > tolerance {
				t.Errorf("XYZ round trip failed: (%f,%f,%f) -> (%f,%f,%f)", r0, g0, b0, r, g, b)
			}
		})
	}
}

func TestAlphaPreservation(t *testing.T) {
	alphas := []float64{0.0, 0.25, 0.5, 0.75, 1.0}
	base := RGB(0.5, 0.3, 0.7)

	for _, alpha := range alphas {
		t.Run(RGBToHex(base.WithAlpha(alpha)), func(t *testing.T) {
			c := base.WithAlpha(alpha)

			// Test that alpha is preserved through conversions
			oklch := ToOKLCH(c)
			if math.Abs(oklch.Alpha()-alpha) > 0.01 {
				t.Errorf("Alpha not preserved in OKLCH: got %f, want %f", oklch.Alpha(), alpha)
			}

			hsl := ToHSL(c)
			if math.Abs(hsl.Alpha()-alpha) > 0.01 {
				t.Errorf("Alpha not preserved in HSL: got %f, want %f", hsl.Alpha(), alpha)
			}

			xyz := ToXYZ(c)
			if math.Abs(xyz.Alpha()-alpha) > 0.01 {
				t.Errorf("Alpha not preserved in XYZ: got %f, want %f", xyz.Alpha(), alpha)
			}
		})
	}
}

func TestNaNAndInfHandling(t *testing.T) {
	// Test that NaN and Inf don't cause panics
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panic occurred with NaN/Inf values: %v", r)
		}
	}()

	// These should not panic, though results may be clamped
	_ = NewRGBA(math.NaN(), 0.5, 0.5, 1.0)
	_ = NewRGBA(math.Inf(1), 0.5, 0.5, 1.0)
	_ = NewRGBA(math.Inf(-1), 0.5, 0.5, 1.0)

	_ = NewOKLCH(math.NaN(), 0.2, 180, 1.0)
	_ = NewHSL(math.NaN(), 0.5, 0.5, 1.0)
}
