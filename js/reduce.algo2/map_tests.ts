// Usage: `$ npm test` (relies on AVA, as defined in package.json)
//
// Update: `$ npx ava reduce.algo2/map_tests.ts --update-snapshots`

import test, { ExecutionContext } from "ava"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { map, EntréeMap, CléSortieMap, SortieMap } from "./map"
import { objects as testData } from "../test/data/objects"
import { naf as nafValues } from "../test/data/naf"
import { runMongoMap, indexMapResultsByKey } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"
import { Scope, EntréeDelai } from "../RawDataTypes"

// initialisation des paramètres globaux de reduce.algo2
function initGlobalParams(dateDebut: Date, dateFin: Date) {
  setGlobals({
    naf: nafValues,
    actual_batch: "2002_1",
    date_debut: dateDebut,
    date_fin: dateFin,
    serie_periode: generatePeriodSerie(dateDebut, dateFin),
    offset_effectif: 2,
    includes: { all: true },
  })
}

test("map() retourne les même données que d'habitude", (t) => {
  const DATE_DEBUT = new Date("2014-01-01")
  const DATE_FIN = new Date("2016-01-01")
  initGlobalParams(DATE_DEBUT, DATE_FIN)
  const results = indexMapResultsByKey(runMongoMap(map, testData))
  t.snapshot(results)
})

test("delai_deviation_remboursement est calculé à partir d'un débit et d'une demande de délai de règlement de cotisations sociales", (t: ExecutionContext) => {
  const DAY_IN_MS = 24 * 60 * 60 * 1000
  const dateDebut = new Date("2018-01-01")
  const dateFin = new Date("2018-02-01")
  initGlobalParams(dateDebut, dateFin)
  const input = (function makeInput(
    dateDebut: Date,
    dateFin: Date,
    delaiOverrides?: Partial<EntréeDelai>
  ) {
    const siret = "12345678901234"
    const duréeDelai = 60 // en jours
    return {
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
                  end: dateFin,
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
                ...delaiOverrides,
              },
            },
          },
        },
      },
    }
  })(dateDebut, dateFin)

  const [res] = runMongoMap<EntréeMap, CléSortieMap, SortieMap>(map, [input])
  const finalCompanyData = Object.values(res?.value ?? {})[0]
  t.is(typeof finalCompanyData?.delai_deviation_remboursement, "number")
  t.is(finalCompanyData?.delai_deviation_remboursement, -0.4)
})
