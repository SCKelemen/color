package main

import (
	"image"
	"image/color"
	"os"

	"golang.org/x/image/font/opentype"
)

// Shared constants
const (
	scale          = 3  // Render at 3x for better quality
	fontSize       = 16 // Font size in points
	gradientWidth  = 830
	gradientHeight = 50
	cornerRadius   = 8
	padding        = 8  // For gradients
	paddingStops   = 20 // For stops
)

var defaultTT *opentype.Font

func init() {
	// Try to load Roboto font
	robotoFonts := []string{
		"fonts/Roboto/static/Roboto-Regular.ttf",
		"fonts/Roboto/static/Roboto-Medium.ttf",
		"fonts/Roboto/Roboto-VariableFont_wdth,wght.ttf",
		"fonts/Roboto/Roboto-Regular.ttf",
	}

	for _, fontPath := range robotoFonts {
		if _, err := os.Stat(fontPath); err == nil {
			fontData, err := os.ReadFile(fontPath)
			if err == nil {
				tt, err := opentype.Parse(fontData)
				if err == nil {
					defaultTT = tt
					return
				}
			}
		}
	}
}

// Utility functions
func clamp255(v float64) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return int(v)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func findNearestColor(img *image.RGBA, x, y, radius int) color.RGBA {
	for r := 1; r <= radius; r++ {
		for dy := -r; dy <= r; dy++ {
			for dx := -r; dx <= r; dx++ {
				if dx*dx+dy*dy > r*r {
					continue
				}
				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < img.Bounds().Dx() && ny >= 0 && ny < img.Bounds().Dy() {
					c := img.RGBAAt(nx, ny)
					if c.A > 0 {
						return c
					}
				}
			}
		}
	}
	return color.RGBA{0, 0, 0, 0}
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, col color.Color, thickness int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	x, y := x0, y0
	for {
		for ty := -thickness / 2; ty <= thickness/2; ty++ {
			for tx := -thickness / 2; tx <= thickness/2; tx++ {
				nx, ny := x+tx, y+ty
				if nx >= 0 && nx < img.Bounds().Dx() && ny >= 0 && ny < img.Bounds().Dy() {
					img.Set(nx, ny, col)
				}
			}
		}
		if x == x1 && y == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}
