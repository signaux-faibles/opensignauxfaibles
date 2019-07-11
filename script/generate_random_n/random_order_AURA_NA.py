#!/usr/bin/env python3
import sys
import random

f = sys.stdin
paths = ['../../data-raw/sirene/juin2019/NA/sirene_NA.csv',\
        '../../data-raw/sirene/juin2019/AURA/sirene_AURA.csv']
seed = 7541
for path in paths:
    f = open(path)
    random.seed(seed)
    periods= ["2015-01-01", "2015-02-01", "2015-03-01", "2015-04-01",\
            "2015-05-01", "2015-06-01", "2015-07-01", "2015-08-01", "2015-09-01",\
            "2015-10-01", "2015-11-01", "2015-12-01", "2016-01-01", "2016-02-01",\
            "2016-03-01", "2016-04-01", "2016-05-01", "2016-06-01", "2016-07-01",\
            "2016-08-01", "2016-09-01", "2016-10-01", "2016-11-01", "2016-12-01"]
    for line in f:
        fields = line.strip().split(",")
        if fields[0] == "SIREN":
            continue
        for period in periods:
            print(fields[0] + fields[1], period, random.random(), sep=',')
