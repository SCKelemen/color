package color

import (
	"math"
	"testing"
)

func TestContrastRatio(t *testing.T) {
	tests := []struct {
		name     string
		c1       Color
		c2       Color
		minRatio float64
		maxRatio float64
	}{
		{
			name:     "Black on white (maximum contrast)",
			c1:       RGB(0, 0, 0),
			c2:       RGB(1, 1, 1),
			minRatio: 20.9,
			maxRatio: 21.1,
		},
		{
			name:     "White on black (should be same as black on white)",
			c1:       RGB(1, 1, 1),
			c2:       RGB(0, 0, 0),
			minRatio: 20.9,
			maxRatio: 21.1,
		},
		{
			name:     "Same color (minimum contrast)",
			c1:       RGB(0.5, 0.5, 0.5),
			c2:       RGB(0.5, 0.5, 0.5),
			minRatio: 0.99,
			maxRatio: 1.01,
		},
		{
			name:     "Dark gray on black",
			c1:       RGB(0.2, 0.2, 0.2),
			c2:       RGB(0, 0, 0),
			minRatio: 1.5,
			maxRatio: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ratio := ContrastRatio(tt.c1, tt.c2)

			if ratio < tt.minRatio || ratio > tt.maxRatio {
				t.Errorf("ContrastRatio() = %f, want between %f and %f", ratio, tt.minRatio, tt.maxRatio)
			}
		})
	}
}

func TestContrastRatioSymmetry(t *testing.T) {
	// Contrast ratio should be symmetric
	c1 := RGB(1, 0, 0)
	c2 := RGB(0, 0, 1)

	ratio1 := ContrastRatio(c1, c2)
	ratio2 := ContrastRatio(c2, c1)

	if math.Abs(ratio1-ratio2) > 0.001 {
		t.Errorf("ContrastRatio not symmetric: %f vs %f", ratio1, ratio2)
	}
}

func TestCheckContrast(t *testing.T) {
	tests := []struct {
		name       string
		fg         Color
		bg         Color
		wantAA     bool
		wantAAA    bool
		wantAALg   bool
		wantAAALg  bool
		wantUI     bool
	}{
		{
			name:       "Black on white - passes everything",
			fg:         RGB(0, 0, 0),
			bg:         RGB(1, 1, 1),
			wantAA:     true,
			wantAAA:    true,
			wantAALg:   true,
			wantAAALg:  true,
			wantUI:     true,
		},
		{
			name:       "Light gray on white - fails most",
			fg:         RGB(0.8, 0.8, 0.8),
			bg:         RGB(1, 1, 1),
			wantAA:     false,
			wantAAA:    false,
			wantAALg:   false,
			wantAAALg:  false,
			wantUI:     false,
		},
		{
			name:       "Medium contrast",
			fg:         RGB(0.4, 0.4, 0.4),
			bg:         RGB(1, 1, 1),
			wantAA:     true,  // Should pass AA normal
			wantAAA:    false, // Should fail AAA normal
			wantAALg:   true,
			wantAAALg:  true,  // Should pass AAA large
			wantUI:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckContrast(tt.fg, tt.bg)

			if result.AANormal != tt.wantAA {
				t.Errorf("AANormal = %v, want %v (ratio: %f)", result.AANormal, tt.wantAA, result.Ratio)
			}
			if result.AAANormal != tt.wantAAA {
				t.Errorf("AAANormal = %v, want %v (ratio: %f)", result.AAANormal, tt.wantAAA, result.Ratio)
			}
			if result.AALarge != tt.wantAALg {
				t.Errorf("AALarge = %v, want %v (ratio: %f)", result.AALarge, tt.wantAALg, result.Ratio)
			}
			if result.AAALarge != tt.wantAAALg {
				t.Errorf("AAALarge = %v, want %v (ratio: %f)", result.AAALarge, tt.wantAAALg, result.Ratio)
			}
			if result.UIComponents != tt.wantUI {
				t.Errorf("UIComponents = %v, want %v (ratio: %f)", result.UIComponents, tt.wantUI, result.Ratio)
			}
		})
	}
}

func TestIsAccessible(t *testing.T) {
	tests := []struct {
		name     string
		fg       Color
		bg       Color
		level    WCAGLevel
		textSize TextSize
		want     bool
	}{
		{
			name:     "Black on white, AA normal",
			fg:       RGB(0, 0, 0),
			bg:       RGB(1, 1, 1),
			level:    WCAGAA,
			textSize: NormalText,
			want:     true,
		},
		{
			name:     "Black on white, AAA normal",
			fg:       RGB(0, 0, 0),
			bg:       RGB(1, 1, 1),
			level:    WCAGAAA,
			textSize: NormalText,
			want:     true,
		},
		{
			name:     "Light gray on white, AA normal",
			fg:       RGB(0.7, 0.7, 0.7),
			bg:       RGB(1, 1, 1),
			level:    WCAGAA,
			textSize: NormalText,
			want:     false,
		},
		{
			name:     "Medium gray on white, AA large",
			fg:       RGB(0.5, 0.5, 0.5),
			bg:       RGB(1, 1, 1),
			level:    WCAGAA,
			textSize: LargeText,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAccessible(tt.fg, tt.bg, tt.level, tt.textSize)

			if result != tt.want {
				ratio := ContrastRatio(tt.fg, tt.bg)
				t.Errorf("IsAccessible() = %v, want %v (ratio: %f)", result, tt.want, ratio)
			}
		})
	}
}

func TestSuggestAccessibleForeground(t *testing.T) {
	tests := []struct {
		name     string
		base     Color
		bg       Color
		level    WCAGLevel
		textSize TextSize
	}{
		{
			name:     "Suggest for white background, AA normal",
			base:     RGB(0.5, 0.3, 0.7),
			bg:       RGB(1, 1, 1),
			level:    WCAGAA,
			textSize: NormalText,
		},
		{
			name:     "Suggest for black background, AA normal",
			base:     RGB(0.5, 0.3, 0.7),
			bg:       RGB(0, 0, 0),
			level:    WCAGAA,
			textSize: NormalText,
		},
		{
			name:     "Suggest for gray background, AAA normal",
			base:     RGB(1, 0.5, 0),
			bg:       RGB(0.5, 0.5, 0.5),
			level:    WCAGAAA,
			textSize: NormalText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggested, success := SuggestAccessibleForeground(tt.base, tt.bg, tt.level, tt.textSize)

			// Verify the suggestion is accessible
			isOk := IsAccessible(suggested, tt.bg, tt.level, tt.textSize)

			if success && !isOk {
				ratio := ContrastRatio(suggested, tt.bg)
				t.Errorf("Suggested color is not accessible (ratio: %f)", ratio)
			}

			// Suggested color should preserve the hue
			baseOKLCH := ToOKLCH(tt.base)
			suggestedOKLCH := ToOKLCH(suggested)

			hueDiff := math.Abs(baseOKLCH.H - suggestedOKLCH.H)
			if hueDiff > 180 {
				hueDiff = 360 - hueDiff
			}

			if hueDiff > 30 {
				t.Errorf("Suggested color changed hue too much: %f -> %f (diff: %f)",
					baseOKLCH.H, suggestedOKLCH.H, hueDiff)
			}
		})
	}
}

func TestSuggestAccessibleBackground(t *testing.T) {
	fg := RGB(0.2, 0.2, 0.2)
	baseBg := RGB(0.3, 0.3, 0.3)

	suggested, success := SuggestAccessibleBackground(fg, baseBg, WCAGAA, NormalText)

	if !success {
		t.Skip("Could not find accessible background")
	}

	if !IsAccessible(fg, suggested, WCAGAA, NormalText) {
		ratio := ContrastRatio(fg, suggested)
		t.Errorf("Suggested background is not accessible (ratio: %f)", ratio)
	}
}

func TestSimulateColorBlindness(t *testing.T) {
	red := RGB(1, 0, 0)

	types := []ColorBlindnessType{
		Protanopia,
		Protanomaly,
		Deuteranopia,
		Deuteranomaly,
		Tritanopia,
		Tritanomaly,
		Achromatopsia,
		Achromatomaly,
	}

	for _, cvdType := range types {
		t.Run(cvdType.String(), func(t *testing.T) {
			simulated := SimulateColorBlindness(red, cvdType)

			// Result should be valid
			r, g, b, a := simulated.RGBA()
			if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
				t.Errorf("Simulated color out of range: (%f, %f, %f, %f)", r, g, b, a)
			}

			// Alpha should be preserved
			_, _, _, origA := red.RGBA()
			if math.Abs(a-origA) > 0.01 {
				t.Errorf("Alpha not preserved: %f -> %f", origA, a)
			}

			// For achromatopsia, result should be gray
			if cvdType == Achromatopsia {
				if math.Abs(r-g) > 0.01 || math.Abs(g-b) > 0.01 {
					t.Errorf("Achromatopsia should produce gray, got (%f, %f, %f)", r, g, b)
				}
			}
		})
	}
}

func TestSimulateColorBlindnessPreservesGray(t *testing.T) {
	gray := RGB(0.5, 0.5, 0.5)

	types := []ColorBlindnessType{
		Protanopia,
		Deuteranopia,
		Tritanopia,
		Achromatopsia,
	}

	for _, cvdType := range types {
		simulated := SimulateColorBlindness(gray, cvdType)
		r, g, b, _ := simulated.RGBA()

		// Gray should remain gray for all CVD types
		if math.Abs(r-g) > 0.05 || math.Abs(g-b) > 0.05 {
			t.Errorf("%v: Gray changed to (%f, %f, %f)", cvdType, r, g, b)
		}
	}
}

func TestIsColorBlindSafe(t *testing.T) {
	red := RGB(1, 0, 0)
	green := RGB(0, 1, 0)

	// Red and green should NOT be safe for protanopia/deuteranopia
	if IsColorBlindSafe(red, green, Protanopia, 3.0) {
		t.Error("Red and green should not be safe for protanopia")
	}

	if IsColorBlindSafe(red, green, Deuteranopia, 3.0) {
		t.Error("Red and green should not be safe for deuteranopia")
	}

	// Blue and yellow should be safer
	blue := RGB(0, 0, 1)
	yellow := RGB(1, 1, 0)

	if !IsColorBlindSafe(blue, yellow, Protanopia, 2.0) {
		t.Error("Blue and yellow should be safe for protanopia")
	}

	// Black and white should always be safe
	black := RGB(0, 0, 0)
	white := RGB(1, 1, 1)

	for _, cvdType := range []ColorBlindnessType{Protanopia, Deuteranopia, Tritanopia} {
		if !IsColorBlindSafe(black, white, cvdType, 4.5) {
			t.Errorf("Black and white should be safe for %v", cvdType)
		}
	}
}

func TestCheckColorBlindSafety(t *testing.T) {
	red := RGB(1, 0, 0)
	green := RGB(0, 1, 0)

	results := CheckColorBlindSafety(red, green, 3.0)

	// Should have results for all common types
	expectedTypes := []ColorBlindnessType{
		Protanopia,
		Protanomaly,
		Deuteranopia,
		Deuteranomaly,
		Tritanopia,
		Tritanomaly,
	}

	for _, cvdType := range expectedTypes {
		if _, ok := results[cvdType]; !ok {
			t.Errorf("Missing result for %v", cvdType)
		}
	}

	// Red and green should fail for red-green color blindness
	if results[Protanopia] {
		t.Error("Red and green should not be safe for Protanopia")
	}
	if results[Deuteranopia] {
		t.Error("Red and green should not be safe for Deuteranopia")
	}
}

func TestRelativeLuminance(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		expected float64
		tolerance float64
	}{
		{"Black", RGB(0, 0, 0), 0.0, 0.001},
		{"White", RGB(1, 1, 1), 1.0, 0.001},
		{"Red", RGB(1, 0, 0), 0.2126, 0.01},
		{"Green", RGB(0, 1, 0), 0.7152, 0.01},
		{"Blue", RGB(0, 0, 1), 0.0722, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lum := relativeLuminance(tt.color)

			if math.Abs(lum-tt.expected) > tt.tolerance {
				t.Errorf("relativeLuminance() = %f, want %f ± %f", lum, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestSRGBLinearConversion(t *testing.T) {
	// Test round-trip
	values := []float64{0, 0.04045, 0.5, 1.0}

	for _, v := range values {
		linear := sRGBToLinear(v)
		srgb := linearToSRGB(linear)

		if math.Abs(srgb-v) > 0.001 {
			t.Errorf("Round trip failed: %f -> %f -> %f", v, linear, srgb)
		}
	}
}

// Add String method for ColorBlindnessType for better test output
func (cvd ColorBlindnessType) String() string {
	names := map[ColorBlindnessType]string{
		Protanopia:     "Protanopia",
		Protanomaly:    "Protanomaly",
		Deuteranopia:   "Deuteranopia",
		Deuteranomaly:  "Deuteranomaly",
		Tritanopia:     "Tritanopia",
		Tritanomaly:    "Tritanomaly",
		Achromatopsia:  "Achromatopsia",
		Achromatomaly:  "Achromatomaly",
	}
	if name, ok := names[cvd]; ok {
		return name
	}
	return "Unknown"
}
