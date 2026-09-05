package dataframe

import (
	"strings"
	"testing"
)

func selectTestFrame() DataFrame {
	return NewBaseDataFrame(testCtx).
		AddSeriesFromStrings("Car", []string{"a", "b"}, nil, false).
		AddSeriesFromStrings("CarOrigin", []string{"a - US", "b - EU"}, nil, false).
		AddSeriesFromStrings("Origin", []string{"US", "EU"}, nil, false).
		AddSeriesFromInt64s("Stat", []int64{1, 2}, nil, false)
}

// The point of splitting the API: a plain column name selects that column and
// nothing else. "Car" used to drag "CarOrigin" along with it.
func TestSelect_ExactNamesDoNotOverSelect(t *testing.T) {
	got := selectTestFrame().Select("Car", "Origin", "Stat")
	if got.GetError() != nil {
		t.Fatalf("unexpected error: %v", got.GetError())
	}

	want := []string{"Car", "Origin", "Stat"}
	if names := got.Names(); !equalStrings(names, want) {
		t.Fatalf("got %v, want %v", names, want)
	}
}

// Selector order is preserved, and a repeated name is kept once.
func TestSelect_OrderAndDuplicates(t *testing.T) {
	got := selectTestFrame().Select("Stat", "Car", "Stat")
	if got.GetError() != nil {
		t.Fatalf("unexpected error: %v", got.GetError())
	}

	want := []string{"Stat", "Car"}
	if names := got.Names(); !equalStrings(names, want) {
		t.Fatalf("got %v, want %v", names, want)
	}
}

// A typo must be an error rather than a silently missing column.
func TestSelect_UnknownColumnIsAnError(t *testing.T) {
	got := selectTestFrame().Select("Car", "Orgin")
	if got.GetError() == nil {
		t.Fatal("expected an error for an unknown column name, got nil")
	}
	if msg := got.GetError().Error(); !strings.Contains(msg, "Orgin") {
		t.Fatalf("error should name the missing column, got %q", msg)
	}
}

// An exact name is not treated as a pattern.
func TestSelect_NameIsNotAPattern(t *testing.T) {
	got := selectTestFrame().Select("Car.")
	if got.GetError() == nil {
		t.Fatal("expected \"Car.\" to be an unknown column, not a pattern matching CarOrigin")
	}
}

func TestSelectMatching_StillMatchesPatterns(t *testing.T) {
	got := selectTestFrame().SelectMatching("^Car")
	if got.GetError() != nil {
		t.Fatalf("unexpected error: %v", got.GetError())
	}

	want := []string{"Car", "CarOrigin"}
	if names := got.Names(); !equalStrings(names, want) {
		t.Fatalf("got %v, want %v", names, want)
	}
}

// A pattern matching nothing is empty, not an error: a pattern is not a claim
// that any particular column exists.
func TestSelectMatching_NoMatchIsNotAnError(t *testing.T) {
	got := selectTestFrame().SelectMatching("^nothing$")
	if got.GetError() != nil {
		t.Fatalf("unexpected error: %v", got.GetError())
	}
	if n := len(got.Names()); n != 0 {
		t.Fatalf("expected no columns, got %d", n)
	}
}

func TestSelectMatching_InvalidPatternIsAnError(t *testing.T) {
	got := selectTestFrame().SelectMatching("^(unclosed")
	if got.GetError() == nil {
		t.Fatal("expected an error for an invalid regular expression, got nil")
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
