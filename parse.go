package color

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParseColor parses a color string in various formats and returns a Color.
//
// Supported formats (all major CSS color formats):
//   - Hex: "#FF0000", "#F00", "#FF000080" (with alpha)
//   - RGB: "rgb(255, 0, 0)", "rgb(100%, 0%, 0%)", "rgb(255 0 0)" (modern syntax)
//   - RGBA: "rgba(255, 0, 0, 0.5)", "rgb(255 0 0 / 0.5)" (modern syntax)
//   - HSL: "hsl(0, 100%, 50%)", "hsl(0 100% 50%)" (modern syntax)
//   - HSLA: "hsla(0, 100%, 50%, 0.5)", "hsl(0 100% 50% / 0.5)" (modern syntax)
//   - HWB: "hwb(0 0% 0%)", "hwb(0 0% 0% / 0.5)" (Hue, Whiteness, Blackness)
//   - HSV: "hsv(0, 100%, 100%)" (not in CSS spec, but commonly used)
//   - HSVA: "hsva(0, 100%, 100%, 0.5)"
//   - LAB: "lab(50% 20 30)" or "lab(50 20 30)" (CIE 1976 L*a*b*)
//   - OKLAB: "oklab(0.6 0.1 -0.1)" (perceptually uniform)
//   - LCH: "lch(70% 50 180)" or "lch(70 50 180)" (CIE LCH from LAB)
//   - OKLCH: "oklch(0.7 0.2 120)" (perceptually uniform)
//   - XYZ: "color(xyz 0.5 0.5 0.5)" or "color(xyz-d65 0.5 0.5 0.5)" (CIE 1931 XYZ)
//   - Wide-gamut RGB via color() function:
//   - "color(srgb 1 0 0)" (sRGB, same as rgb())
//   - "color(srgb-linear 1 0 0)" (linear sRGB, no gamma)
//   - "color(display-p3 1 0 0)" (Display P3, wide gamut)
//   - "color(a98-rgb 1 0 0)" (Adobe RGB 1998)
//   - "color(prophoto-rgb 1 0 0)" (ProPhoto RGB)
//   - "color(rec2020 1 0 0)" (Rec. 2020, UHDTV)
//   - Named colors: "red", "blue", "transparent", etc.
//
// CIE (Commission Internationale de l'Éclairage) color spaces are fully supported:
//   - XYZ: CIE 1931 XYZ color space (via color() function)
//   - LAB: CIE 1976 L*a*b* color space
//   - LCH: Polar representation of CIE LAB
func ParseColor(s string) (Color, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, &ParseError{input: s, reason: "empty string"}
	}

	// Try hex first (most common)
	if strings.HasPrefix(s, "#") || isHexString(s) {
		return HexToRGB(s)
	}

	// Try named colors
	if c, ok := parseNamedColor(s); ok {
		return c, nil
	}

	// Try function formats (rgb, rgba, hsl, etc.)
	if strings.Contains(s, "(") {
		return parseFunctionColor(s)
	}

	return nil, &ParseError{input: s, reason: "unknown format"}
}

// ParseError represents an error parsing a color string.
type ParseError struct {
	input  string
	reason string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("cannot parse color %q: %s", e.input, e.reason)
}

// isHexString checks if a string looks like a hex color (without #).
func isHexString(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Remove # if present
	if s[0] == '#' {
		s = s[1:]
	}
	// Check if all characters are hex digits
	matched, _ := regexp.MatchString(`^[0-9A-Fa-f]{3,8}$`, s)
	return matched && (len(s) == 3 || len(s) == 4 || len(s) == 6 || len(s) == 8)
}

// parseFunctionColor parses CSS function-style color strings.
func parseFunctionColor(s string) (Color, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// Extract function name and arguments
	// Handle both simple functions (rgb(...)) and color() function (color(xyz ...))
	re := regexp.MustCompile(`^(\w+)(?:-[\w]+)?\(([^)]+)\)$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		return nil, &ParseError{input: s, reason: "invalid function format"}
	}

	funcName := strings.ToLower(matches[1])
	args := matches[2]

	// Parse arguments (split by comma, handle spaces)
	argList := parseArgs(args)

	// Handle color() function specially - it has format color(xyz 0.5 0.5 0.5)
	// The color space name is the first argument
	if funcName == "color" {
		return parseColorFunction(argList)
	}

	switch funcName {
	case "rgb":
		return parseRGB(argList, false)
	case "rgba":
		return parseRGB(argList, true)
	case "hsl":
		return parseHSL(argList, false)
	case "hsla":
		return parseHSL(argList, true)
	case "hwb":
		return parseHWB(argList)
	case "hsv":
		return parseHSV(argList, false)
	case "hsva":
		return parseHSV(argList, true)
	case "lab":
		return parseLAB(argList)
	case "oklab":
		return parseOKLAB(argList)
	case "lch":
		return parseLCH(argList)
	case "oklch":
		return parseOKLCH(argList)
	case "color":
		return parseColorFunction(argList)
	case "xyz":
		return parseXYZ(argList)
	default:
		return nil, &ParseError{input: s, reason: fmt.Sprintf("unknown function: %s", funcName)}
	}
}

// parseArgs splits function arguments, handling commas, spaces, and slashes.
// Modern CSS syntax uses spaces and "/" for alpha: "rgb(255 0 0 / 0.5)"
// For LAB/OKLAB/LCH/OKLCH, spaces are used. For others, commas or spaces.
func parseArgs(s string) []string {
	// Remove extra whitespace
	s = strings.TrimSpace(s)

	// Handle alpha with slash (modern syntax): "255 0 0 / 0.5"
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) == 2 {
			// Split main args and alpha
			mainArgs := parseArgs(parts[0])
			alpha := strings.TrimSpace(parts[1])
			return append(mainArgs, alpha)
		}
	}

	// Check if it contains commas (legacy RGB, HSL, HSV use commas)
	if strings.Contains(s, ",") {
		// Split by comma, then trim each part
		parts := strings.Split(s, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		return result
	}

	// For space-separated (LAB, OKLAB, LCH, OKLCH, modern RGB/HSL), split by whitespace
	parts := strings.Fields(s)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && part != "/" {
			result = append(result, part)
		}
	}
	return result
}

// parseNumber parses a number that can be:
// - An integer: "255"
// - A float: "0.5"
// - A percentage: "50%" (returns 0.5)
// - A float with unit: "50deg" (returns 50, unit ignored for now)
func parseNumber(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// Remove common units (deg, etc.)
	s = strings.TrimSuffix(s, "deg")
	s = strings.TrimSpace(s)

	// Check for percentage
	if strings.HasSuffix(s, "%") {
		val, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return 0, err
		}
		return val / 100.0, nil
	}

	// Regular number
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// parseRGB parses RGB/RGBA arguments.
// Supports both legacy (comma-separated) and modern (space-separated) syntax.
func parseRGB(args []string, hasAlpha bool) (Color, error) {
	// Check for too many arguments - max 4 (rgb + alpha)
	if len(args) > 4 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "RGB/RGBA requires at most 4 arguments"}
	}

	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "RGB requires at least 3 arguments"}
	}

	r, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}
	// Validate range before conversion
	if r < 0 {
		return nil, &ParseError{input: args[0], reason: "RGB red component cannot be negative"}
	}
	// If > 1, assume 0-255 range, validate and convert to 0-1
	if r > 1 {
		if r > 255 {
			return nil, &ParseError{input: args[0], reason: "RGB red component out of range (0-255)"}
		}
		r = r / 255.0
	}

	g, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}
	if g < 0 {
		return nil, &ParseError{input: args[1], reason: "RGB green component cannot be negative"}
	}
	if g > 1 {
		if g > 255 {
			return nil, &ParseError{input: args[1], reason: "RGB green component out of range (0-255)"}
		}
		g = g / 255.0
	}

	b, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}
	if b < 0 {
		return nil, &ParseError{input: args[2], reason: "RGB blue component cannot be negative"}
	}
	if b > 1 {
		if b > 255 {
			return nil, &ParseError{input: args[2], reason: "RGB blue component out of range (0-255)"}
		}
		b = b / 255.0
	}

	var a float64 = 1.0
	// Check for alpha (either explicit hasAlpha flag or 4th argument)
	if hasAlpha || len(args) >= 4 {
		if len(args) < 4 {
			return nil, &ParseError{input: strings.Join(args, ","), reason: "RGBA requires 4 arguments"}
		}
		a, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
		// Alpha is always 0-1, not 0-255
	}

	return NewRGBA(r, g, b, a), nil
}

// parseHSL parses HSL/HSLA arguments.
// Supports both legacy (comma-separated) and modern (space-separated) syntax.
func parseHSL(args []string, hasAlpha bool) (Color, error) {
	// Check for too many arguments - max 4 (hsl + alpha)
	if len(args) > 4 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "HSL/HSLA requires at most 4 arguments"}
	}

	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "HSL requires at least 3 arguments"}
	}

	h, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}
	// Hue is in degrees, validate range 0-360
	if h < 0 || h > 360 {
		return nil, &ParseError{input: args[0], reason: "HSL hue out of range (0-360 degrees)"}
	}

	s, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}
	// Saturation is 0-1 (or 0-100% which parseNumber handles)
	if s < 0 || s > 1 {
		return nil, &ParseError{input: args[1], reason: "HSL saturation out of range (0-100%)"}
	}

	l, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}
	// Lightness is 0-1 (or 0-100% which parseNumber handles)
	if l < 0 || l > 1 {
		return nil, &ParseError{input: args[2], reason: "HSL lightness out of range (0-100%)"}
	}

	var a float64 = 1.0
	// Check for alpha (either explicit hasAlpha flag or 4th argument)
	if hasAlpha || len(args) >= 4 {
		if len(args) < 4 {
			return nil, &ParseError{input: strings.Join(args, ","), reason: "HSLA requires 4 arguments"}
		}
		a, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	return NewHSL(h, s, l, a), nil
}

// parseHSV parses HSV/HSVA arguments.
func parseHSV(args []string, hasAlpha bool) (Color, error) {
	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "HSV requires at least 3 arguments"}
	}

	h, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}

	s, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}

	v, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}

	var a float64 = 1.0
	if hasAlpha {
		if len(args) < 4 {
			return nil, &ParseError{input: strings.Join(args, ","), reason: "HSVA requires 4 arguments"}
		}
		a, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	return NewHSV(h, s, v, a), nil
}

// parseLAB parses LAB arguments.
func parseLAB(args []string) (Color, error) {
	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "LAB requires 3 arguments"}
	}

	l, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}
	// L can be 0-100 or 0-1 (percentage)
	if l <= 1 {
		l = l * 100 // Convert to 0-100 range
	}

	a, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}

	b, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}

	alpha := 1.0
	if len(args) >= 4 {
		alpha, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	return NewLAB(l, a, b, alpha), nil
}

// parseOKLAB parses OKLAB arguments.
func parseOKLAB(args []string) (Color, error) {
	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "OKLAB requires 3 arguments"}
	}

	l, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}
	// OKLAB L is 0-1

	a, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}

	b, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}

	alpha := 1.0
	if len(args) >= 4 {
		alpha, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	return NewOKLAB(l, a, b, alpha), nil
}

// parseLCH parses LCH arguments.
func parseLCH(args []string) (Color, error) {
	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "LCH requires 3 arguments"}
	}

	l, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}
	// L can be 0-100 or 0-1 (percentage)
	if l <= 1 {
		l = l * 100 // Convert to 0-100 range
	}

	c, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}

	h, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}

	alpha := 1.0
	if len(args) >= 4 {
		alpha, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	return NewLCH(l, c, h, alpha), nil
}

// parseOKLCH parses OKLCH arguments.
func parseOKLCH(args []string) (Color, error) {
	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, ","), reason: "OKLCH requires 3 arguments"}
	}

	l, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}
	// OKLCH L is 0-1
	if l < 0 || l > 1 {
		return nil, &ParseError{input: args[0], reason: "OKLCH lightness out of range (0-1)"}
	}

	c, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}
	// Chroma must be non-negative
	if c < 0 {
		return nil, &ParseError{input: args[1], reason: "OKLCH chroma cannot be negative"}
	}

	h, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}

	alpha := 1.0
	if len(args) >= 4 {
		alpha, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	return NewOKLCH(l, c, h, alpha), nil
}

// parseHWB parses HWB (Hue, Whiteness, Blackness) arguments.
func parseHWB(args []string) (Color, error) {
	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, " "), reason: "HWB requires at least 3 arguments"}
	}

	h, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}
	// Hue is in degrees

	w, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}
	// Whiteness is 0-1 (or 0-100% which parseNumber handles)

	b, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}
	// Blackness is 0-1 (or 0-100% which parseNumber handles)

	alpha := 1.0
	if len(args) >= 4 {
		alpha, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	// Convert HWB to RGB
	// HWB is similar to HSV but uses whiteness and blackness
	// Algorithm: w + b cannot exceed 1, if it does, scale both down
	sum := w + b
	if sum > 1 {
		w = w / sum
		b = b / sum
	}

	// Convert to RGB via HSV-like calculation
	// This is a simplified conversion
	v := 1 - b
	s := 0.0
	if v > 0 {
		s = 1 - (w / v)
	}

	// Create HSV and convert to RGB
	hsv := NewHSV(h, s, v, alpha)
	return hsv, nil
}

// parseXYZ parses XYZ color space arguments.
func parseXYZ(args []string) (Color, error) {
	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, " "), reason: "XYZ requires 3 arguments"}
	}

	x, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}

	y, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}

	z, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}

	alpha := 1.0
	if len(args) >= 4 {
		alpha, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	return NewXYZ(x, y, z, alpha), nil
}

// parseColorFunction parses the CSS color() function.
// Format: color(xyz 0.5 0.5 0.5) or color(display-p3 1 0 0)
//
// Supported color spaces in color() function:
//   - xyz, xyz-d50, xyz-d65 (CIE 1931 XYZ)
//   - srgb (same as rgb())
//   - srgb-linear (linear sRGB, no gamma encoding)
//   - display-p3 (wide gamut RGB, P3 display)
//   - a98-rgb (Adobe RGB 1998)
//   - prophoto-rgb (ProPhoto RGB)
//   - rec2020 (Rec. 2020, UHDTV)
func parseColorFunction(args []string) (Color, error) {
	if len(args) < 4 {
		return nil, &ParseError{input: strings.Join(args, " "), reason: "color() function requires color space name and 3 values"}
	}

	colorSpace := strings.ToLower(strings.TrimSpace(args[0]))

	// Handle CIE XYZ color spaces (xyz, xyz-d50, xyz-d65)
	// Note: Currently all XYZ variants use D65 white point
	if strings.HasPrefix(colorSpace, "xyz") {
		// XYZ values (x, y, z) - CIE 1931 XYZ color space
		return parseXYZ(args[1:])
	}

	// Handle RGB color spaces (srgb, display-p3, a98-rgb, etc.)
	if rgbSpace := getRGBColorSpace(colorSpace); rgbSpace != nil {
		return parseRGBColorSpace(args[1:], rgbSpace)
	}

	return nil, &ParseError{input: strings.Join(args, " "), reason: fmt.Sprintf("unsupported color space in color() function: %s (supported: xyz, xyz-d50, xyz-d65, srgb, srgb-linear, display-p3, a98-rgb, prophoto-rgb, rec2020)", colorSpace)}
}

// parseRGBColorSpace parses RGB values for a specific RGB color space.
func parseRGBColorSpace(args []string, space *RGBColorSpace) (Color, error) {
	if len(args) < 3 {
		return nil, &ParseError{input: strings.Join(args, " "), reason: fmt.Sprintf("%s requires 3 arguments", space.Name)}
	}

	r, err := parseNumber(args[0])
	if err != nil {
		return nil, err
	}
	// If > 1, assume 0-255 range, convert to 0-1
	if r > 1 {
		r = r / 255.0
	}

	g, err := parseNumber(args[1])
	if err != nil {
		return nil, err
	}
	if g > 1 {
		g = g / 255.0
	}

	b, err := parseNumber(args[2])
	if err != nil {
		return nil, err
	}
	if b > 1 {
		b = b / 255.0
	}

	alpha := 1.0
	if len(args) >= 4 {
		alpha, err = parseNumber(args[3])
		if err != nil {
			return nil, err
		}
	}

	// For wide-gamut color spaces, we need to convert to XYZ then to sRGB
	// Convert the RGB values in the source color space to XYZ
	xyz := space.ConvertRGBToXYZ(r, g, b, alpha)

	// Convert XYZ to sRGB for the Color interface
	// (All colors must be convertible to RGBA via the Color interface)
	// We return an RGBA color that represents the closest sRGB approximation
	r2, g2, b2, a2 := xyz.RGBA()
	return NewRGBA(r2, g2, b2, a2), nil
}

// parseNamedColor parses CSS named colors.
func parseNamedColor(s string) (Color, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	if color, ok := namedColors[s]; ok {
		return color, true
	}
	return nil, false
}

// namedColors contains CSS named colors.
var namedColors = map[string]Color{
	"transparent": NewRGBA(0, 0, 0, 0),
	"black":       RGB(0, 0, 0),
	"white":       RGB(1, 1, 1),
	"red":         RGB(1, 0, 0),
	"green":       RGB(0, 0.5, 0), // CSS green is #008000
	"blue":        RGB(0, 0, 1),
	"yellow":      RGB(1, 1, 0),
	"cyan":        RGB(0, 1, 1),
	"magenta":     RGB(1, 0, 1),
	"orange":      RGB(1, 0.647, 0), // #FFA500
	"purple":      RGB(0.5, 0, 0.5),
	"pink":        RGB(1, 0.753, 0.796),     // #FFC0CB
	"brown":       RGB(0.647, 0.165, 0.165), // #A52A2A
	"gray":        RGB(0.5, 0.5, 0.5),
	"grey":        RGB(0.5, 0.5, 0.5),
	"lime":        RGB(0, 1, 0),
	"navy":        RGB(0, 0, 0.5),
	"olive":       RGB(0.5, 0.5, 0),
	"teal":        RGB(0, 0.5, 0.5),
	"aqua":        RGB(0, 1, 1),
	"fuchsia":     RGB(1, 0, 1),
	"maroon":      RGB(0.5, 0, 0),
	"silver":      RGB(0.753, 0.753, 0.753),
	"darkred":     RGB(0.545, 0, 0),
	"darkgreen":   RGB(0, 0.392, 0),
	"darkblue":    RGB(0, 0, 0.545),
	"darkgray":    RGB(0.663, 0.663, 0.663),
	"darkgrey":    RGB(0.663, 0.663, 0.663),
	"lightgray":   RGB(0.827, 0.827, 0.827),
	"lightgrey":   RGB(0.827, 0.827, 0.827),
}
