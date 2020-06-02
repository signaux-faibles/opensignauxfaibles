// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import "../globals"
import { map } from "./map.js"
import { flatten } from "./flatten.js"
import { effectifs } from "./effectifs.js"
import { iterable } from "./iterable.js"
import { sirene } from "./sirene.js"
import { cotisations } from "./cotisations.js"
import { debits } from "./debits.js"
import { apconso } from "./apconso.js"
import { delai } from "./delai.js"
import { compte } from "./compte.js"
import { dealWithProcols } from "./dealWithProcols.js"
import { reduce } from "./reduce"
import { finalize } from "./finalize"

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any
global.f = {
  flatten,
  effectifs,
  iterable,
  sirene,
  cotisations,
  debits,
  apconso,
  delai,
  compte,
  dealWithProcols,
}

const ISODate = (date: string): Date => new Date(date)

const runMongoMap = (mapFct: () => void, keyVal: object): object => {
  const results: { [key: string]: any } = {}
  global.emit = (key: any, value: any): void => {
    results[key] = value
  }
  mapFct.call(keyVal)
  return results
}

// test data inspired by test-api.sh
const SIREN_LENGTH = 9
const siret: SiretOrSiren = "01234567891011"
const scope: Scope = "etablissement"
const batchKey = "1910"
const dates = [
  ISODate("2015-12-01T00:00:00.000+0000"),
  ISODate("2016-01-01T00:00:00.000+0000"),
]
global.actual_batch = batchKey // used by map()
global.serie_periode = dates // used by effectifs(), which is called by map()

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

const etablissementKey = scope + "_" + siret

const expectedMapResults = {
  // TODO: structure et valeurs à confirmer
  [etablissementKey]: {
    apconso: [],
    apdemande: [],
    batch: batchKey,
    compte: undefined,
    cotisation: [0, 0],
    debit: [
      { part_ouvriere: 0, part_patronale: 0 },
      { part_ouvriere: 0, part_patronale: 0 },
    ],
    delai: [],
    dernier_effectif: undefined,
    effectif: [],
    idEntreprise: "entreprise_" + siret.substr(0, SIREN_LENGTH),
    key: siret,
    last_procol: {
      etat: "in_bonis",
    },
    procol: undefined,
    sirene: {},
  },
}

const expectedReduceResults = expectedMapResults[etablissementKey] // TODO: à confirmer

// TODO: à comparer avec la sortie de l'API /public, définie dans test-api.sh
const expectedFinalizeResultValue = {
  _id: "etablissement_01234567891011",
  value: {
    key: "01234567891011",
    batch: "1910",
    effectif: [],
    dernier_effectif: undefined,
    sirene: {},
    cotisation: [0, 0],
    debit: [
      {
        part_ouvriere: 0,
        part_patronale: 0,
      },
      {
        part_ouvriere: 0,
        part_patronale: 0,
      },
    ],
    apconso: [],
    apdemande: [],
    delai: [],
    compte: undefined,
    procol: undefined,
    last_procol: {
      etat: "in_bonis",
    },
    idEntreprise: "entreprise_012345678",
  },
}

// exécution complète de la chaine "public"

test.serial(
  `public.map() retourne les propriétés d'établissement présentées sur le frontal`,
  (t: ExecutionContext) => {
    const mapResults = runMongoMap(map, { value: rawData })
    t.deepEqual(mapResults, expectedMapResults)
  }
)

test.serial(
  `public.reduce() retourne les propriétés d'établissement, telles quelles`,
  (t: ExecutionContext) => {
    const reduceValues = [expectedMapResults[etablissementKey]]
    const reduceResults = reduce({ scope }, reduceValues)
    t.deepEqual(reduceResults, expectedReduceResults)
  }
)

test.serial(`public.finalize()`, (t: ExecutionContext) => {
  const finalizeResultValue = finalize({ scope }, expectedReduceResults)
  const finalizeResult = { _id: etablissementKey, value: finalizeResultValue }
  t.deepEqual(finalizeResult, expectedFinalizeResultValue)
})
