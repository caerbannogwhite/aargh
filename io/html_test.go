package io

import (
	"os"

	"testing"

	"github.com/caerbannogwhite/enchanter/series"
)

func Test_IoHtml_ValidWrite(t *testing.T) {
	iod := NewIoData(ctx)

	iod.AddSeries(series.NewSeriesFloat64([]float64{1, 2, 3}, nil, false, ctx), SeriesMeta{Name: "a"})
	iod.AddSeries(series.NewSeriesString([]string{"a", "b", "c"}, nil, false, ctx), SeriesMeta{Name: "b"})

	err := iod.ToHtml().
		SetPath("test.html").
		Write()

	if err != nil {
		t.Error(err.Error())
	}

	_, err = os.Stat("test.html")
	if err != nil {
		t.Error(err.Error())
	}

	err = os.Remove("test.html")
	if err != nil {
		t.Error(err.Error())
	}
}
