"use strict";

var test_cases = [
  {
    ////////////////////////////////////////////////////////
    test_case_name : "Exemple1: complete type deletion",
    completeTypes :  { "1902": ["apconso"] },
    batchKey : "1902",
    types : "",
    batches : ["1901", "1902"],
    reduce_key : "123",
    reduce_values : [
      {"batch":
        {
          "1901":
          {
            "apconso": {
              "a":{"bonjour":3, "aurevoir":4},
              "b":{"bonjour":5, "aurevoir":6}
            }
          }
        },
        "scope":"etablissement"
      },
      {"batch":
        {
          "1902":
          {
            "apconso": {
              "a":{"bonjour":3, "aurevoir":4},
              "c":{"bonjour":7, "aurevoir":8}
            }
          }
        },
        "scope":"etablissement"
      }],
    expected: {
      "key": "123",
      "scope": "etablissement",
      "batch": {
        "1901": {
          "apconso": {
            "a": {
              "bonjour": 3,
              "aurevoir": 4
            },
            "b": {
              "bonjour": 5,
              "aurevoir": 6
            }
          }
        },
        "1902": {
          "apconso": {
            "c": {
              "bonjour": 7,
              "aurevoir": 8
            }
          },
          "compact": {
            "delete": {
              "apconso": [
                "b"
              ]
            }
          }
        }
      }
    }
  },
  {
    ////////////////////////////////////////////////////////
    test_case_name : "Exemple2: order independence",
    completeTypes :  { "1902": ["apconso"] },
    batchKey : "1902",
    types : "",
    batches : ["1901", "1902"],
    reduce_key : "123",
    reduce_values : [
      {"batch":
        {
          "1902":
          {
            "apconso": {
              "a":{"bonjour":3, "aurevoir":4},
              "c":{"bonjour":7, "aurevoir":8}
            }
          }
        },
        "scope":"etablissement"
      },
      {"batch":
        {
          "1901":
          {
            "apconso": {
              "a":{"bonjour":3, "aurevoir":4},
              "b":{"bonjour":5, "aurevoir":6}
            }
          }
        },
        "scope":"etablissement"
      }]
    ,
    expected: {
      "key": "123",
      "scope": "etablissement",
      "batch": {
        "1901": {
          "apconso": {
            "a": {
              "bonjour": 3,
              "aurevoir": 4
            },
            "b": {
              "bonjour": 5,
              "aurevoir": 6
            }
          }
        },
        "1902": {
          "apconso": {
            "c": {
              "bonjour": 7,
              "aurevoir": 8
            }
          },
          "compact": {
            "delete": {
              "apconso": [
                "b"
              ]
            }
          }
        }
      }
    }
  },
  {
    ////////////////////////////////////////////////////////
    test_case_name : "Exemple3: batch insertion between preexisting",
    completeTypes :  {
      "1901": ["apconso"],
      "1902": ["apconso"]
    },
    batchKey : "1901",
    types : "",
    batches : ["1812", "1901", "1902"],
    reduce_key : "123",
    reduce_values : [
      {"batch":
        {"1812":
          {
            "apconso": {
              "deleteme":{"bonjour":1, "aurevoir":2}
            }
          },
          "1902":
          {
            "apconso": {
              "a":{"bonjour":3, "aurevoir":4},
              "c":{"bonjour":7, "aurevoir":8}
            },
            "compact": {
              "delete": {
                "apconso":["deleteme"]
              }
            }
          }
        },
        "scope":"etablissement"
      },
      {"batch":
        { "1901":
          {
            "apconso": {
              "a":{"bonjour":3, "aurevoir":4},
              "b":{"bonjour":5, "aurevoir":6}
            }
          }
        },
        "scope":"etablissement"
      }]
    ,
    expected: {
      "key": "123",
      "scope": "etablissement",
      "batch": {
        "1812": {
          "apconso": {
            "deleteme": {
              "bonjour": 1,
              "aurevoir": 2
            }
          }
        },
        "1901": {
          "apconso": {
            "a": {
              "bonjour": 3,
              "aurevoir": 4
            },
            "b": {
              "bonjour": 5,
              "aurevoir": 6
            }
          },
          "compact": {
            "delete": {
              "apconso": [
                "deleteme"
              ]
            }
          }
        },
        "1902": {
          "apconso": {
            "c": {
              "bonjour": 7,
              "aurevoir": 8
            }
          },
          "compact": {
            "delete": {
              "apconso": [
                "b"
              ]
            }
          }
        }
      }
    }

  },
  {
    ////////////////////////////////////////////////////////
    test_case_name : "Exemple4: added after removed same key",
    completeTypes :  { "1901": ["apconso"] },
    batchKey : "1901",
    types : "",
    batches : ["1812", "1901"],
    reduce_key : "123",
    reduce_values : [
      {"batch":
        {"1812":
          {
            "apconso": {
              "deleteme":{"bonjour":1, "aurevoir":2}
            }
          },
          "1901":
          {
            "compact": {
              "delete": {
                "apconso":["deleteme"]
              }
            }
          }
        },
        "scope":"etablissement"
      },
      {"batch":
        { "1901":
          {
            "apconso": {
              "deleteme":{"bonjour":1, "aurevoir":2}
            }
          }
        },
        "scope":"etablissement"
      }]
    ,
    expected: {
      "key": "123",
      "scope": "etablissement",
      "batch": {
        "1812": {
          "apconso": {
            "deleteme": {
              "bonjour": 1,
              "aurevoir": 2
            }
          }
        }
      }
    }
  },
  {
    ////////////////////////////////////////////////////////
    test_case_name : "Exemple5: deletion without complete types",
    completeTypes :  { "1901": ["apconso"] },
    batchKey : "1901",
    types : "",
    batches : ["1812", "1901"],
    reduce_key : "123",
    reduce_values : [
      {"batch":
        {"1812":
          {
            "apconso": {
              "deleteme":{"bonjour":1, "aurevoir":2}
            }
          },
          "1901":
          {
            "compact": {
              "delete": {
                "apconso":["deleteme"]
              }
            }
          }
        },
        "scope":"etablissement"
      },
      {"batch":
        { "1901":
          {
            "apconso": {
              "deleteme":{"bonjour":1, "aurevoir":2}
            }
          }
        },
        "scope":"etablissement"
      }]
    ,
    expected: {
      "key": "123",
      "scope": "etablissement",
      "batch": {
        "1812": {
          "apconso": {
            "deleteme": {
              "bonjour": 1,
              "aurevoir": 2
            }
          }
        }
      }
    }
  },
  {
    ////////////////////////////////////////////////////////
    test_case_name : "Exemple6: only one batch",
    completeTypes :  {"1901": []},
    batchKey : "1901",
    types : "",
    batches : ["1901"],
    reduce_key : "123",
    reduce_values : [
      {"batch":
        {"1901":
          {
            "apconso": {
              "uneconso":{"bonjour":1, "aurevoir":2}
            }
          },
        },
        "scope":"etablissement"
      },
      {"batch":
        { "1901":
          {
            "apdemande": {
              "unedemande":{"bonjour":2, "aurevoir":1}
            }
          }
        },
        "scope":"etablissement"
      }]
    ,
    expected: {
      "key": "123",
      "scope": "etablissement",
      "batch": {
        "1901": {
          "apconso": {
            "uneconso": {
              "bonjour": 1,
              "aurevoir": 2
            }
          },
          "apdemande": {
            "unedemande": {
              "bonjour": 2,
              "aurevoir": 1
            }
          }
        }
      }
    }
  }
]

const f = {
  currentState: this.currentState
}

var test_results = test_cases.map(function(tc, id) {
  const globals = {
    f,
    completeTypes: tc.completeTypes,
    batchKey: tc.batchKey,
    types: tc.types,
    batches: tc.batches,
  }
  var actual = reduce.call(globals, tc.reduce_key, tc.reduce_values)
  // print(JSON.stringify(actual, null, 2))
  return(compare(actual, tc.expected))
})

print(test_results.every(t => t))
