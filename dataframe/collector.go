package dataframe

import (
	"math"
	"sort"
)

type collector interface {
	collect(v float64)
	merge(other collector)
	result() (float64, bool)
}

// quantileOf computes the p-quantile of a SORTED non-empty slice using the
// given interpolation (numpy/pandas method set). p is clamped to [0,1].
func quantileOf(sorted []float64, p float64, interp Interpolation) float64 {
	n := len(sorted)
	if n == 1 {
		return sorted[0]
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[n-1]
	}
	rank := p * float64(n-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	frac := rank - float64(lo)
	switch interp {
	case Lower:
		return sorted[lo]
	case Higher:
		return sorted[hi]
	case Nearest:
		// round half to even
		r := math.RoundToEven(rank)
		return sorted[int(r)]
	case Midpoint:
		return (sorted[lo] + sorted[hi]) / 2
	default: // Linear
		return sorted[lo] + (sorted[hi]-sorted[lo])*frac
	}
}

type quantileCollector struct {
	vals   []float64
	p      float64
	interp Interpolation
}

func newQuantileCollector(p float64, interp Interpolation) collector {
	return &quantileCollector{p: p, interp: interp}
}
func (c *quantileCollector) collect(v float64) { c.vals = append(c.vals, v) }
func (c *quantileCollector) merge(o collector) {
	c.vals = append(c.vals, o.(*quantileCollector).vals...)
}
func (c *quantileCollector) result() (float64, bool) {
	if len(c.vals) == 0 {
		return 0, true
	}
	sort.Float64s(c.vals)
	return quantileOf(c.vals, c.p, c.interp), false
}
