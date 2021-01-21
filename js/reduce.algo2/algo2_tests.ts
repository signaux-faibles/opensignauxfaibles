// Ces tests visent à couvrir toutes les fonctions invoquées par map(),
// reduce() et finalize(), en leur fournissant un jeu de données minimal
// mais incluant tous les types d'entrée inclus dans RawData.
//
// => Penser à ajouter les nouveaux types dans rawEtabData et rawEntrData.

import { ExecutionContext, serial as test } from "ava"
import { map, EntréeMap, CléSortieMap, SortieMap } from "./map"
import { reduce } from "./reduce"
import { finalize, Clé } from "./finalize"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { objects as testCases } from "../test/data/objects"
import { naf as nafValues } from "../test/data/naf"
import { reducer, invertedReducer } from "../test/helpers/reducers"
import { runMongoMap, indexMapResultsByKey } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"
import {
  EntréeBdf,
  EntréeDelai,
  EntréeDebit,
  EntrepriseBatchProps,
} from "../RawDataTypes"

const DAY_IN_MS = 24 * 60 * 60 * 1000

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

const makeInput = (
  dateDebut: Date,
  dateFin: Date,
  delaiOverrides?: Partial<EntréeDelai>
): EntréeMap => {
  const siret = "12345678901234"
  const duréeDelai = 60 // en jours
  return {
    _id: siret,
    value: {
      key: siret,
      scope: "etablissement",
      batch: {
        "1905": {
          cotisation: {},
          debit: {
            hashDette: {
              periode: {
                start: dateDebut,
                end: dateFin,
              },
              numero_ecart_negatif: 1,
              numero_historique: 2,
              numero_compte: "",
              date_traitement: dateDebut,
              debit_suivant: "",
              part_ouvriere: 60,
              part_patronale: 0,
            } as EntréeDebit,
          },
          delai: {
            hashDelai: {
              date_creation: dateDebut,
              date_echeance: new Date(
                dateDebut.getTime() + duréeDelai * DAY_IN_MS
              ),
              duree_delai: duréeDelai,
              montant_echeancier: 100,
              ...delaiOverrides,
            } as EntréeDelai,
          },
        },
      },
    },
  }
}

// Tests

test("l'ordre de traitement des données n'influe pas sur les résultats", (t: ExecutionContext) => {
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
})

test("delai_deviation_remboursement est calculé à partir d'un débit et d'une demande de délai de règlement de cotisations sociales", (t: ExecutionContext) => {
  const dateDebut = new Date("2018-01-01")
  const dateFin = new Date("2018-02-01")
  initGlobalParams(dateDebut, dateFin)
  const input = makeInput(dateDebut, dateFin)

  const [res] = runMongoMap<EntréeMap, CléSortieMap, SortieMap>(map, [input])
  const finalCompanyData = Object.values(res?.value ?? {})[0]
  t.is(typeof finalCompanyData?.delai_deviation_remboursement, "number")
  t.is(finalCompanyData?.delai_deviation_remboursement, -0.4)
})

test(`algo2.map() retourne les propriétés d'entreprise attendues par l'algo d'apprentissage`, (t: ExecutionContext) => {
  const siren = "012345678"
  const batchKey = "1910"
  const dateDébut = new Date("2015-12-01T00:00:00.000+0000")
  const dateFin = new Date("2016-01-01T00:00:00.000+0000")
  const serie_periode = [dateDébut, dateFin]
  setGlobals({
    actual_batch: batchKey,
    serie_periode,
    includes: { all: true },
  })
  // Cet objet sera a compléter, au fur et à mesure qu'on ajoutera des props dans EntrepriseBatchProps
  const rawEntrData: EntrepriseBatchProps = {
    reporder: {},
    paydex: {
      decembre: { date_valeur: new Date("2015-12-15T00:00Z"), nb_jours: 1 },
      janvier: { date_valeur: new Date("2016-01-15T00:00Z"), nb_jours: 2 },
    },
  }
  // Note: l'entrée bdf est nécéssaire pour map() émette les données
  const entréeBdf = {
    arrete_bilan_bdf: new Date(
      dateDébut.getTime() - 7 * 31 * 24 * 60 * 60 * 1000 // 7 mois avant date_debut
    ),
  } as EntréeBdf
  const mapResults = runMongoMap<EntréeMap, CléSortieMap, SortieMap>(map, [
    {
      _id: siren,
      value: {
        scope: "entreprise",
        key: siren,
        batch: {
          [batchKey]: {
            ...rawEntrData,
            bdf: { entréeBdf },
          },
        },
      },
    },
  ])
  t.snapshot(mapResults)
})

// TODO: écrire un test équivalent pour les établissements
