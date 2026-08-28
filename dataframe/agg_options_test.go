package dataframe

import "testing"

func TestAggregatorOptions(t *testing.T) {
	// Defaults
	s := Std("v1")
	if s.ddof != 0 || s.ddofSet || s.type_ != AGGREGATE_STD || s.newName != "std(v1)" {
		t.Fatalf("Std defaults wrong: %+v", s)
	}
	// WithDDoF
	v := Variance("v1", WithDDoF(1))
	if v.ddof != 1 || !v.ddofSet || v.type_ != AGGREGATE_VARIANCE {
		t.Fatalf("Variance WithDDoF wrong: %+v", v)
	}
	// Quantile p + interpolation
	q := Quantile("v1", 0.25, WithInterpolation(Higher))
	if q.p != 0.25 || q.interp != Higher || !q.interpSet || q.type_ != AGGREGATE_QUANTILE {
		t.Fatalf("Quantile options wrong: %+v", q)
	}
	if !q.type_.isHolistic() || AGGREGATE_MEDIAN.isHolistic() != true || AGGREGATE_SUM.isHolistic() {
		t.Fatalf("isHolistic classification wrong")
	}
	// Median default interpolation is Linear
	if m := Median("v1"); m.interp != Linear || m.type_ != AGGREGATE_MEDIAN {
		t.Fatalf("Median defaults wrong: %+v", m)
	}
}
