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

	// Create a global palette from all frames using median cut for better color consistency
	palette := createGlobalPalette(frames)
	fmt.Printf("Created palette with %d colors\n", len(palette))

	for i, frame := range frames {
		// Convert to paletted image with Floyd-Steinberg dithering
		bounds := frame.Bounds()
		paletted := image.NewPaletted(bounds, palette)

		// Apply Floyd-Steinberg dithering for smoother gradients
		applyFloydSteinbergDithering(frame, paletted, palette)

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

// applyFloydSteinbergDithering applies Floyd-Steinberg dithering to reduce color banding
func applyFloydSteinbergDithering(src image.Image, dst *image.Paletted, palette color.Palette) {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create error buffer for dithering (RGBA components)
	errors := make([][]struct{ r, g, b int }, height)
	for i := range errors {
		errors[i] = make([]struct{ r, g, b int }, width)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get original color with accumulated error
			oldR, oldG, oldB, oldA := src.At(x, y).RGBA()

			// Skip transparent pixels
			if oldA == 0 {
				dst.Set(x, y, palette[0]) // First color is transparent
				continue
			}

			// Convert to 8-bit and add accumulated error
			r := int(oldR>>8) + errors[y-bounds.Min.Y][x-bounds.Min.X].r
			g := int(oldG>>8) + errors[y-bounds.Min.Y][x-bounds.Min.X].g
			b := int(oldB>>8) + errors[y-bounds.Min.Y][x-bounds.Min.X].b

			// Clamp to valid range
			if r < 0 { r = 0 }
			if r > 255 { r = 255 }
			if g < 0 { g = 0 }
			if g > 255 { g = 255 }
			if b < 0 { b = 0 }
			if b > 255 { b = 255 }

			oldColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}

			// Find nearest color in palette
			newColor := findNearestColorInPalette(oldColor, palette)
			dst.Set(x, y, newColor)

			// Calculate quantization error
			newR, newG, newB, _ := newColor.RGBA()
			errR := r - int(newR>>8)
			errG := g - int(newG>>8)
			errB := b - int(newB>>8)

			// Distribute error to neighboring pixels (Floyd-Steinberg)
			// Format: [y offset][x offset] = error * multiplier / 16
			if x+1 < bounds.Max.X {
				errors[y-bounds.Min.Y][x+1-bounds.Min.X].r += errR * 7 / 16
				errors[y-bounds.Min.Y][x+1-bounds.Min.X].g += errG * 7 / 16
				errors[y-bounds.Min.Y][x+1-bounds.Min.X].b += errB * 7 / 16
			}
			if y+1 < bounds.Max.Y {
				if x-1 >= bounds.Min.X {
					errors[y+1-bounds.Min.Y][x-1-bounds.Min.X].r += errR * 3 / 16
					errors[y+1-bounds.Min.Y][x-1-bounds.Min.X].g += errG * 3 / 16
					errors[y+1-bounds.Min.Y][x-1-bounds.Min.X].b += errB * 3 / 16
				}
				errors[y+1-bounds.Min.Y][x-bounds.Min.X].r += errR * 5 / 16
				errors[y+1-bounds.Min.Y][x-bounds.Min.X].g += errG * 5 / 16
				errors[y+1-bounds.Min.Y][x-bounds.Min.X].b += errB * 5 / 16
				if x+1 < bounds.Max.X {
					errors[y+1-bounds.Min.Y][x+1-bounds.Min.X].r += errR * 1 / 16
					errors[y+1-bounds.Min.Y][x+1-bounds.Min.X].g += errG * 1 / 16
					errors[y+1-bounds.Min.Y][x+1-bounds.Min.X].b += errB * 1 / 16
				}
			}
		}
	}
}

func findNearestColorInPalette(c color.Color, palette color.Palette) color.Color {
	if len(palette) == 0 {
		return c
	}
	
	r1, g1, b1, a1 := c.RGBA()
	
	// Skip transparent pixels - return transparent color
	if a1 == 0 {
		return palette[0] // First color is transparent
	}
	
	minDist := uint32(^uint32(0))
	bestColor := palette[0]
	
	// Use perceptual color distance (weighted RGB)
	// Human eye is more sensitive to green, so weight it more
	for _, p := range palette {
		r2, g2, b2, a2 := p.RGBA()
		
		// Skip transparent colors in palette (except index 0)
		if a2 == 0 && p != palette[0] {
			continue
		}
		
		// Calculate perceptual color distance
		// Weight green more (human eye is more sensitive to green)
		dr := int32(r1) - int32(r2)
		dg := int32(g1) - int32(g2)
		db := int32(b1) - int32(b2)
		
		// Perceptual weights: R=0.3, G=0.59, B=0.11 (approximate)
		// But for simplicity, we'll use squared distance with green weighted more
		dist := uint32(dr*dr*3 + dg*dg*6 + db*db*1) // Green weighted 2x
		
		if dist < minDist {
			minDist = dist
			bestColor = p
		}
	}
	
	return bestColor
}

// createGlobalPalette creates an optimized palette using median cut algorithm
// This samples colors from all frames to create a palette that best represents the animation
func createGlobalPalette(frames []image.Image) color.Palette {
	// Sample colors from all frames (every 4th pixel to reduce processing time)
	colorSamples := make([]color.RGBA, 0, 10000)

	for _, frame := range frames {
		bounds := frame.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y += 4 {
			for x := bounds.Min.X; x < bounds.Max.X; x += 4 {
				r, g, b, a := frame.At(x, y).RGBA()
				// Skip transparent pixels
				if a > 0 {
					colorSamples = append(colorSamples, color.RGBA{
						R: uint8(r >> 8),
						G: uint8(g >> 8),
						B: uint8(b >> 8),
						A: 255,
					})
				}
			}
		}
	}

	fmt.Printf("Sampled %d colors from %d frames\n", len(colorSamples), len(frames))

	// Use median cut to reduce to 255 colors (1 slot reserved for transparency)
	if len(colorSamples) > 255 {
		colorSamples = medianCut(colorSamples, 255)
	}

	// Build palette with transparent color first
	palette := color.Palette{color.RGBA{0, 0, 0, 0}}
	for _, c := range colorSamples {
		palette = append(palette, c)
	}

	// Fill remaining slots if needed
	for len(palette) < 256 {
		palette = append(palette, color.RGBA{0, 0, 0, 255})
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

