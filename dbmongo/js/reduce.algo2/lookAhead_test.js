
var test_cases = [
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

var test_results = test_cases.map(function(tc, id){
  var actual = lookAhead(tc["data"], tc["attr_name"], tc["n_months"], tc["past"])

  var test_passes = compare(actual, tc["expected"])
  if (!test_passes){
    console.log("Test fails: " + id)
    console.log("actual:" + JSON.stringify(actual))
    console.log("expected: " + JSON.stringify(tc["expected"]))
  }
  return(test_passes)
})
console.log("lookAhead_test.js", test_results)
