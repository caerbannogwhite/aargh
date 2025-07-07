package dataframe

import (
	"fmt"
	"strings"

	"github.com/caerbannogwhite/aargh/series"
)

// PivotLongerParams contains parameters for pivot_longer operation
type PivotLongerParams struct {
	// Columns to pivot into longer format. Can be column names or indices.
	Cols []string
	// Name of the column to create from the column names of cols
	NamesTo string
	// Name of the column to create from the data values in cols
	ValuesTo string
	// Prefix to strip from column names when creating the names column
	NamesPrefix string
	// Separator to split column names on
	NamesSep string
}

// PivotWiderParams contains parameters for pivot_wider operation
type PivotWiderParams struct {
	// Column(s) that uniquely identify each observation
	IdCols []string
	// Column name whose values will become column names in the output
	NamesFrom string
	// Column name whose values will populate the new columns
	ValuesFrom string
	// String to add to the start of every variable name
	NamesPrefix string
	// Separator to use between names_from values
	NamesSep string
	// Value to fill in missing combinations
	ValuesFill interface{}
}

// NewPivotLongerParams creates default parameters for pivot_longer
func NewPivotLongerParams() PivotLongerParams {
	return PivotLongerParams{
		NamesTo:  "name",
		ValuesTo: "value",
		NamesSep: "_",
	}
}

// NewPivotWiderParams creates default parameters for pivot_wider
func NewPivotWiderParams() PivotWiderParams {
	return PivotWiderParams{
		NamesSep: "_",
	}
}

// SetCols sets the columns to pivot for pivot_longer
func (p PivotLongerParams) SetCols(cols ...string) PivotLongerParams {
	p.Cols = cols
	return p
}

// SetNamesTo sets the name of the column for column names
func (p PivotLongerParams) SetNamesTo(name string) PivotLongerParams {
	p.NamesTo = name
	return p
}

// SetValuesTo sets the name of the column for values
func (p PivotLongerParams) SetValuesTo(name string) PivotLongerParams {
	p.ValuesTo = name
	return p
}

// SetNamesPrefix sets the prefix to strip from column names
func (p PivotLongerParams) SetNamesPrefix(prefix string) PivotLongerParams {
	p.NamesPrefix = prefix
	return p
}

// SetNamesSep sets the separator for column names
func (p PivotLongerParams) SetNamesSep(sep string) PivotLongerParams {
	p.NamesSep = sep
	return p
}

// SetIdCols sets the ID columns for pivot_wider
func (p PivotWiderParams) SetIdCols(cols ...string) PivotWiderParams {
	p.IdCols = cols
	return p
}

// SetNamesFrom sets the column for new column names
func (p PivotWiderParams) SetNamesFrom(name string) PivotWiderParams {
	p.NamesFrom = name
	return p
}

// SetValuesFrom sets the column for values
func (p PivotWiderParams) SetValuesFrom(name string) PivotWiderParams {
	p.ValuesFrom = name
	return p
}

// SetNamesPrefix sets the prefix for new column names
func (p PivotWiderParams) SetNamesPrefix(prefix string) PivotWiderParams {
	p.NamesPrefix = prefix
	return p
}

// SetNamesSep sets the separator for new column names
func (p PivotWiderParams) SetNamesSep(sep string) PivotWiderParams {
	p.NamesSep = sep
	return p
}

// SetValuesFill sets the fill value for missing combinations
func (p PivotWiderParams) SetValuesFill(fill interface{}) PivotWiderParams {
	p.ValuesFill = fill
	return p
}

// PivotLonger transforms data from wide to long format
func (df BaseDataFrame) PivotLonger(params PivotLongerParams) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.PivotLonger: cannot pivot a grouped dataframe")
		return df
	}

	if len(params.Cols) == 0 {
		df.err = fmt.Errorf("BaseDataFrame.PivotLonger: must specify columns to pivot")
		return df
	}

	// Validate column names exist
	pivotIndices := make([]int, len(params.Cols))
	for i, colName := range params.Cols {
		idx := df.GetSeriesIndex(colName)
		if idx == -1 {
			df.err = fmt.Errorf("BaseDataFrame.PivotLonger: column '%s' not found", colName)
			return df
		}
		pivotIndices[i] = idx
	}

	// Identify non-pivot columns (ID columns)
	idIndices := make([]int, 0)
	for i := range df.names {
		isPivot := false
		for _, pivotIdx := range pivotIndices {
			if i == pivotIdx {
				isPivot = true
				break
			}
		}
		if !isPivot {
			idIndices = append(idIndices, i)
		}
	}

	nRows := df.NRows()
	nPivotCols := len(params.Cols)
	newNRows := nRows * nPivotCols

	// Create result dataframe
	result := NewBaseDataFrame(df.ctx).(BaseDataFrame)

	// Add ID columns (repeated for each pivot column)
	for _, idIdx := range idIndices {
		originalSeries := df.series[idIdx]
		newData := make([]string, newNRows)
		newNullMask := make([]bool, newNRows)

		for i := 0; i < nRows; i++ {
			for j := 0; j < nPivotCols; j++ {
				newIdx := i*nPivotCols + j
				newData[newIdx] = originalSeries.GetAsString(i)
				newNullMask[newIdx] = originalSeries.IsNull(i)
			}
		}

		// Convert back to original type
		switch originalSeries.Type().String() {
		case "Bool":
			boolData := make([]bool, newNRows)
			for i, str := range newData {
				if str == "true" {
					boolData[i] = true
				}
			}
			result = result.AddSeries(df.names[idIdx], series.NewSeriesBool(boolData, newNullMask, false, df.ctx)).(BaseDataFrame)
		case "Int":
			intData := make([]int, newNRows)
			for i, str := range newData {
				if !newNullMask[i] && str != "" {
					fmt.Sscanf(str, "%d", &intData[i])
				}
			}
			result = result.AddSeries(df.names[idIdx], series.NewSeriesInt(intData, newNullMask, false, df.ctx)).(BaseDataFrame)
		case "Int64":
			int64Data := make([]int64, newNRows)
			for i, str := range newData {
				if !newNullMask[i] && str != "" {
					fmt.Sscanf(str, "%d", &int64Data[i])
				}
			}
			result = result.AddSeries(df.names[idIdx], series.NewSeriesInt64(int64Data, newNullMask, false, df.ctx)).(BaseDataFrame)
		case "Float64":
			float64Data := make([]float64, newNRows)
			for i, str := range newData {
				if !newNullMask[i] && str != "" {
					fmt.Sscanf(str, "%f", &float64Data[i])
				}
			}
			result = result.AddSeries(df.names[idIdx], series.NewSeriesFloat64(float64Data, newNullMask, false, df.ctx)).(BaseDataFrame)
		default:
			result = result.AddSeries(df.names[idIdx], series.NewSeriesString(newData, newNullMask, false, df.ctx)).(BaseDataFrame)
		}
	}

	// Add names column
	namesData := make([]string, newNRows)
	for i := 0; i < nRows; i++ {
		for j, colName := range params.Cols {
			newIdx := i*nPivotCols + j
			processedName := colName
			if params.NamesPrefix != "" {
				processedName = strings.TrimPrefix(processedName, params.NamesPrefix)
			}
			namesData[newIdx] = processedName
		}
	}
	result = result.AddSeries(params.NamesTo, series.NewSeriesString(namesData, nil, false, df.ctx)).(BaseDataFrame)

	// Add values column - combine all pivot columns into one
	valuesData := make([]string, newNRows)
	valuesNullMask := make([]bool, newNRows)
	for i := 0; i < nRows; i++ {
		for j, pivotIdx := range pivotIndices {
			newIdx := i*nPivotCols + j
			valuesData[newIdx] = df.series[pivotIdx].GetAsString(i)
			valuesNullMask[newIdx] = df.series[pivotIdx].IsNull(i)
		}
	}
	result = result.AddSeries(params.ValuesTo, series.NewSeriesString(valuesData, valuesNullMask, false, df.ctx)).(BaseDataFrame)

	return result
}

// PivotWider transforms data from long to wide format
func (df BaseDataFrame) PivotWider(params PivotWiderParams) DataFrame {
	if df.err != nil {
		return df
	}

	if df.isGrouped {
		df.err = fmt.Errorf("BaseDataFrame.PivotWider: cannot pivot a grouped dataframe")
		return df
	}

	if params.NamesFrom == "" {
		df.err = fmt.Errorf("BaseDataFrame.PivotWider: must specify names_from column")
		return df
	}

	if params.ValuesFrom == "" {
		df.err = fmt.Errorf("BaseDataFrame.PivotWider: must specify values_from column")
		return df
	}

	// Find column indices
	namesFromIdx := df.GetSeriesIndex(params.NamesFrom)
	if namesFromIdx == -1 {
		df.err = fmt.Errorf("BaseDataFrame.PivotWider: names_from column '%s' not found", params.NamesFrom)
		return df
	}

	valuesFromIdx := df.GetSeriesIndex(params.ValuesFrom)
	if valuesFromIdx == -1 {
		df.err = fmt.Errorf("BaseDataFrame.PivotWider: values_from column '%s' not found", params.ValuesFrom)
		return df
	}

	// Determine ID columns
	var idIndices []int
	if len(params.IdCols) > 0 {
		// Use specified ID columns
		idIndices = make([]int, len(params.IdCols))
		for i, colName := range params.IdCols {
			idx := df.GetSeriesIndex(colName)
			if idx == -1 {
				df.err = fmt.Errorf("BaseDataFrame.PivotWider: id column '%s' not found", colName)
				return df
			}
			idIndices[i] = idx
		}
	} else {
		// Use all columns except names_from and values_from as ID columns
		for i, colName := range df.names {
			if colName != params.NamesFrom && colName != params.ValuesFrom {
				idIndices = append(idIndices, i)
			}
		}
	}

	// Get unique values from names_from column
	namesFromSeries := df.series[namesFromIdx]
	uniqueNames := make(map[string]bool)
	for i := 0; i < namesFromSeries.Len(); i++ {
		if !namesFromSeries.IsNull(i) {
			colName := namesFromSeries.GetAsString(i)
			if params.NamesPrefix != "" {
				colName = params.NamesPrefix + colName
			}
			uniqueNames[colName] = true
		}
	}

	// Convert to sorted slice for consistent ordering
	newColNames := make([]string, 0, len(uniqueNames))
	for colName := range uniqueNames {
		newColNames = append(newColNames, colName)
	}

	// Create a map for quick lookup of unique ID combinations
	idGroups := make(map[string][]int) // key -> row indices
	idOrder := make([]string, 0)       // to maintain order

	for i := 0; i < df.NRows(); i++ {
		keyParts := make([]string, len(idIndices))
		for j, idIdx := range idIndices {
			keyParts[j] = df.series[idIdx].GetAsString(i)
		}
		key := strings.Join(keyParts, "|||") // Use unlikely separator

		if _, exists := idGroups[key]; !exists {
			idOrder = append(idOrder, key)
			idGroups[key] = make([]int, 0)
		}
		idGroups[key] = append(idGroups[key], i)
	}

	// Create result dataframe
	result := NewBaseDataFrame(df.ctx).(BaseDataFrame)
	newNRows := len(idOrder)

	// Add ID columns to result
	for _, idIdx := range idIndices {
		idData := make([]string, newNRows)
		idNullMask := make([]bool, newNRows)

		for i, key := range idOrder {
			firstRowIdx := idGroups[key][0]
			idData[i] = df.series[idIdx].GetAsString(firstRowIdx)
			idNullMask[i] = df.series[idIdx].IsNull(firstRowIdx)
		}

		// Convert back to original type
		originalSeries := df.series[idIdx]
		switch originalSeries.Type().String() {
		case "Bool":
			boolData := make([]bool, newNRows)
			for i, str := range idData {
				if str == "true" {
					boolData[i] = true
				}
			}
			result = result.AddSeries(df.names[idIdx], series.NewSeriesBool(boolData, idNullMask, false, df.ctx)).(BaseDataFrame)
		case "Int":
			intData := make([]int, newNRows)
			for i, str := range idData {
				if !idNullMask[i] && str != "" {
					fmt.Sscanf(str, "%d", &intData[i])
				}
			}
			result = result.AddSeries(df.names[idIdx], series.NewSeriesInt(intData, idNullMask, false, df.ctx)).(BaseDataFrame)
		case "Int64":
			int64Data := make([]int64, newNRows)
			for i, str := range idData {
				if !idNullMask[i] && str != "" {
					fmt.Sscanf(str, "%d", &int64Data[i])
				}
			}
			result = result.AddSeries(df.names[idIdx], series.NewSeriesInt64(int64Data, idNullMask, false, df.ctx)).(BaseDataFrame)
		case "Float64":
			float64Data := make([]float64, newNRows)
			for i, str := range idData {
				if !idNullMask[i] && str != "" {
					fmt.Sscanf(str, "%f", &float64Data[i])
				}
			}
			result = result.AddSeries(df.names[idIdx], series.NewSeriesFloat64(float64Data, idNullMask, false, df.ctx)).(BaseDataFrame)
		default:
			result = result.AddSeries(df.names[idIdx], series.NewSeriesString(idData, idNullMask, false, df.ctx)).(BaseDataFrame)
		}
	}

	// Add new columns from pivot operation
	valuesFromSeries := df.series[valuesFromIdx]
	for _, newColName := range newColNames {
		colData := make([]string, newNRows)
		colNullMask := make([]bool, newNRows)

		// Initialize with fill values or nulls
		for i := range colData {
			colNullMask[i] = true
			if params.ValuesFill != nil {
				colData[i] = fmt.Sprintf("%v", params.ValuesFill)
				colNullMask[i] = false
			}
		}

		// Fill with actual values
		targetName := newColName
		if params.NamesPrefix != "" {
			targetName = strings.TrimPrefix(targetName, params.NamesPrefix)
		}

		for i := 0; i < df.NRows(); i++ {
			namesValue := namesFromSeries.GetAsString(i)
			if namesValue == targetName {
				// Find which ID group this row belongs to
				keyParts := make([]string, len(idIndices))
				for j, idIdx := range idIndices {
					keyParts[j] = df.series[idIdx].GetAsString(i)
				}
				key := strings.Join(keyParts, "|||")

				// Find the position in idOrder
				for j, orderedKey := range idOrder {
					if orderedKey == key {
						colData[j] = valuesFromSeries.GetAsString(i)
						colNullMask[j] = valuesFromSeries.IsNull(i)
						break
					}
				}
			}
		}

		// Add column with appropriate type based on values_from column
		switch valuesFromSeries.Type().String() {
		case "Bool":
			boolData := make([]bool, newNRows)
			for i, str := range colData {
				if str == "true" {
					boolData[i] = true
				}
			}
			result = result.AddSeries(newColName, series.NewSeriesBool(boolData, colNullMask, false, df.ctx)).(BaseDataFrame)
		case "Int":
			intData := make([]int, newNRows)
			for i, str := range colData {
				if !colNullMask[i] && str != "" {
					fmt.Sscanf(str, "%d", &intData[i])
				}
			}
			result = result.AddSeries(newColName, series.NewSeriesInt(intData, colNullMask, false, df.ctx)).(BaseDataFrame)
		case "Int64":
			int64Data := make([]int64, newNRows)
			for i, str := range colData {
				if !colNullMask[i] && str != "" {
					fmt.Sscanf(str, "%d", &int64Data[i])
				}
			}
			result = result.AddSeries(newColName, series.NewSeriesInt64(int64Data, colNullMask, false, df.ctx)).(BaseDataFrame)
		case "Float64":
			float64Data := make([]float64, newNRows)
			for i, str := range colData {
				if !colNullMask[i] && str != "" {
					fmt.Sscanf(str, "%f", &float64Data[i])
				}
			}
			result = result.AddSeries(newColName, series.NewSeriesFloat64(float64Data, colNullMask, false, df.ctx)).(BaseDataFrame)
		default:
			result = result.AddSeries(newColName, series.NewSeriesString(colData, colNullMask, false, df.ctx)).(BaseDataFrame)
		}
	}

	return result
}
