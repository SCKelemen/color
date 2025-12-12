package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	col "github.com/SCKelemen/color"
)

func main() {
	outputDir := "docs/chromaticity"
	os.MkdirAll(outputDir, 0755)

	// Generate chromaticity diagrams for different RGB spaces
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
		img := generateChromaticityDiagram(s.fn, 1000, 1000)
		filename := filepath.Join(outputDir, fmt.Sprintf("chromaticity_%s.png", s.name))
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		png.Encode(f, img)
		f.Close()
		fmt.Printf("Generated %s\n", filename)
	}
}

func generateChromaticityDiagram(createColor func(r, g, b float64) col.Color, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	// Map chromaticity coordinates to image coordinates
	// CIE xy chromaticity: x in [0, 0.8], y in [0, 0.9] approximately
	// We'll use a larger range to show the full diagram
	minX, maxX := 0.0, 0.8
	minY, maxY := 0.0, 0.9
	margin := 50.0

	scaleX := (float64(width) - 2*margin) / (maxX - minX)
	scaleY := (float64(height) - 2*margin) / (maxY - minY)

	// First, fill the visible gamut area with colors
	fillChromaticityColors(img, margin, margin, scaleX, scaleY, minX, minY, maxY, createColor)

	// Draw spectral locus (horseshoe shape) on top
	drawSpectralLocus(img, margin, margin, scaleX, scaleY, minX, minY, maxY)

	// Draw gamut boundary for this RGB space
	drawGamutBoundary(img, margin, margin, scaleX, scaleY, minX, minY, maxY, createColor)

	// Draw white point (D65)
	drawWhitePoint(img, margin, margin, scaleX, scaleY, minX, minY, maxY)

	return img
}

func fillChromaticityColors(img *image.RGBA, offsetX, offsetY, scaleX, scaleY, minX, minY, maxY float64, createColor func(r, g, b float64) col.Color) {
	// Fill the visible gamut area by sampling RGB colors and mapping them to xy coordinates
	// Create a lookup map from xy coordinates to RGB colors
	xyToColor := make(map[[2]int]color.RGBA)

	// Sample RGB space densely
	step := 0.02
	for r := 0.0; r <= 1.0; r += step {
		for g := 0.0; g <= 1.0; g += step {
			for b := 0.0; b <= 1.0; b += step {
				c := createColor(r, g, b)
				xyz := col.ToXYZ(c)
				sum := xyz.X + xyz.Y + xyz.Z
				if sum > 0.001 {
					x := xyz.X / sum
					y := xyz.Y / sum

					// Map to pixel coordinates
					px := int(offsetX + scaleX*(x-minX))
					py := int(offsetY + scaleY*(maxY-y))

					if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
						// Store the color for this pixel
						r8, g8, b8, _ := c.RGBA()
						key := [2]int{px, py}
						// Use the most saturated color for each pixel (or average)
						if existing, ok := xyToColor[key]; !ok {
							xyToColor[key] = color.RGBA{
								uint8(clamp255(r8 * 255)),
								uint8(clamp255(g8 * 255)),
								uint8(clamp255(b8 * 255)),
								255,
							}
						} else {
							// Average with existing (for smoother result)
							xyToColor[key] = color.RGBA{
								uint8((int(existing.R) + int(clamp255(r8*255))) / 2),
								uint8((int(existing.G) + int(clamp255(g8*255))) / 2),
								uint8((int(existing.B) + int(clamp255(b8*255))) / 2),
								255,
							}
						}
					}
				}
			}
		}
	}

	// Fill pixels with colors, using flood fill for smooth coverage
	for key, col := range xyToColor {
		img.Set(key[0], key[1], col)
	}

	// Fill gaps by interpolating nearby colors
	for py := 0; py < img.Bounds().Dy(); py++ {
		for px := 0; px < img.Bounds().Dx(); px++ {
			if img.RGBAAt(px, py).A == 0 {
				// Find nearest colored pixel
				nearest := findNearestColor(img, px, py, 10)
				if nearest.A > 0 {
					img.Set(px, py, nearest)
				}
			}
		}
	}
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

func absFloat(x float64) float64 {
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

func drawSpectralLocus(img *image.RGBA, offsetX, offsetY, scaleX, scaleY, minX, minY, maxY float64) {
	// Spectral locus: wavelengths from 380nm to 780nm
	// Approximate xy coordinates for spectral colors
	xyCoords := []struct{ x, y float64 }{
		{0.1741, 0.0050}, // 380nm
		{0.1738, 0.0049}, // 400nm
		{0.1510, 0.0480}, // 450nm
		{0.0082, 0.5384}, // 500nm
		{0.3135, 0.6900}, // 550nm
		{0.6270, 0.3725}, // 600nm
		{0.7347, 0.2653}, // 650nm
		{0.7347, 0.2653}, // 700nm
		{0.7347, 0.2653}, // 780nm
	}

	// Draw line of purples (from 380nm to 780nm)
	drawLine(img, int(offsetX+scaleX*(xyCoords[0].x-minX)), int(offsetY+scaleY*(maxY-xyCoords[0].y)),
		int(offsetX+scaleX*(xyCoords[len(xyCoords)-1].x-minX)), int(offsetY+scaleY*(maxY-xyCoords[len(xyCoords)-1].y)),
		color.RGBA{100, 100, 100, 255}, 1)

	// Draw spectral locus curve
	for i := 0; i < len(xyCoords)-1; i++ {
		x1 := int(offsetX + scaleX*(xyCoords[i].x-minX))
		y1 := int(offsetY + scaleY*(maxY-xyCoords[i].y))
		x2 := int(offsetX + scaleX*(xyCoords[i+1].x-minX))
		y2 := int(offsetY + scaleY*(maxY-xyCoords[i+1].y))
		drawLine(img, x1, y1, x2, y2, color.RGBA{150, 150, 150, 255}, 2)
	}
}

func drawGamutBoundary(img *image.RGBA, offsetX, offsetY, scaleX, scaleY, minX, minY, maxY float64, createColor func(r, g, b float64) col.Color) {
	// Sample RGB cube edges to find gamut boundary
	// We'll sample the RGB primaries and their combinations
	primaries := []struct {
		r, g, b float64
		name    string
	}{
		{1, 0, 0, "R"},
		{0, 1, 0, "G"},
		{0, 0, 1, "B"},
		{1, 1, 0, "Y"},
		{0, 1, 1, "C"},
		{1, 0, 1, "M"},
		{1, 1, 1, "W"},
		{0, 0, 0, "K"},
	}

	points := make([]struct{ x, y float64 }, 0)

	for _, p := range primaries {
		c := createColor(p.r, p.g, p.b)
		xyz := col.ToXYZ(c)
		sum := xyz.X + xyz.Y + xyz.Z
		if sum > 0 {
			x := xyz.X / sum
			y := xyz.Y / sum
			points = append(points, struct{ x, y float64 }{x, y})
		}
	}

	// Draw gamut triangle (R-G-B)
	if len(points) >= 3 {
		// Red
		x1 := int(offsetX + scaleX*(points[0].x-minX))
		y1 := int(offsetY + scaleY*(maxY-points[0].y))
		// Green
		x2 := int(offsetX + scaleX*(points[1].x-minX))
		y2 := int(offsetY + scaleY*(maxY-points[1].y))
		// Blue
		x3 := int(offsetX + scaleX*(points[2].x-minX))
		y3 := int(offsetY + scaleY*(maxY-points[2].y))

		// Draw triangle edges
		drawLine(img, x1, y1, x2, y2, color.RGBA{255, 255, 255, 255}, 2)
		drawLine(img, x2, y2, x3, y3, color.RGBA{255, 255, 255, 255}, 2)
		drawLine(img, x3, y3, x1, y1, color.RGBA{255, 255, 255, 255}, 2)
	}
}

func drawWhitePoint(img *image.RGBA, offsetX, offsetY, scaleX, scaleY, minX, minY, maxY float64) {
	// D65 white point: x=0.3127, y=0.3290
	wpX, wpY := 0.3127, 0.3290
	x := int(offsetX + scaleX*(wpX-minX))
	y := int(offsetY + scaleY*(maxY-wpY))

	// Draw white point as a small circle
	radius := 5
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < img.Bounds().Dx() && ny >= 0 && ny < img.Bounds().Dy() {
					img.Set(nx, ny, color.RGBA{255, 255, 255, 255})
				}
			}
		}
	}
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
