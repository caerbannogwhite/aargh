package dataframe

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/arrowutil"
	"github.com/caerbannogwhite/aargh/series"
)

// ArrowSchema returns the Arrow schema corresponding to this DataFrame's columns.
func (df BaseDataFrame) ArrowSchema() *arrow.Schema {
	fields := make([]arrow.Field, len(df.series))
	for i, s := range df.series {
		dt := arrowutil.BaseTypeToArrowType(s.Type())
		if dt == nil {
			dt = arrow.Null
		}
		fields[i] = arrow.Field{
			Name:     df.names[i],
			Type:     dt,
			Nullable: s.IsNullable(),
		}
	}
	return arrow.NewSchema(fields, nil)
}

// ToArrowRecord converts this DataFrame to an Arrow Record (columnar batch).
// The caller should call Release() on the returned record when done.
func (df BaseDataFrame) ToArrowRecord() arrow.Record {
	schema := df.ArrowSchema()
	cols := make([]arrow.Array, len(df.series))
	for i, s := range df.series {
		arr := s.ArrowArray()
		if arr == nil {
			// Build a null array for unsupported types
			builder := array.NewNullBuilder(memory.DefaultAllocator)
			for j := 0; j < s.Len(); j++ {
				builder.AppendNull()
			}
			cols[i] = builder.NewNullArray()
			builder.Release()
		} else {
			cols[i] = arr
		}
	}
	return array.NewRecord(schema, cols, int64(df.NRows()))
}

// NewBaseDataFrameFromArrowRecord creates a BaseDataFrame from an Arrow Record.
// Each column in the record becomes a Series in the DataFrame.
func NewBaseDataFrameFromArrowRecord(record arrow.Record, ctx *aargh.Context) DataFrame {
	if ctx == nil {
		return BaseDataFrame{err: fmt.Errorf("NewBaseDataFrameFromArrowRecord: context is nil")}
	}

	df := NewBaseDataFrame(ctx).(BaseDataFrame)

	schema := record.Schema()
	for i := 0; i < int(record.NumCols()); i++ {
		col := record.Column(i)
		name := schema.Field(i).Name
		s := series.ArrowArrayToSeries(col, ctx)
		if s.IsError() {
			df.err = fmt.Errorf("NewBaseDataFrameFromArrowRecord: column %q: %s", name, s.GetError())
			return df
		}
		df = df.AddSeries(name, s).(BaseDataFrame)
	}

	return df
}
