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

	// Draw RGB cube in isometric view
	// Cube vertices: (0,0,0) to (1,1,1)
	centerX := float64(width) / 2
	centerY := float64(height) / 2
	scale := float64(width) * 0.3

	// Isometric projection angles
	angle := math.Pi / 6 // 30 degrees

	// Project 3D point to 2D
	project := func(x, y, z float64) (float64, float64) {
		// Isometric projection
		px := (x - z) * math.Cos(angle) * scale
		py := (x + z) * math.Sin(angle) * scale - y*scale
		return centerX + px, centerY + py
	}

	// Draw cube edges
	edges := [][]float64{
		// Bottom face
		{0, 0, 0}, {1, 0, 0}, // R edge
		{0, 0, 0}, {0, 0, 1}, // B edge
		{1, 0, 0}, {1, 0, 1}, // R-B edge
		{0, 0, 1}, {1, 0, 1}, // B-R edge
		// Top face
		{0, 1, 0}, {1, 1, 0}, // R edge
		{0, 1, 0}, {0, 1, 1}, // B edge
		{1, 1, 0}, {1, 1, 1}, // R-B edge
		{0, 1, 1}, {1, 1, 1}, // B-R edge
		// Vertical edges
		{0, 0, 0}, {0, 1, 0}, // G edge
		{1, 0, 0}, {1, 1, 0}, // R-G edge
		{0, 0, 1}, {0, 1, 1}, // B-G edge
		{1, 0, 1}, {1, 1, 1}, // R-B-G edge
	}

	for i := 0; i < len(edges); i += 2 {
		x1, y1 := project(edges[i][0], edges[i][1], edges[i][2])
		x2, y2 := project(edges[i+1][0], edges[i+1][1], edges[i+1][2])
		drawLine(img, int(x1), int(y1), int(x2), int(y2), color.RGBA{255, 255, 255, 255}, 2)
	}

	// Draw colored faces (semi-transparent)
	// Front face (y=1)
	x1, y1 := project(0, 1, 0)
	x2, y2 := project(1, 1, 0)
	x3, y3 := project(1, 1, 1)
	x4, y4 := project(0, 1, 1)
	drawQuad(img, x1, y1, x2, y2, x3, y3, x4, y4, color.RGBA{255, 255, 255, 128})

	// Right face (r=1)
	x1, y1 = project(1, 0, 0)
	x2, y2 = project(1, 1, 0)
	x3, y3 = project(1, 1, 1)
	x4, y4 = project(1, 0, 1)
	drawQuad(img, x1, y1, x2, y2, x3, y3, x4, y4, color.RGBA{255, 0, 0, 128})

	// Top face (g=1)
	x1, y1 = project(0, 1, 0)
	x2, y2 = project(1, 1, 0)
	x3, y3 = project(1, 0, 0)
	x4, y4 = project(0, 0, 0)
	drawQuad(img, x1, y1, x2, y2, x3, y3, x4, y4, color.RGBA{0, 255, 0, 128})

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
	radius := float64(width) * 0.3
	heightScale := float64(height) * 0.4

	// Draw cylinder outline
	// Top circle
	for angle := 0.0; angle < 2*math.Pi; angle += 0.01 {
		x := centerX + radius*math.Cos(angle)
		y := centerY - heightScale/2
		if x >= 0 && x < float64(width) && y >= 0 && y < float64(height) {
			img.Set(int(x), int(y), color.RGBA{255, 255, 255, 255})
		}
	}

	// Bottom circle
	for angle := 0.0; angle < 2*math.Pi; angle += 0.01 {
		x := centerX + radius*math.Cos(angle)
		y := centerY + heightScale/2
		if x >= 0 && x < float64(width) && y >= 0 && y < float64(height) {
			img.Set(int(x), int(y), color.RGBA{255, 255, 255, 255})
		}
	}

	// Vertical lines
	for i := 0; i < 8; i++ {
		angle := float64(i) * 2 * math.Pi / 8
		x1 := centerX + radius*math.Cos(angle)
		y1 := centerY - heightScale/2
		x2 := centerX + radius*math.Cos(angle)
		y2 := centerY + heightScale/2
		drawLine(img, int(x1), int(y1), int(x2), int(y2), color.RGBA{255, 255, 255, 255}, 1)
	}

	// Fill with HSL colors (sliced at L=0.5)
	l := 0.5
	for s := 0.0; s <= 1.0; s += 0.02 {
		for h := 0.0; h < 1.0; h += 0.01 {
			c := col.NewHSL(h*360, s, l, 1.0)
			r, g, b, _ := c.RGBA()
			angle := h * 2 * math.Pi
			x := centerX + radius*s*math.Cos(angle)
			y := centerY - heightScale/2
			if x >= 0 && x < float64(width) && y >= 0 && y < float64(height) {
				img.Set(int(x), int(y), color.RGBA{
					uint8(clamp255(r * 255)),
					uint8(clamp255(g * 255)),
					uint8(clamp255(b * 255)),
					255,
				})
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
	// A and B range from -128 to 127, but we'll use normalized [-1, 1]
	l := 0.5
	for a := -1.0; a <= 1.0; a += 0.02 {
		for b := -1.0; b <= 1.0; b += 0.02 {
			// Convert LAB to RGB
			lab := col.NewLAB(l*100, a*100, b*100, 1.0)
			r, g, b, _ := lab.RGBA()
			
			// Map to image coordinates
			px := int(centerX + a*scale)
			py := int(centerY - b*scale) // Invert Y axis
			
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

	// Draw axes
	drawLine(img, int(centerX-scale), int(centerY), int(centerX+scale), int(centerY), color.RGBA{255, 255, 255, 255}, 1)
	drawLine(img, int(centerX), int(centerY-scale), int(centerX), int(centerY+scale), color.RGBA{255, 255, 255, 255}, 1)

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
	l := 0.5
	for c := 0.0; c <= 0.4; c += 0.005 {
		for h := 0.0; h < 360.0; h += 1.0 {
			oklch := col.NewOKLCH(l, c, h, 1.0)
			r, g, b, _ := oklch.RGBA()
			
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

