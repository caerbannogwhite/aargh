// Package io implements the file readers and writers behind the dataframe
// I/O chains: CSV, XLSX, XPT (SAS transport, versions 5-9), JSON, HTML,
// Markdown, Parquet and Arrow IPC (Feather). A SAS7BDAT reader is under
// development (header parsing only).
//
// Readers follow a builder pattern and produce an [IoData] — series plus
// per-column and per-file metadata — which dataframes consume:
//
//	iod := io.FromParquet(ctx).SetPath("people.parquet").Read()
//
// Most users reach this package through the wrappers on DataFrame
// (FromCsv, ToParquet, ...) rather than directly.
package io
