var test_cases = [
  {
    data: {"2015-01-01": {any_value: true}},
    data_to_add: {},
    error_expected: true,
    expected: {}
  },
  {
    data: {"2015-01-01": {any_value: true}},
    data_to_add: {"2015-01-01":{}},
    error_expected: false,
    expected: {"2015-01-01":{any_value: true}}
  },
  {
    data: {"2015-01-01": {any_value: true}},
    data_to_add: {"2015-01-01":{any_value: false}},
    error_expected: false,
    expected: {"2015-01-01":{any_value: false}}
  },
  {
    data: {"2015-01-01": {any_value: true}},
    data_to_add: {"2015-01-01":{other_value: false}},
    error_expected: false,
    expected: {"2015-01-01":{any_value: true, other_value: false}}
  },
  {
    data: {
      "2015-01-01": {any_value: true},
      "2015-02-01": {any_value: true}
    },
    data_to_add: {
      "2015-01-01":{other_value: false},
      "2015-02-01":{other_value: false},
    },
    error_expected: false,
    expected: {
      "2015-01-01":{any_value: true, other_value: false},
      "2015-02-01":{any_value: true, other_value: false},
    }
  }
]

var test_results = test_cases.map(function(tc, id){
  var test_passes
  if (tc.error_expected) {
    try {
      add(tc["data_to_add"], tc["data"])
    } catch (e) {
      var rightError = (e instanceof EvalError)
    }
    test_passes = rightError
    if (!test_passes){
      console.log("Test fails: error expected. " + id)
    }

  } else {
    add(tc["data_to_add"], tc["data"])
    test_passes = compare(tc["data"], tc["expected"])
    if (!test_passes){
      console.log("Test fails: " + id)
      console.log("actual:" + JSON.stringify(tc["data"]))
      console.log("expected: " + JSON.stringify(tc["expected"]))
    }

  }
  return(test_passes)
})

console.log("add.js",test_results)

