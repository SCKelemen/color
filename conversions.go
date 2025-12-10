package color

// ConvertToRGBSpace converts a color to a specific RGB color space.
// This allows converting colors to wide-gamut RGB spaces like display-p3, a98-rgb, etc.
//
// Example:
//   rgb := color.RGB(1, 0, 0)
//   displayP3 := color.ConvertToRGBSpace(rgb, "display-p3")
func ConvertToRGBSpace(c Color, spaceName string) (*RGBA, error) {
	space := getRGBColorSpace(spaceName)
	if space == nil {
		return nil, &ParseError{input: spaceName, reason: "unknown RGB color space"}
	}

	// Convert color to XYZ first
	xyz := ToXYZ(c)

	// Convert XYZ to the target RGB color space
	return space.ConvertXYZToRGB(xyz), nil
}

// ConvertFromRGBSpace converts a color from a specific RGB color space to the standard Color interface.
// This is useful when you have RGB values in a wide-gamut space and want to work with them.
//
// Example:
//   // You have display-p3 RGB values
//   displayP3RGB := color.ConvertFromRGBSpace(1.0, 0.0, 0.0, 1.0, "display-p3")
//   // Now convert to OKLCH for manipulation
//   oklch := color.ToOKLCH(displayP3RGB)
func ConvertFromRGBSpace(r, g, b, a float64, spaceName string) (Color, error) {
	space := getRGBColorSpace(spaceName)
	if space == nil {
		return nil, &ParseError{input: spaceName, reason: "unknown RGB color space"}
	}

	// Convert RGB in the source space to XYZ
	xyz := space.ConvertRGBToXYZ(r, g, b, a)

	// Convert XYZ to sRGB (standard Color interface)
	r2, g2, b2, a2 := xyz.RGBA()
	return NewRGBA(r2, g2, b2, a2), nil
}

