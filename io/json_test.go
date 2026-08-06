package io

import (
	"os"

	"testing"

	"github.com/caerbannogwhite/enchanter/series"
)

func Test_IoJson_ValidWrite(t *testing.T) {
	iod := NewIoData(ctx)

	iod.AddSeries(series.NewSeriesFloat64([]float64{1, 2, 3}, nil, false, ctx), SeriesMeta{Name: "a"})
	iod.AddSeries(series.NewSeriesString([]string{"a", "b", "c"}, nil, false, ctx), SeriesMeta{Name: "b"})

	err := iod.ToJson().
		SetPath("test.json").
		Write()

	if err != nil {
		t.Error(err.Error())
	}

	_, err = os.Stat("test.json")
	if err != nil {
		t.Error(err.Error())
	}

	err = os.Remove("test.json")
	if err != nil {
		t.Error(err.Error())
	}
}
