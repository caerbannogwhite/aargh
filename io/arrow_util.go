package io

import (
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
)

const FILE_FORMAT_PARQUET FileFormat = "PARQUET"
const FILE_FORMAT_ARROW_IPC FileFormat = "ARROW_IPC"

// makeRecord creates an Arrow record batch from a schema and columns.
func makeRecord(schema *arrow.Schema, cols []arrow.Array, nrows int64) arrow.RecordBatch {
	return array.NewRecordBatch(schema, cols, nrows)
}

// recordToTable wraps an Arrow record batch in a Table (single-chunk columns).
func recordToTable(rec arrow.RecordBatch) arrow.Table {
	return array.NewTableFromRecords(rec.Schema(), []arrow.RecordBatch{rec})
}
