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
	// We need to access the RGBColorSpace structs directly to use their conversion matrices
	spaces := []struct {
		name      string
		colorName string // For color coding
		convert   func(r, g, b float64) *col.XYZ
	}{
		{
			name:      "sRGB",
			colorName: "sRGB",
			convert: func(r, g, b float64) *col.XYZ {
				// Convert sRGB to XYZ using ToXYZ (which handles sRGB correctly)
				c := col.RGB(r, g, b)
				return col.ToXYZ(c)
			},
		},
		{
			name:      "Display P3",
			colorName: "DisplayP3",
			convert: func(r, g, b float64) *col.XYZ {
				return convertRGBSpaceToXYZ(r, g, b, "display-p3")
			},
		},
		{
			name:      "Adobe RGB",
			colorName: "AdobeRGB",
			convert: func(r, g, b float64) *col.XYZ {
				return convertRGBSpaceToXYZ(r, g, b, "a98-rgb")
			},
		},
		{
			name:      "ProPhoto RGB",
			colorName: "ProPhotoRGB",
			convert: func(r, g, b float64) *col.XYZ {
				return convertRGBSpaceToXYZ(r, g, b, "prophoto-rgb")
			},
		},
		{
			name:      "Rec. 2020",
			colorName: "Rec2020",
			convert: func(r, g, b float64) *col.XYZ {
				return convertRGBSpaceToXYZ(r, g, b, "rec2020")
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

	// Recalculate ranges after padding
	rangeX = maxX - minX
	rangeY = maxY - minY
	rangeZ = maxZ - minZ

	// Calculate scale to fit in image
	// Use a uniform scale based on the largest dimension to maintain aspect ratio
	maxRange := math.Max(rangeX, math.Max(rangeY, rangeZ))
	uniformScale := math.Min(float64(scaledWidth)*0.7, float64(scaledHeight)*0.6) / maxRange

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
			// First scale to normalized coordinates centered at origin
			xNorm := (xyz.X - (minX+maxX)/2) * uniformScale
			yNorm := (xyz.Y - (minY+maxY)/2) * uniformScale
			zNorm := (xyz.Z - (minZ+maxZ)/2) * uniformScale
			// Then apply 3D rotation
			xRot, yRot, zRot := rotate3D(xNorm, yNorm, zNorm, angleY, angleX, angleZ)
			// Project to 2D (isometric)
			px := centerX + xRot
			py := centerY - yRot - zRot*0.5 // Isometric projection
			projected[i] = struct {
				x, y float64
				z    float64
			}{px, py, zRot}
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
			// Draw much thicker lines for better visibility
			drawLine(img, int(p1.x), int(p1.y), int(p2.x), int(p2.y), gamutColor, 5)
		}

		// Draw filled volume by sampling more densely
		step := 0.05 // Increased density
		for r := 0.0; r <= 1.0; r += step {
			for g := 0.0; g <= 1.0; g += step {
				for b := 0.0; b <= 1.0; b += step {
					xyz := space.convert(r, g, b)
					// First scale to normalized coordinates centered at origin
					xNorm := (xyz.X - (minX+maxX)/2) * uniformScale
					yNorm := (xyz.Y - (minY+maxY)/2) * uniformScale
					zNorm := (xyz.Z - (minZ+maxZ)/2) * uniformScale
					// Then apply 3D rotation
					xRot, yRot, zRot := rotate3D(xNorm, yNorm, zNorm, angleY, angleX, angleZ)
					// Project to 2D (isometric)
					px := int(centerX + xRot)
					py := int(centerY - yRot - zRot*0.5)

					if px >= 0 && px < scaledWidth && py >= 0 && py < scaledHeight {
						// Draw a larger point for better visibility
						pointSize := 2
						for dy := -pointSize; dy <= pointSize; dy++ {
							for dx := -pointSize; dx <= pointSize; dx++ {
								if dx*dx+dy*dy <= pointSize*pointSize {
									nx, ny := px+dx, py+dy
									if nx >= 0 && nx < scaledWidth && ny >= 0 && ny < scaledHeight {
										// Blend with existing color for transparency effect
										existing := img.RGBAAt(nx, ny)
										if existing.A == 0 {
											img.Set(nx, ny, gamutColor)
										} else {
											// Blend colors
											alpha := float64(gamutColor.A) / 255.0
											r := uint8(float64(existing.R)*(1-alpha) + float64(gamutColor.R)*alpha)
											g := uint8(float64(existing.G)*(1-alpha) + float64(gamutColor.G)*alpha)
											b := uint8(float64(existing.B)*(1-alpha) + float64(gamutColor.B)*alpha)
											a := uint8(math.Max(float64(existing.A), float64(gamutColor.A)))
											img.Set(nx, ny, color.RGBA{r, g, b, a})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Draw axes - use normalized coordinates
	axisLength := maxRange * 0.3 * uniformScale
	// X axis (red)
	xNorm := axisLength
	yNorm := 0.0
	zNorm := 0.0
	xRot, yRot, zRot := rotate3D(xNorm, yNorm, zNorm, angleY, angleX, angleZ)
	drawLine(img, int(centerX), int(centerY), int(centerX+xRot), int(centerY-yRot-zRot*0.5), color.RGBA{255, 100, 100, 255}, 3)
	// Y axis (green)
	xNorm = 0.0
	yNorm = axisLength
	zNorm = 0.0
	xRot, yRot, zRot = rotate3D(xNorm, yNorm, zNorm, angleY, angleX, angleZ)
	drawLine(img, int(centerX), int(centerY), int(centerX+xRot), int(centerY-yRot-zRot*0.5), color.RGBA{100, 255, 100, 255}, 3)
	// Z axis (blue)
	xNorm = 0.0
	yNorm = 0.0
	zNorm = axisLength
	xRot, yRot, zRot = rotate3D(xNorm, yNorm, zNorm, angleY, angleX, angleZ)
	drawLine(img, int(centerX), int(centerY), int(centerX+xRot), int(centerY-yRot-zRot*0.5), color.RGBA{100, 100, 255, 255}, 3)

	// Draw legend - position it higher to avoid being cut off
	legendY := int(float64(scaledHeight) - labelReserve*0.5)
	legendX := int(float64(scaledWidth) * 0.05)
	legendSpacing := int(35 * float64(scale)) // Increased spacing between legend items

	for i, space := range spaces {
		gamutColor := gamutColors[space.colorName]
		if gamutColor.A == 0 {
			gamutColor = color.RGBA{200, 200, 200, 255}
		}

		// Draw color square
		squareSize := int(20 * float64(scale))
		squareY := legendY + i*legendSpacing
		for y := 0; y < squareSize; y++ {
			for x := 0; x < squareSize; x++ {
				px := legendX + x
				py := squareY + y
				if px >= 0 && px < scaledWidth && py >= 0 && py < scaledHeight {
					img.Set(px, py, gamutColor)
				}
			}
		}

		// Draw label
		labelX := legendX + squareSize + int(10*float64(scale))
		labelY := squareY + squareSize/2
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
	// Make sure to include the legend area in the bounding box
	minXImg, minYImg, maxXImg, maxYImg := findBoundingBox(img)
	
	// Ensure legend is included - check if legend area extends beyond current bounds
	legendX := int(float64(scaledWidth) * 0.05)
	legendStartY := int(float64(scaledHeight) - labelReserve*0.5)
	legendEndY := legendStartY + len(spaces)*int(35*float64(scale))
	
	// Expand bounding box to include legend if needed
	if legendX < minXImg {
		minXImg = legendX
	}
	// Estimate legend width (square + text)
	estimatedLegendWidth := int(20*float64(scale)) + int(10*float64(scale)) + 200*scale // square + spacing + text estimate
	if legendX+estimatedLegendWidth > maxXImg {
		maxXImg = legendX + estimatedLegendWidth
	}
	if legendStartY < minYImg {
		minYImg = legendStartY
	}
	if legendEndY > maxYImg {
		maxYImg = legendEndY
	}
	
	imgPadding := 20 * scale
	croppedWidth := (maxXImg - minXImg) + (imgPadding * 2)
	croppedHeight := (maxYImg - minYImg) + (imgPadding * 2)

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
					croppedImg.Set(x-minXImg+imgPadding, y-minYImg+imgPadding, c)
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

// convertRGBSpaceToXYZ converts RGB values in a specific color space directly to XYZ
// This avoids the double conversion through sRGB that ConvertFromRGBSpace does
func convertRGBSpaceToXYZ(r, g, b float64, spaceName string) *col.XYZ {
	// We'll manually implement the conversion using the matrices from rgb_spaces.go
	// Since we can't access the private RGBColorSpace structs, we'll use ConvertFromRGBSpace
	// which does RGB->XYZ->sRGB, then convert the sRGB result back to XYZ
	// Actually, a better approach: create an RGB color in the space, convert to XYZ
	
	// The issue is ConvertFromRGBSpace converts RGB->XYZ->sRGB, losing the original XYZ
	// We need to access the RGBColorSpace's ConvertRGBToXYZ directly
	// Since it's not exported, let's use a workaround: create a temporary color and extract XYZ
	
	// Actually, the simplest: use ConvertFromRGBSpace to get the color, then convert back to XYZ
	// But that's lossy. Instead, let's manually implement the conversion using known matrices.
	
	// For now, let's use the conversion matrices directly (hardcoded from rgb_spaces.go)
	switch spaceName {
	case "display-p3":
		// Display P3: apply inverse sRGB transfer, then matrix
		linearR := inverseSRGBTransfer(r)
		linearG := inverseSRGBTransfer(g)
		linearB := inverseSRGBTransfer(b)
		// Display P3 RGBToXYZMatrix
		x := 0.4865709486482162*linearR + 0.26566769316909306*linearG + 0.1982172852343625*linearB
		y := 0.2289745640697488*linearR + 0.6917385218365064*linearG + 0.079286914093745*linearB
		z := 0.000000000000000*linearR + 0.04511338185890264*linearG + 1.043944368900976*linearB
		return col.NewXYZ(x, y, z, 1.0)
	case "a98-rgb":
		// Adobe RGB: gamma 2.2, then matrix
		linearR := math.Pow(r, 2.2)
		linearG := math.Pow(g, 2.2)
		linearB := math.Pow(b, 2.2)
		// Adobe RGB RGBToXYZMatrix
		x := 0.5766690429101305*linearR + 0.1855582379065463*linearG + 0.1882286462349947*linearB
		y := 0.29734497525053605*linearR + 0.6273635662554661*linearG + 0.07529145849399788*linearB
		z := 0.02703136138641234*linearR + 0.07068885253582723*linearG + 0.9913375368376388*linearB
		return col.NewXYZ(x, y, z, 1.0)
	case "prophoto-rgb":
		// ProPhoto RGB: gamma 1.8, then matrix
		linearR := math.Pow(r, 1.8)
		linearG := math.Pow(g, 1.8)
		linearB := math.Pow(b, 1.8)
		// ProPhoto RGB RGBToXYZMatrix
		x := 0.7976749*linearR + 0.1351917*linearG + 0.0313534*linearB
		y := 0.2880402*linearR + 0.7118741*linearG + 0.0000857*linearB
		z := 0.0000000*linearR + 0.0000000*linearG + 0.8252100*linearB
		return col.NewXYZ(x, y, z, 1.0)
	case "rec2020":
		// Rec. 2020: gamma 2.4, then matrix
		linearR := math.Pow(r, 2.4)
		linearG := math.Pow(g, 2.4)
		linearB := math.Pow(b, 2.4)
		// Rec. 2020 RGBToXYZMatrix
		x := 0.6369580483012914*linearR + 0.14461690358620832*linearG + 0.1688809751641721*linearB
		y := 0.262704531669281*linearR + 0.6779980715188708*linearG + 0.05930171646986196*linearB
		z := 0.000000000000000*linearR + 0.028072693049087428*linearG + 1.060985057710791*linearB
		return col.NewXYZ(x, y, z, 1.0)
	default:
		// Fallback to sRGB
		c := col.RGB(r, g, b)
		return col.ToXYZ(c)
	}
}

func inverseSRGBTransfer(encoded float64) float64 {
	if encoded <= 0.04045 {
		return encoded / 12.92
	}
	return math.Pow((encoded+0.055)/1.055, 2.4)
}

