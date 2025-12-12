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

// createGlobalPalette creates a palette using median cut algorithm
// This ensures good color distribution across the actual colors in the frames
func createGlobalPalette(frames []image.Image) color.Palette {
	// Add transparent color first
	palette := color.Palette{color.RGBA{0, 0, 0, 0}}
	
	// Collect all colors from frames (sample every pixel for accuracy)
	allColors := make([]color.RGBA, 0)
	step := 2 // Sample every 2nd pixel for speed
	
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
				
				allColors = append(allColors, color.RGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(a >> 8),
				})
			}
		}
	}
	
	if len(allColors) == 0 {
		// Fallback: uniform grid
		return createUniformPalette()
	}
	
	// Use median cut to select 255 colors
	selectedColors := medianCut(allColors, 255)
	
	// Add selected colors to palette
	for _, c := range selectedColors {
		palette = append(palette, c)
	}
	
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

// medianCut implements the median cut algorithm to select representative colors
func medianCut(colors []color.RGBA, targetCount int) []color.RGBA {
	if len(colors) <= targetCount {
		return colors
	}
	
	// Start with one box containing all colors
	boxes := []*colorBox{{colors: colors}}
	
	// Split boxes until we have enough
	for len(boxes) < targetCount && len(boxes) > 0 {
		// Find the box with the largest range
		largestIdx := 0
		largestRange := 0.0
		
		for i, box := range boxes {
			range_ := box.range_()
			if range_ > largestRange {
				largestRange = range_
				largestIdx = i
			}
		}
		
		// Split the largest box
		box := boxes[largestIdx]
		left, right := box.split()
		
		// Replace the box with its two halves
		boxes = append(boxes[:largestIdx], boxes[largestIdx+1:]...)
		boxes = append(boxes, left, right)
	}
	
	// Get the average color from each box
	result := make([]color.RGBA, 0, len(boxes))
	for _, box := range boxes {
		result = append(result, box.average())
	}
	
	return result
}

type colorBox struct {
	colors []color.RGBA
}

func (b *colorBox) range_() float64 {
	if len(b.colors) == 0 {
		return 0
	}
	
	minR, maxR := 255, 0
	minG, maxG := 255, 0
	minB, maxB := 255, 0
	
	for _, c := range b.colors {
		if int(c.R) < minR {
			minR = int(c.R)
		}
		if int(c.R) > maxR {
			maxR = int(c.R)
		}
		if int(c.G) < minG {
			minG = int(c.G)
		}
		if int(c.G) > maxG {
			maxG = int(c.G)
		}
		if int(c.B) < minB {
			minB = int(c.B)
		}
		if int(c.B) > maxB {
			maxB = int(c.B)
		}
	}
	
	// Return the range of the channel with the largest range
	rangeR := float64(maxR - minR)
	rangeG := float64(maxG - minG)
	rangeB := float64(maxB - minB)
	
	if rangeR >= rangeG && rangeR >= rangeB {
		return rangeR
	}
	if rangeG >= rangeB {
		return rangeG
	}
	return rangeB
}

func (b *colorBox) split() (*colorBox, *colorBox) {
	if len(b.colors) == 0 {
		return &colorBox{colors: []color.RGBA{}}, &colorBox{colors: []color.RGBA{}}
	}
	
	// Find the channel with the largest range
	minR, maxR := 255, 0
	minG, maxG := 255, 0
	minB, maxB := 255, 0
	
	for _, c := range b.colors {
		if int(c.R) < minR {
			minR = int(c.R)
		}
		if int(c.R) > maxR {
			maxR = int(c.R)
		}
		if int(c.G) < minG {
			minG = int(c.G)
		}
		if int(c.G) > maxG {
			maxG = int(c.G)
		}
		if int(c.B) < minB {
			minB = int(c.B)
		}
		if int(c.B) > maxB {
			maxB = int(c.B)
		}
	}
	
	rangeR := maxR - minR
	rangeG := maxG - minG
	rangeB := maxB - minB
	
	// Sort by the channel with the largest range
	if rangeR >= rangeG && rangeR >= rangeB {
		// Sort by R
		sort.Slice(b.colors, func(i, j int) bool {
			return b.colors[i].R < b.colors[j].R
		})
	} else if rangeG >= rangeB {
		// Sort by G
		sort.Slice(b.colors, func(i, j int) bool {
			return b.colors[i].G < b.colors[j].G
		})
	} else {
		// Sort by B
		sort.Slice(b.colors, func(i, j int) bool {
			return b.colors[i].B < b.colors[j].B
		})
	}
	
	// Split at median (colors are already sorted)
	mid := len(b.colors) / 2
	left := b.colors[:mid]
	right := b.colors[mid:]
	
	return &colorBox{colors: left}, &colorBox{colors: right}
}

func (b *colorBox) average() color.RGBA {
	if len(b.colors) == 0 {
		return color.RGBA{0, 0, 0, 255}
	}
	
	var sumR, sumG, sumB int
	for _, c := range b.colors {
		sumR += int(c.R)
		sumG += int(c.G)
		sumB += int(c.B)
	}
	
	return color.RGBA{
		R: uint8(sumR / len(b.colors)),
		G: uint8(sumG / len(b.colors)),
		B: uint8(sumB / len(b.colors)),
		A: 255,
	}
}

// createUniformPalette creates a uniform grid palette as fallback
func createUniformPalette() color.Palette {
	palette := color.Palette{color.RGBA{0, 0, 0, 0}}
	colors := make([]color.Color, 0, 255)
	
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 7; b++ {
				R := uint8((r * 255) / 5)
				G := uint8((g * 255) / 5)
				B := uint8((b * 255) / 6)
				colors = append(colors, color.RGBA{R, G, B, 255})
			}
		}
	}
	
	palette = append(palette, colors...)
	
	var lastColor color.Color = color.RGBA{0, 0, 0, 0}
	if len(palette) > 1 {
		lastColor = palette[len(palette)-1]
	}
	for len(palette) < 256 {
		palette = append(palette, lastColor)
	}
	
	return palette
}

