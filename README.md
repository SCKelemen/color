# Color - Professional Color Manipulation for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/SCKelemen/color.svg)](https://pkg.go.dev/github.com/SCKelemen/color)
[![Go Report Card](https://goreportcard.com/badge/github.com/SCKelemen/color)](https://goreportcard.com/report/github.com/SCKelemen/color)

A comprehensive Go library for **perceptually uniform color manipulation** with full support for modern wide-gamut color spaces and professional LOG formats for cinema cameras.

## Why This Library?

Most color libraries treat RGB as the only color space, causing three fundamental problems:

| Problem | Standard Libraries | This Library |
|---------|-------------------|--------------|
| **Perceptual Operations** | Lightening blue → cyan | Lightening blue → lighter blue ✅ |
| **Gradients** | Muddy midpoints, uneven steps | Perceptually smooth ✅ |
| **Wide Gamut** | Force through sRGB, lose vibrancy | Preserve Display P3/Rec.2020 ✅ |
| **Color Science** | Limited or incorrect | Industry-standard algorithms ✅ |

```go
// The difference:

// ❌ Standard RGB manipulation
lighter := RGB{blue.R + 0.2, blue.G + 0.2, blue.B + 0.2} // Looks cyan-ish

// ✅ This library (perceptually uniform)
lighter := color.Lighten(blue, 0.2) // Actually looks lighter blue!
```

## Key Features

### 🎨 Perceptually Uniform Operations
Operations work in OKLCH space where "20% lighter" actually **looks** 20% lighter to human eyes.

```go
blue := color.RGB(0, 0, 1)

// All operations are perceptually uniform:
lighter := color.Lighten(blue, 0.2)    // Looks evenly lighter
darker := color.Darken(blue, 0.2)      // Looks evenly darker
vivid := color.Saturate(blue, 0.3)     // Looks more saturated
muted := color.Desaturate(blue, 0.3)   // Looks less saturated
```

### 🌈 Smooth Gradients
Generate gradients that actually look smooth, not muddy.

<table>
<tr>
<td width="50%">

**❌ RGB Gradient** (muddy middle)
<img src="docs/gradients/gradient_rgb_black.png" alt="RGB gradient with muddy middle" />

</td>
<td width="50%">

**✅ OKLCH Gradient** (smooth & uniform)
<img src="docs/gradients/gradient_oklch_black.png" alt="OKLCH gradient smooth" />

</td>
</tr>
</table>

```go
red := color.RGB(1, 0, 0)
blue := color.RGB(0, 0, 1)

// Smooth, perceptually uniform gradient
gradient := color.Gradient(red, blue, 20)
```

### 🖥️ Wide-Gamut Color Spaces

Full support for modern display color spaces:

- **Display P3** (26% more colors than sRGB) - iPhone X+, iPad Pro, Mac displays
- **DCI-P3** - Digital cinema
- **Adobe RGB** (44% more colors) - Professional photography
- **Rec.2020** (73% more colors) - UHDTV, HDR
- **ProPhoto RGB** (189% more colors) - RAW photo editing

```go
// Create in Display P3 (wider gamut)
p3Red, _ := color.ConvertFromRGBSpace(1, 0, 0, 1, "display-p3")

// Manipulate in perceptual space
lighter := color.Lighten(p3Red, 0.2)

// Convert back to Display P3 (preserves vibrancy!)
result, _ := color.ConvertToRGBSpace(lighter, "display-p3")
```

### 🎯 Accurate Color Matching

Industry-standard color difference metrics:

```go
color1 := color.ParseColor("#FF6B6B")
color2 := color.ParseColor("#FF6D6C")

// How different do they look to humans?
diff := color.DeltaE2000(color1, color2) // Industry standard

if diff < 1.0 {
    fmt.Println("Imperceptible difference")
} else if diff < 2.0 {
    fmt.Println("Small difference")
} else {
    fmt.Println("Noticeable difference")
}
```

### 🔧 Professional Gamut Mapping

When converting from wide to narrow gamuts, choose your mapping strategy:

```go
// Vivid Display P3 color (out of sRGB gamut)
vividColor := NewSpaceColor(DisplayP3Space, []float64{1.1, 0.3, 0.2}, 1.0)

// Choose your mapping strategy:
clipped := MapToGamut(vividColor, GamutClip)                // Fast
lightness := MapToGamut(vividColor, GamutPreserveLightness) // Keep brightness ⭐
chroma := MapToGamut(vividColor, GamutPreserveChroma)       // Keep saturation
best := MapToGamut(vividColor, GamutProject)                // Best quality
```

### 🎬 Cinema Camera LOG Support

Professional LOG color spaces for cinema camera workflows with HDR support:

```go
// Load LOG footage from cinema camera
slog3 := color.NewSpaceColor(color.SLog3Space,
    []float64{0.41, 0.39, 0.35}, 1.0)

// Process in HDR workflow
hdr := slog3.ConvertTo(color.Rec2020Space)      // HDR mastering
web := slog3.ConvertTo(color.SRGBSpace)         // Web delivery

// Supported: Canon C-Log, Sony S-Log3, Panasonic V-Log,
//            Arri LogC, Red Log3G10, Blackmagic Film
```

## Quick Start

### Installation

```bash
go get github.com/SCKelemen/color
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/color"
)

func main() {
    // Parse any CSS color format
    blue, _ := color.ParseColor("#0000FF")

    // Perceptually uniform operations
    lighter := color.Lighten(blue, 0.2)
    darker := color.Darken(blue, 0.2)
    vivid := color.Saturate(blue, 0.3)

    // Generate smooth gradient
    red := color.RGB(1, 0, 0)
    gradient := color.Gradient(red, blue, 10)

    // Convert to hex
    for _, c := range gradient {
        fmt.Println(color.RGBToHex(c))
    }
}
```

## Comprehensive Color Space Support

### RGB Color Spaces
- **sRGB** - Standard web colors
- **sRGB-linear** - Linear RGB for correct blending
- **Display P3** - Modern Apple devices, HDR displays
- **DCI-P3** - Digital cinema
- **Adobe RGB 1998** - Professional photography
- **ProPhoto RGB** - RAW photo editing (widest gamut)
- **Rec.2020** - UHDTV, future displays
- **Rec.709** - HDTV

### LOG Color Spaces (Cinema Cameras)
Professional logarithmic color spaces for cinema camera workflows with HDR support:

- **Canon C-Log** - Cinema EOS cameras (C300, C500, etc.) with Cinema Gamut
- **Sony S-Log3** - Sony cinema cameras (FX6, FX9, Venice, etc.) with S-Gamut3
- **Panasonic V-Log** - Panasonic cameras (GH5, S1H, EVA1, etc.) with V-Gamut
- **Arri LogC** - Arri cameras (Alexa, Amira) with Arri Wide Gamut
- **Red Log3G10** - Red cameras (Komodo, V-Raptor) with RedWideGamutRGB
- **Blackmagic Film** - Blackmagic cameras (Pocket, URSA) with wide gamut

```go
// Load S-Log3 footage from Sony camera
slog3 := color.NewSpaceColor(color.SLog3Space,
    []float64{0.41, 0.39, 0.35}, 1.0)

// Convert to Rec.2020 for HDR delivery
hdr := slog3.ConvertTo(color.Rec2020Space)

// Convert to sRGB for web
web := slog3.ConvertTo(color.SRGBSpace)
```

### Perceptually Uniform Spaces
- **OKLCH** ⭐ - Modern, recommended (cylindrical)
- **OKLAB** - Modern, recommended (rectangular)
- **CIELAB** - Industry standard (rectangular)
- **CIELCH** - Industry standard (cylindrical)
- **CIELUV** - For emissive displays
- **LCHuv** - CIELUV cylindrical

### Intuitive Spaces
- **HSL** - Hue, Saturation, Lightness
- **HSV/HSB** - Hue, Saturation, Value
- **HWB** - Hue, Whiteness, Blackness (CSS Level 4)

### Reference Space
- **XYZ** - CIE 1931 (conversion hub)

## Complete Feature Set

### Color Creation & Parsing
```go
// Create colors
red := color.RGB(1, 0, 0)
oklch := color.NewOKLCH(0.7, 0.2, 30, 1.0)

// Parse CSS colors
parsed, _ := color.ParseColor("#FF0000")
parsed, _ := color.ParseColor("rgb(255, 0, 0)")
parsed, _ := color.ParseColor("hsl(0, 100%, 50%)")
parsed, _ := color.ParseColor("oklch(0.7 0.2 30)")
parsed, _ := color.ParseColor("color(display-p3 1 0 0)")
```

### Color Manipulation
```go
// Perceptually uniform
lighter := color.Lighten(c, 0.2)
darker := color.Darken(c, 0.2)
vivid := color.Saturate(c, 0.3)
muted := color.Desaturate(c, 0.3)

// Hue operations
shifted := color.AdjustHue(c, 60)
complement := color.Complement(c)

// Mixing
mixed := color.MixOKLCH(c1, c2, 0.5)

// Other
inverted := color.Invert(c)
gray := color.Grayscale(c)
```

### Gradients
```go
// Simple gradient (OKLCH, perceptually uniform)
gradient := color.Gradient(start, end, 20)

// Choose interpolation space
gradient := color.GradientInSpace(start, end, 20, color.GradientOKLCH)

// Multi-stop gradient
stops := []color.GradientStop{
    {Color: red, Position: 0.0},
    {Color: yellow, Position: 0.5},
    {Color: blue, Position: 1.0},
}
gradient := color.GradientMultiStop(stops, 30, color.GradientOKLCH)

// With easing
gradient := color.GradientWithEasing(start, end, 20,
    color.GradientOKLCH, color.EaseInOutCubic)
```

### Color Difference
```go
// Modern, fast
diff := color.DeltaEOK(c1, c2)

// Industry standard (slower, most accurate)
diff := color.DeltaE2000(c1, c2)

// Classic formula
diff := color.DeltaE76(c1, c2)
```

### Wide-Gamut Workflows
```go
// Convert to wide-gamut space
displayP3, _ := color.ConvertToRGBSpace(c, "display-p3")

// Create from wide-gamut values
p3Color, _ := color.ConvertFromRGBSpace(1.0, 0.5, 0.0, 1.0, "display-p3")

// Check if in gamut
if color.InGamut(c) {
    fmt.Println("Fits in sRGB")
}

// Map to gamut with strategy
mapped := color.MapToGamut(c, color.GamutPreserveLightness)
```

### Space System (Advanced)
```go
// Create color in specific space
p3Color := color.NewSpaceColor(
    color.DisplayP3Space,
    []float64{1.0, 0.5, 0.0}, // R, G, B
    1.0, // alpha
)

// Convert between spaces (preserves gamut)
rec2020Color := p3Color.ConvertTo(color.Rec2020Space)

// Get metadata
metadata := color.Metadata(color.DisplayP3Space)
fmt.Printf("Gamut: %.2f× sRGB\n", metadata.GamutVolumeRelativeToSRGB)
```

## Use Cases

### ✅ Perfect For

- **Design Systems** - Generate perceptually uniform color palettes
- **Data Visualization** - Smooth, accurate heatmaps and gradients
- **Photo Editing** - Professional wide-gamut workflows
- **UI Frameworks** - Consistent hover/active/disabled states
- **Color Tools** - Pickers, analyzers, converters
- **Brand Guidelines** - Color consistency checking
- **Accessibility** - Perceptual contrast calculations

### ⚠️ Consider Alternatives For

- **Static hex colors** - Overkill if you just need `#FF0000`
- **Per-pixel GPU operations** - Use shaders instead
- **Simple RGB-only needs** - Standard library may suffice

## Documentation

- **[WHEN_TO_USE.md](WHEN_TO_USE.md)** - Decision trees and use cases
- **[COLOR_PRIMER.md](COLOR_PRIMER.md)** - Comprehensive color theory guide
- **[GRADIENTS.md](GRADIENTS.md)** - Gradient examples and techniques
- **[COLOR_SPACE_ARCHITECTURE.md](COLOR_SPACE_ARCHITECTURE.md)** - Technical design
- **[API Reference](https://pkg.go.dev/github.com/SCKelemen/color)** - Full function docs

## Examples

### Generate UI Theme Palette

```go
brand := color.NewOKLCH(0.55, 0.15, 230, 1.0) // Blue

palette := map[string]string{
    "50":  color.RGBToHex(color.Lighten(brand, 0.45)),
    "100": color.RGBToHex(color.Lighten(brand, 0.35)),
    "200": color.RGBToHex(color.Lighten(brand, 0.25)),
    "300": color.RGBToHex(color.Lighten(brand, 0.15)),
    "400": color.RGBToHex(color.Lighten(brand, 0.05)),
    "500": color.RGBToHex(brand), // Base
    "600": color.RGBToHex(color.Darken(brand, 0.05)),
    "700": color.RGBToHex(color.Darken(brand, 0.15)),
    "800": color.RGBToHex(color.Darken(brand, 0.25)),
    "900": color.RGBToHex(color.Darken(brand, 0.35)),
}
```

### Professional Photo Workflow

```go
// RAW photo in ProPhoto RGB (widest gamut)
raw, _ := color.ConvertFromRGBSpace(0.9, 0.2, 0.1, 1.0, "prophoto-rgb")

// Edit in perceptual space
edited := color.Saturate(raw, 0.15)
edited = color.Lighten(edited, 0.05)

// Export for Display P3 (modern screens)
p3, _ := color.ConvertToRGBSpace(edited, "display-p3")

// Or map to sRGB with quality gamut mapping
srgb := color.MapToGamut(edited, color.GamutProject)
```

### Accessibility Color Checker

```go
func checkContrast(fg, bg Color) {
    diff := color.DeltaE2000(fg, bg)

    if diff < 30 {
        fmt.Println("⚠️ Poor contrast - may fail accessibility")
    } else if diff < 50 {
        fmt.Println("✅ Acceptable contrast")
    } else {
        fmt.Println("✅✅ Excellent contrast")
    }
}
```

## Performance

**Fast enough for:**
- UI operations (< 1000 colors at 60fps) ✅
- Palette generation ✅
- Real-time color pickers ✅
- Batch processing (< 100k colors) ✅

**Optimize when:**
- Processing millions of colors per second
- Per-pixel image operations (batch convert instead)
- GPU shader operations (implement there instead)

## Comparison

| Feature | This Library | `image/color` | CSS |
|---------|--------------|---------------|-----|
| Color spaces | 15+ | 1 (RGB) | ~10 |
| Perceptually uniform | ✅ | ❌ | ⚠️ |
| Wide gamut | ✅ | ❌ | ✅ |
| Gradients | ✅ Smooth | ❌ | Browser only |
| Color difference | ✅ ΔE2000 | ❌ | ❌ |
| Gamut mapping | ✅ 4 strategies | ❌ | ⚠️ Basic |
| Programmatic | ✅ | ✅ | ❌ |

## FAQ

**Q: Do I need to understand color science?**
A: No! Just use `color.Gradient()` and `color.Lighten()` - they do the right thing automatically.

**Q: Is this faster than RGB operations?**
A: Slightly slower (conversions needed), but usually imperceptible. The visual quality improvement is worth it.

**Q: Can I use this in the browser?**
A: Not directly (Go library), but you can compile to WASM or generate colors server-side.

**Q: Does this work with the standard library?**
A: Yes! Converts to/from `image/color.Color` interface.

**Q: Why not just use CSS color functions?**
A: CSS only works in browsers. This works anywhere Go runs - CLIs, backends, image processing, etc.

## Documentation

### 📚 Complete Documentation

**[📖 Documentation Hub](docs/README.md)** - Start here for organized documentation

### By Category

- **[📘 API Reference](docs/reference/api-overview.md)** - Complete API documentation
- **[🎨 Color Space List](docs/reference/color-space-list.md)** - All supported color spaces
- **[🚀 Quickstart Guide](QUICKSTART.md)** - Get started in 5 minutes

### Guides

- **[Gradient Generation](docs/guides/gradients.md)** - Perceptually uniform gradients
- **[LOG Workflows](docs/guides/log-workflows.md)** - Cinema camera color spaces
- **[Lipgloss Integration](docs/guides/lipgloss-integration.md)** - Terminal UI styling

### Learn More

- **[Color Primer](docs/theory/color-primer.md)** - Color science fundamentals
- **[Why Use This Library](docs/theory/why-use-this.md)** - Decision guide
- **[Visual Comparisons](docs/theory/visual-comparisons.md)** - See the difference

## Contributing

Contributions welcome! Areas where help would be appreciated:
- Additional color spaces (e.g., CMYK, HSLuv)
- Performance optimizations
- More examples and documentation
- Visualization tools

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Citation

If you use this library in academic work:

```bibtex
@software{color_go,
  author = {Kelemen, Samuel},
  title = {Color: Professional Color Manipulation for Go},
  year = {2024},
  url = {https://github.com/SCKelemen/color}
}
```

## Acknowledgments

- **Björn Ottosson** - OKLAB and OKLCH color spaces
- **CIE** - LAB, LUV, XYZ color spaces
- **W3C** - CSS Color specifications
- **Go community** - Feedback and contributions

---

**Made with 🎨 by developers who care about color science**

[⭐ Star us on GitHub](https://github.com/SCKelemen/color) | [📖 Read the docs](https://pkg.go.dev/github.com/SCKelemen/color) | [💬 Discuss](https://github.com/SCKelemen/color/discussions)
