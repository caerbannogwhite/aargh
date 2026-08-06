package io

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"io"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/meta"
)

type CsvReader struct {
	header           bool
	rows             int
	delimiter        rune
	guessDataTypeLen int
	path             string
	nullValues       bool
	reader           io.Reader
	schema           *meta.Schema
	ctx              *enchanter.Context
}

func NewCsvReader(ctx *enchanter.Context) *CsvReader {
	return &CsvReader{
		header:           enchanter.CSV_READER_DEFAULT_HEADER,
		rows:             -1,
		delimiter:        enchanter.CSV_READER_DEFAULT_DELIMITER,
		guessDataTypeLen: enchanter.CSV_READER_DEFAULT_GUESS_DATA_TYPE_LEN,
		path:             "",
		nullValues:       false,
		reader:           nil,
		schema:           nil,
		ctx:              ctx,
	}
}

func (r *CsvReader) SetHeader(header bool) *CsvReader {
	r.header = header
	return r
}

func (r *CsvReader) SetDelimiter(delimiter rune) *CsvReader {
	r.delimiter = delimiter
	return r
}

func (r *CsvReader) SetGuessDataTypeLen(guessDataTypeLen int) *CsvReader {
	r.guessDataTypeLen = guessDataTypeLen
	return r
}

func (r *CsvReader) SetRows(rows int) *CsvReader {
	r.rows = rows
	return r
}

func (r *CsvReader) SetPath(path string) *CsvReader {
	r.path = path
	return r
}

func (r *CsvReader) SetNullValues(nullValues bool) *CsvReader {
	r.nullValues = nullValues
	return r
}

func (r *CsvReader) SetReader(reader io.Reader) *CsvReader {
	r.reader = reader
	return r
}

func (r *CsvReader) SetSchema(schema *meta.Schema) *CsvReader {
	r.schema = schema
	return r
}

func (r *CsvReader) SetContext(ctx *enchanter.Context) *CsvReader {
	r.ctx = ctx
	return r
}

func (r *CsvReader) Read() *IoData {
	if r.path != "" {
		file, err := os.OpenFile(r.path, os.O_RDONLY, 0666)
		if err != nil {
			return &IoData{Error: err}
		}
		defer file.Close()
		r.reader = file
	}

	if r.reader == nil {
		return &IoData{Error: fmt.Errorf("CsvReader: no reader specified")}
	}

	if r.ctx == nil {
		return &IoData{Error: fmt.Errorf("CsvReader: no context specified")}
	}

	return r.readCsv()
}

// readCsv reads a CSV file and returns a IoData.
func (r *CsvReader) readCsv() *IoData {

	// TODO: Optimize null masks (use bit vectors)?
	// TODO: Try to optimize this function by using goroutines: read the rows (like 1000)
	//		and guess the data types in parallel

	if r.ctx == nil {
		return &IoData{Error: fmt.Errorf("readCsv: no context specified")}
	}

	var err error
	var fileMeta FileMeta

	if r.path != "" {
		fileInfo, err := os.Stat(r.path)
		if err != nil {
			return &IoData{Error: fmt.Errorf("readCsv: %w", err)}
		}

		fileMeta.FileSize = fileInfo.Size()
		fileMeta.FileName = filepath.Base(r.path)
		fileMeta.FilePath = filepath.Dir(r.path)
		fileMeta.FileExt = filepath.Ext(r.path)
		fileMeta.FileFormat = FILE_FORMAT_CSV
	} else {
		fileMeta.FileSize = 0
		fileMeta.FileName = "Unknown"
		fileMeta.FilePath = "Unknown"
		fileMeta.FileExt = ".csv"
		fileMeta.FileFormat = FILE_FORMAT_CSV
	}

	// Initialize CSV reader
	csvReader := csv.NewReader(r.reader)
	csvReader.Comma = r.delimiter
	csvReader.FieldsPerRecord = -1

	// Read header if present
	var seriesMeta []SeriesMeta
	if r.header {
		names, err := csvReader.Read()
		if err != nil {
			return &IoData{Error: err}
		}
		for _, name := range names {
			seriesMeta = append(seriesMeta, SeriesMeta{
				Name: name,
			})
		}
	}

	series, err := readRowData(csvReader, r.nullValues, r.guessDataTypeLen, r.rows, r.schema, r.ctx)
	if err != nil {
		return &IoData{Error: err}
	}

	// Generate names if not present
	if !r.header {
		for i := 0; i < len(series); i++ {
			seriesMeta = append(seriesMeta, SeriesMeta{
				Name: fmt.Sprintf("Column %d", i+1),
			})
		}
	}

	for i, s := range series {
		seriesMeta[i].Type = s.Type()
	}

	return &IoData{
		FileMeta:   fileMeta,
		SeriesMeta: seriesMeta,
		Series:     series,
		ctx:        r.ctx,
	}
}

type CsvQuotingType int

const (
	CsvQuotingNone CsvQuotingType = iota
	CsvQuotingAll
	CsvQuotingNeeded
	CsvQuotingNonNumeric
)

type CsvWriter struct {
	delimiter              rune
	header                 bool
	format                 bool // TODO: Implement this
	useParamNaText         bool
	useParamDateTimeFormat bool
	useParamEol            bool
	useParamQuote          bool
	path                   string
	naText                 string
	dateTimeFormat         string
	eol                    string
	quote                  string
	quoting                CsvQuotingType
	writer                 io.Writer
	ioData                 *IoData
}

func NewCsvWriter() *CsvWriter {
	return &CsvWriter{
		delimiter:              enchanter.CSV_READER_DEFAULT_DELIMITER,
		header:                 enchanter.CSV_READER_DEFAULT_HEADER,
		format:                 true,
		useParamNaText:         false,
		useParamDateTimeFormat: false,
		useParamEol:            false,
		useParamQuote:          false,
		path:                   "",
		naText:                 enchanter.NA_TEXT,
		dateTimeFormat:         enchanter.DATE_TIME_FORMAT,
		eol:                    enchanter.EOL,
		quote:                  enchanter.QUOTE,
		quoting:                CsvQuotingNeeded,
		writer:                 nil,
		ioData:                 nil,
	}
}

func (w *CsvWriter) SetDelimiter(delimiter rune) *CsvWriter {
	w.delimiter = delimiter
	return w
}

func (w *CsvWriter) SetHeader(header bool) *CsvWriter {
	w.header = header
	return w
}

func (w *CsvWriter) SetFormat(format bool) *CsvWriter {
	w.format = format
	return w
}

func (w *CsvWriter) SetPath(path string) *CsvWriter {
	w.path = path
	return w
}

func (w *CsvWriter) SetNaText(naText string) *CsvWriter {
	w.useParamNaText = true
	w.naText = naText
	return w
}

func (w *CsvWriter) SetDateTimeFormat(dateTimeFormat string) *CsvWriter {
	w.useParamDateTimeFormat = true
	w.dateTimeFormat = dateTimeFormat
	return w
}

func (w *CsvWriter) SetEol(eol string) *CsvWriter {
	w.useParamEol = true
	w.eol = eol
	return w
}

func (w *CsvWriter) SetQuote(quote string) *CsvWriter {
	w.useParamQuote = true
	w.quote = quote
	return w
}

func (w *CsvWriter) SetQuoting(quoting CsvQuotingType) *CsvWriter {
	w.quoting = quoting
	return w
}

func (w *CsvWriter) SetWriter(writer io.Writer) *CsvWriter {
	w.writer = writer
	return w
}

func (w *CsvWriter) SetIoData(ioData *IoData) *CsvWriter {
	w.ioData = ioData
	return w
}

func (w *CsvWriter) Write() error {
	if w.ioData == nil {
		return fmt.Errorf("CsvWriter: no ioData specified")
	}

	if w.ioData.Error != nil {
		return w.ioData.Error
	}

	if !w.useParamEol {
		w.eol = w.ioData.ctx.GetEol()
	}

	if !w.useParamNaText {
		w.naText = w.ioData.ctx.GetNaText()
	}

	if !w.useParamDateTimeFormat {
		w.dateTimeFormat = w.ioData.ctx.GetDateTimeFormat()
	}

	if !w.useParamQuote {
		w.quote = w.ioData.ctx.GetQuote()
	}

	if w.path != "" {
		file, err := os.OpenFile(w.path, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		w.writer = file
	}

	if w.writer == nil {
		return fmt.Errorf("CsvWriter: no writer specified")
	}

	err := w.writeCsv()
	if err != nil {
		return err
	}

	return nil
}

func (w *CsvWriter) writeCsv() error {
	if w.header {
		for i, meta := range w.ioData.SeriesMeta {
			if i > 0 {
				fmt.Fprintf(w.writer, "%c", w.delimiter)
			}
			fmt.Fprintf(w.writer, "%s", meta.Name)
		}

		fmt.Fprintf(w.writer, "%s", w.eol)
	}

	for i := 0; i < w.ioData.NRows(); i++ {
		for j, s := range w.ioData.Series {
			if j > 0 {
				fmt.Fprintf(w.writer, "%c", w.delimiter)
			}

			if s.IsNull(i) {
				fmt.Fprintf(w.writer, "%s", w.naText)
			} else {
				switch w.quoting {
				case CsvQuotingNone:
					fmt.Fprintf(w.writer, "%s", s.GetAsString(i))

				case CsvQuotingAll:
					fmt.Fprintf(w.writer, "%s", w.quote+s.GetAsString(i)+w.quote)

				case CsvQuotingNeeded:
					str := s.GetAsString(i)
					if strings.Contains(str, w.quote) {
						str = strings.ReplaceAll(str, w.quote, w.quote+w.quote)
						fmt.Fprintf(w.writer, "%s", w.quote+str+w.quote)
					} else if strings.Contains(str, string(w.delimiter)) || strings.Contains(str, w.eol) {
						fmt.Fprintf(w.writer, "%s", w.quote+str+w.quote)
					} else {
						fmt.Fprintf(w.writer, "%s", str)
					}

				case CsvQuotingNonNumeric:
					if s.Type() == meta.Float64Type || s.Type() == meta.Int64Type || s.Type() == meta.IntType {
						fmt.Fprintf(w.writer, "%s", s.GetAsString(i))
					} else {
						fmt.Fprintf(w.writer, "%s", w.quote+s.GetAsString(i)+w.quote)
					}

				default:
					return fmt.Errorf("writeCsv: invalid quoting type")
				}
			}
		}

		fmt.Fprintf(w.writer, "%s", w.eol)
	}

	return nil
}
