package series

import (
	"context"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/compute"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

// ArrowFloat64s is a MEASUREMENT-ONLY prototype of an Arrow-native Float64
// series: it stores an arrow.Array as the source of truth instead of a Go
// slice. It exists solely to benchmark the "arrow-native" storage regime
// against the current Go-slice path (see arrow_bench_test.go). It is NOT wired
// into the code generator or the DataFrame dispatch, implements only the
// read/compute methods the benchmarks need, and treats the array as immutable
// (mutation is out of scope; see docs/superpowers/specs for the 0.3.0 design).
type ArrowFloat64s struct {
	arr   *array.Float64
	alloc memory.Allocator
}

// NewArrowFloat64s builds an Arrow-backed Float64 series from a Go slice.
func NewArrowFloat64s(data []float64, alloc memory.Allocator) ArrowFloat64s {
	if alloc == nil {
		alloc = memory.DefaultAllocator
	}
	b := array.NewFloat64Builder(alloc)
	defer b.Release()
	b.AppendValues(data, nil)
	return ArrowFloat64s{arr: b.NewFloat64Array(), alloc: alloc}
}

func (s ArrowFloat64s) Len() int { return s.arr.Len() }

func (s ArrowFloat64s) Get(i int) float64 { return s.arr.Value(i) }

// ArrowArray returns the underlying Arrow array (borrowed, not retained).
func (s ArrowFloat64s) ArrowArray() arrow.Array { return s.arr }

// Release frees the underlying Arrow array.
func (s ArrowFloat64s) Release() { s.arr.Release() }

// Add returns s + other element-wise, computed on the Arrow arrays.
func (s ArrowFloat64s) Add(other ArrowFloat64s) ArrowFloat64s {
	ld := &compute.ArrayDatum{Value: s.arr.Data()}
	rd := &compute.ArrayDatum{Value: other.arr.Data()}
	res, err := compute.Add(context.Background(), compute.ArithmeticOptions{}, ld, rd)
	if err != nil {
		panic(err)
	}
	out := res.(*compute.ArrayDatum).MakeArray().(*array.Float64)
	return ArrowFloat64s{arr: out, alloc: s.alloc}
}
