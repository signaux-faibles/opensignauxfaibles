import test from "ava"
import { nbDays } from "./nbDays"
import { delais, ChampsEntréeDelai, ChampsDettes, SortieDelais } from "./delais"
import { ParHash } from "../RawDataTypes"
import { ParPériode, makePeriodeMap } from "../common/makePeriodeMap"

const fevrier = new Date("2014-02-01")
const mars = new Date("2014-03-01")

const dummyPeriod = new Date("2010-01-01").getTime()

const makeDelai = (firstDate: Date, secondDate: Date): ChampsEntréeDelai => ({
  date_creation: firstDate,
  date_echeance: secondDate,
  duree_delai: nbDays(firstDate, secondDate),
  montant_echeancier: 1000,
})

const makeDebitParPériode = ({
  montant_part_patronale = 0,
  montant_part_ouvriere = 0,
} = {}): ChampsDettes => ({
  montant_part_patronale,
  montant_part_ouvriere,
})

const runDelais = (debits?: ChampsDettes): ParPériode<SortieDelais> => {
  const delaiTest = makeDelai(new Date("2014-01-03"), new Date("2014-04-05"))
  const delaiMap: ParHash<ChampsEntréeDelai> = {
    [dummyPeriod]: delaiTest,
  }
  const debitParPériode = makePeriodeMap<ChampsDettes>()
  if (debits) {
    debitParPériode.set(fevrier, makeDebitParPériode(debits))
    debitParPériode.set(mars, makeDebitParPériode(debits))
  }
  return delais(delaiMap, debitParPériode, {
    premièreDate: fevrier,
    dernièreDate: mars,
  })
}

test("la propriété delai_nb_jours_restants représente le nombre de jours restants du délai", (t: test) => {
  const outputDelai = runDelais()
  t.is(
    outputDelai.get(fevrier)?.["delai_nb_jours_restants"],
    nbDays(new Date("2014-02-01"), new Date("2014-04-05"))
  )
  t.is(
    outputDelai.get(mars)?.["delai_nb_jours_restants"],
    nbDays(new Date("2014-03-01"), new Date("2014-04-05"))
  )
})

test("la propriété delai_nb_jours_total représente la durée totale en jours du délai", (t: test) => {
  const dureeEnJours = nbDays(new Date("2014-01-03"), new Date("2014-04-05"))
  const outputDelai = runDelais()
  t.is(outputDelai.get(fevrier)?.["delai_nb_jours_total"], dureeEnJours)
  t.is(outputDelai.get(mars)?.["delai_nb_jours_total"], dureeEnJours)
})

test("la propriété delai_montant_echeancier représente le montant en euros des cotisations sociales couvertes par le délai", (t: test) => {
  const outputDelai = runDelais()
  t.is(outputDelai.get(fevrier)?.["delai_montant_echeancier"], 1000)
  t.is(outputDelai.get(mars)?.["delai_montant_echeancier"], 1000)
})

test(
  "la propriété delai_deviation_remboursement représente:\n" +
    "(dette actuelle - dette hypothétique en cas de remboursement linéaire) / dette initiale\n" +
    "Elle représente la déviation par rapport à un remboursement linéaire de la dette " +
    "en pourcentage de la dette initialement dû",
  (t: test) => {
    const expectedFebruary = -0.0848
    const expectedMarch = 0.22
    const debits = { montant_part_patronale: 600, montant_part_ouvriere: 0 }
    const outputDelai = runDelais(debits as ChampsDettes)
    const tolerance = 10e-3
    const ratioFebruary = outputDelai.get(fevrier)
      ?.delai_deviation_remboursement
    const ratioMarch = outputDelai.get(mars)?.delai_deviation_remboursement
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

test("un délai en dehors de la période d'intérêt est ignorée", (t: test) => {
  const delaiTest = makeDelai(new Date("2013-01-03"), new Date("2013-03-05"))
  const delaiMap: ParHash<ChampsEntréeDelai> = {
    [dummyPeriod]: delaiTest,
  }
  const donnéesParPériode = makePeriodeMap<ChampsDettes>()
  donnéesParPériode.set(fevrier, makeDebitParPériode())
  const périodesComplétées = delais(delaiMap, donnéesParPériode, {
    premièreDate: fevrier,
    dernièreDate: mars,
  })
  t.deepEqual(périodesComplétées, makePeriodeMap<SortieDelais>())
})
