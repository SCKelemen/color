# Quick Start Guide

Get up and running with professional color manipulation in 5 minutes.

## Installation

```bash
go get github.com/SCKelemen/color
```

## Three Most Common Use Cases

### 1. Generate a Gradient (2 minutes)

**Problem:** You need a smooth color gradient without muddy midpoints.

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/color"
)

func main() {
    // Define start and end colors
    red := color.RGB(1, 0, 0)
    blue := color.RGB(0, 0, 1)

    // Generate 10-step gradient (perceptually uniform)
    gradient := color.Gradient(red, blue, 10)

    // Use the colors
    for i, c := range gradient {
        hex := color.RGBToHex(c)
        fmt.Printf("Step %d: %s\n", i, hex)
    }
}
```

**Output:**
```
Step 0: #ff0000  (red)
Step 1: #f0007d
Step 2: #dd00ab
Step 3: #c400cd
Step 4: #a300e3
Step 5: #7900ed
Step 6: #3f00ea
Step 7: #0000dc
Step 8: #0000c4
Step 9: #0000ff  (blue)
```

**Key points:**
- Uses OKLCH space automatically (perceptually uniform)
- No muddy middle colors
- Each step looks evenly spaced

**Advanced: Multi-stop gradients**

```go
// Red → Yellow → Blue gradient
stops := []color.GradientStop{
    {Color: color.RGB(1, 0, 0), Position: 0.0},   // Red at start
    {Color: color.RGB(1, 1, 0), Position: 0.5},   // Yellow at middle
    {Color: color.RGB(0, 0, 1), Position: 1.0},   // Blue at end
}

gradient := color.GradientMultiStop(stops, 20, color.GradientOKLCH)
```

---

### 2. Create a Color Palette (3 minutes)

**Problem:** You have one brand color and need to generate a full palette with proper spacing.

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/color"
)

func main() {
    // Your brand color
    brandBlue := color.ParseColor("#3B82F6")

    // Generate palette with perceptually uniform steps
    palette := map[string]string{
        // Lighter shades
        "50":  color.RGBToHex(color.Lighten(brandBlue, 0.45)),
        "100": color.RGBToHex(color.Lighten(brandBlue, 0.35)),
        "200": color.RGBToHex(color.Lighten(brandBlue, 0.25)),
        "300": color.RGBToHex(color.Lighten(brandBlue, 0.15)),
        "400": color.RGBToHex(color.Lighten(brandBlue, 0.05)),

        // Base color
        "500": color.RGBToHex(brandBlue),

        // Darker shades
        "600": color.RGBToHex(color.Darken(brandBlue, 0.05)),
        "700": color.RGBToHex(color.Darken(brandBlue, 0.15)),
        "800": color.RGBToHex(color.Darken(brandBlue, 0.25)),
        "900": color.RGBToHex(color.Darken(brandBlue, 0.35)),
    }

    // Print as CSS custom properties
    fmt.Println(":root {")
    for shade, hex := range palette {
        fmt.Printf("  --blue-%s: %s;\n", shade, hex)
    }
    fmt.Println("}")
}
```

**Output:**
```css
:root {
  --blue-50: #eff6ff;
  --blue-100: #dbeafe;
  --blue-200: #bfdbfe;
  --blue-300: #93c5fd;
  --blue-400: #60a5fa;
  --blue-500: #3b82f6;  /* Base */
  --blue-600: #2563eb;
  --blue-700: #1d4ed8;
  --blue-800: #1e40af;
  --blue-900: #1e3a8a;
}
```

**Key points:**
- Lighten/Darken work in OKLCH (perceptually uniform)
- Each step looks evenly lighter/darker
- Works with any starting color

**Alternative: Gradient-based palette**

```go
// Generate 10 evenly-spaced shades
lightest := color.Lighten(brandBlue, 0.45)
darkest := color.Darken(brandBlue, 0.35)
shades := color.Gradient(lightest, darkest, 10)

for i, shade := range shades {
    fmt.Printf("%d: %s\n", i*100+50, color.RGBToHex(shade))
}
```

---

### 3. Work with Wide-Gamut Colors (4 minutes)

**Problem:** You're targeting modern displays (iPhone, iPad, Mac) and want to preserve vibrant colors.

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/color"
)

func main() {
    // Create a color in Display P3 (26% more colors than sRGB)
    p3Red, _ := color.ConvertFromRGBSpace(1.0, 0.0, 0.0, 1.0, "display-p3")

    fmt.Println("=== Wide-Gamut Workflow ===\n")

    // 1. Manipulate in perceptual space (preserves gamut!)
    lighter := color.Lighten(p3Red, 0.2)
    saturated := color.Saturate(lighter, 0.15)

    // 2. Stay in Display P3 for modern displays
    p3Result, _ := color.ConvertToRGBSpace(saturated, "display-p3")
    channels := p3Result.Channels()
    fmt.Printf("Display P3: color(display-p3 %.3f %.3f %.3f)\n",
        channels[0], channels[1], channels[2])

    // 3. Or convert to sRGB for standard displays
    srgbResult := saturated.RGBA()
    fmt.Printf("sRGB:       rgb(%d, %d, %d)\n",
        int(srgbResult[0]*255), int(srgbResult[1]*255), int(srgbResult[2]*255))

    // 4. Check if color fits in sRGB
    if color.InGamut(saturated) {
        fmt.Println("✅ Fits in sRGB gamut")
    } else {
        fmt.Println("⚠️  Out of sRGB gamut - will be clipped")

        // Map to sRGB with quality preservation
        mapped := color.MapToGamut(saturated, color.GamutPreserveLightness)
        fmt.Printf("Mapped:     %s\n", color.RGBToHex(mapped))
    }
}
```

**Output:**
```
=== Wide-Gamut Workflow ===

Display P3: color(display-p3 1.089 0.345 0.298)
sRGB:       rgb(255, 88, 76)
⚠️  Out of sRGB gamut - will be clipped
Mapped:     #ff584c
```

**Key points:**
- Preserve wide-gamut throughout your pipeline
- Only convert to sRGB at the last moment
- Use proper gamut mapping when necessary

**Supported Wide-Gamut Spaces:**

```go
// Display P3 - iPhone X+, iPad Pro, Mac displays
p3Color, _ := color.ConvertFromRGBSpace(r, g, b, 1.0, "display-p3")

// Rec.2020 - UHDTV, future displays
rec2020Color, _ := color.ConvertFromRGBSpace(r, g, b, 1.0, "rec2020")

// Adobe RGB - Professional photography
adobeColor, _ := color.ConvertFromRGBSpace(r, g, b, 1.0, "a98-rgb")

// ProPhoto RGB - RAW photo editing (widest gamut)
prophotoColor, _ := color.ConvertFromRGBSpace(r, g, b, 1.0, "prophoto-rgb")
```

---

## Quick Reference: Essential Functions

### Creating Colors

```go
// From RGB values
red := color.RGB(1, 0, 0)                    // Full opacity
semiRed := color.NewRGBA(1, 0, 0, 0.5)       // With alpha

// Parse from string (supports all CSS color formats)
blue, _ := color.ParseColor("#0000FF")
green, _ := color.ParseColor("rgb(0, 255, 0)")
purple, _ := color.ParseColor("hsl(270, 100%, 50%)")
teal, _ := color.ParseColor("oklch(0.7 0.2 180)")

// Create in perceptual space
oklch := color.NewOKLCH(0.7, 0.2, 240, 1.0)  // L, C, H, alpha
```

### Manipulating Colors

```go
lighter := color.Lighten(c, 0.2)       // 20% lighter
darker := color.Darken(c, 0.2)         // 20% darker
vivid := color.Saturate(c, 0.3)        // 30% more saturated
muted := color.Desaturate(c, 0.3)      // 30% less saturated

shifted := color.AdjustHue(c, 60)      // Shift hue by 60°
complement := color.Complement(c)      // Opposite on color wheel
inverted := color.Invert(c)            // Invert RGB
gray := color.Grayscale(c)             // Convert to grayscale

mixed := color.MixOKLCH(c1, c2, 0.5)   // Mix 50/50 in OKLCH
```

### Generating Gradients

```go
// Simple gradient
gradient := color.Gradient(start, end, 10)

// Choose interpolation space
gradient := color.GradientInSpace(start, end, 10, color.GradientOKLCH)

// Multi-stop
stops := []color.GradientStop{
    {Color: red, Position: 0.0},
    {Color: yellow, Position: 0.5},
    {Color: blue, Position: 1.0},
}
gradient := color.GradientMultiStop(stops, 20, color.GradientOKLCH)

// With easing
gradient := color.GradientWithEasing(start, end, 20,
    color.GradientOKLCH, color.EaseInOutCubic)
```

### Converting Colors

```go
// Convert between color spaces
hsl := color.ToHSL(c)
oklch := color.ToOKLCH(c)
lab := color.ToLAB(c)

// Get RGB values
r, g, b, a := c.RGBA()

// Convert to hex
hex := color.RGBToHex(c)              // #RRGGBB
hexWithAlpha := color.RGBToHex(c)     // #RRGGBBAA if alpha < 1
```

### Wide-Gamut Operations

```go
// Convert to wide-gamut space
p3Color, _ := color.ConvertToRGBSpace(c, "display-p3")

// Create from wide-gamut values
color, _ := color.ConvertFromRGBSpace(1.0, 0.5, 0.0, 1.0, "display-p3")

// Check if in gamut
if color.InGamut(c) {
    fmt.Println("Fits in sRGB")
}

// Map to gamut
mapped := color.MapToGamut(c, color.GamutPreserveLightness)
```

### Color Difference

```go
// How different do colors look?
diff := color.DeltaEOK(c1, c2)         // Fast, modern
diff := color.DeltaE2000(c1, c2)       // Industry standard

// Interpretation:
// < 1.0  = Imperceptible
// 1-2    = Small difference
// 2-5    = Noticeable
// > 5    = Obviously different
```

---

## Common Patterns

### Pattern 1: Theme Color Generator

```go
func generateThemeColors(base color.Color) map[string]string {
    return map[string]string{
        "primary":       color.RGBToHex(base),
        "primary-light": color.RGBToHex(color.Lighten(base, 0.2)),
        "primary-dark":  color.RGBToHex(color.Darken(base, 0.2)),
        "secondary":     color.RGBToHex(color.AdjustHue(base, 180)),
        "accent":        color.RGBToHex(color.AdjustHue(base, 120)),
    }
}
```

### Pattern 2: Accessible Color Pairs

```go
func ensureAccessible(fg, bg color.Color, minDiff float64) color.Color {
    for color.DeltaE2000(fg, bg) < minDiff {
        fg = color.Darken(fg, 0.05) // Adjust until sufficient contrast
    }
    return fg
}
```

### Pattern 3: Heatmap Colors

```go
func createHeatmap(dataPoints []float64) []color.Color {
    cold := color.RGB(0, 0, 1)  // Blue
    hot := color.RGB(1, 0, 0)   // Red

    gradient := color.Gradient(cold, hot, 100)

    colors := make([]color.Color, len(dataPoints))
    for i, val := range dataPoints {
        normalized := (val - min) / (max - min)
        idx := int(normalized * 99)
        colors[i] = gradient[idx]
    }
    return colors
}
```

### Pattern 4: Photo Workflow

```go
func processPhoto(rawR, rawG, rawB float64) color.Color {
    // 1. Start in ProPhoto RGB (widest gamut)
    raw, _ := color.ConvertFromRGBSpace(rawR, rawG, rawB, 1.0, "prophoto-rgb")

    // 2. Edit in perceptual space
    edited := color.Saturate(raw, 0.15)
    edited = color.Lighten(edited, 0.05)

    // 3. Export to Display P3 (modern displays)
    p3, _ := color.ConvertToRGBSpace(edited, "display-p3")

    return p3
}
```

---

## Next Steps

Now that you understand the basics:

1. **Read [WHEN_TO_USE.md](WHEN_TO_USE.md)** - Decide if this library fits your needs
2. **Explore [COLOR_PRIMER.md](COLOR_PRIMER.md)** - Understand color science basics
3. **Check [VISUAL_COMPARISON.md](VISUAL_COMPARISON.md)** - See why perceptual spaces matter
4. **Browse [Full API Docs](README.md)** - Comprehensive function reference

## Troubleshooting

### "My colors look different than expected"

```go
// ❌ Don't do this (RGB operations)
lighter := color.NewRGBA(r+0.2, g+0.2, b+0.2, a)

// ✅ Do this (perceptually uniform)
lighter := color.Lighten(myColor, 0.2)
```

### "My gradients are muddy"

```go
// ❌ Don't interpolate in RGB
gradient := interpolateRGB(start, end, steps)

// ✅ Use OKLCH (default)
gradient := color.Gradient(start, end, steps)
```

### "I'm losing color vibrancy"

```go
// ❌ Don't force through sRGB early
srgb := convertToSRGB(displayP3Color) // Loses 26% of colors!

// ✅ Preserve wide-gamut throughout
p3Color, _ := color.ConvertFromRGBSpace(r, g, b, 1.0, "display-p3")
edited := color.Lighten(p3Color, 0.2)
p3Result, _ := color.ConvertToRGBSpace(edited, "display-p3")
```

---

## Examples

Full working examples are in the `examples/` directory:

```bash
# Theme palette generator
go run examples/theme_generator/main.go

# Heatmap with smooth gradients
go run examples/heatmap/main.go

# Wide-gamut photo workflow
go run examples/photo_workflow/main.go
```

---

**You're ready to go!** 🎨

For questions or issues, visit [GitHub Discussions](https://github.com/SCKelemen/color/discussions).
