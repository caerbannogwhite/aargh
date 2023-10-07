package gandalff

import (
	"math/rand"
	"testing"
	"typesys"
)

var stringPool *StringPool

func init() {
	stringPool = NewStringPool()
}

func Test_SeriesString_Base(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	mask := []bool{false, false, true, false, false, true, false, false, true, false}

	// Create a new SeriesString.
	s := NewSeriesString(data, mask, true, stringPool)

	// Check the length.
	if s.Len() != 10 {
		t.Errorf("Expected length of 10, got %d", s.Len())
	}

	// Check the type.
	if s.Type() != typesys.StringType {
		t.Errorf("Expected type of StringType, got %s", s.Type().ToString())
	}

	// Check the data.
	for i, v := range s.Data().([]string) {
		if v != data[i] {
			t.Errorf("Expected data of []string{\"a\", \"b\", \"c\", \"d\", \"e\", \"f\", \"g\", \"h\", \"i\", \"j\"}, got %v", s.Data())
		}
	}

	// Check the nullability.
	if !s.IsNullable() {
		t.Errorf("Expected IsNullable() to be true, got false")
	}

	// Check the null mask.
	for i, v := range s.GetNullMask() {
		if v != mask[i] {
			t.Errorf("Expected nullMask of []bool{false, false, false, false, true, false, true, false, false, true}, got %v", s.GetNullMask())
		}
	}

	// Check the null values.
	for i := range s.Data().([]string) {
		if s.IsNull(i) != mask[i] {
			t.Errorf("Expected IsNull(%d) to be %t, got %t", i, mask[i], s.IsNull(i))
		}
	}

	// Check the null count.
	if s.NullCount() != 3 {
		t.Errorf("Expected NullCount() to be 3, got %d", s.NullCount())
	}

	// Check the HasNull method.
	if !s.HasNull() {
		t.Errorf("Expected HasNull() to be true, got false")
	}

	// Check the SetNull() method.
	for i := range s.Data().([]string) {
		s.SetNull(i)
	}

	// Check the null values.
	for i := range s.Data().([]string) {
		if !s.IsNull(i) {
			t.Errorf("Expected IsNull(%d) to be true, got false", i)
		}
	}

	// Check the null count.
	if s.NullCount() != 10 {
		t.Errorf("Expected NullCount() to be 10, got %d", s.NullCount())
	}

	// Check the Get method.
	for i, v := range s.Data().([]string) {
		if s.Get(i) != v {
			t.Errorf("Expected Get(%d) to be %s, got %s", i, v, s.Get(i))
		}
	}

	// Check the Set method.
	for i, v := range s.Data().([]string) {
		s.Set(i, v+"x")
	}

	// Check the data.
	for i, v := range s.Data().([]string) {
		if v != data[i]+"x" {
			t.Errorf("Expected data of []string{\"ax\", \"bx\", \"cx\", \"dx\", \"ex\", \"fx\", \"gx\", \"hx\", \"ix\", \"jx\"}, got %v", s.Data())
		}
	}

	// Check the MakeNullable() method.
	p := NewSeriesString(data, nil, true, stringPool)

	// Check the nullability.
	if p.IsNullable() {
		t.Errorf("Expected IsNullable() to be false, got true")
	}

	// Set values to null.
	p.SetNull(1)
	p.SetNull(3)
	p.SetNull(5)

	// Check the null count.
	if p.NullCount() != 0 {
		t.Errorf("Expected NullCount() to be 0, got %d", p.NullCount())
	}

	// Make the series nullable.
	p = p.MakeNullable().(SeriesString)

	// Check the nullability.
	if !p.IsNullable() {
		t.Errorf("Expected IsNullable() to be true, got false")
	}

	// Check the null count.
	if p.NullCount() != 0 {
		t.Errorf("Expected NullCount() to be 0, got %d", p.NullCount())
	}

	p.SetNull(1)
	p.SetNull(3)
	p.SetNull(5)

	// Check the null count.
	if p.NullCount() != 3 {
		t.Errorf("Expected NullCount() to be 3, got %d", p.NullCount())
	}
}

func Test_SeriesString_Append(t *testing.T) {
	dataA := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	dataB := []string{"k", "l", "m", "n", "o", "p", "q", "r", "s", "t"}
	dataC := []string{"u", "v", "w", "x", "y", "z", "1", "2", "3", "4"}

	maskA := []bool{false, false, true, false, false, true, false, false, true, false}
	maskB := []bool{false, false, false, false, true, false, true, false, false, true}
	maskC := []bool{true, true, true, true, true, true, true, true, true, true}

	// Create two new series.
	sA := NewSeriesString(dataA, maskA, true, stringPool)
	sB := NewSeriesString(dataB, maskB, true, stringPool)
	sC := NewSeriesString(dataC, maskC, true, stringPool)

	// Append the series.
	result := sA.Append(sB).Append(sC)

	// Check the length.
	if result.Len() != 30 {
		t.Errorf("Expected length of 30, got %d", result.Len())
	}

	// Check the data.
	for i, v := range result.Data().([]string) {
		if i < 10 {
			if v != dataA[i] {
				t.Errorf("Expected %s, got %s at index %d", dataA[i], v, i)
			}
		} else if i < 20 {
			if v != dataB[i-10] {
				t.Errorf("Expected %s, got %s at index %d", dataB[i-10], v, i)
			}
		} else {
			if v != dataC[i-20] {
				t.Errorf("Expected %s, got %s at index %d", dataC[i-20], v, i)
			}
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if i < 10 {
			if v != maskA[i] {
				t.Errorf("Expected nullMask %t, got %t at index %d", maskA[i], v, i)
			}
		} else if i < 20 {
			if v != maskB[i-10] {
				t.Errorf("Expected nullMask %t, got %t at index %d", maskB[i-10], v, i)
			}
		} else {
			if v != maskC[i-20] {
				t.Errorf("Expected nullMask %t, got %t at index %d", maskC[i-20], v, i)
			}
		}
	}

	// Append random values.
	dataD := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	sD := NewSeriesString(dataD, nil, true, stringPool).MakeNullable().(SeriesString)

	// Check the original data.
	for i, v := range sD.Data().([]string) {
		if v != dataD[i] {
			t.Errorf("Expected %s, got %s at index %d", dataD[i], v, i)
		}
	}

	alpha := "abcdefghijklmnopqrstuvwxyz0123456789"

	for i := 0; i < 100; i++ {
		r := string(alpha[rand.Intn(len(alpha))])
		switch i % 4 {
		case 0:
			sD = sD.Append(r).(SeriesString)
		case 1:
			sD = sD.Append([]string{r}).(SeriesString)
		case 2:
			sD = sD.Append(NullableString{true, r}).(SeriesString)
		case 3:
			sD = sD.Append([]NullableString{{false, r}}).(SeriesString)
		}

		if sD.Get(i+10) != r {
			t.Errorf("Expected %t, got %t at index %d (case %d)", true, sD.Get(i+10), i+10, i%4)
		}
	}
}

func Test_SeriesString_Cast(t *testing.T) {
	data := []string{"true", "false", "0", "3", "4", "5", "hello", "7", "8", "world"}
	mask := []bool{false, false, true, false, false, true, false, false, true, false}

	// Create a new series.
	s := NewSeriesString(data, mask, true, stringPool)

	// Cast to bool.
	resBool := s.Cast(typesys.BoolType)

	// Check the data.
	for i, v := range resBool.Data().([]bool) {
		switch i {
		case 0:
			if v != true {
				t.Errorf("Expected %t, got %t at index %d", true, v, i)
			}
		default:
			if v != false {
				t.Errorf("Expected %t, got %t at index %d", false, v, i)
			}
		}
	}

	// Check the null mask.
	for i, v := range resBool.GetNullMask() {
		switch i {
		case 0, 1:
			if v != false {
				t.Errorf("Expected nullMask %t, got %t at index %d", false, v, i)
			}
		default:
			if v != true {
				t.Errorf("Expected nullMask %t, got %t at index %d", true, v, i)
			}
		}
	}

	// Cast to int32.
	resInt := s.Cast(typesys.Int32Type)
	expectedInt32 := []int32{0, 0, 0, 3, 4, 0, 0, 7, 0, 0}

	// Check the data.
	for i, v := range resInt.Data().([]int32) {
		if v != expectedInt32[i] {
			t.Errorf("Expected %d, got %d at index %d", expectedInt32[i], v, i)
		}
	}

	expectedMask := []bool{true, true, true, false, false, true, true, false, true, true}

	// Check the null mask.
	for i, v := range resInt.GetNullMask() {
		if v != expectedMask[i] {
			t.Errorf("Expected nullMask %t, got %t at index %d", expectedMask[i], v, i)
		}
	}

	// Cast to int64.
	resInt64 := s.Cast(typesys.Int64Type)
	expectedInt64 := []int64{0, 0, 0, 3, 4, 0, 0, 7, 0, 0}

	// Check the data.
	for i, v := range resInt64.Data().([]int64) {
		if v != expectedInt64[i] {
			t.Errorf("Expected %d, got %d at index %d", expectedInt64[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range resInt64.GetNullMask() {
		if v != expectedMask[i] {
			t.Errorf("Expected nullMask %t, got %t at index %d", expectedMask[i], v, i)
		}
	}

	// Cast to float64.
	resFloat64 := s.Cast(typesys.Float64Type)
	expectedFloat64 := []float64{0, 0, 0, 3, 4, 0, 0, 7, 0, 0}

	// Check the data.
	for i, v := range resFloat64.Data().([]float64) {
		if v != expectedFloat64[i] {
			t.Errorf("Expected %f, got %f at index %d", expectedFloat64[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range resFloat64.GetNullMask() {
		if v != expectedMask[i] {
			t.Errorf("Expected nullMask %t, got %t at index %d", expectedMask[i], v, i)
		}
	}

	// Cast to error.
	castError := s.Cast(typesys.ErrorType)

	// Check the message.
	if castError.(SeriesError).msg != "SeriesString.Cast: invalid type Error" {
		t.Errorf("Expected error, got %v", castError)
	}
}

func Test_SeriesString_Filter(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t"}
	mask := []bool{false, true, false, false, true, false, false, true, false, false, true, false, false, true, false, false, true, false, false, true}

	// Create a new series.
	s := NewSeriesString(data, mask, true, stringPool)

	// Filter mask.
	filterMask := []bool{true, false, true, true, false, true, true, false, true, true, true, false, true, true, false, true, true, false, true, true}
	filterIndeces := []int{0, 2, 3, 5, 6, 8, 9, 10, 12, 13, 15, 16, 18, 19}

	result := []string{"a", "c", "d", "f", "g", "i", "j", "k", "m", "n", "p", "q", "s", "t"}
	resultMask := []bool{false, false, false, false, false, false, false, true, false, true, false, true, false, true}

	/////////////////////////////////////////////////////////////////////////////////////
	// 							Check the Filter() method.
	filtered := s.Filter(NewSeriesBool(filterMask, nil, true, nil))

	// Check the length.
	if filtered.Len() != len(result) {
		t.Errorf("Expected length of %d, got %d", len(result), filtered.Len())
	}

	// Check the data.
	for i, v := range filtered.Data().([]string) {
		if v != result[i] {
			t.Errorf("Expected %v, got %v at index %d", result[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range filtered.GetNullMask() {
		if v != resultMask[i] {
			t.Errorf("Expected nullMask of %v, got %v at index %d", resultMask[i], v, i)
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////
	// 							Check the Filter() method.
	filtered = s.Filter(filterMask)

	// Check the length.
	if filtered.Len() != len(result) {
		t.Errorf("Expected length of %d, got %d", len(result), filtered.Len())
	}

	// Check the data.
	for i, v := range filtered.Data().([]string) {
		if v != result[i] {
			t.Errorf("Expected %v, got %v at index %d", result[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range filtered.GetNullMask() {
		if v != resultMask[i] {
			t.Errorf("Expected nullMask of %v, got %v at index %d", resultMask[i], v, i)
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////
	// 							Check the Filter() method.
	filtered = s.Filter(filterIndeces)

	// Check the length.
	if filtered.Len() != len(result) {
		t.Errorf("Expected length of %d, got %d", len(result), filtered.Len())
	}

	// Check the data.
	for i, v := range filtered.Data().([]string) {
		if v != result[i] {
			t.Errorf("Expected %v, got %v at index %d", result[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range filtered.GetNullMask() {
		if v != resultMask[i] {
			t.Errorf("Expected nullMask of %v, got %v at index %d", resultMask[i], v, i)
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////

	// try to filter by a series with a different length.
	filtered = filtered.Filter(filterMask)

	if e, ok := filtered.(SeriesError); !ok || e.GetError() != "SeriesString.FilterByMask: mask length (20) does not match series length (14)" {
		t.Errorf("Expected SeriesError, got %v", filtered)
	}

	// Another test.
	data = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w"}
	mask = []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}

	// Create a new series.
	s = NewSeriesString(data, mask, true, stringPool)

	// Filter mask.
	filterMask = []bool{true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, true}
	filterIndeces = []int{0, 15, 22}

	result = []string{"a", "p", "w"}

	/////////////////////////////////////////////////////////////////////////////////////
	// 							Check the Filter() method.
	filtered = s.Filter(filterMask)

	// Check the length.
	if filtered.Len() != 3 {
		t.Errorf("Expected length of 3, got %d", filtered.Len())
	}

	// Check the data.
	for i, v := range filtered.Data().([]string) {
		if v != result[i] {
			t.Errorf("Expected %v, got %v at index %d", result[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range filtered.GetNullMask() {
		if v != true {
			t.Errorf("Expected nullMask of %v, got %v at index %d", true, v, i)
		}
	}

	/////////////////////////////////////////////////////////////////////////////////////
	// 							Check the Filter() method.
	filtered = s.Filter(filterIndeces)

	// Check the length.
	if filtered.Len() != 3 {
		t.Errorf("Expected length of 3, got %d", filtered.Len())
	}

	// Check the data.
	for i, v := range filtered.Data().([]string) {
		if v != result[i] {
			t.Errorf("Expected %v, got %v at index %d", result[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range filtered.GetNullMask() {
		if v != true {
			t.Errorf("Expected nullMask of %v, got %v at index %d", true, v, i)
		}
	}
}

func Test_SeriesString_Map(t *testing.T) {
	data := []string{"", "hello", "world", "this", "is", "a", "test", "of", "the", "map", "function", "in", "the", "series", "", "this", "is", "a", "", "test"}
	nullMask := []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, false, true, true, true, false, true, false}

	// Create a new series.
	s := NewSeriesString(data, nullMask, true, NewStringPool())

	// Map the series to bool.
	resBool := s.Map(func(v any) any {
		return v.(string) == ""
	})

	expectedBool := []bool{true, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, true, false}
	for i, v := range resBool.Data().([]bool) {
		if v != expectedBool[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedBool[i], v, i)
		}
	}

	// Map the series to int32.
	resInt := s.Map(func(v any) any {
		return int32(len(v.(string)))
	})

	expectedInt32 := []int32{0, 5, 5, 4, 2, 1, 4, 2, 3, 3, 8, 2, 3, 6, 0, 4, 2, 1, 0, 4}
	for i, v := range resInt.Data().([]int32) {
		if v != expectedInt32[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedInt32[i], v, i)
		}
	}

	// Map the series to int64.
	resInt64 := s.Map(func(v any) any {
		return int64(len(v.(string)))
	})

	expectedInt64 := []int64{0, 5, 5, 4, 2, 1, 4, 2, 3, 3, 8, 2, 3, 6, 0, 4, 2, 1, 0, 4}
	for i, v := range resInt64.Data().([]int64) {
		if v != expectedInt64[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedInt64[i], v, i)
		}
	}

	// Map the series to float64.
	resFloat64 := s.Map(func(v any) any {
		return -float64(len(v.(string)))
	})

	expectedFloat64 := []float64{-0, -5, -5, -4, -2, -1, -4, -2, -3, -3, -8, -2, -3, -6, -0, -4, -2, -1, -0, -4}
	for i, v := range resFloat64.Data().([]float64) {
		if v != expectedFloat64[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedFloat64[i], v, i)
		}
	}

	// Map the series to string.
	resString := s.Map(func(v any) any {
		if v.(string) == "" {
			return "empty"
		}
		return ""
	})

	expectedString := []string{"empty", "", "", "", "", "", "", "", "", "", "", "", "", "", "empty", "", "", "", "empty", ""}
	for i, v := range resString.Data().([]string) {
		if v != expectedString[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedString[i], v, i)
		}
	}
}

func Test_SeriesString_Group(t *testing.T) {
	var partMap map[int64][]int
	pool := NewStringPool()

	data1 := []string{"foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo"}
	data1Mask := []bool{false, false, false, false, false, true, true, true, true, true}
	data2 := []string{"foo", "foo", "bar", "bar", "foo", "foo", "bar", "bar", "foo", "foo"}
	data2Mask := []bool{false, true, false, true, false, true, false, true, false, true}
	data3 := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	data3Mask := []bool{false, false, false, false, false, true, true, true, true, true}

	// Test 1
	s1 := NewSeriesString(data1, data1Mask, true, pool).
		group()

	p1 := s1.GetPartition().getMap()
	if len(p1) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(p1))
	}

	partMap = map[int64][]int{
		0: {0, 1, 2, 3, 4},
		1: {5, 6, 7, 8, 9},
	}
	if !checkEqPartitionMap(p1, partMap, nil, "String Group") {
		t.Errorf("Expected partition map of %v, got %v", partMap, p1)
	}

	// Test 2
	s2 := NewSeriesString(data2, data2Mask, true, pool).
		GroupBy(s1.GetPartition())

	p2 := s2.GetPartition().getMap()
	if len(p2) != 6 {
		t.Errorf("Expected 6 groups, got %d", len(p2))
	}

	partMap = map[int64][]int{
		0: {0, 4},
		1: {1, 3},
		2: {2},
		3: {5, 7, 9},
		4: {6},
		5: {8},
	}
	if !checkEqPartitionMap(p2, partMap, nil, "String Group") {
		t.Errorf("Expected partition map of %v, got %v", partMap, p2)
	}

	// Test 3
	s3 := NewSeriesString(data3, data3Mask, true, pool).
		GroupBy(s2.GetPartition())

	p3 := s3.GetPartition().getMap()
	if len(p3) != 8 {
		t.Errorf("Expected 8 groups, got %d", len(p3))
	}

	partMap = map[int64][]int{
		0: {0},
		1: {1},
		2: {2},
		3: {3},
		4: {4},
		5: {5, 7, 9},
		6: {6},
		7: {8},
	}
	if !checkEqPartitionMap(p3, partMap, nil, "String Group") {
		t.Errorf("Expected partition map of %v, got %v", partMap, p3)
	}

	// debugPrintPartition(s1.GetPartition(), s1)
	// debugPrintPartition(s2.GetPartition(), s1, s2)
	// debugPrintPartition(s3.GetPartition(), s1, s2, s3)

	partMap = nil
}

func Test_SeriesString_Sort(t *testing.T) {
	data := []string{"w", "w", "d", "y", "b", "e", "a", "e", "e", "b", "l", "u", "a", "g", "w", "u", "{", "x", "t", "h"}
	mask := []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}

	// Create a new series.
	s := NewSeriesString(data, nil, true, NewStringPool())

	// Sort the series.
	sorted := s.Sort()

	// Check the data.
	expected := []string{"a", "a", "b", "b", "d", "e", "e", "e", "g", "h", "l", "t", "u", "u", "w", "w", "w", "x", "y", "{"}
	if !checkEqSliceString(sorted.Data().([]string), expected, nil, "") {
		t.Errorf("SeriesString.Sort() failed, expecting %v, got %v", expected, sorted.Data().([]string))
	}

	// Create a new series.
	s = NewSeriesString(data, mask, true, NewStringPool())

	// Sort the series.
	sorted = s.Sort()

	// Check the data.
	expected = []string{"a", "a", "b", "d", "e", "l", "t", "w", "w", "{", "e", "u", "b", "g", "e", "u", "y", "x", "w", "h"}
	if !checkEqSliceString(sorted.Data().([]string), expected, nil, "") {
		t.Errorf("SeriesString.Sort() failed, expecting %v, got %v", expected, sorted.Data().([]string))
	}

	// Check the null mask.
	expectedMask := []bool{false, false, false, false, false, false, false, false, false, false, true, true, true, true, true, true, true, true, true, true}
	if !checkEqSliceBool(sorted.GetNullMask(), expectedMask, nil, "") {
		t.Errorf("SeriesString.Sort() failed, expecting %v, got %v", expectedMask, sorted.GetNullMask())
	}
}

func Test_SeriesString_Arithmetic_Add(t *testing.T) {
	pool := NewStringPool()

	bools := NewSeriesBool([]bool{true}, nil, true, nil)
	boolv := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, nil)
	bools_ := NewSeriesBool([]bool{true}, nil, true, nil).SetNullMask([]bool{true})
	boolv_ := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, nil).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i32s := NewSeriesInt32([]int32{2}, nil, true, nil)
	i32v := NewSeriesInt32([]int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, nil)
	i32s_ := NewSeriesInt32([]int32{2}, nil, true, nil).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt32([]int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, nil).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i64s := NewSeriesInt64([]int64{2}, nil, true, nil)
	i64v := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, nil)
	i64s_ := NewSeriesInt64([]int64{2}, nil, true, nil).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, nil).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	f64s := NewSeriesFloat64([]float64{2}, nil, true, nil)
	f64v := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, nil)
	f64s_ := NewSeriesFloat64([]float64{2}, nil, true, nil).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, nil).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	ss := NewSeriesString([]string{"2"}, nil, true, pool)
	sv := NewSeriesString([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, nil, true, pool)
	ss_ := NewSeriesString([]string{"2"}, nil, true, pool).SetNullMask([]bool{true})
	sv_ := NewSeriesString([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, nil, true, pool).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	// scalar | bool
	if !checkEqSlice(ss.Add(bools).Data().([]string), []string{"2true"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(boolv).Data().([]string), []string{"2true", "2false", "2true", "2false", "2true", "2false", "2true", "2true", "2false", "2false"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(bools_).GetNullMask(), []bool{true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}

	// scalar | int32
	if !checkEqSlice(ss.Add(i32s).Data().([]string), []string{"22"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(i32v).Data().([]string), []string{"21", "22", "23", "24", "25", "26", "27", "28", "29", "210"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(i32s_).GetNullMask(), []bool{true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}

	// scalar | int64
	if !checkEqSlice(ss.Add(i64s).Data().([]string), []string{"22"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(i64v).Data().([]string), []string{"21", "22", "23", "24", "25", "26", "27", "28", "29", "210"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(i64s_).GetNullMask(), []bool{true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}

	// scalar | float64
	if !checkEqSlice(ss.Add(f64s).Data().([]string), []string{"22"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(f64v).Data().([]string), []string{"21", "22", "23", "24", "25", "26", "27", "28", "29", "210"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(f64s_).GetNullMask(), []bool{true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}

	// scalar | string
	if !checkEqSlice(ss.Add(ss).Data().([]string), []string{"22"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(sv).Data().([]string), []string{"21", "22", "23", "24", "25", "26", "27", "28", "29", "210"}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(ss_).GetNullMask(), []bool{true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
	if !checkEqSlice(ss.Add(sv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "String Add") {
		t.Errorf("Error in String Add")
	}
}

func Test_SeriesString_Logical_Eq(t *testing.T) {
	pool := NewStringPool()
	ss := NewSeriesString([]string{"1"}, nil, true, pool)
	ss_ := NewSeriesString([]string{"1"}, nil, true, pool).SetNullMask([]bool{true})
	sv := NewSeriesString([]string{"1", "2", "3"}, nil, true, pool)
	sv_ := NewSeriesString([]string{"1", "2", "3"}, nil, true, pool).SetNullMask([]bool{true, true, false})

	// scalar | scalar
	res := ss.Eq(ss)
	if res.Data().([]bool)[0] != true {
		t.Errorf("Expected %v, got %v", true, res.Data().([]bool)[0])
	}

	res = ss.Eq(ss_)
	if res.IsNull(0) == false {
		t.Errorf("Expected %v, got %v", true, res.IsNull(0))
	}

	// scalar | vector
	res = ss.Eq(sv)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{true, false, false}, res.Data().([]bool))
	}

	res = ss.Eq(sv_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}

	// vector | scalar
	res = sv.Eq(ss)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{true, false, false}, res.Data().([]bool))
	}

	res = sv_.Eq(ss)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}

	// vector | vector
	res = sv.Eq(sv)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.Data().([]bool))
	}

	res = sv.Eq(sv_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}

	res = sv_.Eq(sv)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}

	res = sv_.Eq(sv_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}
}

func Test_SeriesString_Logical_Ne(t *testing.T) {
	// TODO: add tests for all types
}

func Test_SeriesString_Logical_Lt(t *testing.T) {
	pool := NewStringPool()
	ss := NewSeriesString([]string{"1"}, nil, true, pool)
	ss_ := NewSeriesString([]string{"1"}, nil, true, pool).SetNullMask([]bool{true})
	sv := NewSeriesString([]string{"1", "2", "3"}, nil, true, pool)
	sv_ := NewSeriesString([]string{"1", "2", "3"}, nil, true, pool).SetNullMask([]bool{true, true, false})

	// scalar | scalar
	res := ss.Lt(ss)
	if res.Data().([]bool)[0] != false {
		t.Errorf("Expected %v, got %v", false, res.Data().([]bool)[0])
	}

	res = ss.Lt(ss_)
	if res.IsNull(0) == false {
		t.Errorf("Expected %v, got %v", true, res.IsNull(0))
	}

	// scalar | vector
	res = ss.Lt(sv)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{false, true, true}, res.Data().([]bool))
	}

	res = ss.Lt(sv_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}

	// vector | scalar
	res = sv.Lt(ss)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.Data().([]bool))
	}

	res = sv_.Lt(ss)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}

	// vector | vector
	res = sv.Lt(sv)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.Data().([]bool))
	}

	res = sv.Lt(sv_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}

	res = sv_.Lt(sv)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}

	res = sv_.Lt(sv_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, []bool{res.IsNull(0), res.IsNull(1), res.IsNull(2)})
	}
}

func Test_SeriesString_Logical_Le(t *testing.T) {
	// TODO: add tests for all types
}

func Test_SeriesString_Logical_Gt(t *testing.T) {
	// TODO: add tests for all types
}

func Test_SeriesString_Logical_Ge(t *testing.T) {
	// TODO: add tests for all types
}
