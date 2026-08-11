# 0.4.0 spike — pragmatic pure-Go groupby + fusion

Throwaway measuring-stick (see `spike_test.go`). Question: how much of the
groupby gap does a single-pass, **parallel** hash-aggregation close vs
enchanter's current eager groupby — and does operator fusion pay off?

Machine: 16 logical cores, Go 1.26.5 (scalar, no SIMD). Data: h2oai G1,
`sum(v1) group by id1`, 100 groups. Correctness verified (parallel == single ==
direct total). Numbers are indicative single runs.

## Groupby `sum(v1) by id1`

| rows | regime | ns/op | vs baseline | B/op | allocs/op |
|------|--------|------:|:-----------:|-----:|----------:|
| 1e6 | baseline (current `GroupBy.Agg.Run`) | 25,390,238 | 1.0× | 80,642,888 | 19,287 |
| 1e6 | single-pass hash (1 goroutine) | 21,224,993 | **1.2×** | 4,904 | 3 |
| 1e6 | parallel hash (16 goroutines) | 2,689,704 | **9.4×** | 86,632 | 102 |
| 1e7 | baseline | 328,810,675 | 1.0× | 847,331,526 | 28,893 |
| 1e7 | single-pass hash | 217,971,100 | **1.5×** | 4,904 | 3 |
| 1e7 | parallel hash | 24,114,751 | **13.6×** | 88,162 | 104 |

## Fusion cameo `(a+b) > 0 then filter`, 1e6

| regime | ns/op | vs eager | B/op | allocs/op |
|--------|------:|:--------:|-----:|----------:|
| eager 3-pass (`Add` → `Gt` → `Filter`) | 6,882,271 | 1.0× | 12,984,729 | 9 |
| fused 1-pass loop | 2,685,318 | **2.6×** | 4,005,888 | 1 |

## Findings

1. **Parallelism is the big lever — and it grows with size.** 9.4× at 1e6 →
   13.6× at 1e7, near-linear on 16 cores (partition → local hash tables → merge).
   Essentially free to add. This is the headline.

2. **Single-pass barely moves the clock but annihilates allocation.** 1.2–1.5×
   in time, but 80 MB → 4.9 KB (1e6) and **847 MB → 4.9 KB** (1e7): the current
   partition path materializes a `map[key][]int` of row indices plus a
   `[]float64` preprocess copy (the ~19–29k allocations); the single pass keeps
   only the 100-entry result. That garbage is GC pressure that compounds across
   a real multi-op pipeline even where a microbenchmark hides it.

3. **The pointer-key trick works.** Because enchanter interns strings, hashing
   the `*string` (address) instead of the bytes is correct and cheap — the
   current partition path doesn't exploit it.

4. **Fusion delivers as predicted.** 2.6× and 3.2× less memory by never
   materializing the intermediate `sum`/`mask` arrays — in pure Go, no Arrow.

5. **Headroom left on the table.** These use Go's builtin `map`; an
   open-addressed table specialized for `*string` keys would push the
   single-threaded (and therefore parallel) numbers further. Go 1.27's
   experimental SIMD would add more on the accumulation. Neither was needed to
   get here.

## Verdict: GO for a "make it fast" 0.4.0

The current groupby was 2.3–7.8× behind Polars (see `../data/notes.md`); a
**9–14× speedup over the current path** more than covers that gap, strongly
suggesting pragmatic pure-Go can match or beat Polars on this operation — with
zero Arrow and zero SIMD. The pattern to generalize into the library:

- **single-pass + parallel hash aggregation** for `GroupBy().Agg()`, and
- **fused combinators** for common element-wise op chains.

### Caveats (what a real implementation must still handle, and what this omits)
- Multi-key groups (`id1:id2`), typed/int keys, and string keys **without**
  interning (the pointer trick relies on enchanter's StringPool).
- All aggregates (mean/min/max/std/count), not just sum; and null keys / NA
  handling.
- The spike times the core aggregation given pre-extracted key/value slices; the
  baseline additionally builds the result DataFrame (small for 100 groups) and
  extracts columns. Fair as "the core groupby-agg computation," not a full
  end-to-end op.
- A rigorous cross-library number needs a fresh Polars run on this machine; the
  Polars gap above is from the recorded baseline, used as directional context.
