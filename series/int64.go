package series

import (
	"fmt"
	"sort"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/meta"
)

// Int64s represents a series of ints.
type Int64s struct {
	IsNullable_ bool
	Sorted_     aargh.SeriesSortOrder
	Data_       []int64
	NullMask_   []uint8
	Partition_  *SeriesInt64Partition
	Ctx_        *aargh.Context
	arr_        arrow.Array
}

func (s Int64s) ArrowArray() arrow.Array {
	if s.arr_ != nil {
		return s.arr_
	}
	return buildArrowInt64(s.Ctx_.Allocator, s.Data_, s.IsNullable_, s.NullMask_)
}

func (s *Int64s) invalidateArrow() { s.arr_ = nil }

// Get the element at index i as a string.
func (s Int64s) GetAsString(i int) string {
	if s.IsNullable_ && s.IsNull(i) {
		return aargh.NA_TEXT
	}
	return intToString(s.Data_[i])
}

// Set the element at index i. The value v can be any belonging to types:
// int8, int16, int, int, int64 and their nullable versions.
func (s Int64s) Set(i int, v any) Series {
	if s.Partition_ != nil {
		return Errors{"Int64s.Set: cannot set values on a grouped Series"}
	}

	switch val := v.(type) {
	case nil:
		s = s.MakeNullable().(Int64s)
		s.NullMask_[i>>3] |= 1 << uint(i%8)

	case int8:
		s.Data_[i] = int64(val)

	case int16:
		s.Data_[i] = int64(val)

	case int:
		s.Data_[i] = int64(val)

	case int32:
		s.Data_[i] = int64(val)

	case int64:
		s.Data_[i] = val

	case aargh.NullableInt8:
		s = s.MakeNullable().(Int64s)
		if v.(aargh.NullableInt8).Valid {
			s.Data_[i] = int64(val.Value)
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	case aargh.NullableInt16:
		s = s.MakeNullable().(Int64s)
		if v.(aargh.NullableInt16).Valid {
			s.Data_[i] = int64(val.Value)
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	case aargh.NullableInt:
		s = s.MakeNullable().(Int64s)
		if v.(aargh.NullableInt).Valid {
			s.Data_[i] = int64(val.Value)
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	case aargh.NullableInt64:
		s = s.MakeNullable().(Int64s)
		if v.(aargh.NullableInt64).Valid {
			s.Data_[i] = val.Value
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	default:
		return Errors{fmt.Sprintf("Int64s.Set: invalid type %T", v)}
	}

	s.Sorted_ = aargh.SORTED_NONE
	return s
}

////////////////////////			ALL DATA ACCESSORS

// Return the underlying Data_ as a slice of int64.
func (s Int64s) Int64s() []int64 {
	return s.Data_
}

// Return the underlying Data_ as a slice of NullableInt64.
func (s Int64s) DataAsNullable() any {
	Data_ := make([]aargh.NullableInt64, len(s.Data_))
	for i, v := range s.Data_ {
		Data_[i] = aargh.NullableInt64{Valid: !s.IsNull(i), Value: v}
	}
	return Data_
}

// Return the underlying Data_ as a slice of strings.
func (s Int64s) DataAsString() []string {
	Data_ := make([]string, len(s.Data_))
	if s.IsNullable_ {
		for i, v := range s.Data_ {
			if s.IsNull(i) {
				Data_[i] = aargh.NA_TEXT
			} else {
				Data_[i] = intToString(v)
			}
		}
	} else {
		for i, v := range s.Data_ {
			Data_[i] = intToString(v)
		}
	}
	return Data_
}

// Casts the series to a given type.
func (s Int64s) Cast(t meta.BaseType) Series {
	switch t {
	case meta.BoolType:
		Data_ := make([]bool, len(s.Data_))
		for i, v := range s.Data_ {
			Data_[i] = v != 0
		}

		return Bools{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.IntType:
		Data_ := make([]int, len(s.Data_))
		for i, v := range s.Data_ {
			Data_[i] = int(v)
		}

		return Ints{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.Int64Type:
		return s

	case meta.Float64Type:
		Data_ := make([]float64, len(s.Data_))
		for i, v := range s.Data_ {
			Data_[i] = float64(v)
		}

		return Float64s{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.StringType:
		if s.Ctx_.StringPool == nil {
			return Errors{"Int64s.Cast: StringPool is nil"}
		}

		Data_ := make([]*string, len(s.Data_))
		if s.IsNullable_ {
			for i, v := range s.Data_ {
				if s.IsNull(i) {
					Data_[i] = s.Ctx_.StringPool.Put(aargh.NA_TEXT)
				} else {
					Data_[i] = s.Ctx_.StringPool.Put(intToString(v))
				}
			}
		} else {
			for i, v := range s.Data_ {
				Data_[i] = s.Ctx_.StringPool.Put(intToString(v))
			}
		}

		return Strings{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.TimeType:
		Data_ := make([]time.Time, len(s.Data_))
		for i, v := range s.Data_ {
			Data_[i] = time.Unix(0, v)
		}

		return Times{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.DurationType:
		Data_ := make([]time.Duration, len(s.Data_))
		for i, v := range s.Data_ {
			Data_[i] = time.Duration(v)
		}

		return Durations{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	default:
		return Errors{fmt.Sprintf("Int64s.Cast: invalid type %s", t.String())}
	}
}

////////////////////////			GROUPING OPERATIONS

// A SeriesInt64Partition is a Partition_ of a Int64s.
// Each key is a hash of a bool value, and each value is a slice of indices
// of the original series that are set to that value.
type SeriesInt64Partition struct {
	Partition_           map[int64][]int
	isDense              bool
	Partition_DenseMin   int64
	Partition_Dense      [][]int
	Partition_DenseNulls []int
}

func (gp *SeriesInt64Partition) GetSize() int {
	if gp.isDense {
		if gp.Partition_DenseNulls != nil && len(gp.Partition_DenseNulls) > 0 {
			return len(gp.Partition_Dense) + 1
		}
		return len(gp.Partition_Dense)
	}
	return len(gp.Partition_)
}

func (gp *SeriesInt64Partition) GetMap() map[int64][]int {
	if gp.isDense {
		map_ := make(map[int64][]int, len(gp.Partition_Dense))
		for i, part := range gp.Partition_Dense {
			map_[int64(i)+gp.Partition_DenseMin] = part
		}

		// Merge the nulls to the map
		if gp.Partition_DenseNulls != nil && len(gp.Partition_DenseNulls) > 0 {
			nullKey := __series_get_nullkey(map_, aargh.HASH_NULL_KEY)
			map_[nullKey] = gp.Partition_DenseNulls
		}

		return map_
	}

	return gp.Partition_
}

func (s Int64s) Group() Series {
	var useDenseMap bool
	var min, max int64
	var Partition_ SeriesInt64Partition

	// If the number of elements is small,
	// look for the minimum and maximum values
	// if len(s.Data_) < MINIMUM_PARALLEL_SIZE_2 {
	// 	useDenseMap = true
	// 	max = s.Data_[0]
	// 	min = s.Data_[0]
	// 	for _, v := range s.Data_ {
	// 		if v > max {
	// 			max = v
	// 		}
	// 		if v < min {
	// 			min = v
	// 		}
	// 	}
	// }

	// If the difference between the maximum and minimum values is acceptable,
	// then we can use a dense map, otherwise we use a sparse map
	if useDenseMap && (max-min >= aargh.MINIMUM_PARALLEL_SIZE_1) {
		useDenseMap = false
	}

	// TODO: FIX DENSE MAP
	// if useDenseMap {
	// 	var nulls []int
	// 	map_ := make([][]int, max-min+1)
	// 	for i := 0; i < len(map_); i++ {
	// 		map_[i] = make([]int, 0, DEFAULT_DENSE_MAP_ARRAY_INITIAL_CAPACITY)
	// 	}

	// 	if s.HasNull() {
	// 		nulls = make([]int, 0, DEFAULT_DENSE_MAP_ARRAY_INITIAL_CAPACITY)
	// 		for i, v := range s.Data_ {
	// 			if s.IsNull(i) {
	// 				nulls = append(nulls, i)
	// 			} else {
	// 				map_[v-min] = append(map_[v-min], i)
	// 			}
	// 		}
	// 	} else {
	// 		for i, v := range s.Data_ {
	// 			map_[v-min] = append(map_[v-min], i)
	// 		}
	// 	}

	// 	Partition_ = SeriesInt64Partition{
	// 		isDense:             true,
	// 		Partition_DenseMin:   min,
	// 		Partition_Dense:      map_,
	// 		Partition_DenseNulls: nulls,
	// 	}
	// } else

	// SPARSE MAP
	{
		// Define the worker callback
		worker := func(threadNum, start, end int, map_ map[int64][]int) {
			up := end - ((end - start) % 8)
			for i := start; i < up; {
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				i++
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				i++
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				i++
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				i++
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				i++
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				i++
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				i++
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				i++
			}

			for i := up; i < end; i++ {
				map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
			}
		}

		// Define the worker callback for nulls
		workerNulls := func(threadNum, start, end int, map_ map[int64][]int, nulls *[]int) {
			for i := start; i < end; i++ {
				if s.IsNull(i) {
					(*nulls) = append((*nulls), i)
				} else {
					map_[s.Data_[i]] = append(map_[s.Data_[i]], i)
				}
			}
		}

		Partition_ = SeriesInt64Partition{
			isDense: false,
			Partition_: __series_groupby(
				aargh.THREADS_NUMBER, aargh.MINIMUM_PARALLEL_SIZE_2, len(s.Data_), s.HasNull(),
				worker, workerNulls),
		}
	}

	s.Partition_ = &Partition_

	return s
}

func (s Int64s) GroupBy(Partition_ SeriesPartition) Series {
	if Partition_ == nil {
		return s
	}

	// collect all keys
	otherIndeces := Partition_.GetMap()
	keys := make([]int64, len(otherIndeces))
	i := 0
	for k := range otherIndeces {
		keys[i] = k
		i++
	}

	// Define the worker callback
	worker := func(threadNum, start, end int, map_ map[int64][]int) {
		var newHash int64
		for _, h := range keys[start:end] { // keys is defined outside the function
			for _, index := range otherIndeces[h] { // otherIndeces is defined outside the function
				newHash = s.Data_[index] + aargh.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
				map_[newHash] = append(map_[newHash], index)
			}
		}
	}

	// Define the worker callback for nulls
	workerNulls := func(threadNum, start, end int, map_ map[int64][]int, nulls *[]int) {
		var newHash int64
		for _, h := range keys[start:end] { // keys is defined outside the function
			for _, index := range otherIndeces[h] { // otherIndeces is defined outside the function
				if s.IsNull(index) {
					newHash = aargh.HASH_MAGIC_NUMBER_NULL + (h << 13) + (h >> 4)
				} else {
					newHash = s.Data_[index] + aargh.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
				}
				map_[newHash] = append(map_[newHash], index)
			}
		}
	}

	newPartition := SeriesInt64Partition{
		Partition_: __series_groupby(
			aargh.THREADS_NUMBER, aargh.MINIMUM_PARALLEL_SIZE_2, len(keys), s.HasNull(),
			worker, workerNulls),
	}

	s.Partition_ = &newPartition

	return s
}

////////////////////////			SORTING OPERATIONS

func (s Int64s) Less(i, j int) bool {
	if s.IsNullable_ {
		if s.NullMask_[i>>3]&(1<<uint(i%8)) > 0 {
			return false
		}
		if s.NullMask_[j>>3]&(1<<uint(j%8)) > 0 {
			return true
		}
	}

	return s.Data_[i] < s.Data_[j]
}

func (s Int64s) Equal(i, j int) bool {
	if s.IsNullable_ {
		if (s.NullMask_[i>>3] & (1 << uint(i%8))) > 0 {
			return (s.NullMask_[j>>3] & (1 << uint(j%8))) > 0
		}
		if (s.NullMask_[j>>3] & (1 << uint(j%8))) > 0 {
			return false
		}
	}

	return s.Data_[i] == s.Data_[j]
}

func (s Int64s) Swap(i, j int) {
	if s.IsNullable_ {
		// i is null, j is not null
		if s.NullMask_[i>>3]&(1<<uint(i%8)) > 0 && s.NullMask_[j>>3]&(1<<uint(j%8)) == 0 {
			s.NullMask_[i>>3] &= ^(1 << uint(i%8))
			s.NullMask_[j>>3] |= 1 << uint(j%8)
		} else

		// i is not null, j is null
		if s.NullMask_[i>>3]&(1<<uint(i%8)) == 0 && s.NullMask_[j>>3]&(1<<uint(j%8)) > 0 {
			s.NullMask_[i>>3] |= 1 << uint(i%8)
			s.NullMask_[j>>3] &= ^(1 << uint(j%8))
		}
	}

	s.Data_[i], s.Data_[j] = s.Data_[j], s.Data_[i]
}

func (s Int64s) Sort() Series {
	if s.Sorted_ != aargh.SORTED_ASC {
		sort.Sort(s)
		s.Sorted_ = aargh.SORTED_ASC
	}
	return s
}

func (s Int64s) SortRev() Series {
	if s.Sorted_ != aargh.SORTED_DESC {
		sort.Sort(sort.Reverse(s))
		s.Sorted_ = aargh.SORTED_DESC
	}
	return s
}
