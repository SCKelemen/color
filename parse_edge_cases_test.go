package color

import (
	"testing"
)

// Test edge cases for color parsing

func TestParseColorMalformedHex(t *testing.T) {
	malformed := []string{
		"#GG0000",        // Invalid hex digit
		"#12345",         // Wrong length
		"#1234567",       // Wrong length
		"##FF0000",       // Double hash
		"FF00GG",         // Invalid digit without hash
		"",               // Empty string
		"#",              // Just hash
		"notacolor",      // Random string
		"#ZZZZZZ",        // All invalid digits
	}

	for _, input := range malformed {
		_, err := ParseColor(input)
		if err == nil {
			t.Errorf("Expected error for malformed hex: %q", input)
		}
	}
}

func TestParseColorValidHexFormats(t *testing.T) {
	valid := map[string][3]float64{
		"#F00":     {1.0, 0, 0},       // Short form red
		"#FF0000":  {1.0, 0, 0},       // Long form red
		"#00FF00":  {0, 1.0, 0},       // Green
		"#0000FF":  {0, 0, 1.0},       // Blue
		"#FFF":     {1.0, 1.0, 1.0},   // White short
		"#FFFFFF":  {1.0, 1.0, 1.0},   // White long
		"#000":     {0, 0, 0},         // Black short
		"#000000":  {0, 0, 0},         // Black long
		"FF0000":   {1.0, 0, 0},       // Without hash
		"F00":      {1.0, 0, 0},       // Short without hash
	}

	for input, expected := range valid {
		color, err := ParseColor(input)
		if err != nil {
			t.Errorf("Failed to parse valid hex %q: %v", input, err)
			continue
		}

		r, g, b, _ := color.RGBA()
		if abs(r-expected[0]) > 0.01 || abs(g-expected[1]) > 0.01 || abs(b-expected[2]) > 0.01 {
			t.Errorf("Hex %q: got (%f, %f, %f), want (%f, %f, %f)",
				input, r, g, b, expected[0], expected[1], expected[2])
		}
	}
}

func TestParseColorRGBMalformed(t *testing.T) {
	malformed := []string{
		"rgb()",                    // Empty
		"rgb(255)",                 // Missing components
		"rgb(255, 128)",            // Missing component
		"rgb(255, 128, 64, 128)",   // Too many components (without alpha)
		"rgb(256, 128, 64)",        // Out of range
		"rgb(-1, 128, 64)",         // Negative
		"rgb(a, b, c)",             // Non-numeric
		"rgb(255 128 64)",          // Missing commas
		"rgb(255, , 64)",           // Empty component
	}

	for _, input := range malformed {
		_, err := ParseColor(input)
		if err == nil {
			t.Errorf("Expected error for malformed RGB: %q", input)
		}
	}
}

func TestParseColorHSLMalformed(t *testing.T) {
	malformed := []string{
		"hsl()",                      // Empty
		"hsl(180)",                   // Missing components
		"hsl(180, 50%)",              // Missing component
		"hsl(180, 50%, 50%, 50%)",    // Too many components
		"hsl(400, 50%, 50%)",         // Hue out of range
		"hsl(180, 150%, 50%)",        // Saturation out of range
		"hsl(180, 50%, 150%)",        // Lightness out of range
		"hsl(abc, 50%, 50%)",         // Non-numeric hue
	}

	for _, input := range malformed {
		_, err := ParseColor(input)
		if err == nil {
			t.Errorf("Expected error for malformed HSL: %q", input)
		}
	}
}

func TestParseColorOKLCHMalformed(t *testing.T) {
	malformed := []string{
		"oklch()",                    // Empty
		"oklch(0.5)",                 // Missing components
		"oklch(0.5, 0.2)",            // Missing component
		"oklch(2.0 0.2 180)",         // Lightness out of range
		"oklch(0.5 -0.1 180)",        // Negative chroma
		"oklch(a b c)",               // Non-numeric
	}

	for _, input := range malformed {
		_, err := ParseColor(input)
		if err == nil {
			t.Errorf("Expected error for malformed OKLCH: %q", input)
		}
	}
}

func TestParseColorCaseSensitivity(t *testing.T) {
	// These should all parse successfully
	variants := []string{
		"rgb(255, 0, 0)",
		"RGB(255, 0, 0)",
		"Rgb(255, 0, 0)",
		"hsl(0, 100%, 50%)",
		"HSL(0, 100%, 50%)",
		"Hsl(0, 100%, 50%)",
		"#FF0000",
		"#ff0000",
		"#Ff0000",
	}

	for _, input := range variants {
		_, err := ParseColor(input)
		if err != nil {
			t.Errorf("Failed to parse case variant %q: %v", input, err)
		}
	}
}

func TestParseColorWhitespace(t *testing.T) {
	// Test various whitespace scenarios
	withWhitespace := []string{
		" rgb(255, 0, 0) ",           // Leading/trailing spaces
		"rgb( 255 , 0 , 0 )",         // Extra spaces
		"rgb(255,0,0)",                // No spaces
		"  #FF0000  ",                 // Spaces around hex
		"hsl(180 , 50% , 50%)",       // Mixed spacing
	}

	for _, input := range withWhitespace {
		color, err := ParseColor(input)
		if err != nil {
			t.Errorf("Failed to parse with whitespace %q: %v", input, err)
			continue
		}

		// Just verify it parsed successfully
		_, _, _, _ = color.RGBA()
	}
}

func TestParseColorRGBAAlpha(t *testing.T) {
	tests := []struct {
		input         string
		expectedAlpha float64
	}{
		{"rgba(255, 0, 0, 1.0)", 1.0},
		{"rgba(255, 0, 0, 0.5)", 0.5},
		{"rgba(255, 0, 0, 0.0)", 0.0},
		{"rgba(255, 0, 0, 0)", 0.0},
		{"rgba(255, 0, 0, 1)", 1.0},
	}

	for _, tt := range tests {
		color, err := ParseColor(tt.input)
		if err != nil {
			t.Errorf("Failed to parse %q: %v", tt.input, err)
			continue
		}

		alpha := color.Alpha()
		if abs(alpha-tt.expectedAlpha) > 0.01 {
			t.Errorf("Input %q: got alpha %f, want %f", tt.input, alpha, tt.expectedAlpha)
		}
	}
}

func TestParseColorHexAlpha(t *testing.T) {
	tests := []struct {
		input         string
		expectedAlpha float64
	}{
		{"#FF0000FF", 1.0},     // Full alpha
		{"#FF000080", 0.5},     // Half alpha (128/255)
		{"#FF000000", 0.0},     // Zero alpha
		{"#F00F", 1.0},         // Short form full alpha
		{"#F008", 0.53},        // Short form half alpha (8/15 ≈ 0.53)
		{"#F000", 0.0},         // Short form zero alpha
	}

	for _, tt := range tests {
		color, err := ParseColor(tt.input)
		if err != nil {
			t.Errorf("Failed to parse %q: %v", tt.input, err)
			continue
		}

		alpha := color.Alpha()
		if abs(alpha-tt.expectedAlpha) > 0.02 {
			t.Errorf("Input %q: got alpha %f, want %f", tt.input, alpha, tt.expectedAlpha)
		}
	}
}

func TestHexToRGBRoundTrip(t *testing.T) {
	colors := []struct {
		hex string
	}{
		{"#FF0000"},
		{"#00FF00"},
		{"#0000FF"},
		{"#FFFFFF"},
		{"#000000"},
		{"#FF8800"},
		{"#8800FF"},
		{"#F08"},  // Short form
		{"#ABC"},  // Short form
	}

	for _, tc := range colors {
		// Parse
		color, err := HexToRGB(tc.hex)
		if err != nil {
			t.Errorf("Failed to parse %q: %v", tc.hex, err)
			continue
		}

		// Convert back
		hex := RGBToHex(color)

		// Parse again
		color2, err := HexToRGB(hex)
		if err != nil {
			t.Errorf("Failed to parse generated hex %q: %v", hex, err)
			continue
		}

		// Compare
		r1, g1, b1, a1 := color.RGBA()
		r2, g2, b2, a2 := color2.RGBA()

		if abs(r1-r2) > 0.01 || abs(g1-g2) > 0.01 || abs(b1-b2) > 0.01 || abs(a1-a2) > 0.01 {
			t.Errorf("Round trip failed for %q: (%f,%f,%f,%f) -> %q -> (%f,%f,%f,%f)",
				tc.hex, r1, g1, b1, a1, hex, r2, g2, b2, a2)
		}
	}
}

func TestParseColorUnknownFormats(t *testing.T) {
	unknown := []string{
		"cmyk(0, 100, 100, 0)",       // Unsupported format
		"device-cmyk(0, 1, 1, 0)",    // Unsupported format
		"random string",               // Random text
		"12345",                       // Just numbers
		"color",                       // Keyword without value
	}

	for _, input := range unknown {
		_, err := ParseColor(input)
		if err == nil {
			t.Errorf("Expected error for unknown format: %q", input)
		}
	}
}

func TestParseColorBoundaryValues(t *testing.T) {
	boundary := []string{
		"rgb(0, 0, 0)",              // Min
		"rgb(255, 255, 255)",        // Max
		"hsl(0, 0%, 0%)",            // HSL min
		"hsl(359, 100%, 100%)",      // HSL max
		"oklch(0 0 0)",              // OKLCH min
		"oklch(1 0.4 359)",          // OKLCH near max
	}

	for _, input := range boundary {
		color, err := ParseColor(input)
		if err != nil {
			t.Errorf("Failed to parse boundary value %q: %v", input, err)
			continue
		}

		// Verify RGB is in valid range
		r, g, b, a := color.RGBA()
		if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
			t.Errorf("Boundary value %q produced invalid RGB: (%f, %f, %f, %f)", input, r, g, b, a)
		}
	}
}

func TestParseColorNamedColors(t *testing.T) {
	// If named colors are supported
	named := map[string][3]float64{
		"red":   {1.0, 0, 0},
		"green": {0, 0.5, 0},   // CSS green is not pure green
		"blue":  {0, 0, 1.0},
		"white": {1.0, 1.0, 1.0},
		"black": {0, 0, 0},
	}

	for name, expected := range named {
		color, err := ParseColor(name)
		if err != nil {
			// Named colors might not be supported, skip
			t.Logf("Named color %q not supported (this is okay)", name)
			continue
		}

		r, g, b, _ := color.RGBA()
		if abs(r-expected[0]) > 0.1 || abs(g-expected[1]) > 0.1 || abs(b-expected[2]) > 0.1 {
			t.Errorf("Named color %q: got (%f, %f, %f), want (%f, %f, %f)",
				name, r, g, b, expected[0], expected[1], expected[2])
		}
	}
}
