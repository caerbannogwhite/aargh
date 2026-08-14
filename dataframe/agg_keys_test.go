package dataframe

import (
	"testing"

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
