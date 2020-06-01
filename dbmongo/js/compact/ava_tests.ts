// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import "../globals"
import { map } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"

const ISODate = (date: string): Date => new Date(date)

const removeRandomOrder = (obj: object): object => {
  Object.keys(obj).forEach(
    (key) =>
      (key === "random_order" && delete obj[key]) ||
      (typeof obj[key] === "object" && removeRandomOrder(obj[key]))
  )
  return obj
}

const runMongoMap = (mapFct: () => void, keyVal: object): object => {
  const results = {}
  globalThis.emit = (key: string, value: any): void => {
    results[key] = value
  }
  mapFct.call(keyVal)
  return results
}

// test data inspired by test-api.sh

const siret: SiretOrSiren = "01234567891011"
const scope: Scope = "etablissement"
const batchKey = "1910"
const dates = [
  ISODate("2015-12-01T00:00:00.000+0000"),
  ISODate("2016-01-01T00:00:00.000+0000"),
]
const batch: BatchValues = {
  [batchKey]: {} as any,
}

const importedData = {
  _id: "random123abc",
  value: {
    batch,
    scope,
    index: {
      algo2: true,
    },
    key: siret,
  },
}

const expectedMapResults = {
  [siret]: {
    batch,
    index: {
      algo2: true,
    },
    key: siret,
    scope,
  },
}

const expectedReduceResults = {
  batch,
  key: siret,
  scope,
}

const expectedFinalizeResultValue = {
  batch: {
    [batchKey]: {
      reporder: dates.reduce(
        (reporder, date) => ({
          ...reporder,
          [date.toString()]: { periode: date, siret },
        }),
        {}
      ),
    },
  },
  scope,
  index: {
    algo1: false,
    algo2: false,
  },
  key: siret,
} as unknown

// exécution complète de la chaine "compact"

test.serial(`compact.map()`, (t: ExecutionContext) => {
  const mapResults = runMongoMap(map, importedData)
  t.deepEqual(mapResults, expectedMapResults)
})

test.serial(`compact.reduce()`, (t: ExecutionContext) => {
  const reduceValues: CompanyDataValues[] = [expectedMapResults[siret]]
  const reduceResults = reduce(siret, reduceValues)
  t.deepEqual(reduceResults, expectedReduceResults)
})

test.serial(`compact.finalize()`, (t: ExecutionContext) => {
  const global = globalThis as any
  global.serie_periode = dates // used by complete_reporder(), which is called by finalize()
  const index: ReduceIndexFlags = { algo1: true, algo2: true }
  const finalizeValues = { ...expectedReduceResults, index }
  const finalizeResultValue = removeRandomOrder(finalize(siret, finalizeValues))
  t.deepEqual(finalizeResultValue, expectedFinalizeResultValue)
})
