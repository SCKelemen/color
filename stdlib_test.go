package color

import (
	stdcolor "image/color"
	"testing"
)

func TestToStdColor(t *testing.T) {
	// Test conversion to standard library color
	myColor := RGB(1.0, 0.5, 0.0)
	stdColor := ToStdColor(myColor)

	r, g, b, a := stdColor.RGBA()
	// Standard library returns alpha-premultiplied values
	// For RGB(1.0, 0.5, 0.0) with alpha=1.0:
	// Expected: r ≈ 65535, g ≈ 32767, b ≈ 0, a = 65535
	if r < 65500 || r > 65535 {
		t.Errorf("ToStdColor R = %d, want ~65535", r)
	}
	if g < 32700 || g > 32800 {
		t.Errorf("ToStdColor G = %d, want ~32767", g)
	}
	if b != 0 {
		t.Errorf("ToStdColor B = %d, want 0", b)
	}
	if a != 65535 {
		t.Errorf("ToStdColor A = %d, want 65535", a)
	}
}

func TestFromStdColor(t *testing.T) {
	// Test conversion from standard library color
	stdColor := stdcolor.RGBA{R: 255, G: 128, B: 0, A: 255}
	myColor := FromStdColor(stdColor)

	r, g, b, a := myColor.RGBA()
	// Should be approximately RGB(1.0, 0.5, 0.0) with alpha=1.0
	// Note: 128/255 = 0.50196... due to uint8 quantization
	if !floatEqual(r, 1.0) {
		t.Errorf("FromStdColor R = %v, want 1.0", r)
	}
	// Allow tolerance for uint8 quantization (128/255 ≈ 0.502)
	if g < 0.49 || g > 0.51 {
		t.Errorf("FromStdColor G = %v, want ~0.5 (got %v)", g, g)
	}
	if !floatEqual(b, 0.0) {
		t.Errorf("FromStdColor B = %v, want 0.0", b)
	}
	if !floatEqual(a, 1.0) {
		t.Errorf("FromStdColor A = %v, want 1.0", a)
	}
}

func TestStdColorRoundTrip(t *testing.T) {
	// Test round-trip conversion
	original := RGB(0.8, 0.3, 0.7)
	stdColor := ToStdColor(original)
	converted := FromStdColor(stdColor)

	r1, g1, b1, a1 := original.RGBA()
	r2, g2, b2, a2 := converted.RGBA()

	// Allow some tolerance for rounding errors
	if !floatEqual(r1, r2) {
		t.Errorf("Round-trip R: %v != %v", r1, r2)
	}
	if !floatEqual(g1, g2) {
		t.Errorf("Round-trip G: %v != %v", g1, g2)
	}
	if !floatEqual(b1, b2) {
		t.Errorf("Round-trip B: %v != %v", b1, b2)
	}
	if !floatEqual(a1, a2) {
		t.Errorf("Round-trip A: %v != %v", a1, a2)
	}
}

func TestFromStdColorWithAlpha(t *testing.T) {
	// Test conversion with transparency
	stdColor := stdcolor.RGBA{R: 255, G: 0, B: 0, A: 128} // Semi-transparent red
	myColor := FromStdColor(stdColor)

	r, g, b, a := myColor.RGBA()
	// Should be red with alpha ≈ 0.5
	if !floatEqual(r, 1.0) {
		t.Errorf("FromStdColor with alpha R = %v, want 1.0", r)
	}
	if !floatEqual(g, 0.0) {
		t.Errorf("FromStdColor with alpha G = %v, want 0.0", g)
	}
	if !floatEqual(b, 0.0) {
		t.Errorf("FromStdColor with alpha B = %v, want 0.0", b)
	}
	expectedAlpha := 128.0 / 255.0
	if !floatEqual(a, expectedAlpha) {
		t.Errorf("FromStdColor with alpha A = %v, want %v", a, expectedAlpha)
	}
}

func TestToStdColorWithAlpha(t *testing.T) {
	// Test conversion with transparency
	myColor := NewRGBA(1.0, 0.0, 0.0, 0.5) // Semi-transparent red
	stdColor := ToStdColor(myColor)

	r, g, b, a := stdColor.RGBA()
	// Standard library uses alpha-premultiplied values
	// For RGB(1.0, 0.0, 0.0) with alpha=0.5:
	// Expected: r ≈ 32767 (premultiplied), g = 0, b = 0, a = 32767
	expectedR := uint32(32767)
	if r < expectedR-100 || r > expectedR+100 {
		t.Errorf("ToStdColor with alpha R = %d, want ~%d", r, expectedR)
	}
	if g != 0 {
		t.Errorf("ToStdColor with alpha G = %d, want 0", g)
	}
	if b != 0 {
		t.Errorf("ToStdColor with alpha B = %d, want 0", b)
	}
	expectedA := uint32(32767)
	if a < expectedA-100 || a > expectedA+100 {
		t.Errorf("ToStdColor with alpha A = %d, want ~%d", a, expectedA)
	}
}

