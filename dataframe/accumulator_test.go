package dataframe

import (
	"math"
	"testing"
)

func feed(a accumulator, vs ...float64) accumulator {
	for _, v := range vs {
		a.update(v)
	}
	return a
}

func result(t *testing.T, a accumulator) float64 {
	v, null := a.result()
	if null {
		t.Fatalf("unexpected null result")
	}
	return v
}

func TestReducibleAccumulators(t *testing.T) {
	if got := result(t, feed(newSumAcc(), 1, 2, 3, 4)); got != 10 {
		t.Fatalf("sum = %v, want 10", got)
	}
	if got := result(t, feed(newMeanAcc(), 2, 4, 6)); got != 4 {
		t.Fatalf("mean = %v, want 4", got)
	}
	if got := result(t, feed(newMinAcc(), 3, -1, 2)); got != -1 {
		t.Fatalf("min = %v, want -1", got)
	}
	if got := result(t, feed(newMaxAcc(), 3, -1, 2)); got != 3 {
		t.Fatalf("max = %v, want 3", got)
	}
	if _, null := newCountAcc().result(); null {
		t.Fatalf("empty count should be 0, not null")
	}
	if got := result(t, feed(newCountAcc(), 5, 5, 5)); got != 3 {
		t.Fatalf("count = %v, want 3", got)
	}
	// population variance of [2,4,4,4,5,5,7,9] = 4, std = 2
	data := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	if got := result(t, feed(newVarAcc(0, false), data...)); math.Abs(got-4) > 1e-9 {
		t.Fatalf("pop var = %v, want 4", got)
	}
	if got := result(t, feed(newVarAcc(0, true), data...)); math.Abs(got-2) > 1e-9 {
		t.Fatalf("pop std = %v, want 2", got)
	}
	// sample variance of the same = 32/7 ≈ 4.5714
	if got := result(t, feed(newVarAcc(1, false), data...)); math.Abs(got-32.0/7.0) > 1e-9 {
		t.Fatalf("sample var = %v, want 32/7", got)
	}
	// sample std of a 1-element group → null (n <= ddof)
	if _, null := feed(newVarAcc(1, true), 42).result(); !null {
		t.Fatalf("sample std of n=1 should be null")
	}
}

func TestAccumulatorMerge(t *testing.T) {
	// Welford merge parity: split [2,4,4,4,5,5,7,9] and merge halves.
	full := feed(newVarAcc(0, false), 2, 4, 4, 4, 5, 5, 7, 9)
	a := feed(newVarAcc(0, false), 2, 4, 4, 4)
	b := feed(newVarAcc(0, false), 5, 5, 7, 9)
	a.merge(b)
	fv, _ := full.result()
	av, _ := a.result()
	if math.Abs(fv-av) > 1e-9 {
		t.Fatalf("merged var = %v, want %v", av, fv)
	}
}
