package dataframe

import (
	"math"
	"testing"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestUngroupedMeanRemoveNAs(t *testing.T) {
	// values 2, 4, NaN(null) -> mean of non-nulls = 3, not 6/3 = 2
	data := []float64{2, 4, math.NaN()}
	got := __gdl_mean(data, nil, 1, true)
	if len(got) != 1 || !approx(got[0], 3.0) {
		t.Fatalf("mean removeNAs: got %v, want [3]", got)
	}
}

func TestUngroupedStdRemoveNAs(t *testing.T) {
	// values 2, 4, NaN(null); mean 3; population std over 2 non-nulls =
	// sqrt(((2-3)^2 + (4-3)^2)/2) = 1
	data := []float64{2, 4, math.NaN()}
	got := __gdl_std(data, nil, 1, true)
	if len(got) != 1 || !approx(got[0], 1.0) {
		t.Fatalf("std removeNAs: got %v, want [1]", got)
	}
}
