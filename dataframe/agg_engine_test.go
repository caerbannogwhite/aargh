package dataframe

import (
	"math"
	"runtime"
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

// TestAggregateParallelParity asserts that aggregate's parallel path (row
// count above aggMinParallel) produces results identical to aggregateSerial
// across reducible (Sum, Mean, Min, Max, Std) and non-value (Count)
// aggregators. Run with -race to confirm the worker chunks + merge carry no
// data races.
func TestAggregateParallelParity(t *testing.T) {
	ctx := enchanter.NewContext()
	n := 200_000
	keys := make([]string, n)
	vals := make([]float64, n)
	for i := range keys {
		keys[i] = []string{"a", "b", "c", "d"}[i%4]
		vals[i] = float64(i % 97)
	}
	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString(keys, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64(vals, nil, false, ctx)).(BaseDataFrame)

	aggs := []aggregator{Sum("v"), Mean("v"), Min("v"), Max("v"), Std("v"), Count()}
	ser := aggregateSerial(df, []series.Series{df.C("g")}, aggs, true)
	par := aggregate(df, []series.Series{df.C("g")}, aggs, true)

	for _, col := range []string{"sum(v)", "mean(v)", "min(v)", "max(v)", "std(v)"} {
		s, p := f64(t, ser, col), f64(t, par, col)
		if len(s) != len(p) {
			t.Fatalf("%s len mismatch", col)
		}
		for i := range s {
			if math.Abs(s[i]-p[i]) > 1e-9 {
				t.Fatalf("%s[%d] serial=%v parallel=%v", col, i, s[i], p[i])
			}
		}
	}
}

// TestAggregateParallelPoisonCrossChunk exercises the cross-chunk poison
// OR-merge (agg_engine.go's merge loop, which must OR a group's poison flag
// across every worker chunk rather than reset it) and the holistic
// (Median)/Count columns under the parallel path, neither of which
// TestAggregateParallelParity exercises (it only runs removeNAs=true, which
// never poisons anything).
//
// GOMAXPROCS is forced to 4 for the duration of the test so the null row and
// this group's many non-null rows are guaranteed to land in different worker
// chunks regardless of the host's actual CPU count: group "a" appears at
// every 4th row across the full 200,000-row range (50,000 rows total, spread
// evenly across all 4 forced chunks of 50,000 rows each), with a single null
// at row 0 (chunk 0) and ~49,999 non-null rows for "a" spread across every
// chunk, including chunks 1-3 which see no poison locally at all.
func TestAggregateParallelPoisonCrossChunk(t *testing.T) {
	ctx := enchanter.NewContext()
	n := 200_000
	keys := make([]string, n)
	vals := make([]float64, n)
	nulls := make([]bool, n)
	for i := range keys {
		keys[i] = []string{"a", "b", "c", "d"}[i%4]
		vals[i] = float64(i % 97)
	}
	nulls[0] = true // row 0 is key "a"; this is the only null in the frame

	df := NewBaseDataFrame(ctx).
		AddSeries("g", series.NewSeriesString(keys, nil, false, ctx)).
		AddSeries("v", series.NewSeriesFloat64(vals, nulls, false, ctx)).(BaseDataFrame)

	oldProcs := runtime.GOMAXPROCS(4)
	defer runtime.GOMAXPROCS(oldProcs)

	aggs := []aggregator{Sum("v"), Mean("v"), Median("v"), Count()}
	ser := aggregateSerial(df, []series.Series{df.C("g")}, aggs, false)
	par := aggregate(df, []series.Series{df.C("g")}, aggs, false)

	if ser.NRows() != par.NRows() {
		t.Fatalf("row count mismatch: serial=%d parallel=%d", ser.NRows(), par.NRows())
	}

	gcolPar := par.C("g").(series.Strings)
	aIdx := -1
	for i := 0; i < par.NRows(); i++ {
		if gcolPar.GetAsString(i) == "a" {
			aIdx = i
			break
		}
	}
	if aIdx == -1 {
		t.Fatalf("group 'a' not found in parallel result")
	}

	// The poisoned group's reducible (Sum, Mean) and holistic (Median)
	// aggregates must all be non-null NaN, matching the serial engine's
	// poison-overrides-everything semantics.
	for _, col := range []string{"sum(v)", "mean(v)", "median(v)"} {
		s := par.C(col).(series.Float64s)
		if s.IsNull(aIdx) {
			t.Fatalf("parallel %s[a] isNull = true, want non-null NaN (poison)", col)
		}
		if !math.IsNaN(s.Data_[aIdx]) {
			t.Fatalf("parallel %s[a] = %v, want NaN (poisoned)", col, s.Data_[aIdx])
		}
	}

	// Count is unaffected by value nulls: it counts rows, not values.
	wantCount := int64(n / 4)
	cnt := par.C("n").(series.Int64s)
	if cnt.Data_[aIdx] != wantCount {
		t.Fatalf("parallel count[a] = %d, want %d", cnt.Data_[aIdx], wantCount)
	}

	// Parallel must match serial exactly for the SAME inputs: same keys in
	// the same order, and every aggregate column equal (NaN-aware for the
	// float columns; exact for Count).
	gcolSer := ser.C("g").(series.Strings)
	for i := 0; i < ser.NRows(); i++ {
		if gcolSer.GetAsString(i) != gcolPar.GetAsString(i) {
			t.Fatalf("key[%d] serial=%q parallel=%q", i, gcolSer.GetAsString(i), gcolPar.GetAsString(i))
		}
	}
	for _, col := range []string{"sum(v)", "mean(v)", "median(v)"} {
		sSer := ser.C(col).(series.Float64s)
		sPar := par.C(col).(series.Float64s)
		for i := 0; i < ser.NRows(); i++ {
			if sSer.IsNull(i) != sPar.IsNull(i) {
				t.Fatalf("%s[%d] null mask mismatch: serial=%v parallel=%v", col, i, sSer.IsNull(i), sPar.IsNull(i))
			}
			sv, pv := sSer.Data_[i], sPar.Data_[i]
			// NaN-aware equality: NaN != NaN in Go, but poisoned cells on
			// both sides must both be NaN for the rows to agree.
			if math.IsNaN(sv) || math.IsNaN(pv) {
				if !math.IsNaN(sv) || !math.IsNaN(pv) {
					t.Fatalf("%s[%d] NaN mismatch: serial=%v parallel=%v", col, i, sv, pv)
				}
				continue
			}
			if math.Abs(sv-pv) > 1e-9 {
				t.Fatalf("%s[%d] serial=%v parallel=%v", col, i, sv, pv)
			}
		}
	}
	nSer := ser.C("n").(series.Int64s)
	nPar := par.C("n").(series.Int64s)
	for i := 0; i < ser.NRows(); i++ {
		if nSer.Data_[i] != nPar.Data_[i] {
			t.Fatalf("n[%d] serial=%d parallel=%d", i, nSer.Data_[i], nPar.Data_[i])
		}
	}
}
