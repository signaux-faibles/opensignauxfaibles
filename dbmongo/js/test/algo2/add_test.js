var test_cases = [
  {
    data: {"2015-01-01": {any_value: true}},
    data_to_add: {},
    error_expected: true,
    expected: null,
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
    let result = add(tc["data_to_add"], tc["data"])
    test_passes = result == undefined

  } else {
    add(tc["data_to_add"], tc["data"])
    test_passes = compare(tc["data"], tc["expected"])
  }
  return(test_passes)
})

print(test_results.every(t => t))

