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

const siret = "01234567891011"
const scope = "etablissement"
const batchKey = "1910"
const dates = [
  ISODate("2015-12-01T00:00:00.000+0000"),
  ISODate("2016-01-01T00:00:00.000+0000"),
]

const importedData = {
  _id: "random123abc",
  value: {
    batch: {
      [batchKey]: {},
    },
    scope,
    index: {
      algo2: true,
    },
    key: siret,
  },
}

const expectedMapResults = {
  [siret]: {
    batch: {
      [batchKey]: {},
    },
    index: {
      algo2: true,
    },
    key: siret,
    scope,
  },
}

const expectedReduceResults = {
  batch: { [batchKey]: {} },
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

test(`exécution complète de la chaine "compact"`, (t: ExecutionContext) => {
  // 1. map
  const mapResults = runMongoMap(map, importedData)
  t.deepEqual(mapResults, expectedMapResults)

  // 2. reduce
  const reduceValues: CompanyDataValues[] = [mapResults[siret]]
  const reduceResults = reduce(siret, reduceValues)
  t.deepEqual(
    reduceResults,
    /*expectedFinalizeResultValue*/ expectedReduceResults as unknown // TODO: update types to match data
  )

  // 3. finalize
  const global = globalThis as any
  global.serie_periode = dates
  const index: ReduceIndexFlags = { algo1: true, algo2: true }
  const finalizeValues = { ...reduceResults, index }
  const finalizeResultValue = removeRandomOrder(finalize(siret, finalizeValues))
  t.deepEqual(finalizeResultValue, expectedFinalizeResultValue)
})
