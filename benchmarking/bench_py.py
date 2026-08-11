#!/usr/bin/env python
# Cross-library groupby benchmark (Python side) for enchanter.
# Q1 of the h2oai db-benchmark: sum(v1) by id1, on the G1 datasets.
# Self-contained: reads ../testdata/G1_*.csv, records library versions,
# prints tab-separated result rows to stdout. Adaptive best-average timing.
import os
import gc
import time

import polars as pl
import pandas as pd

sizes = {"1e6": "G1_1e6_1e2_0_0", "1e7": "G1_1e7_1e2_0_0"}


def timed(f):
    f()  # warm up
    reps = 1
    while True:
        gc.collect()
        t0 = time.perf_counter()
        for _ in range(reps):
            f()
        dt = time.perf_counter() - t0
        if dt >= 0.5 or reps >= 2_000_000:
            return dt / reps
        reps *= 4


print("solution\tversion\trows\tquestion\ttime_sec")
for nm, base in sizes.items():
    path = os.path.join("..", "testdata", base + ".csv")
    if not os.path.exists(path):
        continue

    pf = pl.read_csv(path)
    n = pf.height
    t_pl = timed(lambda: pf.group_by("id1").agg(pl.col("v1").sum()))
    print(f"polars\t{pl.__version__}\t{n}\tQ1_sum_v1_by_id1\t{t_pl:.6f}", flush=True)
    del pf
    gc.collect()

    pdf = pd.read_csv(path)
    t_pd = timed(lambda: pdf.groupby("id1", sort=False)["v1"].sum())
    print(f"pandas\t{pd.__version__}\t{len(pdf)}\tQ1_sum_v1_by_id1\t{t_pd:.6f}", flush=True)
    del pdf
    gc.collect()
