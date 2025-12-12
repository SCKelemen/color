# Visual Comparisons: Why Color Science Matters

This document shows side-by-side comparisons demonstrating why perceptually uniform color spaces produce better results.

## 1. Lightening Colors: RGB vs OKLCH

### The Problem with RGB

When you "lighten" a color in RGB space by adding the same amount to R, G, and B, the results don't look like what you'd expect:

| Color | RGB Lighten (+0.3) | OKLCH Lighten (20%) | Analysis |
|-------|-------------------|---------------------|----------|
| 🔵 Blue<br/>`RGB(0,0,1)` | 🔷 Cyan-ish<br/>`RGB(0.3,0.3,1.3)` | 🔵 Light Blue<br/>`RGB(0.45,0.45,1)` | RGB adds equal amounts → changes hue!<br/>OKLCH preserves hue → looks lighter |
| 🟢 Green<br/>`RGB(0,1,0)` | 🌿 Lime<br/>`RGB(0.3,1.3,0.3)` | 🟢 Light Green<br/>`RGB(0.3,1,0.3)` | RGB shifts toward yellow<br/>OKLCH stays green |
| 🟡 Yellow<br/>`RGB(1,1,0)` | 🌻 Pale Yellow<br/>`RGB(1.3,1.3,0.3)` | 🟡 Light Yellow<br/>`RGB(1,1,0.4)` | RGB makes it almost white<br/>OKLCH controlled lightening |

### Code Comparison

```go
// ❌ RGB (produces unexpected results)
func lightenRGB(c Color, amount float64) Color {
    r, g, b, a := c.RGBA()
    return NewRGBA(r+amount, g+amount, b+amount, a) // Shifts hue!
}

// ✅ OKLCH (perceptually uniform)
lighter := color.Lighten(c, 0.2) // Actually looks 20% lighter
```

---

## 2. Gradients: The Muddy Middle Problem

### Red to Blue Gradient Comparison

<table>
<tr>
<th>RGB Interpolation</th>
<th>OKLCH Interpolation</th>
</tr>
<tr>
<td>

![RGB Gradient](docs/gradients/gradient_rgb_black.png)

**Problems:**
- Dark/muddy purple in middle
- More purple steps than red/blue
- Uneven perceived brightness

</td>
<td>

![OKLCH Gradient](docs/gradients/gradient_oklch_black.png)

**Benefits:**
- Vibrant purple in middle
- Evenly spaced steps
- Consistent brightness
- Smooth to human eye

</td>
</tr>
</table>

### Why RGB Fails

```
Red → Blue in RGB:
RGB(1, 0, 0) → RGB(0.5, 0, 0.5) → RGB(0, 0, 1)
  Bright    →     Dark!        →    Bright

Lightness graph:
100% ┤╮           ╭
75%  ┤ ╰╮       ╭╯
50%  ┤   ╰────╯      ← Dips in middle!
25%  ┤
     └─────────────
```

### Why OKLCH Works

```
Red → Blue in OKLCH:
L: 0.7 → 0.7 → 0.7     ← Constant lightness!
C: 0.3 → 0.3 → 0.3     ← Constant chroma!
H: 0°  → 180° → 240°   ← Smooth hue transition

Lightness graph:
100% ┤
75%  ┤─────────────   ← Constant!
50%  ┤
25%  ┤
     └─────────────
```

---

## 3. Color Space Gamut Volumes

Visual representation of how much color each space can represent:

```
Relative Gamut Volumes (sRGB = 1.0):

sRGB         ▓▓▓▓▓▓▓▓▓▓ 1.00× (baseline)
Display P3   ▓▓▓▓▓▓▓▓▓▓▓▓▓ 1.26× (+26%)
Adobe RGB    ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 1.44× (+44%)
Rec.2020     ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 1.73× (+73%)
ProPhoto RGB ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 2.89× (+189%)
```

### What This Means

| Space | Additional Colors | Example |
|-------|------------------|---------|
| **sRGB** | Baseline | Standard web red: `rgb(255, 0, 0)` |
| **Display P3** | +26% more colors | iPhone red: 32% more saturated than web red |
| **Rec.2020** | +73% more colors | Future TV red: 56% more saturated than web red |
| **ProPhoto RGB** | +189% more colors | Can represent colors humans can't see! |

### The Problem with Gamut Loss

```go
// Vivid Display P3 red
p3Red := color(display-p3 1 0 0) // Very saturated!

// Force through sRGB (standard libraries)
srgbRed := convertToSRGB(p3Red)  // rgb(255, 0, 0) - lost 26% of vibrancy! 😢

// This library preserves it
p3RedPreserved, _ := color.ConvertToRGBSpace(c, "display-p3") // Still vivid! 🎉
```

---

## 4. Saturation Operations: HSL vs OKLCH

### The HSL Saturation Problem

HSL's "saturation" doesn't actually control perceived colorfulness consistently:

| Base Color | HSL Saturate(+30%) | OKLCH Saturate(+30%) |
|------------|-------------------|---------------------|
| Dark Blue | Becomes brighter AND more colorful | Becomes more colorful only |
| Yellow | Barely changes | Consistently more colorful |
| Cyan | Becomes darker | Pure saturation change |

### Why OKLCH is Better

OKLCH separates:
- **L** (Lightness) - How light/dark
- **C** (Chroma) - How colorful
- **H** (Hue) - Which color

HSL's "S" mixes lightness and chroma together!

```go
// ❌ HSL (saturation affects lightness)
hsl := color.ToHSL(myColor)
hsl.S += 0.3  // Changes perceived brightness too!

// ✅ OKLCH (chroma is independent)
vivid := color.Saturate(myColor, 0.3) // Only affects colorfulness
```

---

## 5. Hue Shifting Example

### Complementary Colors

| Original | RGB "Complement" | OKLCH Complement |
|----------|-----------------|------------------|
| 🔴 Red | 🐟 Cyan (not opposite!) | 🟢 True Green (opposite on color wheel) |
| 🟡 Yellow | 🔵 Blue (close) | 🟣 Purple-Blue (true opposite) |
| 🟠 Orange | ⚪ Pale Blue (washed out) | 💙 Vibrant Blue (same chroma) |

### The Difference

```go
// RGB complement (just inverts RGB values)
rgb := RGB(1, 0.5, 0) // Orange
complement := RGB(0, 0.5, 1) // Pale blue (different lightness/chroma!)

// OKLCH complement (rotates hue 180°, preserves L and C)
orange := color.RGB(1, 0.5, 0)
complement := color.Complement(orange) // Vibrant blue (same lightness/chroma!)
```

---

## 6. Perceptual Uniformity Visualization

### Equal Steps in Different Spaces

```
RGB Space (non-uniform):
Step 1 ██ Looks like 20% change
Step 2 █  Looks like 5% change   ← Uneven!
Step 3 ███ Looks like 30% change
Step 4 █  Looks like 5% change
Step 5 ██ Looks like 20% change

OKLCH Space (uniform):
Step 1 ██ Looks like 20% change
Step 2 ██ Looks like 20% change  ← Even!
Step 3 ██ Looks like 20% change
Step 4 ██ Looks like 20% change
Step 5 ██ Looks like 20% change
```

### Real Example: Blue Palette

```go
blue := color.RGB(0, 0, 1)

// RGB: Add 0.1 to each component
rgb1 := RGB(0.1, 0.1, 1.0)  // Looks VERY different (cyan-ish)
rgb2 := RGB(0.2, 0.2, 1.0)  // Less different
rgb3 := RGB(0.3, 0.3, 1.0)  // Even less different
// Uneven steps!

// OKLCH: Lighten by 0.1
oklch1 := Lighten(blue, 0.1) // Looks 10% lighter
oklch2 := Lighten(blue, 0.2) // Looks 20% lighter
oklch3 := Lighten(blue, 0.3) // Looks 30% lighter
// Even steps!
```

---

## 7. Color Difference: Can Humans Tell Colors Apart?

### DeltaE Visualization

```
ΔE = 0    ┤█ Identical
ΔE < 1    ┤█ Imperceptible (same color to humans)
ΔE 1-2    ┤█░ Barely noticeable (experts only)
ΔE 2-5    ┤█░░ Small difference (most notice)
ΔE 5-10   ┤█░░░ Obvious difference
ΔE > 10   ┤█░░░░ Completely different
```

### Example: Finding Similar Colors

```go
target := color.ParseColor("#FF6B6B")
colors := []Color{
    color.ParseColor("#FF6C6B"), // ΔE = 0.5  ← Almost identical
    color.ParseColor("#FF7676"), // ΔE = 2.1  ← Small difference
    color.ParseColor("#FF0000"), // ΔE = 12.3 ← Very different
}

for _, c := range colors {
    diff := color.DeltaE2000(target, c)
    if diff < 1.0 {
        fmt.Println("Humans can't tell these apart")
    }
}
```

---

## 8. Gamut Mapping Strategies Compared

When converting vivid Display P3 color to sRGB:

| Strategy | Lightness | Chroma | Hue | Use When |
|----------|-----------|--------|-----|----------|
| **Clip** | Changes | Changes | May shift | Speed critical |
| **Preserve Lightness** | ✅ Same | Reduces | ✅ Same | UI backgrounds, text |
| **Preserve Chroma** | Reduces | ✅ Same | ✅ Same | Brand colors, accents |
| **Project** | Slight change | Slight change | ✅ Same | Quality critical |

### Visual Example

Original Display P3 color (out of sRGB gamut):
```
L: 0.7, C: 0.25, H: 150° (vivid teal)
```

Results when mapped to sRGB:

```
Clip:                L: 0.68 ✗  C: 0.22 ✗  H: 145° ✗  (hue shifted!)
Preserve Lightness:  L: 0.70 ✅  C: 0.18 ✗  H: 150° ✅  (less vivid, same brightness)
Preserve Chroma:     L: 0.65 ✗  C: 0.25 ✅  H: 150° ✅  (darker, same saturation)
Project:             L: 0.69 ~  C: 0.20 ~  H: 149° ~  (best overall compromise)
```

---

## 9. Multi-Stop Gradient Comparison

### Red → Yellow → Blue

<table>
<tr>
<th>RGB</th>
<th>HSL</th>
<th>OKLCH ⭐</th>
</tr>
<tr>
<td>

- Muddy brown in red→yellow
- Dark purple in yellow→blue
- Uneven steps

</td>
<td>

- Better hue transition
- Still brightness inconsistency
- Yellow section looks "washed out"

</td>
<td>

- Clean vibrant transitions
- Consistent brightness throughout
- Evenly spaced to human eye

</td>
</tr>
</table>

```go
stops := []color.GradientStop{
    {Color: red,    Position: 0.0},
    {Color: yellow, Position: 0.5},
    {Color: blue,   Position: 1.0},
}

// RGB: muddy transitions
rgbGrad := color.GradientMultiStop(stops, 30, color.GradientRGB)

// OKLCH: clean, vibrant transitions
oklchGrad := color.GradientMultiStop(stops, 30, color.GradientOKLCH)
```

---

## 10. Before & After: Real-World Examples

### Design System Palette Generation

**Before (Manual RGB adjustments):**
```
Base:    #3B82F6
Light 1: #6BA3FF  ← Not uniform
Light 2: #9BC4FF  ← Steps feel uneven
Light 3: #CBE5FF  ← Too light!
```

**After (OKLCH-based generation):**
```go
base := color.ParseColor("#3B82F6")
palette := color.Gradient(
    color.Lighten(base, 0.3),
    color.Darken(base, 0.3),
    7,
)
// Each step looks evenly spaced!
```

### Heatmap Colors

**Before (RGB interpolation):**
- Dark muddy section in middle
- Uneven temperature perception
- Hard to read values

**After (OKLCH interpolation):**
- Smooth, even progression
- Intuitive hot-to-cold perception
- Easy to read precise values

### Photo Editing Workflow

**Before (sRGB pipeline):**
```
RAW → sRGB → Edit → sRGB output
       ↓
    Loses 73% of ProPhoto RGB gamut!
```

**After (Wide-gamut pipeline):**
```
RAW → ProPhoto RGB → Edit in OKLCH → Display P3 output
                                  ↓
                    Preserves vibrant colors!
```

---

## Summary: Why Color Science Matters

| Operation | Standard (RGB) | This Library (OKLCH) |
|-----------|----------------|---------------------|
| Lighten | Changes hue ❌ | Preserves hue ✅ |
| Gradients | Muddy middle ❌ | Smooth, vibrant ✅ |
| Saturate | Affects brightness ❌ | Only affects color ✅ |
| Wide-gamut | Loses vibrancy ❌ | Preserves it ✅ |
| Steps | Uneven perception ❌ | Perceptually uniform ✅ |
| Color matching | Guesswork ❌ | Scientific metrics ✅ |

**The bottom line:** If humans will see your colors, use perceptually uniform color spaces. Your users will notice the difference, even if they can't explain why it looks better.

---

## Try It Yourself

```bash
go run examples/comparison.go
```

Open the generated images side-by-side to see the difference!
