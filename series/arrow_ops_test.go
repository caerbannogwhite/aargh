package series

import (
	"testing"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/meta"
)

func TestArrowAdd(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, false, ctx)
	b := NewSeriesFloat64([]float64{10.0, 20.0, 30.0}, nil, false, ctx)

	result := ArrowAdd(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	if result.Len() != 3 {
		t.Fatalf("expected len 3, got %d", result.Len())
	}
	f := result.(Float64s)
	if f.Data_[0] != 11.0 || f.Data_[1] != 22.0 || f.Data_[2] != 33.0 {
		t.Errorf("add values mismatch: %v", f.Data_)
	}
}

func TestArrowSub(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{10.0, 20.0, 30.0}, nil, false, ctx)
	b := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, false, ctx)

	result := ArrowSub(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	f := result.(Float64s)
	if f.Data_[0] != 9.0 || f.Data_[1] != 18.0 || f.Data_[2] != 27.0 {
		t.Errorf("sub values mismatch: %v", f.Data_)
	}
}

func TestArrowMul(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{2.0, 3.0, 4.0}, nil, false, ctx)
	b := NewSeriesFloat64([]float64{5.0, 6.0, 7.0}, nil, false, ctx)

	result := ArrowMul(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	f := result.(Float64s)
	if f.Data_[0] != 10.0 || f.Data_[1] != 18.0 || f.Data_[2] != 28.0 {
		t.Errorf("mul values mismatch: %v", f.Data_)
	}
}

func TestArrowDiv(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{10.0, 20.0, 30.0}, nil, false, ctx)
	b := NewSeriesFloat64([]float64{2.0, 4.0, 5.0}, nil, false, ctx)

	result := ArrowDiv(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	f := result.(Float64s)
	if f.Data_[0] != 5.0 || f.Data_[1] != 5.0 || f.Data_[2] != 6.0 {
		t.Errorf("div values mismatch: %v", f.Data_)
	}
}

func TestArrowInt64Add(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesInt64([]int64{1, 2, 3}, nil, false, ctx)
	b := NewSeriesInt64([]int64{10, 20, 30}, nil, false, ctx)

	result := ArrowAdd(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	i := result.(Int64s)
	if i.Data_[0] != 11 || i.Data_[1] != 22 || i.Data_[2] != 33 {
		t.Errorf("int64 add values mismatch: %v", i.Data_)
	}
}

func TestArrowEq(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, false, ctx)
	b := NewSeriesFloat64([]float64{1.0, 5.0, 3.0}, nil, false, ctx)

	result := ArrowEq(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	if result.Type() != meta.BoolType {
		t.Fatalf("expected BoolType, got %v", result.Type())
	}
	bools := result.(Bools)
	if bools.Data_[0] != true || bools.Data_[1] != false || bools.Data_[2] != true {
		t.Errorf("eq values mismatch: %v", bools.Data_)
	}
}

func TestArrowLt(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{1.0, 5.0, 3.0}, nil, false, ctx)
	b := NewSeriesFloat64([]float64{2.0, 3.0, 3.0}, nil, false, ctx)

	result := ArrowLt(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	bools := result.(Bools)
	if bools.Data_[0] != true || bools.Data_[1] != false || bools.Data_[2] != false {
		t.Errorf("lt values mismatch: %v", bools.Data_)
	}
}

func TestArrowGe(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{1.0, 5.0, 3.0}, nil, false, ctx)
	b := NewSeriesFloat64([]float64{2.0, 3.0, 3.0}, nil, false, ctx)

	result := ArrowGe(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	bools := result.(Bools)
	if bools.Data_[0] != false || bools.Data_[1] != true || bools.Data_[2] != true {
		t.Errorf("ge values mismatch: %v", bools.Data_)
	}
}

func TestArrowBoolAnd(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesBool([]bool{true, true, false, false}, nil, false, ctx)
	b := NewSeriesBool([]bool{true, false, true, false}, nil, false, ctx)

	result := ArrowAnd(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	bools := result.(Bools)
	if bools.Data_[0] != true || bools.Data_[1] != false || bools.Data_[2] != false || bools.Data_[3] != false {
		t.Errorf("and values mismatch: %v", bools.Data_)
	}
}

func TestArrowBoolOr(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesBool([]bool{true, true, false, false}, nil, false, ctx)
	b := NewSeriesBool([]bool{true, false, true, false}, nil, false, ctx)

	result := ArrowOr(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	bools := result.(Bools)
	if bools.Data_[0] != true || bools.Data_[1] != true || bools.Data_[2] != true || bools.Data_[3] != false {
		t.Errorf("or values mismatch: %v", bools.Data_)
	}
}

func TestArrowNullableAdd(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{1.0, 0, 3.0}, []bool{false, true, false}, false, ctx)
	b := NewSeriesFloat64([]float64{10.0, 20.0, 30.0}, nil, false, ctx)

	result := ArrowAdd(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	if !result.IsNullable() {
		t.Error("result should be nullable")
	}
	if !result.IsNull(1) {
		t.Error("index 1 should be null")
	}
	f := result.(Float64s)
	if f.Data_[0] != 11.0 || f.Data_[2] != 33.0 {
		t.Errorf("nullable add values mismatch: %v", f.Data_)
	}
}

func TestArrowBinaryOpInterface(t *testing.T) {
	ctx := enchanter.NewContext()
	s := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, false, ctx)

	// Test with raw value (will be coerced to series)
	result := ArrowBinaryOp("add", s, []float64{10.0, 20.0, 30.0}, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	f := result.(Float64s)
	if f.Data_[0] != 11.0 || f.Data_[1] != 22.0 || f.Data_[2] != 33.0 {
		t.Errorf("binary op add values mismatch: %v", f.Data_)
	}
}

func TestArrowScalarBroadcast(t *testing.T) {
	ctx := enchanter.NewContext()
	a := NewSeriesFloat64([]float64{1.0, 2.0, 3.0}, nil, false, ctx)
	b := NewSeriesFloat64([]float64{10.0}, nil, false, ctx) // scalar

	result := ArrowMul(a, b, ctx)
	if result.IsError() {
		t.Fatal(result.GetError())
	}
	f := result.(Float64s)
	if f.Data_[0] != 10.0 || f.Data_[1] != 20.0 || f.Data_[2] != 30.0 {
		t.Errorf("scalar broadcast mul values mismatch: %v", f.Data_)
	}
}
