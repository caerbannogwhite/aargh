package dataframe

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/caerbannogwhite/aargh"
)

const (
	NA_TEXT = aargh.NA_TEXT
)

var ctx *aargh.Context
var testDataDir string

func TestMain(m *testing.M) {
	ctx = aargh.NewContext()
	testDataDir = filepath.Join("..", "testdata")

	flag.Parse()

	// The G1 fixtures are large local-only files (up to ~500MB, gitignored)
	// used exclusively by benchmarks, which b.Skip when a frame is nil.
	// Loading them takes ~40s and gigabytes of memory, so skip in -short mode.
	if !testing.Short() {
		read_G1_1e4_1e2_0_0()
		read_G1_1e5_1e2_0_0()
		read_G1_1e6_1e2_0_0()
		read_G1_1e7_1e2_0_0()
		read_G1_1e4_1e2_10_0()
		read_G1_1e5_1e2_10_0()
		read_G1_1e6_1e2_10_0()
		read_G1_1e7_1e2_10_0()
	}

	os.Exit(m.Run())
}
