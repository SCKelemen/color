# Release v1.0.0 - Initial Release

## 🎉 First Release!

This is the initial release of the comprehensive Go color library, providing full CSS color space support, perceptually uniform operations, and advanced gradient generation.

## ✨ Features

### Color Space Support
- **RGB, HSL, HSV** - Standard color spaces
- **LAB, OKLAB** - CIE color spaces for scientific accuracy
- **LCH, OKLCH** - Perceptually uniform polar color spaces
- **XYZ** - CIE 1931 XYZ color space
- **Wide-gamut RGB** - display-p3, a98-rgb, prophoto-rgb, rec2020, srgb-linear

### CSS Color Parsing
- Parse all major CSS color formats:
  - Hex colors (`#FF0000`, `#F00`, `#FF000080`)
  - RGB/RGBA (`rgb(255, 0, 0)`, `rgb(255 0 0 / 0.5)`)
  - HSL/HSLA (`hsl(0, 100%, 50%)`, `hsl(0 100% 50% / 0.5)`)
  - HWB (`hwb(0 0% 0%)`)
  - LAB/OKLAB/LCH/OKLCH (`lab(50 20 30)`, `oklch(0.7 0.2 120)`)
  - XYZ (`color(xyz 0.5 0.5 0.5)`)
  - Wide-gamut RGB (`color(display-p3 1 0 0)`)
  - Named colors (`red`, `blue`, `transparent`)

### Perceptually Uniform Operations
- `Lighten()` / `Darken()` - Natural-looking color adjustments
- `Saturate()` / `Desaturate()` - Saturation control
- `MixOKLCH()` - Perceptually uniform color mixing
- All operations use OKLCH for best visual results

### Advanced Gradients
- **Basic gradients** - Two-color gradients in any color space
- **Multistop gradients** - Gradients with multiple color stops
- **Non-linear gradients** - 10 easing functions (ease-in, ease-out, ease-in-out, etc.)
- **Color space selection** - Generate gradients in RGB, HSL, LAB, OKLAB, LCH, or OKLCH

### Color Manipulation
- `AdjustHue()` - Shift hue
- `Invert()` - Invert RGB values
- `Grayscale()` - Convert to grayscale
- `Complement()` - Get complementary color
- `Opacity()`, `FadeIn()`, `FadeOut()` - Alpha channel control

### Universal Conversions
- Convert between any supported color formats
- All colors implement the `Color` interface
- Seamless conversion chain: RGB → HSL → OKLCH → LAB → XYZ → display-p3

### Integration
- **lipgloss** - Helper functions for terminal UI integration
- Compatible with standard library `image/color` patterns

## 📊 Statistics

- **Test Coverage**: 82.5%
- **Color Spaces**: 8+ supported
- **CSS Formats**: 15+ parsing formats
- **Easing Functions**: 10
- **Lines of Code**: ~5,000

## 🚀 Getting Started

```go
import "github.com/SCKelemen/color"

// Parse CSS colors
red, _ := color.ParseColor("#FF0000")
blue, _ := color.ParseColor("rgb(0, 0, 255)")

// Create perceptually uniform gradients
gradient := color.Gradient(red, blue, 20)

// Manipulate colors
lightRed := color.Lighten(red, 0.3)
mixed := color.MixOKLCH(red, blue, 0.5)
```

## 📚 Documentation

- Full documentation: https://pkg.go.dev/github.com/SCKelemen/color
- README: See README.md
- Examples: See example_test.go
- lipgloss integration: See LIPGLOSS.md
- Gradients: See GRADIENTS.md

## 🔗 Links

- **GitHub**: https://github.com/SCKelemen/color
- **pkg.go.dev**: https://pkg.go.dev/github.com/SCKelemen/color
- **Issues**: https://github.com/SCKelemen/color/issues

## 📝 Requirements

- Go 1.19 or later

## 🙏 Acknowledgments

This library implements the CSS Color Module Level 4 and Level 5 specifications, providing comprehensive color space support for Go applications.

---

**Full Changelog**: This is the initial release.

