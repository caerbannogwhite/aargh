package series

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/meta"
	"github.com/caerbannogwhite/aargh/utils"
)

// Strings represents a series of strings.
type Strings struct {
	IsNullable_ bool
	Sorted_     aargh.SeriesSortOrder
	Data_       []*string
	NullMask_   []uint8
	Partition_  *SeriesStringPartition
	Ctx_        *aargh.Context
	arr_        arrow.Array
}

func (s Strings) ArrowArray() arrow.Array {
	if s.arr_ != nil {
		return s.arr_
	}
	return buildArrowString(s.Ctx_.Allocator, s.Data_, s.IsNullable_, s.NullMask_)
}

func (s *Strings) invalidateArrow() { s.arr_ = nil }

// Get the element at index i as a string.
func (s Strings) GetAsString(i int) string {
	return *s.Data_[i]
}

// Set the element at index i. The value v must be of type string or NullableString.
func (s Strings) Set(i int, v any) Series {
	if s.Partition_ != nil {
		return Errors{"Strings.Set: cannot set values on a grouped Series"}
	}

	switch v := v.(type) {
	case nil:
		s = s.MakeNullable().(Strings)
		s.Data_[i] = s.Ctx_.StringPool.Put(aargh.NA_TEXT)
		s.NullMask_[i>>3] |= 1 << uint(i%8)

	case string:
		s.Data_[i] = s.Ctx_.StringPool.Put(v)

	case aargh.NullableString:
		s = s.MakeNullable().(Strings)
		if v.Valid {
			s.Data_[i] = s.Ctx_.StringPool.Put(v.Value)
		} else {
			s.Data_[i] = s.Ctx_.StringPool.Put(aargh.NA_TEXT)
			s.NullMask_[i>>3] |= 1 << uint(i%8)
		}

	default:
		return Errors{fmt.Sprintf("Strings.Set: invalid type %T", v)}
	}

	s.Sorted_ = aargh.SORTED_NONE
	return s
}

////////////////////////			ALL DATA ACCESSORS

// Return the underlying Data_ as a slice of string.
func (s Strings) Strings() []string {
	Data_ := make([]string, len(s.Data_))
	if s.IsNullable_ {
		for i, v := range s.Data_ {
			if s.IsNull(i) {
				Data_[i] = aargh.NA_TEXT
			} else {
				Data_[i] = *v
			}
		}
	} else {
		for i, v := range s.Data_ {
			Data_[i] = *v
		}
	}
	return Data_
}

// Return the underlying Data_ as a slice of NullableString.
func (s Strings) DataAsNullable() any {
	Data_ := make([]aargh.NullableString, len(s.Data_))
	for i, v := range s.Data_ {
		Data_[i] = aargh.NullableString{Valid: !s.IsNull(i), Value: *v}
	}
	return Data_
}

// Return the underlying Data_ as a slice of string.
func (s Strings) DataAsString() []string {
	Data_ := make([]string, len(s.Data_))
	if s.IsNullable_ {
		for i, v := range s.Data_ {
			if s.IsNull(i) {
				Data_[i] = aargh.NA_TEXT
			} else {
				Data_[i] = *v
			}
		}
	} else {
		for i, v := range s.Data_ {
			Data_[i] = *v
		}
	}
	return Data_
}

func atoBool(s string) (bool, error) {
	trueRegex := regexp.MustCompile(`^[Tt]([Rr][Uu][Ee])?$`)
	falseRegex := regexp.MustCompile(`^[Ff]([Aa][Ll][Ss][Ee])?$`)

	if trueRegex.MatchString(s) {
		return true, nil
	} else if falseRegex.MatchString(s) {
		return false, nil
	}
	return false, fmt.Errorf("cannot convert \"%s\" to bool", s)
}

// Casts the series to a given type.
func (s Strings) Cast(t meta.BaseType) Series {
	switch t {
	case meta.BoolType:
		Data_ := make([]bool, len(s.Data_))
		NullMask_ := utils.BinVecInit(len(s.Data_), false)
		if s.IsNullable_ {
			copy(NullMask_, s.NullMask_)
		}

		if s.IsNullable_ {
			for i, v := range s.Data_ {
				if !s.IsNull(i) {
					b, err := atoBool(*v)
					if err != nil {
						NullMask_[i>>3] |= (1 << uint(i%8))
					}
					Data_[i] = b
				}
			}
		} else {
			for i, v := range s.Data_ {
				b, err := atoBool(*v)
				if err != nil {
					NullMask_[i>>3] |= (1 << uint(i%8))
				}
				Data_[i] = b
			}
		}

		return Bools{
			IsNullable_: true,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.IntType:
		Data_ := make([]int, len(s.Data_))
		NullMask_ := utils.BinVecInit(len(s.Data_), false)
		if s.IsNullable_ {
			copy(NullMask_, s.NullMask_)
		}

		if s.IsNullable_ {
			for i, v := range s.Data_ {
				if !s.IsNull(i) {
					d, err := strconv.Atoi(*v)
					if err != nil {
						NullMask_[i>>3] |= (1 << uint(i%8))
					} else {
						Data_[i] = int(d)
					}
				}
			}
		} else {
			for i, v := range s.Data_ {
				d, err := strconv.Atoi(*v)
				if err != nil {
					NullMask_[i>>3] |= (1 << uint(i%8))
				} else {
					Data_[i] = int(d)
				}
			}
		}

		return Ints{
			IsNullable_: true,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.Int64Type:
		Data_ := make([]int64, len(s.Data_))
		NullMask_ := utils.BinVecInit(len(s.Data_), false)
		if s.IsNullable_ {
			copy(NullMask_, s.NullMask_)
		}

		if s.IsNullable_ {
			for i, v := range s.Data_ {
				if !s.IsNull(i) {
					d, err := strconv.ParseInt(*v, 10, 64)
					if err != nil {
						NullMask_[i>>3] |= (1 << uint(i%8))
					} else {
						Data_[i] = d
					}
				}
			}
		} else {
			for i, v := range s.Data_ {
				d, err := strconv.ParseInt(*v, 10, 64)
				if err != nil {
					NullMask_[i>>3] |= (1 << uint(i%8))
				} else {
					Data_[i] = d
				}
			}
		}

		return Int64s{
			IsNullable_: true,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.Float64Type:
		Data_ := make([]float64, len(s.Data_))
		NullMask_ := utils.BinVecInit(len(s.Data_), false)
		if s.IsNullable_ {
			copy(NullMask_, s.NullMask_)
		}

		if s.IsNullable_ {
			for i, v := range s.Data_ {
				if !s.IsNull(i) {
					f, err := strconv.ParseFloat(*v, 64)
					if err != nil {
						NullMask_[i>>3] |= (1 << uint(i%8))
					} else {
						Data_[i] = f
					}
				}
			}
		} else {
			for i, v := range s.Data_ {
				f, err := strconv.ParseFloat(*v, 64)
				if err != nil {
					NullMask_[i>>3] |= (1 << uint(i%8))
				} else {
					Data_[i] = f
				}
			}
		}

		return Float64s{
			IsNullable_: true,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case meta.StringType:
		return s

	case meta.TimeType:
		return Errors{"Strings.Cast: cannot cast to Time, use Strings.ParseTime(layout) instead"}

	default:
		return Errors{fmt.Sprintf("Strings.Cast: invalid type %s", t.String())}
	}
}

// Parse the series as a time series.
func (s Strings) ParseTime(layout string) Series {
	Data_ := make([]time.Time, len(s.Data_))
	NullMask_ := utils.BinVecInit(len(s.Data_), false)
	if s.IsNullable_ {
		copy(NullMask_, s.NullMask_)
	}

	for i, v := range s.Data_ {
		if s.IsNullable_ && s.IsNull(i) {
			continue
		}
		t, err := time.Parse(layout, *v)
		if err != nil {
			NullMask_[i>>3] |= (1 << uint(i%8))
		} else {
			Data_[i] = t
		}
	}

	return Times{
		IsNullable_: true,
		Sorted_:     aargh.SORTED_NONE,
		Data_:       Data_,
		NullMask_:   NullMask_,
		Partition_:  nil,
		Ctx_:        s.Ctx_,
	}
}

////////////////////////			GROUPING OPERATIONS

// A SeriesStringPartition is a Partition_ of a Strings.
// Each key is a hash of a bool value, and each value is a slice of indices
// of the original series that are set to that value.
type SeriesStringPartition struct {
	Partition_ map[int64][]int
	Ctx_       *aargh.Context
}

func (gp *SeriesStringPartition) GetSize() int {
	return len(gp.Partition_)
}

func (gp *SeriesStringPartition) GetMap() map[int64][]int {
	return gp.Partition_
}

func (s Strings) Group() Series {

	// Define the worker callback
	worker := func(threadNum, start, end int, map_ map[int64][]int) {
		var ptr unsafe.Pointer
		for i := start; i < end; i++ {
			ptr = unsafe.Pointer(s.Data_[i])
			map_[(*(*int64)(unsafe.Pointer(&ptr)))] = append(map_[(*(*int64)(unsafe.Pointer(&ptr)))], i)
		}
	}

	// Define the worker callback for nulls
	workerNulls := func(threadNum, start, end int, map_ map[int64][]int, nulls *[]int) {
		var ptr unsafe.Pointer
		for i := start; i < end; i++ {
			if s.IsNull(i) {
				(*nulls) = append((*nulls), i)
			} else {
				ptr = unsafe.Pointer(s.Data_[i])
				map_[(*(*int64)(unsafe.Pointer(&ptr)))] = append(map_[(*(*int64)(unsafe.Pointer(&ptr)))], i)
			}
		}
	}

	Partition_ := SeriesStringPartition{
		Partition_: __series_groupby(
			aargh.THREADS_NUMBER, aargh.MINIMUM_PARALLEL_SIZE_1, len(s.Data_), s.HasNull(),
			worker, workerNulls),
		Ctx_: s.Ctx_,
	}

	s.Partition_ = &Partition_

	return s
}

func (s Strings) GroupBy(Partition_ SeriesPartition) Series {
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
		var ptr unsafe.Pointer
		for _, h := range keys[start:end] { // keys is defined outside the function
			for _, index := range otherIndeces[h] { // otherIndeces is defined outside the function
				ptr = unsafe.Pointer(s.Data_[index])
				newHash = *(*int64)(unsafe.Pointer(&ptr)) + aargh.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
				map_[newHash] = append(map_[newHash], index)
			}
		}
	}

	// Define the worker callback for nulls
	workerNulls := func(threadNum, start, end int, map_ map[int64][]int, nulls *[]int) {
		var newHash int64
		var ptr unsafe.Pointer
		for _, h := range keys[start:end] { // keys is defined outside the function
			for _, index := range otherIndeces[h] { // otherIndeces is defined outside the function
				if s.IsNull(index) {
					newHash = aargh.HASH_MAGIC_NUMBER_NULL + (h << 13) + (h >> 4)
				} else {
					ptr = unsafe.Pointer(s.Data_[index])
					newHash = *(*int64)(unsafe.Pointer(&ptr)) + aargh.HASH_MAGIC_NUMBER + (h << 13) + (h >> 4)
				}
				map_[newHash] = append(map_[newHash], index)
			}
		}
	}

	newPartition := SeriesStringPartition{
		Partition_: __series_groupby(
			aargh.THREADS_NUMBER, aargh.MINIMUM_PARALLEL_SIZE_1, len(keys), s.HasNull(),
			worker, workerNulls),
		Ctx_: s.Ctx_,
	}

	s.Partition_ = &newPartition

	return s
}

////////////////////////			SORTING OPERATIONS

func (s Strings) Less(i, j int) bool {
	if s.IsNullable_ {
		if s.NullMask_[i>>3]&(1<<uint(i%8)) > 0 {
			return false
		}
		if s.NullMask_[j>>3]&(1<<uint(j%8)) > 0 {
			return true
		}
	}

	return (*s.Data_[i]) < (*s.Data_[j])
}

func (s Strings) Equal(i, j int) bool {
	if s.IsNullable_ {
		if (s.NullMask_[i>>3] & (1 << uint(i%8))) > 0 {
			return (s.NullMask_[j>>3] & (1 << uint(j%8))) > 0
		}
		if (s.NullMask_[j>>3] & (1 << uint(j%8))) > 0 {
			return false
		}
	}

	return (*s.Data_[i]) == (*s.Data_[j])
}

func (s Strings) Swap(i, j int) {
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

func (s Strings) Sort() Series {
	if s.Sorted_ != aargh.SORTED_ASC {
		sort.Sort(s)
		s.Sorted_ = aargh.SORTED_ASC
	}
	return s
}

func (s Strings) SortRev() Series {
	if s.Sorted_ != aargh.SORTED_DESC {
		sort.Sort(sort.Reverse(s))
		s.Sorted_ = aargh.SORTED_DESC
	}
	return s
}

////////////////////////			STRING OPERATIONS

func (s Strings) ToUpper() Series {
	if s.Partition_ != nil {
		return Errors{"Strings.ToUpper() not supported on grouped Series"}
	}

	for i := 0; i < len(s.Data_); i++ {
		s.Data_[i] = s.Ctx_.StringPool.Put(strings.ToUpper(*s.Data_[i]))
	}

	return s
}

func (s Strings) ToLower() Series {
	if s.Partition_ != nil {
		return Errors{"Strings.ToLower() not supported on grouped Series"}
	}

	for i := 0; i < len(s.Data_); i++ {
		s.Data_[i] = s.Ctx_.StringPool.Put(strings.ToLower(*s.Data_[i]))
	}

	return s
}

func (s Strings) TrimSpace() Series {
	if s.Partition_ != nil {
		return Errors{"Strings.TrimSpace() not supported on grouped Series"}
	}

	for i := 0; i < len(s.Data_); i++ {
		s.Data_[i] = s.Ctx_.StringPool.Put(strings.TrimSpace(*s.Data_[i]))
	}

	return s
}

func (s Strings) Trim(cutset string) Series {
	if s.Partition_ != nil {
		return Errors{"Strings.Trim() not supported on grouped Series"}
	}

	for i := 0; i < len(s.Data_); i++ {
		s.Data_[i] = s.Ctx_.StringPool.Put(strings.Trim(*s.Data_[i], cutset))
	}

	return s
}

func (s Strings) Replace(old, new string, n int) Series {
	if s.Partition_ != nil {
		return Errors{"Strings.Replace() not supported on grouped Series"}
	}

	for i := 0; i < len(s.Data_); i++ {
		s.Data_[i] = s.Ctx_.StringPool.Put(strings.Replace(*s.Data_[i], old, new, n))
	}

	return s
}
