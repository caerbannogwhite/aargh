package dataframe

import (
	"fmt"

	"github.com/caerbannogwhite/enchanter/series"
)

type aggregatorBuilder struct {
	df          BaseDataFrame
	removeNAs   bool
	aggregators []aggregator
}

func (ab aggregatorBuilder) RemoveNAs(b bool) aggregatorBuilder {
	ab.removeNAs = b
	return ab
}

func (ab aggregatorBuilder) Run() DataFrame {
	df := ab.df
	if df.err != nil {
		return ab.df
	}

	if len(ab.aggregators) == 0 {
		return df
	}

	// CHECK: aggregators must have unique output names and input series must exist
	aggNewNames := make(map[string]bool)
	for _, agg := range ab.aggregators {

		// Check if output names are unique
		if aggNewNames[agg.newName] {
			df.err = fmt.Errorf("BaseDataFrame.Agg: aggregator output names must be unique")
			return df
		}
		aggNewNames[agg.newName] = true

		// CASE: aggregator count doesn't need an input series
		if agg.type_ != AGGREGATE_COUNT {
			if df.__series(agg.name) == nil {
				df.err = fmt.Errorf("BaseDataFrame.Agg: series \"%s\" not found", agg.name)
				return df
			}
		}
	}

	var result DataFrame
	if df.isGrouped {
		var indeces [][]int
		var flatIndeces []int
		result, indeces, flatIndeces, _ = df.groupHelper()

		groupsNum := len(indeces)

		var _series series.Series
		for _, agg := range ab.aggregators {
			_series = df.__series(agg.name)

			switch agg.type_ {
			case AGGREGATE_COUNT:
				counts := make([]int64, groupsNum)
				for i, group := range indeces {
					counts[i] = int64(len(group))
				}
				result = result.AddSeries(agg.newName, series.NewSeriesInt64(counts, nil, false, df.ctx))

			case AGGREGATE_SUM:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_sum(dataF64, flatIndeces, groupsNum, ab.removeNAs), nil, false, df.ctx))

			case AGGREGATE_MIN:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_min(dataF64, flatIndeces, groupsNum, ab.removeNAs), nil, false, df.ctx))

			case AGGREGATE_MAX:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_max(dataF64, flatIndeces, groupsNum, ab.removeNAs), nil, false, df.ctx))

			case AGGREGATE_MEAN:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_mean(dataF64, flatIndeces, groupsNum, ab.removeNAs), nil, false, df.ctx))

			case AGGREGATE_STD:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_std(dataF64, flatIndeces, groupsNum, ab.removeNAs), nil, false, df.ctx))
			}
		}

		// var wg sync.WaitGroup
		// wg.Add(THREADS_NUMBER)

		// buffer := make(chan __stats_thread_data)
		// for i := 0; i < THREADS_NUMBER; i++ {
		// 	go __stats_worker(&wg, buffer)
		// }

		// for _, agg := range aggregators {
		// 	series := df.__series(agg.name)

		// 	resultData := make([]float64, len(*indeces))
		// 	result = result.AddSeries(agg.name, NewSeriesFloat64(resultData, nil, false, df.ctx))
		// 	for gi, group := range *indeces {
		// 		buffer <- __stats_thread_data{
		// 			op:      agg.type_,
		// 			gi:      gi,
		// 			indeces: group,
		// 			series:  series,
		// 			res:     resultData,
		// 		}
		// 	}
		// }

		// close(buffer)
		// wg.Wait()

	} else {
		result = NewBaseDataFrame(df.ctx)

		var _series series.Series
		for _, agg := range ab.aggregators {
			_series = df.__series(agg.name)

			switch agg.type_ {
			case AGGREGATE_COUNT:
				result = result.AddSeries(agg.newName, series.NewSeriesInt64([]int64{int64(df.NRows())}, nil, false, df.ctx))

			case AGGREGATE_SUM:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_sum(dataF64, nil, 1, ab.removeNAs), nil, false, df.ctx))

			case AGGREGATE_MIN:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_min(dataF64, nil, 1, ab.removeNAs), nil, false, df.ctx))

			case AGGREGATE_MAX:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_max(dataF64, nil, 1, ab.removeNAs), nil, false, df.ctx))

			case AGGREGATE_MEAN:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_mean(dataF64, nil, 1, ab.removeNAs), nil, false, df.ctx))

			case AGGREGATE_STD:
				dataF64 := __gdl_stats_preprocess(_series)
				result = result.AddSeries(agg.newName, series.NewSeriesFloat64(__gdl_std(dataF64, nil, 1, ab.removeNAs), nil, false, df.ctx))
			}
		}
	}

	return result
}

type AggregateType int8

const (
	AGGREGATE_COUNT AggregateType = iota
	AGGREGATE_SUM
	AGGREGATE_MEAN
	AGGREGATE_MEDIAN
	AGGREGATE_MIN
	AGGREGATE_MAX
	AGGREGATE_STD
)

const DEFAULT_COUNT_NAME = "n"

type aggregator struct {
	name    string
	newName string
	type_   AggregateType
}

func Count() aggregator {
	return aggregator{DEFAULT_COUNT_NAME, DEFAULT_COUNT_NAME, AGGREGATE_COUNT}
}

func Sum(name string) aggregator {
	return aggregator{name, fmt.Sprintf("sum(%s)", name), AGGREGATE_SUM}
}

func Mean(name string) aggregator {
	return aggregator{name, fmt.Sprintf("mean(%s)", name), AGGREGATE_MEAN}
}

func Median(name string) aggregator {
	return aggregator{name, fmt.Sprintf("median(%s)", name), AGGREGATE_MEDIAN}
}

func Min(name string) aggregator {
	return aggregator{name, fmt.Sprintf("min(%s)", name), AGGREGATE_MIN}
}

func Max(name string) aggregator {
	return aggregator{name, fmt.Sprintf("max(%s)", name), AGGREGATE_MAX}
}

func Std(name string) aggregator {
	return aggregator{name, fmt.Sprintf("std(%s)", name), AGGREGATE_STD}
}

////////////////////////			SORT

type SortParam struct {
	asc     bool
	name    string
	_series series.Series
}

func Asc(name string) SortParam {
	return SortParam{asc: true, name: name}
}

func Desc(name string) SortParam {
	return SortParam{asc: false, name: name}
}
