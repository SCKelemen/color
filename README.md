# Color

A comprehensive Go library for color space conversions and color manipulation, supporting all major CSS color spaces including RGB, HSL, LAB, OKLAB, LCH, OKLCH, HSV, and XYZ.

## Features

- **Full Color Space Support**: RGB, HSL, HSV, LAB, OKLAB, LCH, OKLCH, and XYZ
- **Wide-Gamut RGB Support**: display-p3, a98-rgb, prophoto-rgb, rec2020, srgb-linear
- **Alpha Channel Support**: All color spaces support transparency
- **Perceptually Uniform Operations**: Lighten, darken, and gradients use OKLCH for perceptually uniform results
- **Color Mixing**: Mix colors in RGB, HSL, LAB, OKLAB, LCH, or OKLCH space
- **Gradient Generation**: Generate smooth gradients in any color space
- **Universal Conversions**: Convert between any supported color formats
- **Utility Functions**: Lighten, darken, saturate, desaturate, invert, grayscale, complement, and more
- **lipgloss Integration**: Helper functions for seamless integration with [lipgloss](https://github.com/charmbracelet/lipgloss) terminal UI library

## Installation

```bash
go get github.com/SCKelemen/color
```

## Usage

### Basic Color Creation

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/color"
)

func main() {
    // Create RGB color
    red := color.RGB(1.0, 0.0, 0.0)
    fmt.Printf("Red: %+v\n", red)
    
    // Create RGB with alpha
    semiTransparentBlue := color.NewRGBA(0.0, 0.0, 1.0, 0.5)
    
    // Create HSL color
    hsl := color.NewHSL(120, 1.0, 0.5, 1.0) // Green
    
    // Create LAB color
    lab := color.NewLAB(50, 20, 30, 1.0)
    
    // Create OKLAB color
    oklab := color.NewOKLAB(0.6, 0.1, -0.1, 1.0)
    
    // Create LCH color
    lch := color.NewLCH(70, 50, 180, 1.0)
    
    // Create OKLCH color
    oklch := color.NewOKLCH(0.7, 0.2, 120, 1.0)
    
    // Create HSV color
    hsv := color.NewHSV(240, 1.0, 1.0, 1.0) // Blue
}
```

### Parsing Colors

```go
// Parse colors from strings (supports all major CSS color formats)
red, _ := color.ParseColor("#FF0000")
blue, _ := color.ParseColor("rgb(0, 0, 255)")
green, _ := color.ParseColor("hsl(120, 100%, 50%)")
purple, _ := color.ParseColor("oklch(0.6 0.2 300)")
displayP3, _ := color.ParseColor("color(display-p3 1 0 0)")
named, _ := color.ParseColor("red")

// Supported formats:
// - Hex: "#FF0000", "#F00", "#FF000080" (with alpha)
// - RGB: "rgb(255, 0, 0)", "rgb(100%, 0%, 0%)", "rgb(255 0 0)" (modern syntax)
// - RGBA: "rgba(255, 0, 0, 0.5)", "rgb(255 0 0 / 0.5)" (modern syntax)
// - HSL: "hsl(0, 100%, 50%)", "hsl(0 100% 50%)" (modern syntax)
// - HSLA: "hsla(0, 100%, 50%, 0.5)", "hsl(0 100% 50% / 0.5)" (modern syntax)
// - HWB: "hwb(0 0% 0%)", "hwb(0 0% 0% / 0.5)" (Hue, Whiteness, Blackness)
// - HSV: "hsv(0, 100%, 100%)" (not in CSS spec, but commonly used)
// - LAB: "lab(50 20 30)" or "lab(50% 20 30)" (CIE 1976 L*a*b*)
// - OKLAB: "oklab(0.6 0.1 -0.1)" (perceptually uniform)
// - LCH: "lch(70 50 180)" (CIE LCH from LAB)
// - OKLCH: "oklch(0.7 0.2 120)" (perceptually uniform)
// - XYZ: "color(xyz 0.5 0.5 0.5)" (CIE 1931 XYZ via color() function)
// - Wide-gamut RGB: "color(display-p3 1 0 0)", "color(a98-rgb 1 0 0)", etc.
// - Named: "red", "blue", "transparent", etc.
```

### Color Conversions

```go
// Convert between color spaces
rgb := color.RGB(1.0, 0.5, 0.0) // Orange

// Convert to different color spaces
hsl := color.ToHSL(rgb)
lab := color.ToLAB(rgb)
oklab := color.ToOKLAB(rgb)
lch := color.ToLCH(rgb)
oklch := color.ToOKLCH(rgb)
hsv := color.ToHSV(rgb)
xyz := color.ToXYZ(rgb)

// All colors implement the Color interface, so you can convert any color
// to any other color space via RGBA
rgb2 := color.ToHSL(oklch) // Convert OKLCH -> HSL

// Convert to wide-gamut RGB spaces
displayP3, _ := color.ConvertToRGBSpace(rgb, "display-p3")
adobeRGB, _ := color.ConvertToRGBSpace(rgb, "a98-rgb")

// Convert from wide-gamut RGB spaces
fromDisplayP3, _ := color.ConvertFromRGBSpace(1.0, 0.0, 0.0, 1.0, "display-p3")
oklchFromP3 := color.ToOKLCH(fromDisplayP3) // Now manipulate in OKLCH
```

### Gradients (Perceptually Uniform)

```go
red := color.RGB(1, 0, 0)
blue := color.RGB(0, 0, 1)

// Generate gradient in OKLCH space (perceptually uniform - recommended)
gradient := color.Gradient(red, blue, 10) // 10 steps

// Generate gradient in specific color space
gradientHSL := color.GradientInSpace(red, blue, 10, color.GradientHSL)
gradientLAB := color.GradientInSpace(red, blue, 10, color.GradientLAB)
gradientOKLCH := color.GradientInSpace(red, blue, 10, color.GradientOKLCH) // Best quality

// Use gradients
for _, c := range gradient {
    hex := color.RGBToHex(c)
    fmt.Printf("Color: %s\n", hex)
}
```

### Color Manipulation

```go
// Lighten a color (perceptually uniform)
lightBlue := color.Lighten(blue, 0.3) // Lighten by 30%

// Darken a color (perceptually uniform)
darkBlue := color.Darken(blue, 0.3) // Darken by 30%

// Saturate a color
vividRed := color.Saturate(red, 0.5) // Increase saturation by 50%

// Desaturate a color
mutedRed := color.Desaturate(red, 0.5) // Decrease saturation by 50%

// Adjust hue
shifted := color.AdjustHue(red, 60) // Shift hue by 60 degrees

// Mix colors (RGB space)
mixed := color.Mix(red, blue, 0.5) // 50% red, 50% blue

// Mix colors (OKLCH space - perceptually uniform)
uniformMixed := color.MixOKLCH(red, blue, 0.5)

// Mix in any color space
mixedHSL := color.MixInSpace(red, blue, 0.5, color.GradientHSL)
mixedLAB := color.MixInSpace(red, blue, 0.5, color.GradientLAB)

// Invert color
inverted := color.Invert(red)

// Convert to grayscale
gray := color.Grayscale(red)

// Get complementary color
complement := color.Complement(red)
```

### Alpha Channel Operations

```go
// Set opacity
semiTransparent := color.Opacity(red, 0.5)

// Fade out (decrease opacity)
faded := color.FadeOut(red, 0.3) // Reduce opacity by 30%

// Fade in (increase opacity)
moreOpaque := color.FadeIn(red, 0.3) // Increase opacity by 30%

// Get alpha value
alpha := red.Alpha()

// Create new color with different alpha
newAlpha := red.WithAlpha(0.8)
```

### Working with lipgloss

This library provides advanced color manipulation that complements [lipgloss](https://github.com/charmbracelet/lipgloss). While lipgloss supports basic RGB colors, this library adds perceptually uniform color spaces and advanced operations.

```go
package main

import (
    "fmt"
    "github.com/charmbracelet/lipgloss"
    "github.com/SCKelemen/color"
)

func main() {
    // Create a perceptually uniform color
    bgColor := color.NewOKLCH(0.2, 0.1, 240, 1.0) // Dark blue
    
    // Convert to lipgloss-compatible hex string
    hexColor := color.ToLipglossColor(bgColor)
    
    // Use with lipgloss
    style := lipgloss.NewStyle().
        Background(lipgloss.Color(hexColor)).
        Foreground(lipgloss.AdaptiveColor{
            Light: "#000000",
            Dark:  "#ffffff",
        })
    
    fmt.Println(style.Render("Hello, World!"))
}
```

See [LIPGLOSS.md](LIPGLOSS.md) for detailed examples and comparison with lipgloss.

### Advanced: Gradients in Perceptually Uniform Space

```go
// Create a gradient that looks smooth to the human eye
start := color.ParseColor("color(display-p3 1 0 0)") // Display P3 red
end := color.ParseColor("color(display-p3 0 0 1)")   // Display P3 blue

// Generate gradient in OKLCH (perceptually uniform)
gradient := color.Gradient(start, end, 20)

// Convert back to display-p3 for output
for i, c := range gradient {
    displayP3, _ := color.ConvertToRGBSpace(c, "display-p3")
    r, g, b, _ := displayP3.RGBA()
    fmt.Printf("Step %d: display-p3(%.3f %.3f %.3f)\n", i, r, g, b)
}
```

### Universal Conversions

Since all colors implement the `Color` interface, you can convert between any formats:

```go
// Start with any color format
displayP3, _ := color.ParseColor("color(display-p3 1 0 0)")

// Convert to any other format
oklch := color.ToOKLCH(displayP3)        // display-p3 -> OKLCH
lab := color.ToLAB(oklch)                // OKLCH -> LAB
hsl := color.ToHSL(lab)                   // LAB -> HSL
xyz := color.ToXYZ(hsl)                   // HSL -> XYZ
adobeRGB, _ := color.ConvertToRGBSpace(xyz, "a98-rgb") // XYZ -> Adobe RGB

// All conversions work seamlessly!
```

## Color Space Details

- **RGB**: Standard RGB color space (sRGB), values [0, 1]
- **HSL**: Hue, Saturation, Lightness - H: [0, 360), S: [0, 1], L: [0, 1]
- **HSV**: Hue, Saturation, Value - H: [0, 360), S: [0, 1], V: [0, 1]
- **LAB**: CIE LAB color space - L: [0, 100], A and B: unbounded
- **OKLAB**: Perceptually uniform LAB variant - L: [0, 1], A and B: unbounded
- **LCH**: Polar representation of LAB - L: [0, 100], C: [0, ~132], H: [0, 360)
- **OKLCH**: Polar representation of OKLAB - L: [0, 1], C: [0, ~0.4], H: [0, 360)
- **XYZ**: CIE XYZ color space (intermediate for conversions)
- **Wide-gamut RGB**: display-p3, a98-rgb, prophoto-rgb, rec2020, srgb-linear

## API Reference

### Color Interface

All color types implement the `Color` interface:

```go
type Color interface {
    RGBA() (r, g, b, a float64)
    Alpha() float64
    WithAlpha(alpha float64) Color
}
```

### Color Creation Functions

- `RGB(r, g, b float64) *RGBA` - Create RGB color (alpha = 1.0)
- `NewRGBA(r, g, b, a float64) *RGBA` - Create RGB color with alpha
- `NewHSL(h, s, l, a float64) *HSL` - Create HSL color
- `NewHSV(h, s, v, a float64) *HSV` - Create HSV color
- `NewLAB(l, a, b, alpha float64) *LAB` - Create LAB color
- `NewOKLAB(l, a, b, alpha float64) *OKLAB` - Create OKLAB color
- `NewLCH(l, c, h, alpha float64) *LCH` - Create LCH color
- `NewOKLCH(l, c, h, alpha float64) *OKLCH` - Create OKLCH color

### Conversion Functions

- `ToHSL(c Color) *HSL` - Convert to HSL
- `ToHSV(c Color) *HSV` - Convert to HSV
- `ToLAB(c Color) *LAB` - Convert to CIE LAB
- `ToOKLAB(c Color) *OKLAB` - Convert to OKLAB
- `ToLCH(c Color) *LCH` - Convert to CIE LCH
- `ToOKLCH(c Color) *OKLCH` - Convert to OKLCH
- `ToXYZ(c Color) *XYZ` - Convert to CIE XYZ
- `ConvertToRGBSpace(c Color, spaceName string) (*RGBA, error)` - Convert to wide-gamut RGB space (display-p3, a98-rgb, etc.)
- `ConvertFromRGBSpace(r, g, b, a float64, spaceName string) (Color, error)` - Convert from wide-gamut RGB space

**Note**: Since all colors implement the `Color` interface, you can convert between any formats:
```go
// Any color can be converted to any other format via the Color interface
rgb := color.RGB(1, 0, 0)
hsl := color.ToHSL(rgb)           // RGB -> HSL
oklch := color.ToOKLCH(hsl)       // HSL -> OKLCH
lab := color.ToLAB(oklch)         // OKLCH -> LAB
xyz := color.ToXYZ(lab)           // LAB -> XYZ
backToRGB := xyz.RGBA()           // XYZ -> RGB (via interface)
```

### Parsing Functions

- `ParseColor(s string) (Color, error)` - Parse color from string (supports all major CSS color formats)

### Gradient Functions

- `Gradient(start, end Color, steps int) []Color` - Generate perceptually uniform gradient in OKLCH space
- `GradientInSpace(start, end Color, steps int, space GradientSpace) []Color` - Generate gradient in specified color space
- `GradientMultiStop(stops []GradientStop, steps int, space GradientSpace) []Color` - Generate gradient with multiple color stops
- `GradientWithEasing(start, end Color, steps int, space GradientSpace, easing EasingFunction) []Color` - Generate gradient with easing function
- `GradientMultiStopWithEasing(stops []GradientStop, steps int, space GradientSpace, easing EasingFunction) []Color` - Multistop gradient with easing
- `MixInSpace(c1, c2 Color, weight float64, space GradientSpace) Color` - Mix colors in specified color space
- `EaseLinear`, `EaseInQuad`, `EaseOutQuad`, `EaseInOutQuad`, `EaseInCubic`, `EaseOutCubic`, `EaseInOutCubic`, `EaseInSine`, `EaseOutSine`, `EaseInOutSine` - Easing functions for non-linear gradients

### Utility Functions

- `Lighten(c Color, amount float64) Color` - Lighten color (amount: [0, 1])
- `Darken(c Color, amount float64) Color` - Darken color (amount: [0, 1])
- `Saturate(c Color, amount float64) Color` - Increase saturation (amount: [0, 1])
- `Desaturate(c Color, amount float64) Color` - Decrease saturation (amount: [0, 1])
- `Mix(c1, c2 Color, weight float64) Color` - Mix colors in RGB space
- `MixOKLCH(c1, c2 Color, weight float64) Color` - Mix colors in OKLCH space (perceptually uniform)
- `MixInSpace(c1, c2 Color, weight float64, space GradientSpace) Color` - Mix colors in specified color space
- `Gradient(start, end Color, steps int) []Color` - Generate gradient in OKLCH space (perceptually uniform)
- `GradientInSpace(start, end Color, steps int, space GradientSpace) []Color` - Generate gradient in specified color space
- `AdjustHue(c Color, degrees float64) Color` - Shift hue
- `Invert(c Color) Color` - Invert RGB values
- `Grayscale(c Color) Color` - Convert to grayscale
- `Complement(c Color) Color` - Get complementary color
- `Opacity(c Color, opacity float64) Color` - Set opacity
- `FadeOut(c Color, amount float64) Color` - Decrease opacity
- `FadeIn(c Color, amount float64) Color` - Increase opacity
- `ToLipglossColor(c Color) string` - Convert color to lipgloss-compatible hex string
- `ToLipglossColorWithAlpha(c Color) string` - Convert color to lipgloss hex with alpha

## License

MIT
