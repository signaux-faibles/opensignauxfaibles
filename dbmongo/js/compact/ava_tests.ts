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

const renderedDate = (d: string): string => new Date(d).toString()

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

// input data from test-api.sh
const importedData = {
  _id: "random123abc",
  value: {
    batch: {
      "1910": {},
    },
    scope: "etablissement",
    index: {
      algo2: true,
    },
    key: "01234567891011",
  },
}

// output data inspired by test-api.sh
const expectedFinalizeResultValue = {
      batch: {
        "1910": {
          reporder: {
            [renderedDate("Tue Dec 01 2015 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2015-12-01T00:00:00Z"),
              siret: "01234567891011",
            },
            [renderedDate("Fri Jan 01 2016 00:00:00 GMT+0000 (UTC)")]: {
              periode: ISODate("2016-01-01T00:00:00Z"),
              siret: "01234567891011",
            },
          },
        },
      },
      scope: "etablissement",
      index: {
        algo1: false,
        algo2: false,
      },
      key: "01234567891011",
} as unknown

test(`exécution complète de la chaine "compact"`, (t: ExecutionContext) => {
  // 1. map
  const mapResults = runMongoMap(map, importedData)
  const expectedMapResults = {
    "01234567891011": {
      batch: {
        1910: {},
      },
      index: {
        algo2: true,
      },
      key: "01234567891011",
      scope: "etablissement",
    },
  }
  t.deepEqual(mapResults, expectedMapResults)

  // 2. reduce
  const reduceKey = importedData.value.key
  const reduceValues = [mapResults[reduceKey]]
  const reduceResults = reduce(reduceKey, reduceValues)
  const expectedReduceResults = {
    batch: {
      1910: {},
    },
    key: "01234567891011",
    scope: "etablissement",
  }
  t.deepEqual(
    reduceResults,
    /*expectedFinalizeResults[0].value*/ (expectedReduceResults as unknown) as CompanyDataValues // TODO: update types to match data
  )

  // 3. finalize
  const global = globalThis as any
  global.serie_periode = [
    ISODate("2015-12-01T00:00:00.000+0000"),
    ISODate("2016-01-01T00:00:00.000+0000"),
  ]
  const index: ReduceIndexFlags = { algo1: true, algo2: true }
  const finalizeKey = reduceKey
  const finalizeValues = { ...reduceResults, index }
  const finalizeResultValue = removeRandomOrder(
    finalize(finalizeKey, finalizeValues)
  )
  t.deepEqual(finalizeResultValue, expectedFinalizeResultValue)
})
