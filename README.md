## Enchanter 🧙‍♂️

[![Go Reference](https://pkg.go.dev/badge/github.com/caerbannogwhite/enchanter.svg)](https://pkg.go.dev/github.com/caerbannogwhite/enchanter)
[![CI](https://github.com/caerbannogwhite/enchanter/actions/workflows/ci.yml/badge.svg)](https://github.com/caerbannogwhite/enchanter/actions/workflows/ci.yml)

> Formerly known as **aargh** — renamed in v0.2.0. The v0.1.x releases remain
> available under the old module path.

Enchanter is a data wrangling library in pure Go, in the spirit of Pandas and
Polars in Python and dplyr in R: typed, nullable columns (series) composed
into DataFrames with select / filter / group by / join / sort / aggregate
pipelines. Series and DataFrames interoperate with
[Apache Arrow](https://arrow.apache.org/), which also powers the Parquet and
Arrow IPC readers and writers.

Enchanter is a work in progress and the API is not stable yet.

### Install

```sh
go get github.com/caerbannogwhite/enchanter
```

Requires Go 1.26+. Pure Go, no cgo.

### Quick start

```go
package main

import (
	"strings"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/dataframe"
)

func main() {
	data1 := `
name,age,weight,junior,department,salary band
Alice C,29,75.0,F,HR,4
John Doe,30,80.5,true,IT,2
Bob,31,85.0,F,IT,4
Jane H,25,60.0,false,IT,4
Mary,28,70.0,false,IT,3
Oliver,32,90.0,true,HR,1
Ursula,27,65.0,f,Business,4
Charlie,33,60.0,t,Business,2
Megan,26,55.0,F,IT,3`

	dataframe.NewBaseDataFrame(enchanter.NewContext()).
		FromCsv().
		SetReader(strings.NewReader(data1)).
		Read().
		Select("department", "age", "weight", "junior").
		GroupBy("department").
		Agg(dataframe.Min("age"), dataframe.Max("weight"), dataframe.Mean("junior"), dataframe.Count()).
		Run().
		PPrint(dataframe.NewPPrintParams().SetUseLipGloss(true))
}

//   BaseDataFrame: 3 rows, 5 columns
// ╭────────────┬──────────┬─────────────┬──────────────┬───────╮
// │ department │ min(age) │ max(weight) │ mean(junior) │ n     │
// ├────────────┼──────────┼─────────────┼──────────────┼───────┤
// │ String     │ Float64  │ Float64     │ Float64      │ Int64 │
// ├────────────┼──────────┼─────────────┼──────────────┼───────┤
// │ HR         │    29.00 │       90.00 │       0.5000 │ 2.000 │
// │ IT         │    25.00 │       85.00 │       0.2000 │ 5.000 │
// │ Business   │    27.00 │       65.00 │       0.5000 │ 2.000 │
// ╰────────────┴──────────┴─────────────┴──────────────┴───────╯
```

More runnable examples live in [examples](examples) and in the package
documentation on [pkg.go.dev](https://pkg.go.dev/github.com/caerbannogwhite/enchanter).

### Reading and writing files

| Format              | Read | Write | Notes                                    |
| ------------------- | :--: | :---: | ---------------------------------------- |
| CSV                 |  ✅  |  ✅   | delimiter, header, type guessing         |
| Parquet             |  ✅  |  ✅   | via Apache Arrow; types + nulls preserved |
| Arrow IPC (Feather) |  ✅  |  ✅   |                                          |
| XPT (SAS transport) |  ✅  |  ✅   | versions 5–9                             |
| XLSX                |  ✅  |  ✅   |                                          |
| JSON                |  ✅  |  ✅   | record-oriented                          |
| HTML                |  ✅  |  ✅   | tables                                   |
| Markdown            |  ✅  |  ✅   | tables                                   |
| SAS7BDAT            |  🚧  |   —   | header parsing only; data reading planned |

All readers and writers share the same builder style:

```go
// Parquet round trip: types and nulls survive, unlike CSV.
err := df.ToParquet().SetPath("people.parquet").Write()

df2 := dataframe.NewBaseDataFrame(ctx).
	FromParquet().
	SetPath("people.parquet").
	Read()
```

See [examples/parquet](examples/parquet/main.go) for a runnable version.

### Apache Arrow interop

Every series can produce an Arrow array and every DataFrame an Arrow record
batch, which makes Enchanter data exchangeable with anything that speaks
Arrow (DuckDB, Polars, pandas, DataFusion, Spark, ...):

```go
rec := df.ToArrowRecord() // freshly built, owned by the record
defer rec.Release()       // optional under the default GC-backed allocator

df2 := dataframe.NewBaseDataFrameFromArrowRecord(rec, ctx)
```

Conversion notes:

- Enchanter null masks map onto Arrow validity bitmaps in both directions.
- `series.ArrowArrayToSeries` accepts more Arrow types than Enchanter has
  series types (all integer widths, Float32, Date32/64, every timestamp
  unit, LargeString, ...) and funnels them onto the native types: integers
  normalize to int64, timestamps to nanoseconds.
- Conversions are materialized copies — no reference to the source Arrow
  memory is retained, so you may release inputs immediately.
- Releasing values returned to you is optional under GC-backed allocators
  (the default); see `Context.Allocator` for the exact ownership contract.

### Supported data types

The data types not checked are not yet supported, but might be in the future.

- [x] Bool
- [ ] Bool (memory optimized, not fully implemented yet)
- [ ] Int16
- [x] Int
- [x] Int64
- [ ] Float32
- [x] Float64
- [ ] Complex64
- [ ] Complex128
- [x] String
- [x] Time
- [x] Duration

### Supported operations for Series

- [x] Filter

  - [x] filter by bool slice
  - [x] filter by int slice
  - [x] filter by bool series
  - [x] filter by int series

- [x] Group

  - [x] Group (with nulls)
  - [x] SubGroup (with nulls)

- [x] Map
- [x] Sort

  - [x] Sort (with nulls)
  - [x] SortRev (with nulls)

- [x] Take

### Supported operations for DataFrame

- [x] Agg
- [x] Filter
- [x] GroupBy
- [x] Join

  - [x] Inner
  - [x] Left
  - [x] Right
  - [x] Outer
  - [x] Inner with nulls
  - [x] Left with nulls
  - [x] Right with nulls
  - [x] Outer with nulls

- [ ] Map
- [x] OrderBy
- [x] Select
- [x] Take
- [ ] Pivot
- [ ] Stack/Append

### Supported stats functions

- [x] Count
- [x] Sum
- [x] Mean
- [ ] Median
- [x] Min
- [x] Max
- [x] StdDev
- [ ] Variance
- [ ] Quantile

### Development

```sh
go test -short ./...     # fast unit tests (skips large local benchmark fixtures)
go test ./...            # full suite
go test -bench . ./...   # benchmarks (need local fixtures in testdata/)
go generate ./series/    # regenerate *_base.go / *_ops.go via the generators module
```

CI runs build, vet, gofmt, a generated-code freshness check, the test suite
and the race detector on Linux and Windows.

### Roadmap

Next (0.3.0):

- [ ] Route series operations through Arrow compute kernels (replacing the generated per-type loops).
- [ ] Aggregations via Arrow compute (Sum, Min, Max, Mean).
- [ ] Dictionary-encoded (factor) strings.
- [ ] SAS7BDAT data reading ([format notes](https://cran.r-project.org/web/packages/sas7bdat/vignettes/sas7bdat.pdf)).

Later:

- [ ] Pivot longer / wider (started on the `dev-pivot` branch).
- [ ] Custom aggregators (prototype archived as `archive/dev-fix-agg`).
- [ ] Stricter CSV type guessing with acceptance threshold (archived as `archive/dev-0.1.3`).
- [ ] Implement chunked series.
- [ ] Implement SPSS reader and writer.
- [ ] Improve filtering interface.
- [ ] Improve dataframe PrettyPrint: add parameters, optimize data display, use lipgloss.
- [ ] Times: set time format.
- [ ] Implement `Set(i []int, v []any) Series`.
- [ ] Add `Slice(i []int) Series` (using filter?).
- [ ] Implement memory optimized Bool series with uint64.
- [ ] Use uint64 for null mask.
- [ ] Optimize XPT reader/writer with float32.
- [ ] Add url resolver to each reader.
- [ ] Add format option to each writer.
- [ ] JSON reader by records.

### Dependencies

Built with:

- [arrow-go](https://github.com/apache/arrow-go)
- [xslx](https://github.com/tealeg/xlsx/tree/master)
- [lipgloss](https://github.com/charmbracelet/lipgloss)

### License

See [LICENSE](LICENSE).
