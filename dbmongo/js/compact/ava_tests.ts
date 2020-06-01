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
import * as f from "../common/generatePeriodSerie.js"

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

// output data from test-api.sh
const expected = [
  {
    _id: "01234567891011",
    value: {
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
    },
  },
]

const runMongoMap = (mapFct: () => void, keyVal: object): object => {
  const results = {}
  globalThis.emit = (key: string, value: any): void => {
    results[key] = value
  }
  mapFct.call(keyVal)
  return results
}

test(`exécution complète de la chaine "compact"`, (t: ExecutionContext) => {
  // 1. map
  const mapResults = runMongoMap(map, importedData)
  const potentialMapResults = {
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
  t.deepEqual(mapResults, potentialMapResults)

  // 2. reduce
  const reduceKey = importedData.value.key
  const reduceValues = [mapResults[reduceKey]]
  const reduceResults = reduce(reduceKey, reduceValues)
  const potentialReduceResults = {
    batch: {
      1910: {},
    },
    key: "01234567891011",
    scope: "etablissement",
  }
  t.deepEqual(
    reduceResults,
    /*expected[0].value*/ (potentialReduceResults as unknown) as CompanyDataValues // TODO: update types to match data
  )

  // 3. finalize
  const global = globalThis as any
  global.serie_periode = f.generatePeriodSerie(
    ISODate("2015-12-01T00:00:00.000+0000"),
    ISODate("2016-02-01T00:00:00.000+0000")
  )
  const index: ReduceIndexFlags = { algo1: true, algo2: true }
  const finalizeKey = reduceKey
  const finalizeValues = { ...reduceResults, index }
  const finalizeResultValue = finalize(finalizeKey, finalizeValues)
  const finalizeResults = [
    { _id: finalizeKey, value: removeRandomOrder(finalizeResultValue) },
  ]
  t.deepEqual(finalizeResults, expected as unknown)
  // => sample of `actual` VS `expected`:
  //   -             'Tue Oct 01 2019 02:00:00 GMT+0200 (GMT+02:00)': {
  //   -               periode: Date 2019-10-01 00:00:00 UTC {},
  //   -               random_order: 0.19479352943685613,
  //   -               siret: '01234567891011',
  //   -             },
  //   -             'Wed Jan 01 2014 01:00:00 GMT+0100 (GMT+01:00)': {
  //   -               periode: Date 2014-01-01 00:00:00 UTC {},
  //   -               random_order: 0.6133162030905268,
  //   -               siret: '01234567891011',
  //   -             },
  //   +             'Fri Apr 01 2016 00:00:00 GMT+0000 (UTC)': {
  //   +               periode: Date 2016-04-01 00:00:00 UTC {},
  //   +               siret: '01234567891011',
  //   +             },
  //   +             'Fri Aug 01 2014 00:00:00 GMT+0000 (UTC)': {
  //   +               periode: Date 2014-08-01 00:00:00 UTC {},
  //   +               siret: '01234567891011',
  //   +             },
})
