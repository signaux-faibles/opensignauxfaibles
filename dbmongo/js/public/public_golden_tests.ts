/*global globalThis*/

import "../globals"
declare function emit(k: unknown, v: unknown): void
declare const f: any
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
import { finalize } from "./finalize"
;(globalThis as any).f = {
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
  finalize,
}

import test from "ava"
import "./reduce"
import "./finalize"
import { objects as testCases } from "../test/data/objects"
import { reducer, invertedReducer } from "../test/helpers/reducers"
import { runMongoMap, indexMapResultsByKey } from "../test/helpers/mongodb"

test("la chaine d'intégration 'public' donne le même résultat que d'habitude", (t) => {
  const jsParams = globalThis as any // => all properties of this object will become global.
  jsParams.actual_batch = "1905"
  jsParams.date_debut = new Date("2014-01-01")
  jsParams.date_fin = new Date("2018-02-01")
  jsParams.serie_periode = f.generatePeriodSerie(
    new Date("2014-01-01"),
    new Date("2018-02-01")
  )
  jsParams.offset_effectif = 2

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
