package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	col "github.com/SCKelemen/color"
)

func main() {
	outputDir := "docs/gamuts"
	os.MkdirAll(outputDir, 0755)

	spaces := []struct {
		name string
		fn   func(r, g, b float64) col.Color
	}{
		{"sRGB", func(r, g, b float64) col.Color { return col.RGB(r, g, b) }},
		{"DisplayP3", func(r, g, b float64) col.Color {
			c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "display-p3")
			return c
		}},
		{"AdobeRGB", func(r, g, b float64) col.Color {
			c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "a98-rgb")
			return c
		}},
		{"Rec2020", func(r, g, b float64) col.Color {
			c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "rec2020")
			return c
		}},
	}

	for _, s := range spaces {
		img := generateGamutIsometric(s.fn, 1000, 800)
		filename := filepath.Join(outputDir, fmt.Sprintf("gamut_%s.png", s.name))
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		png.Encode(f, img)
		f.Close()
		fmt.Printf("Generated %s\n", filename)
	}
}

func generateGamutIsometric(createColor func(r, g, b float64) col.Color, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with dark gray background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{20, 20, 20, 255})
		}
	}

	// Isometric projection: standard isometric angles
	// Isometric projection matrix:
	// x' = (x - y) * cos(30°)
	// y' = (x + y) / 2 - z
	cos30 := math.Cos(math.Pi / 6) // cos(30°)
	sin30 := math.Sin(math.Pi / 6) // sin(30°)
	scale := float64(width) * 0.25

	// Center of image
	centerX := float64(width) / 2
	centerY := float64(height) / 2

	// Convert RGB cube to XYZ for better visualization
	// Sample the gamut boundary more densely
	step := 0.02
	points := make([]struct {
		x, y int
		col  color.RGBA
	}, 0)

	for r := 0.0; r <= 1.0; r += step {
		for g := 0.0; g <= 1.0; g += step {
			for b := 0.0; b <= 1.0; b += step {
				// Create color in the target space
				c := createColor(r, g, b)

				// Convert to XYZ to show actual gamut shape
				xyz := col.ToXYZ(c)

				// Isometric projection of XYZ coordinates
				// Standard isometric: x' = (x - y) * cos(30°), y' = (x + y) / 2 - z
				xProj := (xyz.X - xyz.Y) * cos30 * scale
				yProj := ((xyz.X+xyz.Y)/2 - xyz.Z) * scale

				// Map to image coordinates
				xCoord := int(centerX + xProj)
				yCoord := int(centerY - yProj)

				if xCoord >= 0 && xCoord < width && yCoord >= 0 && yCoord < height {
					// Get RGB for display (convert to sRGB)
					rgbR, rgbG, rgbB, _ := c.RGBA()
					points = append(points, struct {
						x, y int
						col  color.RGBA
					}{
						x: xCoord,
						y: yCoord,
						col: color.RGBA{
							R: uint8(math.Max(0, math.Min(255, rgbR*255))),
							G: uint8(math.Max(0, math.Min(255, rgbG*255))),
							B: uint8(math.Max(0, math.Min(255, rgbB*255))),
							A: 255,
						},
					})
				}
			}
		}
	}

	// Draw points (with slight alpha blending for depth)
	for _, p := range points {
		img.Set(p.x, p.y, p.col)
		// Add slight glow for better visibility
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}
				nx, ny := p.x+dx, p.y+dy
				if nx >= 0 && nx < width && ny >= 0 && ny < height {
					// Blend with existing color
					existing := img.RGBAAt(nx, ny)
					blended := color.RGBA{
						R: uint8((int(existing.R) + int(p.col.R)) / 2),
						G: uint8((int(existing.G) + int(p.col.G)) / 2),
						B: uint8((int(existing.B) + int(p.col.B)) / 2),
						A: 255,
					}
					img.Set(nx, ny, blended)
				}
			}
		}
	}

	// Draw axes in RGB space for reference
	drawAxis(img, centerX, centerY, scale, cos30, sin30, 1.0, 0.0, 0.0, color.RGBA{255, 100, 100, 255}, "R")
	drawAxis(img, centerX, centerY, scale, cos30, sin30, 0.0, 1.0, 0.0, color.RGBA{100, 255, 100, 255}, "G")
	drawAxis(img, centerX, centerY, scale, cos30, sin30, 0.0, 0.0, 1.0, color.RGBA{100, 100, 255, 255}, "B")

	return img
}

func drawAxis(img *image.RGBA, centerX, centerY, scale, cos30, sin30, r, g, b float64, axisColor color.RGBA, label string) {
	// Convert RGB to XYZ for axis endpoint
	c := col.RGB(r, g, b)
	xyz := col.ToXYZ(c)

	// Project to isometric
	xProj := (xyz.X - xyz.Y) * cos30 * scale
	yProj := ((xyz.X+xyz.Y)/2 - xyz.Z) * scale

	xEnd := int(centerX + xProj)
	yEnd := int(centerY - yProj)

	// Draw line from origin
	drawLine(img, int(centerX), int(centerY), xEnd, yEnd, axisColor, 2)

	// Draw label
	drawLabel(img, xEnd+5, yEnd, label, axisColor)
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
		// Draw thick line
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

func drawLabel(img *image.RGBA, x, y int, label string, col color.Color) {
	// Draw label as small squares
	for i := range label {
		px := x + i*10
		py := y
		// Draw a small square for visibility
		for dy := 0; dy < 4; dy++ {
			for dx := 0; dx < 8; dx++ {
				if px+dx >= 0 && px+dx < img.Bounds().Dx() && py+dy >= 0 && py+dy < img.Bounds().Dy() {
					img.Set(px+dx, py+dy, col)
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
