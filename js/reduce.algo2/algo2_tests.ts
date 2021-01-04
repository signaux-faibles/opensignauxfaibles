// Ces tests visent à couvrir toutes les fonctions invoquées par map(),
// reduce() et finalize(), en leur fournissant un jeu de données minimal
// mais incluant tous les types d'entrée inclus dans RawData.
//
// => Penser à ajouter les nouveaux types dans rawEtabData et rawEntrData.

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
import { EntrepriseBatchProps, EntrepriseDataValues } from "../RawDataTypes"

// initialisation des paramètres globaux de reduce.algo2
const initGlobalParams = (dateDebut: Date, dateFin: Date) =>
  setGlobals({
    naf: nafValues,
    actual_batch: "1905",
    date_debut: dateDebut,
    date_fin: dateFin,
    serie_periode: generatePeriodSerie(dateDebut, dateFin),
    offset_effectif: 2,
    includes: { all: true },
  })

// Tests

test.serial(
  `algo2.map() retourne toutes les propriétés d'entreprise attendues par l'algo d'apprentissage`,
  (t: ExecutionContext) => {
    const siren = "012345678"
    const batchKey = "1910"
    setGlobals({
      actual_batch: batchKey,
      serie_periode: [
        new Date("2015-12-01T00:00:00.000+0000"),
        new Date("2016-01-01T00:00:00.000+0000"),
      ],
      includes: { all: true },
    })
    const entréePaydex = { date_valeur: new Date(), nb_jours: 1 }
    const rawEntrData: EntrepriseBatchProps = {
      reporder: {},
      paydex: { somehash: entréePaydex },
    }
    const rawData: EntrepriseDataValues = {
      scope: "entreprise",
      key: siren,
      batch: { [batchKey]: rawEntrData },
    }
    const expectedMapResults = {
      ["entreprise_" + siren]: {
        key: siren,
        batch: batchKey,
        paydex: [entréePaydex],
      },
    }
    const mapResults: Record<string, unknown> = {}
    runMongoMap(map, [{ _id: siren, value: rawData }]).map(
      ({ _id, value }) => (mapResults[_id as string] = value)
    )
    t.deepEqual(mapResults, expectedMapResults)
  }
)

// TODO: écrire un test équivalent pour les établissements

test.serial(
  "l'ordre de traitement des données n'influe pas sur les résultats",
  (t: ExecutionContext) => {
    testCases.forEach(({ _id, value }) => {
      initGlobalParams(new Date("2014-01-01"), new Date("2018-02-01"))

      const values = Object.values(
        indexMapResultsByKey(
          runMongoMap<EntréeMap, CléSortieMap, SortieMap>(map, [{ _id, value }])
        )
      )

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
  }
)
