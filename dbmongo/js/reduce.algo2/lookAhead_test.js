"use strict";

const test_cases = [
  {
    data: {"2015-01-01": {outcome: true}},
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {"2015-01-01":{"time_til_outcome":0,"outcome":true}}
  },
  {
    data: {"2015-01-01": {outcome: false}},
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {}
  },
  {
    data: {
      "2015-01-01": {outcome: false},
      "2015-02-01": {outcome: false},
      "2015-03-01": {outcome: false}
    },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {}
  },
  {
    data: {
      "2015-01-01": {outcome: true},
      "2015-02-01": {outcome: true},
      "2015-03-01": {outcome: true}
    },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {
      "2015-01-01":{"time_til_outcome":0,"outcome":true},
      "2015-02-01":{"time_til_outcome":0,"outcome":true},
      "2015-03-01":{"time_til_outcome":0,"outcome":true}
    }
  },
  {
    data: {
      "2015-01-01": {outcome: false},
      "2015-02-01": {outcome: false},
      "2015-03-01": {outcome: true}
    },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {
      "2015-01-01":{"time_til_outcome":2,"outcome":false},
      "2015-02-01":{"time_til_outcome":1,"outcome":true},
      "2015-03-01":{"time_til_outcome":0,"outcome":true}
    }
  },
  {
    data: {
      "2015-01-01": {outcome: true},
      "2015-02-01": {outcome: false},
      "2015-03-01": {outcome: false}
    },
    attr_name: "outcome",
    n_months: 1,
    past: true,
    expected: {
      "2015-01-01":{"time_til_outcome":0,"outcome":true},
    }
  },
  {
    data: {
      "2015-01-01": {outcome: true},
      "2015-02-01": {outcome: false},
      "2015-03-01": {outcome: false}
    },
    attr_name: "outcome",
    n_months: 1,
    past: false,
    expected: {
      "2015-01-01":{"time_til_outcome":0,"outcome":true},
      "2015-02-01":{"time_til_outcome":1,"outcome":true},
      "2015-03-01":{"time_til_outcome":2,"outcome":false}
    }
  },
  {
    data: {
      "2015-01-01": {outcome: true},
      "2015-02-01": {},
      "2015-03-01": {}
    },
    attr_name: "outcome",
    n_months: 1,
    past: false,
    expected: {
      "2015-01-01":{"time_til_outcome":0,"outcome":true},
      "2015-02-01":{"time_til_outcome":1,"outcome":true},
      "2015-03-01":{"time_til_outcome":2,"outcome":false}
    }
  }
]
Object.freeze(test_cases)

const test_results = test_cases.map(function(tc, id){
  const actual = lookAhead(tc["data"], tc["attr_name"], tc["n_months"], tc["past"])
  const test_passes = compare(actual, tc["expected"])
  return test_passes
})
print(test_results.every(t=>t))
