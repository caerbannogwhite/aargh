package dataframe

import (
	"encoding/binary"
	"time"

	"github.com/caerbannogwhite/enchanter/series"
)

// groupTable assigns dense, stable group ids to rows based on the exact
// composite value of one or more key columns. Two rows map to the same id
// iff every key column has an equal value for those rows, with null==null
// per column (nulls in a given column group together, distinct from every
// non-null value in that column).
type groupTable struct {
	coders []func(row int) uint64 // one per key column
	ids    map[string]int         // packed composite code -> dense group id
	reps   []int                  // representative row per group id (index = gid)
	buf    []byte                 // reused scratch buffer for packing per-row codes
}

// newGroupTable builds a groupTable over the given key columns. Each column
// contributes one dense per-row code via makeCellCoder; idOf packs the
// per-column codes for a row into an exact composite key.
func newGroupTable(keyCols []series.Series) *groupTable {
	g := &groupTable{
		coders: make([]func(int) uint64, len(keyCols)),
		ids:    make(map[string]int, 128),
		buf:    make([]byte, 8*len(keyCols)),
	}
	for i, col := range keyCols {
		g.coders[i] = makeCellCoder(col)
	}
	return g
}

// idOf returns the dense group id for row, assigning a new id (and
// capturing row as that group's representative) the first time this exact
// composite key is seen.
func (g *groupTable) idOf(row int) int {
	for i, code := range g.coders {
		binary.LittleEndian.PutUint64(g.buf[i*8:], code(row))
	}
	key := string(g.buf) // copies bytes; g.buf is safely reused after this
	if id, ok := g.ids[key]; ok {
		return id
	}
	id := len(g.reps)
	g.ids[key] = id
	g.reps = append(g.reps, row)
	return id
}

// numGroups returns the number of distinct groups seen so far.
func (g *groupTable) numGroups() int { return len(g.reps) }

// representativeRows returns the representative row index for each group,
// indexed by group id.
func (g *groupTable) representativeRows() []int { return g.reps }

// makeCellCoder returns a per-row dense code for one key column: 0 means
// null, distinct non-null values get stable codes 1..k in order of first
// appearance. The type switch covers the same concrete series types as the
// former BaseDataFrame.groupHelper.
//
// For interned Strings the map is keyed by the *string pointer (cheap,
// pointer-identity is stable within a StringPool); the other supported key
// types are keyed by their typed value in a dedicated map to avoid `any`
// boxing on the hot path. Unrecognized series types fall back to the
// generic Get(row) any accessor exposed by series.Series.
func makeCellCoder(col series.Series) func(int) uint64 {
	switch c := col.(type) {
	case series.Strings:
		codes := make(map[*string]uint64, 16)
		return func(row int) uint64 {
			if c.IsNull(row) {
				return 0
			}
			p := c.Data_[row]
			if v, ok := codes[p]; ok {
				return v
			}
			v := uint64(len(codes)) + 1
			codes[p] = v
			return v
		}

	case series.Bools:
		codes := make(map[bool]uint64, 2)
		return func(row int) uint64 {
			if c.IsNull(row) {
				return 0
			}
			cell := c.Data_[row]
			if v, ok := codes[cell]; ok {
				return v
			}
			v := uint64(len(codes)) + 1
			codes[cell] = v
			return v
		}

	case series.Ints:
		codes := make(map[int]uint64, 16)
		return func(row int) uint64 {
			if c.IsNull(row) {
				return 0
			}
			cell := c.Data_[row]
			if v, ok := codes[cell]; ok {
				return v
			}
			v := uint64(len(codes)) + 1
			codes[cell] = v
			return v
		}

	case series.Int64s:
		codes := make(map[int64]uint64, 16)
		return func(row int) uint64 {
			if c.IsNull(row) {
				return 0
			}
			cell := c.Data_[row]
			if v, ok := codes[cell]; ok {
				return v
			}
			v := uint64(len(codes)) + 1
			codes[cell] = v
			return v
		}

	case series.Float64s:
		codes := make(map[float64]uint64, 16)
		return func(row int) uint64 {
			if c.IsNull(row) {
				return 0
			}
			cell := c.Data_[row]
			if v, ok := codes[cell]; ok {
				return v
			}
			v := uint64(len(codes)) + 1
			codes[cell] = v
			return v
		}

	case series.Times:
		// Key by UnixNano (exact instant), not the raw time.Time value:
		// time.Time equality via == also compares the monotonic reading and
		// *Location pointer, so two values representing the same instant
		// (t.Equal(t2) == true) can otherwise land in different groups.
		// Mirrors series.Times.Group() (series/time.go), which also keys by
		// UnixNano().
		codes := make(map[int64]uint64, 16)
		return func(row int) uint64 {
			if c.IsNull(row) {
				return 0
			}
			cell := c.Data_[row].UnixNano()
			if v, ok := codes[cell]; ok {
				return v
			}
			v := uint64(len(codes)) + 1
			codes[cell] = v
			return v
		}

	case series.Durations:
		codes := make(map[time.Duration]uint64, 16)
		return func(row int) uint64 {
			if c.IsNull(row) {
				return 0
			}
			cell := c.Data_[row]
			if v, ok := codes[cell]; ok {
				return v
			}
			v := uint64(len(codes)) + 1
			codes[cell] = v
			return v
		}

	default:
		// Unrecognized series type: fall back to the generic any accessor.
		codes := make(map[any]uint64, 16)
		return func(row int) uint64 {
			if col.IsNull(row) {
				return 0
			}
			cell := col.Get(row)
			if v, ok := codes[cell]; ok {
				return v
			}
			v := uint64(len(codes)) + 1
			codes[cell] = v
			return v
		}
	}
}
