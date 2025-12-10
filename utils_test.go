package color

import (
	"math"
	"testing"
)

func TestLighten(t *testing.T) {
	blue := RGB(0, 0, 1)
	lightBlue := Lighten(blue, 0.5)

	// Lightened color should have higher lightness
	oklch1 := ToOKLCH(blue)
	oklch2 := ToOKLCH(lightBlue)

	if oklch2.L <= oklch1.L {
		t.Errorf("Lighten should increase lightness: L=%v -> L=%v", oklch1.L, oklch2.L)
	}

	// Lightening by 0 should not change the color
	unchanged := Lighten(blue, 0)
	r1, g1, b1, _ := blue.RGBA()
	r2, g2, b2, _ := unchanged.RGBA()
	if !rgbaEqual(r1, g1, b1, 1, r2, g2, b2, 1) {
		t.Errorf("Lighten(0) should not change color")
	}
}

func TestDarken(t *testing.T) {
	red := RGB(1, 0, 0)
	darkRed := Darken(red, 0.5)

	// Darkened color should have lower lightness
	oklch1 := ToOKLCH(red)
	oklch2 := ToOKLCH(darkRed)

	if oklch2.L >= oklch1.L {
		t.Errorf("Darken should decrease lightness: L=%v -> L=%v", oklch1.L, oklch2.L)
	}

	// Darkening by 0 should not change the color
	unchanged := Darken(red, 0)
	r1, g1, b1, _ := red.RGBA()
	r2, g2, b2, _ := unchanged.RGBA()
	if !rgbaEqual(r1, g1, b1, 1, r2, g2, b2, 1) {
		t.Errorf("Darken(0) should not change color")
	}
}

func TestSaturate(t *testing.T) {
	// Start with a desaturated color
	gray := RGB(0.5, 0.5, 0.5)
	saturated := Saturate(gray, 0.5)

	oklch1 := ToOKLCH(gray)
	oklch2 := ToOKLCH(saturated)

	// Saturated color should have higher chroma
	if oklch2.C <= oklch1.C {
		t.Errorf("Saturate should increase chroma: C=%v -> C=%v", oklch1.C, oklch2.C)
	}
}

func TestDesaturate(t *testing.T) {
	red := RGB(1, 0, 0)
	desaturated := Desaturate(red, 0.5)

	oklch1 := ToOKLCH(red)
	oklch2 := ToOKLCH(desaturated)

	// Desaturated color should have lower chroma
	if oklch2.C >= oklch1.C {
		t.Errorf("Desaturate should decrease chroma: C=%v -> C=%v", oklch1.C, oklch2.C)
	}

	// Desaturating by 1 should make it grayscale
	grayscale := Desaturate(red, 1.0)
	oklch3 := ToOKLCH(grayscale)
	if oklch3.C > 0.01 { // Allow small tolerance
		t.Errorf("Desaturate(1.0) should make color grayscale, but C=%v", oklch3.C)
	}
}

func TestMix(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Mix 50/50
	mixed := Mix(red, blue, 0.5)
	r, g, b, _ := mixed.RGBA()

	// Should be approximately purple (0.5, 0, 0.5)
	if !floatEqual(r, 0.5) || !floatEqual(b, 0.5) || !floatEqual(g, 0) {
		t.Errorf("Mix(red, blue, 0.5) = RGB(%v, %v, %v), want RGB(0.5, 0, 0.5)", r, g, b)
	}

	// Mix with weight 0 should return first color
	unchanged := Mix(red, blue, 0)
	r1, g1, b1, _ := red.RGBA()
	r2, g2, b2, _ := unchanged.RGBA()
	if !rgbaEqual(r1, g1, b1, 1, r2, g2, b2, 1) {
		t.Errorf("Mix(weight=0) should return first color")
	}

	// Mix with weight 1 should return second color
	unchanged2 := Mix(red, blue, 1)
	r3, g3, b3, _ := blue.RGBA()
	r4, g4, b4, _ := unchanged2.RGBA()
	if !rgbaEqual(r3, g3, b3, 1, r4, g4, b4, 1) {
		t.Errorf("Mix(weight=1) should return second color")
	}
}

func TestMixOKLCH(t *testing.T) {
	red := RGB(1, 0, 0)
	blue := RGB(0, 0, 1)

	// Mix 50/50 in OKLCH space
	mixed := MixOKLCH(red, blue, 0.5)

	// Should be a valid color
	r, g, b, a := mixed.RGBA()
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
		t.Errorf("MixOKLCH produced invalid color: RGB(%v, %v, %v, %v)", r, g, b, a)
	}

	// Alpha should be mixed
	if !floatEqual(a, 1.0) {
		t.Errorf("MixOKLCH alpha = %v, want 1.0", a)
	}
}

func TestAdjustHue(t *testing.T) {
	red := RGB(1, 0, 0)
	shifted := AdjustHue(red, 120) // Shift by 120 degrees (should be green-ish)

	oklch1 := ToOKLCH(red)
	oklch2 := ToOKLCH(shifted)

	// Hue should be shifted (accounting for wraparound)
	expectedHue := math.Mod(oklch1.H+120, 360)
	actualHue := oklch2.H
	diff := math.Abs(expectedHue - actualHue)
	// Account for wraparound - take the smaller of the two possible differences
	if diff > 180 {
		diff = 360 - diff
	}
	if diff > 20 { // Allow tolerance for color space conversion and gamut mapping
		t.Errorf("AdjustHue(120) shifted hue from %v to %v, expected ~%v (diff: %v)", oklch1.H, actualHue, expectedHue, diff)
	}
}

func TestInvert(t *testing.T) {
	white := RGB(1, 1, 1)
	black := Invert(white)
	r, g, b, _ := black.RGBA()
	if !rgbaEqual(r, g, b, 1, 0, 0, 0, 1) {
		t.Errorf("Invert(white) = RGB(%v, %v, %v), want RGB(0, 0, 0)", r, g, b)
	}

	red := RGB(1, 0, 0)
	cyan := Invert(red)
	r2, g2, b2, _ := cyan.RGBA()
	if !floatEqual(r2, 0) || !floatEqual(g2, 1) || !floatEqual(b2, 1) {
		t.Errorf("Invert(red) = RGB(%v, %v, %v), want RGB(0, 1, 1)", r2, g2, b2)
	}
}

func TestGrayscale(t *testing.T) {
	red := RGB(1, 0, 0)
	gray := Grayscale(red)
	r, g, b, _ := gray.RGBA()

	// All channels should be equal (grayscale)
	if !floatEqual(r, g) || !floatEqual(g, b) {
		t.Errorf("Grayscale(red) = RGB(%v, %v, %v), all channels should be equal", r, g, b)
	}

	// Should have reasonable luminance
	if r < 0.2 || r > 0.4 {
		t.Errorf("Grayscale(red) luminance seems off: %v", r)
	}
}

func TestComplement(t *testing.T) {
	red := RGB(1, 0, 0)
	complement := Complement(red)

	oklch1 := ToOKLCH(red)
	oklch2 := ToOKLCH(complement)

	// Complementary color should be 180 degrees away
	expectedHue := math.Mod(oklch1.H+180, 360)
	diff := math.Abs(oklch2.H - expectedHue)
	// Account for wraparound
	if diff > 180 {
		diff = 360 - diff
	}
	if diff > 20 { // Allow tolerance for color space conversion and gamut mapping
		t.Errorf("Complement hue = %v, expected ~%v (180 degrees from %v, diff: %v)", oklch2.H, expectedHue, oklch1.H, diff)
	}
}

func TestOpacity(t *testing.T) {
	red := RGB(1, 0, 0)
	semiTransparent := Opacity(red, 0.5)

	if !floatEqual(semiTransparent.Alpha(), 0.5) {
		t.Errorf("Opacity(0.5).Alpha() = %v, want 0.5", semiTransparent.Alpha())
	}

	// RGB should be unchanged
	r1, g1, b1, _ := red.RGBA()
	r2, g2, b2, _ := semiTransparent.RGBA()
	if !floatEqual(r1, r2) || !floatEqual(g1, g2) || !floatEqual(b1, b2) {
		t.Errorf("Opacity should not change RGB values")
	}
}

func TestFadeOut(t *testing.T) {
	red := NewRGBA(1, 0, 0, 1.0)
	faded := FadeOut(red, 0.3)

	// Should reduce alpha by 30%
	expected := 1.0 * (1 - 0.3)
	if !floatEqual(faded.Alpha(), expected) {
		t.Errorf("FadeOut(0.3).Alpha() = %v, want %v", faded.Alpha(), expected)
	}
}

func TestFadeIn(t *testing.T) {
	red := NewRGBA(1, 0, 0, 0.5)
	moreOpaque := FadeIn(red, 0.3)

	// Should increase alpha
	if moreOpaque.Alpha() <= red.Alpha() {
		t.Errorf("FadeIn should increase alpha: %v -> %v", red.Alpha(), moreOpaque.Alpha())
	}
}

