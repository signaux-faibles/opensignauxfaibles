import test, { ExecutionContext } from "ava"
import "../globals"
import { add } from "./add"

const testCases = [
  {
    data: { "2015-01-01": { any_value: true } },
    data_to_add: {},
    error_expected: true,
    expected: null,
  },
  {
    data: { "2015-01-01": { any_value: true } },
    data_to_add: { "2015-01-01": {} },
    error_expected: false,
    expected: { "2015-01-01": { any_value: true } },
  },
  {
    data: { "2015-01-01": { any_value: true } },
    data_to_add: { "2015-01-01": { any_value: false } },
    error_expected: false,
    expected: { "2015-01-01": { any_value: false } },
  },
  {
    data: { "2015-01-01": { any_value: true } },
    data_to_add: { "2015-01-01": { other_value: false } },
    error_expected: false,
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
    error_expected: false,
    expected: {
      "2015-01-01": { any_value: true, other_value: false },
      "2015-02-01": { any_value: true, other_value: false },
    },
  },
]

testCases.forEach(({ expected, ...testCase }, number) => {
  test.serial(`add(): case #${number}`, (t: ExecutionContext) => {
    if (testCase.error_expected) {
      const result = add(testCase["data_to_add"], testCase["data"])
      t.is(typeof result, "undefined")
    }
    if (expected) {
      add(testCase["data_to_add"], testCase["data"])
      t.true(compare(testCase["data"], expected))
    }
  })
})

function compare(a: any, b: any) {
  if (Object.keys(a).length !== Object.keys(b).length) {
    return false
  }
  const equal = Object.keys(a).every(function (k) {
    return JSON.stringify(a[k]) === JSON.stringify(b[k])
  })
  return equal
}
