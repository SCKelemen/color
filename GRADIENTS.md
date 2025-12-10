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

