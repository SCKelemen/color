# Using with lipgloss

[lipgloss](https://github.com/charmbracelet/lipgloss) is a popular Go library for styling terminal UIs. While lipgloss provides basic RGB color support, this library adds powerful color space conversions and perceptually uniform color manipulation that lipgloss doesn't support.

## What lipgloss supports

lipgloss supports:
- Basic RGB colors via hex strings (e.g., `#FF0000`)
- ANSI color names
- Terminal color support

## What this library adds

This library extends lipgloss with:

1. **Perceptually uniform color spaces** (OKLAB, OKLCH) for better color manipulation
2. **Advanced color space conversions** (LAB, LCH, HSL, HSV, XYZ)
3. **Perceptually uniform operations** - lighten/darken that look natural
4. **Color mixing in perceptually uniform space** - better gradients
5. **Alpha channel support** across all color spaces
6. **Advanced color manipulation** - saturate, desaturate, adjust hue, complement, etc.

## Examples

### Basic Integration

```go
package main

import (
    "fmt"
    "github.com/charmbracelet/lipgloss"
    "github.com/SCKelemen/color"
)

func main() {
    // Create a color using this library
    bgColor := color.NewOKLCH(0.2, 0.1, 240, 1.0) // Dark blue
    
    // Convert to RGB for lipgloss
    r, g, b, _ := bgColor.RGBA()
    
    // Create lipgloss style
    style := lipgloss.NewStyle().
        Background(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
            uint8(r*255), uint8(g*255), uint8(b*255))))
    
    fmt.Println(style.Render("Hello, World!"))
}
```

### Perceptually Uniform Color Schemes

```go
// Create a base color
base := color.NewOKLCH(0.6, 0.2, 180, 1.0) // Teal

// Generate a palette with perceptually uniform lightness
light := color.Lighten(base, 0.3)
dark := color.Darken(base, 0.3)
complement := color.Complement(base)

// Convert to lipgloss colors
colors := map[string]string{
    "base":      toLipglossColor(base),
    "light":     toLipglossColor(light),
    "dark":      toLipglossColor(dark),
    "complement": toLipglossColor(complement),
}

func toLipglossColor(c color.Color) string {
    r, g, b, _ := c.RGBA()
    return fmt.Sprintf("#%02x%02x%02x",
        uint8(r*255), uint8(g*255), uint8(b*255))
}
```

### Creating Gradients

```go
// Create a gradient between two colors in perceptually uniform space
start := color.NewOKLCH(0.3, 0.2, 240, 1.0) // Dark blue
end := color.NewOKLCH(0.8, 0.2, 60, 1.0)    // Light yellow

// Mix in OKLCH space for perceptually uniform gradient
steps := 10
for i := 0; i <= steps; i++ {
    weight := float64(i) / float64(steps)
    mixed := color.MixOKLCH(start, end, weight)
    
    // Use with lipgloss
    style := lipgloss.NewStyle().
        Background(lipgloss.Color(toLipglossColor(mixed)))
    
    fmt.Print(style.Render(" "))
}
```

### Dynamic Theme Generation

```go
// Generate a theme from a single base color
func generateTheme(baseHue float64) map[string]string {
    // Create colors with different lightness but same hue
    bg := color.NewOKLCH(0.15, 0.05, baseHue, 1.0)      // Dark background
    fg := color.NewOKLCH(0.95, 0.05, baseHue, 1.0)       // Light foreground
    accent := color.NewOKLCH(0.6, 0.3, baseHue, 1.0)     // Accent color
    muted := color.NewOKLCH(0.5, 0.1, baseHue, 1.0)      // Muted color
    
    return map[string]string{
        "bg":     toLipglossColor(bg),
        "fg":     toLipglossColor(fg),
        "accent": toLipglossColor(accent),
        "muted":  toLipglossColor(muted),
    }
}

// Use with different themes
blueTheme := generateTheme(240)   // Blue theme
greenTheme := generateTheme(120)  // Green theme
redTheme := generateTheme(0)      // Red theme
```

### Color Manipulation

```go
// Start with a color from hex (lipgloss format)
hexColor, _ := color.HexToRGB("#FF5733") // Orange

// Manipulate it
lighter := color.Lighten(hexColor, 0.2)
darker := color.Darken(hexColor, 0.2)
saturated := color.Saturate(hexColor, 0.3)
desaturated := color.Desaturate(hexColor, 0.3)
shifted := color.AdjustHue(hexColor, 60) // Shift hue by 60 degrees

// Convert back to lipgloss
styles := map[string]lipgloss.Style{
    "original":   lipgloss.NewStyle().Foreground(lipgloss.Color(color.RGBToHex(hexColor))),
    "lighter":    lipgloss.NewStyle().Foreground(lipgloss.Color(color.RGBToHex(lighter))),
    "darker":     lipgloss.NewStyle().Foreground(lipgloss.Color(color.RGBToHex(darker))),
    "saturated":  lipgloss.NewStyle().Foreground(lipgloss.Color(color.RGBToHex(saturated))),
    "desaturated": lipgloss.NewStyle().Foreground(lipgloss.Color(color.RGBToHex(desaturated))),
    "shifted":    lipgloss.NewStyle().Foreground(lipgloss.Color(color.RGBToHex(shifted))),
}
```

### Helper Function

Here's a convenient helper function to convert any color to a lipgloss-compatible hex string:

```go
func ToLipglossColor(c color.Color) string {
    r, g, b, _ := c.RGBA()
    return fmt.Sprintf("#%02x%02x%02x",
        uint8(r*255), uint8(g*255), uint8(b*255))
}

// Usage
bgColor := color.NewOKLCH(0.2, 0.1, 240, 1.0)
style := lipgloss.NewStyle().
    Background(lipgloss.Color(ToLipglossColor(bgColor)))
```

## Why Perceptually Uniform Matters

When you lighten or darken colors in RGB space, the results can look unnatural. For example, lightening a blue might make it look washed out, while lightening a yellow might not change much.

This library uses OKLCH (a perceptually uniform color space) for operations like `Lighten()` and `Darken()`, which means:
- Colors maintain their perceived intensity
- Gradients look smooth and natural
- Color relationships are preserved

## Comparison

| Feature | lipgloss | This Library |
|---------|----------|---------------|
| RGB colors | ✅ | ✅ |
| Hex strings | ✅ | ✅ |
| HSL/HSV | ❌ | ✅ |
| LAB/OKLAB | ❌ | ✅ |
| LCH/OKLCH | ❌ | ✅ |
| Perceptually uniform operations | ❌ | ✅ |
| Color mixing | ❌ | ✅ |
| Advanced manipulation | ❌ | ✅ |
| Alpha channel | ❌ | ✅ |

## Summary

This library complements lipgloss by providing advanced color manipulation capabilities. Use this library to:
- Create perceptually uniform color schemes
- Generate smooth gradients
- Manipulate colors in ways that look natural
- Convert between color spaces
- Work with alpha channels

Then convert the results to lipgloss-compatible hex strings for use in your terminal UI.

