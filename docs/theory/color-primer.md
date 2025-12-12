# Understanding Color: A Comprehensive Guide

## Why This Library Exists

Most color libraries treat RGB as the only color space and force every operation through it. This causes three major problems:

1. **Perceptual incorrectness**: Lightening blue in RGB space makes it cyan, not lighter blue
2. **Wide-gamut loss**: Display P3 reds lose their vibrancy when forced through sRGB
3. **Poor gradients**: RGB gradients have muddy midpoints and uneven steps

This library solves these problems by:
- Using **XYZ as the universal hub** for lossless color space conversions
- Performing operations in **perceptually uniform spaces** (OKLCH) so "lighten 20%" looks like 20% to human eyes
- Preserving **wide-gamut colors** throughout the pipeline until you explicitly convert to sRGB

## Quick Decision Guide

**When should I use this library?**

| You Need... | Use This Library If... |
|------------|------------------------|
| **Perceptually uniform operations** | Lighten/darken should look even to human eyes |
| **Smooth gradients** | RGB gradients look muddy or uneven |
| **Wide-gamut support** | Working with Display P3, DCI-P3, Rec.2020, ProPhoto RGB |
| **Color science accuracy** | Need proper chromatic adaptation, gamut mapping, color difference metrics |
| **Multiple color spaces** | Converting between HSL, LAB, OKLCH, HWB, LUV, etc. |

---

## Color Fundamentals

### What is Color?

- **Physics**: Light is electromagnetic waves with wavelengths between ~380-700nm
- **Biology**: Human eyes have 3 cone types (S, M, L) that respond to different wavelength ranges
- **Math**: Since we only have 3 receptor types, we can represent any color as a 3-component vector

### What is a Color Space?

A **coordinate system** for representing colors. Like GPS coordinates for locations, color spaces define:

- **Axes/primaries**: What the dimensions represent (e.g., Red, Green, Blue)
- **White point**: What counts as "white" (D65 = 6500K daylight, D50 = 5000K horizon)
- **Transfer function**: How to encode values (gamma correction, linear, etc.)
- **Conversion math**: How to translate to other spaces

### What is a Color Gamut?

The **range of colors** a space or device can represent. Think of it as:
- **Volume**: 3D region in a reference space (usually XYZ or OKLCH)
- **Boundary**: Surface defining what's representable
- **Comparison**: Display P3 is 26% larger than sRGB; Rec.2020 is 73% larger

### Visualizing Gamuts

<picture>
  <source media="(prefers-color-scheme: dark)"  srcset="docs/gamuts/gamut_comparison_white.png">
  <source media="(prefers-color-scheme: light)" srcset="docs/gamuts/gamut_comparison_black.png">
  <img alt="Gamut Volume Comparison" src="docs/gamuts/gamut_comparison_white.png">
</picture>

*Relative gamut sizes: sRGB (1.0×), Display P3 (1.26×), Adobe RGB (1.44×), Rec.2020 (1.73×), ProPhoto RGB (2.89×)*

---

## Color Spaces We Support

### RGB Color Spaces

RGB spaces define colors by mixing red, green, and blue light. Each has different **primaries** (which specific reds/greens/blues) and **transfer functions** (how values are encoded).

#### sRGB - Standard RGB (1996)
- **Gamut**: Smallest (reference = 1.0×)
- **When to use**: Web content, standard displays, maximum compatibility
- **Transfer**: Gamma ~2.2 with linear segment
- **White point**: D65 (6500K daylight)
- **Note**: 96%+ of displays support this

#### Display P3 (2015)
- **Gamut**: 26% wider than sRGB
- **When to use**: Modern Apple devices, high-end monitors, vibrant colors needed
- **Transfer**: Same as sRGB (compatibility)
- **White point**: D65
- **Note**: iPhone X and later, iPad Pro, Mac displays

#### DCI-P3 (2005)
- **Gamut**: Similar to Display P3 but different white point
- **When to use**: Digital cinema, professional video
- **Transfer**: Gamma 2.6
- **White point**: DCI (6300K, slightly warmer)
- **Note**: Same primaries as Display P3, different encoding

#### Adobe RGB 1998
- **Gamut**: 44% wider than sRGB, especially in cyans/greens
- **When to use**: Professional photography, print workflows
- **Transfer**: Simple gamma 2.2
- **White point**: D65
- **Note**: Designed for CMYK print color matching

#### ProPhoto RGB (ROMM RGB)
- **Gamut**: 189% wider than sRGB - encompasses nearly all visible colors
- **When to use**: RAW photo editing, archival, professional color grading
- **Transfer**: Gamma 1.8
- **White point**: D50 (5000K - standard illuminant for print/photography)
- **Note**: Can represent colors invisible to humans (useful for future-proofing)

#### Rec.2020 (BT.2020)
- **Gamut**: 73% wider than sRGB
- **When to use**: UHDTV, HDR video, future displays
- **Transfer**: Gamma 2.4 (simplified)
- **White point**: D65
- **Note**: Few current displays can reproduce full gamut

#### Rec.709 (BT.709)
- **Gamut**: ~Same as sRGB
- **When to use**: HDTV, standard definition video
- **Transfer**: Similar to sRGB
- **White point**: D65
- **Note**: Essentially identical to sRGB in practice

### Perceptually Uniform Spaces

These spaces are designed so that mathematical distance correlates with perceived color difference.

#### OKLCH ⭐ **RECOMMENDED**
- **Structure**: Cylindrical (Lightness, Chroma, Hue)
- **Range**: L[0,1], C[0,~0.4], H[0,360°)
- **When to use**:
  - Generating gradients
  - Lighten/darken operations
  - Palette generation
  - Any time "looks uniform" matters
- **Why it's better**: Fixes perceptual issues in CIELAB
- **Conversion**: OKLCH ↔ OKLAB (lossless) ↔ XYZ

#### OKLAB
- **Structure**: Cartesian (L, a, b)
- **Range**: L[0,1], a/b unbounded (typically ±0.4)
- **When to use**: Mathematical color operations, color difference calculations
- **Note**: Rectangular version of OKLCH; use OKLCH for hue-based operations

#### CIELAB (L*a*b*)
- **Structure**: Cartesian (Lightness, red-green, blue-yellow)
- **Range**: L[0,100], a/b[-128,128] typically
- **When to use**:
  - Industry standard color difference (ΔE)
  - Print/industrial applications
  - Legacy compatibility
- **Note**: Older perceptual space; OKLAB is more accurate

#### CIELCH
- **Structure**: Cylindrical version of CIELAB (L, Chroma, Hue)
- **Range**: L[0,100], C[0,~132], H[0,360°)
- **When to use**: Hue-based operations in LAB space
- **Note**: Cylindrical LAB; prefer OKLCH for new work

#### CIELUV & LCHuv
- **Structure**: Alternative to LAB with different uniformity properties
- **Range**: Similar to LAB/LCH
- **When to use**: Emissive displays (where LAB's print-oriented design is suboptimal)
- **Note**: Less common but better for some display applications

### Intuitive Color Spaces

#### HSL - Hue, Saturation, Lightness
- **Structure**: Cylindrical
- **Range**: H[0,360°), S[0,1], L[0,1]
- **When to use**:
  - UI color pickers
  - Simple hue-based operations
  - When "pure saturation" at L=0.5 is desired
- **Limitations**: Not perceptually uniform; lightness can feel inconsistent

#### HSV/HSB - Hue, Saturation, Value/Brightness
- **Structure**: Cylindrical (cone instead of double-cone)
- **Range**: H[0,360°), S[0,1], V[0,1]
- **When to use**:
  - Photoshop-style color pickers
  - "Brightness" more intuitive than "Lightness"
- **Limitations**: Not perceptually uniform

#### HWB - Hue, Whiteness, Blackness (CSS Color Level 4)
- **Structure**: Cylindrical with whiteness/blackness instead of saturation
- **Range**: H[0,360°), W[0,1], B[0,1]
- **When to use**:
  - More intuitive than HSL for beginners
  - "Add white" and "add black" are easy to understand
- **CSS syntax**: `hwb(180 20% 30%)` = hue 180°, 20% white, 30% black

### Reference Space

#### XYZ (CIE 1931)
- **Purpose**: Mathematical reference connecting all color spaces
- **Not for operations**: Use OKLCH for perceptual operations
- **White point in this library**: D65 (we convert from D50 spaces via chromatic adaptation)
- **Why XYZ**:
  - Device-independent
  - Matches human color matching functions
  - Industry-standard conversion hub

---

## Architecture: How This Library Works

### Single Reference Hub: XYZ (D65)

```
┌─────────────┐
│    sRGB     │──┐
└─────────────┘  │
┌─────────────┐  │      ┌─────────────┐
│ Display P3  │──┤      │             │
└─────────────┘  │      │             │
┌─────────────┐  ├─────▶│   XYZ (D65) │
│   Adobe RGB │──┤      │             │
└─────────────┘  │      │  (Reference  │
┌─────────────┐  │      │     Hub)    │
│ ProPhoto RGB│──┤      │             │
│   (D50→D65) │  │      └─────────────┘
└─────────────┘  │             │
┌─────────────┐  │             │
│  Rec.2020   │──┘             │
└─────────────┘                │
                               │
                    ┌──────────┼──────────┐
                    │          │          │
              ┌──────────┐┌────────┐┌─────────┐
              │  OKLAB   ││  LAB   ││   LUV   │
              └──────────┘└────────┘└─────────┘
                    │          │          │
              ┌──────────┐┌────────┐┌─────────┐
              │  OKLCH   ││  LCH   ││  LCHuv  │
              └──────────┘└────────┘└─────────┘
```

**Key principle**: All color space conversions go through XYZ (D65) as the single reference hub. This ensures:
- **Consistency**: Every conversion path uses the same math
- **Accuracy**: Industry-standard conversion matrices
- **Simplicity**: No need for N² conversion functions

### Perceptual Operations: OKLCH

While XYZ is the conversion hub, **OKLCH is the working space for perceptual operations**:

```go
// When you call Lighten(), this happens internally:
color := sRGB(1.0, 0.0, 0.0) // Red

// 1. Convert to OKLCH
oklch := ToOKLCH(color) // Uses sRGB → XYZ → OKLAB → OKLCH

// 2. Operate in OKLCH (perceptually uniform!)
oklch.L += 0.2 // Lightness increase looks uniform to human eyes

// 3. Convert back
result := oklch.RGBA() // OKLCH → OKLAB → XYZ → sRGB
```

**Why separate conversion hub (XYZ) and working space (OKLCH)?**
- XYZ is mathematically simple and industry-standard for conversions
- OKLCH is perceptually uniform for operations
- Best of both worlds: accurate conversions + perceptual operations

### Chromatic Adaptation: Handling Different White Points

ProPhoto RGB uses D50 white point (5000K, horizon daylight), while most other spaces use D65 (6500K, noon daylight). We handle this automatically:

```go
// ProPhoto RGB color
prophotoColor := NewSpaceColor(ProPhotoRGBSpace, []float64{1, 0, 0}, 1.0)

// Internally when converting to XYZ:
// 1. Convert ProPhoto RGB → XYZ(D50)
// 2. Chromatic adaptation: XYZ(D50) → XYZ(D65) using Bradford transform
// 3. Now in our standard D65 reference

// When converting back:
// 1. XYZ(D65) in
// 2. Chromatic adaptation: XYZ(D65) → XYZ(D50)
// 3. XYZ(D50) → ProPhoto RGB
```

This ensures accurate color appearance across different white points!

---

## Gamut Mapping: Handling Out-of-Gamut Colors

When converting from a wider gamut to a narrower one, some colors won't fit. We provide multiple strategies:

### Gamut Mapping Strategies

#### 1. Clip (Fast, Hue-Shifting)
```go
color := NewSpaceColor(DisplayP3Space, []float64{1.2, 0.3, 0.1}, 1.0)
mapped := MapToGamut(color, GamutClip)
```
- **How it works**: Clips RGB values to [0,1]
- **Pro**: Fast
- **Con**: May shift hue (vivid green becomes more yellow)
- **Use when**: Speed matters, slight color shifts acceptable

#### 2. Preserve Lightness (Recommended)
```go
mapped := MapToGamut(color, GamutPreserveLightness)
```
- **How it works**: Reduces chroma (saturation) while keeping lightness constant
- **Pro**: Maintains perceived brightness
- **Con**: Colors become less saturated
- **Use when**: Brightness is important (UI elements, text backgrounds)

#### 3. Preserve Chroma
```go
mapped := MapToGamut(color, GamutPreserveChroma)
```
- **How it works**: Reduces lightness while keeping chroma constant
- **Pro**: Maintains saturation/"punchiness"
- **Con**: Can make colors significantly darker or lighter
- **Use when**: Saturation is critical (brand colors, accent colors)

#### 4. Project (Best Quality)
```go
mapped := MapToGamut(color, GamutProject)
```
- **How it works**: Finds closest in-gamut color using perceptual distance
- **Pro**: Most accurate to original color intent
- **Con**: Slower (uses search algorithm)
- **Use when**: Quality is paramount (photography, professional color work)

### Visual Example: Gamut Mapping

```
Original (Display P3, out of sRGB gamut):
   L: 0.7, C: 0.25, H: 150° (vivid teal)

Clip:                 L: 0.68, C: 0.22, H: 145° (slightly different hue)
Preserve Lightness:   L: 0.70, C: 0.18, H: 150° (same lightness, less vivid)
Preserve Chroma:      L: 0.65, C: 0.25, H: 150° (darker but same saturation)
Project:              L: 0.69, C: 0.20, H: 149° (perceptually closest overall)
```

---

## Color Difference Metrics

Measuring how different two colors look to human eyes:

### DeltaEOK (Fast, Modern)
```go
diff := DeltaEOK(color1, color2)
// Returns: 0 = identical, <0.02 = imperceptible, <0.05 = barely noticeable
```
- **Formula**: Euclidean distance in OKLAB space
- **Speed**: Very fast
- **Accuracy**: Good for modern applications
- **Use when**: Real-time operations, color matching, palette generation

### DeltaE76 (Classic, Simple)
```go
diff := DeltaE76(color1, color2)
// Returns: 0 = identical, <1 = barely perceptible, 1-2 = small, >2 = noticeable
```
- **Formula**: Euclidean distance in CIELAB space
- **Speed**: Very fast
- **Accuracy**: Acceptable but superseded by DeltaE2000
- **Use when**: Legacy compatibility needed

### DeltaE2000 (Industry Standard)
```go
diff := DeltaE2000(color1, color2)
// Returns: 0 = identical, <1 = imperceptible, 1-2 = small, >2 = noticeable
```
- **Formula**: Complex weighted formula accounting for lightness, chroma, hue
- **Speed**: Slower (many calculations)
- **Accuracy**: Industry standard, matches human perception best
- **Use when**: Professional color matching, quality control, scientific applications

### Choosing a Metric

| Metric | Speed | Accuracy | Use Case |
|--------|-------|----------|----------|
| **DeltaEOK** | ⚡⚡⚡ | ⭐⭐⭐ | Modern apps, real-time, good enough |
| DeltaE76 | ⚡⚡⚡ | ⭐⭐ | Legacy compatibility |
| **DeltaE2000** | ⚡ | ⭐⭐⭐⭐⭐ | Professional, quality-critical |

---

## Gradients: Why Color Space Matters

### The Problem with RGB Gradients

RGB gradients have three issues:

1. **Muddy midpoints**: Red to blue goes through dark purple/brown
2. **Uneven steps**: Some sections appear to have more steps than others
3. **Lightness jumps**: Perceived brightness changes non-uniformly

### Comparison: Same Stops, Different Spaces

<picture>
  <source media="(prefers-color-scheme: dark)"  srcset="docs/gradients/gradient_rgb_white.png">
  <source media="(prefers-color-scheme: light)" srcset="docs/gradients/gradient_rgb_black.png">
  <img alt="RGB Gradient" src="docs/gradients/gradient_rgb_white.png">
</picture>

**RGB**: Notice the muddy middle and how it feels like there are more purple steps than red steps

<picture>
  <source media="(prefers-color-scheme: dark)"  srcset="docs/gradients/gradient_hsl_white.png">
  <source media="(prefers-color-scheme: light)" srcset="docs/gradients/gradient_hsl_black.png">
  <img alt="HSL Gradient" src="docs/gradients/gradient_hsl_white.png">
</picture>

**HSL**: Better hue transition (goes around color wheel) but lightness feels uneven

<picture>
  <source media="(prefers-color-scheme: dark)"  srcset="docs/gradients/gradient_lab_white.png">
  <source media="(prefers-color-scheme: light)" srcset="docs/gradients/gradient_lab_black.png">
  <img alt="LAB Gradient" src="docs/gradients/gradient_lab_white.png">
</picture>

**LAB**: Much more perceptually uniform than RGB/HSL

<picture>
  <source media="(prefers-color-scheme: dark)"  srcset="docs/gradients/gradient_oklch_white.png">
  <source media="(prefers-color-scheme: light)" srcset="docs/gradients/gradient_oklch_black.png">
  <img alt="OKLCH Gradient" src="docs/gradients/gradient_oklch_white.png">
</picture>

**OKLCH ⭐**: Smoothest, most uniform - each step appears equally different to human eyes

### Code Example

```go
red := color.RGB(1, 0, 0)
blue := color.RGB(0, 0, 1)

// Bad: Muddy middle
rgbGradient := color.GradientInSpace(red, blue, 10, color.GradientRGB)

// Better: Cleaner but lightness varies
hslGradient := color.GradientInSpace(red, blue, 10, color.GradientHSL)

// Best: Perceptually uniform
oklchGradient := color.Gradient(red, blue, 10) // Uses OKLCH by default
```

### When to Use Each Space for Gradients

| Space | Best For | Avoid For |
|-------|----------|-----------|
| **OKLCH** | Everything (default recommendation) | N/A |
| **OKLAB** | When you want to interpolate in rectangular coordinates | Hue-based transitions |
| LAB/LCH | Legacy compatibility, scientific contexts | Modern applications (use OKLCH) |
| HSL | Quick hue-based transitions, UI pickers | Professional/quality-critical gradients |
| RGB | Performance-critical code where quality doesn't matter | Anything user-facing |

---

## Practical Examples

### Example 1: Creating a UI Color Palette

```go
// Start with your brand color
brandColor := color.NewOKLCH(0.55, 0.15, 230, 1.0) // Blue

// Generate a full palette with proper perceptual spacing
palette := map[string]Color{
    "50":  color.Lighten(brandColor, 0.45),  // Lightest
    "100": color.Lighten(brandColor, 0.35),
    "200": color.Lighten(brandColor, 0.25),
    "300": color.Lighten(brandColor, 0.15),
    "400": color.Lighten(brandColor, 0.05),
    "500": brandColor,                        // Base
    "600": color.Darken(brandColor, 0.05),
    "700": color.Darken(brandColor, 0.15),
    "800": color.Darken(brandColor, 0.25),
    "900": color.Darken(brandColor, 0.35),    // Darkest
}

// Or use a gradient for even steps
shades := color.Gradient(
    color.Lighten(brandColor, 0.45),
    color.Darken(brandColor, 0.35),
    10,
)
```

### Example 2: Wide-Gamut Photography Workflow

```go
// RAW photo in ProPhoto RGB (widest gamut)
rawColor, _ := color.ConvertFromRGBSpace(0.9, 0.2, 0.1, 1.0, "prophoto-rgb")

// Edit in perceptually uniform space
oklch := color.ToOKLCH(rawColor)
adjusted := color.Saturate(oklch, 0.15)
adjusted = color.Lighten(adjusted, 0.05)

// Export for web: convert to sRGB with proper gamut mapping
webColor := color.MapToGamut(adjusted, color.GamutProject) // Best quality mapping
hex := color.RGBToHex(webColor)

// Export for modern displays: convert to Display P3
p3Color, _ := color.ConvertToRGBSpace(adjusted, "display-p3")
// Still preserves more color than sRGB!
```

### Example 3: Checking if Colors Are Distinguishable

```go
color1 := color.ParseColor("#FF6B6B")
color2 := color.ParseColor("#FF6D6C")

// Check if humans can tell them apart
diff := color.DeltaEOK(color1, color2)

if diff < 0.02 {
    fmt.Println("These colors look identical")
} else if diff < 0.05 {
    fmt.Println("Barely distinguishable")
} else if diff < 0.1 {
    fmt.Println("Small difference")
} else {
    fmt.Println("Clearly different")
}
```

---

## Key Takeaways

1. **Use OKLCH for operations**: Lighten, darken, saturate, gradients - always OKLCH unless you have a specific reason not to

2. **XYZ is the conversion hub**: All color space conversions go through XYZ (D65) for accuracy and consistency

3. **Preserve wide-gamut colors**: Don't convert to sRGB until you absolutely have to display/export

4. **Choose the right gamut mapping**: PreserveLightness for most uses, Project for quality-critical work

5. **Measure perceptual difference**: Use DeltaEOK for speed, DeltaE2000 for accuracy

6. **Chromatic adaptation is automatic**: ProPhoto RGB (D50) colors are properly adapted to D65 and back

## When to Use This Library

✅ **Yes, use this library when:**
- You need perceptually uniform color operations
- You're working with wide-gamut displays (Display P3, HDR)
- You're generating gradients that need to look smooth
- You need accurate color space conversions
- You're building professional color tools
- You want color difference metrics

❌ **Consider alternatives when:**
- You only need basic RGB operations
- Performance is critical and perceptual uniformity doesn't matter
- You're only working with 8-bit sRGB
- You don't care about color science accuracy

---

## Further Reading

- [Oklab color space](https://bottosson.github.io/posts/oklab/) - Björn Ottosson's excellent writeup
- [CSS Color Module Level 4](https://www.w3.org/TR/css-color-4/) - Modern web color standards
- [Color Appearance Models](https://www.wiley.com/en-us/Color+Appearance+Models%2C+3rd+Edition-p-9781119967033) - Comprehensive color science
- [Understanding Color Spaces](https://www.cambridgeincolour.com/tutorials/color-spaces.htm) - Accessible introduction
