package color

import "fmt"

// HexToRGB converts a hex color string (with or without #) to RGB.
// Supports formats: #RGB, #RRGGBB, #RGBA, #RRGGBBAA
func HexToRGB(hex string) (*RGBA, error) {
	// Remove # if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	var r, g, b, a float64
	var err error

	switch len(hex) {
	case 3: // #RGB
		r, err = parseHexByte(hex[0], hex[0])
		if err != nil {
			return nil, err
		}
		g, err = parseHexByte(hex[1], hex[1])
		if err != nil {
			return nil, err
		}
		b, err = parseHexByte(hex[2], hex[2])
		if err != nil {
			return nil, err
		}
		a = 1.0
	case 4: // #RGBA
		r, err = parseHexByte(hex[0], hex[0])
		if err != nil {
			return nil, err
		}
		g, err = parseHexByte(hex[1], hex[1])
		if err != nil {
			return nil, err
		}
		b, err = parseHexByte(hex[2], hex[2])
		if err != nil {
			return nil, err
		}
		a, err = parseHexByte(hex[3], hex[3])
		if err != nil {
			return nil, err
		}
	case 6: // #RRGGBB
		r, err = parseHexByte(hex[0], hex[1])
		if err != nil {
			return nil, err
		}
		g, err = parseHexByte(hex[2], hex[3])
		if err != nil {
			return nil, err
		}
		b, err = parseHexByte(hex[4], hex[5])
		if err != nil {
			return nil, err
		}
		a = 1.0
	case 8: // #RRGGBBAA
		r, err = parseHexByte(hex[0], hex[1])
		if err != nil {
			return nil, err
		}
		g, err = parseHexByte(hex[2], hex[3])
		if err != nil {
			return nil, err
		}
		b, err = parseHexByte(hex[4], hex[5])
		if err != nil {
			return nil, err
		}
		a, err = parseHexByte(hex[6], hex[7])
		if err != nil {
			return nil, err
		}
	default:
		return nil, &HexParseError{hex: hex}
	}

	// For cases 3 and 6, alpha is already 1.0 (not divided by 255)
	// For cases 4 and 8, alpha needs to be divided by 255
	if len(hex) == 3 || len(hex) == 6 {
		return NewRGBA(r/255.0, g/255.0, b/255.0, 1.0), nil
	}
	return NewRGBA(r/255.0, g/255.0, b/255.0, a/255.0), nil
}

// parseHexByte parses a single hex byte (one or two hex digits).
func parseHexByte(high, low byte) (float64, error) {
	highVal, err := hexDigitToInt(high)
	if err != nil {
		return 0, err
	}
	lowVal, err := hexDigitToInt(low)
	if err != nil {
		return 0, err
	}
	return float64(highVal*16 + lowVal), nil
}

// hexDigitToInt converts a hex digit to its integer value.
func hexDigitToInt(b byte) (int, error) {
	switch {
	case b >= '0' && b <= '9':
		return int(b - '0'), nil
	case b >= 'a' && b <= 'f':
		return int(b - 'a' + 10), nil
	case b >= 'A' && b <= 'F':
		return int(b - 'A' + 10), nil
	default:
		return 0, &HexParseError{hex: string(b)}
	}
}

// HexParseError represents an error parsing a hex color string.
type HexParseError struct {
	hex string
}

func (e *HexParseError) Error() string {
	return "invalid hex color: " + e.hex
}

// RGBToHex converts an RGB color to a hex string.
// Returns format #RRGGBB or #RRGGBBAA if alpha < 1.0.
func RGBToHex(c Color) string {
	r, g, b, a := c.RGBA()
	
	r8 := uint8(r * 255)
	g8 := uint8(g * 255)
	b8 := uint8(b * 255)
	
	if a >= 1.0 {
		return formatHex(r8, g8, b8)
	}
	
	a8 := uint8(a * 255)
	return formatHexWithAlpha(r8, g8, b8, a8)
}

// formatHex formats RGB bytes as #RRGGBB.
func formatHex(r, g, b uint8) string {
	return "#" + formatHexByte(r) + formatHexByte(g) + formatHexByte(b)
}

// formatHexWithAlpha formats RGBA bytes as #RRGGBBAA.
func formatHexWithAlpha(r, g, b, a uint8) string {
	return "#" + formatHexByte(r) + formatHexByte(g) + formatHexByte(b) + formatHexByte(a)
}

// formatHexByte formats a byte as a two-digit hex string.
func formatHexByte(b uint8) string {
	const hex = "0123456789abcdef"
	return string([]byte{hex[b>>4], hex[b&0x0f]})
}

// formatString formats a color string.
func formatString(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

