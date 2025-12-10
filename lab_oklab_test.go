package color

import (
	"math"
	"testing"
)

// Test round-trip conversions RGB <-> LAB
func TestRGBLABRoundTrip(t *testing.T) {
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
		{"Gray", RGB(0.5, 0.5, 0.5)},
	}

	// LAB round-trips can have larger errors due to gamut mapping
	labEpsilon := 1e-3

	for _, tt := range colors {
		t.Run(tt.name, func(t *testing.T) {
			// RGB -> LAB -> RGB
			lab := ToLAB(tt.rgb)
			r2, g2, b2, a2 := lab.RGBA()

			r1, g1, b1, a1 := tt.rgb.RGBA()

			// Check alpha is preserved
			if !floatEqual(a1, a2) {
				t.Errorf("Alpha not preserved: %v != %v", a1, a2)
			}

			// Check RGB values (allowing for gamut mapping errors)
			dr := math.Abs(r1 - r2)
			dg := math.Abs(g1 - g2)
			db := math.Abs(b1 - b2)

			if dr > labEpsilon || dg > labEpsilon || db > labEpsilon {
				t.Errorf("RGB->LAB->RGB round-trip failed:\n"+
					"  Original: RGB(%v, %v, %v, %v)\n"+
					"  LAB: L=%v, A=%v, B=%v\n"+
					"  Result: RGB(%v, %v, %v, %v)\n"+
					"  Diff: RGB(%v, %v, %v)",
					r1, g1, b1, a1, lab.L, lab.A, lab.B, r2, g2, b2, a2, dr, dg, db)
			}
		})
	}
}

// Test round-trip conversions RGB <-> OKLAB
func TestRGBOKLABRoundTrip(t *testing.T) {
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
		{"Gray", RGB(0.5, 0.5, 0.5)},
		{"Orange", RGB(1, 0.5, 0)},
	}

	// OKLAB should have better round-trip accuracy
	oklabEpsilon := 1e-4

	for _, tt := range colors {
		t.Run(tt.name, func(t *testing.T) {
			// RGB -> OKLAB -> RGB
			oklab := ToOKLAB(tt.rgb)
			r2, g2, b2, a2 := oklab.RGBA()

			r1, g1, b1, a1 := tt.rgb.RGBA()

			// Check alpha is preserved
			if !floatEqual(a1, a2) {
				t.Errorf("Alpha not preserved: %v != %v", a1, a2)
			}

			// Check RGB values
			dr := math.Abs(r1 - r2)
			dg := math.Abs(g1 - g2)
			db := math.Abs(b1 - b2)

			if dr > oklabEpsilon || dg > oklabEpsilon || db > oklabEpsilon {
				t.Errorf("RGB->OKLAB->RGB round-trip failed:\n"+
					"  Original: RGB(%v, %v, %v, %v)\n"+
					"  OKLAB: L=%v, A=%v, B=%v\n"+
					"  Result: RGB(%v, %v, %v, %v)\n"+
					"  Diff: RGB(%v, %v, %v)",
					r1, g1, b1, a1, oklab.L, oklab.A, oklab.B, r2, g2, b2, a2, dr, dg, db)
			}
		})
	}
}

// Test known OKLAB values (from CSS Color 4 spec examples)
func TestOKLABKnownValues(t *testing.T) {
	// These are approximate values - OKLAB is designed for perceptual uniformity
	tests := []struct {
		name     string
		oklab    *OKLAB
		expected [3]float64 // RGB (approximate)
		tolerance float64
	}{
		// White should be L=1, A=0, B=0
		{"White OKLAB", NewOKLAB(1, 0, 0, 1), [3]float64{1, 1, 1}, 0.01},
		// Black should be L=0, A=0, B=0
		{"Black OKLAB", NewOKLAB(0, 0, 0, 1), [3]float64{0, 0, 0}, 0.01},
		// Gray should have A=0, B=0 (but OKLAB L=0.5 doesn't map exactly to RGB 0.5)
		{"Gray OKLAB", NewOKLAB(0.5, 0, 0, 1), [3]float64{0.388, 0.388, 0.388}, 0.05},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b, _ := tt.oklab.RGBA()
			dr := math.Abs(r - tt.expected[0])
			dg := math.Abs(g - tt.expected[1])
			db := math.Abs(b - tt.expected[2])

			if dr > tt.tolerance || dg > tt.tolerance || db > tt.tolerance {
				t.Errorf("OKLAB(%v, %v, %v) = RGB(%v, %v, %v), want RGB(%v, %v, %v) ±%v",
					tt.oklab.L, tt.oklab.A, tt.oklab.B, r, g, b,
					tt.expected[0], tt.expected[1], tt.expected[2], tt.tolerance)
			}
		})
	}
}

// Test LAB <-> LCH round-trip
func TestLABLCHRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		lab  *LAB
	}{
		{"Red LAB", ToLAB(RGB(1, 0, 0))},
		{"Green LAB", ToLAB(RGB(0, 1, 0))},
		{"Blue LAB", ToLAB(RGB(0, 0, 1))},
		{"White LAB", ToLAB(RGB(1, 1, 1))},
		{"Black LAB", ToLAB(RGB(0, 0, 0))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// LAB -> LCH -> LAB
			lch := tt.lab.toLCH()
			lab2 := lch.toLAB()

			// Check values are close
			dl := math.Abs(tt.lab.L - lab2.L)
			da := math.Abs(tt.lab.A - lab2.A)
			db := math.Abs(tt.lab.B - lab2.B)

			if dl > 1e-5 || da > 1e-5 || db > 1e-5 {
				t.Errorf("LAB->LCH->LAB round-trip failed:\n"+
					"  Original: LAB(%v, %v, %v)\n"+
					"  LCH: L=%v, C=%v, H=%v\n"+
					"  Result: LAB(%v, %v, %v)\n"+
					"  Diff: LAB(%v, %v, %v)",
					tt.lab.L, tt.lab.A, tt.lab.B, lch.L, lch.C, lch.H,
					lab2.L, lab2.A, lab2.B, dl, da, db)
			}
		})
	}
}

// Test OKLAB <-> OKLCH round-trip
func TestOKLABOKLCHRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		oklab *OKLAB
	}{
		{"Red OKLAB", ToOKLAB(RGB(1, 0, 0))},
		{"Green OKLAB", ToOKLAB(RGB(0, 1, 0))},
		{"Blue OKLAB", ToOKLAB(RGB(0, 0, 1))},
		{"White OKLAB", ToOKLAB(RGB(1, 1, 1))},
		{"Black OKLAB", ToOKLAB(RGB(0, 0, 0))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// OKLAB -> OKLCH -> OKLAB
			oklch := tt.oklab.toOKLCH()
			oklab2 := oklch.toOKLAB()

			// Check values are close
			dl := math.Abs(tt.oklab.L - oklab2.L)
			da := math.Abs(tt.oklab.A - oklab2.A)
			db := math.Abs(tt.oklab.B - oklab2.B)

			if dl > 1e-5 || da > 1e-5 || db > 1e-5 {
				t.Errorf("OKLAB->OKLCH->OKLAB round-trip failed:\n"+
					"  Original: OKLAB(%v, %v, %v)\n"+
					"  OKLCH: L=%v, C=%v, H=%v\n"+
					"  Result: OKLAB(%v, %v, %v)\n"+
					"  Diff: OKLAB(%v, %v, %v)",
					tt.oklab.L, tt.oklab.A, tt.oklab.B, oklch.L, oklch.C, oklch.H,
					oklab2.L, oklab2.A, oklab2.B, dl, da, db)
			}
		})
	}
}

