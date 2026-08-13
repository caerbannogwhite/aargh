package dataframe

import "math"

type accumulator interface {
	update(v float64)
	merge(other accumulator)
	result() (float64, bool)
}

type countAcc struct{ n int64 }

func newCountAcc() accumulator              { return &countAcc{} }
func (a *countAcc) update(float64)          { a.n++ }
func (a *countAcc) merge(o accumulator)     { a.n += o.(*countAcc).n }
func (a *countAcc) result() (float64, bool) { return float64(a.n), false }

type sumAcc struct {
	s float64
	n int64
}

func newSumAcc() accumulator { return &sumAcc{} }
func (a *sumAcc) update(v float64) {
	a.s += v
	a.n++
}
func (a *sumAcc) merge(o accumulator) {
	b := o.(*sumAcc)
	a.s += b.s
	a.n += b.n
}
func (a *sumAcc) result() (float64, bool) {
	if a.n == 0 {
		return 0, true
	}
	return a.s, false
}

type meanAcc struct {
	s float64
	n int64
}

func newMeanAcc() accumulator { return &meanAcc{} }
func (a *meanAcc) update(v float64) {
	a.s += v
	a.n++
}
func (a *meanAcc) merge(o accumulator) {
	b := o.(*meanAcc)
	a.s += b.s
	a.n += b.n
}
func (a *meanAcc) result() (float64, bool) {
	if a.n == 0 {
		return 0, true
	}
	return a.s / float64(a.n), false
}

type minAcc struct {
	v   float64
	set bool
}

func newMinAcc() accumulator { return &minAcc{} }
func (a *minAcc) update(v float64) {
	if !a.set || v < a.v {
		a.v = v
		a.set = true
	}
}
func (a *minAcc) merge(o accumulator) {
	b := o.(*minAcc)
	if b.set {
		a.update(b.v)
	}
}
func (a *minAcc) result() (float64, bool) { return a.v, !a.set }

type maxAcc struct {
	v   float64
	set bool
}

func newMaxAcc() accumulator { return &maxAcc{} }
func (a *maxAcc) update(v float64) {
	if !a.set || v > a.v {
		a.v = v
		a.set = true
	}
}
func (a *maxAcc) merge(o accumulator) {
	b := o.(*maxAcc)
	if b.set {
		a.update(b.v)
	}
}
func (a *maxAcc) result() (float64, bool) { return a.v, !a.set }

// varAcc: Welford running mean/M2, parallel-mergeable (Chan et al.).
type varAcc struct {
	n    float64
	mean float64
	m2   float64
	ddof int
	std  bool
}

func newVarAcc(ddof int, std bool) accumulator { return &varAcc{ddof: ddof, std: std} }
func (a *varAcc) update(v float64) {
	a.n++
	d := v - a.mean
	a.mean += d / a.n
	a.m2 += d * (v - a.mean)
}
func (a *varAcc) merge(o accumulator) {
	b := o.(*varAcc)
	if b.n == 0 {
		return
	}
	if a.n == 0 {
		a.n, a.mean, a.m2 = b.n, b.mean, b.m2
		return
	}
	delta := b.mean - a.mean
	tot := a.n + b.n
	a.m2 += b.m2 + delta*delta*a.n*b.n/tot
	a.mean += delta * b.n / tot
	a.n = tot
}
func (a *varAcc) result() (float64, bool) {
	denom := a.n - float64(a.ddof)
	if a.n == 0 || denom <= 0 {
		return 0, true
	}
	v := a.m2 / denom
	if a.std {
		return math.Sqrt(v), false
	}
	return v, false
}

func newReducibleAcc(t AggregateType, ddof int) accumulator {
	switch t {
	case AGGREGATE_COUNT:
		return newCountAcc()
	case AGGREGATE_SUM:
		return newSumAcc()
	case AGGREGATE_MEAN:
		return newMeanAcc()
	case AGGREGATE_MIN:
		return newMinAcc()
	case AGGREGATE_MAX:
		return newMaxAcc()
	case AGGREGATE_STD:
		return newVarAcc(ddof, true)
	case AGGREGATE_VARIANCE:
		return newVarAcc(ddof, false)
	}
	panic("newReducibleAcc: not a reducible aggregate")
}
