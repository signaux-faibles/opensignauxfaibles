import "../globals"
import test, { ExecutionContext } from "ava"
import { cotisationsdettes } from "./cotisationsdettes"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

test("La variable cotisation représente les cotisations sociales dues à une période donnée", (t: ExecutionContext) => {
  const date = new Date("2018-01-01")
  const datePlusUnMois = new Date("2018-02-01")

  const v: DonnéesCotisation & DonnéesDebit = {
    cotisation: {
      hash1: {
        periode: { start: date, end: datePlusUnMois },
        du: 100,
      },
    },
    debit: {},
  }

  const actual = cotisationsdettes(v, [date.getTime()])

  const expected = {
    [date.getTime()]: {
      montant_part_ouvriere: 0,
      montant_part_patronale: 0,
      cotisation: 100,
    },
  }

  t.deepEqual(actual, expected)
})

test("Le montant de dette d'une période est rapporté dans les périodes suivantes", (t: ExecutionContext) => {
  const dateDebut = new Date("2018-01-01")
  const periode = generatePeriodSerie(
    dateDebut,
    dateAddMonth(dateDebut, 13)
  ).map((date) => date.getTime())

  const v: DonnéesCotisation & DonnéesDebit = {
    cotisation: {
      hash1: {
        periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
        du: 100,
      },
    },
    debit: {},
  }

  const actual = cotisationsdettes(v, periode)

  const montants = {
    montant_part_ouvriere: 0,
    montant_part_patronale: 0,
  }

  const montantsUnMois = {
    ...montants,
    montant_part_ouvriere_past_1: 0,
    montant_part_patronale_past_1: 0,
  }

  const montantsDeuxMois = {
    ...montantsUnMois,
    montant_part_ouvriere_past_2: 0,
    montant_part_patronale_past_2: 0,
  }

  const montantsTroisMois = {
    ...montantsDeuxMois,
    montant_part_ouvriere_past_3: 0,
    montant_part_patronale_past_3: 0,
  }

  const montantsSixMois = {
    ...montantsTroisMois,
    montant_part_ouvriere_past_6: 0,
    montant_part_patronale_past_6: 0,
  }

  const montantsDouzeMois = {
    ...montantsSixMois,
    montant_part_ouvriere_past_12: 0,
    montant_part_patronale_past_12: 0,
  }

  t.deepEqual(actual[dateAddMonth(dateDebut, 1).getTime()], montantsUnMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 2).getTime()], montantsDeuxMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 3).getTime()], montantsTroisMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 4).getTime()], montantsTroisMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 5).getTime()], montantsTroisMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 6).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 7).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 8).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 9).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 10).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 11).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 12).getTime()], montantsDouzeMois)
})
