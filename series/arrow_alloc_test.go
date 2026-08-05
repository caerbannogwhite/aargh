package series

import (
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/caerbannogwhite/aargh"
)

// Leak tests: run the Arrow interop and compute paths under a CheckedAllocator
// and assert that every internal Arrow reference is released. Values returned
// to callers (Series.ArrowArray) are released explicitly here, per the
// documented ownership contract.

func newCheckedCtx() (*aargh.Context, *memory.CheckedAllocator) {
	mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
	ctx := aargh.NewContext()
	ctx.Allocator = mem
	return ctx, mem
}

func TestArrowInteropNoLeaks(t *testing.T) {
	ctx, mem := newCheckedCtx()
	defer mem.AssertSize(t, 0)

	mask := []bool{false, true, false, false, true}

	roundTrip := func(s Series) {
		t.Helper()
		arr := s.ArrowArray()
		got := ArrowArrayToSeries(arr, ctx)
		arr.Release()
		if got.IsError() {
			t.Fatal(got.GetError())
		}
	}

	roundTrip(NewSeriesFloat64([]float64{1, 2, 3, 4, 5}, mask, false, ctx))
	roundTrip(NewSeriesInt64([]int64{1, 2, 3, 4, 5}, mask, false, ctx))
	roundTrip(NewSeriesInt([]int{1, 2, 3, 4, 5}, mask, false, ctx))
	roundTrip(NewSeriesBool([]bool{true, false, true, false, true}, mask, false, ctx))
	roundTrip(NewSeriesString([]string{"a", "b", "c", "d", "e"}, mask, false, ctx))
	roundTrip(NewSeriesTime(make([]time.Time, 5), mask, false, ctx))
	roundTrip(NewSeriesDuration(make([]time.Duration, 5), mask, false, ctx))
}

func TestArrowOpsNoLeaks(t *testing.T) {
	ctx, mem := newCheckedCtx()
	defer mem.AssertSize(t, 0)

	a := NewSeriesFloat64([]float64{1, 2, 3, 4}, []bool{false, true, false, false}, false, ctx)
	b := NewSeriesFloat64([]float64{10, 20, 30, 40}, nil, false, ctx)

	if res := ArrowAdd(a, b, ctx); res.IsError() {
		t.Fatal(res.GetError())
	}
	if res := ArrowMul(a, a, ctx); res.IsError() {
		t.Fatal(res.GetError())
	}
	if res := ArrowLt(a, b, ctx); res.IsError() {
		t.Fatal(res.GetError())
	}

	// Scalar broadcast path (single-element series -> ScalarDatum).
	two := NewSeriesFloat64([]float64{2}, nil, false, ctx)
	if res := ArrowMul(a, two, ctx); res.IsError() {
		t.Fatal(res.GetError())
	}

	p := NewSeriesBool([]bool{true, false, true, false}, nil, false, ctx)
	q := NewSeriesBool([]bool{true, true, false, false}, []bool{false, false, true, false}, false, ctx)
	if res := ArrowAnd(p, q, ctx); res.IsError() {
		t.Fatal(res.GetError())
	}
	if res := ArrowOr(p, q, ctx); res.IsError() {
		t.Fatal(res.GetError())
	}
}
