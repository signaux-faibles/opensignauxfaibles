"use strict"

// Context: this golden-file-based test runner was designed to prevent
// regressions on the JS functions (common + algo2) used to compute the
// "Features" collection from the "RawData" collection.
//
// It requires the JS functions from common + algo2 (notably: map()),
// and a makeTestData() function to generate a realistic test data set.
//
// Please execute ../test/test_algo2.sh to fill these requirements and
// run the tests.

import { makeTestData } from "./test_algo2_testdata"
import { naf } from "../test/data/naf"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { map } from "../reduce.algo2/map"
import { finalize } from "../reduce.algo2/finalize"
import { reduce } from "../reduce.algo2/reduce"
import { runMongoMap } from "../test/helpers/mongodb"

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any

const f = {
  generatePeriodSerie,
  map,
  finalize,
  reduce,
}

declare const console: any

// Define global parameters that are required by JS functions
const jsParams = global
jsParams.actual_batch = "2002_1"
jsParams.date_debut = new Date("2014-01-01")
jsParams.date_fin = new Date("2016-01-01")
jsParams.serie_periode = f.generatePeriodSerie(
  jsParams.date_debut,
  jsParams.date_fin
)
jsParams.includes = { all: true }
jsParams.offset_effectif = 2
jsParams.naf = naf
;(Object as any).bsonsize = (obj: unknown) => JSON.stringify(obj).length

// Generate a realistic test data set
const testData = makeTestData({
  ISODate: (date: string) => new Date(date.replace("+0000", "+00:00")), // make sure that timezone format complies with the spec
  NumberInt: (int: number) => int,
})

const mapResult = runMongoMap(
  f.map,
  testData as any[] // TODO: as { _id: string; value: CompanyDataValuesWithFlags }[]
) // -> [ { _id, value } ]

// Print the output of the f.map() function
console.log(JSON.stringify(mapResult, null, 2))

const valuesPerKey: Record<string, unknown[]> = {}
mapResult.forEach(({ _id, value }) => {
  const idString = JSON.stringify(_id)
  valuesPerKey[idString] = valuesPerKey[idString] || []
  valuesPerKey[idString].push(value)
})

const finalizeResult = Object.keys(valuesPerKey).map((key) =>
  f.finalize(JSON.parse(key), f.reduce(key, valuesPerKey[key]))
)

// Print the output of the f.finalize() function
console.log(JSON.stringify(finalizeResult, null, 2))
