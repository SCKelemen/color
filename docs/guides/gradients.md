# Gradients and Color Space Conversions

This library supports generating gradients in perceptually uniform color spaces and converting between all supported formats.

## Why Perceptually Uniform Gradients?

When you create a gradient in RGB space, the steps may not appear evenly spaced to the human eye. For example, a gradient from blue to yellow in RGB might look like it has more blue steps and fewer yellow steps, even though they're mathematically equal.

**Perceptually uniform color spaces** (like OKLCH) ensure that equal steps in the color space appear as equal steps to the human eye.

## Gradient Functions

### Basic Gradient (OKLCH - Recommended)

```go
red := color.RGB(1, 0, 0)
blue := color.RGB(0, 0, 1)

// Generate perceptually uniform gradient
gradient := color.Gradient(red, blue, 10) // 10 steps

for i, c := range gradient {
    hex := color.RGBToHex(c)
    fmt.Printf("Step %d: %s\n", i, hex)
}
```

### Gradient in Specific Color Space

```go
// Generate gradient in different color spaces
gradientRGB := color.GradientInSpace(red, blue, 10, color.GradientRGB)
gradientHSL := color.GradientInSpace(red, blue, 10, color.GradientHSL)
gradientLAB := color.GradientInSpace(red, blue, 10, color.GradientLAB)
gradientOKLAB := color.GradientInSpace(red, blue, 10, color.GradientOKLAB)
gradientLCH := color.GradientInSpace(red, blue, 10, color.GradientLCH)
gradientOKLCH := color.GradientInSpace(red, blue, 10, color.GradientOKLCH) // Best quality
```

## Converting Between All Formats

### Universal Conversion

Since all colors implement the `Color` interface, you can convert between any formats:

```go
// Start with any format
displayP3, _ := color.ParseColor("color(display-p3 1 0 0)")

// Convert through any chain of color spaces
oklch := color.ToOKLCH(displayP3)        // display-p3 -> OKLCH
lab := color.ToLAB(oklch)                // OKLCH -> LAB
hsl := color.ToHSL(lab)                  // LAB -> HSL
xyz := color.ToXYZ(hsl)                  // HSL -> XYZ
hsv := color.ToHSV(xyz)                  // XYZ -> HSV (via RGB)
adobeRGB, _ := color.ConvertToRGBSpace(hsv, "a98-rgb") // HSV -> Adobe RGB

// All conversions work seamlessly!
```

### Wide-Gamut RGB Conversions

```go
// Convert to wide-gamut RGB spaces
rgb := color.RGB(1, 0, 0)
displayP3, _ := color.ConvertToRGBSpace(rgb, "display-p3")
adobeRGB, _ := color.ConvertToRGBSpace(rgb, "a98-rgb")
proPhoto, _ := color.ConvertToRGBSpace(rgb, "prophoto-rgb")
rec2020, _ := color.ConvertToRGBSpace(rgb, "rec2020")

// Convert from wide-gamut RGB spaces
fromDisplayP3, _ := color.ConvertFromRGBSpace(1.0, 0.0, 0.0, 1.0, "display-p3")
// Now you can manipulate it in any color space
oklch := color.ToOKLCH(fromDisplayP3)
lightened := color.Lighten(oklch, 0.2)
```

## Complete Workflow Example

```go
// 1. Parse colors from different formats
color1, _ := color.ParseColor("color(display-p3 1 0 0)")
color2, _ := color.ParseColor("color(a98-rgb 0 0 1)")

// 2. Generate gradient in perceptually uniform space
gradient := color.Gradient(color1, color2, 20)

// 3. Convert each step back to display-p3 for output
for i, c := range gradient {
    // Convert to display-p3
    displayP3, _ := color.ConvertToRGBSpace(c, "display-p3")
    r, g, b, a := displayP3.RGBA()
    
    fmt.Printf("Step %d: color(display-p3 %.3f %.3f %.3f / %.2f)\n", 
        i, r, g, b, a)
}
```

## Color Space Comparison for Gradients

| Color Space | Perceptually Uniform | Best For |
|-------------|---------------------|----------|
| RGB | ❌ | Fast, simple gradients |
| HSL | ⚠️ | Hue-based gradients |
| LAB | ✅ | Scientific accuracy |
| OKLAB | ✅✅ | Modern perceptually uniform |
| LCH | ✅ | Hue-based, perceptually uniform |
| **OKLCH** | ✅✅✅ | **Best overall (recommended)** |

## Example: Creating a Theme Palette

```go
// Create a base color
base := color.NewOKLCH(0.6, 0.2, 180, 1.0) // Teal

// Generate variations using gradients
light := color.Lighten(base, 0.3)
dark := color.Darken(base, 0.3)

// Create a gradient for intermediate shades
shades := color.Gradient(dark, light, 5)

// Convert all to hex for use
palette := make([]string, len(shades))
for i, c := range shades {
    palette[i] = color.RGBToHex(c)
}
```

## Example: Converting Wide-Gamut to Standard

```go
// You have a color in display-p3
displayP3Color, _ := color.ParseColor("color(display-p3 1 0.5 0)")

// Convert to OKLCH for manipulation
oklch := color.ToOKLCH(displayP3Color)

// Manipulate (perceptually uniform)
lightened := color.Lighten(oklch, 0.2)
saturated := color.Saturate(lightened, 0.3)

// Convert back to display-p3
backToP3, _ := color.ConvertToRGBSpace(saturated, "display-p3")
r, g, b, _ := backToP3.RGBA()

// Or convert to sRGB for standard display
sRGB := saturated.RGBA() // Already in sRGB via Color interface
```

## Custom Easing Functions

You can easily create your own easing functions! The `EasingFunction` type is simply a function that maps `[0, 1]` to `[0, 1]`.

### Creating Custom Easing Functions

```go
import (
    "math"
    "github.com/SCKelemen/color"
)

// Example 1: Simple exponential easing
myEasing := color.EasingFunction(func(t float64) float64 {
    return 1 - math.Pow(1-t, 3) // Ease-out cubic
})

// Example 2: Elastic easing (bouncy effect)
elasticEasing := color.EasingFunction(func(t float64) float64 {
    if t == 0 || t == 1 {
        return t
    }
    c4 := (2 * math.Pi) / 3
    return math.Pow(2, -10*t) * math.Sin((t*10-0.75)*c4) + 1
})

// Example 3: Bounce easing
bounceEasing := color.EasingFunction(func(t float64) float64 {
    if t < 1/2.75 {
        return 7.5625 * t * t
    } else if t < 2/2.75 {
        t -= 1.5 / 2.75
        return 7.5625*t*t + 0.75
    } else if t < 2.5/2.75 {
        t -= 2.25 / 2.75
        return 7.5625*t*t + 0.9375
    } else {
        t -= 2.625 / 2.75
        return 7.5625*t*t + 0.984375
    }
})

// Use your custom easing function
red := color.RGB(1, 0, 0)
blue := color.RGB(0, 0, 1)
gradient := color.GradientWithEasing(red, blue, 20, color.GradientOKLCH, myEasing)
```

### Easing Function Requirements

Your custom easing function must:
- Accept a `float64` parameter `t` in range `[0, 1]`
- Return a `float64` in range `[0, 1]`
- Map `0` to `0` and `1` to `1`

### Built-in Easing Functions

The library provides 10 built-in easing functions:
- `EaseLinear` - No easing (linear)
- `EaseInQuad`, `EaseOutQuad`, `EaseInOutQuad` - Quadratic curves
- `EaseInCubic`, `EaseOutCubic`, `EaseInOutCubic` - Cubic curves
- `EaseInSine`, `EaseOutSine`, `EaseInOutSine` - Sinusoidal curves

### Example: Custom Easing with Multistop Gradients

```go
// Create a custom easing that emphasizes the middle
middleEmphasis := color.EasingFunction(func(t float64) float64 {
    // Slow at start, fast in middle, slow at end
    if t < 0.5 {
        return 2 * t * t // Ease-in quadratic
    }
    return 1 - 2*(1-t)*(1-t) // Ease-out quadratic
})

stops := []color.GradientStop{
    {Color: color.RGB(1, 0, 0), Position: 0.0},   // Red
    {Color: color.RGB(1, 1, 0), Position: 0.5}, // Yellow
    {Color: color.RGB(0, 0, 1), Position: 1.0}, // Blue
}

gradient := color.GradientMultiStopWithEasing(stops, 20, color.GradientOKLCH, middleEmphasis)
```

