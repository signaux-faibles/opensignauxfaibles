import test, { ExecutionContext } from "ava"
import "../globals"
import {
  delais,
  Delai,
  DelaiMap,
  DelaiComputedValues,
  IndexedOutputExpectedValues,
  IndexedOutputPartial,
} from "./delais"

const fevrier = new Date("2014-02-01")
const mars = new Date("2014-03-01")

const nbDays = (firstDate: Date, secondDate: Date): number => {
  const oneDay = 24 * 60 * 60 * 1000 // hours*minutes*seconds*milliseconds
  return Math.round(
    Math.abs((firstDate.getTime() - secondDate.getTime()) / oneDay)
  )
}

const makeDelai = (firstDate: Date, secondDate: Date): Delai => ({
  date_creation: firstDate,
  date_echeance: secondDate,
  duree_delai: nbDays(firstDate, secondDate),
  montant_echeancier: 1000,
})

const makeOutputIndexed = ({
  montant_part_patronale = 0,
  montant_part_ouvriere = 0,
} = {}): IndexedOutputExpectedValues => ({
  montant_part_patronale,
  montant_part_ouvriere,
})

const testProperty = (
  t: ExecutionContext,
  propertyName: keyof DelaiComputedValues,
  expectedFebruary: number,
  expectedMarch: number
): IndexedOutputPartial => {
  const delaiTest = makeDelai(new Date("2014-01-03"), new Date("2014-04-05"))
  const delaiMap: DelaiMap = {
    abc: delaiTest,
  }
  const output_indexed: IndexedOutputPartial = {}
  output_indexed[fevrier.getTime()] = makeOutputIndexed({
    montant_part_patronale: 600, // TODO: n'inclure ces valeurs que dans les tests qui en ont besoin
    montant_part_ouvriere: 0,
  })
  output_indexed[mars.getTime()] = makeOutputIndexed({
    montant_part_patronale: 600, // TODO: n'inclure ces valeurs que dans les tests qui en ont besoin
    montant_part_ouvriere: 0,
  })
  delais({ delai: delaiMap }, output_indexed)
  t.is(output_indexed[fevrier.getTime()][propertyName], expectedFebruary)
  t.is(output_indexed[mars.getTime()][propertyName], expectedMarch)
  return output_indexed
}

test("la propriété delai représente le nombre de mois complets restants du délai", (t) => {
  testProperty(t, "delai", 2, 1)
})

test("la propriété duree_delai représente la durée totale en jours du délai", (t) => {
  const dureeEnJours = nbDays(new Date("2014-01-03"), new Date("2014-04-05"))
  testProperty(t, "duree_delai", dureeEnJours, dureeEnJours)
})

test("la propriété montant_echeancier représente le montant en euros des cotisations sociales couvertes par le délai", (t) => {
  testProperty(t, "montant_echeancier", 1000, 1000)
})

test("la propriété ratio_dette_delai représente la déviation du remboursement de la dette par rapport à un remboursement linéaire sur la durée du délai", (t) => {
  // TODO: Inclure la formule dans la documentation de ce test
  const expectedFebruary = -0.05217391304347825
  const expectedMarch = 0.2739130434782609
  const output_indexed = testProperty(t, "ratio_dette_delai", expectedFebruary, expectedMarch)

  t.is(output_indexed[fevrier.getTime()]["ratio_dette_delai"], expectedFebruary)
  t.is(output_indexed[mars.getTime()]["ratio_dette_delai"], expectedMarch)
  // TODO: éviter la comparaison de nombres à virgule flottante
})

test("un délai en dehors de la période d'intérêt est ignorée", (t) => {
  const delaiTest = makeDelai(new Date("2013-01-03"), new Date("2013-03-05"))
  const delaiMap: DelaiMap = {
    abc: delaiTest,
  }
  const output_indexed: IndexedOutputPartial = {}
  output_indexed[fevrier.getTime()] = makeOutputIndexed()
  delais({ delai: delaiMap }, output_indexed)
  t.deepEqual(Object.keys(output_indexed), [fevrier.getTime().toString()])
})

// TODO: ajouter des tests sur les autres propriétés retournées
// TODO: ajouter des tests sur les cas limites => table-based testing
