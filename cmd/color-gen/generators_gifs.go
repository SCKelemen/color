package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"sort"
)

func generateGIFs() error {
	animationsDir := "docs/animations"
	outputDir := "docs/models"
	os.MkdirAll(outputDir, 0755)

	// Find all animation frame directories
	entries, err := os.ReadDir(animationsDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modelName := entry.Name()
		frameDir := filepath.Join(animationsDir, modelName)

		// Read all PNG frames
		frameFiles, err := filepath.Glob(filepath.Join(frameDir, "frame_*.png"))
		if err != nil {
			panic(err)
		}

		if len(frameFiles) == 0 {
			continue
		}

		// Sort frames by filename
		sort.Strings(frameFiles)

		// Load all frames
		frames := make([]image.Image, 0, len(frameFiles))
		for _, frameFile := range frameFiles {
			f, err := os.Open(frameFile)
			if err != nil {
				panic(err)
			}
			img, err := png.Decode(f)
			f.Close()
			if err != nil {
				panic(err)
			}
			frames = append(frames, img)
		}

		// Create GIF from frames
		gifFile := filepath.Join(outputDir, fmt.Sprintf("model_%s.gif", modelName))
		if err := createGIF(gifFile, frames); err != nil {
			panic(err)
		}

		fmt.Printf("Generated %s from %d frames\n", gifFile, len(frames))
	}
	return nil
}

func createGIF(filename string, frames []image.Image) error {
	// Convert frames to palette images for GIF
	paletteFrames := make([]*image.Paletted, len(frames))
	delays := make([]int, len(frames))

	if len(frames) == 0 {
		return fmt.Errorf("no frames to encode")
	}

	// Create a global palette from all frames for better color consistency
	palette := createGlobalPalette(frames)

	for i, frame := range frames {
		// Convert to paletted image
		bounds := frame.Bounds()
		paletted := image.NewPaletted(bounds, palette)
		
		// Convert RGBA to palette using nearest color matching
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := frame.At(x, y)
				// Find nearest color in palette
				paletted.Set(x, y, findNearestColorInPalette(c, palette))
			}
		}

		paletteFrames[i] = paletted
		delays[i] = 5 // 50ms delay (5 * 10ms = 50ms) for smooth animation
	}

	// Create GIF
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	g := &gif.GIF{
		Image: paletteFrames,
		Delay: delays,
		// LoopCount: 0 means infinite loop
		LoopCount: 0,
	}

	return gif.EncodeAll(f, g)
}

func findNearestColorInPalette(c color.Color, palette color.Palette) color.Color {
	if len(palette) == 0 {
		return c
	}
	
	r1, g1, b1, a1 := c.RGBA()
	minDist := uint32(^uint32(0))
	bestColor := palette[0]
	
	for _, p := range palette {
		r2, g2, b2, a2 := p.RGBA()
		
		// Calculate color distance (weighted by alpha)
		dr := int32(r1) - int32(r2)
		dg := int32(g1) - int32(g2)
		db := int32(b1) - int32(b2)
		da := int32(a1) - int32(a2)
		
		// Weight RGB more than alpha, and use squared distance
		dist := uint32(dr*dr + dg*dg + db*db + da*da/4)
		
		if dist < minDist {
			minDist = dist
			bestColor = p
		}
	}
	
	return bestColor
}

// createGlobalPalette creates an evenly distributed palette across the RGB gamut
// Since these are color space visualizations, we generate a uniform grid in RGB space
func createGlobalPalette(frames []image.Image) color.Palette {
	// Add transparent color first
	palette := color.Palette{color.RGBA{0, 0, 0, 0}}
	
	// Generate evenly distributed colors across RGB gamut
	// We want 255 colors, so we'll use a 6x6x7 grid (252 colors) + 3 extra
	// This gives us good coverage of the RGB cube
	
	// Create a uniform grid in RGB space
	// Using 6 levels for R and G, 7 for B to get close to 255
	colors := make([]color.Color, 0, 255)
	
	// Generate colors on a uniform grid
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 7; b++ {
				// Map to 0-255 range
				R := uint8((r * 255) / 5)
				G := uint8((g * 255) / 5)
				B := uint8((b * 255) / 6)
				colors = append(colors, color.RGBA{R, G, B, 255})
			}
		}
	}
	
	// Add a few extra colors to fill to 255 (6*6*7 = 252, need 3 more)
	// Add some saturated colors that might be missing
	extraColors := []color.RGBA{
		{255, 255, 255, 255}, // White
		{128, 128, 128, 255}, // Gray
		{0, 0, 0, 255},       // Black (though transparent is already first)
	}
	
	for _, ec := range extraColors {
		if len(colors) < 255 {
			colors = append(colors, ec)
		}
	}
	
	// Add all colors to palette
	palette = append(palette, colors...)
	
	// Pad to 256 colors if needed
	var lastColor color.Color = color.RGBA{0, 0, 0, 0}
	if len(palette) > 1 {
		lastColor = palette[len(palette)-1]
	}
	for len(palette) < 256 {
		palette = append(palette, lastColor)
	}
	
	return palette
}

