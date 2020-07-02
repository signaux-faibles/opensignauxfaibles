import test from "ava"
import "../globals"
import {
  delais,
  DelaiComputedValues,
  DebitComputedValues,
  ParPériode,
} from "./delais"

const fevrier = new Date("2014-02-01")
const mars = new Date("2014-03-01")

const nbDays = (firstDate: Date, secondDate: Date): number => {
  const oneDay = 24 * 60 * 60 * 1000 // hours*minutes*seconds*milliseconds
  return Math.round(
    Math.abs((firstDate.getTime() - secondDate.getTime()) / oneDay)
  )
}

const makeDelai = (firstDate: Date, secondDate: Date): EntréeDelai => ({
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

// TODO: renommer cette fonction
const testProperty = (
  debits?: DebitComputedValues
): ParPériode<DelaiComputedValues> => {
  const delaiTest = makeDelai(new Date("2014-01-03"), new Date("2014-04-05"))
  const delaiMap: ParPériode<EntréeDelai> = {
    abc: delaiTest,
  }
  const output_indexed: ParPériode<DebitComputedValues> = {}
  output_indexed[fevrier.getTime()] = debits ? makeOutputIndexed(debits) : {}
  output_indexed[mars.getTime()] = debits ? makeOutputIndexed(debits) : {}
  return delais({ delai: delaiMap }, output_indexed)
}

test("la propriété delai_nb_jours_restants représente le nombre de jours restants du délai", (t) => {
  const output_indexed = testProperty()
  t.is(output_indexed[fevrier.getTime()]["delai_nb_jours_restants"], 2)
  t.is(output_indexed[mars.getTime()]["delai_nb_jours_restants"], 1)
})

test("la propriété delai_nb_jours_total représente la durée totale en jours du délai", (t) => {
  const dureeEnJours = nbDays(new Date("2014-01-03"), new Date("2014-04-05"))
  const output_indexed = testProperty()
  t.is(output_indexed[fevrier.getTime()]["delai_nb_jours_total"], dureeEnJours)
  t.is(output_indexed[mars.getTime()]["delai_nb_jours_total"], dureeEnJours)
})

test("la propriété delai_montant_echeancier représente le montant en euros des cotisations sociales couvertes par le délai", (t) => {
  const output_indexed = testProperty()
  t.is(output_indexed[fevrier.getTime()]["delai_montant_echeancier"], 1000)
  t.is(output_indexed[mars.getTime()]["delai_montant_echeancier"], 1000)
})

test(
  "la propriété delai_deviation_remboursement représente:\n" +
    "(dette actuelle - dette hypothétique en cas de remboursement linéaire) / dette initiale\n" +
    "Elle représente la déviation par rapport à un remboursement linéaire de la dette " +
    "en pourcentage de la dette initialement dû",
  (t) => {
    const expectedFebruary = -0.052
    const expectedMarch = 0.273
    const debits = { montant_part_patronale: 600, montant_part_ouvriere: 0 }
    const output_indexed = testProperty(debits)
    const tolerance = 10e-3
    const ratioFebruary =
      output_indexed[fevrier.getTime()]["delai_deviation_remboursement"]
    const ratioMarch =
      output_indexed[mars.getTime()]["delai_deviation_remboursement"]
    t.is(typeof ratioFebruary, "number")
    t.is(typeof ratioMarch, "number")
    if (typeof ratioFebruary === "number") {
      t.true(Math.abs(ratioFebruary - expectedFebruary) < tolerance)
    }
    if (typeof ratioMarch === "number") {
      t.true(Math.abs(ratioMarch - expectedMarch) < tolerance)
    }
  }
)

test("la propriété delai_deviation_remboursement n'est pas créée si la durée du délai est nulle", (t) => {
  // const debits = { montant_part_patronale: 600, montant_part_ouvriere: 0 }
  const delaiTest = makeDelai(new Date("2014-02-03"), new Date("2014-02-03"))
  const delaiMap: ParPériode<EntréeDelai> = {
    abc: delaiTest,
  }
  const input_indexed: ParPériode<DebitComputedValues> = {
    [fevrier.getTime()]: {},
  }
  const output_indexed = delais({ delai: delaiMap }, input_indexed)
  t.is(Object.keys(output_indexed).length, 1)
  const ratioFebruary =
    output_indexed[fevrier.getTime()]["delai_deviation_remboursement"]
  t.is(typeof ratioFebruary, "undefined")
})

test("un délai en dehors de la période d'intérêt est ignorée", (t) => {
  const delaiTest = makeDelai(new Date("2013-01-03"), new Date("2013-03-05"))
  const delaiMap: ParPériode<EntréeDelai> = {
    abc: delaiTest,
  }
  const donnéesParPériode: ParPériode<DebitComputedValues> = {}
  donnéesParPériode[fevrier.getTime()] = makeOutputIndexed()
  const périodesComplétées = delais({ delai: delaiMap }, donnéesParPériode)
  t.deepEqual(périodesComplétées, {})
})

// TODO: ajouter des tests sur les cas limites: denominateurs nuls dans calcul de delai_deviation_remboursement
