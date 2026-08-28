package dataframe

import (
	"math"
	"testing"
)

// TestAccumulatorSoAParity guards the deliberate duplication between
// accumulator.go's accumulator interface (countAcc/sumAcc/meanAcc/minAcc/
// maxAcc/varAcc — update/merge/result, the documented 0.5.0 custom-aggregator
// extension point) and agg_engine.go's struct-of-arrays "reducible" lane
// (reducibleState + growReducible/updateReducible/mergeReducible/
// finalizeReducible, the fast lane the built-in engine actually uses). Both
// are meant to reproduce identical behavior; this test feeds the same fixed
// input to both lanes and asserts identical finalized output.
//
// Any/All have no accumulator.go equivalent — the SoA path is their only
// implementation (see reducibleState's doc comment in agg_engine.go) — so
// they are excluded here.
func TestAccumulatorSoAParity(t *testing.T) {
	// Fixed, hand-written input (no randomness): split into two sub-ranges so
	// the test exercises update (within each half), merge (combining the
	// halves) and result/finalize (reading the merged state) in both lanes.
	first := []float64{2, 4, 4, 4}
	second := []float64{5, 5, 7, 9}
	all := append(append([]float64{}, first...), second...)

	const eps = 1e-9

	type lane struct {
		name string
		acc  func() accumulator // accumulator.go lane constructor
		t    AggregateType      // agg_engine.go SoA lane aggregator kind
	}

	lanes := []lane{
		{"Count", newCountAcc, AGGREGATE_COUNT},
		{"Sum", newSumAcc, AGGREGATE_SUM},
		{"Mean", newMeanAcc, AGGREGATE_MEAN},
		{"Min", newMinAcc, AGGREGATE_MIN},
		{"Max", newMaxAcc, AGGREGATE_MAX},
		{"Std", func() accumulator { return newVarAcc(0, true) }, AGGREGATE_STD},
		{"Variance", func() accumulator { return newVarAcc(0, false) }, AGGREGATE_VARIANCE},
	}

	for _, ln := range lanes {
		t.Run(ln.name, func(t *testing.T) {
			// accumulator.go lane: accumulate each half separately, merge b
			// into a, then read the result.
			accA := ln.acc()
			for _, v := range first {
				accA.update(v)
			}
			accB := ln.acc()
			for _, v := range second {
				accB.update(v)
			}
			accA.merge(accB)
			wantV, wantNull := accA.result()

			// agg_engine.go SoA lane: same shape, one group (gid 0) per
			// reducibleState, merged the same way the parallel engine merges
			// worker-chunk state (see aggregate's merge loop).
			stA := &reducibleState{}
			growReducible(stA, ln.t)
			for _, v := range first {
				updateReducible(stA, ln.t, 0, v)
			}
			stB := &reducibleState{}
			growReducible(stB, ln.t)
			for _, v := range second {
				updateReducible(stB, ln.t, 0, v)
			}
			mergeReducible(stA, 0, stB, 0, ln.t)
			gotV, gotNull := finalizeReducible(stA, 0, ln.t, 0)

			if wantNull != gotNull {
				t.Fatalf("%s: isNull mismatch: accumulator.go=%v agg_engine.go=%v", ln.name, wantNull, gotNull)
			}
			if !wantNull && math.Abs(wantV-gotV) > eps {
				t.Fatalf("%s: value mismatch: accumulator.go=%v agg_engine.go=%v", ln.name, wantV, gotV)
			}

			// Cross-check: both lanes' split-then-merge result should also
			// equal a straight-through accumulation of the full input.
			accFull := ln.acc()
			for _, v := range all {
				accFull.update(v)
			}
			fullV, fullNull := accFull.result()
			if fullNull != wantNull || (!fullNull && math.Abs(fullV-wantV) > eps) {
				t.Fatalf("%s: split-then-merge (%v) diverges from straight-through accumulation (%v)", ln.name, wantV, fullV)
			}
		})
	}
}
