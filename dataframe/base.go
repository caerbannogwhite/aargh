package dataframe

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/formatter"
	"github.com/caerbannogwhite/aargh/meta"
	"github.com/caerbannogwhite/aargh/series"
	"github.com/caerbannogwhite/aargh/utils"
)

type BaseDataFramePartitionEntry struct {
	index     int
	name      string
	partition series.SeriesPartition
}

type BaseDataFrame struct {
	isGrouped  bool
	err        error
	names      []string
	series     []series.Series
	partitions []BaseDataFramePartitionEntry
	sortParams []SortParam
	ctx        *aargh.Context
}

func NewBaseDataFrame(ctx *aargh.Context) DataFrame {
	if ctx == nil {
		return BaseDataFrame{err: fmt.Errorf("NewBaseDataFrame: context is nil")}
	}

	return BaseDataFrame{
		series: make([]series.Series, 0),
		ctx:    ctx,
	}
}

////////////////////////			BASIC ACCESSORS

// GetContext returns the context of the dataframe.
func (df BaseDataFrame) GetContext() *aargh.Context {
	return df.ctx
}

// Names returns the names of the series in the dataframe.
func (df BaseDataFrame) Names() []string {
	return df.names
}

// Types returns the types of the series in the dataframe.
func (df BaseDataFrame) Types() []meta.BaseType {
	types := make([]meta.BaseType, len(df.series))
	for i, series := range df.series {
		types[i] = series.Type()
	}
	return types
}

// NCols returns the number of columns in the dataframe.
func (df BaseDataFrame) NCols() int {
	return len(df.series)
}

// NRows returns the number of rows in the dataframe.
func (df BaseDataFrame) NRows() int {
	if len(df.series) == 0 {
		return 0
	}
	return df.series[0].Len()
}

func (df BaseDataFrame) IsErrored() bool {
	return df.err != nil
}

func (df BaseDataFrame) IsGrouped() bool {
	return df.isGrouped
}

func (df BaseDataFrame) GetError() error {
	return df.err
}

func (df BaseDataFrame) GetSeriesIndex(name string) int {
	for i, name_ := range df.names {
		if name_ == name {
			return i
		}
	}
	return -1
}

func (df BaseDataFrame) AddSeries(name string, series series.Series) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.AddSeries: cannot add series to a grouped dataframe")
		return df
	}

	if df.NCols() > 0 && series.Len() != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.AddSeries: series length (%d) does not match dataframe length (%d)", series.Len(), df.NRows())
		return df
	}

	df.names = append(df.names, name)
	df.series = append(df.series, series)

	return df
}

func (df BaseDataFrame) AddSeriesFromBools(name string, data []bool, nullMask []bool, makeCopy bool) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromBools: cannot add series to a grouped dataframe")
		return df
	}

	if df.NCols() > 0 && len(data) != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromBools: series length (%d) does not match dataframe length (%d)", len(data), df.NRows())
		return df
	}

	return df.AddSeries(name, series.NewSeriesBool(data, nullMask, makeCopy, df.ctx))
}

func (df BaseDataFrame) AddSeriesFromInts(name string, data []int, nullMask []bool, makeCopy bool) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromInts: cannot add series to a grouped dataframe")
		return df
	}

	if df.NCols() > 0 && len(data) != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromInts: series length (%d) does not match dataframe length (%d)", len(data), df.NRows())
		return df
	}

	return df.AddSeries(name, series.NewSeriesInt(data, nullMask, makeCopy, df.ctx))
}

func (df BaseDataFrame) AddSeriesFromInt64s(name string, data []int64, nullMask []bool, makeCopy bool) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromInt64s: cannot add series to a grouped dataframe")
		return df
	}

	if df.NCols() > 0 && len(data) != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromInt64s: series length (%d) does not match dataframe length (%d)", len(data), df.NRows())
		return df
	}

	return df.AddSeries(name, series.NewSeriesInt64(data, nullMask, makeCopy, df.ctx))
}

func (df BaseDataFrame) AddSeriesFromFloat64s(name string, data []float64, nullMask []bool, makeCopy bool) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromFloat64s: cannot add series to a grouped dataframe")
		return df
	}

	if df.NCols() > 0 && len(data) != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromFloat64s: series length (%d) does not match dataframe length (%d)", len(data), df.NRows())
		return df
	}

	return df.AddSeries(name, series.NewSeriesFloat64(data, nullMask, makeCopy, df.ctx))
}

func (df BaseDataFrame) AddSeriesFromStrings(name string, data []string, nullMask []bool, makeCopy bool) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromStrings: cannot add series to a grouped dataframe")
		return df
	}

	if df.NCols() > 0 && len(data) != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromStrings: series length (%d) does not match dataframe length (%d)", len(data), df.NRows())
		return df
	}

	return df.AddSeries(name, series.NewSeriesString(data, nullMask, makeCopy, df.ctx))
}

func (df BaseDataFrame) AddSeriesFromTimes(name string, data []time.Time, nullMask []bool, makeCopy bool) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromTimes: cannot add series to a grouped dataframe")
		return df
	}

	if df.NCols() > 0 && len(data) != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromTimes: series length (%d) does not match dataframe length (%d)", len(data), df.NRows())
		return df
	}

	return df.AddSeries(name, series.NewSeriesTime(data, nullMask, makeCopy, df.ctx))
}

func (df BaseDataFrame) AddSeriesFromDurations(name string, data []time.Duration, nullMask []bool, makeCopy bool) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromDurations: cannot add series to a grouped dataframe")
		return df
	}

	if df.NCols() > 0 && len(data) != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.AddSeriesFromDurations: series length (%d) does not match dataframe length (%d)", len(data), df.NRows())
		return df
	}

	return df.AddSeries(name, series.NewSeriesDuration(data, nullMask, makeCopy, df.ctx))
}

func (df BaseDataFrame) Replace(name string, s series.Series) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.Replace: cannot replace series in a grouped dataframe")
		return df
	}

	index := df.GetSeriesIndex(name)
	if index == -1 {
		df.err = fmt.Errorf("BaseDataFrame.Replace: series \"%s\" not found", name)
		return df
	}

	if s.Len() != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.Replace: series length (%d) does not match dataframe length (%d)", s.Len(), df.NRows())
		return df
	}

	df.series[index] = s
	return df
}

// Returns the column with the given name.
func (df BaseDataFrame) C(name string) series.Series {
	for i, name_ := range df.names {
		if name_ == name {
			return df.series[i]
		}
	}

	return series.Errors{Msg_: fmt.Sprintf("BaseDataFrame.C: series \"%s\" not found", name)}
}

// Returns the series with the given name.
// For internal use only: returns nil if the series is not found.
func (df BaseDataFrame) __series(name string) series.Series {
	for i, name_ := range df.names {
		if name_ == name {
			return df.series[i]
		}
	}

	return nil
}

// Returns the series at the given index.
func (df BaseDataFrame) At(index int) series.Series {
	if index < 0 || index >= len(df.series) {
		return series.Errors{Msg_: fmt.Sprintf("BaseDataFrame.SeriesAt: index %d out of bounds", index)}
	}
	return df.series[index]
}

// Returns the series with the given name as a bool series.
func (df BaseDataFrame) NameAt(index int) string {
	if index < 0 || index >= len(df.names) {
		return ""
	}
	return df.names[index]
}

func (df BaseDataFrame) Select(selectors ...string) DataFrame {
	if df.err != nil {
		return df
	}

	regexes := make([]*regexp.Regexp, len(selectors))
	for i, selector := range selectors {
		regex, err := regexp.Compile(selector)
		if err != nil {
			df.err = fmt.Errorf("BaseDataFrame.Select: invalid selector \"%s\"", selector)
			return df
		}
		regexes[i] = regex
	}

	selected := make(map[string]bool)
	for _, name := range df.names {
		selected[name] = false
	}

	names := make([]string, 0)
	seriesList := make([]series.Series, 0)

	for _, regex := range regexes {
		for _, name := range df.names {
			if !selected[name] && regex.MatchString(name) {
				selected[name] = true
				names = append(names, name)
				seriesList = append(seriesList, df.C(name))
			}
		}
	}

	return BaseDataFrame{
		names:  names,
		series: seriesList,
		ctx:    df.ctx,
	}
}

func (df BaseDataFrame) SelectAt(indices ...int) DataFrame {
	if df.err != nil {
		return df
	}

	selected := NewBaseDataFrame(df.ctx)
	for _, index := range indices {
		if index < 0 || index >= len(df.series) {
			selected.AddSeries(df.names[index], df.series[index])
		} else {
			return BaseDataFrame{err: fmt.Errorf("BaseDataFrame.SelectAt: index %d out of bounds", index)}
		}
	}

	return selected
}

func (df BaseDataFrame) Filter(mask any) DataFrame {
	if df.err != nil {
		return df
	}

	var maskSeries series.Bools
	if _, ok := mask.(series.Bools); ok {
		maskSeries = mask.(series.Bools)

	} else {
		if mask, ok := mask.([]bool); ok {
			maskSeries = series.NewSeriesBool(mask, nil, false, df.ctx)
		} else {
			df.err = fmt.Errorf("BaseDataFrame.Filter: mask is not a bool series")
			return df
		}
	}

	if maskSeries.Len() != df.NRows() {
		df.err = fmt.Errorf("BaseDataFrame.Filter: mask length (%d) does not match dataframe length (%d)", maskSeries.Len(), df.NRows())
		return df
	}

	seriesList := make([]series.Series, 0)
	for _, series := range df.series {
		seriesList = append(seriesList, series.Filter(maskSeries))
	}

	return BaseDataFrame{
		names:  df.names,
		series: seriesList,
		ctx:    df.ctx,
	}
}

func (df BaseDataFrame) GroupBy(by ...string) DataFrame {
	if df.err != nil {
		return df
	}

	// Grouping an already grouped dataframe replaces the existing grouping.
	if df.isGrouped {
		df = df.Ungroup().(BaseDataFrame)
	}

	{

		// Check that all the group by columns exist
		for _, name := range by {
			found := false
			for _, name_ := range df.names {
				if name_ == name {
					found = true
					break
				}
			}

			if !found {
				df.err = fmt.Errorf("BaseDataFrame.GroupBy: column \"%s\" not found", name)
				return df
			}
		}

		df.isGrouped = true
		df.partitions = make([]BaseDataFramePartitionEntry, len(by))

		for partitionsIndex, name := range by {
			i := df.GetSeriesIndex(name)
			series := df.series[i]

			// First partition: group the series
			if partitionsIndex == 0 {
				df.partitions[partitionsIndex] = BaseDataFramePartitionEntry{
					index:     i,
					name:      name,
					partition: series.Group().GetPartition(),
				}
			} else

			// Subsequent partitions: sub-group the series
			{
				df.partitions[partitionsIndex] = BaseDataFramePartitionEntry{
					index:     i,
					name:      name,
					partition: series.GroupBy(df.partitions[partitionsIndex-1].partition).GetPartition(),
				}
			}
		}

		return df
	}
}

func (df BaseDataFrame) Ungroup() DataFrame {
	if df.err != nil {
		return df
	}

	df.isGrouped = false
	df.partitions = nil
	return df
}

func (df BaseDataFrame) getPartitions() []series.SeriesPartition {
	if df.err != nil {
		return nil
	}

	if df.isGrouped {
		partitions := make([]series.SeriesPartition, len(df.partitions))
		for i, partition := range df.partitions {
			partitions[i] = partition.partition
		}
		return partitions
	} else {
		return nil
	}
}

func (df BaseDataFrame) groupHelper() (DataFrame, [][]int, []int, []int) {

	// Keep track of which series are not grouped
	seriesIndices := make(map[int]bool)
	for i := 0; i < df.NCols(); i++ {
		seriesIndices[i] = true
	}

	result := NewBaseDataFrame(df.ctx).(BaseDataFrame)

	// The last partition tells us how many groups there are
	// and how many rows are in each group
	indeces := make([][]int, 0, df.partitions[len(df.partitions)-1].partition.GetSize())
	for _, group := range df.partitions[len(df.partitions)-1].partition.GetMap() {
		indeces = append(indeces, group)
	}

	// Keep only the grouped series
	for _, partition := range df.partitions {
		seriesIndices[partition.index] = false
		old := df.series[partition.index]

		result.names = append(result.names, partition.name)

		// Each group is represented by its first row; nulls group together,
		// so the key column keeps the null flag of that representative row.
		var keyNulls []bool
		if old.IsNullable() {
			keyNulls = make([]bool, len(indeces))
			for i, group := range indeces {
				keyNulls[i] = old.IsNull(group[0])
			}
		}

		switch _series := old.(type) {
		case series.Bools:
			values := make([]bool, len(indeces))
			for i, group := range indeces {
				values[i] = _series.Data_[group[0]]
			}
			result.series = append(result.series, series.NewSeriesBool(values, keyNulls, false, _series.GetContext()))

		case series.Ints:
			values := make([]int, len(indeces))
			for i, group := range indeces {
				values[i] = _series.Data_[group[0]]
			}
			result.series = append(result.series, series.NewSeriesInt(values, keyNulls, false, _series.GetContext()))

		case series.Int64s:
			values := make([]int64, len(indeces))
			for i, group := range indeces {
				values[i] = _series.Data_[group[0]]
			}
			result.series = append(result.series, series.NewSeriesInt64(values, keyNulls, false, _series.GetContext()))

		case series.Float64s:
			values := make([]float64, len(indeces))
			for i, group := range indeces {
				values[i] = _series.Data_[group[0]]
			}
			result.series = append(result.series, series.NewSeriesFloat64(values, keyNulls, false, _series.GetContext()))

		case series.Strings:
			values := make([]*string, len(indeces))
			for i, group := range indeces {
				values[i] = _series.Data_[group[0]]
			}
			result.series = append(result.series, series.NewSeriesStringFromPtrs(values, keyNulls, false, _series.GetContext()))

		case series.Times:
			values := make([]time.Time, len(indeces))
			for i, group := range indeces {
				values[i] = _series.Data_[group[0]]
			}
			result.series = append(result.series, series.NewSeriesTime(values, keyNulls, false, _series.GetContext()))

		case series.Durations:
			values := make([]time.Duration, len(indeces))
			for i, group := range indeces {
				values[i] = _series.Data_[group[0]]
			}
			result.series = append(result.series, series.NewSeriesDuration(values, keyNulls, false, _series.GetContext()))

		default:
			// Unsupported key type: keep names and series aligned.
			result.series = append(result.series, series.NewSeriesNA(len(indeces), df.ctx))
		}
	}

	// Get the indices of the ungrouped series
	ungroupedSeriesIndices := make([]int, 0)
	for index, isGrouped := range seriesIndices {
		if isGrouped {
			ungroupedSeriesIndices = append(ungroupedSeriesIndices, index)
		}
	}

	// sort the indices
	sort.Ints(ungroupedSeriesIndices)

	// Flatten the indeces
	flatIndeces := make([]int, df.NRows())
	for i, group := range indeces {
		for _, index := range group {
			flatIndeces[index] = i
		}
	}

	return result, indeces, flatIndeces, ungroupedSeriesIndices
}

func (df BaseDataFrame) Join(how DataFrameJoinType, other DataFrame, on ...string) DataFrame {
	if df.err != nil {
		return df
	}

	// CASE: the dataframes have different contexts
	if df.ctx != other.GetContext() {
		df.err = fmt.Errorf("BaseDataFrame.Join: dataframes have different contexts")
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.Join: cannot join a grouped dataframe")
		return df
	}

	if other.IsGrouped() {
		df.err = fmt.Errorf("BaseDataFrame.Join: cannot join with a grouped dataframe")
		return df
	}

	// CHECK: all the join columns must exist
	// CHECK: all the join columns must have the same type
	types := make([]meta.BaseType, len(on))
	for _, name := range on {

		// Series A
		found := false
		for idx, series := range df.series {
			if df.names[idx] == name {
				found = true

				// keep track of the types
				types = append(types, series.Type())
				break
			}
		}
		if !found {
			df.err = fmt.Errorf("BaseDataFrame.Join: column \"%s\" not found in left dataframe", name)
			return df
		}

		// Series B
		found = false
		otherNames := other.Names()
		for idx, series := range other.(BaseDataFrame).series {
			if idx < len(otherNames) && otherNames[idx] == name {
				found = true

				// CHECK: the types must match
				if types[len(types)-1] != series.Type() {
					df.err = fmt.Errorf("BaseDataFrame.Join: columns \"%s\" have different types", name)
					return df
				}
				break
			}
		}
		if !found {
			df.err = fmt.Errorf("BaseDataFrame.Join: column \"%s\" not found in right dataframe", name)
			return df
		}
	}

	// CASE: on is empty -> use all columns with the same name
	if len(on) == 0 {
		for _, name := range df.Names() {
			if other.GetSeriesIndex(name) != -1 {
				on = append(on, name)
			}
		}
	}

	// CASE: on is still empty -> error
	if len(on) == 0 {
		df.err = fmt.Errorf("BaseDataFrame.Join: no columns to join on")
		return df
	}

	// CHECK: all columns in on must have the same type
	for _, name := range on {
		if df.C(name).Type() != other.C(name).Type() {
			df.err = fmt.Errorf("BaseDataFrame.Join: columns \"%s\" have different types", name)
			return df
		}
	}

	// Group the dataframes by the join columns
	dfGrouped := df.GroupBy(on...).(BaseDataFrame)
	otherGrouped := other.GroupBy(on...).(BaseDataFrame)

	colsDiffA := make([]string, 0)
	colsDiffB := make([]string, 0)

	// Get the columns that are not in the join columns
	for _, name := range df.Names() {
		found := false
		for _, joinName := range on {
			if name == joinName {
				found = true
				break
			}
		}
		if !found {
			colsDiffA = append(colsDiffA, name)
		}
	}

	for _, name := range other.Names() {
		found := false
		for _, joinName := range on {
			if name == joinName {
				found = true
				break
			}
		}
		if !found {
			colsDiffB = append(colsDiffB, name)
		}
	}

	// Get the columns that are in both dataframes
	commonCols := make(map[string]bool)
	for _, name := range df.Names() {
		for _, otherName := range other.Names() {
			if name == otherName {
				commonCols[name] = true
				break
			}
		}
	}

	joined := NewBaseDataFrame(df.ctx)

	pA := dfGrouped.getPartitions()
	pB := otherGrouped.getPartitions()

	// Get the maps, keys and sort them
	mapA := pA[len(pA)-1].GetMap()
	mapB := pB[len(pB)-1].GetMap()

	keysA := make([]int64, 0, len(mapA))
	keysB := make([]int64, 0, len(mapB))

	for key := range mapA {
		keysA = append(keysA, key)
	}

	for key := range mapB {
		keysB = append(keysB, key)
	}

	sort.Slice(keysA, func(i, j int) bool { return keysA[i] < keysA[j] })
	sort.Slice(keysB, func(i, j int) bool { return keysB[i] < keysB[j] })

	// Find the intersection
	keysAOnly := make([]int64, 0, len(keysA))
	keysBOnly := make([]int64, 0, len(keysB))
	keysIntersection := make([]int64, 0, len(keysA))

	var i, j int = 0, 0
	for i < len(keysA) && j < len(keysB) {
		if keysA[i] < keysB[j] {
			keysAOnly = append(keysAOnly, keysA[i])
			i++
		} else if keysA[i] > keysB[j] {
			keysBOnly = append(keysBOnly, keysB[j])
			j++
		} else {
			keysIntersection = append(keysIntersection, keysA[i])
			i++
			j++
		}
	}

	for i < len(keysA) {
		keysAOnly = append(keysAOnly, keysA[i])
		i++
	}

	for j < len(keysB) {
		keysBOnly = append(keysBOnly, keysB[j])
		j++
	}

	switch how {
	case INNER_JOIN:
		// Get indices of the intersection
		indicesA := make([]int, 0, len(keysIntersection))
		indicesB := make([]int, 0, len(keysIntersection))

		for _, key := range keysIntersection {
			for _, indexA := range mapA[key] {
				for _, indexB := range mapB[key] {
					indicesA = append(indicesA, indexA)
					indicesB = append(indicesB, indexB)
				}
			}
		}

		// Join columns
		for i, name := range on {
			joined = joined.AddSeries(name, dfGrouped.C(on[i]).Filter(indicesA))
		}

		// A columns
		var ser_ series.Series
		for _, name := range colsDiffA {
			ser_ = df.C(name).Filter(indicesA)
			if commonCols[name] {
				name += "_x"
			}
			joined = joined.AddSeries(name, ser_)
		}

		// B columns
		for _, name := range colsDiffB {
			ser_ = other.C(name).Filter(indicesB)
			if commonCols[name] {
				name += "_y"
			}
			joined = joined.AddSeries(name, ser_.Filter(indicesB))
		}

	case LEFT_JOIN:
		indicesA := make([]int, 0, len(keysA))
		indicesB := make([]int, 0, len(keysIntersection))

		for _, key := range keysAOnly {
			indicesA = append(indicesA, mapA[key]...)
		}

		for _, key := range keysIntersection {
			for _, indexA := range mapA[key] {
				for _, indexB := range mapB[key] {
					indicesA = append(indicesA, indexA)
					indicesB = append(indicesB, indexB)
				}
			}
		}

		// Join columns
		for i, name := range on {
			joined = joined.AddSeries(name, dfGrouped.C(on[i]).Filter(indicesA))
		}

		// A columns
		var ser_ series.Series
		for _, name := range colsDiffA {
			ser_ = df.C(name).Filter(indicesA)
			if commonCols[name] {
				name += "_x"
			}
			joined = joined.AddSeries(name, ser_)
		}

		padBlen := len(indicesA) - len(indicesB)
		nullMask := make([]bool, padBlen)
		for i := range nullMask {
			nullMask[i] = true
		}

		// B columns
		for _, name := range colsDiffB {
			ser_ = other.C(name).Filter(indicesB)
			switch ser_.Type() {
			case meta.BoolType:
				ser_ = series.NewSeriesBool(make([]bool, padBlen), nullMask, false, df.ctx).
					Append(ser_)

			case meta.IntType:
				ser_ = series.NewSeriesInt(make([]int, padBlen), nullMask, false, df.ctx).
					Append(ser_)

			case meta.Int64Type:
				ser_ = series.NewSeriesInt64(make([]int64, padBlen), nullMask, false, df.ctx).
					Append(ser_)

			case meta.Float64Type:
				ser_ = series.NewSeriesFloat64(make([]float64, padBlen), nullMask, false, df.ctx).
					Append(ser_)

			case meta.StringType:
				ser_ = series.NewSeriesString(make([]string, padBlen), nullMask, false, df.ctx).
					Append(ser_)

			case meta.TimeType:
				ser_ = series.NewSeriesTime(make([]time.Time, padBlen), nullMask, false, df.ctx).
					Append(ser_)

			case meta.DurationType:
				ser_ = series.NewSeriesDuration(make([]time.Duration, padBlen), nullMask, false, df.ctx).
					Append(ser_)
			}

			if commonCols[name] {
				name += "_y"
			}
			joined = joined.AddSeries(name, ser_)
		}

	case RIGHT_JOIN:
		indicesA := make([]int, 0, len(keysIntersection))
		indicesB := make([]int, 0, len(keysB))

		for _, key := range keysIntersection {
			for _, indexA := range mapA[key] {
				for _, indexB := range mapB[key] {
					indicesA = append(indicesA, indexA)
					indicesB = append(indicesB, indexB)
				}
			}
		}

		for _, key := range keysBOnly {
			indicesB = append(indicesB, mapB[key]...)
		}

		// Join columns
		for i, name := range on {
			joined = joined.AddSeries(name, otherGrouped.C(on[i]).Filter(indicesB))
		}

		padAlen := len(indicesB) - len(indicesA)
		nullMask := make([]bool, padAlen)
		for i := range nullMask {
			nullMask[i] = true
		}

		// A columns
		var ser_ series.Series
		for _, name := range colsDiffA {
			ser_ = df.C(name).Filter(indicesA)
			switch ser_.Type() {
			case meta.BoolType:
				ser_ = ser_.(series.Bools).Append(series.NewSeriesBool(make([]bool, padAlen), nullMask, false, df.ctx))

			case meta.IntType:
				ser_ = ser_.(series.Ints).Append(series.NewSeriesInt(make([]int, padAlen), nullMask, false, df.ctx))

			case meta.Int64Type:
				ser_ = ser_.(series.Int64s).Append(series.NewSeriesInt64(make([]int64, padAlen), nullMask, false, df.ctx))

			case meta.Float64Type:
				ser_ = ser_.(series.Float64s).Append(series.NewSeriesFloat64(make([]float64, padAlen), nullMask, false, df.ctx))

			case meta.StringType:
				ser_ = ser_.(series.Strings).Append(series.NewSeriesString(make([]string, padAlen), nullMask, false, df.ctx))

			case meta.TimeType:
				ser_ = ser_.(series.Times).Append(series.NewSeriesTime(make([]time.Time, padAlen), nullMask, false, df.ctx))

			case meta.DurationType:
				ser_ = ser_.(series.Durations).Append(series.NewSeriesDuration(make([]time.Duration, padAlen), nullMask, false, df.ctx))
			}

			if commonCols[name] {
				name += "_x"
			}
			joined = joined.AddSeries(name, ser_)
		}

		// B columns
		for _, name := range colsDiffB {
			ser_ = other.C(name).Filter(indicesB)
			if commonCols[name] {
				name += "_y"
			}
			joined = joined.AddSeries(name, ser_)
		}

	case OUTER_JOIN:
		indicesA := make([]int, 0, len(keysA))
		indicesB := make([]int, 0, len(keysB))

		padAlen := 0
		padBlen := 0

		for _, key := range keysAOnly {
			indicesA = append(indicesA, mapA[key]...)
			padBlen += len(mapA[key])
		}

		intersectionLen := 0
		for _, key := range keysIntersection {
			for _, indexA := range mapA[key] {
				for _, indexB := range mapB[key] {
					indicesA = append(indicesA, indexA)
					indicesB = append(indicesB, indexB)
					intersectionLen++
				}
			}
		}

		for _, key := range keysBOnly {
			indicesB = append(indicesB, mapB[key]...)
			padAlen += len(mapB[key])
		}

		// Join columns
		indicesBOnly := indicesB[intersectionLen:]
		for i, name := range on {
			joined = joined.AddSeries(name,
				dfGrouped.C(on[i]).
					Filter(indicesA).Append(
					otherGrouped.C(on[i]).
						Filter(indicesBOnly)))
		}

		nullMaskA := make([]bool, padAlen)
		for i := range nullMaskA {
			nullMaskA[i] = true
		}

		nullMaskB := make([]bool, padBlen)
		for i := range nullMaskB {
			nullMaskB[i] = true
		}

		// A columns
		var ser_ series.Series
		for _, name := range colsDiffA {
			ser_ = df.C(name).Filter(indicesA)
			switch ser_.Type() {
			case meta.BoolType:
				ser_ = ser_.(series.Bools).Append(series.NewSeriesBool(make([]bool, padAlen), nullMaskA, false, df.ctx))

			case meta.IntType:
				ser_ = ser_.(series.Ints).Append(series.NewSeriesInt(make([]int, padAlen), nullMaskA, false, df.ctx))

			case meta.Int64Type:
				ser_ = ser_.(series.Int64s).Append(series.NewSeriesInt64(make([]int64, padAlen), nullMaskA, false, df.ctx))

			case meta.Float64Type:
				ser_ = ser_.(series.Float64s).Append(series.NewSeriesFloat64(make([]float64, padAlen), nullMaskA, false, df.ctx))

			case meta.StringType:
				ser_ = ser_.(series.Strings).Append(series.NewSeriesString(make([]string, padAlen), nullMaskA, false, df.ctx))

			case meta.TimeType:
				ser_ = ser_.(series.Times).Append(series.NewSeriesTime(make([]time.Time, padAlen), nullMaskA, false, df.ctx))

			case meta.DurationType:
				ser_ = ser_.(series.Durations).Append(series.NewSeriesDuration(make([]time.Duration, padAlen), nullMaskA, false, df.ctx))
			}

			if commonCols[name] {
				name += "_x"
			}
			joined = joined.AddSeries(name, ser_)
		}

		// B columns
		for _, name := range colsDiffB {
			ser_ = other.C(name).Filter(indicesB)
			switch ser_.Type() {
			case meta.BoolType:
				ser_ = series.NewSeriesBool(make([]bool, padBlen), nullMaskB, false, df.ctx).
					Append(ser_)

			case meta.IntType:
				ser_ = series.NewSeriesInt(make([]int, padBlen), nullMaskB, false, df.ctx).
					Append(ser_)

			case meta.Int64Type:
				ser_ = series.NewSeriesInt64(make([]int64, padBlen), nullMaskB, false, df.ctx).
					Append(ser_)

			case meta.Float64Type:
				ser_ = series.NewSeriesFloat64(make([]float64, padBlen), nullMaskB, false, df.ctx).
					Append(ser_)

			case meta.StringType:
				ser_ = series.NewSeriesString(make([]string, padBlen), nullMaskB, false, df.ctx).
					Append(ser_)

			case meta.TimeType:
				ser_ = series.NewSeriesTime(make([]time.Time, padBlen), nullMaskB, false, df.ctx).
					Append(ser_)

			case meta.DurationType:
				ser_ = series.NewSeriesDuration(make([]time.Duration, padBlen), nullMaskB, false, df.ctx).
					Append(ser_)
			}

			if commonCols[name] {
				name += "_y"
			}
			joined = joined.AddSeries(name, ser_)
		}
	}

	return joined
}

func (df BaseDataFrame) Take(params ...int) DataFrame {
	if df.err != nil {
		return df
	}

	indeces, err := series.SeriesTakePreprocess("BaseDataFrame", df.NRows(), params...)
	if err != nil {
		df.err = err
		return df
	}

	taken := NewBaseDataFrame(df.ctx)
	for idx, series := range df.series {
		taken = taken.AddSeries(df.names[idx], series.FilterIntSlice(indeces, false))
	}

	return taken
}

func (df BaseDataFrame) Len() int {
	if df.err != nil || len(df.series) < 1 {
		return 0
	}

	return df.series[0].Len()
}

func (df BaseDataFrame) Less(i, j int) bool {
	for _, param := range df.sortParams {
		if !param._series.Equal(i, j) {
			return (param.asc && param._series.Less(i, j)) || (!param.asc && param._series.Less(j, i))
		}
	}

	return false
}

func (df BaseDataFrame) Swap(i, j int) {
	for _, series := range df.series {
		series.Swap(i, j)
	}
}

func (df BaseDataFrame) OrderBy(params ...SortParam) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.OrderBy: cannot order grouped DataFrame")
		return df
	}

	// CHECK: params must have unique names and names must be valid
	paramNames := make(map[string]bool)
	for i, param := range params {
		if paramNames[param.name] {
			df.err = fmt.Errorf("BaseDataFrame.OrderBy: series names must be unique")
			return df
		}
		paramNames[param.name] = true

		if series := df.__series(param.name); series != nil {
			params[i]._series = series
		} else {
			df.err = fmt.Errorf("BaseDataFrame.OrderBy: series \"%s\" not found", param.name)
			return df
		}
	}

	df.sortParams = params
	sort.Sort(df)
	df.sortParams = nil

	return df
}

////////////////////////			SUMMARY

func (df BaseDataFrame) Agg(aggregators ...aggregator) aggregatorBuilder {
	return aggregatorBuilder{df, false, aggregators}
}

////////////////////////			PRINTING

func (df BaseDataFrame) Describe() string {
	return ""
}

func (df BaseDataFrame) Records(header bool) [][]string {
	var out [][]string
	if header {
		out = make([][]string, df.NRows()+1)
	} else {
		out = make([][]string, df.NRows())
	}

	h := 0
	if header {
		out[0] = make([]string, df.NCols())
		for j := 0; j < df.NCols(); j++ {
			out[0][j] = df.names[j]
		}

		h = 1
	}

	for i := 0 + h; i < df.NRows()+h; i++ {
		out[i] = make([]string, df.NCols())
		for j := 0; j < df.NCols(); j++ {
			out[i][j] = df.series[j].GetAsString(i - h)
		}
	}

	return out
}

// Pretty print the dataframe.
func (df BaseDataFrame) PPrint(params PPrintParams) DataFrame {
	if df.err != nil {
		fmt.Println(df.err)
		return df
	}

	buffer := ""

	// check if the dataframe is empty
	if df.NRows() == 0 {
		buffer += params.indent
		if params.useLipGloss {
			params.styleNames.Render("  Empty DataFrame\n")
		} else {
			buffer += "  Empty DataFrame\n"
		}
		fmt.Println(buffer)
		return df
	}

	// print the shape
	buffer += params.indent
	if params.useLipGloss {
		buffer += params.styleTypes.Render(fmt.Sprintf("  BaseDataFrame: %d rows, %d columns", df.NRows(), df.NCols()))
	} else {
		buffer += fmt.Sprintf("  BaseDataFrame: %d rows, %d columns", df.NRows(), df.NCols())
	}
	buffer += "\n"

	// print the group by columns
	if df.isGrouped {
		buffer += params.indent
		if params.useLipGloss {
			buffer += params.styleTypes.Render("  Grouped by: ")
		} else {
			buffer += "  Grouped by: "
		}
		for i, partition := range df.partitions {
			if params.useLipGloss {
				buffer += params.styleTypes.Render(fmt.Sprintf("%s", partition.name))
			} else {
				buffer += fmt.Sprintf("%s", partition.name)
			}
			if i < len(df.partitions)-1 {
				if params.useLipGloss {
					buffer += params.styleTypes.Render(",")
				} else {
					buffer += ","
				}
			}
		}
		buffer += "\n"
	}

	// check how many variables can fit in the screen
	nColsOut := 0
	actualWidthsSum := 0

	widths := make([]int, df.NCols())
	for i, name := range df.names {
		widths[i] = max(len(df.series[i].Type().String()), len(name))
		actualWidthsSum += widths[i] + 3
		if actualWidthsSum > params.width {
			break
		}
		nColsOut++
	}
	widths = widths[:nColsOut]

	nRowsOut := min(10, df.NRows())
	if params.nrows > 0 {
		nRowsOut = min(params.nrows, df.NRows())
	}

	addTail := false
	if df.NRows() > nRowsOut+params.tailLen {
		addTail = true
		nRowsOut -= params.tailLen
	}

	formatters := make([]formatter.Formatter, nColsOut)
	for i := 0; i < nColsOut; i++ {
		switch df.series[i].Type() {
		case meta.BoolType, meta.StringType, meta.TimeType:
			formatters[i] = formatter.NewStringFormatter().
				SetUseLipGloss(params.useLipGloss)
		case meta.IntType, meta.Int64Type, meta.Float64Type, meta.DurationType:
			formatters[i] = formatter.NewNumericFormatter().
				SetUseLipGloss(params.useLipGloss).
				SetNaText(df.ctx.GetNaText()).
				SetTruncateOutput(true)
		}

		switch s := df.series[i].(type) {
		case series.Bools:
			for _, v := range s.DataAsString()[:nRowsOut] {
				formatters[i].Push(v)
			}

			if addTail {
				for _, v := range s.DataAsString()[df.NRows()-params.tailLen:] {
					formatters[i].Push(v)
				}
			}

		case series.Ints:
			for _, v := range s.Ints()[:nRowsOut] {
				formatters[i].Push(v)
			}

			if addTail {
				for _, v := range s.Ints()[df.NRows()-params.tailLen:] {
					formatters[i].Push(v)
				}
			}

		case series.Int64s:
			for _, v := range s.Int64s()[:nRowsOut] {
				formatters[i].Push(v)
			}

			if addTail {
				for _, v := range s.Int64s()[df.NRows()-params.tailLen:] {
					formatters[i].Push(v)
				}
			}

		case series.Float64s:
			for _, v := range s.Float64s()[:nRowsOut] {
				formatters[i].Push(v)
			}

			if addTail {
				for _, v := range s.Float64s()[df.NRows()-params.tailLen:] {
					formatters[i].Push(v)
				}
			}

		case series.Strings:
			for _, v := range s.Strings()[:nRowsOut] {
				formatters[i].Push(v)
			}

			if addTail {
				for _, v := range s.Strings()[df.NRows()-params.tailLen:] {
					formatters[i].Push(v)
				}
			}

		case series.Times:
			for _, v := range s.DataAsString()[:nRowsOut] {
				formatters[i].Push(v)
			}

			if addTail {
				for _, v := range s.DataAsString()[df.NRows()-params.tailLen:] {
					formatters[i].Push(v)
				}
			}

		case series.Durations:
			for _, v := range s.Data().([]time.Duration)[:nRowsOut] {
				formatters[i].Push(v)
			}

			if addTail {
				for _, v := range s.Data().([]time.Duration)[df.NRows()-params.tailLen:] {
					formatters[i].Push(v)
				}
			}
		}
	}

	// compute the optimal width for each column with the given formatter
	// get the new number of columns that can fit in the screen
	actualWidthsSum = 0
	nColsOut = 0
	for i, f := range formatters {
		f.Compute()
		widths[i] = max(f.GetMaxWidth(), widths[i])
		actualWidthsSum += widths[i] + 3
		if actualWidthsSum > params.width {
			break
		}
		nColsOut++
	}
	widths = widths[:nColsOut]
	formatters = formatters[:nColsOut]

	// header
	buffer += params.indent + "╭"
	for i, w := range widths {
		for j := 0; j < w+2; j++ {
			buffer += "─"
		}
		if i < nColsOut-1 {
			buffer += "┬"
		}
	}
	buffer += "╮\n"

	// column names
	buffer += params.indent + "│"
	if params.useLipGloss {
		for i, name := range df.names[:nColsOut] {
			buffer += params.styleNames.Render(fmt.Sprintf(" %-*s ", widths[i], utils.Truncate(name, widths[i]))) + "│"
		}
	} else {
		for i, name := range df.names[:nColsOut] {
			buffer += fmt.Sprintf(" %-*s ", widths[i], utils.Truncate(name, widths[i])) + "│"
		}
	}
	buffer += "\n"

	// separator
	buffer += params.indent + "├"
	for i, w := range widths {
		for j := 0; j < w+2; j++ {
			buffer += "─"
		}
		if i < nColsOut-1 {
			buffer += "┼"
		}
	}
	buffer += "┤\n"

	// column types
	buffer += params.indent + "│"
	if params.useLipGloss {
		for i, c := range df.series[:nColsOut] {
			buffer += params.styleTypes.Render(fmt.Sprintf(" %-*s ", widths[i], c.Type().String())) + "│"
		}
	} else {
		for i, c := range df.series[:nColsOut] {
			buffer += fmt.Sprintf(" %-*s ", widths[i], c.Type().String()) + "│"
		}
	}
	buffer += "\n"

	// separator
	buffer += params.indent + "├"
	for i, w := range widths {
		for j := 0; j < w+2; j++ {
			buffer += "─"
		}
		if i < nColsOut-1 {
			buffer += "┼"
		}
	}
	buffer += "┤\n"

	// data
	for i := 0; i < nRowsOut; i++ {
		buffer += params.indent + "│"
		for j, c := range df.series[:nColsOut] {
			switch s := c.(type) {
			case series.Bools:
				buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.GetAsString(i), s.IsNull(i))) + "│"
			case series.Ints:
				buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
			case series.Int64s:
				buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
			case series.Float64s:
				buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
			case series.Strings:
				buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
			case series.Times:
				buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.GetAsString(i), s.IsNull(i))) + "│"
			case series.Durations:
				buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
			}
		}
		buffer += "\n"
	}

	if addTail {
		// separator (bottom)
		buffer += params.indent + "┊"
		for j, _ := range df.series[:nColsOut] {
			buffer += utils.Center("⋮", widths[j]+2) + "┊"
		}
		buffer += "\n"

		// tail
		for i := df.NRows() - params.tailLen; i < df.NRows(); i++ {
			buffer += params.indent + "│"
			for j, c := range df.series[:nColsOut] {
				switch s := c.(type) {
				case series.Bools:
					buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.GetAsString(i), s.IsNull(i))) + "│"
				case series.Ints:
					buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
				case series.Int64s:
					buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
				case series.Float64s:
					buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
				case series.Strings:
					buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
				case series.Times:
					buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.GetAsString(i), s.IsNull(i))) + "│"
				case series.Durations:
					buffer += fmt.Sprintf(" %s ", formatters[j].Format(widths[j], s.Get(i), s.IsNull(i))) + "│"
				}
			}
			buffer += "\n"
		}
	}

	// end
	buffer += params.indent + "╰"
	for i, w := range widths {
		for j := 0; j < w+2; j++ {
			buffer += "─"
		}
		if i < nColsOut-1 {
			buffer += "┴"
		}
	}
	buffer += "╯\n"

	// Non-displayed column names
	if df.NCols() > nColsOut {
		for i := nColsOut; i < df.NCols(); i++ {
			buffer += df.names[i] + ", "
		}
	}

	fmt.Println(buffer)

	return df
}
