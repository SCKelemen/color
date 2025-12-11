package color

import "math"

// OKLCHSpace represents the OKLCH color space (perceptually uniform)
var OKLCHSpace Space = &oklchSpace{}

// oklchSpace implements Space for OKLCH
type oklchSpace struct{}

func (s *oklchSpace) Name() string {
	return "OKLCH"
}

func (s *oklchSpace) Channels() int {
	return 3
}

func (s *oklchSpace) ChannelNames() []string {
	return []string{"L", "C", "H"}
}

func (s *oklchSpace) ToXYZ(channels []float64) (x, y, z float64) {
	if len(channels) != 3 {
		panic("OKLCH space requires 3 channels")
	}
	
	okl, c, h := channels[0], channels[1], channels[2]
	
	// Convert OKLCH to OKLAB
	rad := h * math.Pi / 180
	oka := c * math.Cos(rad)
	okb := c * math.Sin(rad)
	
	// Convert OKLAB to XYZ
	return oklabToXYZ(okl, oka, okb)
}

func (s *oklchSpace) FromXYZ(x, y, z float64) []float64 {
	// Convert XYZ to OKLAB
	okl, oka, okb := xyzToOKLAB(x, y, z)
	
	// Convert OKLAB to OKLCH
	c := math.Sqrt(oka*oka + okb*okb)
	h := math.Atan2(okb, oka) * 180 / math.Pi
	h = normalizeHue(h)
	
	return []float64{okl, c, h}
}

// oklabToXYZ converts OKLAB to XYZ
func oklabToXYZ(okl, oka, okb float64) (x, y, z float64) {
	// Convert OKLAB to linear sRGB
	l_ := okl + 0.3963377774*oka + 0.2158037573*okb
	m_ := okl - 0.1055613458*oka - 0.0638541728*okb
	s_ := okl - 0.0894841775*oka - 1.2914855480*okb
	
	l_ = l_ * l_ * l_
	m_ = m_ * m_ * m_
	s_ = s_ * s_ * s_
	
	// Convert LMS to linear RGB
	r := +4.0767416621*l_ - 3.3077115913*m_ + 0.2309699292*s_
	g := -1.2684380046*l_ + 2.6097574011*m_ - 0.3413193965*s_
	b := -0.0041960863*l_ - 0.7034186147*m_ + 1.7076147010*s_
	
	// Convert linear RGB to XYZ
	x = r*0.4124564 + g*0.3575761 + b*0.1804375
	y = r*0.2126729 + g*0.7151522 + b*0.0721750
	z = r*0.0193339 + g*0.1191920 + b*0.9503041
	
	return x, y, z
}

// xyzToOKLAB converts XYZ to OKLAB
func xyzToOKLAB(x, y, z float64) (okl, oka, okb float64) {
	// Convert XYZ to linear RGB
	r := x*3.2404542 - y*1.5371385 - z*0.4985314
	g := -x*0.9692660 + y*1.8760108 + z*0.0415560
	b := x*0.0556434 - y*0.2040259 + z*1.0572252
	
	// Convert linear RGB to LMS
	l_ := 0.4122214708*r + 0.5363325363*g + 0.0514459929*b
	m_ := 0.2119034982*r + 0.6806995451*g + 0.1073969566*b
	s_ := 0.0883024619*r + 0.2817188376*g + 0.6299787005*b
	
	// Apply cube root
	l_ = math.Cbrt(l_)
	m_ = math.Cbrt(m_)
	s_ = math.Cbrt(s_)
	
	// Convert LMS to OKLAB
	okl = 0.2104542553*l_ + 0.7936177850*m_ - 0.0040720468*s_
	oka = 1.9779984951*l_ - 2.4285922050*m_ + 0.4505937099*s_
	okb = 0.0259040371*l_ + 0.7827717662*m_ - 0.8086757660*s_
	
	return okl, oka, okb
}

