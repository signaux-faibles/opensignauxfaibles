// Context: this golden-file-based test runner was designed to prevent
// regressions on the JS functions (common + algo2) used to compute the
// "Features" collection from the "RawData" collection.
//
// It a rewrite of map_test.js in TypeScript and without confidential data.
// Inspired by algo2_tests.ts.
//
// Usage: `$ npm test` (relies on AVA, as defined in package.json)
//
// Update: `$ npx ava reduce.algo2/map_tests.ts --update-snapshots`

import test from "ava"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { map } from "./map"
import { objects as testData } from "../test/data/objects"
import { naf as nafValues } from "../test/data/naf"

// Constantes
const DATE_DEBUT = new Date("2014-01-01")
const DATE_FIN = new Date("2016-01-01")

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
    const id = JSON.stringify(key) //key.siren + key.batch + key.periode.getTime()
    pool[id] = (pool[id] || []).concat([{ key, value }])
  }
  return pool
}

// initialisation des paramètres globaux de reduce.algo2
function initGlobalParams(dateDebut: Date, dateFin: Date) {
  naf = nafValues
  actual_batch = "2002_1"
  date_debut = dateDebut
  date_fin = dateFin
  serie_periode = generatePeriodSerie(dateDebut, dateFin)
  offset_effectif = 2
  includes = { all: true }
}

// Define global parameters that are required by JS functions
const f = { generatePeriodSerie, map }

test("map() retourne les même données que d'habitude", (t) => {
  const results = setupMapCollector()
  testData.forEach(({ _id, value }) => {
    initGlobalParams(DATE_DEBUT, DATE_FIN)
    f.map.call({ _id, value }) // will append to results
    t.log("results:", Object.keys(results))
    return results
  })
  t.snapshot(results)
})
