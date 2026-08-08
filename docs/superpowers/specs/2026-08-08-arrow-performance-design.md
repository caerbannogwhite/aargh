# enchanter 0.3.0 — "measure the Arrow thesis"

Status: approved design
Date: 2026-08-08
Author: caerbannogwhite (with Claude)

## Summary

0.3.0's headline is an empirical question, not a feature: **would an
Arrow-backed storage model actually make enchanter faster in Go?** 0.2.0 added
Apache Arrow as an *interop* layer (series and dataframes convert to and from
Arrow; Parquet and Arrow IPC I/O). The natural next step is to use Arrow as the
*compute* engine — but whether that pays off in Go is unproven and non-obvious.
This release answers the question with measurements and a throwaway prototype
before committing to a large migration, and ships three smaller independent
improvements alongside.

## Background: the arrow-go constraint

Verified against the vendored dependency (`arrow-go v18.5.2`): its `compute`
package registers only scalar **arithmetic, comparison, boolean, cast,
set-lookup** kernels and the vector **hash** (`unique`, `value_counts`,
`dictionary_encode`) and **selection** (`Filter`, `Take`) kernels. It registers
**no aggregation kernels** — there is no `sum`, `mean`, `min`, `max`, `stddev`,
and no grouped/hash aggregation. The only accelerated reduction anywhere in
arrow-go is `math.Float64/Int64/Uint64.Sum` (SIMD sum; Sum only, three types).

Consequences that shape this release:

- "Route aggregations through Arrow compute" is **not achievable** with this
  dependency. The current aggregation code is already single-pass O(n) and Arrow
  cannot beat it without kernels. (Getting real Arrow aggregation would require a
  cgo binding to Arrow C++, which breaks enchanter's pure-Go promise — out of
  scope, permanently.)
- The Arrow operations that *are* available — element-wise compute and
  `Filter`/`Take` — pay a **build → compute → materialize round-trip** when data
  lives in Go slices (which it does today). On a one-off call that is *slower*
  than the current tight Go loop. They only win when data is **already an Arrow
  array and stays Arrow across a chain of operations**.

Therefore the only real Arrow-performance lever is making **Arrow the series
storage backend**. That is a large, non-obvious bet, so 0.3.0 measures it rather
than committing to it.

## Goals

1. Produce a **numbers-backed go/no-go** on whether Arrow-native series storage
   is worth migrating to (decision lands for 0.4.0).
2. Fix known aggregation **correctness bugs**.
3. Clear **foundation debt** (deprecations, CI, lint) at the start of the cycle.

## Non-goals (candidates for 0.4.0+)

- The full storage migration itself.
- Wiring the existing element-wise compute dispatch into the generated ops.
- Dictionary-encoded (factor) strings.
- SAS7BDAT data reading.
- Pivot longer/wider (the `dev-pivot` branch).

---

## Section 1 — Foundation cleanup

Independent and low-risk; lands first as its own PR so the API churn happens once
at the start of the cycle.

- **`workflow_dispatch` CI trigger** — re-adds the change from the closed PR #9,
  so CI can be re-run on demand (e.g. after a hosted-runner outage).
- **`arrow.Record` → `RecordBatch` migration.** `arrow-go v18.5` deprecated
  `arrow.Record`, `array.NewRecord`, and `reader.Record` in favour of the
  `RecordBatch` names. Used in ~9 spots across `io/` and `dataframe/`, including
  the **public** `DataFrame.ToArrowRecord() arrow.Record` and
  `NewBaseDataFrameFromArrowRecord(record arrow.Record, ...)`. This is a
  deliberate breaking change to the public API; it goes in the CHANGELOG under a
  `[Changed]` entry.
- **Lint & generated-code hygiene:** add a `golangci-lint` config and triage the
  legacy findings (fix or explicitly disable, no silent noise); emit the
  canonical `// Code generated ... DO NOT EDIT.` header on the `*_ops.go` files
  (currently only `*_base.go` carry it); remove the dead `printInfo` method from
  the `Series` interface and its generated implementations.

Acceptance: CI green including the new lint job; generated code still
byte-idempotent; public-API rename reflected in docs/CHANGELOG.

## Section 2 — Aggregation correctness fixes

Plain-Go fixes (Arrow offers no help) for two real bugs in `dataframe/stats.go`,
each pinned by a new failing-then-passing test:

- **Grouped `std`** divides the summed squared deviations by
  `len(flatGroupIndeces)/groupsNum` — the *average* group size — instead of each
  group's actual size. Wrong whenever groups are unequal. Fix: divide each
  group's accumulator by that group's own count.
- **Ungrouped `mean` with `removeNAs`** divides the sum by `len(dataF64)` (total
  length) instead of the count of non-NaN values. Fix: divide by the non-null
  count.

While here, add a small correctness test matrix for `Sum/Min/Max/Mean/Std` over
grouped and ungrouped inputs, with and without nulls, so these paths gain the
coverage they currently lack. Note in the CHANGELOG that aggregation results
change for the affected cases (bug fix, not a feature).

## Section 3 — Arrow-performance prototype + measurement (headline)

### Question

Does computing on Arrow arrays beat the current Go-slice path by enough to
justify migrating series storage — accounting for Go interface dispatch, escape
analysis, and the build/materialize round-trip?

### Method: three regimes

Each benchmarked operation is measured in three regimes:

1. **go-slice** — today's `[]float64` + null-mask implementation (baseline).
2. **arrow round-trip** — build an Arrow array from the Go slice, run the compute
   kernel, materialize the result back to a Go slice. Represents the cost of
   naïvely wiring compute into *today's* storage.
3. **arrow-native** — inputs are already `arrow.Array`; the result stays an
   `arrow.Array`; no build, no materialize. Represents what storage-native would
   enable (compute-only cost).

### Operations

Chosen because arrow-go actually provides the kernel and each is a hot path:

- **Add** — element-wise `a + b` (`compute` arithmetic vs Go loop).
- **Gt → Filter** — comparison producing a boolean mask, then selection
  (`greater` + `Filter` vs Go loop + manual take). This is the filter/groupby
  hot path.
- **Sum** — reduction (`math.Float64.Sum` vs Go loop).
- **2-op chain** — `(a + b) > k` then filter. The decisive case: arrow-native
  amortizes its single build/materialize only when data stays Arrow across
  multiple ops, so this is where it should pull ahead if it ever does.

Sizes: 1e4, 1e5, 1e6, 1e7 (1e7 skipped under `-short`). Report `ns/op`, `B/op`,
`allocs/op` — allocations matter because Arrow's likely win is fewer allocations
and its likely loss is the build copy.

### Prototype

`series/arrow_native_float64.go` — an `ArrowFloat64s` type wrapping an
`arrow.Array` (Float64) that implements just enough of the `Series` interface
(`Len`, `Get`, `Add`, `Gt`/comparison, `Filter`, `Sum`, `ArrowArray`) to run the
arrow-native regime **through the real interface**, so the numbers include
interface-dispatch overhead rather than measuring raw kernel calls in isolation.

Constraints:
- Measurement-only. **Not** wired into the code generator or the existing
  dispatch; touches nothing in the shipping code path.
- Mutation is out of scope. Arrow arrays are immutable, so `Set`/`Sort` and other
  mutating methods are left unimplemented (panic or `Errors`) — the immutability
  question is analysed in the design doc, not solved in the prototype.

### Success criterion (explicit go/no-go)

- **GO** (design the full migration for 0.4.0): arrow-native beats go-slice by
  **≥ ~1.5×** on the 2-op chain at **≥ 1e6 rows**, without badly regressing
  single ops.
- **NO-GO**: otherwise. Document why; Arrow stays interop-only; pursue other
  performance avenues.

The threshold is a guide, not a tripwire — the deliverable is a reasoned
recommendation grounded in the table of measurements, not a single pass/fail bit.

### Deliverables

- `series/arrow_native_float64.go` — the prototype.
- `series/arrow_bench_test.go` — the three-regime benchmarks.
- `docs/superpowers/specs/2026-08-08-arrow-native-storage-migration.md` — the
  migration design doc: immutability / copy-on-write strategy for mutating
  methods, how the code generator and the 44K generated lines would be handled, a
  phased migration order, risks, and the measured go/no-go recommendation with the
  benchmark table. Written during Section 3 implementation, once the numbers
  exist.

### What this does not measure

Groupby aggregation end-to-end (arrow-go has no grouped-aggregation kernels, so
there is nothing to compare against), and mutation-heavy workloads (out of scope
for the prototype). Both are noted as open questions in the design doc. The
cross-library comparison harness on the `dev-benchmarking` branch (polars /
pandas / dplyr) is a separate concern and is not entangled with these internal
microbenchmarks.

---

## Risks

- **Prototype not representative.** Running through the `Series` interface and
  reporting allocations mitigates the "microbenchmark lies" failure mode, but the
  prototype omits nulls handling in some paths; the design doc must state which
  costs are and are not captured.
- **Inconclusive result** (win margin near the threshold). Acceptable outcome:
  the doc records it as NO-GO-for-now with the specific conditions under which it
  would flip, rather than forcing a decision.
- **Deprecation rename churn.** The `Record` → `RecordBatch` change is public;
  doing it first and in one PR contains the blast radius.

## Sequencing

1. Section 1 (foundation) — merge first.
2. Section 2 (agg fixes) — independent, any time.
3. Section 3 (prototype + measure + design doc) — the headline; its design doc
   feeds the 0.4.0 go/no-go.
