package io

import (
	"fmt"
	"io"
	"os"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/meta"
)

type MarkDownWriter struct {
	header bool
	index  bool
	types  bool
	path   string
	naText string
	writer io.Writer
	ioData *IoData
}

func NewMarkDownWriter() *MarkDownWriter {
	return &MarkDownWriter{
		header: true,
		index:  false,
		path:   "",
		naText: enchanter.NA_TEXT,
		writer: nil,
		ioData: nil,
	}
}

func (w *MarkDownWriter) SetHeader(header bool) *MarkDownWriter {
	w.header = header
	return w
}

func (w *MarkDownWriter) SetIndex(index bool) *MarkDownWriter {
	w.index = index
	return w
}

func (w *MarkDownWriter) SetPath(path string) *MarkDownWriter {
	w.path = path
	return w
}

func (w *MarkDownWriter) SetNaText(naText string) *MarkDownWriter {
	w.naText = naText
	return w
}

func (w *MarkDownWriter) SetWriter(writer io.Writer) *MarkDownWriter {
	w.writer = writer
	return w
}

func (w *MarkDownWriter) SetIoData(ioData *IoData) *MarkDownWriter {
	w.ioData = ioData
	return w
}

func (w *MarkDownWriter) Write() error {
	if w.ioData == nil {
		return fmt.Errorf("writeMarkDown: no ioData specified")
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
		return fmt.Errorf("writeMarkDown: no writer specified")
	}

	err := writeMarkDown(w.ioData, w.writer, w.header, w.index, w.naText)
	if err != nil {
		return err
	}

	return nil
}

func writeMarkDown(ioData *IoData, writer io.Writer, header, index bool, naText string) error {
	buff := ""
	if header {
		buff += "|"
		if index {
			buff += " |"
		}

		for _, col := range ioData.SeriesMeta {
			buff += "**" + col.Name + "**|"
		}
		buff += "\n"
	}

	buff += "|"
	if index {
		buff += "----:|"
	}

	for i := range ioData.SeriesMeta {
		switch ioData.Types()[i].GetDefaultJustification() {
		case meta.JUSTIFY_LEFT:
			buff += ":----|"
		case meta.JUSTIFY_RIGHT:
			buff += "----:|"
		case meta.JUSTIFY_CENTER:
			buff += ":---:|"
		}
	}
	buff += "\n"

	for i := 0; i < ioData.NRows(); i++ {
		buff += "|"
		if index {
			buff += fmt.Sprintf("%d|", i)
		}

		for j := 0; j < ioData.NCols(); j++ {
			if ioData.At(j).IsNull(i) {
				buff += naText + "|"
			} else {
				buff += ioData.At(j).GetAsString(i) + "|"
			}
		}
		buff += "\n"
	}

	_, err := writer.Write([]byte(buff))
	if err != nil {
		return err
	}

	return nil
}
