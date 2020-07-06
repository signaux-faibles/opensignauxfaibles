// Converted to TS from _test.js

import test, { ExecutionContext } from "ava"
import { map } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { generatePeriodSerie } from "../common/generatePeriodSerie"

// Importation du jeu de données
import { objects as testCases } from "../test/data/objects"
import { naf as nafValues } from "../test/data/naf"

// Paramètres globaux utilisés par "reduce.algo2"
declare let f: unknown //{ generatePeriodSerie: Function; map: Function }
declare let emit: unknown // called by map()
declare let naf: NAF
declare let actual_batch: BatchKey
declare let date_debut: Date
declare let date_fin: Date
declare let serie_periode: Date[]
declare let offset_effectif: number
declare let includes: Record<"all", boolean>

test("algo2_tests depuis objects", (t: ExecutionContext) => {
  testCases.forEach(({ _id, value }) => {
    // preparation de l'environnement d'exécution de map()
    const pool: Record<any, any> = {}
    emit = (key: any, value: any) => {
      const id = key.siren + key.batch + key.periode.getTime()
      pool[id] = (pool[id] || []).concat([{ key, value }])
    }
    // initialisation des paramètres globaux de reduce.algo2
    naf = nafValues
    actual_batch = "1905"
    date_debut = new Date("2014-01-01")
    date_fin = new Date("2018-02-01")
    serie_periode = generatePeriodSerie(
      new Date("2014-01-01"),
      new Date("2018-02-01")
    )
    offset_effectif = 2
    includes = { all: true }
    // exécution du test
    map.call({ _id, value /*, ...f*/ }) // will populate pool
    const intermediateResult = objectValues(pool).map((array) =>
      reducer(array, reduce)
    )

    const invertedIntermediateResult = objectValues(pool).map((array) =>
      invertedReducer(array, reduce)
    )

    const result = intermediateResult.map((r) => finalize(null as any, r))

    const invertedResult = invertedIntermediateResult.map((r) =>
      finalize(null as any, r)
    )

    t.deepEqual(sortObject(result), sortObject(invertedResult))
  })
})

function objectValues<T>(obj: Record<string, T>): T[] {
  return Object.keys(obj).map((key) => obj[key])
}

// from https://gist.github.com/ninapavlich/1697bcc107052f5b884a794d307845fe
function sortObject(object: any): any {
  if (!object) {
    return object
  }

  const isArray = object instanceof Array
  let sortedObj: Record<string, unknown> = {}
  if (isArray) {
    sortedObj = object.map((item: any) => sortObject(item))
  } else {
    const keys = Object.keys(object)
    // console.log(keys);
    keys.sort(function (key1, key2) {
      ;(key1 = key1.toLowerCase()), (key2 = key2.toLowerCase())
      if (key1 < key2) return -1
      if (key1 > key2) return 1
      return 0
    })

    for (const index in keys) {
      const key = keys[index]
      if (typeof object[key] === "object") {
        sortedObj[key] = sortObject(object[key])
      } else {
        sortedObj[key] = object[key]
      }
    }
  }

  return sortedObj
}

// Fonctions globales temporairement importées depuis fakes.js

function reducer(array: any[], reduce: any): any {
  if (array.length == 1) {
    return array[0]
  } else {
    const newVal = reduce(array[0].key, [array[0].value, array[1].value])
    return reducer([newVal].concat(array.slice(2, array.length)), reduce)
  }
}

function invertedReducer(array: any[], reduce: any): any {
  if (array.length == 1) {
    return array[0]
  } else {
    const newVal = reduce(array[0].key, [
      array[array.length - 1].value,
      array[array.length - 2].value,
    ])
    return reducer([newVal].concat(array.slice(0, array.length - 2)), reduce)
  }
}
