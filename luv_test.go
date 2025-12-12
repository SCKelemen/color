package color

import (
	"math"
	"testing"
)

func TestNewLUV(t *testing.T) {
	tests := []struct {
		name string
		l, u, v, a float64
	}{
		{"Black", 0, 0, 0, 1},
		{"White", 100, 0, 0, 1},
		{"Mid-tone", 50, 25, -30, 1},
		{"With alpha", 50, 10, 20, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			luv := NewLUV(tt.l, tt.u, tt.v, tt.a)
			if luv.L != tt.l || luv.U != tt.u || luv.V != tt.v || luv.A_ != tt.a {
				t.Errorf("NewLUV(%f, %f, %f, %f) = %+v", tt.l, tt.u, tt.v, tt.a, luv)
			}
		})
	}
}

func TestLUVRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		color Color
	}{
		{"Red", RGB(1, 0, 0)},
		{"Green", RGB(0, 1, 0)},
		{"Blue", RGB(0, 0, 1)},
		{"Gray", RGB(0.5, 0.5, 0.5)},
		{"Cyan", RGB(0, 1, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to LUV and back
			luv := ToLUV(tt.color)
			r2, g2, b2, a2 := luv.RGBA()

			r1, g1, b1, a1 := tt.color.RGBA()

			if math.Abs(r1-r2) > 0.01 || math.Abs(g1-g2) > 0.01 || math.Abs(b1-b2) > 0.01 {
				t.Errorf("Round trip failed: (%f,%f,%f,%f) -> LUV -> (%f,%f,%f,%f)",
					r1, g1, b1, a1, r2, g2, b2, a2)
			}
		})
	}
}

func TestNewLCHuv(t *testing.T) {
	tests := []struct {
		name string
		l, c, h, a float64
	}{
		{"Black", 0, 0, 0, 1},
		{"White", 100, 0, 0, 1},
		{"Red-ish", 50, 60, 20, 1},
		{"With alpha", 70, 40, 180, 0.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lch := NewLCHuv(tt.l, tt.c, tt.h, tt.a)
			if math.Abs(lch.L-tt.l) > 0.01 || math.Abs(lch.C-tt.c) > 0.01 {
				t.Errorf("NewLCHuv(%f, %f, %f, %f) L,C mismatch = %+v", tt.l, tt.c, tt.h, tt.a, lch)
			}
			// Hue should be normalized to [0, 360)
			if lch.H < 0 || lch.H >= 360 {
				t.Errorf("Hue not normalized: %f", lch.H)
			}
		})
	}
}

func TestLCHuvRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		color Color
	}{
		{"Red", RGB(1, 0, 0)},
		{"Green", RGB(0, 1, 0)},
		{"Blue", RGB(0, 0, 1)},
		{"Yellow", RGB(1, 1, 0)},
		{"Magenta", RGB(1, 0, 1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to LCHuv and back
			lch := ToLCHuv(tt.color)
			r2, g2, b2, a2 := lch.RGBA()

			r1, g1, b1, a1 := tt.color.RGBA()

			if math.Abs(r1-r2) > 0.01 || math.Abs(g1-g2) > 0.01 || math.Abs(b1-b2) > 0.01 {
				t.Errorf("Round trip failed: (%f,%f,%f,%f) -> LCHuv -> (%f,%f,%f,%f)",
					r1, g1, b1, a1, r2, g2, b2, a2)
			}
		})
	}
}

func TestLUVToLCHuvConversion(t *testing.T) {
	// Test that LUV <-> LCHuv conversion is consistent
	luv := NewLUV(50, 30, 40, 1)
	lch := ToLCHuv(luv)

	// Convert back to LUV
	luvBack := NewLUV(lch.L, lch.C*math.Cos(lch.H*math.Pi/180), lch.C*math.Sin(lch.H*math.Pi/180), lch.A_)

	if math.Abs(luv.L-luvBack.L) > 0.01 || math.Abs(luv.U-luvBack.U) > 0.01 || math.Abs(luv.V-luvBack.V) > 0.01 {
		t.Errorf("LUV <-> LCHuv conversion failed: %+v -> %+v -> %+v", luv, lch, luvBack)
	}
}

func TestLUVBlack(t *testing.T) {
	black := RGB(0, 0, 0)
	luv := ToLUV(black)

	if luv.L != 0 {
		t.Errorf("Black L should be 0, got %f", luv.L)
	}
	// U and V should be close to 0 for black (achromatic)
	if math.Abs(luv.U) > 0.01 || math.Abs(luv.V) > 0.01 {
		t.Errorf("Black should have U,V ≈ 0, got U=%f, V=%f", luv.U, luv.V)
	}
}

func TestLUVWhite(t *testing.T) {
	white := RGB(1, 1, 1)
	luv := ToLUV(white)

	if math.Abs(luv.L-100) > 1 {
		t.Errorf("White L should be ≈100, got %f", luv.L)
	}
	// U and V should be close to 0 for white (achromatic)
	if math.Abs(luv.U) > 1 || math.Abs(luv.V) > 1 {
		t.Errorf("White should have U,V ≈ 0, got U=%f, V=%f", luv.U, luv.V)
	}
}

func TestLCHuvHueNormalization(t *testing.T) {
	tests := []struct {
		hue      float64
		expected float64
	}{
		{0, 0},
		{180, 180},
		{360, 0},
		{720, 0},
		{-90, 270},
		{-180, 180},
	}

	for _, tt := range tests {
		lch := NewLCHuv(50, 30, tt.hue, 1)
		if math.Abs(lch.H-tt.expected) > 0.01 {
			t.Errorf("Hue %f normalized to %f, want %f", tt.hue, lch.H, tt.expected)
		}
	}
}

func TestLUVAlpha(t *testing.T) {
	luv := NewLUV(50, 20, 30, 0.5)
	if luv.Alpha() != 0.5 {
		t.Errorf("Alpha() = %f, want 0.5", luv.Alpha())
	}

	withAlpha := luv.WithAlpha(0.8)
	if withAlpha.Alpha() != 0.8 {
		t.Errorf("WithAlpha(0.8).Alpha() = %f, want 0.8", withAlpha.Alpha())
	}
}

func TestLCHuvAlpha(t *testing.T) {
	lch := NewLCHuv(50, 30, 180, 0.6)
	if lch.Alpha() != 0.6 {
		t.Errorf("Alpha() = %f, want 0.6", lch.Alpha())
	}

	withAlpha := lch.WithAlpha(0.9)
	if withAlpha.Alpha() != 0.9 {
		t.Errorf("WithAlpha(0.9).Alpha() = %f, want 0.9", withAlpha.Alpha())
	}
}

func BenchmarkToLUV(b *testing.B) {
	c := RGB(0.5, 0.6, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToLUV(c)
	}
}

func BenchmarkLUVToRGBA(b *testing.B) {
	luv := NewLUV(50, 20, 30, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		luv.RGBA()
	}
}

func BenchmarkToLCHuv(b *testing.B) {
	c := RGB(0.5, 0.6, 0.7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToLCHuv(c)
	}
}

func BenchmarkLCHuvToRGBA(b *testing.B) {
	lch := NewLCHuv(50, 30, 180, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lch.RGBA()
	}
}
