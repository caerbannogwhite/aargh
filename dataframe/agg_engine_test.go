package dataframe

import (
	"math"
	"testing"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/series"
)

func f64(t *testing.T, df DataFrame, col string) []float64 {
	s, ok := df.C(col).(series.Float64s)
	if !ok {
		t.Fatalf("col %q not Float64s", col)
	}
	return s.Data_
}

func TestAggregateSerialSingleKey(t *testing.T) {
	ctx := enchanter.NewContext()
	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString([]string{"b", "a", "b", "a", "b"}, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64([]float64{1, 10, 2, 20, 3}, nil, false, ctx)).(BaseDataFrame)

	out := aggregateSerial(df, []series.Series{df.C("g")}, []aggregator{Sum("v"), Count()}, true)
	// sorted by key: a then b
	if got := out.C("g").(series.Strings).GetAsString(0); got != "a" {
		t.Fatalf("row0 key = %q, want a (sorted)", got)
	}
	sum := f64(t, out, "sum(v)")
	if sum[0] != 30 || sum[1] != 6 { // a: 10+20, b: 1+2+3
		t.Fatalf("sum = %v, want [30 6]", sum)
	}
	n := out.C("n").(series.Int64s).Data_
	if n[0] != 2 || n[1] != 3 {
		t.Fatalf("count = %v, want [2 3]", n)
	}
}

func TestAggregateSerialSkipNullsDefault(t *testing.T) {
	ctx := enchanter.NewContext()
	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString([]string{"a", "a", "a"}, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64([]float64{1, 0, 3}, []bool{false, true, false}, false, ctx)).(BaseDataFrame)

	// removeNAs=true (new default): mean of {1,3} = 2
	out := aggregateSerial(df, []series.Series{df.C("g")}, []aggregator{Mean("v")}, true)
	if got := f64(t, out, "mean(v)")[0]; got != 2 {
		t.Fatalf("skip-null mean = %v, want 2", got)
	}
	// removeNAs=false: poisoned -> NaN
	out2 := aggregateSerial(df, []series.Series{df.C("g")}, []aggregator{Mean("v")}, false)
	if got := f64(t, out2, "mean(v)")[0]; !math.IsNaN(got) {
		t.Fatalf("poison mean = %v, want NaN", got)
	}
}

func TestAggregateSerialMedianMultiKey(t *testing.T) {
	ctx := enchanter.NewContext()
	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString([]string{"a", "a", "a", "a"}, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64([]float64{1, 2, 3, 4}, nil, false, ctx)).(BaseDataFrame)
	out := aggregateSerial(df, []series.Series{df.C("g")}, []aggregator{Median("v"), Quantile("v", 0.25, WithInterpolation(Higher))}, true)
	if got := f64(t, out, "median(v)")[0]; got != 2.5 {
		t.Fatalf("median = %v, want 2.5", got)
	}
	if got := f64(t, out, "quantile_0.25(v)")[0]; got != 2 {
		t.Fatalf("q0.25 higher = %v, want 2", got)
	}
}
