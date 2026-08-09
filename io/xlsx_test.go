package io

import (
	"os"
	"path/filepath"

	"testing"

	"github.com/caerbannogwhite/enchanter/series"
)

func Test_IoXlsx_ValidWrite(t *testing.T) {
	iod := NewIoData(ctx)

	iod.AddSeries(series.NewSeriesFloat64([]float64{1, 2, 3}, nil, false, ctx), SeriesMeta{Name: "a"})
	iod.AddSeries(series.NewSeriesString([]string{"a", "b", "c"}, nil, false, ctx), SeriesMeta{Name: "b"})

	err := iod.ToXlsx().
		SetPath("test.xlsx").
		Write()

	if err != nil {
		t.Error(err.Error())
	}

	_, err = os.Stat("test.xlsx")
	if err != nil {
		t.Error(err.Error())
	}

	err = os.Remove("test.xlsx")
	if err != nil {
		t.Error(err.Error())
	}
}

// Test_IoXlsx_SetRowsLimitsRows is a regression test for XlsxReader.SetRows.
//
// SetRows(k) is meant to cap how many data rows the reader returns, but it was
// silently ignored: readXlsx passed sh.MaxRow to readRowData instead of the
// user-supplied row cap, so every read returned the whole sheet regardless of
// SetRows. This test writes a sheet with more data rows than the requested cap
// and asserts the reader honours the cap.
//
// It fails against the pre-fix behaviour: with the cap ignored, reading with
// SetRows(rowLimit) returned all totalRows rows, so NRows() == totalRows != rowLimit.
func Test_IoXlsx_SetRowsLimitsRows(t *testing.T) {
	const (
		totalRows = 10
		rowLimit  = 4 // strictly fewer than totalRows, so the cap is observable
	)

	data := make([]int64, totalRows)
	for i := range data {
		data[i] = int64(i)
	}

	iod := NewIoData(ctx)
	iod.AddSeries(series.NewSeriesInt64(data, nil, false, ctx), SeriesMeta{Name: "a"})

	path := filepath.Join(t.TempDir(), "setrows.xlsx")

	if err := iod.ToXlsx().SetPath(path).SetSheet("test").Write(); err != nil {
		t.Fatalf("write failed: %s", err.Error())
	}

	// Control: without a row cap the full sheet is returned. This confirms the
	// fixture actually holds totalRows data rows, so the capped assertion below
	// is meaningful (and not passing simply because the sheet is short).
	full := FromXlsx(ctx).SetPath(path).SetSheet("test").Read()
	if full.Error != nil {
		t.Fatalf("uncapped read failed: %s", full.Error.Error())
	}
	if full.NRows() != totalRows {
		t.Fatalf("uncapped read: expected %d rows, got %d", totalRows, full.NRows())
	}

	// The behaviour under test: SetRows must limit the returned rows to rowLimit.
	limited := FromXlsx(ctx).SetPath(path).SetSheet("test").SetRows(rowLimit).Read()
	if limited.Error != nil {
		t.Fatalf("capped read failed: %s", limited.Error.Error())
	}
	if limited.NRows() != rowLimit {
		t.Errorf("SetRows(%d) not honoured: expected %d rows, got %d", rowLimit, rowLimit, limited.NRows())
	}
}
