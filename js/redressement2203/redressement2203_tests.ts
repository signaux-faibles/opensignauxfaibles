// Tests des fonctions map(), reduce() et finalize() de Public,
// en les faisant tourner sur les données de test/data/objects.js.
//
// Important: ce jeu de données ne couvre pas tous les types d'entrées.

import { map, Input, OutKey, OutValue } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { objects as testCases } from "../test/data/objects_redressement2203"
import { setGlobals } from "../test/helpers/setGlobals"
import { runMongoMap, runMongoReduce } from "../test/helpers/mongodb"
import test from "ava"

// initialisation des paramètres globaux de reduce.algo2
const initGlobalParams = (date: string) =>
  setGlobals({
    dateStr: date,
  })

// inspiré par reduce.algo2/map_tests.ts et reduce.algo2/algo2_golden_tests.ts
test("map(), reduce() et finalize() retournent les même données que d'habitude", (t) => {
  initGlobalParams("2016-01-01")

  const mapResult = runMongoMap<Input, OutKey, OutValue>(map, testCases)
  t.snapshot(mapResult)

  const reduceResult = runMongoReduce(reduce, mapResult)
  t.snapshot(reduceResult)

  const finalizeResult = reduceResult.map(({ _id, value }) =>
    finalize(_id, value)
  )
  t.snapshot(finalizeResult)
})
