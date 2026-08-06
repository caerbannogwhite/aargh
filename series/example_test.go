package series_test

import (
	"fmt"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/series"
)

func ExampleNewSeriesInt64() {
	ctx := enchanter.NewContext()

	a := series.NewSeriesInt64([]int64{1, 2, 3}, nil, false, ctx)
	b := series.NewSeriesInt64([]int64{10, 20, 30}, nil, false, ctx)

	sum := a.Add(b)
	fmt.Println(sum.Data().([]int64))
	// Output:
	// [11 22 33]
}

func ExampleNewSeriesFloat64() {
	ctx := enchanter.NewContext()

	// The second element is null; nulls propagate through operations.
	s := series.NewSeriesFloat64([]float64{1.5, 0, 2.5}, []bool{false, true, false}, false, ctx)
	doubled := s.Mul(2.0)

	fmt.Println(doubled.IsNull(0), doubled.IsNull(1), doubled.IsNull(2))
	data := doubled.Data().([]float64)
	fmt.Println(data[0], data[2])
	// Output:
	// false true false
	// 3 5
}
