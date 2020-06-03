import test from "ava"
import "../globals"
import {
  delais,
  Delai,
  DelaiMap,
  DelaiComputedValues,
  DebitComputedValues,
  DebitComputedValuesPerPeriod,
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
} = {}): DebitComputedValues => ({
  montant_part_patronale,
  montant_part_ouvriere,
})

type IndexedOutput = {
  [time: string]: DebitComputedValues & DelaiComputedValues
}

const testProperty = (debits?: DebitComputedValues): IndexedOutput => {
  const delaiTest = makeDelai(new Date("2014-01-03"), new Date("2014-04-05"))
  const delaiMap: DelaiMap = {
    abc: delaiTest,
  }
  const output_indexed: DebitComputedValuesPerPeriod = {}
  output_indexed[fevrier.getTime()] = debits ? makeOutputIndexed(debits) : {}
  output_indexed[mars.getTime()] = debits ? makeOutputIndexed(debits) : {}
  return delais({ delai: delaiMap }, output_indexed)
}

test("la propriété delai représente le nombre de mois complets restants du délai", (t) => {
  const output_indexed = testProperty()
  t.is(output_indexed[fevrier.getTime()]["delai"], 2)
  t.is(output_indexed[mars.getTime()]["delai"], 1)
})

test("la propriété duree_delai représente la durée totale en jours du délai", (t) => {
  const dureeEnJours = nbDays(new Date("2014-01-03"), new Date("2014-04-05"))
  const output_indexed = testProperty()
  t.is(output_indexed[fevrier.getTime()]["duree_delai"], dureeEnJours)
  t.is(output_indexed[mars.getTime()]["duree_delai"], dureeEnJours)
})

test("la propriété montant_echeancier représente le montant en euros des cotisations sociales couvertes par le délai", (t) => {
  const output_indexed = testProperty()
  t.is(output_indexed[fevrier.getTime()]["montant_echeancier"], 1000)
  t.is(output_indexed[mars.getTime()]["montant_echeancier"], 1000)
})

test("la propriété ratio_dette_delai représente la déviation du remboursement de la dette par rapport à un remboursement linéaire sur la durée du délai", (t) => {
  // TODO: Inclure la formule dans la documentation de ce test
  const expectedFebruary = -0.052
  const expectedMarch = 0.273
  const debits = { montant_part_patronale: 600, montant_part_ouvriere: 0 }
  const output_indexed = testProperty(debits)
  const tolerance = 10e-3
  const ratioFebruary = output_indexed[fevrier.getTime()]["ratio_dette_delai"]
  const ratioMarch = output_indexed[mars.getTime()]["ratio_dette_delai"]
  t.is(typeof ratioFebruary, "number")
  t.is(typeof ratioMarch, "number")
  if (typeof ratioFebruary === "number") {
    t.true(Math.abs(ratioFebruary - expectedFebruary) < tolerance)
  }
  if (typeof ratioMarch === "number") {
    t.true(Math.abs(ratioMarch - expectedMarch) < tolerance)
  }
})

test("un délai en dehors de la période d'intérêt est ignorée", (t) => {
  const delaiTest = makeDelai(new Date("2013-01-03"), new Date("2013-03-05"))
  const delaiMap: DelaiMap = {
    abc: delaiTest,
  }
  const output_indexed: DebitComputedValuesPerPeriod = {}
  output_indexed[fevrier.getTime()] = makeOutputIndexed()
  delais({ delai: delaiMap }, output_indexed)
  t.deepEqual(Object.keys(output_indexed), [fevrier.getTime().toString()])
})

// TODO: ajouter des tests sur les autres propriétés retournées
// TODO: ajouter des tests sur les cas limites => table-based testing
