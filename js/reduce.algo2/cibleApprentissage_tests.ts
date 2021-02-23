import test, { ExecutionContext } from "ava"
import { cibleApprentissage } from "./cibleApprentissage"
import { parPériode } from "../test/helpers/parPeriode"

type TestCase = {
  // name: string
  data: Parameters<typeof cibleApprentissage>[0]
  n_months: Parameters<typeof cibleApprentissage>[1]
  expected: Partial<ReturnType<typeof cibleApprentissage>>
}

const testCases: TestCase[] = [
  {
    data: parPériode({ "2015-01-01": { tag_default: true } }),
    n_months: 1,
    expected: parPériode({
      "2015-01-01": { time_til_outcome: 0, outcome: true, time_til_default: 0 },
    }),
  },
  {
    data: parPériode({
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": {},
      "2015-04-01": {},
    }),
    n_months: 1,
    expected: parPériode({
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": {},
      "2015-04-01": {},
    }),
  },
  {
    data: parPériode({
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": { tag_default: true },
      "2015-04-01": {},
    }),
    n_months: 1,
    expected: parPériode({
      "2015-01-01": {
        time_til_outcome: 2,
        outcome: false,
        time_til_default: 2,
      },
      "2015-02-01": { time_til_outcome: 1, outcome: true, time_til_default: 1 },
      "2015-03-01": { time_til_outcome: 0, outcome: true, time_til_default: 0 },
      "2015-04-01": { time_til_outcome: -1, outcome: true },
    }),
  },
  {
    data: parPériode({
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": { tag_failure: true },
      "2015-04-01": {},
    }),
    n_months: 1,
    expected: parPériode({
      "2015-01-01": {
        time_til_outcome: 2,
        outcome: false,
        time_til_failure: 2,
      },
      "2015-02-01": { time_til_outcome: 1, outcome: true, time_til_failure: 1 },
      "2015-03-01": { time_til_outcome: 0, outcome: true, time_til_failure: 0 },
      "2015-04-01": { time_til_outcome: -1, outcome: true },
    }),
  },
  {
    data: parPériode({
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": { tag_failure: true },
      "2015-04-01": { tag_default: true },
    }),
    n_months: 1,
    expected: parPériode({
      "2015-01-01": {
        time_til_outcome: 2,
        outcome: false,
        time_til_default: 3,
        time_til_failure: 2,
      },
      "2015-02-01": {
        time_til_outcome: 1,
        outcome: true,
        time_til_default: 2,
        time_til_failure: 1,
      },
      "2015-03-01": {
        time_til_outcome: 0,
        outcome: true,
        time_til_default: 1,
        time_til_failure: 0,
      },
      "2015-04-01": { time_til_outcome: 0, outcome: true, time_til_default: 0 },
    }),
  },
]

testCases.forEach(({ expected, data, n_months }, name) => {
  test.serial(`cibleApprentissage(): ${name}`, (t: ExecutionContext) => {
    const actual = cibleApprentissage(data, n_months)
    t.deepEqual(actual, expected)
  })
})
