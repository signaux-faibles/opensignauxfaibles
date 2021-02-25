import test, { ExecutionContext } from "ava"
import { add } from "./add"
import { ParPériode } from "../RawDataTypes"
import { parPériode } from "../test/helpers/parPeriode"

const testCases = [
  {
    name: "the data does not change when adding an empty object",
    data: parPériode({ "2015-01-01": { any_value: true } }),
    additions: parPériode({}),
    expected: parPériode({ "2015-01-01": { any_value: true } }),
  },
  {
    name:
      "the data does not change when adding an empty object for a given period",
    data: parPériode({ "2015-01-01": { any_value: true } }),
    additions: parPériode({ "2015-01-01": {} }),
    expected: parPériode({ "2015-01-01": { any_value: true } }),
  },
  {
    name: "the data changes when overwriting a property for a given period",
    data: parPériode({ "2015-01-01": { any_value: true } }),
    additions: parPériode({ "2015-01-01": { any_value: false } }),
    expected: parPériode({ "2015-01-01": { any_value: false } }),
  },
  {
    name: "properties are merged for any given period",
    data: parPériode({ "2015-01-01": { any_value: true } }),
    additions: parPériode({ "2015-01-01": { other_value: false } }),
    expected: parPériode({
      "2015-01-01": { any_value: true, other_value: false },
    }),
  },
  {
    name: "properties are merged for more than one given period",
    data: parPériode({
      "2015-01-01": { any_value: true },
      "2015-02-01": { any_value: true },
    }),
    additions: parPériode({
      "2015-01-01": { other_value: false },
      "2015-02-01": { other_value: false },
    }),
    expected: parPériode({
      "2015-01-01": { any_value: true, other_value: false },
      "2015-02-01": { any_value: true, other_value: false },
    }),
  },
]

type DataEntry = Record<string, unknown>

testCases.forEach(({ name, data, additions, expected }) => {
  test.serial(`add(): ${name}`, (t: ExecutionContext) => {
    add(additions as ParPériode<DataEntry>, data as ParPériode<DataEntry>)
    t.deepEqual(data, expected)
  })
})
