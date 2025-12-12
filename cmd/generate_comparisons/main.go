package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	colorlib "github.com/SCKelemen/color"
)

const (
	swatchWidth  = 80
	swatchHeight = 80
	padding      = 20
	fontSize     = 12
)

func main() {
	// Create output directory
	outDir := "docs/comparisons"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		panic(err)
	}

	fmt.Println("Generating comparison visualizations...")

	// 1. RGB vs OKLCH Lightening
	fmt.Println("  - RGB vs OKLCH lightening examples")
	generateLighteningComparison(filepath.Join(outDir, "lightening_comparison.png"))

	// 2. Saturation comparison
	fmt.Println("  - Saturation comparison")
	generateSaturationComparison(filepath.Join(outDir, "saturation_comparison.png"))

	// 3. Multi-stop gradient comparison
	fmt.Println("  - Multi-stop gradient comparison")
	generateMultiStopComparison(filepath.Join(outDir, "multistop_comparison.png"))

	// 4. Gamut mapping comparison
	fmt.Println("  - Gamut mapping comparison")
	generateGamutMappingComparison(filepath.Join(outDir, "gamut_mapping_comparison.png"))

	fmt.Println("Done! Visualizations saved to", outDir)
}

// generateLighteningComparison creates a side-by-side comparison of RGB vs OKLCH lightening
func generateLighteningComparison(filename string) {
	colors := []struct {
		name string
		base colorlib.Color
	}{
		{"Blue", colorlib.RGB(0, 0, 1)},
		{"Green", colorlib.RGB(0, 1, 0)},
		{"Red", colorlib.RGB(1, 0, 0)},
	}

	// Create image: 3 rows (colors) × 2 columns (RGB vs OKLCH) × 3 swatches each (base, +20%, +40%)
	imgWidth := (swatchWidth*3 + padding*2) * 2
	imgHeight := (swatchHeight + padding) * len(colors)
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// White background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	for i, c := range colors {
		yOffset := i * (swatchHeight + padding)

		// RGB lightening (naive)
		rgbBase := c.base
		r, g, b, a := rgbBase.RGBA()
		rgbLight1 := colorlib.NewRGBA(r+0.2, g+0.2, b+0.2, a)
		rgbLight2 := colorlib.NewRGBA(r+0.4, g+0.4, b+0.4, a)

		// Draw RGB swatches
		drawSwatch(img, 0, yOffset, rgbBase)
		drawSwatch(img, swatchWidth+padding, yOffset, rgbLight1)
		drawSwatch(img, (swatchWidth+padding)*2, yOffset, rgbLight2)

		// OKLCH lightening (perceptually uniform)
		oklchBase := c.base
		oklchLight1 := colorlib.Lighten(oklchBase, 0.2)
		oklchLight2 := colorlib.Lighten(oklchBase, 0.4)

		xOffsetOKLCH := (swatchWidth*3 + padding*2)
		drawSwatch(img, xOffsetOKLCH, yOffset, oklchBase)
		drawSwatch(img, xOffsetOKLCH+swatchWidth+padding, yOffset, oklchLight1)
		drawSwatch(img, xOffsetOKLCH+(swatchWidth+padding)*2, yOffset, oklchLight2)
	}

	saveImage(img, filename)
}

// generateSaturationComparison creates comparison of saturation in different spaces
func generateSaturationComparison(filename string) {
	baseColor := colorlib.RGB(0.5, 0.3, 0.8) // Purple

	// Create a row of 5 swatches showing saturation progression
	imgWidth := (swatchWidth + padding) * 5
	imgHeight := swatchHeight * 2 // Two rows: HSL and OKLCH

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight+padding))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// HSL saturation (top row)
	hsl := colorlib.ToHSL(baseColor)
	for i := 0; i < 5; i++ {
		s := float64(i) * 0.25 // 0%, 25%, 50%, 75%, 100%
		hslCopy := colorlib.NewHSL(hsl.H, s, hsl.L, 1.0)
		drawSwatch(img, i*(swatchWidth+padding), 0, hslCopy)
	}

	// OKLCH saturation (bottom row)
	for i := 0; i < 5; i++ {
		amount := float64(i) * 0.25
		saturated := colorlib.Saturate(baseColor, amount)
		drawSwatch(img, i*(swatchWidth+padding), swatchHeight+padding, saturated)
	}

	saveImage(img, filename)
}

// generateMultiStopComparison creates comparison of multi-stop gradients
func generateMultiStopComparison(filename string) {
	red := colorlib.RGB(1, 0, 0)
	yellow := colorlib.RGB(1, 1, 0)
	blue := colorlib.RGB(0, 0, 1)

	stops := []colorlib.GradientStop{
		{Color: red, Position: 0.0},
		{Color: yellow, Position: 0.5},
		{Color: blue, Position: 1.0},
	}

	steps := 100
	imgHeight := swatchHeight * 3 // 3 rows: RGB, HSL, OKLCH
	imgWidth := steps * 10

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight+padding*2))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// RGB gradient
	rgbGradient := colorlib.GradientMultiStop(stops, steps, colorlib.GradientRGB)
	for i, c := range rgbGradient {
		drawGradientStrip(img, i*10, 0, 10, swatchHeight, c)
	}

	// HSL gradient
	hslGradient := colorlib.GradientMultiStop(stops, steps, colorlib.GradientHSL)
	for i, c := range hslGradient {
		drawGradientStrip(img, i*10, swatchHeight+padding, 10, swatchHeight, c)
	}

	// OKLCH gradient
	oklchGradient := colorlib.GradientMultiStop(stops, steps, colorlib.GradientOKLCH)
	for i, c := range oklchGradient {
		drawGradientStrip(img, i*10, (swatchHeight+padding)*2, 10, swatchHeight, c)
	}

	saveImage(img, filename)
}

// generateGamutMappingComparison shows different gamut mapping strategies
func generateGamutMappingComparison(filename string) {
	// Create a vivid Display P3 color that's out of sRGB gamut
	p3Color := colorlib.NewOKLCH(0.7, 0.25, 150, 1.0) // Vivid teal

	strategies := []struct {
		name     string
		strategy colorlib.GamutMapping
	}{
		{"Original", colorlib.GamutClip}, // Will use as-is for reference
		{"Clip", colorlib.GamutClip},
		{"Preserve Lightness", colorlib.GamutPreserveLightness},
		{"Preserve Chroma", colorlib.GamutPreserveChroma},
		{"Project", colorlib.GamutProject},
	}

	imgWidth := (swatchWidth + padding) * len(strategies)
	imgHeight := swatchHeight

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw original
	drawSwatch(img, 0, 0, p3Color)

	// Draw each mapping strategy
	for i := 1; i < len(strategies); i++ {
		mapped := colorlib.MapToGamut(p3Color, strategies[i].strategy)
		drawSwatch(img, i*(swatchWidth+padding), 0, mapped)
	}

	saveImage(img, filename)
}

// Helper functions

func drawSwatch(img *image.RGBA, x, y int, c colorlib.Color) {
	r, g, b, a := c.RGBA()
	col := color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: uint8(a * 255),
	}

	rect := image.Rect(x, y, x+swatchWidth, y+swatchHeight)
	draw.Draw(img, rect, &image.Uniform{col}, image.Point{}, draw.Src)
}

func drawGradientStrip(img *image.RGBA, x, y, width, height int, c colorlib.Color) {
	r, g, b, a := c.RGBA()
	col := color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: uint8(a * 255),
	}

	rect := image.Rect(x, y, x+width, y+height)
	draw.Draw(img, rect, &image.Uniform{col}, image.Point{}, draw.Src)
}

func saveImage(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
