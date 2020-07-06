import test, { ExecutionContext } from "ava"
import "../globals"
import { add } from "./add"

const testCases = [
  {
    data: { "2015-01-01": { any_value: true } },
    data_to_add: {},
    expected: { "2015-01-01": { any_value: true } },
  },
  {
    data: { "2015-01-01": { any_value: true } },
    data_to_add: { "2015-01-01": {} },
    expected: { "2015-01-01": { any_value: true } },
  },
  {
    data: { "2015-01-01": { any_value: true } },
    data_to_add: { "2015-01-01": { any_value: false } },
    expected: { "2015-01-01": { any_value: false } },
  },
  {
    data: { "2015-01-01": { any_value: true } },
    data_to_add: { "2015-01-01": { other_value: false } },
    expected: { "2015-01-01": { any_value: true, other_value: false } },
  },
  {
    data: {
      "2015-01-01": { any_value: true },
      "2015-02-01": { any_value: true },
    },
    data_to_add: {
      "2015-01-01": { other_value: false },
      "2015-02-01": { other_value: false },
    },
    expected: {
      "2015-01-01": { any_value: true, other_value: false },
      "2015-02-01": { any_value: true, other_value: false },
    },
  },
]

testCases.forEach(({ expected, ...testCase }, number) => {
  test.serial(`add(): case #${number}`, (t: ExecutionContext) => {
    add(testCase["data_to_add"], testCase["data"])
    t.deepEqual(testCase["data"], expected)
  })
})
