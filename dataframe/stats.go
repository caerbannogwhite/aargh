package dataframe

import (
	"math"

	"github.com/caerbannogwhite/enchanter/series"
)

// type __numeric_thread_data struct {
// 	op      AggregateType
// 	gi      int
// 	indices []int
// 	series  Series
// 	res     []float64
// }

// func __numeric_worker(wg *sync.WaitGroup, buffer <-chan __numeric_thread_data) {
// 	for td := range buffer {
// 		switch td.op {
// 		case AGGREGATE_SUM:
// 			switch series := td.series.(type) {
// 			case Bools:

// 			case Ints:
// 				sum_ := int(0)
// 				data := series.GetDataPtr()
// 				for _, i := range td.indices {
// 					sum_ += (*data)[i]
// 				}
// 				td.res[td.gi] = float64(sum_)

// 			case Int64s:
// 				sum_ := int64(0)
// 				data := series.GetDataPtr()
// 				for _, i := range td.indices {
// 					sum_ += (*data)[i]
// 				}
// 				td.res[td.gi] = float64(sum_)

// 			case Float64s:
// 				sum_ := float64(0)
// 				data := series.GetDataPtr()
// 				for _, i := range td.indices {
// 					sum_ += (*data)[i]
// 				}
// 				td.res[td.gi] = sum_
// 			}
// 		case AGGREGATE_MIN:

// 		case AGGREGATE_MAX:

// 		case AGGREGATE_MEAN:
// 			switch series := td.series.(type) {
// 			case Bools:

// 			case Ints:
// 				sum_ := int(0)
// 				data := series.GetDataPtr()
// 				for _, i := range td.indices {
// 					sum_ += (*data)[i]
// 				}
// 				td.res[td.gi] = float64(sum_) / float64(len(td.indices))

// 			case Int64s:
// 				sum_ := int64(0)
// 				data := series.GetDataPtr()
// 				for _, i := range td.indices {
// 					sum_ += (*data)[i]
// 				}
// 				td.res[td.gi] = float64(sum_) / float64(len(td.indices))

// 			case Float64s:
// 				sum_ := float64(0)
// 				data := series.GetDataPtr()
// 				for _, i := range td.indices {
// 					sum_ += (*data)[i]
// 				}
// 				td.res[td.gi] = sum_ / float64(len(td.indices))
// 			}
// 		}
// 	}
// 	wg.Done()
// }

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

func __gdl_sum(dataF64 []float64, flatGroupIndeces []int, groupsNum int, removeNAs bool) []float64 {
	if dataF64 == nil {
		return nil
	}

	// SINGLE THREAD
	// if len(dataF64) < MINIMUM_PARALLEL_SIZE_2 {
	if flatGroupIndeces == nil {
		sum_ := float64(0)
		if removeNAs {
			for _, v := range dataF64 {
				if !math.IsNaN(v) {
					sum_ += v
				}
			}
		} else {
			for _, v := range dataF64 {
				sum_ += v
			}
		}
		return []float64{sum_}
	} else {
		sum := make([]float64, groupsNum)
		if removeNAs {
			for idx, gi := range flatGroupIndeces {
				if !math.IsNaN(dataF64[idx]) {
					sum[gi] += dataF64[idx]
				}
			}
		} else {
			for idx, gi := range flatGroupIndeces {
				sum[gi] += dataF64[idx]
			}
		}
		return sum
	}
	// }
}

func __gdl_min(dataF64 []float64, flatGroupIndeces []int, groupsNum int, removeNAs bool) []float64 {
	if dataF64 == nil {
		return nil
	}

	// SINGLE THREAD
	// if len(dataF64) < MINIMUM_PARALLEL_SIZE_2 {
	if flatGroupIndeces == nil {
		min_ := dataF64[0]
		if removeNAs {
			for _, v := range dataF64 {
				if !math.IsNaN(v) {
					min_ = min(min_, v)
				}
			}
		} else {
			for _, v := range dataF64 {
				min_ = min(min_, v)
			}
		}
		return []float64{min_}
	} else {
		min_ := make([]float64, groupsNum)
		for i := range min_ {
			min_[i] = math.Inf(1)
		}

		if removeNAs {
			for idx, gi := range flatGroupIndeces {
				if !math.IsNaN(dataF64[idx]) {
					min_[gi] = min(min_[gi], dataF64[idx])
				}
			}
		} else {
			for idx, gi := range flatGroupIndeces {
				min_[gi] = min(min_[gi], dataF64[idx])
			}
		}
		return min_
	}
	// }
}

func __gdl_max(dataF64 []float64, flatGroupIndeces []int, groupsNum int, removeNAs bool) []float64 {
	if dataF64 == nil {
		return nil
	}

	// SINGLE THREAD
	// if len(dataF64) < MINIMUM_PARALLEL_SIZE_2 {
	if flatGroupIndeces == nil {
		max_ := dataF64[0]
		if removeNAs {
			for _, v := range dataF64 {
				if !math.IsNaN(v) {
					max_ = max(max_, v)
				}
			}
		} else {
			for _, v := range dataF64 {
				max_ = max(max_, v)
			}
		}
		return []float64{max_}
	} else {
		max_ := make([]float64, groupsNum)
		for i := range max_ {
			max_[i] = math.Inf(-1)
		}

		if removeNAs {
			for idx, gi := range flatGroupIndeces {
				if !math.IsNaN(dataF64[idx]) {
					max_[gi] = max(max_[gi], dataF64[idx])
				}
			}
		} else {
			for idx, gi := range flatGroupIndeces {
				max_[gi] = max(max_[gi], dataF64[idx])
			}
		}
		return max_
	}
	// }
}

func __gdl_mean(dataF64 []float64, flatGroupIndeces []int, groupsNum int, removeNAs bool) []float64 {
	if flatGroupIndeces == nil {
		mean_ := float64(0)
		if removeNAs {
			count := 0
			for _, v := range dataF64 {
				if !math.IsNaN(v) {
					mean_ += v
					count++
				}
			}
			if count == 0 {
				return []float64{math.NaN()}
			}
			return []float64{mean_ / float64(count)}
		} else {
			for _, v := range dataF64 {
				mean_ += v
			}
			return []float64{mean_ / float64(len(dataF64))}
		}
	} else {
		mean_ := make([]float64, groupsNum)
		counts_ := make([]int, groupsNum)
		if removeNAs {
			for idx, gi := range flatGroupIndeces {
				if !math.IsNaN(dataF64[idx]) {
					mean_[gi] += dataF64[idx]
					counts_[gi]++
				}
			}
		} else {
			for idx, gi := range flatGroupIndeces {
				mean_[gi] += dataF64[idx]
				counts_[gi]++
			}
		}
		for i, v := range mean_ {
			mean_[i] = v / float64(counts_[i])
		}
		return mean_
	}
}

func __gdl_std(dataF64 []float64, flatGroupIndeces []int, groupsNum int, removeNAs bool) []float64 {
	mean_ := __gdl_mean(dataF64, flatGroupIndeces, groupsNum, removeNAs)
	if flatGroupIndeces == nil {
		std_ := float64(0)
		count := 0
		if removeNAs {
			for _, v := range dataF64 {
				if !math.IsNaN(v) {
					std_ += (v - mean_[0]) * (v - mean_[0])
					count++
				}
			}
		} else {
			for _, v := range dataF64 {
				std_ += (v - mean_[0]) * (v - mean_[0])
			}
			count = len(dataF64)
		}
		if count == 0 {
			return []float64{math.NaN()}
		}
		return []float64{math.Sqrt(std_ / float64(count))}
	} else {
		std_ := make([]float64, groupsNum)
		if removeNAs {
			for idx, gi := range flatGroupIndeces {
				if !math.IsNaN(dataF64[idx]) {
					std_[gi] += (dataF64[idx] - mean_[gi]) * (dataF64[idx] - mean_[gi])
				}
			}
		} else {
			for idx, gi := range flatGroupIndeces {
				std_[gi] += (dataF64[idx] - mean_[gi]) * (dataF64[idx] - mean_[gi])
			}
		}
		for i, v := range std_ {
			std_[i] = math.Sqrt(v / float64(len(flatGroupIndeces)/groupsNum))
		}
		return std_
	}
}
