# Color Space Reference

Complete reference for all supported color spaces in the library.

## Color Space Categories

- **[RGB Spaces](#rgb-color-spaces)** - Device-dependent, display-oriented
- **[LOG Spaces](#log-color-spaces-cinema)** - Cinema cameras, HDR workflows
- **[Perceptual Spaces](#perceptual-color-spaces)** - Uniform, device-independent
- **[Intuitive Spaces](#intuitive-color-spaces)** - Human-friendly (HSL, HSV, HWB)
- **[Reference Spaces](#reference-color-spaces)** - Interchange (XYZ)

---

## RGB Color Spaces

Device-dependent color spaces for displays and monitors.

### sRGB

**Name:** `srgb`
**Variable:** `color.SRGBSpace`
**Gamut:** Standard (baseline)
**HDR Support:** No
**White Point:** D65

The standard color space for the web and most displays. All colors on the web are assumed to be sRGB unless otherwise specified.

**Use Cases:**
- Web development
- Standard displays
- Default color space

**Gamut:** 1.0× (baseline)

```go
// sRGB is the default for RGB() function
srgb := color.RGB(1, 0, 0)

// Or explicit space
srgb := color.NewSpaceColor(color.SRGBSpace,
    []float64{1, 0, 0}, 1.0)
```

### sRGB Linear

**Name:** `srgb-linear`
**Variable:** `color.SRGBLinearSpace`
**Gamut:** Same as sRGB
**HDR Support:** Yes
**White Point:** D65

Linear light version of sRGB (no gamma encoding). Used for physically correct blending and compositing.

**Use Cases:**
- Light calculations
- Physically correct blending
- HDR values > 1.0

**Gamut:** 1.0× (same primaries as sRGB)

```go
linear := color.NewSpaceColor(color.SRGBLinearSpace,
    []float64{1, 0, 0}, 1.0)
```

### Display P3

**Name:** `display-p3`
**Aliases:** `display-p3-d65`
**Variable:** `color.DisplayP3Space`
**Gamut:** Wide (26% more colors)
**HDR Support:** No
**White Point:** D65

Apple's wide-gamut color space used in iPhone X+, iPad Pro, Mac displays, and modern monitors.

**Use Cases:**
- Modern Apple devices
- Wide-gamut displays
- Vibrant colors

**Gamut:** 1.26× wider than sRGB

```go
p3, _ := color.ConvertFromRGBSpace(1, 0, 0, 1, "display-p3")

// Or direct
p3 := color.NewSpaceColor(color.DisplayP3Space,
    []float64{1, 0, 0}, 1.0)
```

### DCI-P3

**Name:** `dci-p3`
**Aliases:** `dci-p3-d65`
**Variable:** `color.DCIP3Space`
**Gamut:** Wide (26% more colors)
**HDR Support:** No
**White Point:** D65 (adapted from DCI white point)

Digital cinema color space, similar to Display P3 but with different white point adaptation.

**Use Cases:**
- Digital cinema
- Film post-production
- Cinema projection

**Gamut:** 1.26× wider than sRGB

```go
dci := color.NewSpaceColor(color.DCIP3Space,
    []float64{1, 0, 0}, 1.0)
```

### Adobe RGB 1998

**Name:** `a98-rgb`
**Aliases:** `a98rgb`, `adobe-rgb-1998`
**Variable:** `color.A98RGBSpace`
**Gamut:** Wide (44% more colors)
**HDR Support:** No
**White Point:** D65

Professional photography color space with wider gamut than sRGB, especially in cyan/green.

**Use Cases:**
- Professional photography
- Print workflows
- CMYK conversion

**Gamut:** 1.44× wider than sRGB

```go
adobe := color.NewSpaceColor(color.A98RGBSpace,
    []float64{1, 0, 0}, 1.0)
```

### ProPhoto RGB

**Name:** `prophoto-rgb`
**Aliases:** `prophoto`
**Variable:** `color.ProPhotoRGBSpace`
**Gamut:** Very wide (189% more colors)
**HDR Support:** No
**White Point:** D50

Extremely wide-gamut space for RAW photo editing. Contains colors beyond human vision.

**Use Cases:**
- RAW photo editing
- Professional photography
- Archival workflows

**Gamut:** 2.89× wider than sRGB

```go
prophoto := color.NewSpaceColor(color.ProPhotoRGBSpace,
    []float64{1, 0, 0}, 1.0)
```

### Rec.2020

**Name:** `rec2020`
**Aliases:** `rec-2020`
**Variable:** `color.Rec2020Space`
**Gamut:** Wide (73% more colors)
**HDR Support:** No (but used for HDR)
**White Point:** D65

UHDTV and HDR color space. Standard for 4K/8K television and HDR content.

**Use Cases:**
- UHDTV broadcasting
- HDR mastering
- Future-proof content

**Gamut:** 1.73× wider than sRGB

```go
rec2020 := color.NewSpaceColor(color.Rec2020Space,
    []float64{1, 0, 0}, 1.0)
```

### Rec.709

**Name:** `rec709`
**Aliases:** `rec-709`
**Variable:** `color.Rec709Space`
**Gamut:** Standard (same as sRGB)
**HDR Support:** No
**White Point:** D65

HDTV color space with same primaries as sRGB but different transfer function.

**Use Cases:**
- HDTV broadcasting
- Video production
- Broadcast standards

**Gamut:** 1.0× (same primaries as sRGB)

```go
rec709 := color.NewSpaceColor(color.Rec709Space,
    []float64{1, 0, 0}, 1.0)
```

---

## LOG Color Spaces (Cinema)

Professional logarithmic color spaces for cinema cameras with HDR support.

### Canon C-Log

**Name:** `c-log`
**Aliases:** `clog`
**Variable:** `color.CLogSpace`
**Gamut:** Cinema Gamut (56% more colors)
**HDR Support:** Yes
**White Point:** D65
**Dynamic Range:** ~12-14 stops

Canon's LOG curve for Cinema EOS cameras (C300, C500, C700, R5 C).

**Key Characteristics:**
- 18% gray → ~0.34 encoded
- Wide Cinema Gamut primaries
- Good highlight rolloff

**Cameras:** C300, C500, C700, C70, R5 C, R7 C

```go
clog := color.NewSpaceColor(color.CLogSpace,
    []float64{0.45, 0.40, 0.35}, 1.0)
```

### Sony S-Log3

**Name:** `s-log3`
**Aliases:** `slog3`
**Variable:** `color.SLog3Space`
**Gamut:** S-Gamut3 (69% more colors)
**HDR Support:** Yes
**White Point:** D65
**Dynamic Range:** ~14+ stops

Sony's latest LOG curve optimized for HDR workflows.

**Key Characteristics:**
- 18% gray → ~0.41 encoded (41 IRE)
- S-Gamut3 primaries (wider than Rec.2020)
- Excellent for HDR

**Cameras:** FX6, FX9, FX3, Venice, BURANO, A7S III, A1

```go
slog3 := color.NewSpaceColor(color.SLog3Space,
    []float64{0.41, 0.39, 0.35}, 1.0)
```

### Panasonic V-Log

**Name:** `v-log`
**Aliases:** `vlog`
**Variable:** `color.VLogSpace`
**Gamut:** V-Gamut (58% more colors)
**HDR Support:** Yes
**White Point:** D65
**Dynamic Range:** ~14 stops

Panasonic's LOG curve for cinema and hybrid cameras.

**Key Characteristics:**
- 18% gray → ~0.42 encoded
- V-Gamut primaries
- Clean shadows

**Cameras:** GH5, GH6, S1H, S5 II, EVA1, Varicam

```go
vlog := color.NewSpaceColor(color.VLogSpace,
    []float64{0.48, 0.45, 0.40}, 1.0)
```

### Arri LogC

**Name:** `arri-logc`
**Aliases:** `logc`
**Variable:** `color.ArriLogCSpace`
**Gamut:** Arri Wide Gamut (55% more colors)
**HDR Support:** Yes
**White Point:** D65
**Dynamic Range:** ~14 stops

Industry-standard LOG curve from Arri (LogC V3, EI 800).

**Key Characteristics:**
- 18% gray → ~0.38-0.39 encoded
- Arri Wide Gamut primaries
- Excellent highlight handling

**Cameras:** Alexa Mini, Alexa LF, Alexa 35, Amira

```go
logc := color.NewSpaceColor(color.ArriLogCSpace,
    []float64{0.42, 0.38, 0.35}, 1.0)
```

### Red Log3G10

**Name:** `red-log3g10`
**Aliases:** `log3g10`
**Variable:** `color.RedLog3G10Space`
**Gamut:** RedWideGamutRGB (68% more colors)
**HDR Support:** Yes
**White Point:** D65
**Dynamic Range:** ~16+ stops

Red's LOG curve with 10-bit encoding optimization.

**Key Characteristics:**
- Log3G10 (10-bit optimized)
- RedWideGamutRGB primaries
- Massive dynamic range

**Cameras:** Komodo, V-Raptor, Ranger, DSMC2

```go
redlog := color.NewSpaceColor(color.RedLog3G10Space,
    []float64{0.46, 0.42, 0.38}, 1.0)
```

### Blackmagic Film

**Name:** `bmd-film`
**Aliases:** `bmdfilm`
**Variable:** `color.BMDFilmSpace`
**Gamut:** Wide (70% more colors, ~Rec.2020)
**HDR Support:** Yes
**White Point:** D65
**Dynamic Range:** ~13 stops

Blackmagic's film-like LOG curve for cinema cameras.

**Key Characteristics:**
- Wide gamut similar to Rec.2020
- Popular in independent filmmaking
- Good balance of DR and gradability

**Cameras:** Pocket Cinema 4K/6K, URSA Mini Pro

```go
bmdfilm := color.NewSpaceColor(color.BMDFilmSpace,
    []float64{0.44, 0.40, 0.36}, 1.0)
```

---

## Perceptual Color Spaces

Device-independent spaces designed for perceptual uniformity.

### OKLCH

**Name:** `oklch`
**Variable:** `color.OKLCHSpace`
**Type:** Cylindrical (Lightness, Chroma, Hue)
**Perceptually Uniform:** Yes ⭐
**White Point:** D65

Modern, perceptually uniform color space. **Recommended** for all color manipulation.

**Why Use OKLCH:**
- Best perceptual uniformity
- Intuitive cylindrical coordinates
- Modern CSS support
- Hue-based operations

**Channels:**
- L: Lightness (0-1)
- C: Chroma/saturation (0-0.4 typical)
- H: Hue (0-360°)

```go
oklch := color.NewOKLCH(0.7, 0.2, 180, 1.0)

// Convert to OKLCH for manipulation
oklch := color.ToOKLCH(anyColor)
```

### OKLAB

**Name:** `oklab`
**Type:** Rectangular (Lightness, a, b)
**Perceptually Uniform:** Yes ⭐
**White Point:** D65

Rectangular version of OKLCH. Good for blending and direct manipulation.

**Channels:**
- L: Lightness (0-1)
- a: Green-red axis
- b: Blue-yellow axis

```go
oklab := color.NewOKLAB(0.7, 0.1, -0.1, 1.0)
oklab := color.ToOKLAB(anyColor)
```

### CIELAB

**Type:** Rectangular
**Perceptually Uniform:** Yes
**White Point:** D65 (adapted)

Industry-standard perceptual color space. Older than OKLAB but widely supported.

```go
lab := color.ToLAB(anyColor)
```

### CIELCH

**Type:** Cylindrical
**Perceptually Uniform:** Yes
**White Point:** D65

Cylindrical version of CIELAB.

```go
lch := color.ToLCH(anyColor)
```

### CIELUV

**Type:** Rectangular
**Perceptually Uniform:** Yes (for emissive displays)
**White Point:** D65

Alternative to CIELAB, optimized for emissive displays.

```go
luv := color.ToLUV(anyColor)
```

### LCHuv

**Type:** Cylindrical
**Perceptually Uniform:** Yes
**White Point:** D65

Cylindrical version of CIELUV.

```go
lchuv := color.ToLCHuv(anyColor)
```

---

## Intuitive Color Spaces

Human-friendly color spaces (not perceptually uniform).

### HSL

**Type:** Cylindrical (Hue, Saturation, Lightness)
**Perceptually Uniform:** No
**Based on:** sRGB

Common color picker model. Not perceptually uniform but intuitive.

**Channels:**
- H: Hue (0-360°)
- S: Saturation (0-1)
- L: Lightness (0-1)

```go
hsl := color.NewHSL(180, 0.5, 0.5, 1.0)
hsl := color.ToHSL(anyColor)
```

### HSV/HSB

**Type:** Cylindrical (Hue, Saturation, Value/Brightness)
**Perceptually Uniform:** No
**Based on:** sRGB

Alternative to HSL. Value represents brightness differently.

**Channels:**
- H: Hue (0-360°)
- S: Saturation (0-1)
- V: Value/Brightness (0-1)

```go
hsv := color.NewHSV(180, 0.5, 0.8, 1.0)
hsv := color.ToHSV(anyColor)
```

### HWB

**Type:** Cylindrical (Hue, Whiteness, Blackness)
**Perceptually Uniform:** No
**Based on:** sRGB

CSS Color Level 4 format. Often more intuitive than HSL/HSV.

**Channels:**
- H: Hue (0-360°)
- W: Whiteness (0-1)
- B: Blackness (0-1)

```go
hwb := color.NewHWB(180, 0.2, 0.1, 1.0)
```

---

## Reference Color Spaces

### XYZ (CIE 1931)

**Type:** Tristimulus
**Purpose:** Conversion hub
**White Point:** D65

Device-independent reference space. All color conversions go through XYZ.

**Not typically used directly**, but important for understanding conversions.

```go
xyz := color.ToXYZ(anyColor)
```

---

## Space Selection Guide

### For Color Manipulation
**Use:** OKLCH (best perceptual uniformity)
```go
oklch := color.ToOKLCH(myColor)
lighter := color.Lighten(oklch, 0.2)
```

### For Gradients
**Use:** OKLCH or OKLAB (smooth, perceptually uniform)
```go
gradient := color.Gradient(start, end, 20)  // Uses OKLCH internally
```

### For Display
**Use:** sRGB (web), Display P3 (modern devices), Rec.2020 (HDR)
```go
srgb := color.RGB(1, 0, 0)
p3, _ := color.ConvertFromRGBSpace(1, 0, 0, 1, "display-p3")
```

### For Cinema/Video
**Use:** LOG spaces (camera-dependent), Rec.709 (HDTV), Rec.2020 (UHD/HDR)
```go
slog3 := color.NewSpaceColor(color.SLog3Space, vals, 1.0)
rec2020 := slog3.ConvertTo(color.Rec2020Space)
```

### For Color Pickers
**Use:** HSL, HSV, or HWB (intuitive but not perceptual)
```go
hsl := color.NewHSL(180, 0.5, 0.5, 1.0)
```

---

## Gamut Comparison Table

| Space | Gamut vs sRGB | Colors Beyond sRGB | Primary Use |
|-------|---------------|-------------------|-------------|
| sRGB | 1.0× | Baseline | Web, standard displays |
| Display P3 | 1.26× | +26% | Modern Apple devices |
| DCI-P3 | 1.26× | +26% | Digital cinema |
| Adobe RGB | 1.44× | +44% | Photography, print |
| Rec.709 | 1.0× | Same primaries | HDTV broadcast |
| Rec.2020 | 1.73× | +73% | UHDTV, HDR |
| ProPhoto RGB | 2.89× | +189% | RAW editing |
| C-Log (Cinema Gamut) | 1.56× | +56% | Canon cinema |
| S-Log3 (S-Gamut3) | 1.69× | +69% | Sony cinema |
| V-Log (V-Gamut) | 1.58× | +58% | Panasonic cinema |
| Arri LogC (AWG) | 1.55× | +55% | Arri cinema |
| Red Log3G10 | 1.68× | +68% | Red cinema |
| BMD Film | 1.70× | +70% | Blackmagic cinema |

---

## See Also

- **[API Overview](api-overview.md)** - Function reference
- **[Color Primer](../theory/color-primer.md)** - Color theory fundamentals
- **[LOG Workflows](../guides/log-workflows.md)** - Cinema camera usage
