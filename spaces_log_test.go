package color

import (
	"math"
	"testing"
)

// Test LOG transfer function round-trips
func TestCLogRoundTrip(t *testing.T) {
	testCases := []float64{0.0, 0.01, 0.1, 0.18, 0.5, 0.8, 1.0, 2.0, 5.0}

	for _, linear := range testCases {
		encoded := cLogTransfer(linear)
		decoded := cLogInverseTransfer(encoded)

		if !floatNear(linear, decoded, 1e-6) {
			t.Errorf("C-Log round-trip failed: %v -> %v -> %v", linear, encoded, decoded)
		}
	}
}

func TestSLog3RoundTrip(t *testing.T) {
	testCases := []float64{0.0, 0.01, 0.1, 0.18, 0.5, 0.8, 1.0, 2.0, 5.0}

	for _, linear := range testCases {
		encoded := sLog3Transfer(linear)
		decoded := sLog3InverseTransfer(encoded)

		if !floatNear(linear, decoded, 1e-6) {
			t.Errorf("S-Log3 round-trip failed: %v -> %v -> %v", linear, encoded, decoded)
		}
	}
}

func TestVLogRoundTrip(t *testing.T) {
	testCases := []float64{0.0, 0.01, 0.1, 0.18, 0.5, 0.8, 1.0, 2.0, 5.0}

	for _, linear := range testCases {
		encoded := vLogTransfer(linear)
		decoded := vLogInverseTransfer(encoded)

		if !floatNear(linear, decoded, 1e-6) {
			t.Errorf("V-Log round-trip failed: %v -> %v -> %v", linear, encoded, decoded)
		}
	}
}

func TestArriLogCRoundTrip(t *testing.T) {
	testCases := []float64{0.0, 0.01, 0.1, 0.18, 0.5, 0.8, 1.0, 2.0, 5.0}

	for _, linear := range testCases {
		encoded := arriLogCTransfer(linear)
		decoded := arriLogCInverseTransfer(encoded)

		if !floatNear(linear, decoded, 1e-6) {
			t.Errorf("Arri LogC round-trip failed: %v -> %v -> %v", linear, encoded, decoded)
		}
	}
}

func TestRedLog3G10RoundTrip(t *testing.T) {
	testCases := []float64{0.0, 0.01, 0.1, 0.18, 0.5, 0.8, 1.0, 2.0, 5.0}

	for _, linear := range testCases {
		encoded := redLog3G10Transfer(linear)
		decoded := redLog3G10InverseTransfer(encoded)

		if !floatNear(linear, decoded, 1e-6) {
			t.Errorf("Red Log3G10 round-trip failed: %v -> %v -> %v", linear, encoded, decoded)
		}
	}
}

func TestBMDFilmRoundTrip(t *testing.T) {
	testCases := []float64{0.0, 0.01, 0.1, 0.18, 0.5, 0.8, 1.0}

	for _, linear := range testCases {
		encoded := bmdFilmTransfer(linear)
		decoded := bmdFilmInverseTransfer(encoded)

		if !floatNear(linear, decoded, 1e-6) {
			t.Errorf("BMD Film round-trip failed: %v -> %v -> %v", linear, encoded, decoded)
		}
	}
}

// Test LOG transfer functions handle edge cases
func TestCLogEdgeCases(t *testing.T) {
	// Test negative values
	if cLogTransfer(-0.1) != 0.0 {
		t.Errorf("C-Log should return 0 for negative linear values")
	}
	if cLogInverseTransfer(-0.1) != 0.0 {
		t.Errorf("C-Log inverse should return 0 for negative encoded values")
	}

	// Test black (0.0)
	encoded := cLogTransfer(0.0)
	if encoded < 0 {
		t.Errorf("C-Log encoding of black should not be negative: %v", encoded)
	}
}

func TestSLog3EdgeCases(t *testing.T) {
	// Test negative values - S-Log3 can handle negative linear values
	if sLog3Transfer(-0.1) < 0 {
		t.Errorf("S-Log3 should handle negative linear values gracefully")
	}

	// Test cut point behavior (linear near black)
	linearCut := 0.01125
	encoded := sLog3Transfer(linearCut)
	decoded := sLog3InverseTransfer(encoded)
	if !floatNear(linearCut, decoded, 1e-5) {
		t.Errorf("S-Log3 cut point behavior incorrect: %v -> %v -> %v", linearCut, encoded, decoded)
	}
}

func TestVLogEdgeCases(t *testing.T) {
	// Test negative values - V-Log can handle negative linear values
	if vLogTransfer(-0.1) < 0 {
		t.Errorf("V-Log should handle negative linear values gracefully")
	}

	// Test cut point behavior
	cut := 0.01
	encoded := vLogTransfer(cut)
	decoded := vLogInverseTransfer(encoded)
	if !floatNear(cut, decoded, 1e-5) {
		t.Errorf("V-Log cut point behavior incorrect: %v -> %v -> %v", cut, encoded, decoded)
	}
}

// Test LOG color space conversions
func TestCLogSpaceConversion(t *testing.T) {
	// Create a color in C-Log space
	clogColor := NewSpaceColor(CLogSpace, []float64{0.5, 0.4, 0.3}, 1.0)

	if clogColor.Space().Name() != "c-log" {
		t.Errorf("Expected space 'c-log', got '%s'", clogColor.Space().Name())
	}

	// Convert to sRGB
	srgbColor := clogColor.ConvertTo(SRGBSpace)
	if srgbColor.Space().Name() != "sRGB" {
		t.Errorf("Expected space 'sRGB', got '%s'", srgbColor.Space().Name())
	}

	// Convert back to C-Log
	clogBack := srgbColor.ConvertTo(CLogSpace)

	// Should be close to original (allowing for gamut mapping and rounding)
	// Note: Wider gamut LOG spaces may not round-trip perfectly through sRGB
	origChannels := clogColor.Channels()
	backChannels := clogBack.Channels()

	for i := 0; i < 3; i++ {
		// Use larger tolerance due to gamut clipping (Cinema Gamut is very wide)
		if !floatNear(origChannels[i], backChannels[i], 0.06) {
			t.Errorf("C-Log round-trip channel %d: %v -> %v", i, origChannels[i], backChannels[i])
		}
	}
}

func TestSLog3SpaceConversion(t *testing.T) {
	// Create a color in S-Log3 space
	slog3Color := NewSpaceColor(SLog3Space, []float64{0.5, 0.4, 0.3}, 1.0)

	if slog3Color.Space().Name() != "s-log3" {
		t.Errorf("Expected space 's-log3', got '%s'", slog3Color.Space().Name())
	}

	// Convert to sRGB
	srgbColor := slog3Color.ConvertTo(SRGBSpace)

	// Convert back to S-Log3
	slog3Back := srgbColor.ConvertTo(SLog3Space)

	// Should be close to original (allowing for gamut mapping)
	origChannels := slog3Color.Channels()
	backChannels := slog3Back.Channels()

	for i := 0; i < 3; i++ {
		// Use larger tolerance due to gamut clipping
		if !floatNear(origChannels[i], backChannels[i], 0.35) {
			t.Errorf("S-Log3 round-trip channel %d: %v -> %v", i, origChannels[i], backChannels[i])
		}
	}
}

func TestVLogSpaceConversion(t *testing.T) {
	// Create a color in V-Log space
	vlogColor := NewSpaceColor(VLogSpace, []float64{0.5, 0.4, 0.3}, 1.0)

	if vlogColor.Space().Name() != "v-log" {
		t.Errorf("Expected space 'v-log', got '%s'", vlogColor.Space().Name())
	}

	// Convert to Display P3
	p3Color := vlogColor.ConvertTo(DisplayP3Space)

	// Convert back to V-Log
	vlogBack := p3Color.ConvertTo(VLogSpace)

	// Should be close to original (allowing for gamut mapping)
	origChannels := vlogColor.Channels()
	backChannels := vlogBack.Channels()

	for i := 0; i < 3; i++ {
		// Use larger tolerance due to gamut clipping
		if !floatNear(origChannels[i], backChannels[i], 0.05) {
			t.Errorf("V-Log round-trip channel %d: %v -> %v", i, origChannels[i], backChannels[i])
		}
	}
}

func TestArriLogCSpaceConversion(t *testing.T) {
	// Create a color in Arri LogC space
	logcColor := NewSpaceColor(ArriLogCSpace, []float64{0.5, 0.4, 0.3}, 1.0)

	if logcColor.Space().Name() != "arri-logc" {
		t.Errorf("Expected space 'arri-logc', got '%s'", logcColor.Space().Name())
	}

	// Convert to sRGB
	srgbColor := logcColor.ConvertTo(SRGBSpace)

	// Convert back to Arri LogC
	logcBack := srgbColor.ConvertTo(ArriLogCSpace)

	// Should be close to original (allowing for gamut mapping)
	origChannels := logcColor.Channels()
	backChannels := logcBack.Channels()

	for i := 0; i < 3; i++ {
		// Use larger tolerance due to gamut clipping
		if !floatNear(origChannels[i], backChannels[i], 0.05) {
			t.Errorf("Arri LogC round-trip channel %d: %v -> %v", i, origChannels[i], backChannels[i])
		}
	}
}

func TestRedLog3G10SpaceConversion(t *testing.T) {
	// Create a color in Red Log3G10 space
	redlogColor := NewSpaceColor(RedLog3G10Space, []float64{0.5, 0.4, 0.3}, 1.0)

	if redlogColor.Space().Name() != "red-log3g10" {
		t.Errorf("Expected space 'red-log3g10', got '%s'", redlogColor.Space().Name())
	}

	// Convert to Rec.2020
	rec2020Color := redlogColor.ConvertTo(Rec2020Space)

	// Convert back to Red Log3G10
	redlogBack := rec2020Color.ConvertTo(RedLog3G10Space)

	// Should be close to original (allowing for gamut mapping)
	origChannels := redlogColor.Channels()
	backChannels := redlogBack.Channels()

	for i := 0; i < 3; i++ {
		// Use larger tolerance due to gamut clipping
		if !floatNear(origChannels[i], backChannels[i], 0.05) {
			t.Errorf("Red Log3G10 round-trip channel %d: %v -> %v", i, origChannels[i], backChannels[i])
		}
	}
}

func TestBMDFilmSpaceConversion(t *testing.T) {
	// Create a color in BMD Film space
	bmdColor := NewSpaceColor(BMDFilmSpace, []float64{0.5, 0.4, 0.3}, 1.0)

	if bmdColor.Space().Name() != "bmd-film" {
		t.Errorf("Expected space 'bmd-film', got '%s'", bmdColor.Space().Name())
	}

	// Convert to sRGB
	srgbColor := bmdColor.ConvertTo(SRGBSpace)

	// Convert back to BMD Film
	bmdBack := srgbColor.ConvertTo(BMDFilmSpace)

	// Should be close to original (allowing for gamut mapping)
	origChannels := bmdColor.Channels()
	backChannels := bmdBack.Channels()

	for i := 0; i < 3; i++ {
		// Use larger tolerance due to gamut clipping
		if !floatNear(origChannels[i], backChannels[i], 0.05) {
			t.Errorf("BMD Film round-trip channel %d: %v -> %v", i, origChannels[i], backChannels[i])
		}
	}
}

// Test registry lookups
func TestLOGSpaceRegistry(t *testing.T) {
	tests := []struct {
		name     string
		expected Space
	}{
		{"c-log", CLogSpace},
		{"clog", CLogSpace},
		{"s-log3", SLog3Space},
		{"slog3", SLog3Space},
		{"v-log", VLogSpace},
		{"vlog", VLogSpace},
		{"arri-logc", ArriLogCSpace},
		{"logc", ArriLogCSpace},
		{"red-log3g10", RedLog3G10Space},
		{"log3g10", RedLog3G10Space},
		{"bmd-film", BMDFilmSpace},
		{"bmdfilm", BMDFilmSpace},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			space, ok := GetSpace(tt.name)
			if !ok {
				t.Errorf("Failed to get space '%s' from registry", tt.name)
			}
			if space != tt.expected {
				t.Errorf("Registry returned wrong space for '%s'", tt.name)
			}
		})
	}
}

// Test case-insensitive lookups
func TestLOGSpaceRegistryCaseInsensitive(t *testing.T) {
	tests := []string{
		"C-Log", "CLOG", "c-LOG",
		"S-Log3", "SLOG3", "s-LOG3",
		"V-Log", "VLOG", "v-LOG",
		"Arri-LogC", "LOGC", "arri-LOGC",
		"Red-Log3G10", "LOG3G10", "red-LOG3G10",
		"BMD-Film", "BMDFILM", "bmd-FILM",
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			space, ok := GetSpace(name)
			if !ok {
				t.Errorf("Failed to get space '%s' from registry (case-insensitive)", name)
			}
			if space == nil {
				t.Errorf("Registry returned nil space for '%s'", name)
			}
		})
	}
}

// Test known reference values from camera specs
func TestCLogReferenceValues(t *testing.T) {
	// Canon C-Log reference: 18% gray should encode to approximately 0.34
	gray18Linear := 0.18
	encoded := cLogTransfer(gray18Linear)

	// Allow some tolerance for specification variations
	if encoded < 0.30 || encoded > 0.38 {
		t.Errorf("C-Log 18%% gray: expected ~0.34, got %v", encoded)
	}
}

func TestSLog3ReferenceValues(t *testing.T) {
	// Sony S-Log3 reference: 18% gray should encode to approximately 0.41 (41 IRE)
	gray18Linear := 0.18
	encoded := sLog3Transfer(gray18Linear)

	// Allow tolerance for specification variations
	if encoded < 0.38 || encoded > 0.44 {
		t.Errorf("S-Log3 18%% gray: expected ~0.41, got %v", encoded)
	}
}

func TestVLogReferenceValues(t *testing.T) {
	// Panasonic V-Log reference: 18% gray should encode to approximately 0.42
	gray18Linear := 0.18
	encoded := vLogTransfer(gray18Linear)

	// Allow tolerance for specification variations
	if encoded < 0.38 || encoded > 0.46 {
		t.Errorf("V-Log 18%% gray: expected ~0.42, got %v", encoded)
	}
}

// Test HDR values (> 1.0)
func TestLOGSpaceHDRValues(t *testing.T) {
	// LOG spaces should handle HDR values (> 1.0) gracefully
	hdrValues := []float64{2.0, 5.0, 10.0, 100.0}

	for _, hdr := range hdrValues {
		// Test C-Log
		encoded := cLogTransfer(hdr)
		if math.IsNaN(encoded) || math.IsInf(encoded, 0) {
			t.Errorf("C-Log failed on HDR value %v: got %v", hdr, encoded)
		}
		decoded := cLogInverseTransfer(encoded)
		if !floatNear(hdr, decoded, 1e-4) {
			t.Errorf("C-Log HDR round-trip failed: %v -> %v -> %v", hdr, encoded, decoded)
		}

		// Test S-Log3
		encoded = sLog3Transfer(hdr)
		if math.IsNaN(encoded) || math.IsInf(encoded, 0) {
			t.Errorf("S-Log3 failed on HDR value %v: got %v", hdr, encoded)
		}
		decoded = sLog3InverseTransfer(encoded)
		if !floatNear(hdr, decoded, 1e-4) {
			t.Errorf("S-Log3 HDR round-trip failed: %v -> %v -> %v", hdr, encoded, decoded)
		}

		// Test V-Log
		encoded = vLogTransfer(hdr)
		if math.IsNaN(encoded) || math.IsInf(encoded, 0) {
			t.Errorf("V-Log failed on HDR value %v: got %v", hdr, encoded)
		}
		decoded = vLogInverseTransfer(encoded)
		if !floatNear(hdr, decoded, 1e-4) {
			t.Errorf("V-Log HDR round-trip failed: %v -> %v -> %v", hdr, encoded, decoded)
		}
	}
}

// Test alpha channel preservation
func TestLOGSpaceAlpha(t *testing.T) {
	alphaValues := []float64{0.0, 0.5, 1.0}

	for _, alpha := range alphaValues {
		clogColor := NewSpaceColor(CLogSpace, []float64{0.5, 0.4, 0.3}, alpha)

		if clogColor.Alpha() != alpha {
			t.Errorf("Alpha not preserved: expected %v, got %v", alpha, clogColor.Alpha())
		}

		// Convert to sRGB and back
		srgbColor := clogColor.ConvertTo(SRGBSpace)
		clogBack := srgbColor.ConvertTo(CLogSpace)

		if !floatNear(clogBack.Alpha(), alpha, 1e-10) {
			t.Errorf("Alpha not preserved through conversion: expected %v, got %v", alpha, clogBack.Alpha())
		}
	}
}

// Helper function for float comparison with tolerance
func floatNear(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
