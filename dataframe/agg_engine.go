package dataframe

import (
	"math"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/caerbannogwhite/enchanter/series"
)

// aggMinParallel is the row-count threshold at or above which aggregate
// spawns worker goroutines instead of running the single-chunk core inline.
// Below it, the fixed cost of partitioning rows and merging worker state
// outweighs the benefit of parallelism.
const aggMinParallel = 1 << 16

// aggregateSerial runs the fused single-pass aggregation over df in a single
// chunk covering every row. It shares its accumulation core (accumulateChunk)
// and result-building core (finalizeAggregate) with aggregate's parallel
// workers; see aggregate for the parallel entry point and merge semantics.
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
	nRows := df.NRows()
	valcols, isHol := prepAggValueColumns(df, aggs)

	gt, accs, cols, poisoned := accumulateChunk(keyCols, aggs, valcols, isHol, removeNAs, 0, nRows)

	return finalizeAggregate(df, keyCols, aggs, isHol, removeNAs, gt, accs, cols, poisoned)
}

// aggregate is the public single-pass aggregation entry point. Below
// aggMinParallel rows it runs the single-chunk core inline (identical output
// to aggregateSerial). At or above the threshold, it partitions the rows into
// runtime.GOMAXPROCS(0) contiguous chunks, accumulates each chunk in its own
// goroutine into a chunk-local groupTable and accumulator/collector state
// (accumulateChunk over disjoint row ranges — no shared mutable state between
// workers), then merges every chunk's state into one global groupTable and
// finalizes exactly as the serial path.
//
// Merge correctness relies on two points:
//
//   - Chunk-local group ids are NOT comparable across chunks: groupTable
//     assigns dense ids in first-appearance order, which differs per chunk, so
//     the same logical group can get different ids in different chunks. Each
//     chunk-local group is therefore re-keyed into the global groupTable by
//     replaying its representative row through global.idOf, which re-reads
//     that row's actual key-column values, so equal groups collapse to the
//     same global id regardless of which chunk (or which local id) they came
//     from.
//
//   - Under removeNAs == false, poison flags are OR-merged per (aggregator,
//     group): a group poisoned in any one chunk is poisoned in the merged
//     result. This matches the serial engine's per-row poisoning exactly,
//     since a poisoned group there stays poisoned once any row in it is null.
func aggregate(df BaseDataFrame, keyCols []series.Series, aggs []aggregator, removeNAs bool) DataFrame {
	nRows := df.NRows()
	valcols, isHol := prepAggValueColumns(df, aggs)

	if nRows < aggMinParallel {
		gt, accs, cols, poisoned := accumulateChunk(keyCols, aggs, valcols, isHol, removeNAs, 0, nRows)
		return finalizeAggregate(df, keyCols, aggs, isHol, removeNAs, gt, accs, cols, poisoned)
	}

	nWorkers := runtime.GOMAXPROCS(0)
	if nWorkers < 1 {
		nWorkers = 1
	}
	if nWorkers > nRows {
		nWorkers = nRows
	}
	chunkSize := (nRows + nWorkers - 1) / nWorkers

	type chunkState struct {
		gt       *groupTable
		accs     [][]accumulator
		cols     [][]collector
		poisoned [][]bool
	}

	// Each goroutine writes only its own index w; disjoint slice elements are
	// safe to write concurrently. wg.Wait() below establishes the
	// happens-before edge before the main goroutine reads any of them.
	chunks := make([]chunkState, nWorkers)
	var wg sync.WaitGroup
	for w := 0; w < nWorkers; w++ {
		lo := w * chunkSize
		hi := lo + chunkSize
		if hi > nRows {
			hi = nRows
		}
		if lo >= hi {
			continue
		}
		wg.Add(1)
		go func(w, lo, hi int) {
			defer wg.Done()
			gt, accs, cols, poisoned := accumulateChunk(keyCols, aggs, valcols, isHol, removeNAs, lo, hi)
			chunks[w] = chunkState{gt, accs, cols, poisoned}
		}(w, lo, hi)
	}
	wg.Wait()

	// Merge every chunk's local state into one global groupTable, re-keying
	// each chunk-local group by its representative row (see doc comment
	// above) and OR-merging poison flags. Chunks are merged in worker order
	// for determinism; the merge itself is order-independent — accumulator
	// and collector merge are commutative, and a group's representative row
	// always carries the same key values no matter which chunk supplied it.
	global := newGroupTable(keyCols)
	gAccs := make([][]accumulator, len(aggs))
	gCols := make([][]collector, len(aggs))
	gPoisoned := make([][]bool, len(aggs))

	nGroups := 0
	for _, chunk := range chunks {
		if chunk.gt == nil { // empty chunk (lo >= hi above)
			continue
		}
		reps := chunk.gt.representativeRows()
		for localGid, row := range reps {
			gid := global.idOf(row)

			if gid == nGroups {
				for j, agg := range aggs {
					if isHol[j] {
						gCols[j] = append(gCols[j], newQuantileCollector(agg.p, agg.interp))
					} else {
						gAccs[j] = append(gAccs[j], newReducibleAcc(agg.type_, agg.ddof))
					}
					gPoisoned[j] = append(gPoisoned[j], false)
				}
				nGroups++
			}

			for j := range aggs {
				if isHol[j] {
					gCols[j][gid].merge(chunk.cols[j][localGid])
				} else {
					gAccs[j][gid].merge(chunk.accs[j][localGid])
				}
				if chunk.poisoned[j][localGid] {
					gPoisoned[j][gid] = true
				}
			}
		}
	}

	return finalizeAggregate(df, keyCols, aggs, isHol, removeNAs, global, gAccs, gCols, gPoisoned)
}

// prepAggValueColumns coerces each aggregator's value column to []float64
// once (NaN == null); Count needs no value column. The returned slices are
// read-only from this point on and safe to share across accumulateChunk
// calls running concurrently over disjoint row ranges.
func prepAggValueColumns(df BaseDataFrame, aggs []aggregator) (valcols [][]float64, isHol []bool) {
	valcols = make([][]float64, len(aggs))
	isHol = make([]bool, len(aggs))
	for j, agg := range aggs {
		isHol[j] = agg.type_.isHolistic()
		if agg.type_ == AGGREGATE_COUNT {
			continue
		}
		valcols[j] = __gdl_stats_preprocess(df.C(agg.name))
	}
	return valcols, isHol
}

// accumulateChunk is the single-chunk accumulation core shared by
// aggregateSerial and aggregate's parallel workers. It builds a fresh
// groupTable over rows [lo, hi) of keyCols, with one accumulator (reducible)
// or collector (holistic) per aggregator per group encountered in that row
// range, and returns the per-(aggregator, group) poison flags for
// removeNAs == false (see aggregateSerial's doc comment for poison
// semantics). valcols/isHol come from prepAggValueColumns and are read-only.
//
// The returned groupTable, accumulators, collectors and poison flags are
// local to [lo, hi): group ids are dense in first-appearance order within
// this chunk only and are not meaningful outside it (see aggregate's doc
// comment on merging).
func accumulateChunk(keyCols []series.Series, aggs []aggregator, valcols [][]float64, isHol []bool, removeNAs bool, lo, hi int) (*groupTable, [][]accumulator, [][]collector, [][]bool) {
	gt := newGroupTable(keyCols)

	// Per-aggregator, per-gid state; grown lazily as new gids appear.
	accs := make([][]accumulator, len(aggs)) // reducible (incl. Count)
	cols := make([][]collector, len(aggs))   // holistic
	poisoned := make([][]bool, len(aggs))    // reducible poison (removeNAs == false)

	nGroups := 0
	for row := lo; row < hi; row++ {
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

	return gt, accs, cols, poisoned
}

// finalizeAggregate builds the sorted result DataFrame from a fully-populated
// (or fully-merged) group table and its per-aggregator per-group state: the
// key columns (one row per group, taken from each group's representative
// row) followed by one aggregate column per aggregator, sorted by key
// ascending, nulls last, lexicographic across keyCols. Count is emitted as
// Int64s; every other aggregate as Float64s with a null mask from the
// per-row isNull flags (poisoned groups under removeNAs == false surface as
// non-null NaN, matching aggregateSerial).
func finalizeAggregate(df BaseDataFrame, keyCols []series.Series, aggs []aggregator, isHol []bool, removeNAs bool, gt *groupTable, accs [][]accumulator, cols [][]collector, poisoned [][]bool) DataFrame {
	ctx := df.GetContext()
	nGroups := gt.numGroups()
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
// equal. The type switch covers the same key types the group-key coder supports.
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
// order. The type switch covers the same key types as the former groupHelper;
// the representative row's null flag is carried into the emitted column.
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
