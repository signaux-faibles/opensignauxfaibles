import test, { ExecutionContext } from "ava"
import { add } from "./add"

const testCases = [
  {
    name: "the data does not change when adding an empty object",
    data: { "2015-01-01": { any_value: true } },
    additions: {},
    expected: { "2015-01-01": { any_value: true } },
  },
  {
    name:
      "the data does not change when adding an empty object for a given period",
    data: { "2015-01-01": { any_value: true } },
    additions: { "2015-01-01": {} },
    expected: { "2015-01-01": { any_value: true } },
  },
  {
    name: "the data changes when overwriting a property for a given period",
    data: { "2015-01-01": { any_value: true } },
    additions: { "2015-01-01": { any_value: false } },
    expected: { "2015-01-01": { any_value: false } },
  },
  {
    name: "properties are merged for any given period",
    data: { "2015-01-01": { any_value: true } },
    additions: { "2015-01-01": { other_value: false } },
    expected: { "2015-01-01": { any_value: true, other_value: false } },
  },
  {
    name: "properties are merged for more than one given period",
    data: {
      "2015-01-01": { any_value: true },
      "2015-02-01": { any_value: true },
    },
    additions: {
      "2015-01-01": { other_value: false },
      "2015-02-01": { other_value: false },
    },
    expected: {
      "2015-01-01": { any_value: true, other_value: false },
      "2015-02-01": { any_value: true, other_value: false },
    },
  },
]

testCases.forEach(({ name, data, additions, expected }) => {
  test.serial(`add(): ${name}`, (t: ExecutionContext) => {
    add(additions, data)
    t.deepEqual(data, expected)
  })
})
