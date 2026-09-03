package io

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/caerbannogwhite/enchanter/series"
)

// failingWriter fails every write, standing in for a full disk or a
// disconnected network share.
type failingWriter struct{ err error }

func (f *failingWriter) Write(p []byte) (int, error) { return 0, f.err }

func oneCellIoData() *IoData {
	iod := NewIoData(ctx)
	iod.AddSeries(series.NewSeriesString([]string{"fresh"}, nil, false, ctx), SeriesMeta{Name: "value"})
	return iod
}

// The CSV writer used to discard the error from every write, so a write that
// never reached the destination was reported as a success.
func Test_IoCsv_WriteReportsWriteErrors(t *testing.T) {
	want := errors.New("no space left on device")

	err := oneCellIoData().ToCsv().SetWriter(&failingWriter{err: want}).Write()
	if err == nil {
		t.Fatal("expected an error when every write fails, got nil")
	}
	if !errors.Is(err, want) {
		t.Fatalf("expected the underlying write error, got %v", err)
	}
}

// A configuration error must be caught before the destination is opened:
// opening truncates, so validating late would destroy the previous contents
// on behalf of a write that could never succeed.
func Test_IoCsv_InvalidQuotingLeavesFileUntouched(t *testing.T) {
	path := filepath.Join(t.TempDir(), "keep.csv")
	original := "irreplaceable,production,data\n1,2,3\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}

	err := oneCellIoData().ToCsv().SetPath(path).SetQuoting(CsvQuotingType(99)).Write()
	if err == nil {
		t.Fatal("expected an error for an invalid quoting type, got nil")
	}

	got, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != original {
		t.Fatalf("file was modified despite the error:\n got %q\nwant %q", got, original)
	}
}

// Every writer that takes a path must truncate: writing a small dataset over a
// larger file previously left the tail of the old one in place, which is
// invalid output for each of these formats.
func Test_IoWriters_TruncateOnOverwrite(t *testing.T) {
	const staleMarker = "STALE-MUST-NOT-SURVIVE"
	stale := []byte(strings.Repeat(staleMarker+"\n", 512))

	writers := []struct {
		name  string
		ext   string
		write func(iod *IoData, path string) error
	}{
		{"csv", ".csv", func(iod *IoData, p string) error {
			return iod.ToCsv().SetPath(p).SetEol("\n").Write()
		}},
		{"json", ".json", func(iod *IoData, p string) error {
			return iod.ToJson().SetPath(p).Write()
		}},
		{"markdown", ".md", func(iod *IoData, p string) error {
			return iod.ToMarkdown().SetPath(p).Write()
		}},
		{"html", ".html", func(iod *IoData, p string) error {
			return iod.ToHtml().SetPath(p).Write()
		}},
	}

	for _, w := range writers {
		t.Run(w.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "out"+w.ext)
			if err := os.WriteFile(path, stale, 0600); err != nil {
				t.Fatal(err)
			}

			if err := w.write(oneCellIoData(), path); err != nil {
				t.Fatalf("write: %v", err)
			}

			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(got), staleMarker) {
				t.Fatalf("previous contents survived the overwrite: %q", truncateForMsg(string(got)))
			}
			if len(got) >= len(stale) {
				t.Fatalf("file was not truncated: %d bytes written over %d", len(got), len(stale))
			}
		})
	}
}

func truncateForMsg(s string) string {
	if len(s) > 160 {
		return s[:160] + "..."
	}
	return s
}
