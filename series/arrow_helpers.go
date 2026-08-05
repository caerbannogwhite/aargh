package series

import (
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
)

// buildArrowFloat64 builds an Arrow Float64 array from a Go slice and aargh null mask.
func buildArrowFloat64(alloc memory.Allocator, data []float64, isNullable bool, nullMask []uint8) *array.Float64 {
	builder := array.NewFloat64Builder(alloc)
	defer builder.Release()
	builder.Reserve(len(data))
	if isNullable && len(nullMask) > 0 {
		for i, v := range data {
			if nullMask[i>>3]&(1<<uint(i%8)) != 0 {
				builder.AppendNull()
			} else {
				builder.Append(v)
			}
		}
	} else {
		builder.AppendValues(data, nil)
	}
	return builder.NewFloat64Array()
}

// buildArrowInt64 builds an Arrow Int64 array from a Go int64 slice and aargh null mask.
func buildArrowInt64(alloc memory.Allocator, data []int64, isNullable bool, nullMask []uint8) *array.Int64 {
	builder := array.NewInt64Builder(alloc)
	defer builder.Release()
	builder.Reserve(len(data))
	if isNullable && len(nullMask) > 0 {
		for i, v := range data {
			if nullMask[i>>3]&(1<<uint(i%8)) != 0 {
				builder.AppendNull()
			} else {
				builder.Append(v)
			}
		}
	} else {
		builder.AppendValues(data, nil)
	}
	return builder.NewInt64Array()
}

// buildArrowInt64FromInts builds an Arrow Int64 array from a Go int slice and aargh null mask.
// Arrow has no platform-dependent int type, so we store as int64.
func buildArrowInt64FromInts(alloc memory.Allocator, data []int, isNullable bool, nullMask []uint8) *array.Int64 {
	builder := array.NewInt64Builder(alloc)
	defer builder.Release()
	builder.Reserve(len(data))
	if isNullable && len(nullMask) > 0 {
		for i, v := range data {
			if nullMask[i>>3]&(1<<uint(i%8)) != 0 {
				builder.AppendNull()
			} else {
				builder.Append(int64(v))
			}
		}
	} else {
		for _, v := range data {
			builder.Append(int64(v))
		}
	}
	return builder.NewInt64Array()
}

// buildArrowBoolean builds an Arrow Boolean array from a Go bool slice and aargh null mask.
func buildArrowBoolean(alloc memory.Allocator, data []bool, isNullable bool, nullMask []uint8) *array.Boolean {
	builder := array.NewBooleanBuilder(alloc)
	defer builder.Release()
	builder.Reserve(len(data))
	if isNullable && len(nullMask) > 0 {
		for i, v := range data {
			if nullMask[i>>3]&(1<<uint(i%8)) != 0 {
				builder.AppendNull()
			} else {
				builder.Append(v)
			}
		}
	} else {
		builder.AppendValues(data, nil)
	}
	return builder.NewBooleanArray()
}

// buildArrowString builds an Arrow String array from a Go []*string slice and aargh null mask.
func buildArrowString(alloc memory.Allocator, data []*string, isNullable bool, nullMask []uint8) *array.String {
	builder := array.NewStringBuilder(alloc)
	defer builder.Release()
	builder.Reserve(len(data))
	if isNullable && len(nullMask) > 0 {
		for i, v := range data {
			if nullMask[i>>3]&(1<<uint(i%8)) != 0 {
				builder.AppendNull()
			} else {
				builder.Append(*v)
			}
		}
	} else {
		for _, v := range data {
			builder.Append(*v)
		}
	}
	return builder.NewStringArray()
}

// buildArrowTimestamp builds an Arrow Timestamp array from a Go time.Time slice and aargh null mask.
// Arrow timestamps are stored as int64 nanoseconds since epoch.
func buildArrowTimestamp(alloc memory.Allocator, data []time.Time, isNullable bool, nullMask []uint8) *array.Timestamp {
	builder := array.NewTimestampBuilder(alloc, &arrow.TimestampType{Unit: arrow.Nanosecond})
	defer builder.Release()
	builder.Reserve(len(data))
	if isNullable && len(nullMask) > 0 {
		for i, v := range data {
			if nullMask[i>>3]&(1<<uint(i%8)) != 0 {
				builder.AppendNull()
			} else {
				builder.Append(arrow.Timestamp(v.UnixNano()))
			}
		}
	} else {
		for _, v := range data {
			builder.Append(arrow.Timestamp(v.UnixNano()))
		}
	}
	return builder.NewTimestampArray()
}

// buildArrowDuration builds an Arrow Duration array from a Go time.Duration slice and aargh null mask.
// Arrow durations are stored as int64 nanoseconds.
func buildArrowDuration(alloc memory.Allocator, data []time.Duration, isNullable bool, nullMask []uint8) *array.Duration {
	builder := array.NewDurationBuilder(alloc, &arrow.DurationType{Unit: arrow.Nanosecond})
	defer builder.Release()
	builder.Reserve(len(data))
	if isNullable && len(nullMask) > 0 {
		for i, v := range data {
			if nullMask[i>>3]&(1<<uint(i%8)) != 0 {
				builder.AppendNull()
			} else {
				builder.Append(arrow.Duration(v.Nanoseconds()))
			}
		}
	} else {
		for _, v := range data {
			builder.Append(arrow.Duration(v.Nanoseconds()))
		}
	}
	return builder.NewDurationArray()
}
