print("# filter-pandas.py", flush=True)

import os
import gc
import timeit
import pandas as pd

from statistics import mean
from helpers import write_log, memory_usage, make_chk

# other questions ans info here
# https://github.com/h2oai/db-benchmark/tree/master

ver = pd.__version__
git = ""
task = "groupby"
solution = "pandas"
fun = ".groupby"
cache = "TRUE"
on_disk = "FALSE"
data_names = [
    "G1_1e4_1e2_0_0", "G1_1e5_1e2_0_0", "G1_1e6_1e2_0_0", "G1_1e7_1e2_0_0",
    "G1_1e4_1e2_10_0", "G1_1e5_1e2_10_0", "G1_1e6_1e2_10_0", "G1_1e7_1e2_10_0"
]

for data_name in data_names:
    filepath = os.path.join("..", "testdata", data_name+".csv")
    print("loading dataset %s" % data_name, flush=True)

    x = pd.read_csv(filepath)

    ###################     QUESTION 1   ###################

    question = "sum v1 by id1" # q1
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id1', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'sum'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id1', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'sum'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 2   ###################

    question = "sum v1 by id1:id2" # q2
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby(['id1','id2'], as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'sum'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby(['id1','id2'], as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'sum'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 3   ###################

    question = "sum v1 mean v3 by id3" # q3
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id3', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'sum', 'v3':'mean'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum(), ans['v3'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id3', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'sum', 'v3':'mean'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum(), ans['v3'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 4   ###################

    question = "mean v1:v3 by id4" # q4
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id4', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'mean', 'v2':'mean', 'v3':'mean'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum(), ans['v2'].sum(), ans['v3'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id4', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'mean', 'v2':'mean', 'v3':'mean'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum(), ans['v2'].sum(), ans['v3'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 5   ###################

    question = "sum v1:v3 by id6" # q5
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id6', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'sum', 'v2':'sum', 'v3':'sum'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum(), ans['v2'].sum(), ans['v3'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id6', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'sum', 'v2':'sum', 'v3':'sum'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v1'].sum(), ans['v2'].sum(), ans['v3'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 6   ###################

    question = "median v3 sd v3 by id4 id5" # q6
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby(['id4','id5'], as_index=False, sort=False, observed=True, dropna=False).agg({'v3': ['median','std']})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v3']['median'].sum(), ans['v3']['std'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby(['id4','id5'], as_index=False, sort=False, observed=True, dropna=False).agg({'v3': ['median','std']})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v3']['median'].sum(), ans['v3']['std'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 7   ###################

    question = "max v1 - min v2 by id3" # q7
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id3', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'max', 'v2':'min'}).assign(range_v1_v2=lambda x: x['v1']-x['v2'])[['id3','range_v1_v2']]
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['range_v1_v2'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby('id3', as_index=False, sort=False, observed=True, dropna=False).agg({'v1':'max', 'v2':'min'}).assign(range_v1_v2=lambda x: x['v1']-x['v2'])[['id3','range_v1_v2']]
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['range_v1_v2'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans

    ###################     QUESTION 8   ###################

    # question = "largest two v3 by id6" # q8
    # gc.collect()
    # t_start = timeit.default_timer()
    # ans = x[~x['v3'].isna()][['id6','v3']].sort_values('v3', ascending=False).groupby('id6', as_index=False, sort=False, observed=True, dropna=False).head(2)
    # ans.reset_index(drop=True, inplace=True)
    # print(ans.shape, flush=True)
    # t = timeit.default_timer() - t_start
    # m = memory_usage()
    # t_start = timeit.default_timer()
    # chk = [ans['v3'].sum()]
    # chkt = timeit.default_timer() - t_start
    # write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    # del ans
    # gc.collect()
    # t_start = timeit.default_timer()
    # ans = x[~x['v3'].isna()][['id6','v3']].sort_values('v3', ascending=False).groupby('id6', as_index=False, sort=False, observed=True, dropna=False).head(2)
    # ans.reset_index(drop=True, inplace=True)
    # print(ans.shape, flush=True)
    # t = timeit.default_timer() - t_start
    # m = memory_usage()
    # t_start = timeit.default_timer()
    # chk = [ans['v3'].sum()]
    # chkt = timeit.default_timer() - t_start
    # write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    # print(ans.head(3), flush=True)
    # print(ans.tail(3), flush=True)
    # del ans

    ###################     QUESTION 9   ###################

    # question = "regression v1 v2 by id2 id4" # q9
    # #corr().iloc[0::2][['v2']]**2 # on 1e8,k=1e2 slower, 76s vs 47s
    # gc.collect()
    # t_start = timeit.default_timer()
    # ans = x[['id2','id4','v1','v2']].groupby(['id2','id4'], as_index=False, sort=False, observed=True, dropna=False).apply(lambda x: pd.Series({'r2': x.corr()['v1']['v2']**2}))
    # print(ans.shape, flush=True)
    # t = timeit.default_timer() - t_start
    # m = memory_usage()
    # t_start = timeit.default_timer()
    # chk = [ans['r2'].sum()]
    # chkt = timeit.default_timer() - t_start
    # write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    # del ans
    # gc.collect()
    # t_start = timeit.default_timer()
    # ans = x[['id2','id4','v1','v2']].groupby(['id2','id4'], as_index=False, sort=False, observed=True, dropna=False).apply(lambda x: pd.Series({'r2': x.corr()['v1']['v2']**2}))
    # print(ans.shape, flush=True)
    # t = timeit.default_timer() - t_start
    # m = memory_usage()
    # t_start = timeit.default_timer()
    # chk = [ans['r2'].sum()]
    # chkt = timeit.default_timer() - t_start
    # write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    # print(ans.head(3), flush=True)
    # print(ans.tail(3), flush=True)
    # del ans

    ###################     QUESTION 10   ###################

    question = "sum v3 count by id1:id6" # q10
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby(['id1','id2','id3','id4','id5','id6'], as_index=False, sort=False, observed=True, dropna=False).agg({'v3':'sum', 'v1':'size'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v3'].sum(), ans['v1'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=1, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    del ans
    gc.collect()
    t_start = timeit.default_timer()
    ans = x.groupby(['id1','id2','id3','id4','id5','id6'], as_index=False, sort=False, observed=True, dropna=False).agg({'v3':'sum', 'v1':'size'})
    print(ans.shape, flush=True)
    t = timeit.default_timer() - t_start
    m = memory_usage()
    t_start = timeit.default_timer()
    chk = [ans['v3'].sum(), ans['v1'].sum()]
    chkt = timeit.default_timer() - t_start
    write_log(task=task, data=data_name, in_rows=x.shape[0], question=question, out_rows=ans.shape[0], out_cols=ans.shape[1], solution=solution, version=ver, git=git, fun=fun, run=2, time_sec=t, mem_gb=m, cache=cache, chk=make_chk(chk), chk_time_sec=chkt, on_disk=on_disk)
    print(ans.head(3), flush=True)
    print(ans.tail(3), flush=True)
    del ans