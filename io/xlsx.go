package io

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/caerbannogwhite/aargh"
	"github.com/caerbannogwhite/aargh/meta"
	"github.com/caerbannogwhite/aargh/series"

	"github.com/tealeg/xlsx"
)

type XlsxReader struct {
	path             string
	sheet            string
	header           int
	rows             int
	guessDataTypeLen int
	strict           bool
	schema           *meta.Schema
	ctx              *aargh.Context
}

func NewXlsxReader(ctx *aargh.Context) *XlsxReader {
	return &XlsxReader{
		path:             "",
		sheet:            "",
		header:           0,
		rows:             -1,
		guessDataTypeLen: aargh.XLSX_READER_DEFAULT_GUESS_DATA_TYPE_LEN,
		strict:           true,
		schema:           nil,
		ctx:              ctx,
	}
}

// SetPath sets the path to the XLSX file.
func (r *XlsxReader) SetPath(path string) *XlsxReader {
	r.path = path
	return r
}

// SetSheet sets the sheet to read from the XLSX file.
func (r *XlsxReader) SetSheet(sheet string) *XlsxReader {
	r.sheet = sheet
	return r
}

// SetHeader sets the header row to read from the XLSX file.
func (r *XlsxReader) SetHeader(header int) *XlsxReader {
	r.header = header
	return r
}

// SetRows sets the number of rows to read from the XLSX file.
func (r *XlsxReader) SetRows(rows int) *XlsxReader {
	r.rows = rows
	return r
}

// SetGuessDataTypeLen sets the number of rows to read from the XLSX file for type guessing.
// Default is 1000.
func (r *XlsxReader) SetGuessDataTypeLen(guessDataTypeLen int) *XlsxReader {
	r.guessDataTypeLen = guessDataTypeLen
	return r
}

// SetStrict sets if the strict type guessing is enabled.
// Default is true.
func (r *XlsxReader) SetStrict(strict bool) *XlsxReader {
	r.strict = strict
	return r
}

// SetSchema sets the schema for the XLSX file.
func (r *XlsxReader) SetSchema(schema *meta.Schema) *XlsxReader {
	r.schema = schema
	return r
}

// Read reads the XLSX file and returns a IoData object.
func (r *XlsxReader) Read() *IoData {
	if r.ctx == nil {
		return &IoData{Error: fmt.Errorf("XlsxReader: no context specified")}
	}

	names, series, err := r.readXlsx()

	if err != nil {
		return &IoData{Error: err}
	}

	iod := IoData{
		FileMeta: FileMeta{
			FileName: r.path,
			FilePath: r.path,
		},
		ctx: r.ctx,
	}

	for i, name := range names {
		iod.AddSeries(series[i], SeriesMeta{Name: name})
	}

	return &iod
}

type xlsxRowReader struct {
	sh    *xlsx.Sheet
	row   int
	cols  int
	cells []*xlsx.Cell
}

func (r *xlsxRowReader) Read() ([]string, error) {
	if r.row >= r.sh.MaxRow {
		return nil, io.EOF
	}

	row := r.sh.Row(r.row)
	r.cells = row.Cells
	r.row++

	values := make([]string, len(r.cells))
	for i, cell := range r.cells {
		values[i] = cell.String()
	}

	if len(values) < r.cols {
		for i := len(values); i < r.cols; i++ {
			values = append(values, "")
		}
	}

	return values, nil
}

func (r *XlsxReader) readXlsx() ([]string, []series.Series, error) {
	wb, err := xlsx.OpenFile(r.path)
	if err != nil {
		return nil, nil, err
	}

	sh, ok := wb.Sheet[r.sheet]
	if !ok {
		return nil, nil, fmt.Errorf("Sheet %s not found", r.sheet)
	}

	if r.rows < 0 {
		r.rows = sh.MaxRow
	}

	names := make([]string, len(sh.Row(r.header).Cells))
	for i, cell := range sh.Row(r.header).Cells {
		names[i] = cell.String()
	}

	xlsxRowReader := &xlsxRowReader{
		sh:    sh,
		row:   r.header + 1,
		cols:  len(names),
		cells: nil,
	}

	series, err := readRowData(xlsxRowReader, r.guessDataTypeLen, r.strict, r.rows, r.schema, r.ctx)
	if err != nil {
		return nil, nil, err
	}

	return names, series, nil
}

type XlsxWriter struct {
	path   string
	sheet  string
	naText string
	writer io.Writer
	ioData *IoData
}

func NewXlsxWriter() *XlsxWriter {
	return &XlsxWriter{
		path:   "",
		sheet:  "Sheet1",
		writer: nil,
		ioData: nil,
	}
}

func (w *XlsxWriter) SetPath(path string) *XlsxWriter {
	w.path = path
	return w
}

func (w *XlsxWriter) SetSheet(sheet string) *XlsxWriter {
	w.sheet = sheet
	return w
}

func (w *XlsxWriter) SetNaText(naText string) *XlsxWriter {
	w.naText = naText
	return w
}

func (w *XlsxWriter) SetWriter(writer io.Writer) *XlsxWriter {
	w.writer = writer
	return w
}

func (w *XlsxWriter) SetIoData(ioData *IoData) *XlsxWriter {
	w.ioData = ioData
	return w
}

func (w *XlsxWriter) Write() error {
	if w.ioData == nil {
		return fmt.Errorf("XlsxWriter: no ioData specified")
	}

	if w.path != "" {
		// make sure os.O_WRONLY arg is supplies so that file is not opened in read-only mode
		file, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		w.writer = file
	}

	if w.writer == nil {
		return fmt.Errorf("XlsxWriter: no writer specified")
	}

	err := writeXlsx(w.ioData, w.writer, w.sheet, w.naText)
	if err != nil {
		return err
	}

	return nil
}

func writeXlsx(ioData *IoData, writer io.Writer, sheetName string, naText string) error {

	file := xlsx.NewFile()

	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return err
	}

	// write header
	row := sheet.AddRow()
	for _, meta := range ioData.SeriesMeta {
		cell := row.AddCell()
		cell.Value = meta.Name
	}

	// write data
	for i := 0; i < ioData.NRows(); i++ {
		row := sheet.AddRow()
		for j := range ioData.SeriesMeta {
			cell := row.AddCell()

			switch s := ioData.At(j).(type) {
			case series.Bools:
				if s.IsNull(i) {
					cell.Value = naText
					continue
				}
				cell.SetBool(s.Get(i).(bool))

			case series.Ints:
				if s.IsNull(i) {
					cell.Value = naText
					continue
				}
				cell.SetInt(s.Get(i).(int))

			case series.Int64s:
				if s.IsNull(i) {
					cell.Value = naText
					continue
				}
				cell.SetInt64(s.Get(i).(int64))

			case series.Float64s:
				if s.IsNull(i) {
					cell.Value = naText
					continue
				}
				cell.SetFloat(s.Get(i).(float64))

			case series.Strings:
				if s.IsNull(i) {
					cell.Value = naText
					continue
				}
				cell.Value = s.Get(i).(string)

			case series.Times:
				if s.IsNull(i) {
					cell.Value = naText
					continue
				}
				cell.SetDateTime(s.Get(i).(time.Time))

			case series.Durations:
				if s.IsNull(i) {
					cell.Value = naText
					continue
				}
				cell.SetInt(int(s.Get(i).(time.Duration).Nanoseconds()))

			default:
				cell.Value = naText
			}
		}
	}

	err = file.Write(writer)
	if err != nil {
		return err
	}

	return nil
}
