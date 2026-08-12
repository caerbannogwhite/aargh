#!/usr/bin/env Rscript
# Cross-library groupby benchmark (R side) for enchanter.
# Q1 of the h2oai db-benchmark: sum(v1) by id1, on the G1 datasets.
# Self-contained: reads ../testdata/G1_*.csv, records library versions,
# prints tab-separated result rows to stdout. Best-of-N wall time.
suppressPackageStartupMessages({
  library(dplyr)
  library(data.table)
})

sizes <- c("1e6" = "G1_1e6_1e2_0_0", "1e7" = "G1_1e7_1e2_0_0")

# Adaptive timing: proc.time() on Windows can't resolve sub-~10ms ops, so run
# the operation enough times that total wall time is large, then divide.
timed <- function(f) {
  invisible(f()) # warm up
  reps <- 1
  repeat {
    t0 <- Sys.time()
    for (i in seq_len(reps)) invisible(f())
    dt <- as.numeric(Sys.time() - t0, units = "secs")
    if (dt >= 0.5 || reps >= 2e6) break
    reps <- reps * 4
  }
  dt / reps
}

cat("solution\tversion\trows\tquestion\ttime_sec\n")
dtver <- as.character(packageVersion("data.table"))
dpver <- as.character(packageVersion("dplyr"))

for (nm in names(sizes)) {
  path <- file.path("..", "testdata", paste0(sizes[[nm]], ".csv"))
  if (!file.exists(path)) next
  DT <- fread(path, showProgress = FALSE)
  n <- nrow(DT)

  t_dt <- timed(function() DT[, .(v1 = sum(v1)), by = id1])
  cat(sprintf("data.table\t%s\t%d\tQ1_sum_v1_by_id1\t%.6f\n", dtver, n, t_dt))

  tb <- as_tibble(DT)
  t_dp <- timed(function() tb %>% group_by(id1) %>% summarise(v1 = sum(v1), .groups = "drop"))
  cat(sprintf("dplyr\t%s\t%d\tQ1_sum_v1_by_id1\t%.6f\n", dpver, n, t_dp))
}
