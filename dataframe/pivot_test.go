package dataframe

import (
	"testing"

	"github.com/caerbannogwhite/aargh"
)

func TestPivotLonger(t *testing.T) {
	ctx := aargh.NewContext()

	// Create a wide dataframe for testing
	df := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("id", []string{"A", "B", "C"}, nil, false).
		AddSeriesFromInts("x", []int{1, 2, 3}, nil, false).
		AddSeriesFromInts("y", []int{4, 5, 6}, nil, false).
		AddSeriesFromInts("z", []int{7, 8, 9}, nil, false).(BaseDataFrame)

	// Test basic pivot_longer
	params := NewPivotLongerParams().
		SetCols("x", "y", "z").
		SetNamesTo("variable").
		SetValuesTo("value")

	result := df.PivotLonger(params)

	if result.IsErrored() {
		t.Fatalf("PivotLonger failed: %v", result.GetError())
	}

	// Check dimensions
	expectedRows := 9 // 3 rows * 3 pivot columns
	expectedCols := 3 // id + variable + value
	if result.NRows() != expectedRows {
		t.Errorf("Expected %d rows, got %d", expectedRows, result.NRows())
	}
	if result.NCols() != expectedCols {
		t.Errorf("Expected %d columns, got %d", expectedCols, result.NCols())
	}

	// Check column names
	names := result.Names()
	expectedNames := []string{"id", "variable", "value"}
	for i, expected := range expectedNames {
		if names[i] != expected {
			t.Errorf("Expected column %d to be '%s', got '%s'", i, expected, names[i])
		}
	}

	// Check first few values
	idColumn := result.C("id")
	varColumn := result.C("variable")
	valueColumn := result.C("value")

	// First row should be: id="A", variable="x", value="1"
	if idColumn.GetAsString(0) != "A" {
		t.Errorf("Expected first id to be 'A', got '%s'", idColumn.GetAsString(0))
	}
	if varColumn.GetAsString(0) != "x" {
		t.Errorf("Expected first variable to be 'x', got '%s'", varColumn.GetAsString(0))
	}
	if valueColumn.GetAsString(0) != "1" {
		t.Errorf("Expected first value to be '1', got '%s'", valueColumn.GetAsString(0))
	}
}

func TestPivotLongerWithPrefix(t *testing.T) {
	ctx := aargh.NewContext()

	// Create a dataframe with prefixed columns
	df := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("id", []string{"A", "B"}, nil, false).
		AddSeriesFromInts("measure_x", []int{1, 2}, nil, false).
		AddSeriesFromInts("measure_y", []int{3, 4}, nil, false).(BaseDataFrame)

	params := NewPivotLongerParams().
		SetCols("measure_x", "measure_y").
		SetNamesTo("type").
		SetValuesTo("measurement").
		SetNamesPrefix("measure_")

	result := df.PivotLonger(params)

	if result.IsErrored() {
		t.Fatalf("PivotLonger with prefix failed: %v", result.GetError())
	}

	// Check that prefix was stripped
	typeColumn := result.C("type")
	if typeColumn.GetAsString(0) != "x" {
		t.Errorf("Expected prefix to be stripped, got '%s'", typeColumn.GetAsString(0))
	}
}

func TestPivotWider(t *testing.T) {
	ctx := aargh.NewContext()

	// Create a long dataframe for testing (similar to result of pivot_longer)
	df := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("id", []string{"A", "A", "A", "B", "B", "B"}, nil, false).
		AddSeriesFromStrings("variable", []string{"x", "y", "z", "x", "y", "z"}, nil, false).
		AddSeriesFromInts("value", []int{1, 4, 7, 2, 5, 8}, nil, false).(BaseDataFrame)

	params := NewPivotWiderParams().
		SetNamesFrom("variable").
		SetValuesFrom("value")

	result := df.PivotWider(params)

	if result.IsErrored() {
		t.Fatalf("PivotWider failed: %v", result.GetError())
	}

	// Check dimensions
	expectedRows := 2 // unique id values
	expectedCols := 4 // id + x + y + z
	if result.NRows() != expectedRows {
		t.Errorf("Expected %d rows, got %d", expectedRows, result.NRows())
	}
	if result.NCols() != expectedCols {
		t.Errorf("Expected %d columns, got %d", expectedCols, result.NCols())
	}

	// Check that we have the expected columns
	names := result.Names()
	hasId := false
	hasX := false
	hasY := false
	hasZ := false

	for _, name := range names {
		switch name {
		case "id":
			hasId = true
		case "x":
			hasX = true
		case "y":
			hasY = true
		case "z":
			hasZ = true
		}
	}

	if !hasId || !hasX || !hasY || !hasZ {
		t.Errorf("Missing expected columns in result: %v", names)
	}
}

func TestPivotWiderWithIdCols(t *testing.T) {
	ctx := aargh.NewContext()

	// Create a dataframe with multiple ID columns
	df := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("group", []string{"G1", "G1", "G2", "G2"}, nil, false).
		AddSeriesFromStrings("id", []string{"A", "A", "B", "B"}, nil, false).
		AddSeriesFromStrings("variable", []string{"x", "y", "x", "y"}, nil, false).
		AddSeriesFromInts("value", []int{1, 2, 3, 4}, nil, false).(BaseDataFrame)

	params := NewPivotWiderParams().
		SetIdCols("group", "id").
		SetNamesFrom("variable").
		SetValuesFrom("value")

	result := df.PivotWider(params)

	if result.IsErrored() {
		t.Fatalf("PivotWider with ID cols failed: %v", result.GetError())
	}

	// Should have 2 rows (unique group+id combinations)
	expectedRows := 2
	if result.NRows() != expectedRows {
		t.Errorf("Expected %d rows, got %d", expectedRows, result.NRows())
	}
}

func TestPivotWiderWithPrefix(t *testing.T) {
	ctx := aargh.NewContext()

	df := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("id", []string{"A", "A", "B", "B"}, nil, false).
		AddSeriesFromStrings("measure", []string{"height", "weight", "height", "weight"}, nil, false).
		AddSeriesFromFloat64s("value", []float64{170.5, 65.2, 180.1, 75.8}, nil, false).(BaseDataFrame)

	params := NewPivotWiderParams().
		SetNamesFrom("measure").
		SetValuesFrom("value").
		SetNamesPrefix("m_")

	result := df.PivotWider(params)

	if result.IsErrored() {
		t.Fatalf("PivotWider with prefix failed: %v", result.GetError())
	}

	// Check that columns have the expected prefix
	names := result.Names()
	hasPrefix := false
	for _, name := range names {
		if name == "m_height" || name == "m_weight" {
			hasPrefix = true
			break
		}
	}

	if !hasPrefix {
		t.Errorf("Expected columns with prefix 'm_', got: %v", names)
	}
}

func TestPivotWiderWithFill(t *testing.T) {
	ctx := aargh.NewContext()

	// Create data with missing combinations
	df := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("id", []string{"A", "B"}, nil, false).
		AddSeriesFromStrings("variable", []string{"x", "y"}, nil, false).
		AddSeriesFromInts("value", []int{1, 2}, nil, false).(BaseDataFrame)

	params := NewPivotWiderParams().
		SetNamesFrom("variable").
		SetValuesFrom("value").
		SetValuesFill(0)

	result := df.PivotWider(params)

	if result.IsErrored() {
		t.Fatalf("PivotWider with fill failed: %v", result.GetError())
	}

	// Should fill missing values with 0
	expectedRows := 2
	if result.NRows() != expectedRows {
		t.Errorf("Expected %d rows, got %d", expectedRows, result.NRows())
	}
}

func TestPivotLongerErrors(t *testing.T) {
	ctx := aargh.NewContext()

	df := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("id", []string{"A", "B"}, nil, false).
		AddSeriesFromInts("x", []int{1, 2}, nil, false).(BaseDataFrame)

	// Test with no columns specified
	params := NewPivotLongerParams()
	result := df.PivotLonger(params)

	if !result.IsErrored() {
		t.Error("Expected error when no columns specified for pivot_longer")
	}

	// Test with non-existent column
	params = NewPivotLongerParams().SetCols("nonexistent")
	result = df.PivotLonger(params)

	if !result.IsErrored() {
		t.Error("Expected error when specifying non-existent column")
	}
}

func TestPivotWiderErrors(t *testing.T) {
	ctx := aargh.NewContext()

	df := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("id", []string{"A", "B"}, nil, false).
		AddSeriesFromStrings("variable", []string{"x", "y"}, nil, false).
		AddSeriesFromInts("value", []int{1, 2}, nil, false).(BaseDataFrame)

	// Test with missing names_from
	params := NewPivotWiderParams().SetValuesFrom("value")
	result := df.PivotWider(params)

	if !result.IsErrored() {
		t.Error("Expected error when names_from not specified")
	}

	// Test with missing values_from
	params = NewPivotWiderParams().SetNamesFrom("variable")
	result = df.PivotWider(params)

	if !result.IsErrored() {
		t.Error("Expected error when values_from not specified")
	}

	// Test with non-existent names_from column
	params = NewPivotWiderParams().
		SetNamesFrom("nonexistent").
		SetValuesFrom("value")
	result = df.PivotWider(params)

	if !result.IsErrored() {
		t.Error("Expected error when names_from column doesn't exist")
	}

	// Test with non-existent values_from column
	params = NewPivotWiderParams().
		SetNamesFrom("variable").
		SetValuesFrom("nonexistent")
	result = df.PivotWider(params)

	if !result.IsErrored() {
		t.Error("Expected error when values_from column doesn't exist")
	}
}

func TestPivotRoundTrip(t *testing.T) {
	ctx := aargh.NewContext()

	// Create original wide data
	original := NewBaseDataFrame(ctx).
		AddSeriesFromStrings("id", []string{"A", "B", "C"}, nil, false).
		AddSeriesFromInts("x", []int{1, 2, 3}, nil, false).
		AddSeriesFromInts("y", []int{4, 5, 6}, nil, false).(BaseDataFrame)

	// Pivot to long format
	longerParams := NewPivotLongerParams().
		SetCols("x", "y").
		SetNamesTo("variable").
		SetValuesTo("value")

	longer := original.PivotLonger(longerParams)

	if longer.IsErrored() {
		t.Fatalf("PivotLonger failed: %v", longer.GetError())
	}

	// Pivot back to wide format
	widerParams := NewPivotWiderParams().
		SetNamesFrom("variable").
		SetValuesFrom("value")

	result := longer.(BaseDataFrame).PivotWider(widerParams)

	if result.IsErrored() {
		t.Fatalf("PivotWider failed: %v", result.GetError())
	}

	// Check that we got back to the same dimensions
	if result.NRows() != original.NRows() {
		t.Errorf("Expected %d rows after round trip, got %d", original.NRows(), result.NRows())
	}

	// Note: Column order might be different, so we check that all expected columns exist
	resultNames := result.Names()
	expectedNames := []string{"id", "x", "y"}

	for _, expected := range expectedNames {
		found := false
		for _, actual := range resultNames {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected column '%s' not found after round trip", expected)
		}
	}
}
