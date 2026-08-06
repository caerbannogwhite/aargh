package io

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/meta"
	"github.com/caerbannogwhite/enchanter/series"
)

type typeBucket struct {
	nullCount     int
	boolCount     int
	intCount      int
	floatCount    int
	stringCount   int
	dateCount     int
	datetimeCount int
}

// Get the most common type in the bucket and whether it is the only type
func (tb *typeBucket) getMostCommonType() (meta.BaseType, bool) {
	timeCount := tb.dateCount + tb.datetimeCount
	if tb.boolCount > tb.intCount && tb.boolCount > tb.floatCount && tb.boolCount > tb.stringCount && tb.boolCount > timeCount {
		return meta.BoolType, tb.nullCount+tb.intCount+tb.floatCount+tb.stringCount+timeCount == 0
	} else if tb.intCount > tb.floatCount && tb.intCount > tb.stringCount && tb.intCount > timeCount {
		return meta.Int64Type, tb.nullCount+tb.boolCount+tb.floatCount+tb.stringCount+timeCount == 0
	} else if tb.floatCount > tb.stringCount && tb.floatCount > timeCount {
		return meta.Float64Type, tb.nullCount+tb.boolCount+tb.intCount+tb.stringCount+timeCount == 0
	} else if timeCount > tb.stringCount {
		return meta.TimeType, tb.nullCount+tb.boolCount+tb.intCount+tb.floatCount+tb.stringCount == 0
	}
	return meta.StringType, tb.nullCount+tb.boolCount+tb.intCount+tb.floatCount+timeCount == 0
}

type typeGuesser struct {
	nullValues     bool
	nullRegex      *regexp.Regexp
	boolRegex      *regexp.Regexp
	boolTrueRegex  *regexp.Regexp
	boolFalseRegex *regexp.Regexp
	intRegex       *regexp.Regexp
	floatRegex     *regexp.Regexp
	dateRegex      *regexp.Regexp
	datetimeRegex  *regexp.Regexp

	// For each column, count the number of values that match each type
	typeBuckets []typeBucket
}

// Get the regexes for guessing data types
func newTypeGuesser(nullValues bool) typeGuesser {
	return typeGuesser{
		nullValues,
		regexp.MustCompile(`^([Nn][Uu][Ll][Ll])$|^([Nn][Aa][Nn]?)$|^([Nn]/[Aa])$|^$`),
		regexp.MustCompile(`^[Tt]([Rr][Uu][Ee])?$|^[Ff]([Aa][Ll][Ss][Ee])?$`),
		regexp.MustCompile(`^[Tt]([Rr][Uu][Ee])?$`),
		regexp.MustCompile(`^[Ff]([Aa][Ll][Ss][Ee])?$`),
		regexp.MustCompile(`^[-+]?[0-9]+$`),
		regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$`),
		// Date patterns: YYYY-MM-DD, MM/DD/YYYY, DD/MM/YYYY, MM/DD/YY, DD/MM/YY
		regexp.MustCompile(`^(19|20)\d{2}[-/](0[1-9]|1[0-2])[-/](0[1-9]|[12]\d|3[01])$|^(0[1-9]|1[0-2])[-/](0[1-9]|[12]\d|3[01])[-/]((19|20)\d{2}|\d{2})$|^(0[1-9]|[12]\d|3[01])[-/](0[1-9]|1[0-2])[-/]((19|20)\d{2}|\d{2})$`),
		// Datetime patterns: YYYY-MM-DD HH:MM:SS, MM/DD/YYYY HH:MM:SS, MM/DD/YY HH:MM:SS, etc.
		regexp.MustCompile(`^(19|20)\d{2}[-/](0[1-9]|1[0-2])[-/](0[1-9]|[12]\d|3[01])\s+(0[0-9]|1[0-9]|2[0-3]):([0-5]\d)(:([0-5]\d))?(\s*[AP]M)?$|^(0[1-9]|1[0-2])[-/](0[1-9]|[12]\d|3[01])[-/]((19|20)\d{2}|\d{2})\s+(0[0-9]|1[0-9]|2[0-3]):([0-5]\d)(:([0-5]\d))?(\s*[AP]M)?$|^(0[1-9]|[12]\d|3[01])[-/](0[1-9]|1[0-2])[-/]((19|20)\d{2}|\d{2})\s+(0[0-9]|1[0-9]|2[0-3]):([0-5]\d)(:([0-5]\d))?(\s*[AP]M)?$`),
		nil,
	}
}

func (tg *typeGuesser) setLength(length int) {
	tg.typeBuckets = make([]typeBucket, length)
}

func (tg *typeGuesser) guessType(record string) meta.BaseType {
	if tg.boolRegex.MatchString(record) {
		return meta.BoolType
	} else if tg.intRegex.MatchString(record) {
		return meta.Int64Type
	} else if tg.floatRegex.MatchString(record) {
		return meta.Float64Type
	} else if tg.datetimeRegex.MatchString(record) {
		return meta.TimeType
	} else if tg.dateRegex.MatchString(record) {
		return meta.TimeType
	}
	return meta.StringType
}

func (tg *typeGuesser) guessTypes(records []string) {
	for i, v := range records {
		if tg.boolRegex.MatchString(v) {
			tg.typeBuckets[i].boolCount++
		} else if tg.intRegex.MatchString(v) {
			tg.typeBuckets[i].intCount++
		} else if tg.floatRegex.MatchString(v) {
			tg.typeBuckets[i].floatCount++
		} else if tg.datetimeRegex.MatchString(v) {
			tg.typeBuckets[i].datetimeCount++
		} else if tg.dateRegex.MatchString(v) {
			tg.typeBuckets[i].dateCount++
		} else {
			tg.typeBuckets[i].stringCount++
		}
	}
}

func (tg *typeGuesser) guessTypesNulls(records []string) {
	for i, v := range records {
		if tg.boolRegex.MatchString(v) {
			tg.typeBuckets[i].boolCount++
		} else if tg.intRegex.MatchString(v) {
			tg.typeBuckets[i].intCount++
		} else if tg.floatRegex.MatchString(v) {
			tg.typeBuckets[i].floatCount++
		} else if tg.datetimeRegex.MatchString(v) {
			tg.typeBuckets[i].datetimeCount++
		} else if tg.dateRegex.MatchString(v) {
			tg.typeBuckets[i].dateCount++
		} else if tg.nullRegex.MatchString(v) {
			tg.typeBuckets[i].nullCount++
		} else {
			tg.typeBuckets[i].stringCount++
		}
	}
}

func (tg typeGuesser) getTypes() []meta.BaseType {
	types := make([]meta.BaseType, len(tg.typeBuckets))
	if tg.nullValues {
		for i, v := range tg.typeBuckets {
			types[i], _ = v.getMostCommonType()
		}
	} else {
		var onlyType bool
		for i, v := range tg.typeBuckets {
			types[i], onlyType = v.getMostCommonType()
			if !onlyType {
				types[i] = meta.StringType
			}
		}
	}
	return types
}

func (tg typeGuesser) atoBool(s string) (bool, error) {
	if tg.boolTrueRegex.MatchString(s) {
		return true, nil
	} else if tg.boolFalseRegex.MatchString(s) {
		return false, nil
	}
	return false, fmt.Errorf("cannot convert \"%s\" to bool", s)
}

func (tg typeGuesser) atoTime(s string) (time.Time, error) {
	// Try various date and datetime formats
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"01/02/2006 15:04:05",
		"01/02/2006 15:04",
		"01/02/2006",
		"02/01/2006 15:04:05",
		"02/01/2006 15:04",
		"02/01/2006",
		"01-02-2006 15:04:05",
		"01-02-2006 15:04",
		"01-02-2006",
		"02-01-2006 15:04:05",
		"02-01-2006 15:04",
		"02-01-2006",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		// Handle AM/PM formats
		"2006-01-02 03:04:05 PM",
		"2006-01-02 03:04 PM",
		"01/02/2006 03:04:05 PM",
		"01/02/2006 03:04 PM",
		"02/01/2006 03:04:05 PM",
		"02/01/2006 03:04 PM",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot convert \"%s\" to time", s)
}

type RowDataProvider interface {
	Read() ([]string, error)
}

func readRowData(reader RowDataProvider, nullValues bool, guessDataTypeLen int, maxLen int, schema *meta.Schema, ctx *enchanter.Context) ([]series.Series, error) {
	var dataTypes []meta.BaseType
	var recordsForGuessing [][]string

	// Initialize TypeGuesser
	tg := newTypeGuesser(nullValues)

	if maxLen < 0 {
		maxLen = math.MaxInt32
	} else if guessDataTypeLen > maxLen {
		guessDataTypeLen = maxLen
	}

	counter := 0

	// Guess data types
	if schema == nil {
		recordsForGuessing = make([][]string, guessDataTypeLen)

		// Read first record to get length
		record, err := reader.Read()
		counter++

		if err != nil && err != io.EOF {
			return nil, err
		}
		recordsForGuessing[0] = record

		tg.setLength(len(record))

		if nullValues {
			tg.guessTypesNulls(record)
			for i := 1; i < guessDataTypeLen; i++ {
				record, err := reader.Read()
				counter++

				if err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}
				recordsForGuessing[i] = record
				tg.guessTypesNulls(record)
			}
		} else {
			tg.guessTypes(record)
			for i := 1; i < guessDataTypeLen; i++ {
				record, err := reader.Read()
				counter++

				if err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}
				recordsForGuessing[i] = record
				tg.guessTypes(record)
			}
		}
		dataTypes = tg.getTypes()
	} else

	// Use schema
	{
		dataTypes = schema.GetDataTypes()
	}

	nullMasks := make([][]bool, len(dataTypes))
	if nullValues {
		for i := range nullMasks {
			nullMasks[i] = make([]bool, 0)
		}
	}

	values := make([]interface{}, len(dataTypes))
	for i := range values {
		switch dataTypes[i] {
		case meta.BoolType:
			values[i] = make([]bool, 0)
		case meta.IntType:
			values[i] = make([]int, 0)
		case meta.Int64Type:
			values[i] = make([]int64, 0)
		case meta.Float64Type:
			values[i] = make([]float64, 0)
		case meta.StringType:
			values[i] = make([]*string, 0)
		case meta.TimeType:
			values[i] = make([]time.Time, 0)
		}
	}

	// If no schema: add records for guessing to values
	if schema == nil {
		if nullValues {
			for _, record := range recordsForGuessing {
				for i, v := range record {
					switch dataTypes[i] {
					case meta.BoolType:
						if b, err := tg.atoBool(v); err != nil {
							nullMasks[i] = append(nullMasks[i], true)
							values[i] = append(values[i].([]bool), false)
						} else {
							nullMasks[i] = append(nullMasks[i], false)
							values[i] = append(values[i].([]bool), b)
						}

					case meta.IntType:
						if d, err := strconv.Atoi(v); err != nil {
							nullMasks[i] = append(nullMasks[i], true)
							values[i] = append(values[i].([]int), 0)
						} else {
							nullMasks[i] = append(nullMasks[i], false)
							values[i] = append(values[i].([]int), int(d))
						}

					case meta.Int64Type:
						if d, err := strconv.ParseInt(v, 10, 64); err != nil {
							nullMasks[i] = append(nullMasks[i], true)
							values[i] = append(values[i].([]int64), 0)
						} else {
							nullMasks[i] = append(nullMasks[i], false)
							values[i] = append(values[i].([]int64), d)
						}

					case meta.Float64Type:
						if f, err := strconv.ParseFloat(v, 64); err != nil {
							nullMasks[i] = append(nullMasks[i], true)
							values[i] = append(values[i].([]float64), math.NaN())
						} else {
							nullMasks[i] = append(nullMasks[i], false)
							values[i] = append(values[i].([]float64), f)
						}

					case meta.StringType:
						nullMasks[i] = append(nullMasks[i], false)
						values[i] = append(values[i].([]*string), ctx.StringPool.Put(v))

					case meta.TimeType:
						if t, err := tg.atoTime(v); err != nil {
							nullMasks[i] = append(nullMasks[i], true)
							values[i] = append(values[i].([]time.Time), time.Time{})
						} else {
							nullMasks[i] = append(nullMasks[i], false)
							values[i] = append(values[i].([]time.Time), t)
						}
					}
				}
			}
		} else {
			for _, record := range recordsForGuessing {
				for i, v := range record {
					switch dataTypes[i] {
					case meta.BoolType:
						b, err := tg.atoBool(v)
						if err != nil {
							return nil, err
						}
						values[i] = append(values[i].([]bool), b)

					case meta.IntType:
						d, err := strconv.Atoi(v)
						if err != nil {
							return nil, err
						}
						values[i] = append(values[i].([]int), int(d))

					case meta.Int64Type:
						d, err := strconv.ParseInt(v, 10, 64)
						if err != nil {
							return nil, err
						}
						values[i] = append(values[i].([]int64), d)

					case meta.Float64Type:
						f, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return nil, err
						}
						values[i] = append(values[i].([]float64), f)

					case meta.StringType:
						values[i] = append(values[i].([]*string), ctx.StringPool.Put(v))

					case meta.TimeType:
						t, err := tg.atoTime(v)
						if err != nil {
							return nil, err
						}
						values[i] = append(values[i].([]time.Time), t)
					}
				}
			}
		}
	}

	if nullValues {
		for {
			if counter >= maxLen {
				break
			}

			record, err := reader.Read()
			counter++

			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}

			for i, v := range record {
				switch dataTypes[i] {
				case meta.BoolType:
					if b, err := tg.atoBool(v); err != nil {
						nullMasks[i] = append(nullMasks[i], true)
						values[i] = append(values[i].([]bool), false)
					} else {
						nullMasks[i] = append(nullMasks[i], false)
						values[i] = append(values[i].([]bool), b)
					}

				case meta.IntType:
					if d, err := strconv.Atoi(v); err != nil {
						nullMasks[i] = append(nullMasks[i], true)
						values[i] = append(values[i].([]int), 0)
					} else {
						nullMasks[i] = append(nullMasks[i], false)
						values[i] = append(values[i].([]int), int(d))
					}

				case meta.Int64Type:
					if d, err := strconv.ParseInt(v, 10, 64); err != nil {
						nullMasks[i] = append(nullMasks[i], true)
						values[i] = append(values[i].([]int64), 0)
					} else {
						nullMasks[i] = append(nullMasks[i], false)
						values[i] = append(values[i].([]int64), d)
					}

				case meta.Float64Type:
					if f, err := strconv.ParseFloat(v, 64); err != nil {
						nullMasks[i] = append(nullMasks[i], true)
						values[i] = append(values[i].([]float64), math.NaN())
					} else {
						nullMasks[i] = append(nullMasks[i], false)
						values[i] = append(values[i].([]float64), f)
					}

				case meta.StringType:
					nullMasks[i] = append(nullMasks[i], false)
					values[i] = append(values[i].([]*string), ctx.StringPool.Put(v))

				case meta.TimeType:
					if t, err := tg.atoTime(v); err != nil {
						nullMasks[i] = append(nullMasks[i], true)
						values[i] = append(values[i].([]time.Time), time.Time{})
					} else {
						nullMasks[i] = append(nullMasks[i], false)
						values[i] = append(values[i].([]time.Time), t)
					}
				}
			}
		}
	} else {
		for {
			if counter >= maxLen {
				break
			}

			record, err := reader.Read()
			counter++

			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}

			for i, v := range record {
				switch dataTypes[i] {
				case meta.BoolType:
					b, err := tg.atoBool(v)
					if err != nil {
						return nil, err
					}
					values[i] = append(values[i].([]bool), b)

				case meta.IntType:
					d, err := strconv.Atoi(v)
					if err != nil {
						return nil, err
					}
					values[i] = append(values[i].([]int), int(d))

				case meta.Int64Type:
					d, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						return nil, err
					}
					values[i] = append(values[i].([]int64), d)

				case meta.Float64Type:
					f, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, err
					}
					values[i] = append(values[i].([]float64), f)

				case meta.StringType:
					values[i] = append(values[i].([]*string), ctx.StringPool.Put(v))

				case meta.TimeType:
					t, err := tg.atoTime(v)
					if err != nil {
						return nil, err
					}
					values[i] = append(values[i].([]time.Time), t)
				}
			}
		}
	}

	// Create series
	_series := make([]series.Series, len(dataTypes))
	for i := range dataTypes {
		switch dataTypes[i] {
		case meta.BoolType:
			_series[i] = series.NewSeriesBool(values[i].([]bool), nullMasks[i], false, ctx)

		case meta.IntType:
			_series[i] = series.NewSeriesInt(values[i].([]int), nullMasks[i], false, ctx)

		case meta.Int64Type:
			_series[i] = series.NewSeriesInt64(values[i].([]int64), nullMasks[i], false, ctx)

		case meta.Float64Type:
			_series[i] = series.NewSeriesFloat64(values[i].([]float64), nullMasks[i], false, ctx)

		case meta.StringType:
			_series[i] = series.NewSeriesStringFromPtrs(values[i].([]*string), nullMasks[i], false, ctx)

		case meta.TimeType:
			_series[i] = series.NewSeriesTime(values[i].([]time.Time), nullMasks[i], false, ctx)
		}
	}

	return _series, nil
}
