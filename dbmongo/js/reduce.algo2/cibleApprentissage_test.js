"use strict";

const test_cases = [
  {
    data: {"2015-01-01": {tag_default: true}},
    n_months: 1,
    expected: {"2015-01-01":{"time_til_outcome":0,"outcome":true, "time_til_default":0}}
  },
  {
    data: {
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": {},
      "2015-04-01": {}
    },
    n_months: 1,
    expected: {
      "2015-01-01":{},
      "2015-02-01":{},
      "2015-03-01":{},
      "2015-04-01":{}
    }
  },
  {
    data: {
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": {tag_default: true},
      "2015-04-01": {}
    },
    n_months: 1,
    expected: {
      "2015-01-01":{"time_til_outcome":2, "outcome":false, "time_til_default":2},
      "2015-02-01":{"time_til_outcome":1, "outcome":true, "time_til_default":1},
      "2015-03-01":{"time_til_outcome":0, "outcome":true, "time_til_default":0},
      "2015-04-01":{}
    }
  },
  {
    data: {
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": {tag_failure: true},
      "2015-04-01": {}
    },
    n_months: 1,
    expected: {
      "2015-01-01":{"time_til_outcome":2, "outcome":false, "time_til_failure":2},
      "2015-02-01":{"time_til_outcome":1, "outcome":true, "time_til_failure":1},
      "2015-03-01":{"time_til_outcome":0, "outcome":true, "time_til_failure":0},
      "2015-04-01":{}
    }
  },
  {
    data: {
      "2015-01-01": {},
      "2015-02-01": {},
      "2015-03-01": {tag_failure: true},
      "2015-04-01": {tag_default: true}
    },
    n_months: 1,
    expected: {
      "2015-01-01":{"time_til_outcome":2, "outcome":false, "time_til_default": 3, "time_til_failure": 2},
      "2015-02-01":{"time_til_outcome":1, "outcome":true, "time_til_default": 2, "time_til_failure": 1},
      "2015-03-01":{"time_til_outcome":0, "outcome":true, "time_til_default": 1, "time_til_failure": 0},
      "2015-04-01":{"time_til_outcome":0, "outcome":true, "time_til_default": 0}
    }
  }
]
Object.freeze(test_cases)

// Define f as a namespace that contains all global functions
const f = this // eslint-disable-line @typescript-eslint/no-this-alias

const test_results = test_cases.map(function(tc, id){
  const actual = cibleApprentissage(tc["data"], tc["n_months"])
  const test_passes = compare(actual, tc["expected"])
  return test_passes
})

print(test_results.every(t => t))
