package io

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caerbannogwhite/enchanter/series"
)

func Test_IoAll(t *testing.T) {
	var iod *IoData

	iod = FromXpt(ctx).
		SetPath(filepath.Join(testDataFolder, "xpt_test_mixed.xpt")).
		Read()

	if iod.Error != nil {
		t.Error(iod.Error.Error())
	}

	// to JSON
	iod.ToJson().
		SetPath(filepath.Join(testDataFolder, "test.json")).
		Write()

	if iod.Error != nil {
		t.Error(iod.Error.Error())
	}

	// From JSON
	iod = FromJson(ctx).
		SetPath(filepath.Join(testDataFolder, "test.json")).
		Read()

	if iod.Error != nil {
		t.Error(iod.Error.Error())
	}

	// to CSV
	iod.ToCsv().
		SetPath(filepath.Join(testDataFolder, "test.csv")).
		SetNaText("").
		SetDelimiter(',').
		Write()

	if iod.Error != nil {
		t.Error(iod.Error.Error())
	}

	// From CSV
	iod = FromCsv(ctx).
		SetPath(filepath.Join(testDataFolder, "test.csv")).
		SetDelimiter(',').
		SetHeader(true).
		Read()

	if iod.Error != nil {
		t.Error(iod.Error.Error())
	}

	// to Excel
	iod.ToXlsx().
		SetPath(filepath.Join(testDataFolder, "test.xlsx")).
		SetSheet("test").
		SetNaText("").
		Write()

	if iod.Error != nil {
		t.Error(iod.Error.Error())
	}

	// From Excel
	iod = FromXlsx(ctx).
		SetPath(filepath.Join(testDataFolder, "test.xlsx")).
		SetSheet("test").
		Read()

	if iod.Error != nil {
		t.Error(iod.Error.Error())
	}

	os.Remove(filepath.Join(testDataFolder, "test.json"))
	os.Remove(filepath.Join(testDataFolder, "test.csv"))
	os.Remove(filepath.Join(testDataFolder, "test.xlsx"))

	if iod.NCols() != 4 {
		t.Errorf("expected 4 columns, got %d", iod.NCols())
	}

	if iod.NRows() != 7 {
		t.Errorf("expected 7 rows, got %d", iod.NRows())
	}

	charvar1 := series.NewSeriesString([]string{
		"abcdefghij",
		"wbiwbui749",
		"abcdefghij",
		"wbiwbui749",
		"abcdefghij",
		"wbiwbui749",
		"abcdefghij",
	}, nil, false, ctx)

	charvar2 := series.NewSeriesString([]string{
		"abcdefghijklmnopqrst",
		"nionione983203jnfui2",
		"abcdefghijklmnopqrst",
		"nionione983203jnfui2",
		"abcdefghijklmnopqrst",
		"nionione983203jnfui2",
		"abcdefghijklmnopqrst",
	}, nil, false, ctx)

	numvar1 := series.NewSeriesInt64([]int64{1, 2, 3, 4, 5, 6, 7}, nil, false, ctx)
	numvar2 := series.NewSeriesFloat64([]float64{
		1.2345e2,
		6.5432e1,
		1.2345e0,
		6.5432e-1,
		1.2345e-2,
		6.543e-3, // 6.5432e-3
		1.23e-4,  // 1.2345e-4
	}, nil, false, ctx)

	if _series := iod.ByName("CHARVAR1"); _series == nil {
		t.Error("expected series, got nil")
	} else if !_series.Eq(charvar1).(series.Bools).All() {
		t.Errorf("expected %s, got %s", charvar1.DataAsString(), _series.DataAsString())
	}

	if _series := iod.ByName("CHARVAR2"); _series == nil {
		t.Error("expected series, got nil")
	} else if !_series.Eq(charvar2).(series.Bools).All() {
		t.Errorf("expected %s, got %s", charvar2.DataAsString(), _series.DataAsString())
	}

	if _series := iod.ByName("NUMVAR1"); _series == nil {
		t.Error("expected series, got nil")
	} else if !_series.Eq(numvar1).(series.Bools).All() {
		t.Errorf("expected %s, got %s", numvar1.DataAsString(), _series.DataAsString())
	}

	if _series := iod.ByName("FOO"); _series == nil {
		t.Error("expected series, got nil")
	} else if !_series.Eq(numvar2).(series.Bools).All() {
		t.Errorf("expected %s, got %s", numvar2.DataAsString(), _series.DataAsString())
	}
}
