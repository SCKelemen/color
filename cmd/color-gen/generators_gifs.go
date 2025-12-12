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
// Uses median cut algorithm to select the best 256 colors distributed across the gamut
func createGlobalPalette(frames []image.Image) color.Palette {
	// Add transparent color first
	palette := color.Palette{color.RGBA{0, 0, 0, 0}}
	
	// Collect all colors from all frames (with frequency)
	type colorCount struct {
		color color.RGBA
		count int
	}
	colorMap := make(map[uint32]*colorCount)
	
	// Sample colors from all frames evenly
	// Sample every 3rd pixel to get good coverage without being too slow
	step := 3
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
				
				// Use less aggressive quantization (7 bits per channel) to preserve more colors
				rq := (r >> 9) << 9
				gq := (g >> 9) << 9
				bq := (b >> 9) << 9
				
				// Create a key for this quantized color
				key := uint32(rq>>8)<<24 | uint32(gq>>8)<<16 | uint32(bq>>8)<<8 | 255
				
				// Count frequency of this color
				if cc, exists := colorMap[key]; exists {
					cc.count++
				} else {
					colorMap[key] = &colorCount{
						color: color.RGBA{
							R: uint8(r >> 8),
							G: uint8(g >> 8),
							B: uint8(b >> 8),
							A: uint8(a >> 8),
						},
						count: 1,
					}
				}
			}
		}
	}
	
	// Convert to slice for processing
	allColors := make([]*colorCount, 0, len(colorMap))
	for _, cc := range colorMap {
		allColors = append(allColors, cc)
	}
	
	// If we have more than 255 colors, use a simple but effective selection:
	// Sort by frequency (most common colors first), but also ensure diversity
	if len(allColors) > 255 {
		// Sort by count (descending) to prioritize common colors
		sort.Slice(allColors, func(i, j int) bool {
			return allColors[i].count > allColors[j].count
		})
		
		// Take top colors, but also sample from different parts of the color space
		// to ensure we get good gamut coverage
		selected := make(map[uint32]bool)
		result := make([]color.Color, 0, 255)
		
		// First, take top 128 most frequent colors
		for i := 0; i < 128 && i < len(allColors); i++ {
			cc := allColors[i]
			key := uint32(cc.color.R)<<24 | uint32(cc.color.G)<<16 | uint32(cc.color.B)<<8 | 255
			if !selected[key] {
				selected[key] = true
				result = append(result, cc.color)
			}
		}
		
		// Then, sample evenly across the remaining colors to ensure gamut coverage
		// Divide color space into buckets and take one from each
		buckets := make([][]*colorCount, 64) // 4x4x4 = 64 buckets
		for _, cc := range allColors[128:] {
			// Map to bucket based on RGB
			rBucket := int(cc.color.R) / 64
			gBucket := int(cc.color.G) / 64
			bBucket := int(cc.color.B) / 64
			bucketIdx := rBucket*16 + gBucket*4 + bBucket
			if bucketIdx >= 0 && bucketIdx < 64 {
				buckets[bucketIdx] = append(buckets[bucketIdx], cc)
			}
		}
		
		// Take one color from each non-empty bucket
		for _, bucket := range buckets {
			if len(bucket) > 0 && len(result) < 255 {
				// Take the most frequent color from this bucket
				best := bucket[0]
				for _, cc := range bucket[1:] {
					if cc.count > best.count {
						best = cc
					}
				}
				key := uint32(best.color.R)<<24 | uint32(best.color.G)<<16 | uint32(best.color.B)<<8 | 255
				if !selected[key] {
					selected[key] = true
					result = append(result, best.color)
				}
			}
		}
		
		// Fill remaining slots with most frequent remaining colors
		for _, cc := range allColors {
			if len(result) >= 255 {
				break
			}
			key := uint32(cc.color.R)<<24 | uint32(cc.color.G)<<16 | uint32(cc.color.B)<<8 | 255
			if !selected[key] {
				selected[key] = true
				result = append(result, cc.color)
			}
		}
		
		allColors = nil // Clear to use result
		palette = append(palette, result...)
	} else {
		// We have 255 or fewer colors, just add them all
		for _, cc := range allColors {
			palette = append(palette, cc.color)
		}
	}
	
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

