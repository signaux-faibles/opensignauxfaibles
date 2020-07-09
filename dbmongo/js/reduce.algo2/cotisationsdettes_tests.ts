import "../globals"
import test, { ExecutionContext } from "ava"
import { cotisationsdettes, SortieCotisationsDettes } from "./cotisationsdettes"

test(`La variable cotisation représente les cotisations sociales dues à une période donnée`, (t: ExecutionContext) => {
  const date = new Date("2018-01-01")
  const datePlusUnMois = new Date("2018-02-01")
  const periode = { start: date, end: datePlusUnMois }
  const v: DonnéesCotisation & DonnéesDebit = {
    cotisation: {
      hash1: {
        periode,
        du: 100,
      },
    },
    debit: {},
  }
  const actual: Record<number, SortieCotisationsDettes> = cotisationsdettes(v, [
    date.getTime(),
  ])

  const expected = {
    [date.getTime()]: {
      montant_part_ouvriere: 0,
      montant_part_patronale: 0,
      cotisation: 100,
    },
  }

  t.deepEqual(actual, expected)
})
