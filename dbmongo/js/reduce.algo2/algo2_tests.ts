// Converted to TS from _test.js

import test, { ExecutionContext } from "ava"
import { map } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { objects as testCases } from "../test/data/objects"
import { naf as nafValues } from "../test/data/naf"
import { reducer, invertedReducer } from "../test/helpers/reducers"

// Paramètres globaux utilisés par "reduce.algo2"
declare let emit: unknown // called by map()
declare let naf: NAF
declare let actual_batch: BatchKey
declare let date_debut: Date
declare let date_fin: Date
declare let serie_periode: Date[]
declare let offset_effectif: number
declare let includes: Record<"all", boolean>

// preparation de l'environnement d'exécution de map()
function setupMapCollector() {
  const pool: Record<any, any> = {}
  emit = (key: any, value: any) => {
    const id = key.siren + key.batch + key.periode.getTime()
    pool[id] = (pool[id] || []).concat([{ key, value }])
  }
  return pool
}

// initialisation des paramètres globaux de reduce.algo2
function initGlobalParams(
  dateDebut = new Date("2014-01-01"),
  dateFin = new Date("2018-02-01")
) {
  naf = nafValues
  actual_batch = "1905"
  date_debut = dateDebut
  date_fin = dateFin
  serie_periode = generatePeriodSerie(dateDebut, dateFin)
  offset_effectif = 2
  includes = { all: true }
}

test("l'ordre de traitement des données n'influe pas sur les résultats", (t: ExecutionContext) => {
  testCases.forEach(({ _id, value }) => {
    initGlobalParams()

    const pool = setupMapCollector()
    map.call({ _id, value }) // will populate pool

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

// Helpers

const objectValues = <T>(obj: Record<string, T>): T[] =>
  Object.keys(obj).map((key) => obj[key])

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

test("delai_deviation_remboursement est calculé si un délai de règlement de cotisations sociales a été demandé", (t: ExecutionContext) => {
  const dateDebut = new Date("2018-01-01")
  const datePlusUnMois = new Date("2018-02-01")
  initGlobalParams(dateDebut, datePlusUnMois)

  const input = {
    _id: "012345678901234",
    value: {
      key: "012345678901234",
      scope: "etablissement" as Scope,
      batch: {
        "1905": {
          cotisation: {
            hash0: {
              periode: { start: dateDebut, end: datePlusUnMois },
              du: 100,
            },
          },
          debit: {},
          /*
          dettes: {
            hash1: {
              periode: dateDebut,
              part_ouvriere: 100,
              part_patronale: 0,
            },
          },
          delai: {
            hash2: {
              date_creation: dateDebut,
              date_echeance: new Date(
                dateDebut.getTime() + 2 * 24 * 60 * 60 * 1000
                ),
                duree_delai: 2,
                montant_echeancier: 100,
              },
            },
            */
        },
      },
    } as CompanyDataValues,
  }

  const pool = setupMapCollector()
  map.call(input) // will populate pool

  const values = objectValues(pool)
  t.is(values.length, 1)
  t.deepEqual(values[0], [
    {
      key: {
        batch: "1905",
        periode: dateDebut,
        siren: "012345678",
        type: "other",
      },
      value: {
        "012345678901234": {
          cotisation: 100,
          cotisation_moy12m: 100,
          effectif: null,
          etat_proc_collective: "in_bonis",
          interessante_urssaf: true,
          montant_part_ouvriere: 0,
          montant_part_patronale: 0,
          outcome: false,
          periode: dateDebut,
          ratio_dette: 0,
          ratio_dette_moy12m: 0,
          siret: "012345678901234",
        },
      },
    },
  ])
})
