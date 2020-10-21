// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import { map } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { runMongoMap } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"
import {
  BatchValues,
  CompanyDataValuesWithFlags,
  CompanyDataValues,
  Scope,
  EntréeRepOrder,
  SiretOrSiren,
} from "../RawDataTypes"

const removeRandomOrder = (
  reporderProp: Record<string, Partial<EntréeRepOrder>>
): void =>
  Object.keys(reporderProp).forEach((period) => {
    delete reporderProp[period].random_order
  })

// test data inspired by test-api.sh
const siret: SiretOrSiren = "01234567891011"
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
  index: { algo1: false, algo2: false }, // car il n'y a pas de données justifiant que l'établissement compte 10 employés ou pas
  key: siret,
}

// exécution complète de la chaine "compact"

test.serial(
  `compact.map() groupe les données par siret`,
  (t: ExecutionContext) => {
    const mapResults: Record<string, unknown> = {}
    runMongoMap(map, [
      {
        _id: null,
        value: { ...importedData.value } as CompanyDataValuesWithFlags,
      },
    ]).map(({ _id, value }) => (mapResults[_id as string] = value))
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
    setGlobals({ serie_periode: dates }) // used by complete_reporder(), which is called by finalize()
    const finalizeResult = finalize(siret, expectedReduceResults)
    const { reporder } = finalizeResult.batch[fromBatchKey]
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

test.serial(
  `compact retourne 2 cotisations depuis deux objets importés couvrant le même batch`,
  (t: ExecutionContext) => {
    const siret = ""
    const importedData = [
      {
        _id: "abc",
        value: ({
          scope: "etablissement",
          key: siret,
          batch: {
            "1910": {
              cotisation: {
                f72742994ce361fd830eeee5f43f07fd: {
                  periode: {
                    start: new Date("2014-12-01T00:00:00.000Z"),
                    end: new Date("2015-01-01T00:00:00.000Z"),
                  },
                  du: 64012.0,
                },
              },
            },
          },
        } as CompanyDataValues) as CompanyDataValuesWithFlags,
      },
      {
        _id: "def",
        value: ({
          scope: "etablissement",
          key: siret,
          batch: {
            "1910": {
              cotisation: {
                f72742994ce361fd830eeee5f43f07fe: {
                  periode: {
                    start: new Date("2014-11-01T00:00:00.000Z"),
                    end: new Date("2014-12-01T00:00:00.000Z"),
                  },
                  du: 123.0,
                },
              },
            },
          },
        } as CompanyDataValues) as CompanyDataValuesWithFlags,
      },
    ]
    setGlobals({
      fromBatchKey: "1910",
      batches: ["1910"],
      completeTypes: { "1910": [] },
    })
    const mapResults = runMongoMap(map, importedData).map(({ value }) => value)
    const reduceResults = reduce(siret, mapResults as CompanyDataValues[])
    const finalizeResult = finalize(siret, reduceResults)
    const cotisations = finalizeResult.batch["1910"].cotisation || {}
    t.deepEqual(Object.keys(cotisations), [
      "f72742994ce361fd830eeee5f43f07fe",
      "f72742994ce361fd830eeee5f43f07fd",
    ])
  }
)
