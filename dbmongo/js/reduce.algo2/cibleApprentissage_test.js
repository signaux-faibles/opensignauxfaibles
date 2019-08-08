var test_cases = [
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
      "2015-01-01":{"time_til_outcome":2, "outcome":false, "time_til_failure":2, "time_til_default": 3},
      "2015-02-01":{"time_til_outcome":1, "outcome":true, "time_til_failure":1, "time_til_default": 2},
      "2015-03-01":{"time_til_outcome":0, "outcome":true, "time_til_failure":0, "time_til_default": 1},
      "2015-04-01":{"time_til_outcome":0, "outcome":true, "time_til_default":0}
    }
  }
]

var test_results = test_cases.map(function(tc, id){
  var actual = cibleApprentissage(tc["data"], tc["n_months"])

  var test_passes = compare(actual, tc["expected"])
  if (!test_passes){
    console.log("Test fails: " + id)
    console.log("actual:" + JSON.stringify(actual))
    console.log("expected: " + JSON.stringify(tc["expected"]))
  }
  return(test_passes)
})

console.log("cibleApprentissage_test.js",test_results)
