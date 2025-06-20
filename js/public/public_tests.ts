// Tests des fonctions map(), reduce() et finalize() de Public,
// en les faisant tourner sur les données de test/data/objects.js.
//
// Important: ce jeu de données ne couvre pas tous les types d'entrées.

import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { map, Input, OutKey, OutValue } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { objects as testCases } from "../test/data/objects"
import { setGlobals } from "../test/helpers/setGlobals"
import { reducer, invertedReducer } from "../test/helpers/reducers"
import {
  runMongoMap,
  runMongoReduce,
  indexMapResultsByKey,
} from "../test/helpers/mongodb"
import test from "ava"

// initialisation des paramètres globaux de reduce.algo2
const initGlobalParams = (dateDebut: Date, dateFin: Date) =>
  setGlobals({
    actual_batch: "2002_1",
    date_fin: dateFin,
    serie_periode: generatePeriodSerie(dateDebut, dateFin),
  })

test("l'ordre de traitement des données n'influe pas sur les résultats", (t: test) => {
  initGlobalParams(new Date("2014-01-01"), new Date("2018-02-01"))

  const pool = indexMapResultsByKey(
    runMongoMap<Input, OutKey, OutValue>(map, testCases)
  )

  const intermediateResult = Object.values(pool).map((array) => ({
    key: array[0]?.key,
    value: reducer(array, reduce),
  }))

  const invertedIntermediateResult = Object.values(pool).map((array) => ({
    key: array[0]?.key,
    value: invertedReducer(array, reduce),
  }))

  const result = intermediateResult.map((r) => ({
    _id: r.key,
    value: finalize(r.key, r.value),
  }))

  const invertedResult = invertedIntermediateResult.map((r) => ({
    _id: r.key,
    value: finalize(r.key, r.value),
  }))

  t.deepEqual(result, invertedResult)
})

// inspiré par reduce.algo2/map_tests.ts et reduce.algo2/algo2_golden_tests.ts
test("map(), reduce() et finalize() retournent les même données que d'habitude", (t: test) => {
  initGlobalParams(new Date("2014-01-01"), new Date("2016-01-01"))

  const mapResult = runMongoMap<Input, OutKey, OutValue>(map, testCases)
  t.snapshot(mapResult)

  const reduceResult = runMongoReduce(reduce, mapResult)
  t.snapshot(reduceResult)

  const finalizeResult = reduceResult.map(({ _id, value }) =>
    finalize(_id, value)
  )
  t.snapshot(finalizeResult)
})
