package io

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/meta"
	"github.com/caerbannogwhite/enchanter/series"
)

type JsonReader struct {
	path   string
	reader io.Reader
	schema *meta.Schema
	ctx    *enchanter.Context
}

func NewJsonReader(ctx *enchanter.Context) *JsonReader {
	return &JsonReader{
		path:   "",
		reader: nil,
		schema: nil,
		ctx:    ctx,
	}
}

func (r *JsonReader) SetPath(path string) *JsonReader {
	r.path = path
	return r
}

func (r *JsonReader) SetReader(reader io.Reader) *JsonReader {
	r.reader = reader
	return r
}

func (r *JsonReader) SetSchema(schema *meta.Schema) *JsonReader {
	r.schema = schema
	return r
}

func (r *JsonReader) Read() *IoData {
	if r.path != "" {
		file, err := os.OpenFile(r.path, os.O_RDONLY, 0666)
		if err != nil {
			return &IoData{Error: err}
		}
		defer file.Close()
		r.reader = file
	}

	if r.reader == nil {
		return &IoData{Error: fmt.Errorf("JsonReader: no reader specified")}
	}

	if r.ctx == nil {
		return &IoData{Error: fmt.Errorf("JsonReader: no context specified")}
	}

	names, series, err := readJson(r.reader, r.schema, r.ctx)
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
		iod.AddSeries(series[i], SeriesMeta{
			Name: name,
		})
	}

	return &iod
}

func readJson(reader io.Reader, schema *meta.Schema, ctx *enchanter.Context) ([]string, []series.Series, error) {

	tokens := json.NewDecoder(reader)

	token, err := tokens.Token()
	if err != nil || token != json.Delim('{') {
		return nil, nil, fmt.Errorf("readJson: invalid json, expected '{' at the beginning")
	}

	var names []string
	var _series []series.Series

	// read by column
	for {

		// read column name or end of json
		token, err := tokens.Token()
		if err != nil {
			return nil, nil, fmt.Errorf("readJson: invalid json, expected column name or '}'")
		}

		switch t := token.(type) {
		case json.Delim:
			if t == json.Delim('}') {
				return names, _series, nil
			} else {
				return nil, nil, fmt.Errorf("readJson: invalid json, expected column name or '}'")
			}

		case string:
			names = append(names, t)

		default:
			return nil, nil, fmt.Errorf("readJson: invalid json, expected column name or '}'")
		}

		// read array
		token, err = tokens.Token()
		if err != nil {
			return nil, nil, fmt.Errorf("readJson: invalid json, column '%s' expected array begin, got '%v'", names[len(names)-1], token)
		}

		if token != json.Delim('{') {
			return nil, nil, fmt.Errorf("readJson: invalid json, column '%s' expected array begin, got '%v'", names[len(names)-1], token)
		}

		// read column values
		var type_ meta.BaseType
		var boolValues []bool
		var floatValues []float64
		var stringValues []string

		for tokens.More() {

			// read index
			_, err := tokens.Token()
			if err != nil {
				return nil, nil, fmt.Errorf("readJson: invalid json")
			}

			token, err = tokens.Token()
			if err != nil {
				return nil, nil, fmt.Errorf("readJson: invalid json")
			}

			switch t := token.(type) {
			case bool:
				type_ = meta.BoolType
				boolValues = append(boolValues, t)

			case float64:
				type_ = meta.Float64Type
				floatValues = append(floatValues, t)

			case string:
				type_ = meta.StringType
				stringValues = append(stringValues, t)

			default:
				return nil, nil, fmt.Errorf("readJson: invalid json")
			}
		}

		// read end of array
		token, err = tokens.Token()
		if err != nil {
			return nil, nil, fmt.Errorf("readJson: invalid json, column '%s' expected array end, got '%v'", names[len(names)-1], token)
		}

		if token != json.Delim('}') {
			return nil, nil, fmt.Errorf("readJson: invalid json, column '%s' expected array end, got '%v'", names[len(names)-1], token)
		}

		// create series
		switch type_ {
		case meta.BoolType:
			_series = append(_series, series.NewSeriesBool(boolValues, nil, false, ctx))
		case meta.Float64Type:
			_series = append(_series, series.NewSeriesFloat64(floatValues, nil, false, ctx))
		case meta.StringType:
			_series = append(_series, series.NewSeriesString(stringValues, nil, false, ctx))
		}
	}
}

type JsonWriter struct {
	path    string
	newLine string
	indent  string
	writer  io.Writer
	ioData  *IoData
}

func NewJsonWriter() *JsonWriter {
	return &JsonWriter{
		path:    "",
		newLine: "\n",
		indent:  "\t",
		writer:  nil,
		ioData:  nil,
	}
}

func (w *JsonWriter) SetPath(path string) *JsonWriter {
	w.path = path
	return w
}

func (w *JsonWriter) SetNewLine(newLine string) *JsonWriter {
	w.newLine = newLine
	return w
}

func (w *JsonWriter) SetIndent(indent string) *JsonWriter {
	w.indent = indent
	return w
}

func (w *JsonWriter) SetWriter(writer io.Writer) *JsonWriter {
	w.writer = writer
	return w
}

func (w *JsonWriter) SetIoData(ioData *IoData) *JsonWriter {
	w.ioData = ioData
	return w
}

func (w *JsonWriter) Write() error {
	if w.ioData == nil {
		return fmt.Errorf("writeJson: no ioData specified")
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
		return fmt.Errorf("writeJson: no writer specified")
	}

	err := writeJson(w.ioData, w.writer, w.newLine, w.indent)
	if err != nil {
		return err
	}

	return nil
}

func writeJson(ioData *IoData, writer io.Writer, newLine, indent string) error {
	indent2 := indent + indent

	writer.Write([]byte("{\n"))
	for i, meta := range ioData.SeriesMeta {
		writer.Write([]byte(fmt.Sprintf("%s\"%s\": {%s", indent2, meta.Name, newLine)))

		_series := ioData.Series[i]
		switch ser := _series.(type) {
		case series.Bools:
			for j, b := range ser.Bools() {
				if ser.IsNull(j) {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": null", indent2, j)))
				} else {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": %v", indent2, j, b)))
				}

				if j < ser.Len()-1 {
					writer.Write([]byte(","))
				}
				writer.Write([]byte(newLine))
			}

		case series.Ints:
			for j, n := range ser.Ints() {
				if ser.IsNull(j) {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": null", indent2, j)))
				} else {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": %d", indent2, j, n)))
				}

				if j < ser.Len()-1 {
					writer.Write([]byte(","))
				}
				writer.Write([]byte(newLine))
			}

		case series.Int64s:
			for j, n := range ser.Int64s() {
				if ser.IsNull(j) {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": null", indent2, j)))
				} else {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": %d", indent2, j, n)))
				}

				if j < ser.Len()-1 {
					writer.Write([]byte(","))
				}
				writer.Write([]byte(newLine))
			}

		case series.Float64s:
			for j, f := range ser.Float64s() {
				if ser.IsNull(j) {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": null", indent2, j)))
				} else {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": %f", indent2, j, f)))
				}

				if j < ser.Len()-1 {
					writer.Write([]byte(","))
				}
				writer.Write([]byte(newLine))
			}

		case series.Strings:
			for j, s := range ser.Strings() {
				if ser.IsNull(j) {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": null", indent2, j)))
				} else {
					writer.Write([]byte(fmt.Sprintf("%s\"%d\": \"%s\"", indent2, j, s)))
				}

				if j < ser.Len()-1 {
					writer.Write([]byte(","))
				}
				writer.Write([]byte(newLine))
			}
		}

		writer.Write([]byte(indent + "}"))
		if i < ioData.NCols()-1 {
			writer.Write([]byte(","))
		}
		writer.Write([]byte(newLine))
	}

	writer.Write([]byte("}" + newLine))
	return nil
}
