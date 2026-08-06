package series

import (
	"fmt"
	"sort"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/meta"
)

// Times represents a datetime series.
type Times struct {
	IsNullable_ bool
	Sorted_     enchanter.SeriesSortOrder
	Data_       []time.Time
	NullMask_   []uint8
	Partition_  *SeriesTimePartition
	Ctx_        *enchanter.Context
	timeFormat  string
}

// ArrowArray builds and returns a fresh Arrow array from the series data.
// The caller owns the returned array; releasing it is optional under
// GC-backed allocators (see enchanter.Context.Allocator).
func (s Times) ArrowArray() arrow.Array {
	return buildArrowTimestamp(s.Ctx_.Allocator, s.Data_, s.IsNullable_, s.NullMask_)
}

// Get the time format of the series.
func (s Times) GetTimeFormat() string {
	return s.timeFormat
}

// Set the time format of the series.
func (s Times) SetTimeFormat(format string) Series {
	s.timeFormat = format
	return s
}

// Get the element at index i as a string.
func (s Times) GetAsString(i int) string {
	if s.IsNullable_ && s.NullMask_[i>>3]&(1<<uint(i%8)) != 0 {
		return enchanter.NA_TEXT
	}
	return s.Data_[i].Format(s.timeFormat)
}

// Set the element at index i. The value v must be of type time.Time or NullableTime.
func (s Times) Set(i int, v any) Series {
	if s.Partition_ != nil {
		return Errors{"Times.Set: cannot set values on a grouped Series"}
	}

	switch v := v.(type) {
	case nil:
		s = s.MakeNullable().(Times)
		s.NullMask_[i>>3] |= 1 << uint(i%8)

	case time.Time:
		s.Data_[i] = v

	case enchanter.NullableTime:
		s = s.MakeNullable().(Times)
		if v.Valid {
			s.Data_[i] = v.Value
		} else {
			s.Data_[i] = time.Time{}
			s.NullMask_[i/8] |= 1 << uint(i%8)
		}

	default:
		return Errors{fmt.Sprintf("Times.Set: invalid type %T", v)}
	}

	s.Sorted_ = enchanter.SORTED_NONE
	return s
}

////////////////////////			ALL DATA ACCESSORS

// Return the underlying Data_ as a slice of time.Time.
func (s Times) Times() []time.Time {
	return s.Data_
}

// Return the underlying Data_ as a slice of NullableTime.
func (s Times) DataAsNullable() any {
	Data_ := make([]enchanter.NullableTime, len(s.Data_))
	for i, v := range s.Data_ {
		Data_[i] = enchanter.NullableTime{Valid: !s.IsNull(i), Value: v}
	}
	return Data_
}

// Return the underlying Data_ as a slice of strings.
func (s Times) DataAsString() []string {
	Data_ := make([]string, len(s.Data_))
	if s.IsNullable_ {
		for i, v := range s.Data_ {
			if s.IsNull(i) {
				Data_[i] = enchanter.NA_TEXT
			} else {
				Data_[i] = v.Format(s.timeFormat)
			}
		}
	} else {
		for i, v := range s.Data_ {
			Data_[i] = v.Format(s.timeFormat)
		}
	}
	return Data_
}

// Casts the series to a given type.
func (s Times) Cast(t meta.BaseType) Series {
	switch t {
	case meta.BoolType:
		return Errors{fmt.Sprintf("Times.Cast: cannot cast to %s", t.String())}

	case meta.IntType:
		Data_ := make([]int, len(s.Data_))
		for i, v := range s.Data_ {
			Data_[i] = int(v.UnixNano())
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
			Data_[i] = v.UnixNano()
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
			Data_[i] = float64(v.UnixNano())
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
		if s.IsNullable_ {
			for i, v := range s.Data_ {
				if s.IsNull(i) {
					Data_[i] = s.Ctx_.StringPool.Put(enchanter.NA_TEXT)
				} else {
					Data_[i] = s.Ctx_.StringPool.Put(v.Format(s.timeFormat))
				}
			}
		} else {
			for i, v := range s.Data_ {
				Data_[i] = s.Ctx_.StringPool.Put(v.Format(s.timeFormat))
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

	case meta.TimeType:
		return s

	case meta.DurationType:
		Data_ := make([]time.Duration, len(s.Data_))
		for i, v := range s.Data_ {
			Data_[i] = v.Sub(time.Time{})
		}

		return Durations{
			IsNullable_: s.IsNullable_,
			Sorted_:     s.Sorted_,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	default:
		return Errors{fmt.Sprintf("Times.Cast: invalid type %T", t)}
	}
}

////////////////////////			GROUPING OPERATIONS

// A SeriesTimePartition is a Partition_ of a Times.
// Each key is a hash of a bool value, and each value is a slice of indices
// of the original series that are set to that value.
type SeriesTimePartition struct {
	Partition_ map[int64][]int
}

func (gp *SeriesTimePartition) GetSize() int {
	return len(gp.Partition_)
}

func (gp *SeriesTimePartition) GetMap() map[int64][]int {
	return gp.Partition_
}

func (s Times) Group() Series {

	// Define the worker callback
	worker := func(threadNum, start, end int, map_ map[int64][]int) {
		for i := start; i < end; i++ {
			map_[s.Data_[i].UnixNano()] = append(map_[s.Data_[i].UnixNano()], i)
		}
	}

	// Define the worker callback for nulls
	workerNulls := func(threadNum, start, end int, map_ map[int64][]int, nulls *[]int) {
		for i := start; i < end; i++ {
			if s.IsNull(i) {
				(*nulls) = append((*nulls), i)
			} else {
				map_[s.Data_[i].UnixNano()] = append(map_[s.Data_[i].UnixNano()], i)
			}
		}
	}

	Partition_ := SeriesTimePartition{
		Partition_: __series_groupby(
			enchanter.THREADS_NUMBER, enchanter.MINIMUM_PARALLEL_SIZE_1, s.Len(), s.HasNull(),
			worker, workerNulls),
	}

	s.Partition_ = &Partition_

	return s
}

func (s Times) GroupBy(Partition_ SeriesPartition) Series {
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
				newHash = s.Data_[index].UnixNano() + enchanter.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
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
					newHash = enchanter.HASH_MAGIC_NUMBER_NULL + (h << 13) + (h >> 4)
				} else {
					newHash = s.Data_[index].UnixNano() + enchanter.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
				}
				map_[newHash] = append(map_[newHash], index)
			}
		}
	}

	newPartition := SeriesTimePartition{
		Partition_: __series_groupby(
			enchanter.THREADS_NUMBER, enchanter.MINIMUM_PARALLEL_SIZE_1, len(keys), s.HasNull(),
			worker, workerNulls),
	}

	s.Partition_ = &newPartition

	return s
}

////////////////////////			SORTING OPERATIONS

func (s Times) Less(i, j int) bool {
	if s.IsNullable_ {
		if s.NullMask_[i>>3]&(1<<uint(i%8)) > 0 {
			return false
		}
		if s.NullMask_[j>>3]&(1<<uint(j%8)) > 0 {
			return true
		}
	}
	return s.Data_[i].Compare(s.Data_[j]) < 0
}

func (s Times) Equal(i, j int) bool {
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

func (s Times) Swap(i, j int) {
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

func (s Times) Sort() Series {
	if s.Sorted_ != enchanter.SORTED_ASC {
		sort.Sort(s)
		s.Sorted_ = enchanter.SORTED_ASC
	}
	return s
}

func (s Times) SortRev() Series {
	if s.Sorted_ != enchanter.SORTED_DESC {
		sort.Sort(sort.Reverse(s))
		s.Sorted_ = enchanter.SORTED_DESC
	}
	return s
}
