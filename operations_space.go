package color

// LightenSpace increases the lightness of a color in its native space.
// For perceptually uniform spaces (OKLCH), this works directly.
// For other spaces, converts to OKLCH, operates, and converts back.
func LightenSpace(c SpaceColor, amount float64) SpaceColor {
	amount = clamp01(amount)
	
	space := c.Space()
	
	// If it's OKLCH, operate directly
	if space.Name() == "OKLCH" {
		channels := c.Channels()
		l := channels[0]
		l = clamp01(l + amount*(1-l))
		newChannels := []float64{l, channels[1], channels[2]}
		return NewSpaceColor(space, newChannels, c.Alpha())
	}
	
	// For other spaces, convert to OKLCH, operate, convert back
	oklchColor := c.ConvertTo(OKLCHSpace)
	oklchChannels := oklchColor.Channels()
	oklchChannels[0] = clamp01(oklchChannels[0] + amount*(1-oklchChannels[0]))
	lightenedOKLCH := NewSpaceColor(OKLCHSpace, oklchChannels, c.Alpha())
	return lightenedOKLCH.ConvertTo(space)
}

// DarkenSpace decreases the lightness of a color in its native space.
func DarkenSpace(c SpaceColor, amount float64) SpaceColor {
	amount = clamp01(amount)
	
	space := c.Space()
	
	// If it's OKLCH, operate directly
	if space.Name() == "OKLCH" {
		channels := c.Channels()
		l := channels[0]
		l = clamp01(l * (1 - amount))
		newChannels := []float64{l, channels[1], channels[2]}
		return NewSpaceColor(space, newChannels, c.Alpha())
	}
	
	// For other spaces, convert to OKLCH, operate, convert back
	oklchColor := c.ConvertTo(OKLCHSpace)
	oklchChannels := oklchColor.Channels()
	oklchChannels[0] = clamp01(oklchChannels[0] * (1 - amount))
	darkenedOKLCH := NewSpaceColor(OKLCHSpace, oklchChannels, c.Alpha())
	return darkenedOKLCH.ConvertTo(space)
}

// SaturateSpace increases the saturation of a color in its native space.
func SaturateSpace(c SpaceColor, amount float64) SpaceColor {
	amount = clamp01(amount)
	
	space := c.Space()
	
	// If it's OKLCH, operate directly
	if space.Name() == "OKLCH" {
		channels := c.Channels()
		l, c_val, h := channels[0], channels[1], channels[2]
		maxC := estimateMaxChroma(l, h)
		c_val = clamp(c_val+amount*(maxC-c_val), 0, maxC)
		newChannels := []float64{l, c_val, h}
		return NewSpaceColor(space, newChannels, c.Alpha())
	}
	
	// For other spaces, convert to OKLCH, operate, convert back
	oklchColor := c.ConvertTo(OKLCHSpace)
	oklchChannels := oklchColor.Channels()
	l, c_val, h := oklchChannels[0], oklchChannels[1], oklchChannels[2]
	maxC := estimateMaxChroma(l, h)
	c_val = clamp(c_val+amount*(maxC-c_val), 0, maxC)
	saturatedOKLCH := NewSpaceColor(OKLCHSpace, []float64{l, c_val, h}, c.Alpha())
	return saturatedOKLCH.ConvertTo(space)
}

// DesaturateSpace decreases the saturation of a color in its native space.
func DesaturateSpace(c SpaceColor, amount float64) SpaceColor {
	amount = clamp01(amount)
	
	space := c.Space()
	
	// If it's OKLCH, operate directly
	if space.Name() == "OKLCH" {
		channels := c.Channels()
		c_val := channels[1]
		c_val = c_val * (1 - amount)
		newChannels := []float64{channels[0], c_val, channels[2]}
		return NewSpaceColor(space, newChannels, c.Alpha())
	}
	
	// For other spaces, convert to OKLCH, operate, convert back
	oklchColor := c.ConvertTo(OKLCHSpace)
	oklchChannels := oklchColor.Channels()
	oklchChannels[1] = oklchChannels[1] * (1 - amount)
	desaturatedOKLCH := NewSpaceColor(OKLCHSpace, oklchChannels, c.Alpha())
	return desaturatedOKLCH.ConvertTo(space)
}

