# Color Library Documentation

Complete documentation for the Go color manipulation library with perceptually uniform operations and professional color space support.

## Quick Navigation

### 🚀 Getting Started
- **[Main README](../README.md)** - Project overview and feature highlights
- **[Quickstart Guide](../QUICKSTART.md)** - Get up and running in 5 minutes

### 📚 Guides (How-To)
Practical guides for common tasks:
- **[Gradient Generation](guides/gradients.md)** - Create perceptually uniform gradients
- **[LOG Color Workflows](guides/log-workflows.md)** - Cinema camera LOG color spaces
- **[Lipgloss Integration](guides/lipgloss-integration.md)** - Terminal UI styling integration

### 📖 Reference
Technical documentation and API reference:
- **[API Reference](reference/api-overview.md)** - Complete API documentation
- **[Color Space Reference](reference/color-space-list.md)** - All supported color spaces
- **[Format Support](reference/format-support.md)** - CSS color format coverage
- **[Architecture](reference/architecture.md)** - Color space system design

### 🎓 Theory & Concepts
Educational content and explanations:
- **[Color Primer](theory/color-primer.md)** - Color science fundamentals
- **[Why Use This Library](theory/why-use-this.md)** - Decision guide and use cases
- **[Visual Comparisons](theory/visual-comparisons.md)** - See the difference

### 🛠️ Contributing
For maintainers and contributors:
- **[Release Process](contributing/release-process.md)** - How to publish releases
- **[LOG Implementation Notes](contributing/log-implementation-notes.md)** - Implementation planning

---

## Documentation by User Type

### For Beginners
1. Start with [Main README](../README.md) for an overview
2. Follow [Quickstart Guide](../QUICKSTART.md) for hands-on introduction
3. Read [Color Primer](theory/color-primer.md) to understand color theory
4. Explore [Why Use This Library](theory/why-use-this.md) for use cases

### For Developers Building Applications
1. Review [API Reference](reference/api-overview.md) for available functions
2. Check [Color Space Reference](reference/color-space-list.md) for space support
3. Learn [Gradient Generation](guides/gradients.md) for smooth color transitions
4. See [Format Support](reference/format-support.md) for parsing capabilities

### For Professional Video/Film
1. Read [LOG Color Workflows](guides/log-workflows.md) for cinema camera support
2. Review supported LOG formats and workflows
3. Check [Architecture](reference/architecture.md) for technical details

### For Terminal UI Developers
1. Follow [Lipgloss Integration](guides/lipgloss-integration.md) guide
2. Learn dynamic theme generation
3. Explore color manipulation for terminal styling

### For Contributors
1. Review [Architecture](reference/architecture.md) for system design
2. Follow [Release Process](contributing/release-process.md) for publishing
3. Check implementation notes for specific features

---

## Key Features at a Glance

### Perceptually Uniform Operations
```go
blue := color.RGB(0, 0, 1)
lighter := color.Lighten(blue, 0.2)  // Actually looks 20% lighter!
```

### Smooth Gradients
```go
red := color.RGB(1, 0, 0)
blue := color.RGB(0, 0, 1)
gradient := color.Gradient(red, blue, 20)  // Perceptually smooth
```

### Wide-Gamut Support
```go
// Display P3 color
p3Color, _ := color.ConvertFromRGBSpace(1, 0, 0, 1, "display-p3")
result := color.Lighten(p3Color, 0.2)  // Preserves wide gamut!
```

### Professional LOG Support
```go
// Cinema camera LOG footage
slog3 := color.NewSpaceColor(color.SLog3Space, []float64{0.41, 0.39, 0.35}, 1.0)
hdr := slog3.ConvertTo(color.Rec2020Space)  // HDR mastering
```

---

## Popular Topics

### Color Manipulation
- [Lighten/Darken](../QUICKSTART.md#perceptually-uniform-operations)
- [Saturate/Desaturate](../QUICKSTART.md#perceptually-uniform-operations)
- [Hue Adjustment](../README.md#color-manipulation)
- [Color Mixing](../README.md#color-manipulation)

### Color Spaces
- [sRGB, Display P3, Rec.2020](reference/color-space-list.md)
- [OKLCH, OKLAB, CIELAB](theory/color-primer.md)
- [LOG Spaces (Cinema)](guides/log-workflows.md)

### Gradients
- [Basic Gradients](guides/gradients.md#basic-gradients)
- [Multi-Stop Gradients](guides/gradients.md#multi-stop-gradients)
- [Easing Functions](guides/gradients.md#easing-functions)
- [Color Space Selection](guides/gradients.md#choosing-color-spaces)

### Gamut Mapping
- [What is Gamut Mapping](theory/color-primer.md#gamut-and-gamut-mapping)
- [Mapping Strategies](reference/architecture.md#gamut-mapping)
- [Preserve Lightness vs Chroma](../README.md#professional-gamut-mapping)

---

## External Resources

- **[Go Package Documentation](https://pkg.go.dev/github.com/SCKelemen/color)** - Complete godoc reference
- **[GitHub Repository](https://github.com/SCKelemen/color)** - Source code and issues
- **[Examples](https://pkg.go.dev/github.com/SCKelemen/color#pkg-examples)** - Executable code examples

---

## Need Help?

- **Can't find what you need?** Check the [Main README](../README.md) FAQ section
- **Found a bug?** [Open an issue](https://github.com/SCKelemen/color/issues)
- **Want to contribute?** See [Contributing Documentation](contributing/)

---

## Documentation Status

Last updated: December 2025
Library version: 1.0+
Go version: 1.22+
