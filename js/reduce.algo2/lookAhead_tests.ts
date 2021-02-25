import test, { ExecutionContext } from "ava"
import { lookAhead } from "./lookAhead"
import { parPériode } from "../test/helpers/parPeriode"
import { ParPériode } from "../RawDataTypes"

type TestCase = {
  name: string
  data: Parameters<typeof lookAhead>[0]
  attr_name: Parameters<typeof lookAhead>[1]
  n_months: Parameters<typeof lookAhead>[2]
  past: Parameters<typeof lookAhead>[3]
  expected: ReturnType<typeof lookAhead>
}

const testCases: Array<TestCase> = [
  {
    name:
      "la période indique une échéance immédiate, si celle-ci est deja atteinte dans l'unique période fournie",
    data: parPériode({ "2015-01-01": { outcome: true } }),
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: parPériode({
      "2015-01-01": { time_til_outcome: 0, outcome: true },
    }),
  },
  {
    name:
      "aucune période n'est retournée si l'échéance n'est pas atteinte dans la période fournie",
    data: parPériode({ "2015-01-01": { outcome: false } }),
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: parPériode({}),
  },
  {
    name:
      "aucune période n'est retournée si l'échéance n'est jamais atteinte dans les données fournies",
    data: parPériode({
      "2015-01-01": { outcome: false },
      "2015-02-01": { outcome: false },
      "2015-03-01": { outcome: false },
    }),
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: parPériode({}),
  },
  {
    name:
      "l'échéance est marquée comme atteinte immédiatement pour chaque période, si outcome est vrai pour chaque période fournie",
    data: parPériode({
      "2015-01-01": { outcome: true },
      "2015-02-01": { outcome: true },
      "2015-03-01": { outcome: true },
    }),
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: parPériode({
      "2015-01-01": { time_til_outcome: 0, outcome: true },
      "2015-02-01": { time_til_outcome: 0, outcome: true },
      "2015-03-01": { time_til_outcome: 0, outcome: true },
    }),
  },
  {
    name:
      "compte à rebours jusqu'à l'échéance déclarée en dernière période fournie",
    data: parPériode({
      "2015-01-01": { outcome: false },
      "2015-02-01": { outcome: false },
      "2015-03-01": { outcome: true },
    }),
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: parPériode({
      "2015-01-01": { time_til_outcome: 2, outcome: false },
      "2015-02-01": { time_til_outcome: 1, outcome: true },
      "2015-03-01": { time_til_outcome: 0, outcome: true },
    }),
  },
  {
    name: "les périodes suivant l'échéance ne sont pas retournées",
    data: parPériode({
      "2015-01-01": { outcome: true },
      "2015-02-01": { outcome: false },
      "2015-03-01": { outcome: false },
    }),
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: parPériode({
      "2015-01-01": { time_til_outcome: 0, outcome: true },
    }),
  },
  {
    name:
      "si l'échéance se répercute dans le futur, les périodes suivantes sont retournées",
    data: parPériode({
      "2015-01-01": { outcome: true },
      "2015-02-01": { outcome: false },
      "2015-03-01": { outcome: false },
    }),
    attr_name: "outcome",
    n_months: 1,
    past: false,
    expected: parPériode({
      "2015-01-01": { time_til_outcome: 0, outcome: true },
      "2015-02-01": { time_til_outcome: 1, outcome: true },
      "2015-03-01": { time_til_outcome: 2, outcome: false },
    }),
  },
  {
    name:
      "si l'échéance se répercute dans le futur, les périodes suivantes sont retournées, même si aucune donnée n'est fournie dans les périodes suivant l'échéance",
    data: parPériode({
      "2015-01-01": { outcome: true },
      "2015-02-01": {},
      "2015-03-01": {},
    }) as TestCase["data"],
    attr_name: "outcome",
    n_months: 1,
    past: false,
    expected: parPériode({
      "2015-01-01": { time_til_outcome: 0, outcome: true },
      "2015-02-01": { time_til_outcome: 1, outcome: true },
      "2015-03-01": { time_til_outcome: 2, outcome: false },
    }),
  },
]

testCases.forEach(({ name, expected, ...tc }) => {
  test.serial(`lookAhead(): ${name}`, (t: ExecutionContext) => {
    const actual = lookAhead(
      tc["data"],
      tc["attr_name"],
      tc["n_months"],
      tc["past"]
    )
    const sortedActual = new ParPériode([...actual.entries()].sort())
    t.deepEqual(sortedActual, expected)
  })
})
