package io

import (
	"os"

	"testing"

	"github.com/caerbannogwhite/enchanter/series"
)

func Test_IoMarkDown_ValidWrite(t *testing.T) {
	iod := NewIoData(ctx)

	iod.AddSeries(series.NewSeriesFloat64([]float64{1, 2, 3}, nil, false, ctx), SeriesMeta{Name: "a"})
	iod.AddSeries(series.NewSeriesString([]string{"a", "b", "c"}, nil, false, ctx), SeriesMeta{Name: "b"})

	err := iod.ToMarkdown().
		SetPath("test.md").
		Write()

	if err != nil {
		t.Error(err.Error())
	}

	_, err = os.Stat("test.md")
	if err != nil {
		t.Error(err.Error())
	}

	err = os.Remove("test.md")
	if err != nil {
		t.Error(err.Error())
	}
}
