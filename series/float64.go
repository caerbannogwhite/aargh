package series

import (
	"fmt"
	"sort"
	"time"
	"unsafe"

	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/meta"
)

// Float64s represents a series of floats.
type Float64s struct {
	IsNullable_ bool
	Sorted_     aargh.SeriesSortOrder
	Data_       []float64
	NullMask_   []uint8
	Partition_  *SeriesFloat64Partition
	Ctx_        *aargh.Context
}

// Get the element at index i as a string.
func (s Float64s) GetAsString(i int) string {
	if s.IsNullable_ && s.IsNull(i) {
		return aargh.NA_TEXT
	}
	return floatToString(s.Data_[i])
}

// Set the element at index i. The value v can be any belonging to types:
// int8, int16, int, int, int64, float32, float64 and their nullable versions.
func (s Float64s) Set(i int, v any) Series {
	if s.Partition_ != nil {
		return Errors{"Float64s.Set: cannot set values in a grouped series"}
	}

	switch val := v.(type) {
	case nil:
		s = s.MakeNullable().(Float64s)
		s.NullMask_[i>>3] |= 1 << uint(i%8)

	case int8:
		s.Data_[i] = float64(val)

	case int16:
		s.Data_[i] = float64(val)

	case int:
		s.Data_[i] = float64(val)

	case int32:
		s.Data_[i] = float64(val)

	case int64:
		s.Data_[i] = float64(val)

	case float32:
		s.Data_[i] = float64(val)

	case float64:
		s.Data_[i] = val

	case aargh.NullableInt8:
		s = s.MakeNullable().(Float64s)
		if v.(aargh.NullableInt8).Valid {
			s.Data_[i] = float64(val.Value)
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	case aargh.NullableInt16:
		s = s.MakeNullable().(Float64s)
		if v.(aargh.NullableInt16).Valid {
			s.Data_[i] = float64(val.Value)
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	case aargh.NullableInt:
		s = s.MakeNullable().(Float64s)
		if v.(aargh.NullableInt).Valid {
			s.Data_[i] = float64(val.Value)
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	case aargh.NullableInt64:
		s = s.MakeNullable().(Float64s)
		if v.(aargh.NullableInt64).Valid {
			s.Data_[i] = float64(val.Value)
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	case aargh.NullableFloat32:
		s = s.MakeNullable().(Float64s)
		if v.(aargh.NullableFloat32).Valid {
			s.Data_[i] = float64(val.Value)
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	case aargh.NullableFloat64:
		s = s.MakeNullable().(Float64s)
		if v.(aargh.NullableFloat64).Valid {
			s.Data_[i] = val.Value
		} else {
			s.Data_[i] = 0
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	default:
		return Errors{fmt.Sprintf("Float64s.Set: invalid type %T", v)}
	}

	s.Sorted_ = aargh.SORTED_NONE
	return s
}

////////////////////////			ALL DATA ACCESSORS

// Return the underlying Data_ as a slice of float64.
func (s Float64s) Float64s() []float64 {
	return s.Data_
}

// Return the underlying Data_ as a slice of NullableFloat64.
func (s Float64s) DataAsNullable() any {
	Data_ := make([]aargh.NullableFloat64, len(s.Data_))
	for i, v := range s.Data_ {
		Data_[i] = aargh.NullableFloat64{Valid: !s.IsNull(i), Value: v}
	}
	return Data_
}

// Return the underlying Data_ as a slice of strings.
func (s Float64s) DataAsString() []string {
	Data_ := make([]string, len(s.Data_))
	if s.IsNullable_ {
		for i, v := range s.Data_ {
			if s.IsNull(i) {
				Data_[i] = aargh.NA_TEXT
			} else {
				Data_[i] = floatToString(v)
			}
		}
	} else {
		for i, v := range s.Data_ {
			Data_[i] = floatToString(v)
		}
	}
	return Data_
}

// Casts the series to a given type.
func (s Float64s) Cast(t meta.BaseType) Series {
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
		Data_ := make([]int64, len(s.Data_))
		for i, v := range s.Data_ {
			Data_[i] = int64(v)
		}

		return Int64s{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.Float64Type:
		return s

	case meta.StringType:
		Data_ := make([]*string, len(s.Data_))
		if s.IsNullable_ {
			for i, v := range s.Data_ {
				if s.IsNull(i) {
					Data_[i] = s.Ctx_.StringPool.Put(aargh.NA_TEXT)
				} else {
					Data_[i] = s.Ctx_.StringPool.Put(floatToString(v))
				}
			}
		} else {
			for i, v := range s.Data_ {
				Data_[i] = s.Ctx_.StringPool.Put(floatToString(v))
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
			Data_[i] = time.Unix(0, int64(v))
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
		return Errors{fmt.Sprintf("Float64s.Cast: invalid type %s", t.ToString())}
	}
}

////////////////////////			GROUPING OPERATIONS

// A SeriesFloat64Partition is a Partition_ of a Float64s.
// Each key is a hash of a bool value, and each value is a slice of indices
// of the original series that are set to that value.
type SeriesFloat64Partition struct {
	Partition_   map[int64][]int
	indexToGroup []int
}

func (gp *SeriesFloat64Partition) GetSize() int {
	return len(gp.Partition_)
}

func (gp *SeriesFloat64Partition) GetMap() map[int64][]int {
	return gp.Partition_
}

func (s Float64s) Group() Series {

	// Define the worker callback
	worker := func(threadNum, start, end int, map_ map[int64][]int) {
		for i := start; i < end; i++ {
			map_[*(*int64)(unsafe.Pointer((&s.Data_[i])))] = append(map_[*(*int64)(unsafe.Pointer((&s.Data_[i])))], i)
		}
	}

	// Define the worker callback for nulls
	workerNulls := func(threadNum, start, end int, map_ map[int64][]int, nulls *[]int) {
		for i := start; i < end; i++ {
			if s.IsNull(i) {
				(*nulls) = append((*nulls), i)
			} else {
				map_[*(*int64)(unsafe.Pointer((&s.Data_[i])))] = append(map_[*(*int64)(unsafe.Pointer((&s.Data_[i])))], i)
			}
		}
	}

	Partition_ := SeriesFloat64Partition{
		Partition_: __series_groupby(
			aargh.THREADS_NUMBER, aargh.MINIMUM_PARALLEL_SIZE_2, len(s.Data_), s.HasNull(),
			worker, workerNulls),
	}

	s.Partition_ = &Partition_

	return s
}

func (s Float64s) GroupBy(Partition_ SeriesPartition) Series {
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
				newHash = *(*int64)(unsafe.Pointer((&(s.Data_)[index]))) + aargh.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
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
					newHash = *(*int64)(unsafe.Pointer((&(s.Data_)[index]))) + aargh.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
				}
				map_[newHash] = append(map_[newHash], index)
			}
		}
	}

	newPartition := SeriesFloat64Partition{
		Partition_: __series_groupby(
			aargh.THREADS_NUMBER, aargh.MINIMUM_PARALLEL_SIZE_1, len(keys), s.HasNull(),
			worker, workerNulls),
	}

	s.Partition_ = &newPartition

	return s
}

////////////////////////			SORTING OPERATIONS

func (s Float64s) Less(i, j int) bool {
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

func (s Float64s) Equal(i, j int) bool {
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

func (s Float64s) Swap(i, j int) {
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

func (s Float64s) Sort() Series {
	if s.Sorted_ != aargh.SORTED_ASC {
		sort.Sort(s)
		s.Sorted_ = aargh.SORTED_ASC
	}
	return s
}

func (s Float64s) SortRev() Series {
	if s.Sorted_ != aargh.SORTED_DESC {
		sort.Sort(sort.Reverse(s))
		s.Sorted_ = aargh.SORTED_DESC
	}
	return s
}

////////////////////////			NUMERIC OPERATIONS

func (s Float64s) Min() any {
	if s.IsNullable_ {
		return aargh.NullableFloat64{Valid: false, Value: 0}
	}

	min := s.Data_[0]
	for _, v := range s.Data_ {
		if v < min {
			min = v
		}
	}

	return min
}
