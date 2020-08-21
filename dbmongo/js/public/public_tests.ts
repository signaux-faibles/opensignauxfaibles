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
import { flatten } from "../common/flatten"
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
import {
  runMongoMap,
  runMongoReduce,
  indexMapResultsByKey,
} from "../test/helpers/mongodb"
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

// initialisation des paramètres globaux de reduce.algo2
function initGlobalParams(dateDebut: Date, dateFin: Date) {
  const jsParams = global
  jsParams.offset_effectif = 2
  jsParams.actual_batch = "2002_1"
  jsParams.date_debut = dateDebut
  jsParams.date_fin = dateFin
  jsParams.serie_periode = f.generatePeriodSerie(dateDebut, dateFin)
}

test("l'ordre de traitement des données n'influe pas sur les résultats", (t) => {
  initGlobalParams(new Date("2014-01-01"), new Date("2018-02-01"))

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

// inspiré par reduce.algo2/map_tests.ts et reduce.algo2/algo2_golden_tests.ts
test("map() et reduce() retournent les même données que d'habitude", (t) => {
  initGlobalParams(new Date("2014-01-01"), new Date("2016-01-01"))

  const mapResult = runMongoMap(f.map, testCases)
  t.snapshot(mapResult)

  const reduceResult = runMongoReduce(f.reduce, mapResult)
  t.snapshot(reduceResult)
})
