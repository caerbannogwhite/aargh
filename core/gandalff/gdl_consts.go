package gandalff

import "strconv"

const (
	// The default capacity of a series.
	DEFAULT_SERIES_INITIAL_CAPACITY = 10

	// The default capacity of a hash map.
	DEFAULT_HASH_MAP_INITIAL_CAPACITY = 1024

	// The default capacity of a dense map array.
	DEFAULT_DENSE_MAP_ARRAY_INITIAL_CAPACITY = 64

	// Number of threads to use for parallel operations.
	THREADS_NUMBER = 16

	// Minimum number of elements to use parallel operations.
	MINIMUM_PARALLEL_SIZE_1 = 16_384
	MINIMUM_PARALLEL_SIZE_2 = 131_072

	HASH_MAGIC_NUMBER = int64(0xa8f4979b77e3f93)
	HASH_NULL_KEY     = int64(0x7ff8000000000001)
)

////////////////////////////////			ENUMS

type SeriesSortOrder int16

const (
	// The series is not sorted.
	SORTED_NONE SeriesSortOrder = iota
	// The series is sorted in ascending order.
	SORTED_ASC
	// The series is sorted in descending order.
	SORTED_DESC
)

////////////////////////////////			ERRORS

////////////////////////////////			TO STRING

const NULL_STRING = "NA"
const BOOL_TRUE_STRING = "true"
const BOOL_FALSE_STRING = "false"

func boolToString(b bool) string {
	if b {
		return BOOL_TRUE_STRING
	} else {
		return BOOL_FALSE_STRING
	}
}

func intToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func floatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

type any interface{}

////////////////////////////////			NULLABLE TYPES

type NullableBool struct {
	Valid bool
	Value bool
}

type NullableInt8 struct {
	Valid bool
	Value int8
}

type NullableInt16 struct {
	Valid bool
	Value int16
}

type NullableInt32 struct {
	Valid bool
	Value int32
}

type NullableInt64 struct {
	Valid bool
	Value int64
}

type NullableFloat32 struct {
	Valid bool
	Value float32
}

type NullableFloat64 struct {
	Valid bool
	Value float64
}

type NullableString struct {
	Valid bool
	Value string
}

type GDLMapFunc func(any) any
