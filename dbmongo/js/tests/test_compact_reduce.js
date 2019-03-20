/////////////////
/// Exemple 1 ///
/////////////////
// on simule le scope
completeTypes = {
  
  "1902": ["apconso"]
}

batchKey = "1902"

types = ""
batches = ["1901", "1902"]
f = {}
f.currentState = currentState

// on simule le map
var key = "123"
var values = [
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
  }]
  


print("Exemple retourné")
print(JSON.stringify(reduce(key,values), null, 2))

// Exemple simple dans l'autre sens
var values = [
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
/////////////////
/// Exemple 2 ///
/////////////////
// Modification du passé

batchKey = "1901"
completeTypes = {
  "1901": ["apconso"],  
  "1902": ["apconso"]
}

//on simule le map
var values = [
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


print("Exemple avec insertion au milieu")
print(JSON.stringify(reduce(key,values), null, 2))
