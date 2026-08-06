package series

import "github.com/caerbannogwhite/enchanter"

const (
	NA_TEXT = enchanter.NA_TEXT
)

var ctx *enchanter.Context

func init() {
	ctx = enchanter.NewContext()
}
