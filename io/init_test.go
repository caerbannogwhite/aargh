package io

import (
	"path/filepath"

	"github.com/caerbannogwhite/enchanter"
)

const (
	NA_TEXT = enchanter.NA_TEXT
)

var ctx *enchanter.Context
var testDataFolder = filepath.Join("..", "testdata")

func init() {
	ctx = enchanter.NewContext()
}
