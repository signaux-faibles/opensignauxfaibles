jsParams = {}
jsParams.serie_periode = [new Date("2014-01-01"), new Date("2015-01-01")]
Object.freeze(jsParams);

f = {}
f.complete_reporder = complete_reporder
Object.freeze(f);

var test_cases =
  [
    {
      // Exemple 1: add random_order
      finalize_object: {
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
      },
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
            },
            "reporder": {
              "Wed Jan 01 2014 01:00:00 GMT+0100 (CET)": {
                "random_order": 0.06991004803786005,
                "periode": "2014-01-01T00:00:00.000Z",
                "siret": "123"
              },
              "Thu Jan 01 2015 01:00:00 GMT+0100 (CET)": {
                "random_order": 0.7252145133591734,
                "periode": "2015-01-01T00:00:00.000Z",
                "siret": "123"
              }
            }
          }
        },
        "index": {
          "algo1": false,
          "algo2": false
        }
      }

    },
    // Exemple 2: random_order already present
    {
      finalize_object: {
        "key": "123",
        "scope": "etablissement",
        "batch": {
          "1901": {
            "effectif": {
              "a": {
                "bonjour": 3,
                "aurevoir": 4
              }
            },
            "reporder": {
              "2014": {
                "periode": new Date("2014-01-01")
              },
              "2015": {
                "periode": new Date("2015-01-01")
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
        },
        "index": {
          "algo1": false,
          "algo2": false
        }
      },
      expected: {
        "key": "123",
        "scope": "etablissement",
        "batch": {
          "1901": {
            "effectif": {
              "a": {
                "bonjour": 3,
                "aurevoir": 4
              }
            },
            "reporder": {
              "2014": {
                "periode": "2014-01-01T00:00:00.000Z"
              },
              "2015": {
                "periode": "2015-01-01T00:00:00.000Z"
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
        },
        "index": {
          "algo1": true,
          "algo2": true
        }
      }

    },
    //exemple3: partial random_order
    {
      finalize_object: {
        "key": "123",
        "scope": "etablissement",
        "batch": {
          "1901": {
            "effectif": {
              "a": {
                "bonjour": 3,
                "aurevoir": 4
              }
            },
            "reporder": {
              "2014": {
                "periode": new Date("2014-01-01")
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
        },
        "index": {
          "algo1": false,
          "algo2": false
        }
      },
      expected: {
        "key": "123",
        "scope": "etablissement",
        "batch": {
          "1901": {
            "effectif": {
              "a": {
                "bonjour": 3,
                "aurevoir": 4
              }
            },
            "reporder": {
              "2014": {
                "periode": "2014-01-01T00:00:00.000Z"
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
            "reporder": {
              "Thu Jan 01 2015 01:00:00 GMT+0100 (CET)": {
                "random_order": 0.3914977130034015,
                "periode": "2015-01-01T00:00:00.000Z",
                "siret": "123"
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
        },
        "index": {
          "algo1": true,
          "algo2": true
        }
      }

    },
    //exemple4: partial random_order with batch in between (names with text)
    {
      finalize_object: {
        "key": "123",
        "scope": "etablissement",
        "batch": {
          "1901_1repeatable" : {
            "reporder" : {
              "4124ad3ec7264743785e6a0b107cbc41" : {
                "siret" : "00578004400011",
                "periode" : new Date("2015-01-01"),
                "random_order" : 0.5391696081492233
              },
            }
          },
          "1901_2other": {
            "other_stuff": {
            }
          },
          "1902": {
            "apconso": {
              "c": {
                "bonjour": 7,
                "aurevoir": 8
              }
            },
          }
        },
        "index": {
          "algo1": false,
          "algo2": false
        }
      },
      expected: {
        "key": "123",
        "scope": "etablissement",
        "batch": {
          "1902": {
            "apconso": {
              "c": {
                "bonjour": 7,
                "aurevoir": 8
              }
            },
            "reporder": {
              "Wed Jan 01 2014 01:00:00 GMT+0100 (CET)": {
                "random_order": 0.3914977130034015,
                "periode": "2014-01-01T00:00:00.000Z",
                "siret": "123"
              }
            },
          },
          "1901_1repeatable": {
            "reporder": {
              "4124ad3ec7264743785e6a0b107cbc41": {
                "siret" : "00578004400011",
                "periode": "2015-01-01T00:00:00.000Z",
                "random_order" : 0.5391696081492233
              }
            }
          },
          "1901_2other": {
            "other_stuff": {
            }
          },
        },
        "index": {
          "algo1": true,
          "algo2": true
        }
      }
    },
    //exemple 5: Always keep only first reporder
    {
      finalize_object: {
        "key": "123",
        "scope": "etablissement",
        "batch": {
          "1901" : {
            "reporder" : {
              "4124ad3ec7264743785e6a0b107cbc41" : {
                "siret" : "00578004400011",
                "periode" : new Date("2014-01-01"),
                "random_order" : 0.5391696081492233
              },
            }
          },
          "1902": {
            "apconso": {
              "c": {
                "bonjour": 7,
                "aurevoir": 8
              }
            },
            "reporder" : {
              "4124ad3ec7264743785e6a0b107cbc42" : {
                "siret" : "00578004400011",
                "periode" : new Date("2014-01-01"),
                "random_order" : 0.12
              },
              "4124ad3ec7264743785e6a0b107cbc43" : {
                "siret" : "00578004400011",
                "periode" : new Date("2015-01-01"),
                "random_order" : 0.83
              },
            }
          }
        },
        "index": {
          "algo1": false,
          "algo2": false
        }
      },
      expected: {
        "key": "123",
        "scope": "etablissement",
        "batch": {
          "1901": {
            "reporder": {
              "4124ad3ec7264743785e6a0b107cbc41": {
                "siret" : "00578004400011",
                "periode": "2014-01-01T00:00:00.000Z",
                "random_order" : 0.5391696081492233
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
            "reporder" : {
              "4124ad3ec7264743785e6a0b107cbc43" : {
                "siret" : "00578004400011",
                "periode" : new Date("2015-01-01"),
                "random_order" : 0.12
              },
            },
          },
        },
        "index": {
          "algo1": true,
          "algo2": true
        }
      }
    }
  ]


var test_results = test_cases.map(function(tc, id) {
  var actual = finalize("123", tc.finalize_object)
  // print(JSON.stringify(actual, null, 2))
  return(compareIgnoreRandom(actual, tc.expected))
})

// print(test_results)
print(test_results.every(t => t))
