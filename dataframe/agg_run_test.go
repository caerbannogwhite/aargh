package dataframe

import (
	"math"
	"testing"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/series"
)

func TestAggRunSortedAndSkipNullByDefault(t *testing.T) {
	ctx := enchanter.NewContext()
	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString([]string{"b", "a", "b"}, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64([]float64{1, 5, 0}, []bool{false, false, true}, false, ctx))

	out := df.GroupBy("g").Agg(Sum("v"), Std("v", WithDDoF(1))).Run()
	if out.IsErrored() {
		t.Fatalf("agg errored: %v", out.GetError())
	}
	// sorted: a, b
	if out.C("g").(series.Strings).GetAsString(0) != "a" {
		t.Fatalf("not sorted by key")
	}
	// b's null value skipped by default: sum(b) = 1, sample std of single value -> null
	sum := out.C("sum(v)").(series.Float64s)
	if sum.Data_[1] != 1 {
		t.Fatalf("skip-null sum(b) = %v, want 1", sum.Data_[1])
	}
	std := out.C("std(v)").(series.Float64s)
	if !std.IsNull(1) {
		t.Fatalf("sample std of single non-null should be null")
	}
}

func TestAggRunOptionValidation(t *testing.T) {
	ctx := enchanter.NewContext()
	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString([]string{"a"}, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64([]float64{1}, nil, false, ctx))
	out := df.GroupBy("g").Agg(Sum("v", WithDDoF(1))).Run() // ddof on Sum → error
	if !out.IsErrored() {
		t.Fatalf("expected error for WithDDoF on Sum")
	}
	_ = math.Inf
}
