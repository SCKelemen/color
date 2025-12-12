package color_test

import (
	"fmt"

	"github.com/SCKelemen/color"
)

// Example demonstrating color difference measurement
func ExampleDeltaE2000() {
	color1 := color.RGB(1, 0, 0)     // Red
	color2 := color.RGB(0.95, 0, 0)  // Slightly darker red

	diff := color.DeltaE2000(color1, color2)

	if diff < 1.0 {
		fmt.Println("Colors are perceptually identical")
	} else if diff < 2.0 {
		fmt.Println("Small difference")
	} else {
		fmt.Println("Noticeable difference")
	}
	// Output: Noticeable difference
}

// Example demonstrating gamut mapping
func ExampleMapToGamut() {
	// Create a vivid out-of-gamut color
	vivid := color.NewOKLCH(0.7, 0.35, 150, 1.0)

	// Map to sRGB gamut preserving lightness
	mapped := color.MapToGamut(vivid, color.GamutPreserveLightness)

	// Check if now in gamut
	if color.InGamut(mapped) {
		fmt.Println("Successfully mapped to gamut")
	}
	// Output: Successfully mapped to gamut
}

// Example demonstrating wide-gamut workflow
func ExampleConvertToRGBSpace() {
	// Start with an sRGB color
	srgbColor := color.RGB(1, 0, 0)

	// Convert to Display P3 (wider gamut)
	p3Color, _ := color.ConvertToRGBSpace(srgbColor, "display-p3")

	// Manipulate in perceptual space
	lighter := color.Lighten(p3Color, 0.2)

	// Convert back to Display P3
	result, _ := color.ConvertToRGBSpace(lighter, "display-p3")

	fmt.Printf("Result color space: Display P3\n")
	fmt.Printf("Alpha preserved: %.1f\n", result.Alpha())
	// Output:
	// Result color space: Display P3
	// Alpha preserved: 1.0
}

// Example showing gradient with hue interpolation
func ExampleGradientInSpace() {
	red := color.RGB(1, 0, 0)
	blue := color.RGB(0, 0, 1)

	// Generate smooth gradient in OKLCH space
	gradient := color.GradientInSpace(red, blue, 5, color.GradientOKLCH)

	fmt.Printf("Generated %d colors\n", len(gradient))
	// Output: Generated 5 colors
}

// Example demonstrating theme palette generation
func ExampleLighten_themePalette() {
	brandColor := color.RGB(0.23, 0.51, 0.96) // Blue

	// Generate palette
	palette := map[string]string{
		"light":   color.RGBToHex(color.Lighten(brandColor, 0.2)),
		"base":    color.RGBToHex(brandColor),
		"dark":    color.RGBToHex(color.Darken(brandColor, 0.2)),
	}

	fmt.Printf("Palette generated with %d shades\n", len(palette))
	// Output: Palette generated with 3 shades
}

// Example showing color space registration
func ExampleRegisterSpace() {
	// Register a custom space (in practice, use existing spaces)
	customSpace := color.DisplayP3Space
	color.RegisterSpace("my-custom-space", customSpace)

	// Retrieve it later
	space, ok := color.GetSpace("my-custom-space")
	if ok {
		fmt.Printf("Retrieved space: %s\n", space.Name())
	}
	// Output: Retrieved space: display-p3
}

// Example showing metadata usage
func ExampleMetadata() {
	meta := color.Metadata(color.DisplayP3Space)

	fmt.Printf("Space: %s\n", meta.Name)
	fmt.Printf("Gamut: %.2fx sRGB\n", meta.GamutVolumeRelativeToSRGB)
	fmt.Printf("Is RGB: %v\n", meta.IsRGB)
	// Output:
	// Space: display-p3
	// Gamut: 1.26x sRGB
	// Is RGB: true
}

// Example showing HWB color space
func ExampleNewHWB() {
	// Create color using HWB (Hue, Whiteness, Blackness)
	hwb := color.NewHWB(120, 0.2, 0.1, 1.0) // Green with some white and black

	// Convert to RGB
	r, _, _, _ := hwb.RGBA()
	fmt.Printf("HWB color converted: R=%.2f\n", r)
	// Output: HWB color converted: R=0.20
}

// Example showing LUV color space
func ExampleToLUV() {
	c := color.RGB(0.5, 0.3, 0.7) // Purple

	// Convert to CIELUV
	luv := color.ToLUV(c)

	fmt.Printf("LUV L component: %.1f\n", luv.L)
	// Output: LUV L component: 42.8
}

// Example showing chromatic adaptation
func ExampleAdaptD65ToD50() {
	// XYZ values for a color with D65 white point
	x, y, z := 0.4124, 0.2126, 0.0193

	// Adapt to D50 white point (used by ProPhoto RGB)
	xD50, yD50, zD50 := color.AdaptD65ToD50(x, y, z)

	fmt.Printf("Adapted to D50: X=%.4f Y=%.4f Z=%.4f\n", xD50, yD50, zD50)
	// Output: Adapted to D50: X=0.4360 Y=0.2224 Z=0.0139
}

// Example showing multi-stop gradients
func ExampleGradientMultiStop() {
	stops := []color.GradientStop{
		{Color: color.RGB(1, 0, 0), Position: 0.0},   // Red
		{Color: color.RGB(1, 1, 0), Position: 0.5},   // Yellow
		{Color: color.RGB(0, 0, 1), Position: 1.0},   // Blue
	}

	gradient := color.GradientMultiStop(stops, 10, color.GradientOKLCH)

	fmt.Printf("Multi-stop gradient: %d colors\n", len(gradient))
	// Output: Multi-stop gradient: 10 colors
}

// Example showing saturation adjustment
func ExampleSaturate() {
	dullColor := color.RGB(0.6, 0.5, 0.55) // Grayish pink

	// Make it more vivid
	vivid := color.Saturate(dullColor, 0.3) // 30% more saturated

	oklch := color.ToOKLCH(vivid)
	fmt.Printf("Increased chroma: %.2f\n", oklch.C)
	// Output: Increased chroma: 0.29
}

// Example demonstrating perceptually uniform operations
func ExampleToOKLCH() {
	c := color.RGB(0.5, 0.3, 0.7)

	// Convert to OKLCH for manipulation
	oklch := color.ToOKLCH(c)

	// Adjust in perceptual space
	oklch.L += 0.1 // 10% lighter
	oklch.C *= 1.2 // 20% more saturated

	// Convert back
	r, _, _, _ := oklch.RGBA()
	fmt.Printf("Adjusted color: R=%.2f\n", r)
	// Output: Adjusted color: R=0.64
}

// Example showing checking if color is in gamut
func ExampleInGamut() {
	// sRGB color - should be in gamut
	inGamut := color.RGB(0.5, 0.6, 0.7)
	fmt.Printf("In gamut: %v\n", color.InGamut(inGamut))

	// Out of gamut color (gets clamped by NewRGBA)
	outOfGamut := color.NewRGBA(1.2, 0.5, 0.3, 1.0)
	fmt.Printf("Out of gamut check: %v\n", color.InGamut(outOfGamut))
	// Output:
	// In gamut: true
	// Out of gamut check: true
}

// Example showing color mixing
func ExampleMixOKLCH() {
	red := color.RGB(1, 0, 0)
	blue := color.RGB(0, 0, 1)

	// Mix 50/50 in OKLCH space
	purple := color.MixOKLCH(red, blue, 0.5)

	hex := color.RGBToHex(purple)
	fmt.Printf("Mixed color: %s\n", hex)
	// Output: Mixed color: #ba00c1
}
