package color

import (
	"math"
	"testing"
)

func TestChromaticAdaptationRoundTrip(t *testing.T) {
	// Test that adapting D65 -> D50 -> D65 gives back the original
	tests := []struct {
		name       string
		x, y, z    float64
		tolerance  float64
	}{
		{"White D65", 0.95047, 1.0, 1.08883, 0.0001},
		{"Red", 0.4124, 0.2126, 0.0193, 0.0001},
		{"Green", 0.3576, 0.7152, 0.1192, 0.0001},
		{"Blue", 0.1805, 0.0722, 0.9505, 0.0001},
		{"Gray", 0.2034, 0.2140, 0.2330, 0.0001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// D65 -> D50
			xD50, yD50, zD50 := AdaptD65ToD50(tt.x, tt.y, tt.z)

			// D50 -> D65 (round trip)
			xD65, yD65, zD65 := AdaptD50ToD65(xD50, yD50, zD50)

			// Check they match original
			if math.Abs(xD65-tt.x) > tt.tolerance {
				t.Errorf("X component: got %f, want %f (diff %f)", xD65, tt.x, math.Abs(xD65-tt.x))
			}
			if math.Abs(yD65-tt.y) > tt.tolerance {
				t.Errorf("Y component: got %f, want %f (diff %f)", yD65, tt.y, math.Abs(yD65-tt.y))
			}
			if math.Abs(zD65-tt.z) > tt.tolerance {
				t.Errorf("Z component: got %f, want %f (diff %f)", zD65, tt.z, math.Abs(zD65-tt.z))
			}
		})
	}
}

func TestChromaticAdaptationD65WhitePoint(t *testing.T) {
	// D65 white point should remain approximately white when adapted
	xD65, yD65, zD65 := 0.95047, 1.0, 1.08883

	// Adapt to D50
	xD50, yD50, zD50 := AdaptD65ToD50(xD65, yD65, zD65)

	// Y should remain 1.0 (luminance preserved)
	if math.Abs(yD50-1.0) > 0.0001 {
		t.Errorf("Luminance not preserved: Y = %f, want 1.0", yD50)
	}

	// The adapted values should be close to D50 white point
	// D50 white point is approximately (0.9642, 1.0, 0.8251)
	if math.Abs(xD50-0.9642) > 0.01 {
		t.Errorf("X component: got %f, want ~0.9642", xD50)
	}
	if math.Abs(zD50-0.8251) > 0.01 {
		t.Errorf("Z component: got %f, want ~0.8251", zD50)
	}
}

func TestChromaticAdaptationPreservesBlack(t *testing.T) {
	// Black should remain black
	x, y, z := 0.0, 0.0, 0.0

	xD50, yD50, zD50 := AdaptD65ToD50(x, y, z)
	if xD50 != 0 || yD50 != 0 || zD50 != 0 {
		t.Errorf("Black not preserved: got (%f, %f, %f)", xD50, yD50, zD50)
	}

	xD65, yD65, zD65 := AdaptD50ToD65(x, y, z)
	if xD65 != 0 || yD65 != 0 || zD65 != 0 {
		t.Errorf("Black not preserved: got (%f, %f, %f)", xD65, yD65, zD65)
	}
}

func TestChromaticAdaptationScales(t *testing.T) {
	// Chromatic adaptation should scale linearly
	x, y, z := 0.4124, 0.2126, 0.0193 // Red

	// Adapt at full intensity
	x1, y1, z1 := AdaptD65ToD50(x, y, z)

	// Adapt at half intensity
	x2, y2, z2 := AdaptD65ToD50(x*0.5, y*0.5, z*0.5)

	// The adapted values should also be half
	if math.Abs(x2-x1*0.5) > 0.0001 {
		t.Errorf("X doesn't scale linearly: got %f, want %f", x2, x1*0.5)
	}
	if math.Abs(y2-y1*0.5) > 0.0001 {
		t.Errorf("Y doesn't scale linearly: got %f, want %f", y2, y1*0.5)
	}
	if math.Abs(z2-z1*0.5) > 0.0001 {
		t.Errorf("Z doesn't scale linearly: got %f, want %f", z2, z1*0.5)
	}
}

func TestBradfordMatrixInverse(t *testing.T) {
	// Test that the matrices are properly inverse of each other
	// by checking M_D50toD65 * M_D65toD50 ≈ I

	x, y, z := 0.5, 0.6, 0.7 // Arbitrary color

	// Apply both transformations
	x1, y1, z1 := AdaptD65ToD50(x, y, z)
	x2, y2, z2 := AdaptD50ToD65(x1, y1, z1)

	// Should get back original
	if math.Abs(x2-x) > 0.0001 || math.Abs(y2-y) > 0.0001 || math.Abs(z2-z) > 0.0001 {
		t.Errorf("Inverse transformation failed: (%f,%f,%f) -> (%f,%f,%f) -> (%f,%f,%f)",
			x, y, z, x1, y1, z1, x2, y2, z2)
	}
}

func BenchmarkAdaptD65ToD50(b *testing.B) {
	x, y, z := 0.4124, 0.2126, 0.0193
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AdaptD65ToD50(x, y, z)
	}
}

func BenchmarkAdaptD50ToD65(b *testing.B) {
	x, y, z := 0.4124, 0.2126, 0.0193
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AdaptD50ToD65(x, y, z)
	}
}
