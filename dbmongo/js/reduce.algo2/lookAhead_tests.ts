import test, { ExecutionContext } from "ava"
import "../globals"
import { lookAhead } from "./lookAhead"

type TestCase = {
  data: Parameters<typeof lookAhead>[0]
  attr_name: Parameters<typeof lookAhead>[1]
  n_months: Parameters<typeof lookAhead>[2]
  past: Parameters<typeof lookAhead>[3]
  expected: unknown
}
const testCases: Array<TestCase> = [
  {
    data: { "2015-01-01": { outcome: true } },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: { "2015-01-01": { time_til_outcome: 0, outcome: true } },
  },
  {
    data: { "2015-01-01": { outcome: false } },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {},
  },
  {
    data: {
      "2015-01-01": { outcome: false },
      "2015-02-01": { outcome: false },
      "2015-03-01": { outcome: false },
    },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {},
  },
  {
    data: {
      "2015-01-01": { outcome: true },
      "2015-02-01": { outcome: true },
      "2015-03-01": { outcome: true },
    },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {
      "2015-01-01": { time_til_outcome: 0, outcome: true },
      "2015-02-01": { time_til_outcome: 0, outcome: true },
      "2015-03-01": { time_til_outcome: 0, outcome: true },
    },
  },
  {
    data: {
      "2015-01-01": { outcome: false },
      "2015-02-01": { outcome: false },
      "2015-03-01": { outcome: true },
    },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {
      "2015-01-01": { time_til_outcome: 2, outcome: false },
      "2015-02-01": { time_til_outcome: 1, outcome: true },
      "2015-03-01": { time_til_outcome: 0, outcome: true },
    },
  },
  {
    data: {
      "2015-01-01": { outcome: true },
      "2015-02-01": { outcome: false },
      "2015-03-01": { outcome: false },
    },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {
      "2015-01-01": { time_til_outcome: 0, outcome: true },
    },
  },
  {
    data: {
      "2015-01-01": { outcome: true },
      "2015-02-01": { outcome: false },
      "2015-03-01": { outcome: false },
    },
    attr_name: "outcome",
    n_months: 1,
    past: false,
    expected: {
      "2015-01-01": { time_til_outcome: 0, outcome: true },
      "2015-02-01": { time_til_outcome: 1, outcome: true },
      "2015-03-01": { time_til_outcome: 2, outcome: false },
    },
  },
  {
    data: {
      "2015-01-01": { outcome: true },
      "2015-02-01": {},
      "2015-03-01": {},
    },
    attr_name: "outcome",
    n_months: 1,
    past: false,
    expected: {
      "2015-01-01": { time_til_outcome: 0, outcome: true },
      "2015-02-01": { time_til_outcome: 1, outcome: true },
      "2015-03-01": { time_til_outcome: 2, outcome: false },
    },
  },
]

testCases.forEach(({ expected, ...tc }, number) => {
  test.serial(`add(): case #${number}`, (t: ExecutionContext) => {
    const actual = lookAhead(
      tc["data"],
      tc["attr_name"],
      tc["n_months"],
      tc["past"]
    )
    t.deepEqual(actual, expected)
  })
})
