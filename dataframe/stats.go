package dataframe

import (
	"math"

	"github.com/caerbannogwhite/enchanter/series"
)

func __gdl_stats_preprocess(s series.Series) []float64 {
	dataF64 := make([]float64, s.Len())

	switch _series := s.(type) {
	case series.Bools:
		if s.IsNullable() {
			for i, v := range _series.GetData() {
				if _series.IsNull(i) {
					dataF64[i] = math.NaN()
				} else if v {
					dataF64[i] = 1.0
				}
			}
		} else {
			for i, v := range _series.GetData() {
				if v {
					dataF64[i] = 1.0
				}
			}
		}

	case series.Ints:
		if s.IsNullable() {
			for i, v := range _series.GetData() {
				if _series.IsNull(i) {
					dataF64[i] = math.NaN()
				} else {
					dataF64[i] = float64(v)
				}
			}
		} else {
			for i, v := range _series.GetData() {
				dataF64[i] = float64(v)
			}
		}

	case series.Int64s:
		if s.IsNullable() {
			for i, v := range _series.GetData() {
				if _series.IsNull(i) {
					dataF64[i] = math.NaN()
				} else {
					dataF64[i] = float64(v)
				}
			}
		} else {
			for i, v := range _series.GetData() {
				dataF64[i] = float64(v)
			}
		}

	case series.Float64s:
		if s.IsNullable() {
			for i, v := range _series.GetData() {
				if _series.IsNull(i) {
					dataF64[i] = math.NaN()
				} else {
					dataF64[i] = v
				}
			}
		} else {
			dataF64 = _series.GetData()
		}

	case series.Durations:
		if s.IsNullable() {
			for i, v := range _series.GetData() {
				if _series.IsNull(i) {
					dataF64[i] = math.NaN()
				} else {
					dataF64[i] = float64(v)
				}
			}
		} else {
			for i, v := range _series.GetData() {
				dataF64[i] = float64(v)
			}
		}

	default:
		return nil
	}

	return dataF64
}
