import test, { ExecutionContext } from "ava"
import { nbDays } from "./nbDays"
import { SortieCotisationsDettes } from "./cotisationsdettes"
import { delais, SortieDelais } from "./delais"
import { EntréeDelai, ParPériode } from "../RawDataTypes"

const fevrier = new Date("2014-02-01")
const mars = new Date("2014-03-01")

const makeDelai = (firstDate: Date, secondDate: Date): EntréeDelai => ({
  date_creation: firstDate,
  date_echeance: secondDate,
  duree_delai: nbDays(firstDate, secondDate),
  montant_echeancier: 1000,
})

const makeDebitParPériode = ({
  montant_part_patronale = 0,
  montant_part_ouvriere = 0,
} = {}) =>
  ({
    montant_part_patronale,
    montant_part_ouvriere,
  } as SortieCotisationsDettes)

const runDelais = (
  debits?: SortieCotisationsDettes
): ParPériode<SortieDelais> => {
  const delaiTest = makeDelai(new Date("2014-01-03"), new Date("2014-04-05"))
  const delaiMap: ParPériode<EntréeDelai> = {
    abc: delaiTest,
  }
  const debitParPériode: ParPériode<SortieCotisationsDettes> = {}
  if (debits) {
    debitParPériode[fevrier.getTime()] = makeDebitParPériode(debits)
    debitParPériode[mars.getTime()] = makeDebitParPériode(debits)
  }
  return delais(delaiMap, debitParPériode, {
    premièreDate: fevrier,
    dernièreDate: mars,
  })
}

test("la propriété delai_nb_jours_restants représente le nombre de jours restants du délai", (t: ExecutionContext) => {
  const outputDelai = runDelais()
  t.is(
    outputDelai[fevrier.getTime()]?.["delai_nb_jours_restants"],
    nbDays(new Date("2014-02-01"), new Date("2014-04-05"))
  )
  t.is(
    outputDelai[mars.getTime()]?.["delai_nb_jours_restants"],
    nbDays(new Date("2014-03-01"), new Date("2014-04-05"))
  )
})

test("la propriété delai_nb_jours_total représente la durée totale en jours du délai", (t: ExecutionContext) => {
  const dureeEnJours = nbDays(new Date("2014-01-03"), new Date("2014-04-05"))
  const outputDelai = runDelais()
  t.is(outputDelai[fevrier.getTime()]?.["delai_nb_jours_total"], dureeEnJours)
  t.is(outputDelai[mars.getTime()]?.["delai_nb_jours_total"], dureeEnJours)
})

test("la propriété delai_montant_echeancier représente le montant en euros des cotisations sociales couvertes par le délai", (t: ExecutionContext) => {
  const outputDelai = runDelais()
  t.is(outputDelai[fevrier.getTime()]?.["delai_montant_echeancier"], 1000)
  t.is(outputDelai[mars.getTime()]?.["delai_montant_echeancier"], 1000)
})

test(
  "la propriété delai_deviation_remboursement représente:\n" +
    "(dette actuelle - dette hypothétique en cas de remboursement linéaire) / dette initiale\n" +
    "Elle représente la déviation par rapport à un remboursement linéaire de la dette " +
    "en pourcentage de la dette initialement dû",
  (t: ExecutionContext) => {
    const expectedFebruary = -0.0848
    const expectedMarch = 0.22
    const debits = { montant_part_patronale: 600, montant_part_ouvriere: 0 }
    const outputDelai = runDelais(debits as SortieCotisationsDettes)
    const tolerance = 10e-3
    const ratioFebruary =
      outputDelai[fevrier.getTime()]?.["delai_deviation_remboursement"]
    const ratioMarch =
      outputDelai[mars.getTime()]?.["delai_deviation_remboursement"]
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

test("un délai en dehors de la période d'intérêt est ignorée", (t: ExecutionContext) => {
  const delaiTest = makeDelai(new Date("2013-01-03"), new Date("2013-03-05"))
  const delaiMap: ParPériode<EntréeDelai> = {
    abc: delaiTest,
  }
  const donnéesParPériode: ParPériode<SortieCotisationsDettes> = {}
  donnéesParPériode[fevrier.getTime()] = makeDebitParPériode()
  const périodesComplétées = delais(delaiMap, donnéesParPériode, {
    premièreDate: fevrier,
    dernièreDate: mars,
  })
  t.deepEqual(périodesComplétées, {})
})
