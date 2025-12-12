// Package color provides a comprehensive Go library for color space conversions
// and color manipulation, supporting all major CSS color spaces.
//
// Features:
//
//   - Full Color Space Support: RGB, HSL, HSV, LAB, OKLAB, LCH, OKLCH, and XYZ
//   - Alpha Channel Support: All color spaces support transparency
//   - Perceptually Uniform Operations: Lighten, darken, and other operations use OKLCH
//     for perceptually uniform results
//   - Color Mixing: Mix colors in RGB, HSL, LAB, OKLAB, LCH, or OKLCH space
//   - Gradient Generation: Generate smooth gradients with multiple stops and easing functions
//   - Utility Functions: Lighten, darken, saturate, desaturate, invert, grayscale,
//     complement, and more
//   - Standard Library Compatibility: Convert to/from image/color.Color interface
//
// Basic Usage:
//
//	import "github.com/SCKelemen/color"
//
//	// Create colors
//	red := color.RGB(1.0, 0.0, 0.0)
//	blue := color.NewRGBA(0.0, 0.0, 1.0, 0.5) // Semi-transparent
//
//	// Convert between color spaces
//	hsl := color.ToHSL(red)
//	oklch := color.ToOKLCH(blue)
//
//	// Parse colors from strings
//	parsed, _ := color.ParseColor("#FF0000")
//	parsed2, _ := color.ParseColor("rgb(255, 0, 0)")
//	parsed3, _ := color.ParseColor("oklch(0.7 0.2 120)")
//
//	// Manipulate colors
//	lightRed := color.Lighten(red, 0.3)
//	darkBlue := color.Darken(blue, 0.3)
//	mixed := color.Mix(red, blue, 0.5)
//
// Integration with lipgloss:
//
// This library provides advanced color space conversions that lipgloss doesn't support.
// While lipgloss supports basic RGB colors, this library adds:
//
//   - Perceptually uniform color spaces (OKLAB, OKLCH) for better color manipulation
//   - Advanced color space conversions (LAB, LCH, HSL, HSV)
//   - Perceptually uniform lighten/darken operations
//   - Color mixing in perceptually uniform space
//   - Alpha channel support across all color spaces
//
// Example with lipgloss:
//
//	import (
//	    "github.com/charmbracelet/lipgloss"
//	    "github.com/SCKelemen/color"
//	)
//
//	// Create a color using this library
//	bgColor := color.NewOKLCH(0.2, 0.1, 240, 1.0) // Dark blue
//
//	// Convert to lipgloss color
//	style := lipgloss.NewStyle().
//	    Background(lipgloss.Color(color.ToLipglossColor(bgColor)))
//
// See the README for more examples and documentation.
package color

