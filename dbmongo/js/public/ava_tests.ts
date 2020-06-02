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
// import { reduce } from "./reduce"
// import { finalize } from "./finalize"

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

const expectedMapResults = {}

// exécution complète de la chaine "public"

test.serial(`public.map()`, (t: ExecutionContext) => {
  const mapResults = runMongoMap(map, { value: rawData })
  t.deepEqual(mapResults, expectedMapResults)
})

test.todo(
  `public.reduce()`
  /*
  (t: ExecutionContext) => {
    const reduceValues: CompanyDataValues[] = [expectedMapResults[siret]]
    const reduceResults = reduce(siret, reduceValues)
    t.deepEqual(reduceResults, expectedReduceResults)
  }
  */
)

test.todo(
  `public.finalize()`
  /*
  (t: ExecutionContext) => {
    const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any
    global.serie_periode = dates // used by complete_reporder(), which is called by finalize()
    const finalizeResult = finalize(siret, expectedReduceResults)
    const { reporder } = finalizeResult.batch[batchKey]
    // reporder contient une propriété par periode
    t.is(Object.keys(reporder).length, dates.length)
    Object.keys(reporder).forEach((periodKey) => {
      t.is(typeof reporder[periodKey].random_order, "number")
    })
    // vérification de la structure complète, sans les nombres aléatoires
    removeRandomOrder(finalizeResult.batch[batchKey].reporder) // will mutate finalizeResult
    t.deepEqual(finalizeResult, expectedFinalizeResultValue)
  }
  */
)
