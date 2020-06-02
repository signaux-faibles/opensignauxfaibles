// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import "../globals"
import { map } from "./map"
import { flatten } from "./flatten.js"
import { outputs } from "./outputs.js"
import { repeatable } from "./repeatable.js"
import { add } from "./add.js"
import { defaillances } from "./defaillances.js"
import { dealWithProcols } from "./dealWithProcols.js"
import { populateNafAndApe } from "./populateNafAndApe.js"
// import { reduce } from "./reduce"
// import { finalize } from "./finalize"

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any
global.f = {
  flatten,
  outputs,
  repeatable,
  add,
  defaillances,
  dealWithProcols,
  populateNafAndApe,
}

const ISODate = (date: string): Date => new Date(date)

const runMongoMap = (mapFct: () => void, keyVal: any): any => {
  const results: { [key: string]: any } = {}
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
global.actual_batch = batchKey // used by map()
global.serie_periode = dates // used by map()
global.includes = { all: true } // used by map()
global.naf = {} // used by map()

// même valeur en entrée que pour ../compact/ava_tests.ts
const rawData = {
  batch: {
    [batchKey]: {
      reporder: dates.reduce(
        (reporder, date) => ({
          ...reporder,
          [date.toString()]: { periode: date, siret },
        }),
        {}
      ),
    } as BatchValue, // TODO: rendre optionnelles les props de BatchValues, pour retirer ce cast
  },
  scope,
  index: { algo1: false, algo2: false }, // car il n'y a pas de données justifiant que l'établissement compte 10 employés ou pas
  key: siret,
}

const expectedMapResults = {}
/*
const expectedReduceResults = {}

// extrait de test-api.golden-master.txt, pour les dates spécifiées plus haut
const expectedFinalizeResultValue = [
  {
    _id: {
      batch: "1910",
      siret: "01234567891011",
      periode: ISODate("2015-12-01T00:00:00Z"),
    },
    value: {
      siret: "01234567891011",
      periode: ISODate("2015-12-01T00:00:00Z"),
      effectif: null,
      etat_proc_collective: "in_bonis",
      interessante_urssaf: true,
      outcome: false,
      cotisation_moy12m: 0,
      nbr_etablissements_connus: 1,
    },
  },
  {
    _id: {
      batch: "1910",
      siret: "01234567891011",
      periode: ISODate("2016-01-01T00:00:00Z"),
    },
    value: {
      siret: "01234567891011",
      periode: ISODate("2016-01-01T00:00:00Z"),
      effectif: null,
      etat_proc_collective: "in_bonis",
      interessante_urssaf: true,
      outcome: false,
      cotisation_moy12m: 0,
      nbr_etablissements_connus: 1,
    },
  },
]
*/

// exécution complète de la chaine "reduce.algo2"

test.serial(`reduce.algo2.map()`, (t: ExecutionContext) => {
  const mapResults = runMongoMap(map, { value: rawData })
  t.deepEqual(mapResults, expectedMapResults)
})

test.todo(
  `reduce.algo2.reduce()` /*, (t: ExecutionContext) => {
  const reduceValues: CompanyDataValues[] = [expectedMapResults[siret]]
  const reduceResults = reduce(siret, reduceValues)
  t.deepEqual(reduceResults, expectedReduceResults)
}*/
)

test.todo(
  `reduce.algo2.finalize()` /*, (t: ExecutionContext) => {
  const finalizeResult = finalize(siret, expectedReduceResults)
  t.deepEqual(finalizeResult, expectedFinalizeResultValue)
}*/
)
