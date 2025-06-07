package series

import (
	"math"
	"math/rand"
	"testing"

	"github.com/caerbannogwhite/gandalff"
	"github.com/caerbannogwhite/gandalff/meta"
	"github.com/caerbannogwhite/gandalff/utils"
)

func Test_SeriesInt64_Base(t *testing.T) {

	data := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	mask := []bool{false, false, true, false, false, true, false, false, true, false}

	// Create a new Int64s.
	s := NewSeriesInt64(data, mask, true, ctx)

	// Check the length.
	if s.Len() != 10 {
		t.Errorf("Expected length of 10, got %d", s.Len())
	}

	// Check the type.
	if s.Type() != meta.Int64Type {
		t.Errorf("Expected type of Int64Type, got %s", s.Type().ToString())
	}

	// Check the data.
	for i, v := range s.Data().([]int64) {
		if v != data[i] {
			t.Errorf("Expected data of []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, got %v", s.Data())
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
	for i := range s.Data().([]int64) {
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

	// Check the Set() with a null value.
	for i := range s.Data().([]int64) {
		s.Set(i, nil)
	}

	// Check the null values.
	for i := range s.Data().([]int64) {
		if !s.IsNull(i) {
			t.Errorf("Expected IsNull(%d) to be true, got false", i)
		}
	}

	// Check the null count.
	if s.NullCount() != 10 {
		t.Errorf("Expected NullCount() to be 10, got %d", s.NullCount())
	}

	// Check the Get method.
	for i := range s.Data().([]int64) {
		if s.Get(i).(int64) != data[i] {
			t.Errorf("Expected Get(%d) to be %d, got %d", i, data[i], s.Get(i).(int64))
		}
	}

	// Check the Set method.
	for i := range s.Data().([]int64) {
		s.Set(i, int64(i+10))
	}

	// Check the data.
	for i, v := range s.Data().([]int64) {
		if v != int64(i+10) {
			t.Errorf("Expected data of []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, got %v", s.Data())
		}
	}

	// Check the MakeNullable() method.
	p := NewSeriesInt64(data, nil, true, ctx)

	// Check the nullability.
	if p.IsNullable() {
		t.Errorf("Expected IsNullable() to be false, got true")
	}

	// Set values to null.
	p.Set(1, nil)
	p.Set(3, nil)
	p.Set(5, nil)

	// Check the null count.
	if p.NullCount() != 0 {
		t.Errorf("Expected NullCount() to be 0, got %d", p.NullCount())
	}

	// Make the series nullable.
	p = p.MakeNullable().(Int64s)

	// Check the nullability.
	if !p.IsNullable() {
		t.Errorf("Expected IsNullable() to be true, got false")
	}

	// Check the null count.
	if p.NullCount() != 0 {
		t.Errorf("Expected NullCount() to be 0, got %d", p.NullCount())
	}

	p.Set(1, nil)
	p.Set(3, nil)
	p.Set(5, nil)

	// Check the null count.
	if p.NullCount() != 3 {
		t.Errorf("Expected NullCount() to be 3, got %d", p.NullCount())
	}
}

func Test_SeriesInt64_Take(t *testing.T) {

	data := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	mask := []bool{false, false, true, false, false, true, false, false, true, false}

	// Create a new Int64s.
	s := NewSeriesInt64(data, mask, true, ctx)

	// Take the first 5 values.
	result := s.Take(0, 5, 1)

	// Check the length.
	if result.Len() != 5 {
		t.Errorf("Expected length of 5, got %d", result.Len())
	}

	// Check the data.
	for i, v := range result.Data().([]int64) {
		if v != data[i] {
			t.Errorf("Expected data of []int{1, 2, 3, 4, 5}, got %v", result.Data())
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if v != mask[i] {
			t.Errorf("Expected nullMask of []bool{false, false, false, false, true}, got %v", result.GetNullMask())
		}
	}

	// Take the last 5 values.
	result = s.Take(5, 10, 1)

	// Check the length.
	if result.Len() != 5 {
		t.Errorf("Expected length of 5, got %d", result.Len())
	}

	// Check the data.
	for i, v := range result.Data().([]int64) {
		if v != data[i+5] {
			t.Errorf("Expected data of []int{6, 7, 8, 9, 10}, got %v", result.Data())
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if v != mask[i+5] {
			t.Errorf("Expected nullMask of []bool{true, false, false, true, false}, got %v", result.GetNullMask())
		}
	}

	// Take the first 5 values in steps of 2.
	result = s.Take(0, 6, 2)

	// Check the length.
	if result.Len() != 3 {
		t.Errorf("Expected length of 3, got %d", result.Len())
	}

	// Check the data.
	for i, v := range result.Data().([]int64) {
		if v != data[i*2] {
			t.Errorf("Expected data of []int{1, 3, 5}, got %v", result.Data())
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if v != mask[i*2] {
			t.Errorf("Expected nullMask of []bool{false, false, true}, got %v", result.GetNullMask())
		}
	}

	// Take the last 5 values in steps of 2.
	result = s.Take(5, 11, 2)

	// Check the length.
	if result.Len() != 3 {
		t.Errorf("Expected length of 3, got %d", result.Len())
	}

	// Check the data.
	for i, v := range result.Data().([]int64) {
		if v != data[i*2+5] {
			t.Errorf("Expected data of []int{6, 8, 10}, got %v", result.Data())
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if v != mask[i*2+5] {
			t.Errorf("Expected nullMask of []bool{true, false, false}, got %v", result.GetNullMask())
		}
	}
}

func Test_SeriesInt64_Append(t *testing.T) {
	dataA := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	dataB := []int64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	dataC := []int64{21, 22, 23, 24, 25, 26, 27, 28, 29, 30}

	maskA := []bool{false, false, true, false, false, true, false, false, true, false}
	maskB := []bool{false, false, false, false, true, false, true, false, false, true}
	maskC := []bool{true, true, true, true, true, true, true, true, true, true}

	// Create two new series.
	sA := NewSeriesInt64(dataA, maskA, true, ctx)
	sB := NewSeriesInt64(dataB, maskB, true, ctx)
	sC := NewSeriesInt64(dataC, maskC, true, ctx)

	// Append the series.
	result := sA.Append(sB).Append(sC)

	// Check the length.
	if result.Len() != 30 {
		t.Errorf("Expected length of 30, got %d", result.Len())
	}

	// Check the data.
	for i, v := range result.Data().([]int64) {
		if i < 10 {
			if v != dataA[i] {
				t.Errorf("Expected %d, got %d at index %d", dataA[i], v, i)
			}
		} else if i < 20 {
			if v != dataB[i-10] {
				t.Errorf("Expected %d, got %d at index %d", dataB[i-10], v, i)
			}
		} else {
			if v != dataC[i-20] {
				t.Errorf("Expected %d, got %d at index %d", dataC[i-20], v, i)
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
	s := NewSeriesInt64([]int64{}, nil, true, ctx)

	for i := 0; i < 100; i++ {
		r := int64(rand.Intn(100))
		switch i % 4 {
		case 0:
			s = s.Append(r).(Int64s)
		case 1:
			s = s.Append([]int64{r}).(Int64s)
		case 2:
			s = s.Append(gandalff.NullableInt64{Valid: true, Value: r}).(Int64s)
		case 3:
			s = s.Append([]gandalff.NullableInt64{{Valid: false, Value: r}}).(Int64s)
		}

		if s.Get(i) != r {
			t.Errorf("Expected %t, got %t at index %d (case %d)", true, s.Get(i), i, i%4)
		}
	}

	// Append nil
	s = NewSeriesInt64([]int64{}, nil, true, ctx)

	for i := 0; i < 100; i++ {
		s = s.Append(nil).(Int64s)
		if !s.IsNull(i) {
			t.Errorf("Expected %t, got %t at index %d", true, s.IsNull(i), i)
		}
	}

	// Append NAs
	s = NewSeriesInt64([]int64{}, nil, true, ctx)
	na := NewSeriesNA(10, ctx)

	for i := 0; i < 100; i++ {
		s = s.Append(na).(Int64s)
		if !utils.CheckEqSlice(s.GetNullMask()[s.Len()-10:], na.GetNullMask(), nil, "Int64s.Append") {
			t.Errorf("Expected %v, got %v at index %d", na.GetNullMask(), s.GetNullMask()[s.Len()-10:], i)
		}
	}

	// Append NullableInt64
	s = NewSeriesInt64([]int64{}, nil, true, ctx)

	for i := 0; i < 100; i++ {
		s = s.Append(gandalff.NullableInt64{Valid: false, Value: 1}).(Int64s)
		if !s.IsNull(i) {
			t.Errorf("Expected %t, got %t at index %d", true, s.IsNull(i), i)
		}
	}

	// Append []NullableInt64
	s = NewSeriesInt64([]int64{}, nil, true, ctx)

	for i := 0; i < 100; i++ {
		s = s.Append([]gandalff.NullableInt64{{Valid: false, Value: 1}}).(Int64s)
		if !s.IsNull(i) {
			t.Errorf("Expected %t, got %t at index %d", true, s.IsNull(i), i)
		}
	}

	// Append Int64s
	s = NewSeriesInt64([]int64{}, nil, true, ctx)
	b := NewSeriesInt64(dataA, []bool{true, true, true, true, true, true, true, true, true, true}, true, ctx)

	for i := 0; i < 100; i++ {
		s = s.Append(b).(Int64s)
		if !utils.CheckEqSlice(s.GetNullMask()[s.Len()-10:], b.GetNullMask(), nil, "Int64s.Append") {
			t.Errorf("Expected %v, got %v at index %d", b.GetNullMask(), s.GetNullMask()[s.Len()-10:], i)
		}
	}
}

func Test_SeriesInt64_Cast(t *testing.T) {
	data := []int64{0, 1, 0, 3, 4, 5, 6, 7, 8, 9}
	mask := []bool{false, false, true, false, false, true, false, false, true, false}

	// Create a new series.
	s := NewSeriesInt64(data, mask, true, ctx)

	// Cast to bool.
	result := s.Cast(meta.BoolType)

	// Check the data.
	expected := []bool{false, true, false, true, true, true, true, true, true, true}
	for i, v := range result.Data().([]bool) {
		if v != expected[i] {
			t.Errorf("Expected %t, got %t at index %d", expected[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if v != mask[i] {
			t.Errorf("Expected nullMask of %t, got %t at index %d", mask[i], v, i)
		}
	}

	// Cast to int.
	result = s.Cast(meta.IntType)

	// Check the data.
	expectedInt := []int{0, 1, 0, 3, 4, 5, 6, 7, 8, 9}
	for i, v := range result.Data().([]int) {
		if v != expectedInt[i] {
			t.Errorf("Expected %d, got %d at index %d", expectedInt[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if v != mask[i] {
			t.Errorf("Expected nullMask of %t, got %t at index %d", mask[i], v, i)
		}
	}

	// Cast to float64.
	result = s.Cast(meta.Float64Type)

	// Check the data.
	expectedFloat := []float64{0, 1, 0, 3, 4, 5, 6, 7, 8, 9}
	for i, v := range result.Data().([]float64) {
		if v != expectedFloat[i] {
			t.Errorf("Expected %f, got %f at index %d", expectedFloat[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if v != mask[i] {
			t.Errorf("Expected nullMask of %t, got %t at index %d", mask[i], v, i)
		}
	}

	// Cast to string.
	result = s.Cast(meta.StringType)

	// Check the data.
	expectedString := []string{"0", "1", gandalff.NA_TEXT, "3", "4", gandalff.NA_TEXT, "6", "7", gandalff.NA_TEXT, "9"}

	for i, v := range result.Data().([]string) {
		if v != expectedString[i] {
			t.Errorf("Expected %s, got %s at index %d", expectedString[i], v, i)
		}
	}

	// Check the null mask.
	for i, v := range result.GetNullMask() {
		if v != mask[i] {
			t.Errorf("Expected nullMask of %t, got %t at index %d", mask[i], v, i)
		}
	}

	// Cast to error.
	castError := s.Cast(meta.ErrorType)

	// Check the message.
	if castError.(Errors).Msg_ != "Int64s.Cast: invalid type Error" {
		t.Errorf("Expected error, got %v", castError)
	}
}

func Test_SeriesInt64_Filter(t *testing.T) {
	data := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	mask := []bool{false, true, false, false, true, false, false, true, false, false, true, false, false, true, false, false, true, false, false, true}

	// Create a new series.
	s := NewSeriesInt64(data, mask, true, ctx)

	// Filter mask.
	filterMask := []bool{true, false, true, true, false, true, true, false, true, true, true, false, true, true, false, true, true, false, true, true}
	filterIndeces := []int{0, 2, 3, 5, 6, 8, 9, 10, 12, 13, 15, 16, 18, 19}

	result := []int64{1, 3, 4, 6, 7, 9, 10, 11, 13, 14, 16, 17, 19, 20}
	resultMask := []bool{false, false, false, false, false, false, false, true, false, true, false, true, false, true}

	/////////////////////////////////////////////////////////////////////////////////////
	// 							Check the Filter() method.
	filtered := s.Filter(NewSeriesBool(filterMask, nil, true, ctx))

	// Check the length.
	if filtered.Len() != len(result) {
		t.Errorf("Expected length of %d, got %d", len(result), filtered.Len())
	}

	// Check the data.
	for i, v := range filtered.Data().([]int64) {
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
	for i, v := range filtered.Data().([]int64) {
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
	for i, v := range filtered.Data().([]int64) {
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

	if e, ok := filtered.(Errors); !ok || e.GetError() != "Int64s.Filter: mask length (20) does not match series length (14)" {
		t.Errorf("Expected Errors, got %v", filtered)
	}

	// Another test.
	data = []int64{2, 323, 42, 4, 9, 674, 42, 48, 9811, 79, 3, 12, 492, 47005, -173, -28, 323, 42, 4, 9, 31, 425, 2}
	mask = []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}

	// Create a new series.
	s = NewSeriesInt64(data, mask, true, ctx)

	// Filter mask.
	filterMask = []bool{true, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false, false, false, false, true}
	filterIndeces = []int{0, 15, 22}

	result = []int64{2, -28, 2}

	/////////////////////////////////////////////////////////////////////////////////////
	// 							Check the Filter() method.
	filtered = s.Filter(filterMask)

	// Check the length.
	if filtered.Len() != 3 {
		t.Errorf("Expected length of 3, got %d", filtered.Len())
	}

	// Check the data.
	for i, v := range filtered.Data().([]int64) {
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
	for i, v := range filtered.Data().([]int64) {
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

func Test_SeriesInt64_Map(t *testing.T) {
	data := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, -2, 323, 24, -23, 4, 42, 5, -6, 7}
	nullMask := []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, false, true, true, true, false, true}

	// Create a new series.
	s := NewSeriesInt64(data, nullMask, true, ctx)

	// Map the series to bool.
	resBool := s.Map(func(v any) any {
		if v.(int64) >= 7 && v.(int64) <= 100 {
			return true
		}
		return false
	})

	expectedBool := []bool{false, false, false, false, false, false, true, true, true, true, false, false, true, false, false, true, false, false, true}
	for i, v := range resBool.Data().([]bool) {
		if v != expectedBool[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedBool[i], v, i)
		}
	}

	// Map the series to int.
	resInt := s.Map(func(v any) any {
		if v.(int64) < 0 {
			return int(-(v.(int64)) % 7)
		}
		return int(v.(int64) % 7)
	})

	expectedInt := []int{1, 2, 3, 4, 5, 6, 0, 1, 2, 3, 2, 1, 3, 2, 4, 0, 5, 6, 0}
	for i, v := range resInt.Data().([]int) {
		if v != expectedInt[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedInt[i], v, i)
		}
	}

	// Map the series to float64.
	resFloat64 := s.Map(func(v any) any {
		if v.(int64) >= 0 {
			return float64(-v.(int64))
		}
		return float64(v.(int64))
	})

	expectedFloat64 := []float64{-1, -2, -3, -4, -5, -6, -7, -8, -9, -10, -2, -323, -24, -23, -4, -42, -5, -6, -7}
	for i, v := range resFloat64.Data().([]float64) {
		if v != expectedFloat64[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedFloat64[i], v, i)
		}
	}

	// Map the series to string.
	resString := s.Map(func(v any) any {
		if v.(int64) >= 0 {
			return "pos"
		}
		return "neg"
	})

	expectedString := []string{"pos", "pos", "pos", "pos", "pos", "pos", "pos", "pos", "pos", "pos", "neg", "pos", "pos", "neg", "pos", "pos", "pos", "neg", "pos"}
	for i, v := range resString.Data().([]string) {
		if v != expectedString[i] {
			t.Errorf("Expected %v, got %v at index %d", expectedString[i], v, i)
		}
	}
}

func Test_SeriesInt64_Group(t *testing.T) {
	var partMap map[int64][]int

	data1 := []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	data1Mask := []bool{false, false, false, false, false, true, true, true, true, true}
	data2 := []int64{1, 1, 2, 2, 1, 1, 2, 2, 1, 1}
	data2Mask := []bool{false, true, false, true, false, true, false, true, false, true}
	data3 := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	data3Mask := []bool{false, false, false, false, false, true, true, true, true, true}

	// Test 1
	s1 := NewSeriesInt64(data1, data1Mask, true, ctx).
		Group()

	p1 := s1.GetPartition().GetMap()
	if len(p1) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(p1))
	}

	partMap = map[int64][]int{
		0: {0, 1, 2, 3, 4},
		1: {5, 6, 7, 8, 9},
	}
	if !utils.CheckEqPartitionMap(p1, partMap, nil, "Int64 Group") {
		t.Errorf("Expected partition map of %v, got %v", partMap, p1)
	}

	// Test 2
	s2 := NewSeriesInt64(data2, data2Mask, true, ctx).
		GroupBy(s1.GetPartition())

	p2 := s2.GetPartition().GetMap()
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
	if !utils.CheckEqPartitionMap(p2, partMap, nil, "Int64 Group") {
		t.Errorf("Expected partition map of %v, got %v", partMap, p2)
	}

	// Test 3
	s3 := NewSeriesInt64(data3, data3Mask, true, ctx).
		GroupBy(s2.GetPartition())

	p3 := s3.GetPartition().GetMap()
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
	if !utils.CheckEqPartitionMap(p3, partMap, nil, "Int64 Group") {
		t.Errorf("Expected partition map of %v, got %v", partMap, p3)
	}

	// debugPrintPartition(s1.GetPartition(), s1)
	// debugPrintPartition(s2.GetPartition(), s1, s2)
	// debugPrintPartition(s3.GetPartition(), s1, s2, s3)

	partMap = nil
}

func Test_SeriesInt64_Sort(t *testing.T) {
	data := []int64{821, 258, -547, -624, 337, -909, -715, 317, -827, -103, 271, 159, 230, -346, 471, 897, 801, 492, 45, -70}
	mask := []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}

	// Create a new series.
	s := NewSeriesInt64(data, nil, true, ctx)

	// Sort the series.
	sorted := s.Sort()

	// Check the data.
	expected := []int64{-909, -827, -715, -624, -547, -346, -103, -70, 45, 159, 230, 258, 271, 317, 337, 471, 492, 801, 821, 897}
	if !utils.CheckEqSliceInt64(sorted.Data().([]int64), expected, nil, "") {
		t.Errorf("Int64s.Sort() failed, expecting %v, got %v", expected, sorted.Data().([]int64))
	}

	// Create a new series.
	s = NewSeriesInt64(data, mask, true, ctx)

	// Sort the series.
	sorted = s.Sort()

	// Check the data.
	expected = []int64{-827, -715, -547, 45, 230, 271, 337, 471, 801, 821, -909, 159, -103, -346, 317, 897, -624, 492, 258, -70}
	if !utils.CheckEqSliceInt64(sorted.Data().([]int64), expected, nil, "") {
		t.Errorf("Int64s.Sort() failed, expecting %v, got %v", expected, sorted.Data().([]int64))
	}

	// Check the null mask.
	expectedMask := []bool{false, false, false, false, false, false, false, false, false, false, true, true, true, true, true, true, true, true, true, true}
	if !utils.CheckEqSliceBool(sorted.GetNullMask(), expectedMask, nil, "") {
		t.Errorf("Int64s.Sort() failed, expecting %v, got %v", expectedMask, sorted.GetNullMask())
	}
}

func Test_SeriesInt64_Arithmetic_Mul(t *testing.T) {
	bools := NewSeriesBool([]bool{true}, nil, true, ctx)
	boolv := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx)
	bools_ := NewSeriesBool([]bool{true}, nil, true, ctx).SetNullMask([]bool{true})
	boolv_ := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i32s := NewSeriesInt([]int{2}, nil, true, ctx)
	i32v := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i32s_ := NewSeriesInt([]int{2}, nil, true, ctx).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i64s := NewSeriesInt64([]int64{2}, nil, true, ctx)
	i64v := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i64s_ := NewSeriesInt64([]int64{2}, nil, true, ctx).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	f64s := NewSeriesFloat64([]float64{2}, nil, true, ctx)
	f64v := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	f64s_ := NewSeriesFloat64([]float64{2}, nil, true, ctx).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	// scalar | bool
	if !utils.CheckEqSlice(i64s.Mul(bools).Data().([]int64), []int64{2}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(boolv).Data().([]int64), []int64{2, 0, 2, 0, 2, 0, 2, 2, 0, 0}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(bools_).GetNullMask(), []bool{true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}

	// scalar | int
	if !utils.CheckEqSlice(i64s.Mul(i32s).Data().([]int64), []int64{4}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(i32v).Data().([]int64), []int64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(i32s_).GetNullMask(), []bool{true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}

	// scalar | int64
	if !utils.CheckEqSlice(i64s.Mul(i64s).Data().([]int64), []int64{4}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(i64v).Data().([]int64), []int64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(i64s_).GetNullMask(), []bool{true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}

	// scalar | float64
	if !utils.CheckEqSlice(i64s.Mul(f64s).Data().([]float64), []float64{4}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(f64v).Data().([]float64), []float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(f64s_).GetNullMask(), []bool{true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64s.Mul(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}

	// vector | bool
	if !utils.CheckEqSlice(i64v.Mul(bools).Data().([]int64), []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(boolv).Data().([]int64), []int64{1, 0, 3, 0, 5, 0, 7, 8, 0, 0}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(bools_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}

	// vector | int
	if !utils.CheckEqSlice(i64v.Mul(i32s).Data().([]int64), []int64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(i32v).Data().([]int64), []int64{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(i32s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}

	// vector | int64
	if !utils.CheckEqSlice(i64v.Mul(i64s).Data().([]int64), []int64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(i64v).Data().([]int64), []int64{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(i64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}

	// vector | float64
	if !utils.CheckEqSlice(i64v.Mul(f64s).Data().([]float64), []float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(f64v).Data().([]float64), []float64{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(f64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
	if !utils.CheckEqSlice(i64v.Mul(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mul") {
		t.Errorf("Error in Int64 Mul")
	}
}

func Test_SeriesInt64_Arithmetic_Div(t *testing.T) {
	bools := NewSeriesBool([]bool{true}, nil, true, ctx)
	boolv := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx)
	bools_ := NewSeriesBool([]bool{true}, nil, true, ctx).SetNullMask([]bool{true})
	boolv_ := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i32s := NewSeriesInt([]int{2}, nil, true, ctx)
	i32v := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i32s_ := NewSeriesInt([]int{2}, nil, true, ctx).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i64s := NewSeriesInt64([]int64{2}, nil, true, ctx)
	i64v := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i64s_ := NewSeriesInt64([]int64{2}, nil, true, ctx).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	f64s := NewSeriesFloat64([]float64{2}, nil, true, ctx)
	f64v := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	f64s_ := NewSeriesFloat64([]float64{2}, nil, true, ctx).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	// scalar | bool
	if !utils.CheckEqSlice(i64s.Div(bools).Data().([]float64), []float64{2}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(boolv).Data().([]float64), []float64{2, math.Inf(1), 2, math.Inf(1), 2, math.Inf(1), 2, 2, math.Inf(1), math.Inf(1)}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(bools_).GetNullMask(), []bool{true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}

	// scalar | int
	if !utils.CheckEqSlice(i64s.Div(i32s).Data().([]float64), []float64{1}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(i32v).Data().([]float64), []float64{2, 1, 0.6666666666666666, 0.5, 0.4, 0.3333333333333333, 0.2857142857142857, 0.25, 0.2222222222222222, 0.2}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(i32s_).GetNullMask(), []bool{true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}

	// scalar | int64
	if !utils.CheckEqSlice(i64s.Div(i64s).Data().([]float64), []float64{1}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(i64v).Data().([]float64), []float64{2, 1, 0.6666666666666666, 0.5, 0.4, 0.3333333333333333, 0.2857142857142857, 0.25, 0.2222222222222222, 0.2}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(i64s_).GetNullMask(), []bool{true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}

	// scalar | float64
	if !utils.CheckEqSlice(i64s.Div(f64s).Data().([]float64), []float64{1}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(f64v).Data().([]float64), []float64{2, 1, 0.6666666666666666, 0.5, 0.4, 0.3333333333333333, 0.2857142857142857, 0.25, 0.2222222222222222, 0.2}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(f64s_).GetNullMask(), []bool{true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64s.Div(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}

	// vector | bool
	if !utils.CheckEqSlice(i64v.Div(bools).Data().([]float64), []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(boolv).Data().([]float64), []float64{1, math.Inf(1), 3, math.Inf(1), 5, math.Inf(1), 7, 8, math.Inf(1), math.Inf(1)}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(bools_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}

	// vector | int
	if !utils.CheckEqSlice(i64v.Div(i32s).Data().([]float64), []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(i32v).Data().([]float64), []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(i32s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}

	// vector | int64
	if !utils.CheckEqSlice(i64v.Div(i64s).Data().([]float64), []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(i64v).Data().([]float64), []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(i64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}

	// vector | float64
	if !utils.CheckEqSlice(i64v.Div(f64s).Data().([]float64), []float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(f64v).Data().([]float64), []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(f64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
	if !utils.CheckEqSlice(i64v.Div(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int Div") {
		t.Errorf("Error in Int Div")
	}
}

func Test_SeriesInt64_Arithmetic_Mod(t *testing.T) {
	bools := NewSeriesBool([]bool{true}, nil, true, ctx)
	boolv := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx)
	bools_ := NewSeriesBool([]bool{true}, nil, true, ctx).SetNullMask([]bool{true})
	boolv_ := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i32s := NewSeriesInt([]int{2}, nil, true, ctx)
	i32v := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i32s_ := NewSeriesInt([]int{2}, nil, true, ctx).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i64s := NewSeriesInt64([]int64{2}, nil, true, ctx)
	i64v := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i64s_ := NewSeriesInt64([]int64{2}, nil, true, ctx).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	f64s := NewSeriesFloat64([]float64{2}, nil, true, ctx)
	f64v := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	f64s_ := NewSeriesFloat64([]float64{2}, nil, true, ctx).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	// scalar | bool
	if !utils.CheckEqSlice(i64s.Mod(bools).Data().([]float64), []float64{0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(boolv).Data().([]float64), []float64{0, math.NaN(), 0, math.NaN(), 0, math.NaN(), 0, 0, math.NaN(), math.NaN()}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(bools_).GetNullMask(), []bool{true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}

	// scalar | int
	if !utils.CheckEqSlice(i64s.Mod(i32s).Data().([]float64), []float64{0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(i32v).Data().([]float64), []float64{0, 0, 2, 2, 2, 2, 2, 2, 2, 2}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(i32s_).GetNullMask(), []bool{true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}

	// scalar | int64
	if !utils.CheckEqSlice(i64s.Mod(i64s).Data().([]float64), []float64{0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(i64v).Data().([]float64), []float64{0, 0, 2, 2, 2, 2, 2, 2, 2, 2}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(i64s_).GetNullMask(), []bool{true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}

	// scalar | float64
	if !utils.CheckEqSlice(i64s.Mod(f64s).Data().([]float64), []float64{0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(f64v).Data().([]float64), []float64{0, 0, 2, 2, 2, 2, 2, 2, 2, 2}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(f64s_).GetNullMask(), []bool{true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64s.Mod(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}

	// vector | bool
	if !utils.CheckEqSlice(i64v.Mod(bools).Data().([]float64), []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(boolv).Data().([]float64), []float64{0, math.NaN(), 0, math.NaN(), 0, math.NaN(), 0, 0, math.NaN(), math.NaN()}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(bools_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}

	// vector | int
	if !utils.CheckEqSlice(i64v.Mod(i32s).Data().([]float64), []float64{1, 0, 1, 0, 1, 0, 1, 0, 1, 0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(i32v).Data().([]float64), []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(i32s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}

	// vector | int64
	if !utils.CheckEqSlice(i64v.Mod(i64s).Data().([]float64), []float64{1, 0, 1, 0, 1, 0, 1, 0, 1, 0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(i64v).Data().([]float64), []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(i64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}

	// vector | float64
	if !utils.CheckEqSlice(i64v.Mod(f64s).Data().([]float64), []float64{1, 0, 1, 0, 1, 0, 1, 0, 1, 0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(f64v).Data().([]float64), []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(f64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
	if !utils.CheckEqSlice(i64v.Mod(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Mod") {
		t.Errorf("Error in Int64 Mod")
	}
}

func Test_SeriesInt64_Arithmetic_Exp(t *testing.T) {
	bools := NewSeriesBool([]bool{true}, nil, true, ctx)
	boolv := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx)
	bools_ := NewSeriesBool([]bool{true}, nil, true, ctx).SetNullMask([]bool{true})
	boolv_ := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i32s := NewSeriesInt([]int{2}, nil, true, ctx)
	i32v := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i32s_ := NewSeriesInt([]int{2}, nil, true, ctx).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i64s := NewSeriesInt64([]int64{2}, nil, true, ctx)
	i64v := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i64s_ := NewSeriesInt64([]int64{2}, nil, true, ctx).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	f64s := NewSeriesFloat64([]float64{2}, nil, true, ctx)
	f64v := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	f64s_ := NewSeriesFloat64([]float64{2}, nil, true, ctx).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	// scalar | bool
	if !utils.CheckEqSlice(i64s.Exp(bools).Data().([]int64), []int64{2}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(boolv).Data().([]int64), []int64{2, 1, 2, 1, 2, 1, 2, 2, 1, 1}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(bools_).GetNullMask(), []bool{true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}

	// scalar | int
	if !utils.CheckEqSlice(i64s.Exp(i32s).Data().([]int64), []int64{4}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(i32v).Data().([]int64), []int64{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(i32s_).GetNullMask(), []bool{true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}

	// scalar | int64
	if !utils.CheckEqSlice(i64s.Exp(i64s).Data().([]int64), []int64{4}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(i64v).Data().([]int64), []int64{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(i64s_).GetNullMask(), []bool{true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}

	// scalar | float64
	if !utils.CheckEqSlice(i64s.Exp(f64s).Data().([]float64), []float64{4}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(f64v).Data().([]float64), []float64{2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(f64s_).GetNullMask(), []bool{true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64s.Exp(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}

	// vector | bool
	if !utils.CheckEqSlice(i64v.Exp(bools).Data().([]int64), []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(boolv).Data().([]int64), []int64{1, 1, 3, 1, 5, 1, 7, 8, 1, 1}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(bools_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}

	// vector | int
	if !utils.CheckEqSlice(i64v.Exp(i32s).Data().([]int64), []int64{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(i32v).Data().([]int64), []int64{1, 4, 27, 256, 3125, 46656, 823543, 16777216, 387420489, 10000000000}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(i32s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}

	// vector | int64
	if !utils.CheckEqSlice(i64v.Exp(i64s).Data().([]int64), []int64{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(i64v).Data().([]int64), []int64{1, 4, 27, 256, 3125, 46656, 823543, 16777216, 387420489, 10000000000}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(i64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}

	// vector | float64
	if !utils.CheckEqSlice(i64v.Exp(f64s).Data().([]float64), []float64{1, 4, 9, 16, 25, 36, 49, 64, 81, 100}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(f64v).Data().([]float64), []float64{1, 4, 27, 256, 3125, 46656, 823543, 16777216, 387420489, 10000000000}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(f64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
	if !utils.CheckEqSlice(i64v.Exp(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Exp") {
		t.Errorf("Error in Int64 Exp")
	}
}

func Test_SeriesInt64_Arithmetic_Add(t *testing.T) {
	bools := NewSeriesBool([]bool{true}, nil, true, ctx)
	boolv := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx)
	bools_ := NewSeriesBool([]bool{true}, nil, true, ctx).SetNullMask([]bool{true})
	boolv_ := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i32s := NewSeriesInt([]int{2}, nil, true, ctx)
	i32v := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i32s_ := NewSeriesInt([]int{2}, nil, true, ctx).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i64s := NewSeriesInt64([]int64{2}, nil, true, ctx)
	i64v := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i64s_ := NewSeriesInt64([]int64{2}, nil, true, ctx).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	f64s := NewSeriesFloat64([]float64{2}, nil, true, ctx)
	f64v := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	f64s_ := NewSeriesFloat64([]float64{2}, nil, true, ctx).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	ss := NewSeriesString([]string{"2"}, nil, true, ctx)
	sv := NewSeriesString([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, nil, true, ctx)
	ss_ := NewSeriesString([]string{"2"}, nil, true, ctx).SetNullMask([]bool{true})
	sv_ := NewSeriesString([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	// scalar | bool
	if !utils.CheckEqSlice(i64s.Add(bools).Data().([]int64), []int64{3}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(boolv).Data().([]int64), []int64{3, 2, 3, 2, 3, 2, 3, 3, 2, 2}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(bools_).GetNullMask(), []bool{true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// scalar | int
	if !utils.CheckEqSlice(i64s.Add(i32s).Data().([]int64), []int64{4}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(i32v).Data().([]int64), []int64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(i32s_).GetNullMask(), []bool{true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// scalar | int64
	if !utils.CheckEqSlice(i64s.Add(i64s).Data().([]int64), []int64{4}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(i64v).Data().([]int64), []int64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(i64s_).GetNullMask(), []bool{true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// scalar | float64
	if !utils.CheckEqSlice(i64s.Add(f64s).Data().([]float64), []float64{4}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(f64v).Data().([]float64), []float64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(f64s_).GetNullMask(), []bool{true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// scalar | string
	if !utils.CheckEqSlice(i64s.Add(ss).Data().([]string), []string{"22"}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(sv).Data().([]string), []string{"21", "22", "23", "24", "25", "26", "27", "28", "29", "210"}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(ss_).GetNullMask(), []bool{true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64s.Add(sv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// vector | bool
	if !utils.CheckEqSlice(i64v.Add(bools).Data().([]int64), []int64{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(boolv).Data().([]int64), []int64{2, 2, 4, 4, 6, 6, 8, 9, 9, 10}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(bools_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// vector | int
	if !utils.CheckEqSlice(i64v.Add(i32s).Data().([]int64), []int64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(i32v).Data().([]int64), []int64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(i32s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// vector | int64
	if !utils.CheckEqSlice(i64v.Add(i64s).Data().([]int64), []int64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(i64v).Data().([]int64), []int64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(i64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// vector | float64
	if !utils.CheckEqSlice(i64v.Add(f64s).Data().([]float64), []float64{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(f64v).Data().([]float64), []float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(f64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}

	// vector | string
	if !utils.CheckEqSlice(i64v.Add(ss).Data().([]string), []string{"12", "22", "32", "42", "52", "62", "72", "82", "92", "102"}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(sv).Data().([]string), []string{"11", "22", "33", "44", "55", "66", "77", "88", "99", "1010"}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(ss_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
	if !utils.CheckEqSlice(i64v.Add(sv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Add") {
		t.Errorf("Error in Int64 Add")
	}
}

func Test_SeriesInt64_Arithmetic_Sub(t *testing.T) {
	bools := NewSeriesBool([]bool{true}, nil, true, ctx)
	boolv := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx)
	bools_ := NewSeriesBool([]bool{true}, nil, true, ctx).SetNullMask([]bool{true})
	boolv_ := NewSeriesBool([]bool{true, false, true, false, true, false, true, true, false, false}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i32s := NewSeriesInt([]int{2}, nil, true, ctx)
	i32v := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i32s_ := NewSeriesInt([]int{2}, nil, true, ctx).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	i64s := NewSeriesInt64([]int64{2}, nil, true, ctx)
	i64v := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	i64s_ := NewSeriesInt64([]int64{2}, nil, true, ctx).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	f64s := NewSeriesFloat64([]float64{2}, nil, true, ctx)
	f64v := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx)
	f64s_ := NewSeriesFloat64([]float64{2}, nil, true, ctx).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil, true, ctx).
		SetNullMask([]bool{false, true, false, true, false, true, false, true, false, true})

	// scalar | bool
	if !utils.CheckEqSlice(i64s.Sub(bools).Data().([]int64), []int64{1}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(boolv).Data().([]int64), []int64{1, 2, 1, 2, 1, 2, 1, 1, 2, 2}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(bools_).GetNullMask(), []bool{true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}

	// scalar | int
	if !utils.CheckEqSlice(i64s.Sub(i32s).Data().([]int64), []int64{0}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(i32v).Data().([]int64), []int64{1, 0, -1, -2, -3, -4, -5, -6, -7, -8}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(i32s_).GetNullMask(), []bool{true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}

	// scalar | int64
	if !utils.CheckEqSlice(i64s.Sub(i64s).Data().([]int64), []int64{0}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(i64v).Data().([]int64), []int64{1, 0, -1, -2, -3, -4, -5, -6, -7, -8}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(i64s_).GetNullMask(), []bool{true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}

	// scalar | float64
	if !utils.CheckEqSlice(i64s.Sub(f64s).Data().([]float64), []float64{0}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(f64v).Data().([]float64), []float64{1, 0, -1, -2, -3, -4, -5, -6, -7, -8}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(f64s_).GetNullMask(), []bool{true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64s.Sub(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}

	// vector | bool
	if !utils.CheckEqSlice(i64v.Sub(bools).Data().([]int64), []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(boolv).Data().([]int64), []int64{0, 2, 2, 4, 4, 6, 6, 7, 9, 10}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(bools_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(boolv_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}

	// vector | int
	if !utils.CheckEqSlice(i64v.Sub(i32s).Data().([]int64), []int64{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(i32v).Data().([]int64), []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(i32s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(i32v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}

	// vector | int64
	if !utils.CheckEqSlice(i64v.Sub(i64s).Data().([]int64), []int64{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(i64v).Data().([]int64), []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(i64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(i64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}

	// vector | float64
	if !utils.CheckEqSlice(i64v.Sub(f64s).Data().([]float64), []float64{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(f64v).Data().([]float64), []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(f64s_).GetNullMask(), []bool{true, true, true, true, true, true, true, true, true, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
	if !utils.CheckEqSlice(i64v.Sub(f64v_).GetNullMask(), []bool{false, true, false, true, false, true, false, true, false, true}, nil, "Int64 Sub") {
		t.Errorf("Error in Int64 Sub")
	}
}

func Test_SeriesInt64_Logical_Eq(t *testing.T) {
	// TODO: add tests for all types
}

func Test_SeriesInt64_Logical_Ne(t *testing.T) {
	// TODO: add tests for all types
}

func Test_SeriesInt64_Logical_Lt(t *testing.T) {
	var res Series

	i32s := NewSeriesInt([]int{1}, nil, true, ctx)
	i32v := NewSeriesInt([]int{1, 2, 3}, nil, true, ctx)
	i32s_ := NewSeriesInt([]int{1}, nil, true, ctx).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt([]int{1, 2, 3}, nil, true, ctx).SetNullMask([]bool{true, true, false})

	i64s := NewSeriesInt64([]int64{1}, nil, true, ctx)
	i64v := NewSeriesInt64([]int64{1, 2, 3}, nil, true, ctx)
	i64s_ := NewSeriesInt64([]int64{1}, nil, true, ctx).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3}, nil, true, ctx).SetNullMask([]bool{true, true, false})

	f64s := NewSeriesFloat64([]float64{1.0}, nil, true, ctx)
	f64v := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, true, ctx)
	f64s_ := NewSeriesFloat64([]float64{1.0}, nil, true, ctx).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, true, ctx).SetNullMask([]bool{true, true, false})

	// scalar | int
	res = i64s.Lt(i32s)
	if res.Data().([]bool)[0] != false {
		t.Errorf("Expected %v, got %v", []bool{false}, res.Data().([]bool))
	}

	res = i64s.Lt(i32s_)
	if res.IsNull(0) == false {
		t.Errorf("Expected %v, got %v", []bool{true}, res.GetNullMask())
	}

	res = i64s.Lt(i32v)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{false, true, true}, res.Data().([]bool))
	}

	res = i64s.Lt(i32v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// scalar | int64
	res = i64s.Lt(i64s)
	if res.Data().([]bool)[0] != false {
		t.Errorf("Expected %v, got %v", []bool{false}, res.Data().([]bool))
	}

	res = i64s.Lt(i64s_)
	if res.IsNull(0) == false {
		t.Errorf("Expected %v, got %v", []bool{true}, res.GetNullMask())
	}

	res = i64s.Lt(i64v)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{false, true, true}, res.Data().([]bool))
	}

	res = i64s.Lt(i64v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// scalar | float64
	res = i64s.Lt(f64s)
	if res.Data().([]bool)[0] != false {
		t.Errorf("Expected %v, got %v", []bool{false}, res.Data().([]bool))
	}

	res = i64s.Lt(f64s_)
	if res.IsNull(0) == false {
		t.Errorf("Expected %v, got %v", []bool{true}, res.GetNullMask())
	}

	res = i64s.Lt(f64v)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{false, true, true}, res.Data().([]bool))
	}

	res = i64s.Lt(f64v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// vector | int
	res = i64v.Lt(i32s)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.Data().([]bool))
	}

	res = i64v.Lt(i32s_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == false {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.GetNullMask())
	}

	res = i64v.Lt(i32v)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.Data().([]bool))
	}

	res = i64v.Lt(i32v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// vector | int64
	res = i64v.Lt(i64s)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.Data().([]bool))
	}

	res = i64v.Lt(i64s_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == false {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.GetNullMask())
	}

	res = i64v.Lt(i64v)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.Data().([]bool))
	}

	res = i64v.Lt(i64v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// vector | float64
	res = i64v.Lt(f64s)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.Data().([]bool))
	}

	res = i64v.Lt(f64s_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == false {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.GetNullMask())
	}

	res = i64v.Lt(f64v)
	if res.Data().([]bool)[0] != false || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.Data().([]bool))
	}

	res = i64v.Lt(f64v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}
}

func Test_SeriesInt64_Logical_Le(t *testing.T) {
	var res Series

	i32s := NewSeriesInt([]int{1}, nil, true, ctx)
	i32v := NewSeriesInt([]int{1, 2, 3}, nil, true, ctx)
	i32s_ := NewSeriesInt([]int{1}, nil, true, ctx).SetNullMask([]bool{true})
	i32v_ := NewSeriesInt([]int{1, 2, 3}, nil, true, ctx).SetNullMask([]bool{true, true, false})

	i64s := NewSeriesInt64([]int64{1}, nil, true, ctx)
	i64v := NewSeriesInt64([]int64{1, 2, 3}, nil, true, ctx)
	i64s_ := NewSeriesInt64([]int64{1}, nil, true, ctx).SetNullMask([]bool{true})
	i64v_ := NewSeriesInt64([]int64{1, 2, 3}, nil, true, ctx).SetNullMask([]bool{true, true, false})

	f64s := NewSeriesFloat64([]float64{1.0}, nil, true, ctx)
	f64v := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, true, ctx)
	f64s_ := NewSeriesFloat64([]float64{1.0}, nil, true, ctx).SetNullMask([]bool{true})
	f64v_ := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, true, ctx).SetNullMask([]bool{true, true, false})

	// scalar | int
	res = i64s.Le(i32s)
	if res.Data().([]bool)[0] != true {
		t.Errorf("Expected %v, got %v", []bool{true}, res.Data().([]bool))
	}

	res = i64s.Le(i32s_)
	if res.IsNull(0) == false {
		t.Errorf("Expected %v, got %v", []bool{true}, res.GetNullMask())
	}

	res = i64s.Le(i32v)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.Data().([]bool))
	}

	res = i64s.Le(i32v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// scalar | int64
	res = i64s.Le(i64s)
	if res.Data().([]bool)[0] != true {
		t.Errorf("Expected %v, got %v", []bool{true}, res.Data().([]bool))
	}

	res = i64s.Le(i64s_)
	if res.IsNull(0) == false {
		t.Errorf("Expected %v, got %v", []bool{true}, res.GetNullMask())
	}

	res = i64s.Le(i64v)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.Data().([]bool))
	}

	res = i64s.Le(i64v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// scalar | float64
	res = i64s.Le(f64s)
	if res.Data().([]bool)[0] != true {
		t.Errorf("Expected %v, got %v", []bool{true}, res.Data().([]bool))
	}

	res = i64s.Le(f64s_)
	if res.IsNull(0) == false {
		t.Errorf("Expected %v, got %v", []bool{true}, res.GetNullMask())
	}

	res = i64s.Le(f64v)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.Data().([]bool))
	}

	res = i64s.Le(f64v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// vector | int
	res = i64v.Le(i32s)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{true, false, false}, res.Data().([]bool))
	}

	res = i64v.Le(i32s_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == false {
		t.Errorf("Expected %v, got %v", []bool{false, false, false}, res.GetNullMask())
	}

	res = i64v.Le(i32v)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.Data().([]bool))
	}

	res = i64v.Le(i32v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// vector | int64
	res = i64v.Le(i64s)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{true, false, false}, res.Data().([]bool))
	}

	res = i64v.Le(i64s_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == false {
		t.Errorf("Expected %v, got %v", []bool{true, false, false}, res.GetNullMask())
	}

	res = i64v.Le(i64v)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.Data().([]bool))
	}

	res = i64v.Le(i64v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}

	// vector | float64
	res = i64v.Le(f64s)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != false || res.Data().([]bool)[2] != false {
		t.Errorf("Expected %v, got %v", []bool{true, false, false}, res.Data().([]bool))
	}

	res = i64v.Le(f64s_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == false {
		t.Errorf("Expected %v, got %v", []bool{true, false, false}, res.GetNullMask())
	}

	res = i64v.Le(f64v)
	if res.Data().([]bool)[0] != true || res.Data().([]bool)[1] != true || res.Data().([]bool)[2] != true {
		t.Errorf("Expected %v, got %v", []bool{true, true, true}, res.Data().([]bool))
	}

	res = i64v.Le(f64v_)
	if res.IsNull(0) == false || res.IsNull(1) == false || res.IsNull(2) == true {
		t.Errorf("Expected %v, got %v", []bool{true, true, false}, res.GetNullMask())
	}
}

func Test_SeriesInt64_Logical_Gt(t *testing.T) {
	// TODO: add tests for all types
}

func Test_SeriesInt64_Logical_Ge(t *testing.T) {
	// TODO: add tests for all types
}
