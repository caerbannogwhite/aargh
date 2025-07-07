package io

import (
	"fmt"
	"testing"
	"time"

	"github.com/caerbannogwhite/aargh/meta"
)

func TestTypeGuesser(t *testing.T) {
	// Create a new type guesser.
	tg := newTypeGuesser(false)

	// Test the bool type.
	if tg.guessType("true") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("true").String())
	}

	if tg.guessType("false") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("false").String())
	}

	if tg.guessType("True") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("True").String())
	}

	if tg.guessType("False") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("False").String())
	}

	if tg.guessType("TRUE") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("TRUE").String())
	}

	if tg.guessType("FALSE") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("FALSE").String())
	}

	if tg.guessType("t") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("t").String())
	}

	if tg.guessType("f") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("f").String())
	}

	if tg.guessType("T") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("T").String())
	}

	if tg.guessType("F") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("F").String())
	}

	if tg.guessType("TrUe") != meta.BoolType {
		t.Error("Expected Bool, got", tg.guessType("TrUe").String())
	}

	// Real case: do not remove this test
	if tg.guessType("TLS") != meta.StringType {
		t.Error("Expected String, got", tg.guessType("TLS").String())
	}

	// Test the int type.
	if tg.guessType("0") != meta.Int64Type {
		t.Error("Expected Int64, got", tg.guessType("0").String())
	}

	if tg.guessType("1") != meta.Int64Type {
		t.Error("Expected Int64, got", tg.guessType("1").String())
	}

	if tg.guessType("10000") != meta.Int64Type {
		t.Error("Expected Int64, got", tg.guessType("10000").String())
	}

	if tg.guessType("-1") != meta.Int64Type {
		t.Error("Expected Int64, got", tg.guessType("-1").String())
	}

	if tg.guessType("-10000") != meta.Int64Type {
		t.Error("Expected Int64, got", tg.guessType("-10000").String())
	}

	// Test the float type.
	if tg.guessType("0.0") != meta.Float64Type {
		t.Error("Expected Float64, got", tg.guessType("0.0").String())
	}

	if tg.guessType("1.0") != meta.Float64Type {
		t.Error("Expected Float64, got", tg.guessType("1.0").String())
	}

	if tg.guessType("10000.0") != meta.Float64Type {
		t.Error("Expected Float64, got", tg.guessType("10000.0").String())
	}

	if tg.guessType("-1.0") != meta.Float64Type {
		t.Error("Expected Float64, got", tg.guessType("-1.0").String())
	}

	if tg.guessType("-1e3") != meta.Float64Type {
		t.Error("Expected Float64, got", tg.guessType("-1e3").String())
	}

	if tg.guessType("-1e-3") != meta.Float64Type {
		t.Error("Expected Float64, got", tg.guessType("-1e-3").String())
	}

	if tg.guessType("2.0E4") != meta.Float64Type {
		t.Error("Expected Float64, got", tg.guessType("2.0E4").String())
	}

	if tg.guessType("2.0e4") != meta.Float64Type {
		t.Error("Expected Float64, got", tg.guessType("2.0e4").String())
	}
}

func TestTypeGuesserWithNa(t *testing.T) {
	// Create a new type guesser.
	tg := newTypeGuesser(true)

	tg.setLength(4)

	tg.guessTypes([]string{"t", "-1", "-1e-1", "a"})
	tg.guessTypes([]string{"f", "0", "1E+1", "b"})
	tg.guessTypes([]string{"", "", "", ""})
	tg.guessTypes([]string{"true", "1", "1.23e2", "c"})
	tg.guessTypes([]string{"false", "2", "1.23e-2", "d"})
	tg.guessTypes([]string{"na", "null", "n/a", "e"})

	if tg.typeBuckets[0].boolCount != 4 && tg.typeBuckets[0].nullCount != 2 {
		t.Error("Expected 4 bools and 2 nulls, got", tg.typeBuckets[0].boolCount, tg.typeBuckets[0].nullCount)
	}
	if tg.typeBuckets[1].intCount != 4 && tg.typeBuckets[1].nullCount != 2 {
		t.Error("Expected 4 ints and 2 nulls, got", tg.typeBuckets[1].intCount, tg.typeBuckets[1].nullCount)
	}
	if tg.typeBuckets[2].floatCount != 4 && tg.typeBuckets[2].nullCount != 2 {
		t.Error("Expected 4 floats and 2 nulls, got", tg.typeBuckets[2].floatCount, tg.typeBuckets[2].nullCount)
	}
}

func TestTypeGuesserBasicTypes(t *testing.T) {
	tg := newTypeGuesser(false)
	tg.setLength(1)

	// Test basic type detection
	testCases := []struct {
		input    string
		expected meta.BaseType
	}{
		{"true", meta.BoolType},
		{"false", meta.BoolType},
		{"TRUE", meta.BoolType},
		{"FALSE", meta.BoolType},
		{"123", meta.Int64Type},
		{"-456", meta.Int64Type},
		{"+789", meta.Int64Type},
		{"3.14", meta.Float64Type},
		{"-2.718", meta.Float64Type},
		{"1.23e-10", meta.Float64Type},
		{"hello", meta.StringType},
		{"", meta.StringType},
	}

	for _, tc := range testCases {
		result := tg.guessType(tc.input)
		if result != tc.expected {
			t.Errorf("guessType(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestTypeGuesserDateTypes(t *testing.T) {
	tg := newTypeGuesser(false)
	tg.setLength(1)

	// Test date and datetime detection
	testCases := []struct {
		input    string
		expected meta.BaseType
	}{
		// Date formats
		{"2023-01-15", meta.TimeType},
		{"2023/01/15", meta.TimeType},
		{"01/15/2023", meta.TimeType},
		{"15/01/2023", meta.TimeType},
		{"01-15-2023", meta.TimeType},
		{"15-01-2023", meta.TimeType},
		{"01/15/23", meta.TimeType},
		{"15/01/23", meta.TimeType},

		// Datetime formats
		{"2023-01-15 14:30:00", meta.TimeType},
		{"2023-01-15 14:30", meta.TimeType},
		{"2023/01/15 14:30:00", meta.TimeType},
		{"01/15/2023 14:30:00", meta.TimeType},
		{"15/01/2023 14:30:00", meta.TimeType},
		{"2023-01-15 02:30:00 PM", meta.TimeType},
		{"2023-01-15 02:30 PM", meta.TimeType},
		{"01/15/2023 02:30:00 PM", meta.TimeType},
		{"15/01/2023 02:30:00 PM", meta.TimeType},

		// Non-date strings should be StringType
		{"2023-13-45", meta.StringType}, // Invalid date
		{"25:70:90", meta.StringType},   // Invalid time
		{"not a date", meta.StringType},
	}

	for _, tc := range testCases {
		result := tg.guessType(tc.input)
		if result != tc.expected {
			t.Errorf("guessType(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestTypeGuesserWithNulls(t *testing.T) {
	tg := newTypeGuesser(true)
	tg.setLength(1)

	// Test null value detection
	testCases := []struct {
		input    string
		expected meta.BaseType
	}{
		{"null", meta.StringType}, // null values are handled separately
		{"NULL", meta.StringType},
		{"nan", meta.StringType},
		{"NAN", meta.StringType},
		{"", meta.StringType},
		{"true", meta.BoolType},
		{"123", meta.Int64Type},
	}

	for _, tc := range testCases {
		result := tg.guessType(tc.input)
		if result != tc.expected {
			t.Errorf("guessType(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestTypeBucketGetMostCommonType(t *testing.T) {
	testCases := []struct {
		bucket          typeBucket
		expected        meta.BaseType
		nullsCount      int
		remainingsCount int
	}{
		{
			typeBucket{boolCount: 5, intCount: 2, stringCount: 1, totalCount: 8},
			meta.BoolType,
			0,
			3,
		},
		{
			typeBucket{boolCount: 10, intCount: 0, floatCount: 0, stringCount: 0, dateCount: 0, datetimeCount: 0, totalCount: 10},
			meta.BoolType,
			0,
			0,
		},
		{
			typeBucket{intCount: 8, floatCount: 3, stringCount: 1, totalCount: 12},
			meta.Int64Type,
			0,
			4,
		},
		{
			typeBucket{floatCount: 6, stringCount: 2, totalCount: 8},
			meta.Float64Type,
			0,
			2,
		},
		{
			typeBucket{dateCount: 7, stringCount: 3, datetimeCount: 1, totalCount: 11},
			meta.TimeType,
			0,
			3,
		},
		{
			typeBucket{datetimeCount: 9, stringCount: 2, totalCount: 11},
			meta.TimeType,
			0,
			2,
		},
		{
			typeBucket{stringCount: 10, boolCount: 1, intCount: 1, totalCount: 12},
			meta.StringType,
			0,
			2,
		},
		{
			typeBucket{stringCount: 5, totalCount: 5},
			meta.StringType,
			0,
			0,
		},
	}

	for i, tc := range testCases {
		result, nullsCount, remainingsCount, _ := tc.bucket.getMostCommonType()
		if result != tc.expected {
			t.Errorf("Test case %d: getMostCommonType() = %v, expected %v", i, result, tc.expected)
		}
		if nullsCount != tc.nullsCount {
			t.Errorf("Test case %d: nullsCount = %v, expected %v", i, nullsCount, tc.nullsCount)
		}
		if remainingsCount != tc.remainingsCount {
			t.Errorf("Test case %d: remainingsCount = %v, expected %v", i, remainingsCount, tc.remainingsCount)
		}
	}
}

func TestTypeGuesserGuessTypes(t *testing.T) {
	tg := newTypeGuesser(false)
	tg.setLength(3)

	records := [][]string{
		{"true", "123", "2023-01-15"},
		{"false", "456", "2023-02-20"},
		{"true", "789", "2023-03-25"},
	}

	for _, record := range records {
		tg.guessTypes(record)
	}

	types := tg.getSchema().GetDataTypes()
	expected := []meta.BaseType{meta.BoolType, meta.Int64Type, meta.TimeType}

	if len(types) != len(expected) {
		t.Errorf("Expected %d types, got %d", len(expected), len(types))
		return
	}

	for i, expectedType := range expected {
		if types[i] != expectedType {
			t.Errorf("Column %d: expected %v, got %v", i, expectedType, types[i])
		}
	}
}

func TestTypeGuesserGuessTypesNulls(t *testing.T) {
	tg := newTypeGuesser(true)
	tg.setLength(3)

	records := [][]string{
		{"true", "123", "2023-01-15"},
		{"null", "456", "2023-02-20"},
		{"false", "null", "2023-03-25"},
	}

	for _, record := range records {
		tg.guessTypes(record)
	}

	types := tg.getSchema().GetDataTypes()
	expected := []meta.BaseType{meta.BoolType, meta.Int64Type, meta.TimeType}

	if len(types) != len(expected) {
		t.Errorf("Expected %d types, got %d", len(expected), len(types))
		return
	}

	for i, expectedType := range expected {
		if types[i] != expectedType {
			t.Errorf("Column %d: expected %v, got %v", i, expectedType, types[i])
		}
	}
}

func TestAtoTime(t *testing.T) {
	tg := newTypeGuesser(false)

	testCases := []struct {
		input    string
		expected time.Time
		hasError bool
	}{
		{"2023-01-15", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"2023/01/15", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"01/15/2023", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"15/01/2023", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), false},
		{"2023-01-15 14:30:00", time.Date(2023, 1, 15, 14, 30, 0, 0, time.UTC), false},
		{"2023-01-15 02:30:00 PM", time.Date(2023, 1, 15, 14, 30, 0, 0, time.UTC), false},
		{"2023-01-15 02:30 PM", time.Date(2023, 1, 15, 14, 30, 0, 0, time.UTC), false},
		{"invalid date", time.Time{}, true},
		{"2023-13-45", time.Time{}, true},
		{"25:70:90", time.Time{}, true},
	}

	for _, tc := range testCases {
		result, err := tg.atoTime(tc.input)
		if tc.hasError {
			if err == nil {
				t.Errorf("atoTime(%q) should have returned an error", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("atoTime(%q) returned error: %v", tc.input, err)
			} else if !result.Equal(tc.expected) {
				t.Errorf("atoTime(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		}
	}
}

func TestAtoBool(t *testing.T) {
	tg := newTypeGuesser(false)

	testCases := []struct {
		input    string
		expected bool
		hasError bool
	}{
		{"true", true, false},
		{"TRUE", true, false},
		{"t", true, false},
		{"T", true, false},
		{"false", false, false},
		{"FALSE", false, false},
		{"f", false, false},
		{"F", false, false},
		{"invalid", false, true},
		{"1", false, true},
		{"0", false, true},
	}

	for _, tc := range testCases {
		result, err := tg.atoBool(tc.input)
		if tc.hasError {
			if err == nil {
				t.Errorf("atoBool(%q) should have returned an error", tc.input)
			}
		} else {
			if err != nil {
				t.Errorf("atoBool(%q) returned error: %v", tc.input, err)
			} else if result != tc.expected {
				t.Errorf("atoBool(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		}
	}
}

func TestTypeGuesserMixedData(t *testing.T) {
	tg := newTypeGuesser(false)
	tg.setLength(4)

	// Test with mixed data types including dates
	records := [][]string{
		{"true", "123", "3.14", "2023-01-15"},
		{"false", "456", "2.718", "2023-02-20 14:30:00"},
		{"true", "789", "1.414", "01/15/2023"},
		{"false", "101", "0.577", "15/01/2023 02:30 PM"},
	}

	for _, record := range records {
		tg.guessTypes(record)
	}

	types := tg.getSchema().GetDataTypes()
	expected := []meta.BaseType{meta.BoolType, meta.Int64Type, meta.Float64Type, meta.TimeType}

	if len(types) != len(expected) {
		t.Errorf("Expected %d types, got %d", len(expected), len(types))
		return
	}

	for i, expectedType := range expected {
		if types[i] != expectedType {
			t.Errorf("Column %d: expected %v, got %v", i, expectedType, types[i])
		}
	}
}

func TestTypeGuesserDatePriority(t *testing.T) {
	tg := newTypeGuesser(false)
	tg.setLength(1)

	// Test that date detection takes priority over string when appropriate
	records := [][]string{
		{"2023-01-15"},
		{"2023-02-20"},
		{"2023-03-25"},
		{"2023-04-30"},
		{"some text"}, // This should make it string type
	}

	for _, record := range records {
		tg.guessTypes(record)
	}

	types := tg.getSchema().GetDataTypes()
	if len(types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(types))
		return
	}

	// With mixed data, it should default to string
	if types[0] != meta.StringType {
		t.Errorf("Expected StringType, got %v", types[0])
	}
}

func TestTypeGuesserDateOnly(t *testing.T) {
	tg := newTypeGuesser(false)
	tg.setLength(1)

	// Test with only date data
	records := [][]string{
		{"2023-01-15"},
		{"2023-02-20"},
		{"2023-03-25"},
		{"2023-04-30"},
	}

	for _, record := range records {
		tg.guessTypes(record)
	}

	types := tg.getSchema().GetDataTypes()
	if len(types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(types))
		return
	}

	if types[0] != meta.TimeType {
		t.Errorf("Expected TimeType, got %v", types[0])
	}
}

func TestTypeGuesserDatetimeOnly(t *testing.T) {
	tg := newTypeGuesser(false)
	tg.setLength(1)

	// Test with only datetime data
	records := [][]string{
		{"2023-01-15 14:30:00"},
		{"2023-02-20 15:45:30"},
		{"2023-03-25 09:15:45"},
		{"2023-04-30 22:10:20"},
	}

	for _, record := range records {
		tg.guessTypes(record)
	}

	types := tg.getSchema().GetDataTypes()
	if len(types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(types))
		return
	}

	if types[0] != meta.TimeType {
		t.Errorf("Expected TimeType, got %v", types[0])
	}
}

func TestTypeGuesserDateVsDatetime(t *testing.T) {
	tg := newTypeGuesser(false)
	tg.setLength(1)

	// Test with mixed date and datetime data
	records := [][]string{
		{"2023-01-15"},          // date
		{"2023-02-20 14:30:00"}, // datetime
		{"2023-03-25"},          // date
		{"2023-04-30 15:45:30"}, // datetime
	}

	for _, record := range records {
		tg.guessTypes(record)
	}

	types := tg.getSchema().GetDataTypes()
	if len(types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(types))
		return
	}

	// Both date and datetime should result in TimeType
	if types[0] != meta.TimeType {
		t.Errorf("Expected TimeType, got %v", types[0])
	}
}

// Mock RowDataProvider for testing
type mockRowDataProvider struct {
	data [][]string
	pos  int
}

func (m *mockRowDataProvider) Read() ([]string, error) {
	if m.pos >= len(m.data) {
		return nil, nil // EOF
	}
	row := m.data[m.pos]
	m.pos++
	return row, nil
}

func TestReadRowDataWithDates(t *testing.T) {
	// Skip header row for type detection - only use actual data
	mockData := [][]string{
		{"Alice", "25", "1998-05-15", "2023-01-15 10:30:00"},
		{"Bob", "30", "1993-12-20", "2023-02-20 14:45:30"},
		{"Charlie", "35", "1988-08-10", "2023-03-25 09:15:45"},
	}

	provider := &mockRowDataProvider{data: mockData}

	series, err := readRowData(provider, 3, false, 10, nil, ctx)
	if err != nil {
		t.Fatalf("readRowData failed: %v", err)
	}

	if len(series) != 4 {
		t.Errorf("Expected 4 series, got %d", len(series))
		return
	}

	// Check that the date columns are properly typed
	expectedTypes := []meta.BaseType{meta.StringType, meta.Int64Type, meta.TimeType, meta.TimeType}

	for i, s := range series {
		actualType := s.Type()
		if actualType != expectedTypes[i] {
			t.Errorf("Column %d: expected %v, got %v", i, expectedTypes[i], actualType)
		}
	}
}

func TestReadRowDataWithNullsAndDates(t *testing.T) {
	// Skip header row for type detection - only use actual data
	mockData := [][]string{
		{"Alice", "25", "1998-05-15", "2023-01-15 10:30:00"},
		{"Bob", "null", "1993-12-20", "null"},
		{"Charlie", "35", "null", "2023-03-25 09:15:45"},
	}

	provider := &mockRowDataProvider{data: mockData}

	series, err := readRowData(provider, 3, false, 10, nil, ctx)
	if err != nil {
		t.Fatalf("readRowData failed: %v", err)
	}

	if len(series) != 4 {
		t.Errorf("Expected 4 series, got %d", len(series))
		return
	}

	// Check that the date columns are properly typed
	expectedTypes := []meta.BaseType{meta.StringType, meta.Int64Type, meta.TimeType, meta.TimeType}

	for i, s := range series {
		actualType := s.Type()
		if actualType != expectedTypes[i] {
			t.Errorf("Column %d: expected %v, got %v", i, expectedTypes[i], actualType)
		}
	}
}

func TestReadRowDataStrict(t *testing.T) {
	mockData := [][]string{
		{"t", "1", "0.1"},
		{"f", "2", "0.2"},
		{"t", "3", "0.3"},
		{"f", "4", "0.4"},
		{"t", "5", "0.5"},
		{"f", "6", "0.6"},
		{"t", "7", "0.7"},
		{"f", "8", "0.8"},
		{"a", "b", "c"},
	}

	provider := &mockRowDataProvider{data: mockData}

	series, err := readRowData(provider, 100, false, -1, nil, ctx)
	if err != nil {
		t.Fatalf("readRowData failed: %v", err)
	}

	if len(series) != 3 {
		t.Errorf("Expected 3 series, got %d", len(series))
		return
	}

	fmt.Println(series[0].Type())
	fmt.Println(series[1].Type())
	fmt.Println(series[2].Type())
}
