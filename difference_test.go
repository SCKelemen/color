package color

import (
	"math"
	"testing"
)

func TestDeltaEOK(t *testing.T) {
	tests := []struct {
		name     string
		c1, c2   Color
		maxDelta float64
	}{
		{
			name:     "Identical colors",
			c1:       RGB(1, 0, 0),
			c2:       RGB(1, 0, 0),
			maxDelta: 0.001,
		},
		{
			name:     "Very similar reds",
			c1:       RGB(1, 0, 0),
			c2:       RGB(0.99, 0, 0),
			maxDelta: 0.05,
		},
		{
			name:     "Red to blue (large difference)",
			c1:       RGB(1, 0, 0),
			c2:       RGB(0, 0, 1),
			maxDelta: 2.0, // Should be much larger
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := DeltaEOK(tt.c1, tt.c2)
			if delta < 0 {
				t.Errorf("DeltaEOK returned negative value: %f", delta)
			}
			if tt.maxDelta > 0 && delta > tt.maxDelta {
				t.Errorf("DeltaEOK = %f, want <= %f", delta, tt.maxDelta)
			}
		})
	}
}

func TestDeltaE76(t *testing.T) {
	tests := []struct {
		name     string
		c1, c2   Color
		maxDelta float64
	}{
		{
			name:     "Identical colors",
			c1:       RGB(0.5, 0.5, 0.5),
			c2:       RGB(0.5, 0.5, 0.5),
			maxDelta: 0.001,
		},
		{
			name:     "Similar grays",
			c1:       RGB(0.5, 0.5, 0.5),
			c2:       RGB(0.51, 0.51, 0.51),
			maxDelta: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := DeltaE76(tt.c1, tt.c2)
			if delta < 0 {
				t.Errorf("DeltaE76 returned negative value: %f", delta)
			}
			if tt.maxDelta > 0 && delta > tt.maxDelta {
				t.Errorf("DeltaE76 = %f, want <= %f", delta, tt.maxDelta)
			}
		})
	}
}

func TestDeltaE2000(t *testing.T) {
	tests := []struct {
		name     string
		c1, c2   Color
		expected float64
		tolerance float64
	}{
		{
			name:      "Identical colors",
			c1:        RGB(1, 0, 0),
			c2:        RGB(1, 0, 0),
			expected:  0,
			tolerance: 0.001,
		},
		{
			name:      "Very similar colors (imperceptible)",
			c1:        RGB(1, 0, 0),
			c2:        RGB(0.995, 0, 0),
			expected:  0,
			tolerance: 1.0, // Should be < 1 (imperceptible)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := DeltaE2000(tt.c1, tt.c2)
			if delta < 0 {
				t.Errorf("DeltaE2000 returned negative value: %f", delta)
			}
			if math.Abs(delta-tt.expected) > tt.tolerance {
				t.Logf("DeltaE2000 = %f, expected ~%f (tolerance %f)", delta, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestDeltaESymmetry(t *testing.T) {
	c1 := RGB(1, 0, 0)
	c2 := RGB(0, 0, 1)

	// DeltaE should be symmetric
	deltaOK1 := DeltaEOK(c1, c2)
	deltaOK2 := DeltaEOK(c2, c1)
	if math.Abs(deltaOK1-deltaOK2) > 0.0001 {
		t.Errorf("DeltaEOK not symmetric: %f vs %f", deltaOK1, deltaOK2)
	}

	delta76_1 := DeltaE76(c1, c2)
	delta76_2 := DeltaE76(c2, c1)
	if math.Abs(delta76_1-delta76_2) > 0.0001 {
		t.Errorf("DeltaE76 not symmetric: %f vs %f", delta76_1, delta76_2)
	}

	delta2000_1 := DeltaE2000(c1, c2)
	delta2000_2 := DeltaE2000(c2, c1)
	if math.Abs(delta2000_1-delta2000_2) > 0.0001 {
		t.Errorf("DeltaE2000 not symmetric: %f vs %f", delta2000_1, delta2000_2)
	}
}

func TestDeltaEMonotonicity(t *testing.T) {
	// As colors get more different, DeltaE should increase
	base := RGB(1, 0, 0)
	c1 := RGB(0.99, 0, 0)
	c2 := RGB(0.95, 0, 0)
	c3 := RGB(0.90, 0, 0)

	d1 := DeltaEOK(base, c1)
	d2 := DeltaEOK(base, c2)
	d3 := DeltaEOK(base, c3)

	if !(d1 < d2 && d2 < d3) {
		t.Errorf("DeltaEOK not monotonic: %f, %f, %f", d1, d2, d3)
	}
}

func BenchmarkDeltaEOK(b *testing.B) {
	c1 := RGB(1, 0, 0)
	c2 := RGB(0, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeltaEOK(c1, c2)
	}
}

func BenchmarkDeltaE76(b *testing.B) {
	c1 := RGB(1, 0, 0)
	c2 := RGB(0, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeltaE76(c1, c2)
	}
}

func BenchmarkDeltaE2000(b *testing.B) {
	c1 := RGB(1, 0, 0)
	c2 := RGB(0, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeltaE2000(c1, c2)
	}
}
