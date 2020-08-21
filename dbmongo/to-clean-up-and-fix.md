# Clean-up TODO-list

```sh
$ go run honnef.co/go/tools/cmd/staticcheck --unused.whole-program=true -- ./...
```

## Files to remove?

lib/engine/R.go

## To fix: this value of err is never used (SA4006)

handlers.go:40:2: this value of err is never used (SA4006)
lib/altares/main.go:75:4: this value of err is never used (SA4006)
lib/crp/main.go:79:4: this value of err is never used (SA4006)
lib/engine/data.go:161:2: this value of err is never used (SA4006)
lib/engine/db.go:102:2: this value of err is never used (SA4006)
lib/engine/public.go:24:2: this value of err is never used (SA4006)
lib/marshal/mapping_test.go:41:3: this value of err is never used (SA4006)
lib/urssaf/debit.go:113:5: this value of err is never used (SA4006)
lib/urssaf/delai.go:145:2: this value of err is never used (SA4006)

## To fix: empty branch (SA9003)

lib/apconso/main.go:114:13: empty branch (SA9003)
lib/apdemande/main.go:154:13: empty branch (SA9003)
lib/crp/main.go:110:13: empty branch (SA9003)
lib/engine/main.go:247:5: empty branch (SA9003)
lib/naf/main.go:78:10: empty branch (SA9003)
lib/naf/main.go:98:10: empty branch (SA9003)
lib/naf/main.go:118:10: empty branch (SA9003)
lib/naf/main.go:137:10: empty branch (SA9003)
lib/naf/main.go:156:10: empty branch (SA9003)
lib/urssaf/ccsf.go:128:13: empty branch (SA9003)
lib/urssaf/procol.go:109:13: empty branch (SA9003)

## To fix: should check returned error before deferring *.Close() (SA5001)

lib/diane/main.go:364:2: should check returned error before deferring stdout.Close() (SA5001)
lib/diane/main.go:371:2: should check returned error before deferring stderr.Close() (SA5001)
lib/marshal/filter.go:76:3: should check returned error before deferring file.Close() (SA5001)
lib/marshal/mapping.go:93:2: should check returned error before deferring file.Close() (SA5001)
lib/marshal/marshal.go:107:2: should check returned error before deferring file.Close() (SA5001)
lib/urssaf/cotisation.go:67:4: should check returned error before deferring file.Close() (SA5001)

## To fix: error strings should not be capitalized (ST1005)

lib/engine/adminBatch.go:71:10: error strings should not be capitalized (ST1005)
lib/engine/datapi.go:284:10: error strings should not be capitalized (ST1005)
lib/engine/db.go:48:10: error strings should not be capitalized (ST1005)
lib/engine/event.go:42:10: error strings should not be capitalized (ST1005)
lib/engine/filter.go:48:17: error strings should not be capitalized (ST1005)
lib/engine/reduce.go:50:10: error strings should not be capitalized (ST1005)
lib/engine/reduce.go:258:15: error strings should not be capitalized (ST1005)
lib/marshal/filter.go:49:16: error strings should not be capitalized (ST1005)
lib/marshal/mapping.go:57:16: error strings should not be capitalized (ST1005)
lib/marshal/mapping.go:69:15: error strings should not be capitalized (ST1005)
lib/marshal/marshal.go:196:40: error strings should not be capitalized (ST1005)
lib/marshal/parse.go:295:18: error strings should not be capitalized (ST1005)
lib/marshal/parse.go:300:26: error strings should not be capitalized (ST1005)
lib/marshal/readOptions.go:68:11: error strings should not be capitalized (ST1005)
lib/misc/main.go:62:23: error strings should not be capitalized (ST1005)
lib/sirene/main.go:125:17: error strings should not be capitalized (ST1005)
lib/urssaf/ccsf.go:133:20: error strings should not be capitalized (ST1005)
lib/urssaf/misc.go:47:18: error strings should not be capitalized (ST1005)
lib/urssaf/misc.go:52:26: error strings should not be capitalized (ST1005)

## To fix: should merge variable declaration with assignment on next line (S1021)

lib/marshal/parse.go:120:2: should merge variable declaration with assignment on next line (S1021)
lib/marshal/parse.go:139:2: should merge variable declaration with assignment on next line (S1021)
lib/marshal/parse.go:176:2: should merge variable declaration with assignment on next line (S1021)
lib/marshal/parse.go:195:2: should merge variable declaration with assignment on next line (S1021)
lib/marshal/parse.go:213:2: should merge variable declaration with assignment on next line (S1021)

## Others

lib/diane/main.go:432:2: redundant return statement (S1023)