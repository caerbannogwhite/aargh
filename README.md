## Aargh 🧙‍♂️

[![Go Reference](https://pkg.go.dev/badge/github.com/caerbannogwhite/aargh.svg)](https://pkg.go.dev/github.com/caerbannogwhite/aargh)

### What is it?

Aargh is a library for data wrangling in Go.
The goal is to provide a simple and efficient API for data manipulation in Go,
similar to Pandas or Polars in Python, and Dplyr in R.
It supports nullable types: null data is optimized for memory usage.

Aargh is a work in progress, and the API is not stable yet.
The DataFrame package is still being developed.

However, it already supports the following formats:

- CSV
- XPT (SAS)
- XLSX
- HTML
- Markdown

### Examples

```go
package main

import (
	"strings"

	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/dataframe"
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

	dataframe.NewBaseDataFrame(aargh.NewContext()).
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

### Dependencies

Built with:

- [xslx](https://github.com/tealeg/xlsx/tree/master)
- [lipgloss](https://github.com/charmbracelet/lipgloss)

### TODO

- [ ] Improve filtering interface.
- [ ] Improve dataframe PrettyPrint: add parameters, optimize data display, use lipgloss.
- [ ] Implement string factors.
- [ ] Times: set time format.
- [ ] Implement `Set(i []int, v []any) Series`.
- [ ] Add `Slice(i []int) Series` (using filter?).
- [ ] Implement memory optimized Bool series with uint64.
- [ ] Use uint64 for null mask.
- [ ] Optimize XPT reader/writer with float32.
- [ ] Add url resolver to each reader.
- [ ] Add format option to each writer.
- [ ] JSON reader by records.
- [ ] Implement chunked series.
- [ ] Implement Parquet reader and writer.
- [ ] Implement SPSS reader and writer.
- [ ] Implement SAS7BDAT reader and writer (https://cran.r-project.org/web/packages/sas7bdat/vignettes/sas7bdat.pdf)
