package dataframe

import (
	"math"
	"testing"
)

func TestQuantileOf(t *testing.T) {
	s := []float64{1, 2, 3, 4} // ranks 0..3
	cases := []struct {
		p      float64
		interp Interpolation
		want   float64
	}{
		{0.5, Linear, 2.5}, // 0.5*(n-1)=1.5 -> between 2 and 3
		{0.5, Lower, 2},
		{0.5, Higher, 3},
		{0.5, Nearest, 3}, // rank 1.5 rounds half-to-even to 2 -> sorted[2]=3
		{0.5, Midpoint, 2.5},
		{0.0, Linear, 1},
		{1.0, Linear, 4},
		{0.25, Linear, 1.75},
	}
	for _, c := range cases {
		if got := quantileOf(s, c.p, c.interp); math.Abs(got-c.want) > 1e-9 {
			t.Fatalf("quantileOf(p=%v,%v) = %v, want %v", c.p, c.interp, got, c.want)
		}
	}
	// odd count: median is the middle element (input must already be sorted)
	if got := quantileOf([]float64{1, 3, 5}, 0.5, Linear); got != 3 {
		t.Fatalf("odd-count median = %v, want 3", got)
	}
}

func TestQuantileCollector(t *testing.T) {
	c := newQuantileCollector(0.5, Linear)
	for _, v := range []float64{4, 1, 3, 2} { // unsorted on purpose
		c.collect(v)
	}
	v, null := c.result()
	if null || math.Abs(v-2.5) > 1e-9 {
		t.Fatalf("median = %v null=%v, want 2.5", v, null)
	}
	// merge two halves, result independent of split
	a := newQuantileCollector(0.5, Linear)
	b := newQuantileCollector(0.5, Linear)
	a.collect(4)
	a.collect(1)
	b.collect(3)
	b.collect(2)
	a.merge(b)
	av, _ := a.result()
	if math.Abs(av-2.5) > 1e-9 {
		t.Fatalf("merged median = %v, want 2.5", av)
	}
	// empty -> null
	if _, n := newQuantileCollector(0.5, Linear).result(); !n {
		t.Fatalf("empty collector should be null")
	}
}
