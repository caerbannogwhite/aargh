# enchanter — Arrow-native storage migration: measurements + go/no-go

Status: decision doc (feeds the 0.4.0 go/no-go)
Date: 2026-08-08
Author: caerbannogwhite (with Claude)
Companion to: [`2026-08-08-arrow-performance-design.md`](2026-08-08-arrow-performance-design.md)

## TL;DR — the verdict is **NO-GO for now**

Arrow-native storage does **not** clear the bar the 0.3.0 design set for itself:
arrow-native ≥ ~1.5× on the 2-op chain at ≥ 1e6 rows *without badly regressing
single ops*. The measured chain speed-up is **1.45× at 1e6** and it **decays to
1.21× at 1e7** — the win shrinks exactly as the row count grows toward the sizes
that matter, and single-op `Add` is at parity (1.02–1.16×) or worse (0.61× at
1e4). The one unambiguous win, SIMD `Sum` (3.4× at 1e6), also collapses to 1.08×
at 1e7 once the data no longer fits in cache. The round-trip regime — what you'd
actually get by wiring Arrow compute into today's Go-slice storage — is a **1.9×
to 4.4× regression**.

The Arrow lever is real but narrow: it pays only for long, all-Arrow op chains at
mid sizes, and only if data never leaves Arrow. enchanter's DataFrame model is
mutation-friendly (`Set`, `Sort`, `Append`) and hands `[]T` back to users, so it
cannot keep data Arrow end-to-end without either forbidding those operations or
paying the round-trip on every boundary. **Recommendation: keep Arrow as an
interop layer (Parquet / IPC / zero-copy hand-off), do not migrate series
storage in 0.4.0.** The specific conditions that would flip this are listed at
the end.

---

## 1. The measurements

Full run, all sizes 1e4–1e7 (not `-short`), captured on the development machine:

```
goos: windows
goarch: amd64
pkg: github.com/caerbannogwhite/enchanter/series
cpu: AMD Ryzen 9 6900HX with Radeon Graphics
BenchmarkAdd/goslice/1e4-16         	   70095	     19164 ns/op	   82080 B/op	       3 allocs/op
BenchmarkAdd/roundtrip/1e4-16       	   15063	     84758 ns/op	  368388 B/op	      47 allocs/op
BenchmarkAdd/arrownative/1e4-16     	   39742	     31231 ns/op	   84324 B/op	      33 allocs/op
BenchmarkAdd/goslice/1e5-16         	    5419	    235116 ns/op	  802976 B/op	       3 allocs/op
BenchmarkAdd/roundtrip/1e5-16       	    2470	    503797 ns/op	 2957067 B/op	      48 allocs/op
BenchmarkAdd/arrownative/1e5-16     	   10000	    129422 ns/op	  805390 B/op	      33 allocs/op
BenchmarkAdd/goslice/1e6-16         	     578	   2045849 ns/op	 8003744 B/op	       3 allocs/op
BenchmarkAdd/roundtrip/1e6-16       	     301	   3820886 ns/op	25079911 B/op	      48 allocs/op
BenchmarkAdd/arrownative/1e6-16     	     680	   1768171 ns/op	 8006371 B/op	      33 allocs/op
BenchmarkAdd/goslice/1e7-16         	      60	  19387305 ns/op	80003232 B/op	       3 allocs/op
BenchmarkAdd/roundtrip/1e7-16       	      20	  59250580 ns/op	352669940 B/op	      48 allocs/op
BenchmarkAdd/arrownative/1e7-16     	      62	  19065785 ns/op	80005769 B/op	      33 allocs/op
BenchmarkChain/goslice/1e4-16       	   16802	     76395 ns/op	  133528 B/op	       9 allocs/op
BenchmarkChain/arrownative/1e4-16   	   16272	     75585 ns/op	  131472 B/op	     101 allocs/op
BenchmarkChain/goslice/1e5-16       	    1543	    686801 ns/op	 1311130 B/op	       9 allocs/op
BenchmarkChain/arrownative/1e5-16   	    2421	    425125 ns/op	 1225443 B/op	     102 allocs/op
BenchmarkChain/goslice/1e6-16       	     248	   4891004 ns/op	13009304 B/op	       9 allocs/op
BenchmarkChain/arrownative/1e6-16   	     348	   3365626 ns/op	12140392 B/op	     102 allocs/op
BenchmarkChain/goslice/1e7-16       	      25	  48386644 ns/op	129933720 B/op	       9 allocs/op
BenchmarkChain/arrownative/1e7-16   	      34	  40068288 ns/op	121192150 B/op	     102 allocs/op
BenchmarkSum/goslice/1e4-16         	  337776	      3301 ns/op	       0 B/op	       0 allocs/op
BenchmarkSum/arrownative/1e4-16     	 1893306	       646.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkSum/goslice/1e5-16         	   40359	     26891 ns/op	       0 B/op	       0 allocs/op
BenchmarkSum/arrownative/1e5-16     	  162277	      7690 ns/op	       0 B/op	       0 allocs/op
BenchmarkSum/goslice/1e6-16         	    4333	    272996 ns/op	       0 B/op	       0 allocs/op
BenchmarkSum/arrownative/1e6-16     	   15176	     80890 ns/op	       0 B/op	       0 allocs/op
BenchmarkSum/goslice/1e7-16         	     460	   2827849 ns/op	       0 B/op	       0 allocs/op
BenchmarkSum/arrownative/1e7-16     	     453	   2617399 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/caerbannogwhite/enchanter/series	45.581s
```

Raw output: `.superpowers/sdd/2026-08-08-0.3.0-arrow-performance/bench-0.3.0.txt`.
Reproduce with:

```
go test ./series/ -run '^$' -bench 'BenchmarkAdd|BenchmarkChain|BenchmarkSum' -benchmem -benchtime=1s
```

Three regimes (defined in the 0.3.0 design doc): **go-slice** is today's `[]T` +
null-mask path; **roundtrip** builds an Arrow array from the Go slice, runs the
kernel, materializes back to Go (the cost of naïvely bolting compute onto today's
storage); **arrow-native** takes data already in `arrow.Array`, computes, and
leaves the result Arrow (what storage-native would buy). The prototype
(`series/arrow_native_float64.go`) runs through real methods so the numbers
include Go interface-dispatch and allocation overhead, not raw kernel calls in
isolation.

### Speed-up summary (go-slice time ÷ arrow-native time; >1 = arrow-native faster)

| Op    | 1e4  | 1e5  | 1e6      | 1e7      |
| ----- | ---- | ---- | -------- | -------- |
| Add   | 0.61 | 1.82 | 1.16     | 1.02     |
| Chain | 1.01 | 1.62 | **1.45** | **1.21** |
| Sum   | 5.11 | 3.50 | 3.37     | 1.08     |

Add round-trip vs go-slice (regression factor, higher = worse): 4.42× / 2.14× /
1.87× / 3.06× at 1e4 / 1e5 / 1e6 / 1e7.

### Read per operation

**Add (single element-wise op).** arrow-native is *slower* than the Go loop at
1e4 (0.61×) — the Arrow builder allocation, the `compute.Add` datum wrapping, and
interface dispatch cost more than the whole 10k-element add. It pulls ahead at
1e5 (1.82×) where the vectorized inner loop dominates, then the margin **erodes
back toward parity** as size grows: 1.16× at 1e6, 1.02× at 1e7. By 1e7 the
operation is memory-bandwidth bound (80 MB per operand) and Arrow's SIMD inner
loop no longer matters — both paths are waiting on RAM. The **round-trip regime
is a disaster at every size** (1.9–4.4× slower, 3–4.5× the bytes, 47–48 allocs
vs 3): building the Arrow array and materializing the result back to a Go slice
is pure overhead. This is the single most important line in the table, because
round-trip — not arrow-native — is what you actually get if you wire compute into
storage that still hands users `[]float64`.

**Chain (the decision op: `(a+b) > 0` then filter).** This is where arrow-native
is supposed to shine: one build amortized across three kernels (add, greater,
filter), no intermediate materialization. It does win — but modestly and
non-monotonically: 1.01× (1e4, break-even), **1.62× (1e5, the peak)**, **1.45×
(1e6)**, **1.21× (1e7)**. The win **peaks at 1e5 and decays through the sizes we
actually care about.** At 1e6 it grazes the ~1.5× guide; at 1e7 it is clearly
under it. Allocation-wise arrow-native is *worse* on count — **~101–102 allocs vs
9** — from the per-kernel Arrow buffers and datum objects, though total bytes are
slightly lower (fewer full-width Go slices). The alloc-count blow-up is a
warning: every kernel in a real pipeline mints Arrow `Data`/`Buffer`/`Datum`
objects, and a longer chain than this two-comparison toy would amplify that, not
amortize it.

**Sum (SIMD reduction).** The clearest arrow win and the clearest cautionary
tale. `arrmath.Float64.Sum` (hand-written AVX) beats the scalar Go loop **5.1× at
1e4, 3.5× at 1e5, 3.4× at 1e6** — all zero-alloc on both sides. Then at 1e7 it
**collapses to 1.08×**: at 80 MB the sum is memory-bandwidth bound and SIMD width
stops helping. So even enchanter's single best Arrow primitive delivers its
advantage only while data fits in cache, and evaporates at the largest size. And
critically (see §2) arrow-go provides SIMD `Sum` for exactly three types and
**no grouped aggregation at all** — the operation that actually dominates the
cross-library gap.

### What the table means, in one sentence

The Arrow storage advantage is real, small, size-fragile (peaks at 1e5, decays by
1e7), confined to multi-op all-Arrow chains, and comes with a ~10× higher
allocation count — while the regime you'd actually ship on today's storage
(round-trip) is a 2–4× regression.

---

## 2. Cross-library framing (motivation, not a comparand)

The reason storage performance was worth measuring at all is the baseline in
[`benchmarking/data/notes.md`](../../../benchmarking/data/notes.md): on the h2oai
db-benchmark `groupby` task, enchanter's ancestor (Gandalff) ran **2.3–7.8×
slower than Polars**, worst at the largest size (7.8× at 1e7 rows in the initial
baseline; ~2.1–2.5× at 1e7 after hand-tuning the group/aggregate loops and going
multi-goroutine). Polars is Arrow-native and Rust; the working hypothesis behind
0.3.0 was that adopting Arrow storage would close that gap.

The measurements **do not support that hypothesis for the operations Arrow can
actually accelerate in Go**, and — more decisively — Arrow-native storage cannot
touch the operation where the gap is largest:

- **The gap is a `groupby`-aggregation gap.** All five benchmark questions (Q1–Q5)
  are `sum/mean … by …`. Polars' advantage there is its parallel hash-aggregation
  engine, not element-wise arithmetic.
- **arrow-go v18.5.2 has no grouped-aggregation kernel** — no hash-aggregate, and
  no scalar `sum`/`mean`/`min`/`max`/`stddev` reductions either. The only
  accelerated reduction in the whole library is `math.{Float64,Int64,Uint64}.Sum`
  (SIMD, `Sum` only, three types). Switching series storage to Arrow would give
  enchanter's `GroupBy`/aggregation code **nothing new to call**: it would still
  hand-roll the group hash map and the per-group accumulation loop exactly as it
  does today, just reading values out of an `arrow.Array` instead of a `[]T`
  (arguably *slower*, via `arr.Value(i)` bounds checks vs raw slice indexing).
- Even the parts Arrow *can* accelerate don't scale into the regime that dominates
  the Polars gap: the gap is worst at 1e7, and that is precisely where every
  arrow-native win in the table collapses toward parity (Add 1.02×, Chain 1.21×,
  Sum 1.08×) because the workload becomes memory-bandwidth bound.

So closing the Polars gap is a **grouped-aggregation and parallelism** problem
(better hash tables, goroutine fan-out, cache-aware layout — the avenues the
`notes.md` log was already pursuing), **not** a storage-format problem that
Arrow-in-Go can solve with the kernels available. Matching Polars' actual engine
would require either re-implementing hash-aggregation from scratch (independent of
storage format) or a cgo binding to Arrow C++ (which breaks enchanter's pure-Go
promise and is permanently out of scope). Arrow-native storage is neither
necessary nor sufficient for that work. This is the motivation reframed by the
data, not a head-to-head claim: the microbenchmark is the evidence; the Polars gap
is the "why", and the data says storage-format is the wrong lever for it.

---

## 3. Verdict against the criterion

**Criterion (from the 0.3.0 design):** GO if arrow-native ≥ ~1.5× on the 2-op
chain at ≥ 1e6 rows, without badly regressing single ops. The threshold is an
explicit "reasoned guide, not a tripwire."

**Measured against it:**

| Sub-criterion                                   | Result                                        | Pass?          |
| ----------------------------------------------- | --------------------------------------------- | -------------- |
| Chain ≥ ~1.5× at 1e6                            | 1.45×                                          | Marginal miss  |
| Chain holds up at 1e7                           | 1.21× (decaying, not growing)                 | Miss           |
| Single ops not badly regressed (arrow-native)   | Add 1.02–1.16× at ≥1e6; 0.61× at 1e4          | Parity, not a win |
| Single ops on realistic (round-trip) storage    | 1.9–4.4× **regression**                       | Fail           |

**Decision: NO-GO for now.** Reasoning, weighing the table rather than the single
1e6 number:

1. **The chain win misses and, worse, decays.** 1.45× at 1e6 is a marginal miss
   of the ~1.5× guide, and if the guide were a tripwire that would end it. But the
   more important fact is the *direction*: the advantage peaks at 1e5 (1.62×) and
   falls to 1.21× at 1e7. A migration justified by "wins get bigger with data"
   is contradicted — this one gets smaller. The sizes where enchanter most needs
   to be fast (1e6–1e7, per the Polars gap) are the sizes where Arrow helps least.
2. **Single ops don't win, and the shippable regime loses badly.** arrow-native
   single-op Add is parity at best (1.02× at 1e7) and a loss at 1e4 (0.61×). And
   the regime you would actually ship without a *total* storage rewrite —
   round-trip — regresses 1.9–4.4×. A migration only reaches the arrow-native
   column if *every* op in *every* pipeline stays Arrow with zero materialization;
   the moment a user calls `.Data()`, sorts, or sets a value, you pay round-trip.
3. **The cost is disproportionate to the benefit.** Realizing even the modest
   chain win requires (per §4) reworking all ~49K generated lines, redesigning
   every mutating method around Arrow immutability, and inverting the null-mask
   convention — a multi-release, high-risk migration behind the `Series`
   interface — to buy a 1.2–1.6× speed-up on a narrow op shape that isn't the
   bottleneck. The bottleneck (grouped aggregation vs Polars) gets **no** benefit
   because arrow-go has no kernel for it.
4. **The one big win is size-fragile and type-narrow.** SIMD `Sum` (3.4× at 1e6)
   is genuinely nice, collapses at 1e7, and exists for three numeric types only —
   and it can already be captured today *without* a storage migration by having
   the aggregators call `arrmath.*.Sum` on a temporary Arrow array where profitable
   (a self-contained optimization, not a rewrite).

Per the design doc's own guidance ("record it as NO-GO-for-now with the specific
conditions under which it would flip, rather than forcing a decision"), this is
the honest reading. **Keep Arrow interop-only.**

---

## 4. What dominates the cost (the NO-GO analysis)

The three regimes localize where the time goes.

**Build + materialize is the killer (round-trip regime).** Round-trip Add costs
1.9–4.4× the go-slice time and allocates 47–48 objects and 3–4.5× the bytes,
against go-slice's 3 allocs. That entire delta is `array.NewFloat64Builder` +
`Reserve` + `AppendValues` on the way in, and `Float64Values()` copy on the way
out. Any storage design that still exposes `[]T` to users (enchanter's does —
`Data() any`, `DataAsNullable`, direct `Data_` field access across `dataframe/`
and `io/`) pays this at every Go-boundary crossing. The arrow-native column only
exists because the benchmark deliberately keeps data Arrow across the whole chain;
a real DataFrame cannot make that guarantee.

**Allocation count, not just time.** arrow-native Chain runs ~101–102 allocs/op
vs go-slice's 9. Every kernel invocation allocates an `arrow.Data`, its `Buffer`s,
and `compute.Datum` wrappers. Two ops already yield an 11× alloc-count blow-up; a
realistic 6–10 op pipeline would multiply GC pressure further. Arrow's expected
*advantage* was fewer allocations; measured, it's the opposite on count.

**Interface dispatch and memory bandwidth cap the ceiling.** At 1e4 the fixed
per-op costs (dispatch, datum construction, one builder alloc) swamp the actual
arithmetic — hence Add's 0.61×. At 1e7 the operation is bandwidth-bound, so SIMD
and layout stop mattering — hence Add 1.02×, Sum 1.08×. The Arrow win lives only
in the mid-size band (1e5–1e6) where data is big enough to amortize fixed costs
but small enough to stay cache-resident. That band is narrow and is not where the
Polars gap lives.

**What would have to change to flip the answer to GO:**

- **arrow-go gains a grouped hash-aggregation kernel** (and scalar reductions for
  all types). This is the decisive one — it's what would let a storage migration
  actually attack the Polars gap. Not on the arrow-go roadmap today; would need
  upstream work or a pure-Go reimplementation independent of storage.
- **The chain win becomes large and monotonic in size** — e.g. ≥ 2× and *growing*
  through 1e7 — which would require the fixed per-kernel allocation overhead to be
  eliminated (kernel result reuse / arena allocation, buffer pooling) so longer
  chains amortize instead of accumulate.
- **enchanter stops handing users raw Go slices** (a `Series` API where data can
  stay Arrow across an entire analytical pipeline, materializing only at final
  output), removing the round-trip tax at op boundaries. This is a large public-API
  change orthogonal to storage.

Absent those, Arrow's value to enchanter is what 0.2.0 already delivered:
zero-copy interchange (Parquet, Arrow IPC, hand-off to Arrow-native consumers),
not an internal compute engine.

---

## 5. What the benchmark does **not** capture

- **Nulls in the compute paths.** Every benchmarked series is fully valid. The
  go-slice path skips its null-mask work when `!IsNullable_`; the arrow-native
  path skips validity-bitmap handling. Real nullable data adds validity buffers,
  branch-per-element cost, and — importantly — the enchanter↔Arrow null-convention
  translation (see §6) on every build. Nullable workloads would move *both* columns,
  in unknown proportion; the relative verdict could shift either way and is
  unmeasured.
- **Mutation-heavy workloads.** The prototype is read/compute only. `Set`, `Sort`,
  `Append`, `MakeNullable`, `Cast` are exactly the methods Arrow immutability makes
  expensive (each becomes a full rebuild — see §6). A workload dominated by those
  would favor the mutable Go-slice representation more than any benchmark here shows.
- **Groupby aggregation end-to-end.** The headline cross-library gap. Not
  benchmarked because **arrow-go has no grouped-aggregation kernel to compare
  against** — there is literally nothing to call. Storage format is not the lever
  here (§2).
- **Longer / branchier op chains.** Only a 2-op chain was measured. Given the ~11×
  alloc-count penalty already visible at two ops, longer chains are the scenario
  most likely to change the picture — plausibly *against* Arrow (accumulating
  per-kernel allocations) unless buffer reuse is added.
- **Multi-column / RecordBatch operations, string/dictionary columns, and
  parallelism.** Single Float64 column, single-threaded. Polars' real advantage is
  parallel multi-column execution; none of that is in scope here.
- **Steady-state vs cold allocator.** Uses `memory.DefaultAllocator` (GC-backed).
  A pooled/arena allocator would cut the arrow-native allocation cost and is the
  most likely thing to improve the arrow-native column — untested.

---

## 6. If it were GO: the migration design (for reference)

Recorded so a future revisit (e.g. after arrow-go ships aggregation kernels)
starts from a plan, not a blank page. **Not** being executed for 0.4.0.

### 6.1 `Float64s` holding an `arrow.Array` as source of truth

Today (`series/float64.go`):

```go
type Float64s struct {
    IsNullable_ bool
    Sorted_     enchanter.SeriesSortOrder
    Data_       []float64
    NullMask_   []uint8 // bit-packed, bit-SET = null
    Partition_  *SeriesFloat64Partition
    Ctx_        *enchanter.Context
}
```

Arrow-native shape:

```go
type Float64s struct {
    Sorted_    enchanter.SeriesSortOrder
    arr        *array.Float64   // source of truth; nullability + validity live here
    Partition_ *SeriesFloat64Partition
    Ctx_       *enchanter.Context
}
```

`IsNullable_` folds into `arr.NullN()`/the validity buffer; `Data_` and
`NullMask_` disappear. The good news: enchanter's series are already **value
types with copy-on-write semantics** — every method takes `s Float64s` *by value*
and returns a new `Series` (`func (s Float64s) Set(i int, v any) Series`). The
public contract is already "operations produce a new series," which is exactly
what Arrow immutability demands. The migration doesn't change the API shape; it
changes what the value wraps.

### 6.2 Copy-on-write / rebuild for the mutating methods

Arrow arrays are immutable, so each mutating method becomes builder-based rebuild:

- **`Set(i, v)`** — cannot poke `arr.Value(i)`. Rebuild: `array.NewFloat64Builder`,
  copy `[0,i)` from `arr`, append `v`, copy `(i,len)`, `NewFloat64Array()`. O(n)
  per set (today it's O(1) on `Data_[i]`). A batched `SetMany`/mutation-buffer API
  would be needed to keep row-wise editing loops from going quadratic — this is the
  single biggest ergonomic regression of the migration.
- **`Sort`/`SortRev`** — compute the permutation (or use `compute.Take` with sorted
  indices), producing a new array. Roughly parity with today (already allocates a
  sorted copy), but must recompute/relocate the validity bitmap.
- **`Append(v)`** — Arrow arrays don't grow. Either keep a *builder* alongside the
  array as a tail buffer and finalize lazily on first read, or rebuild
  `old + new`. A lazy builder tail is the only way to keep append-in-a-loop from
  being O(n²).
- **`MakeNullable`/`MakeNonNullable`** — today flips a bool and (de)allocates the
  mask. With Arrow, `MakeNullable` is a no-op if `arr` already carries a validity
  buffer, else rebuild with an all-valid bitmap; `MakeNonNullable` rebuilds only if
  `NullN() == 0` (else it's an error/no-op, same as today).
- **`Cast`** — `compute.CastArray` replaces the hand-rolled per-type conversion.

The through-line: mutation moves from O(1) in-place to O(n) rebuild. Acceptable
for occasional edits, quadratic for row-wise loops; a mutation-buffer/builder
escape hatch is mandatory, and its absence is a real risk to existing callers.

### 6.3 Migrating the ~49K generated lines, type-by-type, behind `Series`

The `*_base.go` + `*_ops.go` files total ~49K generated lines from
`generators/templates.go` + `generators/main.go` (`go generate ./series/`). The
`Series` interface (`series/series.go`) is the seam — nothing outside `series/`
depends on `Data_`/`NullMask_` *through the interface* (though `dataframe/` and
`io/` do reach into fields directly in places; those call sites must be inventoried
and moved onto `ArrowArray()`/`Data()` first). Plan:

1. **Ops become kernel dispatch.** The generator already emits manual per-type
   loops for `Add`/`Sub`/`Eq`/… A `series/arrow_ops.go` dispatch layer
   (`ArrowAdd`/`ArrowSub`/… over `compute`) already exists from the 0.2.0 Phase-4
   work but is not wired in. The generator's op templates would emit calls to that
   layer instead of hand-rolled loops, collapsing large swaths of the 49K lines.
2. **Type-by-type, one series at a time.** Migrate `Float64s` first (kernels exist
   and are best-covered), then `Int64s`/`Ints`, `Bools`, and last `Strings`
   (needs dictionary arrays — a separate epic) and `Times`/`Durations` (timestamp/
   duration arrays). Each type ships behind the unchanged interface so `dataframe/`
   sees no change; the generated golden-file diff is reviewed per type.
3. **Keep the generator idempotent.** `go generate` must stay a no-op diff after
   each step; the templates carry the migration, not hand-edits to generated files.

Feasible but large: it's a rewrite of the compute core plus a redesign of every
mutating method, staged across multiple releases. The interface seam makes it
*safe*; it does not make it *small*.

### 6.4 Null-handling story: opposite conventions

The load-bearing subtlety. enchanter uses a **bit-packed `[]uint8` where bit-SET =
null**; Arrow uses a **validity bitmap where bit-SET = valid** (the exact
inverse), byte-aligned per array with its own length/offset semantics. Today the
two are bridged by a per-element loop in `series/arrow_helpers.go`
(`buildArrowFloat64` et al.) and by `arrowutil` on the way back — an O(n) branch
per element on every build, and a correctness-critical inversion. Migrating
storage to Arrow **removes** this translation on the internal path (validity would
be native), which is a genuine simplification — but it forces every remaining
place that reads `NullMask_` directly (bit-set-null assumptions in
`Set`/`Append`/`GroupBy`/`io`) to be rewritten to Arrow validity semantics.
Getting the inversion wrong is a silent data-corruption class of bug (nulls
flipping to valid and vice-versa), so this is the highest-care part of the
migration and needs exhaustive round-trip null tests per type before any cutover.

### 6.5 Risks (if pursued)

- **Mutation quadratics.** O(n) rebuild per `Set`/`Append` silently turning
  row-wise loops in existing user code from O(n) into O(n²). Needs a batched
  mutation API and a deprecation/migration path for element-wise callers.
- **Null-convention inversion bugs.** Silent corruption if the bit-set-null →
  bit-set-valid flip is wrong anywhere. Highest-care; exhaustive per-type tests.
- **Allocation regression under real pipelines.** The measured ~11× alloc-count
  penalty on a 2-op chain could dominate GC in long pipelines; requires buffer
  pooling / arena allocation to be viable, which is itself a project.
- **Direct field access outside `series/`.** `dataframe/` and `io/` read `Data_`/
  `NullMask_` directly in spots; all such call sites must migrate to
  `ArrowArray()`/`Data()` before the fields can be removed.
- **String/dictionary and datetime types** are materially harder than the numeric
  types and would trail the migration by releases, leaving a mixed representation
  (some series Arrow-backed, some slice-backed) that the DataFrame layer must
  tolerate throughout.

---

## 7. Recommendation

**NO-GO for a storage migration in 0.4.0.** Keep Apache Arrow as the interop layer
it already is (Parquet, Arrow IPC, zero-copy hand-off), which is where it pays
unambiguously. The measured internal advantage is too small, too size-fragile, and
too confined to all-Arrow op chains to justify rewriting the compute core and every
mutating method, and it does not touch the grouped-aggregation gap that actually
separates enchanter from Polars — because arrow-go has no kernel for it.

**Cheaper wins to pursue instead** (independent of any storage migration):

- Capture SIMD `Sum` (and `Int64`/`Uint64` `Sum`) in the aggregators by summing
  over a temporary Arrow array where the input is large and cache-resident —
  self-contained, no rewrite.
- Attack the Polars gap where it lives: faster grouped-aggregation (better hash
  tables, the array-vs-map trick already noted in `notes.md` for dense integer
  keys) and goroutine parallelism across groups/columns.

**Revisit the go/no-go when** arrow-go ships a grouped hash-aggregation kernel (and
full-type scalar reductions), or when enchanter adopts a pipeline API that keeps
data Arrow end-to-end so the round-trip tax disappears. Either changes the inputs
to this decision materially; until then, the numbers say interop-only.
