# When to Use This Color Library

## Quick Decision Tree

```
Are you working with colors?
│
├─ YES → Are you doing more than displaying static hex colors?
│         │
│         ├─ NO → Basic library is probably enough
│         │
│         └─ YES → Choose your path below ↓
│
└─ NO → You probably don't need this library 😉
```

## Choose Your Path

### Path 1: "My gradients look terrible"

**Symptoms:**
- RGB gradients have muddy midpoints
- Steps don't look evenly spaced
- Users complain colors "pop" in some sections and fade in others

**Solution:**
```go
// Instead of this (muddy):
gradient := interpolateRGB(red, blue, 10)

// Use this (smooth):
gradient := color.Gradient(red, blue, 10) // Uses OKLCH
```

**Why it works:** OKLCH is perceptually uniform - equal mathematical steps look equal to human eyes.

**Real-world use cases:**
- Data visualization gradients
- UI theme colors
- Heatmaps
- Progress indicators
- Animated color transitions

---

### Path 2: "I need to lighten/darken colors properly"

**Symptoms:**
- Lightening blue turns it cyan
- Darkening yellow turns it brown/green
- Saturation changes feel inconsistent

**Problem:**
```go
// This looks wrong:
blue := RGB(0, 0, 1)
lighter := RGB(0.2, 0.2, 1.2) // Oops, cyan-ish now
```

**Solution:**
```go
blue := color.RGB(0, 0, 1)
lighter := color.Lighten(blue, 0.2) // Actually looks like lighter blue!
```

**Why it works:** Operations happen in OKLCH where "lighter" means lighter, not "add RGB values."

**Real-world use cases:**
- Hover states (lighten 10%)
- Disabled states (darken 30%)
- Active states (saturate 20%)
- Theme generation (create palette from single color)

---

### Path 3: "I'm working with Display P3 / wide-gamut colors"

**Symptoms:**
- Colors lose vibrancy when converted
- Can't take advantage of modern displays
- Photos look dull compared to Apple Photos

**Problem:**
Most libraries force everything through sRGB, losing 26% of Display P3's gamut:
```go
// Display P3 red (very vivid)
p3Red := displayP3.Color(1, 0, 0)

// Forced through sRGB
standardLib.Convert(p3Red) // Loses vibrancy! 😢
```

**Solution:**
```go
// Create in Display P3
p3Red, _ := color.ConvertFromRGBSpace(1, 0, 0, 1, "display-p3")

// Operate in perceptual space
lighter := color.Lighten(p3Red, 0.2)

// Convert back to Display P3
result, _ := color.ConvertToRGBSpace(lighter, "display-p3")
// Still Display P3! Still vivid! 🎉
```

**Why it works:** We preserve the color space throughout, only converting at the last moment.

**Real-world use cases:**
- iOS/macOS apps targeting modern devices
- Professional photo editing
- Digital cinema workflows
- HDR content creation
- Future-proof color pipelines

---

### Path 4: "I need to match colors accurately"

**Symptoms:**
- Need to know if two colors look the same
- Building color-matching tools
- Quality control / color consistency checking

**Solution:**
```go
color1 := color.ParseColor("#FF6B6B")
color2 := color.ParseColor("#FF6D6C")

// Industry-standard perceptual difference
diff := color.DeltaE2000(color1, color2)

if diff < 1.0 {
    fmt.Println("Humans can't tell these apart")
}
```

**Real-world use cases:**
- Palette generators (ensure colors are distinguishable)
- Print color matching
- Brand color consistency
- Accessibility (ensuring sufficient contrast)
- Color naming ("is this red or orange?")

---

### Path 5: "I'm building professional color tools"

**You need this library if you're building:**

- **Photo Editing Software**
  ```go
  // RAW in ProPhoto RGB → Edit in OKLCH → Export to Display P3
  raw, _ := color.ConvertFromRGBSpace(r, g, b, 1, "prophoto-rgb")
  edited := color.Saturate(color.Lighten(raw, 0.1), 0.15)
  p3, _ := color.ConvertToRGBSpace(edited, "display-p3")
  ```

- **Design Tools**
  ```go
  // Generate harmonious palettes
  base := color.NewOKLCH(0.6, 0.2, 200, 1.0)
  analogous := []Color{
      color.AdjustHue(base, -30),
      base,
      color.AdjustHue(base, 30),
  }
  ```

- **Data Visualization**
  ```go
  // Perceptually uniform gradients for heatmaps
  coldToHot := color.Gradient(
      color.RGB(0, 0, 1),   // Blue
      color.RGB(1, 0, 0),   // Red
      100,
  )
  ```

- **Color Analysis Tools**
  ```go
  // Find most different color for maximum contrast
  bestContrast := findMaxDeltaE(color1, candidateColors)
  ```

---

## Comparison with Other Libraries

### vs. Standard `image/color`

| Feature | `image/color` | This Library |
|---------|---------------|--------------|
| Color spaces | RGB only | RGB, HSL, HSV, LAB, OKLAB, LCH, OKLCH, HWB, LUV, XYZ |
| Gradients | ❌ | ✅ Perceptually uniform |
| Lighten/darken | ❌ | ✅ Perceptually correct |
| Wide gamut | ❌ | ✅ Display P3, Rec.2020, ProPhoto RGB |
| Color difference | ❌ | ✅ DeltaE76, DeltaE2000, DeltaEOK |
| Gamut mapping | ❌ | ✅ 4 strategies |

**Use `image/color` when:**
- Basic image manipulation
- Only need RGB
- Performance is critical
- Color science doesn't matter

**Use this library when:**
- Need perceptually uniform operations
- Working with multiple color spaces
- Quality matters

### vs. CSS Color Functions

| Feature | CSS | This Library |
|---------|-----|--------------|
| Parse CSS colors | ✅ | ✅ |
| Generate gradients | ✅ (browser only) | ✅ (in code) |
| Color manipulation | ⚠️ (basic) | ✅✅ (advanced) |
| Programmatic access | ❌ | ✅ |
| Custom algorithms | ❌ | ✅ |

**Use CSS when:**
- Working in the browser
- Static colors
- No computation needed

**Use this library when:**
- Server-side rendering
- Dynamic palette generation
- Programmatic color manipulation
- Need full control

---

## Real-World Use Cases

### 1. Theme Generator for Design Systems

```go
// Input: Brand color
brand := color.ParseColor("#3B82F6") // Blue

// Generate full palette
palette := struct {
    Lighter []string
    Base    string
    Darker  []string
}{
    Lighter: hexGradient(color.Lighten(brand, 0.3), brand, 4),
    Base:    color.RGBToHex(brand),
    Darker:  hexGradient(brand, color.Darken(brand, 0.3), 4),
}

// All steps are perceptually uniform!
```

**Why this library:**
- Perceptually uniform steps
- Professional-quality palettes
- Consistent lightness/darkness across hues

### 2. Heatmap with Proper Color Progression

```go
// Bad: RGB gradient (muddy middle, uneven steps)
heatmap := linearInterpolate(blue, red, 100)

// Good: OKLCH gradient (smooth, even steps)
heatmap := color.Gradient(
    color.RGB(0, 0, 1),   // Blue (cold)
    color.RGB(1, 0, 0),   // Red (hot)
    100,
)

// Map data to colors
for _, dataPoint := range data {
    normalized := (dataPoint - min) / (max - min)
    idx := int(normalized * 99)
    pixelColor := heatmap[idx]
}
```

**Why this library:**
- Smooth transitions
- No muddy midpoints
- Accurate data → color mapping

### 3. Accessible Color Contrast Tool

```go
func findAccessiblePair(bg Color, candidates []Color) Color {
    for _, fg := range candidates {
        // Ensure sufficient perceptual difference
        if color.DeltaE2000(bg, fg) > 50 { // Simplified
            // Also check WCAG contrast
            if meetsWCAG(bg, fg) {
                return fg
            }
        }
    }
    return nil
}
```

**Why this library:**
- Perceptual color difference metrics
- Accurate contrast calculations
- Color space flexibility

### 4. Photo Filter Pipeline

```go
// Start with wide gamut
raw, _ := color.ConvertFromRGBSpace(r, g, b, 1, "prophoto-rgb")

// Apply filters in perceptual space
filtered := raw
filtered = color.Saturate(filtered, 0.2)    // Vibrance
filtered = color.Lighten(filtered, 0.05)    // Exposure
filtered = color.AdjustHue(filtered, 10)    // Warmth

// Export for target device
p3, _ := color.ConvertToRGBSpace(filtered, "display-p3")
// or
srgb := color.MapToGamut(filtered, color.GamutProject)
```

**Why this library:**
- Wide gamut preservation
- Perceptually uniform edits
- Professional gamut mapping

### 5. Brand Color Consistency Checker

```go
func checkBrandColors(designed, actual Color) {
    diff := color.DeltaE2000(designed, actual)

    switch {
    case diff < 1:
        fmt.Println("✅ Perfect match")
    case diff < 2:
        fmt.Println("✅ Acceptable (barely noticeable)")
    case diff < 5:
        fmt.Printf("⚠️  Noticeable difference (ΔE=%.2f)\n", diff)
    default:
        fmt.Printf("❌ Significant difference (ΔE=%.2f)\n", diff)
    }
}
```

**Why this library:**
- Industry-standard color difference
- Quantifiable matching
- Professional quality control

---

## Performance Considerations

### When This Library Is Fast Enough

- UI operations (< 1000 colors at 60fps)
- Batch processing (< 100k colors)
- Single-threaded palette generation
- Most real-time applications

### When to Optimize

If you're processing millions of colors per second:

```go
// Slower (converts for each operation)
for _, c := range millionsOfColors {
    result := color.Lighten(color.Saturate(c, 0.5), 0.3)
}

// Faster (batch convert)
oklchs := convertAllToOKLCH(millionsOfColors)
for _, oklch := range oklchs {
    oklch.L += 0.3
    oklch.C *= 1.5
}
results := convertAllToRGB(oklchs)
```

---

## FAQ

### "I only need RGB hex colors. Do I need this?"

**No.** If you're just passing hex strings around, this library is overkill. Use a simple hex parser.

### "Can I use this for games/real-time graphics?"

**Yes, usually.** For UI elements, absolutely. For shaders/per-pixel operations, consider caching converted colors or doing conversions on GPU.

### "Will this work on the web?"

**Not directly** - this is a Go library. But you can:
- Generate palettes server-side, export to CSS
- Use it in Go WASM modules
- Use it in backend APIs that serve colors to frontend

### "Do I need to understand color science?"

**No.** The quick decision tree and examples should be enough. But understanding helps you make better decisions!

### "What's the learning curve?"

**5 minutes to productive:**
```go
// This is all you need to know:
gradient := color.Gradient(startColor, endColor, steps)
lighter := color.Lighten(myColor, 0.2)
hex := color.RGBToHex(result)
```

**1 hour to proficient:**
- Understand OKLCH basics
- Know when to use each color space
- Use gamut mapping correctly

**Ongoing:** Color science is deep! But you don't need to know it all.

---

## Decision Matrix

| Your Situation | Use This Library? | Alternative |
|----------------|-------------------|-------------|
| Building color picker UI | ⚠️ Maybe | Consider CSS or basic RGB |
| Generating theme palettes | ✅ Yes | - |
| Lightening/darkening colors | ✅ Yes | - |
| Working with images | ⚠️ Depends | Use with `image` package |
| Wide-gamut displays | ✅ Absolutely | No good alternative |
| Professional color tools | ✅ Absolutely | - |
| Data visualization | ✅ Yes | - |
| Basic hex color storage | ❌ No | Simple hex parser |
| Per-pixel GPU operations | ❌ No | Implement in shader |
| Color science research | ✅ Yes | But may need more |

---

## Getting Started

If you've decided this library is right for you:

1. **Read the [Quick Start](#)** - 5 minute tutorial
2. **Try the examples** - Copy-paste and modify
3. **Consult [COLOR_PRIMER.md](COLOR_PRIMER.md)** - When you need deeper understanding
4. **Check [API docs](README.md)** - Full function reference

Welcome to professional color handling! 🎨
