package series

import (
	"fmt"
	"sort"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/meta"
)

// Bools represents a series of bools.
type Bools struct {
	IsNullable_ bool
	Sorted_     aargh.SeriesSortOrder
	Data_       []bool
	NullMask_   []uint8
	Partition_  *SeriesBoolPartition
	Ctx_        *aargh.Context
	arr_        arrow.Array
}

func (s Bools) ArrowArray() arrow.Array {
	if s.arr_ != nil {
		return s.arr_
	}
	return buildArrowBoolean(s.Ctx_.Allocator, s.Data_, s.IsNullable_, s.NullMask_)
}

func (s *Bools) invalidateArrow() { s.arr_ = nil }

// Get the element at index i as a string.
func (s Bools) GetAsString(i int) string {
	if s.IsNullable_ && s.NullMask_[i>>3]&(1<<uint(i%8)) != 0 {
		return aargh.NA_TEXT
	} else if s.Data_[i] {
		return aargh.BOOL_TRUE_TEXT
	} else {
		return aargh.BOOL_FALSE_TEXT
	}
}

// Set the element at index i. The value must be of type bool or NullableBool.
func (s Bools) Set(i int, v any) Series {
	if s.Partition_ != nil {
		return Errors{"Bools.Set: cannot set values in a grouped series"}
	}

	switch v := v.(type) {
	case nil:
		s = s.MakeNullable().(Bools)
		s.NullMask_[i>>3] |= 1 << uint(i%8)

	case bool:
		s.Data_[i] = v

	case aargh.NullableBool:
		s = s.MakeNullable().(Bools)
		if v.Valid {
			s.Data_[i] = v.Value
		} else {
			s.NullMask_[i>>3] |= 1 << uint(i%8)
			s.Data_[i] = false
		}

	default:
		return Errors{fmt.Sprintf("Bools.Set: invalid type %T", v)}
	}

	s.Sorted_ = aargh.SORTED_NONE
	return s
}

////////////////////////			ALL DATA ACCESSORS

// Return the underlying Data_ as a slice of bools.
func (s Bools) Bools() []bool {
	return s.Data_
}

// Return the underlying Data_ as a slice of NullableBool.
func (s Bools) DataAsNullable() any {
	Data_ := make([]aargh.NullableBool, len(s.Data_))
	for i, v := range s.Data_ {
		Data_[i] = aargh.NullableBool{Valid: !s.IsNull(i), Value: v}
	}
	return Data_
}

// Return the Data_ as a slice of strings.
func (s Bools) DataAsString() []string {
	Data_ := make([]string, len(s.Data_))
	if s.IsNullable_ {
		for i, v := range s.Data_ {
			if s.IsNull(i) {
				Data_[i] = aargh.NA_TEXT
			} else if v {
				Data_[i] = aargh.BOOL_TRUE_TEXT
			} else {
				Data_[i] = aargh.BOOL_FALSE_TEXT
			}
		}
	} else {
		for i, v := range s.Data_ {
			if v {
				Data_[i] = aargh.BOOL_TRUE_TEXT
			} else {
				Data_[i] = aargh.BOOL_FALSE_TEXT
			}
		}
	}
	return Data_
}

// Cast the series to a given type.
func (s Bools) Cast(t meta.BaseType) Series {
	switch t {
	case meta.BoolType:
		return s

	case meta.IntType:
		Data_ := make([]int, len(s.Data_))
		for i, v := range s.Data_ {
			if v {
				Data_[i] = 1
			}
		}

		return Ints{
			IsNullable_: s.IsNullable_,
			Sorted_:     s.Sorted_,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.Int64Type:
		Data_ := make([]int64, len(s.Data_))
		for i, v := range s.Data_ {
			if v {
				Data_[i] = 1
			}
		}

		return Int64s{
			IsNullable_: s.IsNullable_,
			Sorted_:     s.Sorted_,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.Float64Type:
		Data_ := make([]float64, len(s.Data_))
		for i, v := range s.Data_ {
			if v {
				Data_[i] = 1
			}
		}

		return Float64s{
			IsNullable_: s.IsNullable_,
			Sorted_:     s.Sorted_,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.StringType:
		Data_ := make([]*string, len(s.Data_))

		naTextPtr := s.Ctx_.StringPool.Put(aargh.NA_TEXT)
		trueTextPtr := s.Ctx_.StringPool.Put(aargh.BOOL_TRUE_TEXT)
		falseTextPtr := s.Ctx_.StringPool.Put(aargh.BOOL_FALSE_TEXT)

		if s.IsNullable_ {
			for i, v := range s.Data_ {
				if s.IsNull(i) {
					Data_[i] = naTextPtr
				} else if v {
					Data_[i] = trueTextPtr
				} else {
					Data_[i] = falseTextPtr
				}
			}
		} else {
			for i, v := range s.Data_ {
				if v {
					Data_[i] = trueTextPtr
				} else {
					Data_[i] = falseTextPtr
				}
			}
		}

		return Strings{
			IsNullable_: s.IsNullable_,
			Sorted_:     s.Sorted_,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	default:
		return Errors{fmt.Sprintf("Bools.Cast: invalid type %s", t.String())}
	}
}

////////////////////////			GROUPING OPERATIONS

// A SeriesBoolPartition is a Partition_ of a Bools.
// Each key is a hash of a bool value, and each value is a slice of indices
// of the original series that are set to that value.
type SeriesBoolPartition struct {
	Partition_ map[int64][]int
}

func (gp *SeriesBoolPartition) GetSize() int {
	return len(gp.Partition_)
}

func (gp *SeriesBoolPartition) GetMap() map[int64][]int {
	return gp.Partition_
}

func (s Bools) Group() Series {

	// Define the worker callback
	worker := func(threadNum, start, end int, map_ map[int64][]int) {
		for i := start; i < end; i++ {
			if s.Data_[i] {
				map_[1] = append(map_[1], i)
			} else {
				map_[0] = append(map_[0], i)
			}
		}
	}

	// Define the worker callback for nulls
	workerNulls := func(threadNum, start, end int, map_ map[int64][]int, nulls *[]int) {
		for i := start; i < end; i++ {
			if s.IsNull(i) {
				(*nulls) = append((*nulls), i)
			} else if s.Data_[i] {
				map_[1] = append(map_[1], i)
			} else {
				map_[0] = append(map_[0], i)
			}

		}
	}

	Partition_ := SeriesBoolPartition{
		Partition_: __series_groupby(
			aargh.THREADS_NUMBER, aargh.MINIMUM_PARALLEL_SIZE_1, s.Len(), s.HasNull(),
			worker, workerNulls),
	}

	s.Partition_ = &Partition_

	return s
}

func (s Bools) GroupBy(Partition_ SeriesPartition) Series {
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
				if s.Data_[index] {
					newHash = (1 + aargh.HASH_MAGIC_NUMBER) + (h << 13) + (h >> 4)
				} else {
					newHash = aargh.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
				}
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
				} else if s.Data_[index] {
					newHash = (1 + aargh.HASH_MAGIC_NUMBER) + (h << 13) + (h >> 4)
				} else {
					newHash = aargh.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
				}
				map_[newHash] = append(map_[newHash], index)
			}
		}
	}

	newPartition := SeriesBoolPartition{
		Partition_: __series_groupby(
			aargh.THREADS_NUMBER, aargh.MINIMUM_PARALLEL_SIZE_1, len(keys), s.HasNull(),
			worker, workerNulls),
	}

	s.Partition_ = &newPartition

	return s
}

////////////////////////			SORTING OPERATIONS

func (s Bools) Less(i, j int) bool {
	if s.IsNullable_ {
		if s.NullMask_[i>>3]&(1<<uint(i%8)) > 0 {
			return false
		}
		if s.NullMask_[j>>3]&(1<<uint(j%8)) > 0 {
			return true
		}
	}
	return !s.Data_[i] && s.Data_[j]
}

func (s Bools) Equal(i, j int) bool {
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

func (s Bools) Swap(i, j int) {
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

func (s Bools) Sort() Series {
	if s.Sorted_ != aargh.SORTED_ASC {
		sort.Sort(s)
		s.Sorted_ = aargh.SORTED_ASC
	}
	return s
}

func (s Bools) SortRev() Series {
	if s.Sorted_ != aargh.SORTED_DESC {
		sort.Sort(sort.Reverse(s))
		s.Sorted_ = aargh.SORTED_DESC
	}
	return s
}
