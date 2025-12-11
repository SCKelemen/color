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

func generateXYZGamuts() error {
	outputDir := "docs/gamuts"
	os.MkdirAll(outputDir, 0755)

	// Define all RGB color spaces we support
	spaces := []struct {
		name      string
		colorName string // For color coding
		convert   func(r, g, b float64) *col.XYZ
	}{
		{
			name:      "sRGB",
			colorName: "sRGB",
			convert: func(r, g, b float64) *col.XYZ {
				// Convert sRGB to XYZ
				c := col.RGB(r, g, b)
				return col.ToXYZ(c)
			},
		},
		{
			name:      "Display P3",
			colorName: "DisplayP3",
			convert: func(r, g, b float64) *col.XYZ {
				// Convert Display P3 RGB to XYZ
				c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "display-p3")
				return col.ToXYZ(c)
			},
		},
		{
			name:      "Adobe RGB",
			colorName: "AdobeRGB",
			convert: func(r, g, b float64) *col.XYZ {
				// Convert Adobe RGB to XYZ
				c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "a98-rgb")
				return col.ToXYZ(c)
			},
		},
		{
			name:      "ProPhoto RGB",
			colorName: "ProPhotoRGB",
			convert: func(r, g, b float64) *col.XYZ {
				// Convert ProPhoto RGB to XYZ
				c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "prophoto-rgb")
				return col.ToXYZ(c)
			},
		},
		{
			name:      "Rec. 2020",
			colorName: "Rec2020",
			convert: func(r, g, b float64) *col.XYZ {
				// Convert Rec. 2020 RGB to XYZ
				c, _ := col.ConvertFromRGBSpace(r, g, b, 1.0, "rec2020")
				return col.ToXYZ(c)
			},
		},
	}

	// Generate both black and white text versions
	for _, textColor := range []string{"black", "white"} {
		img := generateXYZGamutComparison(spaces, textColor, 1200, 1000)
		suffix := ""
		if textColor == "black" {
			suffix = "_black"
		} else {
			suffix = "_white"
		}
		filename := filepath.Join(outputDir, fmt.Sprintf("gamut_xyz_comparison%s.png", suffix))
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		png.Encode(f, img)
		f.Close()
		fmt.Printf("Generated %s\n", filename)
	}

	return nil
}

func generateXYZGamutComparison(spaces []struct {
	name      string
	colorName string
	convert   func(r, g, b float64) *col.XYZ
}, textColor string, width, height int) *image.RGBA {
	scaledWidth := width * scale
	scaledHeight := height * scale
	img := image.NewRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))

	// Fill with transparent background
	for y := 0; y < scaledHeight; y++ {
		for x := 0; x < scaledWidth; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	// 3D rotation angles for isometric view
	angleY := 45.0 * math.Pi / 180.0
	angleX := 30.0 * math.Pi / 180.0
	angleZ := 0.0

	// Find the bounds of all gamuts in XYZ space
	minX, maxX := math.Inf(1), math.Inf(-1)
	minY, maxY := math.Inf(1), math.Inf(-1)
	minZ, maxZ := math.Inf(1), math.Inf(-1)

	// Sample all gamuts to find bounds
	step := 0.1
	for _, space := range spaces {
		for r := 0.0; r <= 1.0; r += step {
			for g := 0.0; g <= 1.0; g += step {
				for b := 0.0; b <= 1.0; b += step {
					xyz := space.convert(r, g, b)
					if xyz.X < minX {
						minX = xyz.X
					}
					if xyz.X > maxX {
						maxX = xyz.X
					}
					if xyz.Y < minY {
						minY = xyz.Y
					}
					if xyz.Y > maxY {
						maxY = xyz.Y
					}
					if xyz.Z < minZ {
						minZ = xyz.Z
					}
					if xyz.Z > maxZ {
						maxZ = xyz.Z
					}
				}
			}
		}
	}

	// Add some padding
	rangeX := maxX - minX
	rangeY := maxY - minY
	rangeZ := maxZ - minZ
	padding := 0.1
	minX -= rangeX * padding
	maxX += rangeX * padding
	minY -= rangeY * padding
	maxY += rangeY * padding
	minZ -= rangeZ * padding
	maxZ += rangeZ * padding

	// Calculate scale to fit in image
	scaleX := float64(scaledWidth) * 0.7 / rangeX
	scaleY := float64(scaledHeight) * 0.6 / rangeY
	scaleZ := math.Min(scaleX, scaleY) * 0.8

	centerX := float64(scaledWidth) / 2
	labelReserve := float64(scaledHeight) * 0.15
	centerY := (float64(scaledHeight) - labelReserve) / 2

	// Color coding for each gamut
	gamutColors := map[string]color.RGBA{
		"sRGB":         {255, 0, 0, 200},      // Red
		"DisplayP3":    {0, 255, 0, 200},      // Green
		"AdobeRGB":     {0, 0, 255, 200},      // Blue
		"ProPhotoRGB":  {255, 255, 0, 200},    // Yellow
		"Rec2020":      {255, 0, 255, 200},    // Magenta
	}

	// Draw each gamut as a wireframe
	for _, space := range spaces {
		gamutColor := gamutColors[space.colorName]
		if gamutColor.A == 0 {
			gamutColor = color.RGBA{200, 200, 200, 200} // Default gray
		}

		// Sample RGB cube edges and corners
		points := []struct{ r, g, b float64 }{
			// Corners
			{0, 0, 0}, {1, 0, 0}, {0, 1, 0}, {0, 0, 1},
			{1, 1, 0}, {1, 0, 1}, {0, 1, 1}, {1, 1, 1},
			// Edge midpoints
			{0.5, 0, 0}, {1, 0.5, 0}, {0.5, 1, 0}, {0, 0.5, 0},
			{0, 0, 0.5}, {1, 0, 0.5}, {0, 1, 0.5}, {1, 1, 0.5},
			{0.5, 0, 1}, {1, 0.5, 1}, {0.5, 1, 1}, {0, 0.5, 1},
		}

		// Convert to XYZ and project to 2D
		projected := make([]struct {
			x, y float64
			z    float64 // depth
		}, len(points))

		for i, p := range points {
			xyz := space.convert(p.r, p.g, p.b)
			// Apply 3D rotation
			x, y, z := rotate3D(xyz.X, xyz.Y, xyz.Z, angleY, angleX, angleZ)
			// Project to 2D (isometric)
			px := centerX + (x-minX)*scaleX
			py := centerY - (y-minY)*scaleY - (z-minZ)*scaleZ
			projected[i] = struct {
				x, y float64
				z    float64
			}{px, py, z}
		}

		// Draw edges of RGB cube in XYZ space
		edges := [][]int{
			{0, 1}, {0, 2}, {0, 3}, // From black corner
			{1, 4}, {1, 5},         // From red corner
			{2, 4}, {2, 6},         // From green corner
			{3, 5}, {3, 6},         // From blue corner
			{4, 7}, {5, 7}, {6, 7}, // To white corner
		}

		for _, edge := range edges {
			p1 := projected[edge[0]]
			p2 := projected[edge[1]]
			drawLine(img, int(p1.x), int(p1.y), int(p2.x), int(p2.y), gamutColor, 2)
		}

		// Draw some interior points to show the volume
		step := 0.2
		for r := 0.0; r <= 1.0; r += step {
			for g := 0.0; g <= 1.0; g += step {
				for b := 0.0; b <= 1.0; b += step {
					xyz := space.convert(r, g, b)
					x, y, z := rotate3D(xyz.X, xyz.Y, xyz.Z, angleY, angleX, angleZ)
					px := int(centerX + (x-minX)*scaleX)
					py := int(centerY - (y-minY)*scaleY - (z-minZ)*scaleZ)

					if px >= 0 && px < scaledWidth && py >= 0 && py < scaledHeight {
						// Draw a small point
						for dy := -1; dy <= 1; dy++ {
							for dx := -1; dx <= 1; dx++ {
								if dx*dx+dy*dy <= 1 {
									nx, ny := px+dx, py+dy
									if nx >= 0 && nx < scaledWidth && ny >= 0 && ny < scaledHeight {
										img.Set(nx, ny, gamutColor)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Draw axes
	axisLength := math.Max(rangeX, math.Max(rangeY, rangeZ)) * 0.3
	drawAxis3D(img, centerX, centerY, axisLength, angleY, angleX, angleZ, 1, 0, 0, color.RGBA{255, 100, 100, 255}, "X")
	drawAxis3D(img, centerX, centerY, axisLength, angleY, angleX, angleZ, 0, 1, 0, color.RGBA{100, 255, 100, 255}, "Y")
	drawAxis3D(img, centerX, centerY, axisLength, angleY, angleX, angleZ, 0, 0, 1, color.RGBA{100, 100, 255, 255}, "Z")

	// Draw legend
	legendY := int(float64(scaledHeight) - labelReserve*0.3)
	legendX := int(float64(scaledWidth) * 0.05)

	for i, space := range spaces {
		gamutColor := gamutColors[space.colorName]
		if gamutColor.A == 0 {
			gamutColor = color.RGBA{200, 200, 200, 255}
		}

		// Draw color square
		squareSize := int(20 * float64(scale))
		for y := 0; y < squareSize; y++ {
			for x := 0; x < squareSize; x++ {
				px := legendX + x
				py := legendY + y + i*int(30*float64(scale))
				if px >= 0 && px < scaledWidth && py >= 0 && py < scaledHeight {
					img.Set(px, py, gamutColor)
				}
			}
		}

		// Draw label
		labelX := legendX + squareSize + int(10*float64(scale))
		labelY := legendY + squareSize/2 + i*int(30*float64(scale))
		var shadowColor string
		if textColor == "white" {
			shadowColor = "black"
		}
		drawTextScaled(img, labelX, labelY, space.name, textColor, shadowColor, false, scale)
	}

	// Draw title
	titleY := int(float64(scaledHeight) * 0.05)
	var shadowColor string
	if textColor == "white" {
		shadowColor = "black"
	}
	drawTextScaled(img, scaledWidth/2, titleY, "RGB Gamuts in XYZ Color Space", textColor, shadowColor, true, scale)

	// Find bounding box and crop
	minXImg, minYImg, maxXImg, maxYImg := findBoundingBox(img)
	var padding int = 20 * scale
	croppedWidth := (maxXImg - minXImg) + (padding * 2)
	croppedHeight := (maxYImg - minYImg) + (padding * 2)

	croppedImg := image.NewRGBA(image.Rect(0, 0, croppedWidth, croppedHeight))
	for y := 0; y < croppedHeight; y++ {
		for x := 0; x < croppedWidth; x++ {
			croppedImg.Set(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	for y := minYImg; y < maxYImg; y++ {
		for x := minXImg; x < maxXImg; x++ {
			if x >= 0 && x < scaledWidth && y >= 0 && y < scaledHeight {
				c := img.RGBAAt(x, y)
				if c.A > 0 {
					croppedImg.Set(x-minXImg+padding, y-minYImg+padding, c)
				}
			}
		}
	}

	return scaleDown(croppedImg, croppedWidth/scale, croppedHeight/scale)
}

func rotate3D(x, y, z, angleY, angleX, angleZ float64) (float64, float64, float64) {
	// Rotate around Y axis
	x1 := x*math.Cos(angleY) + z*math.Sin(angleY)
	z1 := -x*math.Sin(angleY) + z*math.Cos(angleY)

	// Rotate around X axis
	y1 := y*math.Cos(angleX) - z1*math.Sin(angleX)
	z2 := y*math.Sin(angleX) + z1*math.Cos(angleX)

	// Rotate around Z axis
	x2 := x1*math.Cos(angleZ) - y1*math.Sin(angleZ)
	y2 := x1*math.Sin(angleZ) + y1*math.Cos(angleZ)

	return x2, y2, z2
}

