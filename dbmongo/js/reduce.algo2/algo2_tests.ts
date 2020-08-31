// Converted to TS from _test.js

import test, { ExecutionContext } from "ava"
import { map, EntréeMap, CléSortieMap, SortieMap } from "./map"
import { reduce } from "./reduce"
import { finalize, Clé } from "./finalize"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { objects as testCases } from "../test/data/objects"
import { naf as nafValues } from "../test/data/naf"
import { reducer, invertedReducer } from "../test/helpers/reducers"
import { runMongoMap, indexMapResultsByKey } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"

const DAY_IN_MS = 24 * 60 * 60 * 1000

// initialisation des paramètres globaux de reduce.algo2
function initGlobalParams(
  dateDebut = new Date("2014-01-01"),
  dateFin = new Date("2018-02-01")
) {
  setGlobals({
    naf: nafValues,
    actual_batch: "1905",
    date_debut: dateDebut,
    date_fin: dateFin,
    serie_periode: generatePeriodSerie(dateDebut, dateFin),
    offset_effectif: 2,
    includes: { all: true },
  })
}

const objectValues = <T>(obj: Record<string, T>): T[] =>
  Object.keys(obj).map((key) => obj[key])

// Tests

test("l'ordre de traitement des données n'influe pas sur les résultats", (t: ExecutionContext) => {
  testCases.forEach(({ _id, value }) => {
    initGlobalParams()

    const flatValues = runMongoMap<EntréeMap, CléSortieMap, SortieMap>(map, [
      { _id, value },
    ])
    const groupedValues = indexMapResultsByKey(flatValues)
    const values = objectValues(groupedValues)

    const intermediateResult = values.map((array) => reducer(array, reduce))

    const invertedIntermediateResult = values.map((array) =>
      invertedReducer(array, reduce)
    )

    const result = intermediateResult.map((r) => finalize({} as Clé, r))

    const invertedResult = invertedIntermediateResult.map((r) =>
      finalize({} as Clé, r)
    )

    t.deepEqual(result, invertedResult)
  })
})

test("delai_deviation_remboursement est calculé à partir d'un débit et d'une demande de délai de règlement de cotisations sociales", (t: ExecutionContext) => {
  const dateDebut = new Date("2018-01-01")
  const datePlusUnMois = new Date("2018-02-01")
  initGlobalParams(dateDebut, datePlusUnMois)
  const siret = "12345678901234"
  const duréeDelai = 60 // en jours
  const input = {
    _id: siret,
    value: {
      key: siret,
      scope: "etablissement" as Scope,
      batch: {
        "1905": {
          cotisation: {},
          debit: {
            hashDette: {
              periode: {
                start: dateDebut,
                end: datePlusUnMois,
              },
              numero_ecart_negatif: 1,
              numero_historique: 2,
              numero_compte: "",
              date_traitement: dateDebut,
              debit_suivant: "",
              part_ouvriere: 60,
              part_patronale: 0,
            },
          },
          delai: {
            hashDelai: {
              date_creation: dateDebut,
              date_echeance: new Date(
                dateDebut.getTime() + duréeDelai * DAY_IN_MS
              ),
              duree_delai: duréeDelai,
              montant_echeancier: 100,
            },
          },
        },
      },
    },
  }

  const flatValues = runMongoMap(map, [input]) as {
    _id: unknown
    value: Record<SiretOrSiren, { delai_deviation_remboursement: number }>
  }[]

  const groupedValues = indexMapResultsByKey(flatValues)

  const values = objectValues(groupedValues)

  t.is(values.length, 1)
  t.is(values[0].length, 1)
  t.deepEqual(Object.keys(values[0][0].value), [siret])
  const finalCompanyData = values[0][0].value[siret]
  t.is(typeof finalCompanyData.delai_deviation_remboursement, "number")
  t.is(finalCompanyData.delai_deviation_remboursement, -0.4)
})
