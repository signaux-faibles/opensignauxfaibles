// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test.sh.

import test, { ExecutionContext } from "ava"
import { map } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { runMongoMap } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"
import {
  BatchValues,
  CompanyDataValuesWithFlags,
  Scope,
  EntréeRepOrder,
  Siret,
} from "../RawDataTypes"

const removeRandomOrder = (
  reporderProp: Record<string, Partial<EntréeRepOrder>>
): void =>
  Object.keys(reporderProp).forEach((period) => {
    delete reporderProp[period]?.random_order
  })

// test data inspired by test.sh
const siret: Siret = "01234567891011"
const scope: Scope = "etablissement"
const fromBatchKey = "1910"
const dates = [
  new Date("2015-12-01T00:00:00.000+0000"),
  new Date("2016-01-01T00:00:00.000+0000"),
]
const batch: BatchValues = {
  [fromBatchKey]: {},
}

const importedData = {
  _id: "random123abc",
  value: {
    batch,
    scope,
    key: siret,
  },
}

const expectedMapResult = {
  _id: siret,
  value: {
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
    [fromBatchKey]: {
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
  index: { algo2: false }, // car il n'y a pas de données justifiant que l'établissement compte 10 employés ou pas
  key: siret,
}

// exécution complète de la chaine "compact"

test.serial(
  `compact.map() groupe les données par siret`,
  (t: ExecutionContext) => {
    const mapResults = runMongoMap(map, [
      {
        _id: null,
        value: { ...importedData.value } as CompanyDataValuesWithFlags,
      },
    ])
    t.deepEqual(mapResults, [expectedMapResult])
  }
)

test.serial(
  `compact.reduce() agrège les données par entreprise`,
  (t: ExecutionContext) => {
    const reduceValues = [expectedMapResult.value]
    const reduceResults = reduce(siret, reduceValues)
    t.deepEqual(reduceResults, expectedReduceResults)
  }
)

test.serial(
  `compact.finalize() intègre des clés d'échantillonage pour chaque période`,
  (t: ExecutionContext) => {
    setGlobals({ serie_periode: dates }) // used by complete_reporder(), which is called by finalize()
    const finalizeResult = finalize(siret, expectedReduceResults)
    const { reporder } = finalizeResult.batch[fromBatchKey] || {}
    t.is(typeof reporder, "object")
    // reporder contient une propriété par periode
    t.is(Object.keys(reporder || {}).length, dates.length)
    Object.keys(reporder || {}).forEach((periodKey) => {
      t.is(typeof reporder?.[periodKey]?.random_order, "number")
    })
    // vérification de la structure complète, sans les nombres aléatoires
    removeRandomOrder(reporder || {}) // will mutate finalizeResult
    t.deepEqual(finalizeResult, expectedFinalizeResultValue)
  }
)
