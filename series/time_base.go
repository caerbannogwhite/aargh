package series

import (
	"fmt"
	"time"

	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/meta"
	"github.com/caerbannogwhite/aargh/utils"
)

func (s Times) printInfo() {
	fmt.Println("Times")
	fmt.Println("==========")
	fmt.Println("IsNullable:", s.IsNullable_)
	fmt.Println("Sorted:    ", s.Sorted_)
	fmt.Println("Data:      ", s.Data_)
	fmt.Println("NullMask:  ", s.NullMask_)
	fmt.Println("Partition: ", s.Partition_)
	fmt.Println("Context:   ", s.Ctx_)
}

////////////////////////			BASIC ACCESSORS

// Return the context of the series.
func (s Times) GetContext() *aargh.Context {
	return s.Ctx_
}

// Return the number of elements in the series.
func (s Times) Len() int {
	return len(s.Data_)
}

// Return the type of the series.
func (s Times) Type() meta.BaseType {
	return meta.TimeType
}

// Return the type and cardinality of the series.
func (s Times) TypeCard() meta.BaseTypeCard {
	return meta.BaseTypeCard{Base: meta.TimeType, Card: s.Len()}
}

// Return if the series is grouped.
func (s Times) IsGrouped() bool {
	return s.Partition_ != nil
}

// Return if the series admits null values.
func (s Times) IsNullable() bool {
	return s.IsNullable_
}

// Return if the series is Sorted_.
func (s Times) IsSorted() aargh.SeriesSortOrder {
	return s.Sorted_
}

// Return if the series is error.
func (s Times) IsError() bool {
	return false
}

// Return the error message of the series.
func (s Times) GetError() string {
	return ""
}

// Return the Partition_ of the series.
func (s Times) GetPartition() SeriesPartition {
	return s.Partition_
}

// Return if the series has null values.
func (s Times) HasNull() bool {
	for _, v := range s.NullMask_ {
		if v != 0 {
			return true
		}
	}
	return false
}

// Return the number of null values in the series.
func (s Times) NullCount() int {
	count := 0
	for _, x := range s.NullMask_ {
		for ; x != 0; x >>= 1 {
			count += int(x & 1)
		}
	}
	return count
}

// Return if the element at index i is null.
func (s Times) IsNull(i int) bool {
	if s.IsNullable_ {
		return s.NullMask_[i>>3]&(1<<uint(i%8)) != 0
	}
	return false
}

// Return the null mask of the series.
func (s Times) GetNullMask() []bool {
	mask := make([]bool, len(s.Data_))
	idx := 0
	for _, v := range s.NullMask_ {
		for i := 0; i < 8 && idx < len(s.Data_); i++ {
			mask[idx] = v&(1<<uint(i)) != 0
			idx++
		}
	}
	return mask
}

// Set the null mask of the series.
func (s Times) SetNullMask(mask []bool) Series {
	if s.Partition_ != nil {
		return Errors{"Times.SetNullMask: cannot set values on a grouped series"}
	}

	if s.IsNullable_ {
		for k, v := range mask {
			if v {
				s.NullMask_[k>>3] |= 1 << uint(k%8)
			} else {
				s.NullMask_[k>>3] &= ^(1 << uint(k%8))
			}
		}
		return s
	} else {
		NullMask_ := utils.BinVecInit(len(s.Data_), false)
		for k, v := range mask {
			if v {
				NullMask_[k>>3] |= 1 << uint(k%8)
			} else {
				NullMask_[k>>3] &= ^(1 << uint(k%8))
			}
		}

		s.IsNullable_ = true
		s.NullMask_ = NullMask_

		return s
	}
}

// Make the series nullable.
func (s Times) MakeNullable() Series {
	if !s.IsNullable_ {
		s.IsNullable_ = true
		s.NullMask_ = utils.BinVecInit(len(s.Data_), false)
	}
	return s
}

// Make the series non-nullable.
func (s Times) MakeNonNullable() Series {
	if s.IsNullable_ {
		s.IsNullable_ = false
		s.NullMask_ = make([]uint8, 0)
	}
	return s
}

// Get the element at index i.
func (s Times) Get(i int) any {
	return s.Data_[i]
}

// Append appends a value or a slice of values to the series.
func (s Times) Append(v any) Series {
	if s.Partition_ != nil {
		return Errors{"Times.Append: cannot append values to a grouped series"}
	}

	switch v := v.(type) {
	case nil:
		s.Data_ = append(s.Data_, time.Time{})
		s = s.MakeNullable().(Times)
		if len(s.Data_) > len(s.NullMask_)<<3 {
			s.NullMask_ = append(s.NullMask_, 0)
		}
		s.NullMask_[(len(s.Data_)-1)>>3] |= 1 << uint8((len(s.Data_)-1)%8)

	case NAs:
		s.IsNullable_, s.NullMask_ = utils.MergeNullMasks(len(s.Data_), s.IsNullable_, s.NullMask_, v.Len(), true, utils.BinVecInit(v.Len(), true))
		s.Data_ = append(s.Data_, make([]time.Time, v.Len())...)

	case time.Time:
		s.Data_ = append(s.Data_, v)
		if s.IsNullable_ && len(s.Data_) > len(s.NullMask_)<<3 {
			s.NullMask_ = append(s.NullMask_, 0)
		}

	case []time.Time:
		s.Data_ = append(s.Data_, v...)
		if s.IsNullable_ && len(s.Data_) > len(s.NullMask_)<<3 {
			s.NullMask_ = append(s.NullMask_, make([]uint8, (len(s.Data_)>>3)-len(s.NullMask_))...)
		}

	case aargh.NullableTime:
		s.Data_ = append(s.Data_, v.Value)
		s = s.MakeNullable().(Times)
		if len(s.Data_) > len(s.NullMask_)<<3 {
			s.NullMask_ = append(s.NullMask_, 0)
		}
		if !v.Valid {
			s.NullMask_[(len(s.Data_)-1)>>3] |= 1 << uint8((len(s.Data_)-1)%8)
		}

	case []aargh.NullableTime:
		ssize := len(s.Data_)
		s.Data_ = append(s.Data_, make([]time.Time, len(v))...)
		s = s.MakeNullable().(Times)
		if len(s.Data_) > len(s.NullMask_)<<3 {
			s.NullMask_ = append(s.NullMask_, make([]uint8, (len(s.Data_)>>3)-len(s.NullMask_)+1)...)
		}
		for i, b := range v {
			s.Data_[ssize+i] = b.Value
			if !b.Valid {
				s.NullMask_[(ssize+i)>>3] |= 1 << uint8((ssize+i)%8)
			}
		}

	case Times:
		if s.Ctx_ != v.Ctx_ {
			return Errors{"Times.Append: cannot append Times from different contexts"}
		}

		s.IsNullable_, s.NullMask_ = utils.MergeNullMasks(len(s.Data_), s.IsNullable_, s.NullMask_, len(v.Data_), v.IsNullable_, v.NullMask_)
		s.Data_ = append(s.Data_, v.Data_...)

	default:
		return Errors{fmt.Sprintf("Times.Append: invalid type %T", v)}
	}

	s.Sorted_ = aargh.SORTED_NONE
	return s
}

// Take the elements according to the given interval.
func (s Times) Take(params ...int) Series {
	indeces, err := SeriesTakePreprocess("Times", s.Len(), params...)
	if err != nil {
		return Errors{err.Error()}
	}
	return s.FilterIntSlice(indeces, false)
}

// Return the elements of the series as a slice.
func (s Times) Data() any {
	return s.Data_
}

// Copy the series.
func (s Times) Copy() Series {
	Data_ := make([]time.Time, len(s.Data_))
	copy(Data_, s.Data_)
	NullMask_ := make([]uint8, len(s.NullMask_))
	copy(NullMask_, s.NullMask_)

	return Times{
		IsNullable_: s.IsNullable_,
		Sorted_:     s.Sorted_,
		Data_:       Data_,
		NullMask_:   NullMask_,
		Partition_:  s.Partition_,
		Ctx_:        s.Ctx_,
	}
}

func (s Times) GetData() []time.Time {
	return s.Data_
}

// Ungroup the series.
func (s Times) UnGroup() Series {
	s.Partition_ = nil
	return s
}

////////////////////////			FILTER OPERATIONS

// Filters out the elements by the given mask.
// Mask can be Bools, Ints, bool slice or a int slice.
func (s Times) Filter(mask any) Series {
	switch mask := mask.(type) {
	case Bools:
		return s.filterBoolSlice(mask.Data_)
	case Ints:
		return s.FilterIntSlice(mask.Data_, true)
	case []bool:
		return s.filterBoolSlice(mask)
	case []int:
		return s.FilterIntSlice(mask, true)
	default:
		return Errors{fmt.Sprintf("Times.Filter: invalid type %T", mask)}
	}
}

func (s Times) filterBoolSlice(mask []bool) Series {
	if len(mask) != len(s.Data_) {
		return Errors{fmt.Sprintf("Times.Filter: mask length (%d) does not match series length (%d)", len(mask), len(s.Data_))}
	}

	elementCount := 0
	for _, v := range mask {
		if v {
			elementCount++
		}
	}

	var Data_ []time.Time
	var NullMask_ []uint8

	Data_ = make([]time.Time, elementCount)

	if s.IsNullable_ {
		NullMask_ = utils.BinVecInit(elementCount, false)
		dstIdx := 0
		for srcIdx, v := range mask {
			if v {
				Data_[dstIdx] = s.Data_[srcIdx]
				if srcIdx%8 > dstIdx%8 {
					NullMask_[dstIdx>>3] |= ((s.NullMask_[srcIdx>>3] & (1 << uint(srcIdx%8))) >> uint(srcIdx%8-dstIdx%8))
				} else {
					NullMask_[dstIdx>>3] |= ((s.NullMask_[srcIdx>>3] & (1 << uint(srcIdx%8))) << uint(dstIdx%8-srcIdx%8))
				}
				dstIdx++
			}
		}
	} else {
		NullMask_ = make([]uint8, 0)
		dstIdx := 0
		for srcIdx, v := range mask {
			if v {
				Data_[dstIdx] = s.Data_[srcIdx]
				dstIdx++
			}
		}
	}

	s.Data_ = Data_
	s.NullMask_ = NullMask_

	return s
}

func (s Times) FilterIntSlice(indexes []int, check bool) Series {
	if len(indexes) == 0 {
		s.Data_ = make([]time.Time, 0)
		s.NullMask_ = make([]uint8, 0)
		return s
	}

	// check if indexes are in range
	if check {
		for _, v := range indexes {
			if v < 0 || v >= len(s.Data_) {
				return Errors{fmt.Sprintf("Times.Filter: index %d is out of range", v)}
			}
		}
	}

	var Data_ []time.Time
	var NullMask_ []uint8

	size := len(indexes)
	Data_ = make([]time.Time, size)

	if s.IsNullable_ {
		NullMask_ = utils.BinVecInit(size, false)
		for dstIdx, srcIdx := range indexes {
			Data_[dstIdx] = s.Data_[srcIdx]
			if srcIdx%8 > dstIdx%8 {
				NullMask_[dstIdx>>3] |= ((s.NullMask_[srcIdx>>3] & (1 << uint(srcIdx%8))) >> uint(srcIdx%8-dstIdx%8))
			} else {
				NullMask_[dstIdx>>3] |= ((s.NullMask_[srcIdx>>3] & (1 << uint(srcIdx%8))) << uint(dstIdx%8-srcIdx%8))
			}
		}
	} else {
		NullMask_ = make([]uint8, 0)
		for dstIdx, srcIdx := range indexes {
			Data_[dstIdx] = s.Data_[srcIdx]
		}
	}

	s.Data_ = Data_
	s.NullMask_ = NullMask_

	return s
}

// Apply the given function to each element of the series.
func (s Times) Map(f aargh.MapFunc) Series {
	if len(s.Data_) == 0 {
		return s
	}

	v := f(s.Get(0))
	switch v.(type) {
	case bool:
		Data_ := make([]bool, len(s.Data_))
		for i := 0; i < len(s.Data_); i++ {
			Data_[i] = f(s.Data_[i]).(bool)
		}

		return Bools{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case int:
		Data_ := make([]int, len(s.Data_))
		for i := 0; i < len(s.Data_); i++ {
			Data_[i] = f(s.Data_[i]).(int)
		}

		return Ints{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case int64:
		Data_ := make([]int64, len(s.Data_))
		for i := 0; i < len(s.Data_); i++ {
			Data_[i] = f(s.Data_[i]).(int64)
		}

		return Int64s{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case float64:
		Data_ := make([]float64, len(s.Data_))
		for i := 0; i < len(s.Data_); i++ {
			Data_[i] = f(s.Data_[i]).(float64)
		}

		return Float64s{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case string:
		Data_ := make([]*string, len(s.Data_))
		for i := 0; i < len(s.Data_); i++ {
			Data_[i] = s.Ctx_.StringPool.Put(f(s.Data_[i]).(string))
		}

		return Strings{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case time.Time:
		Data_ := make([]time.Time, len(s.Data_))
		for i := 0; i < len(s.Data_); i++ {
			Data_[i] = f(s.Data_[i]).(time.Time)
		}

		return Times{
			IsNullable_: s.IsNullable_,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   s.NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case time.Duration:
		Data_ := make([]time.Duration, len(s.Data_))
		for i := 0; i < len(s.Data_); i++ {
			Data_[i] = f(s.Data_[i]).(time.Duration)
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
		return Errors{fmt.Sprintf("Times.Map: Unsupported type %T", v)}
	}
}

// Apply the given function to each element of the series.
func (s Times) MapNull(f aargh.MapFuncNull) Series {
	if len(s.Data_) == 0 {
		return s
	}

	if !s.IsNullable_ {
		return Errors{"Times.MapNull: series is not nullable"}
	}

	v, isNull := f(s.Get(0), s.IsNull(0))
	switch v.(type) {
	case bool:
		Data_ := make([]bool, len(s.Data_))
		NullMask_ := make([]uint8, len(s.NullMask_))
		for i := 0; i < len(s.Data_); i++ {
			v, isNull = f(s.Data_[i], s.IsNull(i))
			Data_[i] = v.(bool)
			if isNull {
				NullMask_[i>>3] |= 1 << uint(i%8)
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

	case int:
		Data_ := make([]int, len(s.Data_))
		NullMask_ := make([]uint8, len(s.NullMask_))
		for i := 0; i < len(s.Data_); i++ {
			v, isNull = f(s.Data_[i], s.IsNull(i))
			Data_[i] = v.(int)
			if isNull {
				NullMask_[i>>3] |= 1 << uint(i%8)
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

	case int64:
		Data_ := make([]int64, len(s.Data_))
		NullMask_ := make([]uint8, len(s.NullMask_))
		for i := 0; i < len(s.Data_); i++ {
			v, isNull = f(s.Data_[i], s.IsNull(i))
			Data_[i] = v.(int64)
			if isNull {
				NullMask_[i>>3] |= 1 << uint(i%8)
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

	case float64:
		Data_ := make([]float64, len(s.Data_))
		NullMask_ := make([]uint8, len(s.NullMask_))
		for i := 0; i < len(s.Data_); i++ {
			v, isNull = f(s.Data_[i], s.IsNull(i))
			Data_[i] = v.(float64)
			if isNull {
				NullMask_[i>>3] |= 1 << uint(i%8)
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

	case string:
		Data_ := make([]*string, len(s.Data_))
		NullMask_ := make([]uint8, len(s.NullMask_))
		for i := 0; i < len(s.Data_); i++ {
			v, isNull = f(s.Data_[i], s.IsNull(i))
			Data_[i] = s.Ctx_.StringPool.Put(v.(string))
			if isNull {
				NullMask_[i>>3] |= 1 << uint(i%8)
			}
		}

		return Strings{
			IsNullable_: true,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	case time.Time:
		Data_ := make([]time.Time, len(s.Data_))
		NullMask_ := make([]uint8, len(s.NullMask_))
		for i := 0; i < len(s.Data_); i++ {
			v, isNull = f(s.Data_[i], s.IsNull(i))
			Data_[i] = v.(time.Time)
			if isNull {
				NullMask_[i>>3] |= 1 << uint(i%8)
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

	case time.Duration:
		Data_ := make([]time.Duration, len(s.Data_))
		NullMask_ := make([]uint8, len(s.NullMask_))
		for i := 0; i < len(s.Data_); i++ {
			v, isNull = f(s.Data_[i], s.IsNull(i))
			Data_[i] = v.(time.Duration)
			if isNull {
				NullMask_[i>>3] |= 1 << uint(i%8)
			}
		}

		return Durations{
			IsNullable_: true,
			Sorted_:     aargh.SORTED_NONE,
			Data_:       Data_,
			NullMask_:   NullMask_,
			Partition_:  nil,
			Ctx_:        s.Ctx_,
		}

	default:
		return Errors{fmt.Sprintf("Times.MapNull: Unsupported type %T", v)}
	}
}
