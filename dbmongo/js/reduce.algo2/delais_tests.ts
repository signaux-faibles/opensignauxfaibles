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
  numero_compte: "__DUMMY__",
  numero_contentieux: "__DUMMY__",
  date_creation: firstDate,
  date_echeance: secondDate,
  duree_delai: nbDays(firstDate, secondDate),
  denomination: "__DUMMY__",
  indic_6m: "__DUMMY__",
  annee_creation: 2000, // Dummy value
  montant_echeancier: 1000,
  stade: "__DUMMY__",
  action: "__DUMMY__",
})

const testProperty = (
  t: any,
  propertyName: string,
  expectedFebruary: number,
  expectedMarch: number
): void => {
  const delaiTest = makeDelai(new Date("2014-01-03"), new Date("2014-03-05"))
  const delaiMap: DelaiMap = {
    abc: delaiTest,
  }
  const output_indexed = {}
  output_indexed[fevrier.getTime()] = {
    montant_part_patronale: 600, // TODO: n'inclure ces valeurs que dans les tests qui en ont besoin
    montant_part_ouvriere: 0,
  }
  output_indexed[mars.getTime()] = {
    montant_part_patronale: 600,
    montant_part_ouvriere: 0,
  }
  delais({ delai: delaiMap }, output_indexed)
  t.is(output_indexed[fevrier.getTime()][propertyName], expectedFebruary)
  t.is(output_indexed[mars.getTime()][propertyName], expectedMarch)
}

test("la propriété delai représente le nombre de mois complets restants du délai", (t) => {
  testProperty(t, "delai", 1, 0)
})

test("la propriété duree_delai représente la durée totale en jours du délai", (t) => {
  const dureeEnJours = nbDays(new Date("2014-01-03"), new Date("2014-03-05"))
  testProperty(t, "duree_delai", dureeEnJours, dureeEnJours)
})

test("la propriété montant_echeancier représente le montant en euros des cotisations sociales couvertes par le délai", (t) => {
  testProperty(t, "montant_echeancier", 1000, 1000)
})

test.todo(
  "la propriété ratio_dette_delai représente la déviation du remboursement de la dette par rapport à un remboursement linéaire sur la durée du délai"
  /*(t) => {
    // TODO: Inclure la formule dans la documentation de ce test
    // TODO: populer montant_part_patronale et montant_part_ouvriere dans output_indexed
    testProperty(t, "ratio_dette_delai", 1000, 1000)
  }*/
)

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
