package dataframe

import (
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/meta"
	"github.com/caerbannogwhite/enchanter/series"
)

type DataFrameJoinType int8

const (
	INNER_JOIN DataFrameJoinType = iota
	LEFT_JOIN
	RIGHT_JOIN
	OUTER_JOIN
)

type DataFrame interface {

	// Basic accessors.

	// GetContext returns the context of the dataframe.
	GetContext() *enchanter.Context

	// Names returns the names of the series in the dataframe.
	Names() []string
	// Types returns the types of the series in the dataframe.
	Types() []meta.BaseType
	// NCols returns the number of columns in the dataframe.
	NCols() int
	// NRows returns the number of rows in the dataframe.
	NRows() int

	IsErrored() bool

	IsGrouped() bool

	GetError() error

	GetSeriesIndex(name string) int

	// Add new series to the dataframe.

	// AddSeries adds a generic series to the dataframe.
	AddSeries(name string, series series.Series) DataFrame
	// AddSeriesFromBools adds a series of bools to the dataframe.
	AddSeriesFromBools(name string, data []bool, nullMask []bool, makeCopy bool) DataFrame
	// AddSeriesFromInt32s adds a series of ints to the dataframe.
	AddSeriesFromInts(name string, data []int, nullMask []bool, makeCopy bool) DataFrame
	// AddSeriesFromInt64s adds a series of ints to the dataframe.
	AddSeriesFromInt64s(name string, data []int64, nullMask []bool, makeCopy bool) DataFrame
	// AddSeriesFromFloat64s adds a series of floats to the dataframe.
	AddSeriesFromFloat64s(name string, data []float64, nullMask []bool, makeCopy bool) DataFrame
	// AddSeriesFromStrings adds a series of strings to the dataframe.
	AddSeriesFromStrings(name string, data []string, nullMask []bool, makeCopy bool) DataFrame
	// AddSeriesFromTimes adds a series of times to the dataframe.
	AddSeriesFromTimes(name string, data []time.Time, nullMask []bool, makeCopy bool) DataFrame
	// AddSeriesFromDurations adds a series of durations to the dataframe.
	AddSeriesFromDurations(name string, data []time.Duration, nullMask []bool, makeCopy bool) DataFrame

	// Replace the series with the given name.
	Replace(name string, s series.Series) DataFrame

	// Returns the column with the given name.
	C(name string) series.Series

	// Returns the series at the given index.
	At(index int) series.Series

	// Returns the series with the given name as a bool series.
	NameAt(index int) string

	// Select keeps the columns matched by selectors, which are regular
	// expressions matched unanchored: "Car" also selects "CarOrigin". Use
	// "^Car$" to match a single column exactly.
	Select(selectors ...string) DataFrame

	SelectAt(indices ...int) DataFrame

	Filter(mask any) DataFrame

	GroupBy(by ...string) DataFrame

	Ungroup() DataFrame

	getPartitions() []series.SeriesPartition

	Join(how DataFrameJoinType, other DataFrame, on ...string) DataFrame

	Take(params ...int) DataFrame

	Agg(aggregators ...aggregator) aggregatorBuilder

	// Sort the dataframe.
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
	OrderBy(params ...SortParam) DataFrame

	// IO

	Describe() string
	Records(header bool) [][]string

	// Pretty print the dataframe.
	PPrint(params PPrintParams) DataFrame

	FromCsv() *csvReaderWrapper
	ToCsv() *csvWriterWrapper

	FromJson() *jsonReaderWrapper
	ToJson() *jsonWriterWrapper

	FromXpt() *xptReaderWrapper
	ToXpt() *xptWriterWrapper

	FromXlsx() *xlsxReaderWrapper
	ToXlsx() *xlsxWriterWrapper

	ToHtml() *htmlWriterWrapper
	ToMarkDown() *markDownWriterWrapper

	FromParquet() *parquetReaderWrapper
	ToParquet() *parquetWriterWrapper

	FromArrowIPC() *arrowIPCReaderWrapper
	ToArrowIPC() *arrowIPCWriterWrapper

	// Arrow interop.
	ArrowSchema() *arrow.Schema
	ToArrowRecord() arrow.RecordBatch
}
