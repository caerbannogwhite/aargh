package io

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
)

const FILE_FORMAT_PARQUET FileFormat = "PARQUET"
const FILE_FORMAT_ARROW_IPC FileFormat = "ARROW_IPC"

// makeRecord creates an Arrow Record from a schema and columns.
func makeRecord(schema *arrow.Schema, cols []arrow.Array, nrows int64) arrow.Record {
	return array.NewRecord(schema, cols, nrows)
}

// recordToTable wraps an Arrow Record in a Table (single-chunk columns).
func recordToTable(rec arrow.Record) arrow.Table {
	return array.NewTableFromRecords(rec.Schema(), []arrow.Record{rec})
}
