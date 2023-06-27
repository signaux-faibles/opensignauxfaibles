// Ces tests visent à couvrir toutes les fonctions invoquées par map(),
// reduce() et finalize(), en leur fournissant un jeu de données minimal
// mais incluant tous les types d'entrée inclus dans RawData.
//
// => Penser à ajouter les nouveaux types dans rawEtabData et rawEntrData.

import test from "ava" // Tests à exécuter de manière séquentielle, pour éviter effets de bords liés aux variables globales, dont fonction emit() initialisée par runMongoMap()
import { map, EntréeMap, CléSortieMap, SortieMap } from "./map"
import { reduce } from "./reduce"
import { finalize, Clé } from "./finalize"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { objects as testCases } from "../test/data/objects"
import { naf as nafValues } from "../test/data/naf"
import { reducer, invertedReducer } from "../test/helpers/reducers"
import { runMongoMap, indexMapResultsByKey } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"
import { EntréeBdf, EntréeDebit, EntréeDelai } from "../GeneratedTypes"
import { EntrepriseBatchProps, EtablissementBatchProps } from "../RawDataTypes"

const DAY_IN_MS = 24 * 60 * 60 * 1000

const batchKey = "1905"

// initialisation des paramètres globaux de reduce.algo2
const initGlobalParams = (dateDebut: Date, dateFin: Date) =>
  setGlobals({
    naf: nafValues,
    actual_batch: batchKey,
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
        [batchKey]: {
          cotisation: {},
          debit: {
            hashDette: {
              periode: {
                start: dateDebut,
                end: dateFin,
              },
              numero_ecart_negatif: "1",
              numero_historique: 2,
              numero_compte: "",
              date_traitement: dateDebut,
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
test.serial(
  "l'ordre de traitement des données n'influe pas sur les résultats",
  (t: test) => {
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

test.serial(
  "delai_deviation_remboursement est calculé à partir d'un débit et d'une demande de délai de règlement de cotisations sociales",
  (t: test) => {
    const dateDebut = new Date("2018-01-01")
    const dateFin = new Date("2018-02-01")
    initGlobalParams(dateDebut, dateFin)
    const input = makeInput(dateDebut, dateFin)

    const [res] = runMongoMap<EntréeMap, CléSortieMap, SortieMap>(map, [input])
    const finalCompanyData = Object.values(res?.value ?? {})[0]
    t.is(typeof finalCompanyData?.delai_deviation_remboursement, "number")
    t.is(finalCompanyData?.delai_deviation_remboursement, -0.4)
  }
)

test.serial(
  `algo2.map() retourne les propriétés d'entreprise attendues par l'algo d'apprentissage`,
  (t: test) => {
    const siren = "012345678"
    const dateDébut = new Date("2015-12-01T00:00:00.000+0000")
    const dateFin = new Date("2016-02-01T00:00:00.000+0000")
    initGlobalParams(dateDébut, dateFin)

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
  }
)

test.serial(
  `algo2.map() retourne les propriétés d'établissements attendues par l'algo d'apprentissage`,
  (t: test) => {
    const siret = "012345678901234"
    const dateDébut = new Date("2015-12-01T00:00:00.000+0000")
    const dateFin = new Date("2016-02-01T00:00:00.000+0000")
    initGlobalParams(dateDébut, dateFin)

    // Cet objet sera a compléter, au fur et à mesure qu'on ajoutera des props dans EtablissementBatchProps
    const rawEtabData: EtablissementBatchProps = {
      reporder: {},
      apdemande: {
        b88032c02d1724d92279a47599b112dd: {
          id_demande: "S044130482",
          effectif_entreprise: 10,
          effectif: 10,
          date_statut: dateDébut,
          periode: {
            start: dateDébut,
            end: dateFin,
          },
          hta: 1500,
          mta: 10000,
          effectif_autorise: 30,
          motif_recours_se: 2,
          heure_consommee: 0,
          montant_consommee: 0,
          effectif_consomme: 0,
          perimetre: 5,
        },
      },
      apconso: {
        e3af9bbf7d37e62fbbf88efdae464746: {
          id_conso: "S044130482",
          heure_consomme: 326.0,
          montant: 2530.13,
          effectif: 20,
          periode: dateDébut,
        },
      },
    }
    const mapResults = runMongoMap<EntréeMap, CléSortieMap, SortieMap>(map, [
      {
        _id: siret,
        value: {
          scope: "etablissement",
          key: siret,
          batch: {
            [batchKey]: {
              ...rawEtabData,
            },
          },
        },
      },
    ])
    t.snapshot(mapResults)
  }
)
