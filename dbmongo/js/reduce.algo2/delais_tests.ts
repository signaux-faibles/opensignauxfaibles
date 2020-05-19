import test from "ava"
import "../globals"
import { delais, Delai, DelaiMap } from "./delais"

const janvier = new Date("2014-01-01")
const fevrier = new Date("2014-02-01")
const mars = new Date("2014-03-01")

f = {
  generatePeriodSerie: function (/*date_creation: Date, date_echeance: Date*/): Date[] {
    return [janvier, fevrier, mars]
  },
}

test("delais est défini", (t) => {
  t.is(typeof delais, "function")
})

test("la propriété delai représente le nombre de mois complets restants du délai", (t) => {
  const delaiTest: Delai = {
    numero_compte: undefined,
    numero_contentieux: undefined,
    date_creation: new Date("2014-01-03"),
    date_echeance: new Date("2014-03-05"),
    duree_delai: 61, // nombre de jours entre date_creation et date_echeance
    denomination: undefined,
    indic_6m: undefined,
    annee_creation: undefined,
    montant_echeancier: 1000,
    stade: undefined,
    action: undefined,
  }
  const delaiMap: DelaiMap = {
    abc: delaiTest,
  }
  const output_indexed = {}
  output_indexed[fevrier.getTime()] = {}
  output_indexed[mars.getTime()] = {}
  delais({ delai: delaiMap }, output_indexed)
  t.is(output_indexed[fevrier.getTime()].delai, 1) // nombre de mois complets restants
  t.is(output_indexed[mars.getTime()].delai, 0) // moins d'un mois
})

test("un délai en dehors de la période d'intérêt est ignorée", (t) => {
  const delaiTest: Delai = {
    numero_compte: undefined,
    numero_contentieux: undefined,
    date_creation: new Date("2013-01-03"),
    date_echeance: new Date("2013-03-05"),
    duree_delai: 61, // nombre de jours entre date_creation et date_echeance
    denomination: undefined,
    indic_6m: undefined,
    annee_creation: undefined,
    montant_echeancier: 1000,
    stade: undefined,
    action: undefined,
  }
  const delaiMap: DelaiMap = {
    abc: delaiTest,
  }
  const output_indexed = {}
  output_indexed[fevrier.getTime()] = {}
  delais({ delai: delaiMap }, output_indexed)
  t.deepEqual(Object.keys(output_indexed), [fevrier.getTime().toString()])
})

// TODO: ajouter des tests sur les autres propriétés retournées
// TODO: ajouter des tests sur les cas limites => table-based testing
