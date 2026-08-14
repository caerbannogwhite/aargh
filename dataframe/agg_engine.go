package dataframe

import (
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/caerbannogwhite/enchanter/series"
)

// aggregateSerial runs the fused single-pass aggregation over df in one worker.
//
// It builds one group per distinct composite key of keyCols (via groupTable),
// holding one accumulator (reducible) or collector (holistic) per aggregator
// per group. A single pass over the rows updates the per-group state; nulls are
// coerced to NaN by __gdl_stats_preprocess and detected with math.IsNaN.
//
//   - removeNAs == true  (the new default): null values are skipped for every
//     aggregate; nothing is poisoned.
//   - removeNAs == false: a null value poisons that (group, aggregator) so its
//     result becomes NaN — for every aggregate alike, reducible (Mean, ...) and
//     holistic (Median, Quantile). Count is the sole exception: it counts rows
//     and is independent of value nulls. Collectors still never receive NaN; the
//     poison flag overrides the collector's result at finalize.
//
// The result is the key columns (one value per group, taken from each group's
// representative row) followed by one aggregate column per aggregator, with the
// rows sorted by key ascending, nulls last, lexicographically across keyCols.
// Count is emitted as Int64s; every other aggregate is emitted as Float64s with
// a null mask from the finalize isNull flags.
//
// len(keyCols) == 0 is handled as a single group with no key columns.
func aggregateSerial(df BaseDataFrame, keyCols []series.Series, aggs []aggregator, removeNAs bool) DataFrame {
	ctx := df.GetContext()
	nRows := df.NRows()

	// Coerce each aggregator's value column to []float64 once (NaN == null).
	// Count needs no value column.
	valcols := make([][]float64, len(aggs))
	isHol := make([]bool, len(aggs))
	for j, agg := range aggs {
		isHol[j] = agg.type_.isHolistic()
		if agg.type_ == AGGREGATE_COUNT {
			continue
		}
		valcols[j] = __gdl_stats_preprocess(df.C(agg.name))
	}

	gt := newGroupTable(keyCols)

	// Per-aggregator, per-gid state; grown lazily as new gids appear.
	accs := make([][]accumulator, len(aggs)) // reducible (incl. Count)
	cols := make([][]collector, len(aggs))   // holistic
	poisoned := make([][]bool, len(aggs))    // reducible poison (removeNAs == false)

	nGroups := 0
	for row := 0; row < nRows; row++ {
		gid := gt.idOf(row)

		// New group: append fresh state for every aggregator.
		if gid == nGroups {
			for j, agg := range aggs {
				if isHol[j] {
					cols[j] = append(cols[j], newQuantileCollector(agg.p, agg.interp))
				} else {
					accs[j] = append(accs[j], newReducibleAcc(agg.type_, agg.ddof))
				}
				poisoned[j] = append(poisoned[j], false)
			}
			nGroups++
		}

		for j, agg := range aggs {
			if agg.type_ == AGGREGATE_COUNT {
				accs[j][gid].update(0) // value ignored; Count counts rows
				continue
			}
			vj := valcols[j][row]
			if math.IsNaN(vj) { // null
				if !removeNAs {
					// Poisons every aggregate for this group under RemoveNAs(false),
					// reducible and holistic alike (Count never reaches here).
					poisoned[j][gid] = true
				}
				continue
			}
			if isHol[j] {
				cols[j][gid].collect(vj)
			} else {
				accs[j][gid].update(vj)
			}
		}
	}

	reps := gt.representativeRows()
	order := sortGroupOrder(keyCols, reps)

	// Build the result: key columns first, then one aggregate column per agg.
	result := NewBaseDataFrame(ctx)
	result = appendKeyColumns(result, df, keyCols, reps, order)

	for j, agg := range aggs {
		if agg.type_ == AGGREGATE_COUNT {
			data := make([]int64, nGroups)
			for i, gid := range order {
				v, _ := accs[j][gid].result()
				data[i] = int64(v)
			}
			result = result.AddSeries(agg.newName, series.NewSeriesInt64(data, nil, false, ctx))
			continue
		}

		data := make([]float64, nGroups)
		var nulls []bool
		for i, gid := range order {
			var v float64
			var isNull bool
			switch {
			case !removeNAs && poisoned[j][gid]:
				// Poison overrides both reducible and holistic results.
				v, isNull = math.NaN(), false
			case isHol[j]:
				v, isNull = cols[j][gid].result()
			default:
				v, isNull = accs[j][gid].result()
			}
			data[i] = v
			if isNull {
				if nulls == nil {
					nulls = make([]bool, nGroups)
				}
				nulls[i] = true
			}
		}
		result = result.AddSeries(agg.newName, series.NewSeriesFloat64(data, nulls, false, ctx))
	}

	return result
}

// sortGroupOrder returns the group ids permuted into sorted key order:
// ascending, nulls last, lexicographic across the key columns (compared at each
// group's representative row). With no key columns the identity order (group id
// order, i.e. first-appearance) is returned.
func sortGroupOrder(keyCols []series.Series, reps []int) []int {
	order := make([]int, len(reps))
	for i := range order {
		order[i] = i
	}
	if len(keyCols) == 0 {
		return order
	}
	sort.SliceStable(order, func(i, j int) bool {
		ri, rj := reps[order[i]], reps[order[j]]
		for _, col := range keyCols {
			if c := compareKeyCells(col, ri, rj); c != 0 {
				return c < 0
			}
		}
		return false
	})
	return order
}

// compareKeyCells compares the cells of col at rows ra and rb, returning -1, 0
// or +1. A null cell sorts after any non-null cell (nulls last); two nulls are
// equal. The type switch mirrors groupHelper's supported key types.
func compareKeyCells(col series.Series, ra, rb int) int {
	na, nb := col.IsNull(ra), col.IsNull(rb)
	if na || nb {
		switch {
		case na && nb:
			return 0
		case na:
			return 1 // a is null -> sorts after b
		default:
			return -1 // b is null -> a sorts before
		}
	}

	switch c := col.(type) {
	case series.Bools:
		return compareBool(c.Data_[ra], c.Data_[rb])
	case series.Ints:
		return cmpOrdered(c.Data_[ra], c.Data_[rb])
	case series.Int64s:
		return cmpOrdered(c.Data_[ra], c.Data_[rb])
	case series.Float64s:
		return cmpOrdered(c.Data_[ra], c.Data_[rb])
	case series.Strings:
		return cmpOrdered(*c.Data_[ra], *c.Data_[rb])
	case series.Times:
		return cmpOrdered(c.Data_[ra].UnixNano(), c.Data_[rb].UnixNano())
	case series.Durations:
		return cmpOrdered(int64(c.Data_[ra]), int64(c.Data_[rb]))
	default:
		return cmpOrdered(col.GetAsString(ra), col.GetAsString(rb))
	}
}

// cmpOrdered is a three-way comparison for the ordered key cell types.
func cmpOrdered[T int | int64 | float64 | string](a, b T) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// compareBool orders false before true.
func compareBool(a, b bool) int {
	switch {
	case a == b:
		return 0
	case !a:
		return -1
	default:
		return 1
	}
}

// appendKeyColumns emits one key column per keyCol, one value per group taken
// from the group's representative row, with the groups laid out in the sorted
// order. The type switch mirrors groupHelper (base.go); the representative
// row's null flag is carried into the emitted column.
func appendKeyColumns(result DataFrame, df BaseDataFrame, keyCols []series.Series, reps, order []int) DataFrame {
	ctx := df.GetContext()
	n := len(order)

	for k, col := range keyCols {
		name := keyColumnName(df, col, k)

		var keyNulls []bool
		if col.IsNullable() {
			keyNulls = make([]bool, n)
			for i, gid := range order {
				keyNulls[i] = col.IsNull(reps[gid])
			}
		}

		switch c := col.(type) {
		case series.Bools:
			vals := make([]bool, n)
			for i, gid := range order {
				vals[i] = c.Data_[reps[gid]]
			}
			result = result.AddSeries(name, series.NewSeriesBool(vals, keyNulls, false, ctx))

		case series.Ints:
			vals := make([]int, n)
			for i, gid := range order {
				vals[i] = c.Data_[reps[gid]]
			}
			result = result.AddSeries(name, series.NewSeriesInt(vals, keyNulls, false, ctx))

		case series.Int64s:
			vals := make([]int64, n)
			for i, gid := range order {
				vals[i] = c.Data_[reps[gid]]
			}
			result = result.AddSeries(name, series.NewSeriesInt64(vals, keyNulls, false, ctx))

		case series.Float64s:
			vals := make([]float64, n)
			for i, gid := range order {
				vals[i] = c.Data_[reps[gid]]
			}
			result = result.AddSeries(name, series.NewSeriesFloat64(vals, keyNulls, false, ctx))

		case series.Strings:
			vals := make([]*string, n)
			for i, gid := range order {
				vals[i] = c.Data_[reps[gid]]
			}
			result = result.AddSeries(name, series.NewSeriesStringFromPtrs(vals, keyNulls, false, ctx))

		case series.Times:
			vals := make([]time.Time, n)
			for i, gid := range order {
				vals[i] = c.Data_[reps[gid]]
			}
			result = result.AddSeries(name, series.NewSeriesTime(vals, keyNulls, false, ctx))

		case series.Durations:
			vals := make([]time.Duration, n)
			for i, gid := range order {
				vals[i] = c.Data_[reps[gid]]
			}
			result = result.AddSeries(name, series.NewSeriesDuration(vals, keyNulls, false, ctx))

		default:
			// Unsupported key type: keep names and rows aligned.
			result = result.AddSeries(name, series.NewSeriesNA(n, ctx))
		}
	}

	return result
}

// keyColumnName recovers the dataframe name of a key column. The engine is
// called with keyCols drawn from df.series (see buildGroupKeyCols), so the
// column is matched back to its name by the identity of its backing data array.
// The k-based fallback is defensive: it is unreachable when keyCols come from df.
func keyColumnName(df BaseDataFrame, col series.Series, k int) string {
	if p := seriesBackingPtr(col); p != nil {
		for i, s := range df.series {
			if seriesBackingPtr(s) == p {
				return df.names[i]
			}
		}
	}
	return "key_" + strconv.Itoa(k)
}

// seriesBackingPtr returns a comparable identity for a series' backing data
// array (the address of its first element), or nil when the series is empty or
// of an unrecognized type. Two series that share a backing array (e.g. the
// stored column and the value returned by df.C) yield equal identities.
func seriesBackingPtr(s series.Series) any {
	switch c := s.(type) {
	case series.Bools:
		if len(c.Data_) == 0 {
			return nil
		}
		return &c.Data_[0]
	case series.Ints:
		if len(c.Data_) == 0 {
			return nil
		}
		return &c.Data_[0]
	case series.Int64s:
		if len(c.Data_) == 0 {
			return nil
		}
		return &c.Data_[0]
	case series.Float64s:
		if len(c.Data_) == 0 {
			return nil
		}
		return &c.Data_[0]
	case series.Strings:
		if len(c.Data_) == 0 {
			return nil
		}
		return &c.Data_[0]
	case series.Times:
		if len(c.Data_) == 0 {
			return nil
		}
		return &c.Data_[0]
	case series.Durations:
		if len(c.Data_) == 0 {
			return nil
		}
		return &c.Data_[0]
	default:
		return nil
	}
}
