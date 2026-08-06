package series

import (
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/caerbannogwhite/enchanter"
)

func TestFloat64sArrowArray(t *testing.T) {
	ctx := enchanter.NewContext()
	s := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, false, ctx)
	arr := s.ArrowArray()
	if arr == nil {
		t.Fatal("ArrowArray returned nil")
	}
	if arr.Len() != 3 {
		t.Fatalf("expected len 3, got %d", arr.Len())
	}
	f64arr := arr.(*array.Float64)
	if f64arr.Value(0) != 1.0 || f64arr.Value(1) != 2.0 || f64arr.Value(2) != 3.0 {
		t.Errorf("values mismatch: got %v %v %v", f64arr.Value(0), f64arr.Value(1), f64arr.Value(2))
	}
	if arr.NullN() != 0 {
		t.Errorf("expected 0 nulls, got %d", arr.NullN())
	}
}

func TestFloat64sArrowArrayNullable(t *testing.T) {
	ctx := enchanter.NewContext()
	s := NewSeriesFloat64([]float64{1.0, 0, 3.0}, []bool{false, true, false}, false, ctx)
	arr := s.ArrowArray()
	if arr.Len() != 3 {
		t.Fatalf("expected len 3, got %d", arr.Len())
	}
	if arr.NullN() != 1 {
		t.Fatalf("expected 1 null, got %d", arr.NullN())
	}
	if !arr.IsNull(1) {
		t.Error("index 1 should be null")
	}
	if arr.IsNull(0) || arr.IsNull(2) {
		t.Error("indices 0 and 2 should not be null")
	}
}

func TestInt64sArrowArray(t *testing.T) {
	ctx := enchanter.NewContext()
	s := NewSeriesInt64([]int64{10, 20, 30}, nil, false, ctx)
	arr := s.ArrowArray()
	if arr.Len() != 3 {
		t.Fatalf("expected len 3, got %d", arr.Len())
	}
	i64arr := arr.(*array.Int64)
	if i64arr.Value(0) != 10 || i64arr.Value(1) != 20 || i64arr.Value(2) != 30 {
		t.Error("values mismatch")
	}
}

func TestIntsArrowArray(t *testing.T) {
	ctx := enchanter.NewContext()
	s := NewSeriesInt([]int{100, 200}, nil, false, ctx)
	arr := s.ArrowArray()
	if arr.Len() != 2 {
		t.Fatalf("expected len 2, got %d", arr.Len())
	}
	// Ints are stored as int64 in Arrow
	i64arr := arr.(*array.Int64)
	if i64arr.Value(0) != 100 || i64arr.Value(1) != 200 {
		t.Error("values mismatch")
	}
}

func TestBoolsArrowArray(t *testing.T) {
	ctx := enchanter.NewContext()
	s := NewSeriesBool([]bool{true, false, true}, nil, false, ctx)
	arr := s.ArrowArray()
	if arr.Len() != 3 {
		t.Fatalf("expected len 3, got %d", arr.Len())
	}
	barr := arr.(*array.Boolean)
	if barr.Value(0) != true || barr.Value(1) != false || barr.Value(2) != true {
		t.Error("values mismatch")
	}
}

func TestStringsArrowArray(t *testing.T) {
	ctx := enchanter.NewContext()
	s := NewSeriesString([]string{"hello", "world"}, nil, false, ctx)
	arr := s.ArrowArray()
	if arr.Len() != 2 {
		t.Fatalf("expected len 2, got %d", arr.Len())
	}
	sarr := arr.(*array.String)
	if sarr.Value(0) != "hello" || sarr.Value(1) != "world" {
		t.Error("values mismatch")
	}
}

func TestTimesArrowArray(t *testing.T) {
	ctx := enchanter.NewContext()
	t1 := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 6, 15, 12, 30, 0, 0, time.UTC)
	s := NewSeriesTime([]time.Time{t1, t2}, nil, false, ctx)
	arr := s.ArrowArray()
	if arr.Len() != 2 {
		t.Fatalf("expected len 2, got %d", arr.Len())
	}
	tarr := arr.(*array.Timestamp)
	// Verify nanosecond round-trip
	got1 := time.Unix(0, int64(tarr.Value(0)))
	got2 := time.Unix(0, int64(tarr.Value(1)))
	if !got1.Equal(t1) {
		t.Errorf("time 0: expected %v, got %v", t1, got1)
	}
	if !got2.Equal(t2) {
		t.Errorf("time 1: expected %v, got %v", t2, got2)
	}
}

func TestDurationsArrowArray(t *testing.T) {
	ctx := enchanter.NewContext()
	d1 := 5 * time.Second
	d2 := 100 * time.Millisecond
	s := NewSeriesDuration([]time.Duration{d1, d2}, nil, false, ctx)
	arr := s.ArrowArray()
	if arr.Len() != 2 {
		t.Fatalf("expected len 2, got %d", arr.Len())
	}
	darr := arr.(*array.Duration)
	if time.Duration(darr.Value(0)) != d1 {
		t.Errorf("duration 0: expected %v, got %v", d1, time.Duration(darr.Value(0)))
	}
	if time.Duration(darr.Value(1)) != d2 {
		t.Errorf("duration 1: expected %v, got %v", d2, time.Duration(darr.Value(1)))
	}
}

func TestArrowArrayViaInterface(t *testing.T) {
	ctx := enchanter.NewContext()
	var s Series = NewSeriesFloat64([]float64{1, 2, 3}, nil, false, ctx)
	arr := s.ArrowArray()
	if arr == nil {
		t.Fatal("ArrowArray via Series interface returned nil")
	}
	if arr.Len() != 3 {
		t.Errorf("expected len 3, got %d", arr.Len())
	}
}
