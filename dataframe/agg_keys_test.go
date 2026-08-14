package dataframe

import (
	"testing"
	"time"

	"github.com/caerbannogwhite/enchanter"
	"github.com/caerbannogwhite/enchanter/series"
)

func TestGroupTableSingleAndMulti(t *testing.T) {
	ctx := enchanter.NewContext()
	// single string key: a,b,a,a,b -> 2 groups, ids stable by first appearance
	k := series.NewSeriesString([]string{"a", "b", "a", "a", "b"}, nil, false, ctx)
	g := newGroupTable([]series.Series{k})
	ids := make([]int, 5)
	for i := 0; i < 5; i++ {
		ids[i] = g.idOf(i)
	}
	if g.numGroups() != 2 {
		t.Fatalf("groups = %d, want 2", g.numGroups())
	}
	if !(ids[0] == ids[2] && ids[0] == ids[3] && ids[1] == ids[4] && ids[0] != ids[1]) {
		t.Fatalf("bad id assignment: %v", ids)
	}

	// multi-key: (int, string) — exact composite equality, no collisions
	ci := series.NewSeriesInt64([]int64{1, 1, 2, 1}, nil, false, ctx)
	cs := series.NewSeriesString([]string{"x", "y", "x", "x"}, nil, false, ctx)
	g2 := newGroupTable([]series.Series{ci, cs})
	m := []int{g2.idOf(0), g2.idOf(1), g2.idOf(2), g2.idOf(3)}
	// rows 0 and 3 are (1,"x"); 1 is (1,"y"); 2 is (2,"x")
	if !(m[0] == m[3] && m[0] != m[1] && m[0] != m[2] && m[1] != m[2]) {
		t.Fatalf("multi-key ids wrong: %v", m)
	}
	if g2.numGroups() != 3 {
		t.Fatalf("multi groups = %d, want 3", g2.numGroups())
	}
}

func TestGroupTableNullKeys(t *testing.T) {
	ctx := enchanter.NewContext()
	// nulls at rows 1 and 3 form ONE group distinct from non-null "a"
	k := series.NewSeriesString([]string{"a", "", "a", ""}, []bool{false, true, false, true}, false, ctx)
	g := newGroupTable([]series.Series{k})
	ids := []int{g.idOf(0), g.idOf(1), g.idOf(2), g.idOf(3)}
	if !(ids[0] == ids[2] && ids[1] == ids[3] && ids[0] != ids[1]) {
		t.Fatalf("null-key grouping wrong: %v", ids)
	}
	if g.numGroups() != 2 {
		t.Fatalf("groups = %d, want 2", g.numGroups())
	}
}

func TestGroupTableTimesSameInstantDifferentLocation(t *testing.T) {
	ctx := enchanter.NewContext()
	// Same instant, different *Location: raw time.Time equality (what a
	// map[time.Time]key would use) compares wall/monotonic reading AND
	// Location, so tm == tm2 is false even though tm.Equal(tm2) is true.
	// The group-key table must key on instant equality (UnixNano), not raw
	// time.Time equality, so these two rows must land in ONE group.
	tm := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	tm2 := tm.In(time.FixedZone("plus", 0))
	k := series.NewSeriesTime([]time.Time{tm, tm2}, nil, false, ctx)
	g := newGroupTable([]series.Series{k})
	id0 := g.idOf(0)
	id1 := g.idOf(1)
	if id0 != id1 {
		t.Fatalf("same-instant Times with different Location got different ids: %d vs %d", id0, id1)
	}
	if g.numGroups() != 1 {
		t.Fatalf("groups = %d, want 1", g.numGroups())
	}
}
