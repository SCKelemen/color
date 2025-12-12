package color

// ConvertToRGBSpace converts a color to a specific RGB color space using the new Space system.
// This allows converting colors to wide-gamut RGB spaces like display-p3, a98-rgb, etc.
// Returns a SpaceColor that preserves the target color space information.
//
// Example:
//   rgb := color.RGB(1, 0, 0)
//   displayP3, _ := color.ConvertToRGBSpace(rgb, "display-p3")
//   // displayP3 is now a SpaceColor in Display P3 space
func ConvertToRGBSpace(c Color, spaceName string) (SpaceColor, error) {
	space := getSpaceByName(spaceName)
	if space == nil {
		return nil, &ParseError{input: spaceName, reason: "unknown RGB color space"}
	}

	// Convert color to XYZ first
	xyz := ToXYZ(c)

	// Convert XYZ to the target RGB color space
	channels := space.FromXYZ(xyz.X, xyz.Y, xyz.Z)
	return NewSpaceColor(space, channels, c.Alpha()), nil
}

// ConvertFromRGBSpace converts RGB values from a specific color space to an RGBA color.
// This is useful when you have RGB values in a wide-gamut space and want to work with them.
//
// Example:
//   // You have display-p3 RGB values
//   displayP3RGB, _ := color.ConvertFromRGBSpace(1.0, 0.0, 0.0, 1.0, "display-p3")
//   // Now convert to OKLCH for manipulation
//   oklch := color.ToOKLCH(displayP3RGB)
func ConvertFromRGBSpace(r, g, b, a float64, spaceName string) (Color, error) {
	space := getSpaceByName(spaceName)
	if space == nil {
		return nil, &ParseError{input: spaceName, reason: "unknown RGB color space"}
	}

	// Create a SpaceColor in the source space
	spaceColor := NewSpaceColor(space, []float64{r, g, b}, a)

	// Convert to RGBA (sRGB)
	return spaceColor.ToRGBA(), nil
}

// getSpaceByName returns a Space by name string using the registry.
func getSpaceByName(name string) Space {
	space, _ := GetSpace(name)
	return space
}

