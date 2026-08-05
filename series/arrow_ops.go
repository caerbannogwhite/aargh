package series

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/compute"
	"github.com/apache/arrow-go/v18/arrow/scalar"
	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/meta"
)

// arrowBinaryArith performs a binary arithmetic operation using Arrow compute.
// It handles type promotion and scalar broadcasting automatically.
// Returns a new Series with the result.
func arrowBinaryArith(opName string, left, right Series, ctx *aargh.Context) Series {
	lArr := left.ArrowArray()
	rArr := right.ArrowArray()
	if lArr == nil || rArr == nil {
		return Errors{fmt.Sprintf("arrowBinaryArith(%s): nil Arrow array", opName)}
	}

	var lDatum, rDatum compute.Datum
	if left.Len() == 1 {
		lDatum = scalarDatum(lArr)
	} else {
		lDatum = &compute.ArrayDatum{Value: lArr.Data()}
	}
	if right.Len() == 1 {
		rDatum = scalarDatum(rArr)
	} else {
		rDatum = &compute.ArrayDatum{Value: rArr.Data()}
	}

	var result compute.Datum
	var err error
	opts := compute.ArithmeticOptions{}

	switch opName {
	case "add":
		result, err = compute.Add(context.Background(), opts, lDatum, rDatum)
	case "subtract":
		result, err = compute.Subtract(context.Background(), opts, lDatum, rDatum)
	case "multiply":
		result, err = compute.Multiply(context.Background(), opts, lDatum, rDatum)
	case "divide":
		result, err = compute.Divide(context.Background(), opts, lDatum, rDatum)
	case "power":
		result, err = compute.Power(context.Background(), opts, lDatum, rDatum)
	default:
		return Errors{fmt.Sprintf("arrowBinaryArith: unknown op %q", opName)}
	}

	if err != nil {
		return Errors{fmt.Sprintf("arrowBinaryArith(%s): %v", opName, err)}
	}

	return datumToSeries(result, ctx)
}

// arrowBinaryCompare performs a binary comparison using Arrow compute.
// Always returns a Bool Series.
func arrowBinaryCompare(opName string, left, right Series, ctx *aargh.Context) Series {
	lArr := left.ArrowArray()
	rArr := right.ArrowArray()
	if lArr == nil || rArr == nil {
		return Errors{fmt.Sprintf("arrowBinaryCompare(%s): nil Arrow array", opName)}
	}

	var lDatum, rDatum compute.Datum
	if left.Len() == 1 {
		lDatum = scalarDatum(lArr)
	} else {
		lDatum = &compute.ArrayDatum{Value: lArr.Data()}
	}
	if right.Len() == 1 {
		rDatum = scalarDatum(rArr)
	} else {
		rDatum = &compute.ArrayDatum{Value: rArr.Data()}
	}

	result, err := compute.CallFunction(context.Background(), opName, nil, lDatum, rDatum)
	if err != nil {
		return Errors{fmt.Sprintf("arrowBinaryCompare(%s): %v", opName, err)}
	}

	return datumToSeries(result, ctx)
}

// arrowBooleanOp performs a binary boolean operation (and, or, xor).
func arrowBooleanOp(opName string, left, right Series, ctx *aargh.Context) Series {
	lArr := left.ArrowArray()
	rArr := right.ArrowArray()
	if lArr == nil || rArr == nil {
		return Errors{fmt.Sprintf("arrowBooleanOp(%s): nil Arrow array", opName)}
	}

	lDatum := &compute.ArrayDatum{Value: lArr.Data()}
	rDatum := &compute.ArrayDatum{Value: rArr.Data()}

	result, err := compute.CallFunction(context.Background(), opName, nil, lDatum, rDatum)
	if err != nil {
		return Errors{fmt.Sprintf("arrowBooleanOp(%s): %v", opName, err)}
	}

	return datumToSeries(result, ctx)
}

// scalarDatum extracts the first element of an Arrow array as a scalar Datum.
func scalarDatum(arr arrow.Array) compute.Datum {
	sc, err := scalar.GetScalar(arr, 0)
	if err != nil {
		return &compute.ArrayDatum{Value: arr.Data()}
	}
	return compute.NewDatum(sc)
}

// datumToSeries converts a compute.Datum result back to a Series.
func datumToSeries(d compute.Datum, ctx *aargh.Context) Series {
	switch dt := d.(type) {
	case *compute.ArrayDatum:
		arr := dt.MakeArray()
		return ArrowArrayToSeries(arr, ctx)
	case *compute.ScalarDatum:
		// Wrap scalar into a single-element array
		arr, err := scalar.MakeArrayFromScalar(dt.Value, 1, ctx.Allocator)
		if err != nil {
			return Errors{fmt.Sprintf("datumToSeries: %v", err)}
		}
		return ArrowArrayToSeries(arr, ctx)
	default:
		return Errors{fmt.Sprintf("datumToSeries: unsupported datum kind %v", d.Kind())}
	}
}

// ArrowAdd performs addition using Arrow compute.
func ArrowAdd(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryArith("add", left, right, ctx)
}

// ArrowSub performs subtraction using Arrow compute.
func ArrowSub(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryArith("subtract", left, right, ctx)
}

// ArrowMul performs multiplication using Arrow compute.
func ArrowMul(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryArith("multiply", left, right, ctx)
}

// ArrowDiv performs division using Arrow compute.
func ArrowDiv(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryArith("divide", left, right, ctx)
}

// ArrowPow performs exponentiation using Arrow compute.
func ArrowPow(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryArith("power", left, right, ctx)
}

// ArrowEq performs equality comparison using Arrow compute.
func ArrowEq(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryCompare("equal", left, right, ctx)
}

// ArrowNe performs not-equal comparison using Arrow compute.
func ArrowNe(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryCompare("not_equal", left, right, ctx)
}

// ArrowLt performs less-than comparison using Arrow compute.
func ArrowLt(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryCompare("less", left, right, ctx)
}

// ArrowLe performs less-or-equal comparison using Arrow compute.
func ArrowLe(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryCompare("less_equal", left, right, ctx)
}

// ArrowGt performs greater-than comparison using Arrow compute.
func ArrowGt(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryCompare("greater", left, right, ctx)
}

// ArrowGe performs greater-or-equal comparison using Arrow compute.
func ArrowGe(left, right Series, ctx *aargh.Context) Series {
	return arrowBinaryCompare("greater_equal", left, right, ctx)
}

// ArrowAnd performs boolean AND using Arrow compute.
func ArrowAnd(left, right Series, ctx *aargh.Context) Series {
	return arrowBooleanOp("and", left, right, ctx)
}

// ArrowOr performs boolean OR using Arrow compute.
func ArrowOr(left, right Series, ctx *aargh.Context) Series {
	return arrowBooleanOp("or", left, right, ctx)
}

// coerceToSeries normalizes an `any` value to a Series for use in arrow ops.
func coerceToSeries(val any, ctx *aargh.Context) Series {
	if s, ok := val.(Series); ok {
		return s
	}
	return NewSeries(val, nil, false, false, ctx)
}

// ArrowBinaryOp is the unified entry point for Arrow-based binary operations.
// It normalizes `other` to a Series, checks contexts, and dispatches to the
// appropriate Arrow compute function.
func ArrowBinaryOp(op string, left Series, other any, ctx *aargh.Context) Series {
	right := coerceToSeries(other, ctx)
	if right.IsError() {
		return right
	}

	// Verify contexts match (unless right is Error/NA)
	if right.Type() != meta.ErrorType && right.Type() != meta.NullType {
		if ctx != right.GetContext() {
			return Errors{fmt.Sprintf("ArrowBinaryOp(%s): cannot operate on series with different contexts", op)}
		}
	}

	switch op {
	case "add":
		return ArrowAdd(left, right, ctx)
	case "sub":
		return ArrowSub(left, right, ctx)
	case "mul":
		return ArrowMul(left, right, ctx)
	case "div":
		return ArrowDiv(left, right, ctx)
	case "pow":
		return ArrowPow(left, right, ctx)
	case "eq":
		return ArrowEq(left, right, ctx)
	case "ne":
		return ArrowNe(left, right, ctx)
	case "lt":
		return ArrowLt(left, right, ctx)
	case "le":
		return ArrowLe(left, right, ctx)
	case "gt":
		return ArrowGt(left, right, ctx)
	case "ge":
		return ArrowGe(left, right, ctx)
	case "and":
		return ArrowAnd(left, right, ctx)
	case "or":
		return ArrowOr(left, right, ctx)
	default:
		return Errors{fmt.Sprintf("ArrowBinaryOp: unknown op %q", op)}
	}
}
