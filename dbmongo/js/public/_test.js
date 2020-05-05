"use strict";

const globals = this; // global functions and parameters

globals.actual_batch = "1905"
globals.date_fin = new Date("2018-02-01")
globals.serie_periode = f.generatePeriodSerie(new Date("2014-01-01"), new Date("2018-02-01"))

objects.forEach(({ value, _id }) =>
  ({ ...globals, _id, value }).map()
)

var intermediateResult = Object.values(pool).map(array => ({
  key: array[0].key,
  value: reducer(array, f.reduce)
}))

var invertedIntermediateResult = Object.values(pool).map(array => ({
  key: array[0].key,
  value: invertedReducer(array, f.reduce)
}))

var result = intermediateResult.map(r => ({
  _id: r.key,
  value: f.finalize(r.key, r.value)
}))

var invertedResult = invertedIntermediateResult.map(r => ({
  _id: r.key,
  value: f.finalize(r.key, r.value)
}))

print(JSON.stringify(result) == JSON.stringify(invertedResult))
