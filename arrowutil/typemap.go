package arrowutil

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/caerbannogwhite/aargh/meta"
)

// BaseTypeToArrowType maps an aargh BaseType to the corresponding Arrow DataType.
// Returns nil for types that have no direct Arrow mapping.
func BaseTypeToArrowType(bt meta.BaseType) arrow.DataType {
	switch bt {
	case meta.BoolType:
		return arrow.FixedWidthTypes.Boolean
	case meta.IntType:
		// Go int is platform-dependent; map to int64 for Arrow.
		return arrow.PrimitiveTypes.Int64
	case meta.Int64Type:
		return arrow.PrimitiveTypes.Int64
	case meta.Float32Type:
		return arrow.PrimitiveTypes.Float32
	case meta.Float64Type:
		return arrow.PrimitiveTypes.Float64
	case meta.StringType:
		return arrow.BinaryTypes.String
	case meta.TimeType:
		return arrow.FixedWidthTypes.Timestamp_ns
	case meta.DurationType:
		return &arrow.DurationType{Unit: arrow.Nanosecond}
	default:
		return nil
	}
}

// ArrowTypeToBaseType maps an Arrow DataType to the corresponding aargh BaseType.
// Returns meta.ErrorType for unsupported types.
func ArrowTypeToBaseType(dt arrow.DataType) meta.BaseType {
	switch dt.ID() {
	case arrow.BOOL:
		return meta.BoolType
	case arrow.INT8, arrow.INT16, arrow.INT32:
		return meta.IntType
	case arrow.INT64:
		return meta.Int64Type
	case arrow.UINT8, arrow.UINT16, arrow.UINT32, arrow.UINT64:
		return meta.Int64Type
	case arrow.FLOAT32:
		return meta.Float32Type
	case arrow.FLOAT64:
		return meta.Float64Type
	case arrow.STRING, arrow.LARGE_STRING:
		return meta.StringType
	case arrow.TIMESTAMP:
		return meta.TimeType
	case arrow.DURATION:
		return meta.DurationType
	default:
		return meta.ErrorType
	}
}
