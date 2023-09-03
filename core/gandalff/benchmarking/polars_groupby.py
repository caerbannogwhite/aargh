print("# filter-polars.py", flush=True)

import os
import gc
import timeit
import polars as pl

from statistics import mean
from helpers import write_log, memory_usage, make_chk

# other questions ans info here
# https://github.com/h2oai/db-benchmark/tree/master

ver = pl.__version__
git = ""
task = "group_by"
solution = "polars"
fun = ".group_by"
cache = "TRUE"
on_disk = "FALSE"
data_names = [
    "G1_1e4_1e2_0_0", "G1_1e5_1e2_0_0", "G1_1e6_1e2_0_0", "G1_1e7_1e2_0_0",
    "G1_1e4_1e2_10_0", "G1_1e5_1e2_10_0", "G1_1e6_1e2_10_0", "G1_1e7_1e2_10_0"
]

for data_name in data_names:
    filepath = os.path.join("..", "testdata", data_name+".csv")
    print("loading dataset %s" % data_name, flush=True)

    with pl.StringCache():
        x = pl.read_csv(filepath, dtypes={"id4":pl.Int32, "id5":pl.Int32, "id6":pl.Int32, "v1":pl.Int32, "v2":pl.Int32, "v3":pl.Float64}, low_memory=True)
        # x["id1"] = x["id1"].cast(pl.Categorical)
        # x["id1"].shrink_to_fit(in_place=True)
        # x["id2"] = x["id2"].cast(pl.Categorical)
        # x["id2"].shrink_to_fit(in_place=True)
        # x["id3"] = x["id3"].cast(pl.Categorical)
        # x["id3"].shrink_to_fit(in_place=True)

    in_rows = x.shape[0]
    x = x.lazy()

    print(len(x.collect()), flush=True)

    task_init = timeit.default_timer()
    print("grouping...", flush=True)


    ###################     QUESTION 1   ###################

    question = "sum v1 by id1" # q1
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id1").agg(pl.sum("v1")).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans["v1"].cast(pl.Int64).sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id1").agg(pl.sum("v1")).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans["v1"].cast(pl.Int64).sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 2   ###################

    question = "sum v1 by id1:id2" # q2
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by(["id1","id2"]).agg(pl.sum("v1")).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans["v1"].cast(pl.Int64).sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by(["id1","id2"]).agg(pl.sum("v1")).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans["v1"].cast(pl.Int64).sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 3   ###################

    question = "sum v1 mean v3 by id3" # q3
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id3").agg([pl.sum("v1"), pl.mean("v3")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v1").cast(pl.Int64).sum(), pl.col("v3").sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id3").agg([pl.sum("v1"), pl.mean("v3")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v1").cast(pl.Int64).sum(), pl.col("v3").sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 4   ###################

    question = "mean v1:v3 by id4" # q4
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id4").agg([pl.mean("v1"), pl.mean("v2"), pl.mean("v3")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v1").sum(), pl.col("v2").sum(), pl.col("v3").sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id4").agg([pl.mean("v1"), pl.mean("v2"), pl.mean("v3")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v1").sum(), pl.col("v2").sum(), pl.col("v3").sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 5   ###################

    question = "sum v1:v3 by id6" # q5
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id6").agg([pl.sum("v1"), pl.sum("v2"), pl.sum("v3")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v1").cast(pl.Int64).sum(), pl.col("v2").cast(pl.Int64).sum(), pl.col("v3").sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id6").agg([pl.sum("v1"), pl.sum("v2"), pl.sum("v3")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v1").cast(pl.Int64).sum(), pl.col("v2").cast(pl.Int64).sum(), pl.col("v3").sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 6   ###################

    question = "median v3 sd v3 by id4 id5" # q6
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by(["id4","id5"]).agg([pl.median("v3").alias("v3_median"), pl.std("v3").alias("v3_std")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v3_median").sum(), pl.col("v3_std").sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by(["id4","id5"]).agg([pl.median("v3").alias("v3_median"), pl.std("v3").alias("v3_std")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v3_median").sum(), pl.col("v3_std").sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 7   ###################

    question = "max v1 - min v2 by id3" # q7
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id3").agg([(pl.max("v1") - pl.min("v2")).alias("range_v1_v2")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans["range_v1_v2"].cast(pl.Int64).sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by("id3").agg([(pl.max("v1") - pl.min("v2")).alias("range_v1_v2")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans["range_v1_v2"].cast(pl.Int64).sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 8   ###################

    # question = "largest two v3 by id6" # q8
    # gc.collect()
    # t_start = timeit.default_timer()
    # ans = x.drop_nulls("v3").sort("v3", descending=True).group_by("id6").agg(col("v3").head(2).alias("largest2_v3")).explode("largest2_v3").collect()
    # print(ans.shape, flush=True)
    # t = timeit.default_timer() - t_start
    # m = memory_usage()
    # t_start = timeit.default_timer()
    # chk = [ans["largest2_v3"].sum()]
    # chkt = timeit.default_timer() - t_start
    # write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    # del ans
    # gc.collect()
    # t_start = timeit.default_timer()
    # ans = x.drop_nulls("v3").sort("v3", descending=True).group_by("id6").agg(col("v3").head(2).alias("largest2_v3")).explode("largest2_v3").collect()
    # print(ans.shape, flush=True)
    # t = timeit.default_timer() - t_start
    # m = memory_usage()
    # t_start = timeit.default_timer()
    # chk = [ans["largest2_v3"].sum()]
    # chkt = timeit.default_timer() - t_start
    # write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    # print(ans.head(3), flush=True)
    # print(ans.tail(3), flush=True)
    # del ans

    ###################     QUESTION 9   ###################

    # question = "regression v1 v2 by id2 id4" # q9
    # gc.collect()
    # t_start = timeit.default_timer()
    # ans = x.group_by(["id2","id4"]).agg((pl.pearson_corr("v1","v2")**2).alias("r2")).collect()
    # print(ans.shape, flush=True)
    # t = timeit.default_timer() - t_start
    # m = memory_usage()
    # t_start = timeit.default_timer()
    # chk = [ans["r2"].sum()]
    # chkt = timeit.default_timer() - t_start
    # write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    # del ans
    # gc.collect()
    # t_start = timeit.default_timer()
    # ans = x.group_by(["id2","id4"]).agg((pl.pearson_corr("v1","v2")**2).alias("r2")).collect()
    # print(ans.shape, flush=True)
    # t = timeit.default_timer() - t_start
    # m = memory_usage()
    # t_start = timeit.default_timer()
    # chk = [ans["r2"].sum()]
    # chkt = timeit.default_timer() - t_start
    # write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    # print(ans.head(3), flush=True)
    # print(ans.tail(3), flush=True)
    # del ans

    ###################     QUESTION 10   ###################

    question = "sum v3 count by id1:id6" # q10
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by(["id1","id2","id3","id4","id5","id6"]).agg([pl.sum("v3").alias("v3"), pl.count("v1").alias("count")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v3").sum(), pl.col("count").cast(pl.Int64).sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.group_by(["id1","id2","id3","id4","id5","id6"]).agg([pl.sum("v3").alias("v3"), pl.count("v1").alias("count")]).collect()
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = ans.lazy().select([pl.col("v3").sum(), pl.col("count").cast(pl.Int64).sum()]).collect().to_numpy()[0]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=in_rows, question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans
