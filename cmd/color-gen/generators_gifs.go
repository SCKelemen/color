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

// createGlobalPalette creates a color palette from all frames
// Uses a better approach: collect all unique colors first, then select the best 256
func createGlobalPalette(frames []image.Image) color.Palette {
	// Add transparent color first
	palette := color.Palette{color.RGBA{0, 0, 0, 0}}
	
	// Collect all unique colors from all frames
	colorMap := make(map[uint32]color.RGBA)
	
	// Sample colors from all frames more evenly
	// Use a larger step to sample more frames before hitting the limit
	step := 4
	for _, img := range frames {
		bounds := img.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y += step {
			for x := bounds.Min.X; x < bounds.Max.X; x += step {
				c := img.At(x, y)
				r, g, b, a := c.RGBA()
				
				// Skip fully transparent pixels
				if a == 0 {
					continue
				}
				
				// Quantize to reduce similar colors (reduce to 6 bits per channel for better coverage)
				// This helps group similar colors together while preserving more variety
				rq := (r >> 10) << 10
				gq := (g >> 10) << 10
				bq := (b >> 10) << 10
				
				// Create a key for this quantized color
				key := uint32(rq>>8)<<24 | uint32(gq>>8)<<16 | uint32(bq>>8)<<8 | 255
				
				// Store the actual color (not quantized) for better quality
				if _, exists := colorMap[key]; !exists {
					colorMap[key] = color.RGBA{
						R: uint8(r >> 8),
						G: uint8(g >> 8),
						B: uint8(b >> 8),
						A: uint8(a >> 8),
					}
				}
			}
		}
	}
	
	// Convert map to slice
	allColors := make([]color.RGBA, 0, len(colorMap))
	for _, c := range colorMap {
		allColors = append(allColors, c)
	}
	
	// If we have more than 255 colors, we need to select the best ones
	// Use a simple approach: sort by frequency or use a more sophisticated method
	if len(allColors) > 255 {
		// For now, just take the first 255 unique colors
		// In the future, we could use median cut or k-means clustering
		allColors = allColors[:255]
	}
	
	// Add all collected colors to palette
	palette = append(palette, allColors...)
	
	// Pad to 256 colors if needed (use last color or transparent)
	var lastColor color.Color = color.RGBA{0, 0, 0, 0}
	if len(palette) > 1 {
		lastColor = palette[len(palette)-1]
	}
	for len(palette) < 256 {
		palette = append(palette, lastColor)
	}
	
	return palette
}

