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
	outputDir := "docs/models"
	os.MkdirAll(outputDir, 0755)

	// Generate color model diagrams
	models := []struct {
		name string
		fn   func(width, height int) *image.RGBA
	}{
		{"rgb_cube", generateRGBCube},
		{"hsl_cylinder", generateHSLCylinder},
		{"lab_space", generateLABSpace},
		{"oklch_space", generateOKLCHSpace},
	}

	for _, m := range models {
		img := m.fn(800, 800)
		filename := filepath.Join(outputDir, fmt.Sprintf("model_%s.png", m.name))
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		png.Encode(f, img)
		f.Close()
		fmt.Printf("Generated %s\n", filename)
	}
}

func generateRGBCube(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	// Draw RGB cube showing a slice through the cube
	// Show a slice at G=0.5 to see R and B variations
	centerX := float64(width) / 2
	centerY := float64(height) / 2
	scale := float64(width) * 0.4

	// Sample RGB colors in a 2D slice
	g := 0.5 // Fixed green value
	step := 0.01
	for r := 0.0; r <= 1.0; r += step {
		for b := 0.0; b <= 1.0; b += step {
			c := col.RGB(r, g, b)
			rgbR, rgbG, rgbB, _ := c.RGBA()

			// Map to image coordinates (R on X axis, B on Y axis)
			px := int(centerX + (r-0.5)*scale)
			py := int(centerY - (b-0.5)*scale) // Invert Y axis

			if px >= 0 && px < width && py >= 0 && py < height {
				img.Set(px, py, color.RGBA{
					uint8(clamp255(rgbR * 255)),
					uint8(clamp255(rgbG * 255)),
					uint8(clamp255(rgbB * 255)),
					255,
				})
			}
		}
	}

	// Draw axes
	drawLine(img, int(centerX-scale), int(centerY), int(centerX+scale), int(centerY), color.RGBA{255, 255, 255, 255}, 2)
	drawLine(img, int(centerX), int(centerY-scale), int(centerX), int(centerY+scale), color.RGBA{255, 255, 255, 255}, 2)

	return img
}

func generateHSLCylinder(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	centerX := float64(width) / 2
	centerY := float64(height) / 2
	radius := float64(width) * 0.35

	// Fill with HSL colors showing a circular slice (top-down view of cylinder)
	// Show hue around the circle, saturation as radius, at fixed lightness
	l := 0.5
	for s := 0.0; s <= 1.0; s += 0.005 {
		for h := 0.0; h < 1.0; h += 0.005 {
			c := col.NewHSL(h*360, s, l, 1.0)
			r, g, b, _ := c.RGBA()
			angle := h * 2 * math.Pi
			// Map to circular area centered in the image
			x := centerX + radius*s*math.Cos(angle)
			y := centerY + radius*s*math.Sin(angle) // Use sin for Y, not fixed height
			if x >= 0 && x < float64(width) && y >= 0 && y < float64(height) {
				// Check if point is within the circle
				dx := x - centerX
				dy := y - centerY
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist <= radius {
					img.Set(int(x), int(y), color.RGBA{
						uint8(clamp255(r * 255)),
						uint8(clamp255(g * 255)),
						uint8(clamp255(b * 255)),
						255,
					})
				}
			}
		}
	}

	return img
}

func generateLABSpace(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	centerX := float64(width) / 2
	centerY := float64(height) / 2
	scale := float64(width) * 0.3

	// Draw LAB space as a slice at L=50
	// A and B range from -128 to 127, but we'll use a more limited range that's closer to sRGB gamut
	l := 50.0
	step := 0.3
	for a := -80.0; a <= 80.0; a += step {
		for b := -80.0; b <= 80.0; b += step {
			// Convert LAB to RGB
			lab := col.NewLAB(l, a, b, 1.0)
			r, g, b, _ := lab.RGBA()

			// Map to image coordinates
			px := int(centerX + (a/100.0)*scale)
			py := int(centerY - (b/100.0)*scale) // Invert Y axis

			if px >= 0 && px < width && py >= 0 && py < height {
				// Always draw, even if out of gamut (it will be clamped)
				img.Set(px, py, color.RGBA{
					uint8(clamp255(r * 255)),
					uint8(clamp255(g * 255)),
					uint8(clamp255(b * 255)),
					255,
				})
			}
		}
	}

	// Draw axes
	drawLine(img, int(centerX-scale), int(centerY), int(centerX+scale), int(centerY), color.RGBA{255, 255, 255, 255}, 2)
	drawLine(img, int(centerX), int(centerY-scale), int(centerX), int(centerY+scale), color.RGBA{255, 255, 255, 255}, 2)

	return img
}

func generateOKLCHSpace(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	centerX := float64(width) / 2
	centerY := float64(height) / 2
	maxRadius := float64(width) * 0.35

	// Draw OKLCH space as a slice at L=0.5
	// C ranges from 0 to ~0.4, H ranges from 0 to 360
	// Use a more limited chroma range that stays within sRGB gamut
	l := 0.5
	for c := 0.0; c <= 0.3; c += 0.003 {
		for h := 0.0; h < 360.0; h += 0.5 {
			oklch := col.NewOKLCH(l, c, h, 1.0)
			r, g, b, _ := oklch.RGBA()

			// Only draw if color is valid and within sRGB gamut
			if r >= 0 && r <= 1 && g >= 0 && g <= 1 && b >= 0 && b <= 1 {
				// Convert polar to Cartesian
				angle := h * math.Pi / 180.0
				px := int(centerX + c*maxRadius*math.Cos(angle))
				py := int(centerY - c*maxRadius*math.Sin(angle)) // Invert Y axis

				if px >= 0 && px < width && py >= 0 && py < height {
					img.Set(px, py, color.RGBA{
						uint8(clamp255(r * 255)),
						uint8(clamp255(g * 255)),
						uint8(clamp255(b * 255)),
						255,
					})
				}
			}
		}
	}

	// Draw center point
	for dy := -3; dy <= 3; dy++ {
		for dx := -3; dx <= 3; dx++ {
			if dx*dx+dy*dy <= 9 {
				img.Set(int(centerX)+dx, int(centerY)+dy, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	return img
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

func drawQuad(img *image.RGBA, x1, y1, x2, y2, x3, y3, x4, y4 float64, col color.Color) {
	// Simple triangle fill for quad
	drawTriangle(img, x1, y1, x2, y2, x3, y3, col)
	drawTriangle(img, x1, y1, x3, y3, x4, y4, col)
}

func drawTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 float64, col color.Color) {
	// Find bounding box
	minX := int(math.Min(math.Min(x1, x2), x3))
	maxX := int(math.Max(math.Max(x1, x2), x3))
	minY := int(math.Min(math.Min(y1, y2), y3))
	maxY := int(math.Max(math.Max(y1, y2), y3))

	// Check each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInTriangle(float64(x), float64(y), x1, y1, x2, y2, x3, y3) {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, col)
				}
			}
		}
	}
}

func pointInTriangle(px, py, x1, y1, x2, y2, x3, y3 float64) bool {
	d1 := sign(px, py, x1, y1, x2, y2)
	d2 := sign(px, py, x2, y2, x3, y3)
	d3 := sign(px, py, x3, y3, x1, y1)
	return (d1 >= 0 && d2 >= 0 && d3 >= 0) || (d1 <= 0 && d2 <= 0 && d3 <= 0)
}

func sign(x1, y1, x2, y2, x3, y3 float64) float64 {
	return (x1-x3)*(y2-y3) - (x2-x3)*(y1-y3)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func clamp255(v float64) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return int(v)
}
