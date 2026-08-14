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

// Under RemoveNAs(false) a null in a group must poison EVERY aggregate to NaN,
// holistic (Median) as well as reducible (Mean).
func TestAggregateSerialPoisonAllAggregates(t *testing.T) {
	ctx := enchanter.NewContext()
	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString([]string{"a", "a", "a"}, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64([]float64{1, 0, 3}, []bool{false, true, false}, false, ctx)).(BaseDataFrame)

	out := aggregateSerial(df, []series.Series{df.C("g")}, []aggregator{Mean("v"), Median("v")}, false)
	if got := f64(t, out, "mean(v)")[0]; !math.IsNaN(got) {
		t.Fatalf("poison mean = %v, want NaN", got)
	}
	if got := f64(t, out, "median(v)")[0]; !math.IsNaN(got) {
		t.Fatalf("poison median = %v, want NaN", got)
	}
}

// A nullable key column: the null-key group sorts LAST (after every non-null
// group), which come out ascending.
func TestAggregateSerialNullKeysSortLast(t *testing.T) {
	ctx := enchanter.NewContext()
	// keys: b, <null>, a, <null>, b  ->  groups b, <null>, a
	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString([]string{"b", "", "a", "", "b"}, []bool{false, true, false, true, false}, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64([]float64{1, 2, 3, 4, 5}, nil, false, ctx)).(BaseDataFrame)

	out := aggregateSerial(df, []series.Series{df.C("g")}, []aggregator{Sum("v")}, true)
	g := out.C("g").(series.Strings)
	if g.IsNull(0) || g.GetAsString(0) != "a" {
		t.Fatalf("row0 = %q null=%v, want a", g.GetAsString(0), g.IsNull(0))
	}
	if g.IsNull(1) || g.GetAsString(1) != "b" {
		t.Fatalf("row1 = %q null=%v, want b", g.GetAsString(1), g.IsNull(1))
	}
	if !g.IsNull(2) {
		t.Fatalf("row2 null=%v, want null last", g.IsNull(2))
	}
	// sums follow the sorted rows: a=3, b=1+5=6, null=2+4=6
	sum := f64(t, out, "sum(v)")
	if sum[0] != 3 || sum[1] != 6 || sum[2] != 6 {
		t.Fatalf("sum = %v, want [3 6 6]", sum)
	}
}

// Two key columns: result rows are ordered lexicographically by (k1, k2) with
// k1 primary; rows sharing k1 but differing on k2 exercise the tie-break.
func TestAggregateSerialTwoKeyLexicographic(t *testing.T) {
	ctx := enchanter.NewContext()
	// (k1,k2): (b,y),(a,y),(b,x),(a,x) -> 4 distinct groups
	df := NewBaseDataFrame(ctx).
		AddSeries("k1", series.NewSeriesString([]string{"b", "a", "b", "a"}, nil, false, ctx)).
		AddSeries("k2", series.NewSeriesString([]string{"y", "y", "x", "x"}, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64([]float64{1, 2, 3, 4}, nil, false, ctx)).(BaseDataFrame)

	out := aggregateSerial(df, []series.Series{df.C("k1"), df.C("k2")}, []aggregator{Sum("v")}, true)
	k1 := out.C("k1").(series.Strings)
	k2 := out.C("k2").(series.Strings)
	// lexicographic: (a,x),(a,y),(b,x),(b,y)
	wantK1 := []string{"a", "a", "b", "b"}
	wantK2 := []string{"x", "y", "x", "y"}
	for i := range wantK1 {
		if k1.GetAsString(i) != wantK1[i] || k2.GetAsString(i) != wantK2[i] {
			t.Fatalf("row%d = (%q,%q), want (%q,%q)", i, k1.GetAsString(i), k2.GetAsString(i), wantK1[i], wantK2[i])
		}
	}
}
