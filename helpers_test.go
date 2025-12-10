package color

import (
	"testing"
)

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected [4]float64 // RGBA
		wantErr  bool
	}{
		{"Red #FF0000", "#FF0000", [4]float64{1, 0, 0, 1}, false},
		{"Green #00FF00", "#00FF00", [4]float64{0, 1, 0, 1}, false},
		{"Blue #0000FF", "#0000FF", [4]float64{0, 0, 1, 1}, false},
		{"White #FFFFFF", "#FFFFFF", [4]float64{1, 1, 1, 1}, false},
		{"Black #000000", "#000000", [4]float64{0, 0, 0, 1}, false},
		{"Cyan #00FFFF", "#00FFFF", [4]float64{0, 1, 1, 1}, false},
		{"Magenta #FF00FF", "#FF00FF", [4]float64{1, 0, 1, 1}, false},
		{"Yellow #FFFF00", "#FFFF00", [4]float64{1, 1, 0, 1}, false},
		{"Red short #F00", "#F00", [4]float64{1, 0, 0, 1}, false},
		{"Green short #0F0", "#0F0", [4]float64{0, 1, 0, 1}, false},
		{"Blue short #00F", "#00F", [4]float64{0, 0, 1, 1}, false},
		{"With alpha #FF000080", "#FF000080", [4]float64{1, 0, 0, 0.5019607843137255}, false}, // 128/255
		{"Short with alpha #F008", "#F008", [4]float64{1, 0, 0, 0.5333333333333333}, false}, // 136/255
		{"No hash F00", "F00", [4]float64{1, 0, 0, 1}, false},
		{"Invalid hex #GGG", "#GGG", [4]float64{}, true},
		{"Invalid length #FF", "#FF", [4]float64{}, true},
		{"Empty", "", [4]float64{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgb, err := HexToRGB(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToRGB(%q) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				r, g, b, a := rgb.RGBA()
				if !rgbaEqual(r, g, b, a, tt.expected[0], tt.expected[1], tt.expected[2], tt.expected[3]) {
					t.Errorf("HexToRGB(%q) = RGB(%v, %v, %v, %v), want RGB(%v, %v, %v, %v)",
						tt.hex, r, g, b, a, tt.expected[0], tt.expected[1], tt.expected[2], tt.expected[3])
				}
			}
		})
	}
}

func TestRGBToHex(t *testing.T) {
	tests := []struct {
		name     string
		rgb      *RGBA
		expected string
	}{
		{"Red", RGB(1, 0, 0), "#ff0000"},
		{"Green", RGB(0, 1, 0), "#00ff00"},
		{"Blue", RGB(0, 0, 1), "#0000ff"},
		{"White", RGB(1, 1, 1), "#ffffff"},
		{"Black", RGB(0, 0, 0), "#000000"},
		{"Cyan", RGB(0, 1, 1), "#00ffff"},
		{"Magenta", RGB(1, 0, 1), "#ff00ff"},
		{"Yellow", RGB(1, 1, 0), "#ffff00"},
		{"Gray", RGB(0.5, 0.5, 0.5), "#7f7f7f"}, // 0.5 * 255 = 127.5 -> 127 = 0x7f
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hex := RGBToHex(tt.rgb)
			if hex != tt.expected {
				t.Errorf("RGBToHex(%v) = %q, want %q", tt.rgb, hex, tt.expected)
			}
		})
	}
}

func TestHexRoundTrip(t *testing.T) {
	hexStrings := []string{
		"#FF0000",
		"#00FF00",
		"#0000FF",
		"#FFFFFF",
		"#000000",
		"#FF5733",
		"#808080",
	}

	for _, hex := range hexStrings {
		t.Run(hex, func(t *testing.T) {
			rgb, err := HexToRGB(hex)
			if err != nil {
				t.Fatalf("HexToRGB(%q) error = %v", hex, err)
			}

			hex2 := RGBToHex(rgb)
			// Convert back to compare
			rgb2, err := HexToRGB(hex2)
			if err != nil {
				t.Fatalf("HexToRGB(%q) error = %v", hex2, err)
			}

			r1, g1, b1, a1 := rgb.RGBA()
			r2, g2, b2, a2 := rgb2.RGBA()

			// Allow some tolerance for rounding
			if !rgbaEqual(r1, g1, b1, a1, r2, g2, b2, a2) {
				t.Errorf("Hex round-trip failed:\n"+
					"  Original: %q -> RGB(%v, %v, %v, %v)\n"+
					"  Round-trip: %q -> RGB(%v, %v, %v, %v)",
					hex, r1, g1, b1, a1, hex2, r2, g2, b2, a2)
			}
		})
	}
}

