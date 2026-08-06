package io

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/caerbannogwhite/enchanter/meta"
	"github.com/caerbannogwhite/enchanter/series"
)

func Test_FromCsv(t *testing.T) {

	data := `name,age,weight,junior
Alice C,29,75.0,F
John Doe,30,80.5,true
Bob,31,85.0,T
Jane H,25,60.0,false
Mary,28,70.0,false
Oliver,32,90.0,true
Ursula,27,65.0,f
Charlie,33,95.0,t
`

	// Create a new dataframe from the CSV data.
	iod := FromCsv(ctx).
		SetReader(strings.NewReader(data)).
		SetDelimiter(',').
		SetHeader(true).
		SetGuessDataTypeLen(3).
		Read()

	if iod.Error != nil {
		t.Error(iod.Error)
	}

	// Check the number of rows.
	if iod.NRows() != 8 {
		t.Error("Expected 8 rows, got", iod.NRows())
	}

	// Check the number of columns.
	if iod.NCols() != 4 {
		t.Error("Expected 4 columns, got", iod.NCols())
	}

	// Check the column names.
	if iod.SeriesMetaAt(0).Name != "name" {
		t.Error("Expected 'name', got", iod.SeriesMetaAt(0).Name)
	}

	if iod.SeriesMetaAt(1).Name != "age" {
		t.Error("Expected 'age', got", iod.SeriesMetaAt(1).Name)
	}

	if iod.SeriesMetaAt(2).Name != "weight" {
		t.Error("Expected 'weight', got", iod.SeriesMetaAt(2).Name)
	}

	if iod.SeriesMetaAt(3).Name != "junior" {
		t.Error("Expected 'junior', got", iod.SeriesMetaAt(3).Name)
	}

	// Check the column types.
	if iod.Types()[0] != meta.StringType {
		t.Error("Expected String, got", iod.Types()[0].String())
	}

	if iod.Types()[1] != meta.Int64Type {
		t.Error("Expected Int64, got", iod.Types()[1].String())
	}

	if iod.Types()[2] != meta.Float64Type {
		t.Error("Expected Float64, got", iod.Types()[2].String())
	}

	if iod.Types()[3] != meta.BoolType {
		t.Error("Expected Bool, got", iod.Types()[3].String())
	}

	// Check the values.
	if iod.At(0).Data().([]string)[0] != "Alice C" {
		t.Error("Expected 'Alice C', got", iod.At(0).Data().([]string)[0])
	}

	if iod.At(0).Data().([]string)[1] != "John Doe" {
		t.Error("Expected 'John Doe', got", iod.At(0).Data().([]string)[1])
	}

	if iod.At(0).Data().([]string)[2] != "Bob" {
		t.Error("Expected 'Bob', got", iod.At(0).Data().([]string)[2])
	}

	if iod.At(0).Data().([]string)[3] != "Jane H" {
		t.Error("Expected 'Jane H', got", iod.At(0).Data().([]string)[3])
	}

	if iod.At(1).Data().([]int64)[4] != 28 {
		t.Error("Expected 28, got", iod.At(1).Data().([]int64)[4])
	}

	if iod.At(1).Data().([]int64)[5] != 32 {
		t.Error("Expected 32, got", iod.At(1).Data().([]int64)[5])
	}

	if iod.At(1).Data().([]int64)[6] != 27 {
		t.Error("Expected 27, got", iod.At(1).Data().([]int64)[6])
	}

	if iod.At(1).Data().([]int64)[7] != 33 {
		t.Error("Expected 33, got", iod.At(1).Data().([]int64)[7])
	}

	if iod.At(2).Data().([]float64)[0] != 75.0 {
		t.Error("Expected 75.0, got", iod.At(2).Data().([]float64)[0])
	}

	if iod.At(2).Data().([]float64)[1] != 80.5 {
		t.Error("Expected 80.5, got", iod.At(2).Data().([]float64)[1])
	}

	if iod.At(2).Data().([]float64)[2] != 85.0 {
		t.Error("Expected 85.0, got", iod.At(2).Data().([]float64)[2])
	}

	if iod.At(2).Data().([]float64)[3] != 60.0 {
		t.Error("Expected 60.0, got", iod.At(2).Data().([]float64)[3])
	}

	if iod.At(3).Data().([]bool)[4] != false {
		t.Error("Expected false, got", iod.At(3).Data().([]bool)[4])
	}

	if iod.At(3).Data().([]bool)[5] != true {
		t.Error("Expected true, got", iod.At(3).Data().([]bool)[5])
	}

	if iod.At(3).Data().([]bool)[6] != false {
		t.Error("Expected false, got", iod.At(3).Data().([]bool)[6])
	}

	if iod.At(3).Data().([]bool)[7] != true {
		t.Error("Expected true, got", iod.At(3).Data().([]bool)[7])
	}
}

func Test_FromCsvWithDates(t *testing.T) {

	data := `name,last name,DOB,last seen
Jackson,Lamb,1950-01-01,2025-01-01 12:00:00
River,Cartwright,1991-02-02,2025-02-02 12:00:00
Catherine,Standish,1992-03-03,2025-03-03 12:00:00
Louisa,Guy,1993-04-04,2025-04-04 12:00:00
`

	// Create a new dataframe from the CSV data.
	iod := FromCsv(ctx).
		SetReader(strings.NewReader(data)).
		SetGuessDataTypeLen(3).
		Read()

	if iod.Error != nil {
		t.Error(iod.Error)
	}

	if iod.NCols() != 4 {
		t.Error("Expected 4 columns, got", iod.NCols())
	}

	if iod.NRows() != 4 {
		t.Error("Expected 4 rows, got", iod.NRows())
	}

	if iod.SeriesMetaAt(0).Name != "name" {
		t.Error("Expected 'name', got", iod.SeriesMetaAt(0).Name)
	}

	if iod.SeriesMetaAt(0).Type != meta.StringType {
		t.Error("Expected String, got", iod.SeriesMetaAt(0).Type.String())
	}

	if iod.SeriesMetaAt(1).Name != "last name" {
		t.Error("Expected 'last name', got", iod.SeriesMetaAt(1).Name)
	}

	if iod.SeriesMetaAt(1).Type != meta.StringType {
		t.Error("Expected String, got", iod.SeriesMetaAt(1).Type.String())
	}

	if iod.SeriesMetaAt(2).Name != "DOB" {
		t.Error("Expected 'DOB', got", iod.SeriesMetaAt(2).Name)
	}

	if iod.SeriesMetaAt(2).Type != meta.TimeType {
		t.Error("Expected Time, got", iod.SeriesMetaAt(2).Type.String())
	}

	if iod.SeriesMetaAt(3).Name != "last seen" {
		t.Error("Expected 'last seen', got", iod.SeriesMetaAt(3).Name)
	}

	if iod.SeriesMetaAt(3).Type != meta.TimeType {
		t.Error("Expected Time, got", iod.SeriesMetaAt(3).Type.String())
	}
}

func Benchmark_FromCsv_100000Rows(b *testing.B) {

	// Create a new dataframe from the CSV data.
	var iod *IoData

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, err := os.OpenFile(filepath.Join(testDataFolder, "organizations-100000.csv"), os.O_RDONLY, 0666)
		if err != nil {
			b.Skip("benchmark fixture not available:", err)
		}

		iod = FromCsv(ctx).
			SetReader(f).
			SetDelimiter(',').
			SetHeader(true).
			SetGuessDataTypeLen(100).
			Read()

		if iod.Error != nil {
			b.Error(iod.Error)
		}

		f.Close()
	}
	b.StopTimer()

	// Check the number of rows.
	if iod.NRows() != 100000 {
		b.Error("Expected 100000 rows, got", iod.NRows())
	}

	// Check the number of columns.
	if iod.NCols() != 9 {
		b.Error("Expected 9 columns, got", iod.NCols())
	}

	names := []string{"Index", "Organization Id", "Name", "Website", "Country", "Description", "Founded", "Industry", "Number of employees"}

	// Check the column names.
	for i := 0; i < len(names); i++ {
		if iod.SeriesMetaAt(i).Name != names[i] {
			b.Error("Expected ", names[i], ", got", iod.SeriesMetaAt(i).Name)
		}
	}

	// Check the column types.
	if iod.Types()[0] != meta.Int64Type {
		b.Error("Expected Int64, got", iod.Types()[0].String())
	}

	if iod.Types()[1] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[1].String())
	}

	if iod.Types()[2] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[2].String())
	}

	if iod.Types()[3] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[3].String())
	}

	if iod.Types()[4] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[4].String())
	}

	if iod.Types()[5] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[5].String())
	}

	if iod.Types()[6] != meta.Int64Type {
		b.Error("Expected Int64, got", iod.Types()[6].String())
	}

	if iod.Types()[7] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[7].String())
	}

	if iod.Types()[8] != meta.Int64Type {
		b.Error("Expected Int64, got", iod.Types()[8].String())
	}
}

func Benchmark_FromCsv_500000Rows(b *testing.B) {
	// Create a new dataframe from the CSV data.
	var iod *IoData

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, err := os.OpenFile(filepath.Join(testDataFolder, "organizations-500000.csv"), os.O_RDONLY, 0666)
		if err != nil {
			b.Skip("benchmark fixture not available:", err)
		}

		iod = FromCsv(ctx).
			SetReader(f).
			SetDelimiter(',').
			SetHeader(true).
			SetGuessDataTypeLen(100).
			Read()

		if iod.Error != nil {
			b.Error(iod.Error)
		}

		f.Close()
	}
	b.StopTimer()

	// Check the number of rows.
	if iod.NRows() != 500000 {
		b.Error("Expected 100000 rows, got", iod.NRows())
	}

	// Check the number of columns.
	if iod.NCols() != 9 {
		b.Error("Expected 9 columns, got", iod.NCols())
	}

	names := []string{"Index", "Organization Id", "Name", "Website", "Country", "Description", "Founded", "Industry", "Number of employees"}

	// Check the column names.
	for i := 0; i < len(names); i++ {
		if iod.SeriesMetaAt(i).Name != names[i] {
			b.Error("Expected ", names[i], ", got", iod.SeriesMetaAt(i).Name)
		}
	}

	// Check the column types.
	if iod.Types()[0] != meta.Int64Type {
		b.Error("Expected Int64, got", iod.Types()[0].String())
	}

	if iod.Types()[1] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[1].String())
	}

	if iod.Types()[2] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[2].String())
	}

	if iod.Types()[3] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[3].String())
	}

	if iod.Types()[4] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[4].String())
	}

	if iod.Types()[5] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[5].String())
	}

	if iod.Types()[6] != meta.Int64Type {
		b.Error("Expected Int64, got", iod.Types()[6].String())
	}

	if iod.Types()[7] != meta.StringType {
		b.Error("Expected String, got", iod.Types()[7].String())
	}

	if iod.Types()[8] != meta.Int64Type {
		b.Error("Expected Int64, got", iod.Types()[8].String())
	}
}

func Test_IoCsv_ValidWrite(t *testing.T) {
	iod := NewIoData(ctx)

	iod.AddSeries(series.NewSeriesFloat64([]float64{1, 2, 3}, []bool{true, false, true}, false, ctx), SeriesMeta{Name: "a"})
	iod.AddSeries(series.NewSeriesString([]string{"a", "b", "c"}, []bool{true, false, true}, false, ctx), SeriesMeta{Name: "b"})

	err := iod.ToCsv().
		SetPath("test.csv").
		Write()

	if err != nil {
		t.Error(err.Error())
	}

	_, err = os.Stat("test.csv")
	if err != nil {
		t.Error(err.Error())
	}

	err = os.Remove("test.csv")
	if err != nil {
		t.Error(err.Error())
	}
}
