package color

import "fmt"

// ToLipglossColor converts any color to a lipgloss-compatible hex string.
// This is a convenience function for use with the lipgloss library.
//
// Example:
//
//	import (
//	    "github.com/charmbracelet/lipgloss"
//	    "github.com/SCKelemen/color"
//	)
//
//	bgColor := color.NewOKLCH(0.2, 0.1, 240, 1.0)
//	style := lipgloss.NewStyle().
//	    Background(lipgloss.Color(color.ToLipglossColor(bgColor)))
func ToLipglossColor(c Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x",
		uint8(r*255), uint8(g*255), uint8(b*255))
}

// ToLipglossColorWithAlpha converts any color to a lipgloss-compatible hex string
// including alpha channel (if supported by lipgloss version).
//
// Note: Some versions of lipgloss may not support alpha in hex colors.
// Use ToLipglossColor() if you don't need alpha.
func ToLipglossColorWithAlpha(c Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x%02x",
		uint8(r*255), uint8(g*255), uint8(b*255), uint8(a*255))
}

