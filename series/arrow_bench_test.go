package series

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/caerbannogwhite/enchanter"
)

func benchSizes(short bool) []int {
	if short {
		return []int{1e4, 1e5}
	}
	return []int{1e4, 1e5, 1e6, 1e7}
}

func makeF64(n int) []float64 {
	d := make([]float64, n)
	for i := range d {
		d[i] = float64(i%1000) - 500
	}
	return d
}

// ---- Add ----

func BenchmarkAdd(b *testing.B) {
	ctx := enchanter.NewContext()
	for _, n := range benchSizes(testing.Short()) {
		da, db := makeF64(n), makeF64(n)

		b.Run("goslice/"+itoa(n), func(b *testing.B) {
			left := NewSeriesFloat64(da, nil, false, ctx)
			right := NewSeriesFloat64(db, nil, false, ctx)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = left.Add(right)
			}
		})

		b.Run("roundtrip/"+itoa(n), func(b *testing.B) {
			// Data lives in Go slices; each op builds Arrow, computes, materializes.
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				l := NewArrowFloat64s(da, memory.DefaultAllocator)
				r := NewArrowFloat64s(db, memory.DefaultAllocator)
				out := l.Add(r)
				_ = out.arr.Float64Values() // materialize back to Go
				out.Release()
				l.Release()
				r.Release()
			}
		})

		b.Run("arrownative/"+itoa(n), func(b *testing.B) {
			// Data already Arrow; result stays Arrow. Build cost excluded.
			l := NewArrowFloat64s(da, memory.DefaultAllocator)
			r := NewArrowFloat64s(db, memory.DefaultAllocator)
			defer l.Release()
			defer r.Release()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				out := l.Add(r)
				out.Release()
			}
		})
	}
}

// ---- 2-op chain: (a+b) > 0, then filter ----

func BenchmarkChain(b *testing.B) {
	ctx := enchanter.NewContext()
	for _, n := range benchSizes(testing.Short()) {
		da, db := makeF64(n), makeF64(n)

		b.Run("goslice/"+itoa(n), func(b *testing.B) {
			left := NewSeriesFloat64(da, nil, false, ctx)
			right := NewSeriesFloat64(db, nil, false, ctx)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sum := left.Add(right)
				mask := sum.Gt(0.0)
				_ = sum.Filter(mask)
			}
		})

		b.Run("arrownative/"+itoa(n), func(b *testing.B) {
			l := NewArrowFloat64s(da, memory.DefaultAllocator)
			r := NewArrowFloat64s(db, memory.DefaultAllocator)
			defer l.Release()
			defer r.Release()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sum := l.Add(r)
				mask := sum.GreaterThan(0)
				out := sum.Filter(mask)
				out.Release()
				mask.Release()
				sum.Release()
			}
		})
	}
}

// ---- Sum ----

// goSumFloat64 is the genuine Go-slice baseline for the Sum benchmark: a plain
// loop over the []float64, i.e. what a user would write. The library has no
// scalar Float64s.Sum method, so this loop (not a stub) is the honest go-slice
// counterpart to the arrow-native SIMD sum (arrmath.Float64.Sum) it races.
func goSumFloat64(d []float64) float64 {
	var s float64
	for _, v := range d {
		s += v
	}
	return s
}

func BenchmarkSum(b *testing.B) {
	for _, n := range benchSizes(testing.Short()) {
		da := makeF64(n)

		b.Run("goslice/"+itoa(n), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = goSumFloat64(da)
			}
		})

		b.Run("arrownative/"+itoa(n), func(b *testing.B) {
			s := NewArrowFloat64s(da, memory.DefaultAllocator)
			defer s.Release()
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = s.Sum()
			}
		})
	}
}

// small local itoa to avoid strconv noise in bench names
func itoa(n int) string {
	switch {
	case n >= 1e7:
		return "1e7"
	case n >= 1e6:
		return "1e6"
	case n >= 1e5:
		return "1e5"
	default:
		return "1e4"
	}
}
