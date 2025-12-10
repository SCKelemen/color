package color_test

import (
	"fmt"
	"github.com/SCKelemen/color"
)

func ExampleRGB() {
	// Create a red color
	red := color.RGB(1.0, 0.0, 0.0)
	r, g, b, a := red.RGBA()
	fmt.Printf("Red: R=%.2f, G=%.2f, B=%.2f, A=%.2f\n", r, g, b, a)
	// Output: Red: R=1.00, G=0.00, B=0.00, A=1.00
}

func ExampleToHSL() {
	// Create RGB color
	rgb := color.RGB(1.0, 0.5, 0.0) // Orange

	// Convert to different color spaces
	hsl := color.ToHSL(rgb)
	fmt.Printf("HSL: H=%.1f, S=%.2f, L=%.2f\n", hsl.H, hsl.S, hsl.L)

	oklch := color.ToOKLCH(rgb)
	fmt.Printf("OKLCH: L=%.2f, C=%.2f, H=%.1f\n", oklch.L, oklch.C, oklch.H)
}

func ExampleLighten() {
	blue := color.RGB(0.0, 0.0, 1.0)
	lightBlue := color.Lighten(blue, 0.3)
	r, g, b, _ := lightBlue.RGBA()
	fmt.Printf("Light blue: R=%.2f, G=%.2f, B=%.2f\n", r, g, b)
}

func ExampleDarken() {
	red := color.RGB(1.0, 0.0, 0.0)
	darkRed := color.Darken(red, 0.3)
	r, g, b, _ := darkRed.RGBA()
	fmt.Printf("Dark red: R=%.2f, G=%.2f, B=%.2f\n", r, g, b)
}

func ExampleMix() {
	red := color.RGB(1.0, 0.0, 0.0)
	blue := color.RGB(0.0, 0.0, 1.0)
	purple := color.Mix(red, blue, 0.5)
	r, g, b, _ := purple.RGBA()
	fmt.Printf("Purple (mixed): R=%.2f, G=%.2f, B=%.2f\n", r, g, b)
}

func ExampleHexToRGB() {
	// Parse hex color
	rgb, err := color.HexToRGB("#FF5733")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	r, g, b, _ := rgb.RGBA()
	fmt.Printf("Hex #FF5733: R=%.2f, G=%.2f, B=%.2f\n", r, g, b)
}

func ExampleRGBToHex() {
	rgb := color.RGB(1.0, 0.34, 0.2) // Similar to #FF5733
	hex := color.RGBToHex(rgb)
	fmt.Printf("Hex: %s\n", hex)
}

func ExampleNewRGBA() {
	// Create color with transparency
	semiTransparent := color.NewRGBA(1.0, 0.0, 0.0, 0.5)
	fmt.Printf("Alpha: %.2f\n", semiTransparent.Alpha())

	// Change alpha
	moreOpaque := semiTransparent.WithAlpha(0.8)
	fmt.Printf("New alpha: %.2f\n", moreOpaque.Alpha())
}

