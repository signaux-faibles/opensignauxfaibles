/*global globalThis*/

// Version TypeScript / AVA de public/_test.js, précedemment exécuté par jsc,
// lors de l'appel à `go test`, via dbmongo/js/test/test_public.sh.
//
// Usage: $ npx ava public/public_tests.ts
//     ou $ npm test

declare const f: any
const global = globalThis as any

import "../globals"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { flatten } from "./flatten"
import { diane } from "./diane"
import { bdf } from "./bdf"
import { iterable } from "./iterable"
import { effectifs } from "./effectifs"
import { sirene } from "./sirene"
import { dateAddMonth } from "../common/dateAddMonth"
import { cotisations } from "./cotisations"
import { dateAddDay } from "./dateAddDay"
import { debits } from "./debits"
import { apconso } from "./apconso"
import { delai } from "./delai"
import { compte } from "./compte"
import { procolToHuman } from "../common/procolToHuman"
import { dealWithProcols } from "./dealWithProcols"
import { map } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { objects as testCases } from "../test/data/objects"
import { reducer, invertedReducer } from "../test/helpers/reducers"
import { runMongoMap, indexMapResultsByKey } from "../test/helpers/mongodb"
import test from "ava"

global.f = {
  generatePeriodSerie,
  flatten,
  diane,
  bdf,
  iterable,
  effectifs,
  sirene,
  dateAddMonth,
  dateAddDay,
  debits,
  cotisations,
  apconso,
  delai,
  compte,
  procolToHuman,
  dealWithProcols,
  map,
  reduce,
  finalize,
}

test("la chaine d'intégration 'public' donne le même résultat que d'habitude", (t) => {
  const jsParams = global
  jsParams.offset_effectif = 2
  jsParams.actual_batch = "1905"
  jsParams.date_debut = new Date("2014-01-01")
  jsParams.date_fin = new Date("2018-02-01")
  jsParams.serie_periode = f.generatePeriodSerie(
    jsParams.date_debut,
    jsParams.date_fin
  )

  const pool = indexMapResultsByKey(runMongoMap(f.map, testCases))

  const intermediateResult = Object.values(pool).map((array) => ({
    key: array[0].key,
    value: reducer(array, f.reduce),
  }))

  const invertedIntermediateResult = Object.values(pool).map((array) => ({
    key: array[0].key,
    value: invertedReducer(array, f.reduce),
  }))

  const result = intermediateResult.map((r) => ({
    _id: r.key,
    value: f.finalize(r.key as { scope: Scope }, r.value),
  }))

  const invertedResult = invertedIntermediateResult.map((r) => ({
    _id: r.key,
    value: f.finalize(r.key as { scope: Scope }, r.value),
  }))

  t.deepEqual(result, invertedResult)
})
