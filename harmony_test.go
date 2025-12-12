package color

import (
	"math"
	"testing"
)

func TestComplementary(t *testing.T) {
	tests := []struct {
		name  string
		color Color
	}{
		{"Red", RGB(1, 0, 0)},
		{"Blue", RGB(0, 0, 1)},
		{"Green", RGB(0, 1, 0)},
		{"Yellow", RGB(1, 1, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comp := Complementary(tt.color)
			compOKLCH := ToOKLCH(comp)
			origOKLCH := ToOKLCH(tt.color)

			// Check hue difference is approximately 180°
			hueDiff := math.Abs(compOKLCH.H - origOKLCH.H)
			if hueDiff > 180 {
				hueDiff = 360 - hueDiff
			}

			// Allow wide tolerance since OKLCH hue != HSL hue
			// For some colors, the perceptual complement can differ significantly
			if math.Abs(hueDiff-180) > 60 {
				t.Errorf("Complementary hue difference = %f, want ~180 (orig: %f, comp: %f)",
					hueDiff, origOKLCH.H, compOKLCH.H)
			}

			// Lightness and chroma should be similar
			if math.Abs(compOKLCH.L-origOKLCH.L) > 0.1 {
				t.Errorf("Lightness changed too much: %f -> %f", origOKLCH.L, compOKLCH.L)
			}
		})
	}
}

func TestTriadic(t *testing.T) {
	red := RGB(1, 0, 0)
	colors := Triadic(red)

	if len(colors) != 3 {
		t.Fatalf("Expected 3 colors, got %d", len(colors))
	}

	origOKLCH := ToOKLCH(red)

	// Check that hues are approximately 120° apart
	for i := 0; i < 3; i++ {
		oklch := ToOKLCH(colors[i])
		expectedHue := normalizeHue(origOKLCH.H + float64(i*120))

		hueDiff := math.Abs(oklch.H - expectedHue)
		if hueDiff > 360 {
			hueDiff = 360 - hueDiff
		}

		// Allow wider tolerance for OKLCH hue variations
		if hueDiff > 15 {
			t.Errorf("Color %d: hue = %f, want ~%f (diff = %f)", i, oklch.H, expectedHue, hueDiff)
		}
	}
}

func TestTetradic(t *testing.T) {
	green := RGB(0, 1, 0)
	colors := Tetradic(green)

	if len(colors) != 4 {
		t.Fatalf("Expected 4 colors, got %d", len(colors))
	}

	// All colors should be in gamut
	for i, c := range colors {
		if !InGamut(c) {
			t.Errorf("Color %d is out of gamut", i)
		}
	}

	// First color should be the original
	r1, g1, b1, _ := green.RGBA()
	r2, g2, b2, _ := colors[0].RGBA()

	if math.Abs(r1-r2) > 0.01 || math.Abs(g1-g2) > 0.01 || math.Abs(b1-b2) > 0.01 {
		t.Error("First color should be the original")
	}
}

func TestSquare(t *testing.T) {
	// Square should be same as Tetradic
	green := RGB(0, 1, 0)
	square := Square(green)
	tetradic := Tetradic(green)

	if len(square) != len(tetradic) {
		t.Fatalf("Square and Tetradic should return same number of colors")
	}

	for i := range square {
		r1, g1, b1, _ := square[i].RGBA()
		r2, g2, b2, _ := tetradic[i].RGBA()

		if math.Abs(r1-r2) > 0.01 || math.Abs(g1-g2) > 0.01 || math.Abs(b1-b2) > 0.01 {
			t.Errorf("Square and Tetradic differ at index %d", i)
		}
	}
}

func TestAnalogous(t *testing.T) {
	yellow := RGB(1, 1, 0)
	colors := Analogous(yellow)

	if len(colors) != 3 {
		t.Fatalf("Expected 3 colors, got %d", len(colors))
	}

	// All colors should have similar hues
	hues := make([]float64, len(colors))
	for i, c := range colors {
		oklch := ToOKLCH(c)
		hues[i] = oklch.H
	}

	// Check that hues are within ±30° of center
	centerHue := hues[1]
	for i, h := range hues {
		diff := math.Abs(h - centerHue)
		if diff > 180 {
			diff = 360 - diff
		}
		if diff > 35 {
			t.Errorf("Color %d: hue %f is too far from center %f (diff = %f)", i, h, centerHue, diff)
		}
	}
}

func TestAnalogousN(t *testing.T) {
	orange := RGB(1, 0.5, 0)

	tests := []struct {
		n     int
		angle float64
	}{
		{1, 30},
		{2, 30},
		{5, 45},
		{7, 60},
	}

	for _, tt := range tests {
		colors := AnalogousN(orange, tt.n, tt.angle)

		if len(colors) != tt.n {
			t.Errorf("AnalogousN(%d, %f): got %d colors, want %d", tt.n, tt.angle, len(colors), tt.n)
		}

		if tt.n > 1 {
			// Check spread
			oklch1 := ToOKLCH(colors[0])
			oklchN := ToOKLCH(colors[tt.n-1])

			spread := math.Abs(oklchN.H - oklch1.H)
			if spread > 180 {
				spread = 360 - spread
			}

			expectedSpread := 2 * tt.angle
			// Allow wider tolerance for OKLCH hue variations
			if math.Abs(spread-expectedSpread) > 15 {
				t.Errorf("AnalogousN(%d, %f): spread = %f, want ~%f", tt.n, tt.angle, spread, expectedSpread)
			}
		}
	}
}

func TestSplitComplementary(t *testing.T) {
	purple := RGB(0.5, 0, 0.5)
	colors := SplitComplementary(purple)

	if len(colors) != 3 {
		t.Fatalf("Expected 3 colors, got %d", len(colors))
	}

	origOKLCH := ToOKLCH(purple)
	complementHue := normalizeHue(origOKLCH.H + 180)

	// Colors 1 and 2 should flank the complement
	oklch1 := ToOKLCH(colors[1])
	oklch2 := ToOKLCH(colors[2])

	diff1 := math.Abs(oklch1.H - complementHue)
	if diff1 > 180 {
		diff1 = 360 - diff1
	}
	diff2 := math.Abs(oklch2.H - complementHue)
	if diff2 > 180 {
		diff2 = 360 - diff2
	}

	// Allow wider tolerance for OKLCH hue variations
	if diff1 < 10 || diff1 > 50 {
		t.Errorf("Split complement 1 should be ~30° from complement, got %f", diff1)
	}
	if diff2 < 10 || diff2 > 50 {
		t.Errorf("Split complement 2 should be ~30° from complement, got %f", diff2)
	}
}

func TestMonochromatic(t *testing.T) {
	teal := RGB(0, 0.5, 0.5)

	tests := []int{1, 3, 5, 7}

	for _, n := range tests {
		colors := Monochromatic(teal, n)

		if len(colors) != n {
			t.Errorf("Monochromatic(%d): got %d colors", n, len(colors))
		}

		if n > 1 {
			// All should have same hue
			origOKLCH := ToOKLCH(teal)
			for i, c := range colors {
				oklch := ToOKLCH(c)

				hueDiff := math.Abs(oklch.H - origOKLCH.H)
				if hueDiff > 180 {
					hueDiff = 360 - hueDiff
				}

				if hueDiff > 10 {
					t.Errorf("Color %d: hue changed from %f to %f", i, origOKLCH.H, oklch.H)
				}
			}

			// Lightness should vary
			firstL := ToOKLCH(colors[0]).L
			lastL := ToOKLCH(colors[n-1]).L
			if math.Abs(lastL-firstL) < 0.3 {
				t.Errorf("Lightness range too small: %f to %f", firstL, lastL)
			}
		}
	}
}

func TestMonochromaticCentered(t *testing.T) {
	pink := RGB(1, 0.75, 0.8)
	n := 5

	colors := MonochromaticCentered(pink, n)

	if len(colors) != n {
		t.Fatalf("Expected %d colors, got %d", n, len(colors))
	}

	// Middle color should be closest to original
	origOKLCH := ToOKLCH(pink)
	middleOKLCH := ToOKLCH(colors[n/2])

	lDiff := math.Abs(middleOKLCH.L - origOKLCH.L)
	if lDiff > 0.1 {
		t.Errorf("Middle color lightness should be close to original: %f vs %f", middleOKLCH.L, origOKLCH.L)
	}
}

func TestShades(t *testing.T) {
	red := RGB(1, 0, 0)
	n := 5

	shades := Shades(red, n)

	if len(shades) != n {
		t.Fatalf("Expected %d shades, got %d", n, len(shades))
	}

	// Lightness should decrease
	prevL := 1.0
	for i, shade := range shades {
		oklch := ToOKLCH(shade)
		if oklch.L > prevL {
			t.Errorf("Shade %d: lightness should decrease, got %f after %f", i, oklch.L, prevL)
		}
		prevL = oklch.L
	}

	// Last shade should be much darker
	lastL := ToOKLCH(shades[n-1]).L
	if lastL > 0.3 {
		t.Errorf("Last shade should be dark, got lightness %f", lastL)
	}
}

func TestTints(t *testing.T) {
	blue := RGB(0, 0, 1)
	n := 5

	tints := Tints(blue, n)

	if len(tints) != n {
		t.Fatalf("Expected %d tints, got %d", n, len(tints))
	}

	// Lightness should increase
	prevL := 0.0
	for i, tint := range tints {
		oklch := ToOKLCH(tint)
		if oklch.L < prevL {
			t.Errorf("Tint %d: lightness should increase, got %f after %f", i, oklch.L, prevL)
		}
		prevL = oklch.L
	}

	// Last tint should be much lighter
	lastL := ToOKLCH(tints[n-1]).L
	if lastL < 0.7 {
		t.Errorf("Last tint should be light, got lightness %f", lastL)
	}
}

func TestTones(t *testing.T) {
	green := RGB(0, 1, 0)
	n := 5

	tones := Tones(green, n)

	if len(tones) != n {
		t.Fatalf("Expected %d tones, got %d", n, len(tones))
	}

	// Chroma should decrease
	prevC := 1.0
	for i, tone := range tones {
		oklch := ToOKLCH(tone)
		if oklch.C > prevC {
			t.Errorf("Tone %d: chroma should decrease, got %f after %f", i, oklch.C, prevC)
		}
		prevC = oklch.C
	}

	// Last tone should be nearly gray
	lastC := ToOKLCH(tones[n-1]).C
	if lastC > 0.05 {
		t.Errorf("Last tone should be gray, got chroma %f", lastC)
	}
}

func TestRectangle(t *testing.T) {
	magenta := RGB(1, 0, 1)
	colors := Rectangle(magenta, 60)

	if len(colors) != 4 {
		t.Fatalf("Expected 4 colors, got %d", len(colors))
	}

	// All colors should be in gamut
	for i, c := range colors {
		if !InGamut(c) {
			t.Errorf("Color %d is out of gamut", i)
		}
	}

	// First color should be the original
	r1, g1, b1, _ := magenta.RGBA()
	r2, g2, b2, _ := colors[0].RGBA()

	if math.Abs(r1-r2) > 0.01 || math.Abs(g1-g2) > 0.01 || math.Abs(b1-b2) > 0.01 {
		t.Error("First color should be the original")
	}
}

func TestDoubleSplitComplementary(t *testing.T) {
	cyan := RGB(0, 1, 1)
	colors := DoubleSplitComplementary(cyan)

	if len(colors) != 6 {
		t.Fatalf("Expected 6 colors, got %d", len(colors))
	}

	// All colors should be in gamut
	for i, c := range colors {
		if !InGamut(c) {
			t.Errorf("Color %d is out of gamut", i)
		}
	}
}

func TestHarmonyEdgeCases(t *testing.T) {
	// Test with zero steps
	colors := Monochromatic(RGB(1, 0, 0), 0)
	if len(colors) != 0 {
		t.Errorf("Monochromatic(0) should return empty slice")
	}

	colors = AnalogousN(RGB(0, 1, 0), 0, 30)
	if len(colors) != 0 {
		t.Errorf("AnalogousN(0) should return empty slice")
	}

	// Test with one step
	colors = Monochromatic(RGB(0, 0, 1), 1)
	if len(colors) != 1 {
		t.Errorf("Monochromatic(1) should return one color")
	}
}
