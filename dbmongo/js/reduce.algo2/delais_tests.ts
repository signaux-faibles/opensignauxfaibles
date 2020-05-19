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

const nbDays = (firstDate: Date, secondDate: Date): number => {
  const oneDay = 24 * 60 * 60 * 1000 // hours*minutes*seconds*milliseconds
  return Math.round(
    Math.abs((firstDate.getTime() - secondDate.getTime()) / oneDay)
  )
}

const makeDelai = (firstDate: Date, secondDate: Date): Delai => ({
  numero_compte: undefined,
  numero_contentieux: undefined,
  date_creation: firstDate,
  date_echeance: secondDate,
  duree_delai: nbDays(firstDate, secondDate),
  denomination: undefined,
  indic_6m: undefined,
  annee_creation: undefined,
  montant_echeancier: 1000,
  stade: undefined,
  action: undefined,
})

test("delais est défini", (t) => {
  t.is(typeof delais, "function")
})

test("la propriété delai représente le nombre de mois complets restants du délai", (t) => {
  const delaiTest = makeDelai(new Date("2014-01-03"), new Date("2014-03-05"))
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

test("la propriété duree_delai représente la durée totale en jours du délai", (t) => {
  const delaiTest = makeDelai(new Date("2014-01-03"), new Date("2014-03-05"))
  const delaiMap: DelaiMap = {
    abc: delaiTest,
  }
  const output_indexed = {}
  output_indexed[fevrier.getTime()] = {}
  output_indexed[mars.getTime()] = {}
  delais({ delai: delaiMap }, output_indexed)
  t.is(
    output_indexed[fevrier.getTime()].duree_delai,
    nbDays(new Date("2014-01-03"), new Date("2014-03-05"))
  )
  t.is(
    output_indexed[mars.getTime()].duree_delai,
    nbDays(new Date("2014-01-03"), new Date("2014-03-05"))
  )
})

test("un délai en dehors de la période d'intérêt est ignorée", (t) => {
  const delaiTest = makeDelai(new Date("2013-01-03"), new Date("2013-03-05"))
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
