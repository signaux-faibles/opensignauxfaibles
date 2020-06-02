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

const removeRandomOrder = (reporderProp: { [key: string]: RepOrder }): void =>
  Object.keys(reporderProp).forEach((period) => {
    delete reporderProp[period].random_order
  })

const runMongoMap = (mapFct: () => void, keyVal: object): object => {
  const results: { [key: string]: any } = {}
  globalThis.emit = (key: string, value: CompanyDataValuesWithFlags): void => {
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
  [batchKey]: {} as BatchValue, // TODO: rendre optionnelles les props de BatchValues, pour retirer ce cast
}

const importedData = {
  _id: "random123abc",
  value: {
    batch,
    scope,
    key: siret,
  },
}

const expectedMapResults = {
  [siret]: {
    batch,
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
    } as BatchValue, // TODO: rendre optionnelles les props de BatchValues, pour retirer ce cast
  },
  scope,
  index: { algo1: false, algo2: false }, // car il n'y a pas de données justifiant que l'établissement compte 10 employés ou pas
  key: siret,
}

// exécution complète de la chaine "compact"

test.serial(
  `compact.map() groupe les données par siret`,
  (t: ExecutionContext) => {
    const mapResults = runMongoMap(map, importedData)
    t.deepEqual(mapResults, expectedMapResults)
  }
)

test.serial(
  `compact.reduce() agrège les données par entreprise`,
  (t: ExecutionContext) => {
    const reduceValues: CompanyDataValues[] = [expectedMapResults[siret]]
    const reduceResults = reduce(siret, reduceValues)
    t.deepEqual(reduceResults, expectedReduceResults)
  }
)

test.serial(
  `compact.finalize() intègre des clés d'échantillonage pour chaque période`,
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
)
