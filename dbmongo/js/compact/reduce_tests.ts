import test, { ExecutionContext } from "ava"
import "../globals"
import { reduce } from "./reduce"

// Paramètres globaux utilisés par "compact"
const jsParams = globalThis as any // => all properties of this object will become global. TODO: remove this when merging namespace (https://github.com/signaux-faibles/opensignauxfaibles/pull/40)
/*
declare const completeTypes: Record<BatchKey, DataType[]>
declare const serie_periode: Date[]
declare const batches: BatchKey[]
declare const batchKey: BatchKey // TODO: à renommer, après fusion de https://github.com/signaux-faibles/opensignauxfaibles/pull/90
*/

const AP_CONSO = {
  id_conso: "",
  periode: new Date(),
  heure_consomme: 0,
}

const AP_DEMANDE = {
  id_demande: "",
  periode: { start: new Date(), end: new Date() },
  hta: null,
  motif_recours_se: null,
}

type TestCase = {
  testCaseName: string
  completeTypes: unknown
  batchKey: string
  types: string
  batches: string[]
  reduce_key: string
  reduce_values: CompanyDataValues[]
  expected: unknown
}

const testCases: TestCase[] = [
  {
    ////////////////////////////////////////////////////////
    testCaseName: "Exemple1: complete type deletion",
    completeTypes: { "1902": ["apconso"] },
    batchKey: "1902",
    types: "",
    batches: ["1901", "1902"],
    reduce_key: "123",
    reduce_values: [
      {
        key: "123",
        batch: {
          "1901": {
            apconso: {
              a: AP_CONSO,
              b: AP_CONSO,
            },
          },
        },
        scope: "etablissement",
      },
      {
        key: "123",
        batch: {
          "1902": {
            apconso: {
              a: AP_CONSO,
              c: AP_CONSO,
            },
          },
        },
        scope: "etablissement",
      },
    ],
    expected: {
      key: "123",
      scope: "etablissement",
      batch: {
        "1901": {
          apconso: {
            a: AP_CONSO,
            b: AP_CONSO,
          },
        },
        "1902": {
          apconso: {
            c: AP_CONSO,
          },
          compact: {
            delete: {
              apconso: ["b"],
            },
          },
        },
      },
    },
  },
  {
    ////////////////////////////////////////////////////////
    testCaseName: "Exemple2: order independence",
    completeTypes: { "1902": ["apconso"] },
    batchKey: "1902",
    types: "",
    batches: ["1901", "1902"],
    reduce_key: "123",
    reduce_values: [
      {
        key: "123",
        batch: {
          "1902": {
            apconso: {
              a: AP_CONSO,
              c: AP_CONSO,
            },
          },
        },
        scope: "etablissement",
      },
      {
        key: "123",
        batch: {
          "1901": {
            apconso: {
              a: AP_CONSO,
              b: AP_CONSO,
            },
          },
        },
        scope: "etablissement",
      },
    ],
    expected: {
      key: "123",
      scope: "etablissement",
      batch: {
        "1901": {
          apconso: {
            a: AP_CONSO,
            b: AP_CONSO,
          },
        },
        "1902": {
          apconso: {
            c: AP_CONSO,
          },
          compact: {
            delete: {
              apconso: ["b"],
            },
          },
        },
      },
    },
  },
  {
    ////////////////////////////////////////////////////////
    testCaseName: "Exemple3: batch insertion between preexisting",
    completeTypes: {
      "1901": ["apconso"],
      "1902": ["apconso"],
    },
    batchKey: "1901",
    types: "",
    batches: ["1812", "1901", "1902"],
    reduce_key: "123",
    reduce_values: [
      {
        key: "123",
        batch: {
          "1812": {
            apconso: {
              deleteme: AP_CONSO,
            },
          },
          "1902": {
            apconso: {
              a: AP_CONSO,
              c: AP_CONSO,
            },
            compact: {
              delete: {
                apconso: ["deleteme"],
              },
            },
          },
        },
        scope: "etablissement",
      },
      {
        key: "123",
        batch: {
          "1901": {
            apconso: {
              a: AP_CONSO,
              b: AP_CONSO,
            },
          },
        },
        scope: "etablissement",
      },
    ],
    expected: {
      key: "123",
      scope: "etablissement",
      batch: {
        "1812": {
          apconso: {
            deleteme: AP_CONSO,
          },
        },
        "1901": {
          apconso: {
            a: AP_CONSO,
            b: AP_CONSO,
          },
          compact: {
            delete: {
              apconso: ["deleteme"],
            },
          },
        },
        "1902": {
          apconso: {
            c: AP_CONSO,
          },
          compact: {
            delete: {
              apconso: ["b"],
            },
          },
        },
      },
    },
  },
  {
    ////////////////////////////////////////////////////////
    testCaseName: "Exemple4: added after removed same key",
    completeTypes: { "1901": ["apconso"] },
    batchKey: "1901",
    types: "",
    batches: ["1812", "1901"],
    reduce_key: "123",
    reduce_values: [
      {
        key: "123",
        batch: {
          "1812": {
            apconso: {
              deleteme: AP_CONSO,
            },
          },
          "1901": {
            compact: {
              delete: {
                apconso: ["deleteme"],
              },
            },
          },
        },
        scope: "etablissement",
      },
      {
        key: "123",
        batch: {
          "1901": {
            apconso: {
              deleteme: AP_CONSO,
            },
          },
        },
        scope: "etablissement",
      },
    ],
    expected: {
      key: "123",
      scope: "etablissement",
      batch: {
        "1812": {
          apconso: {
            deleteme: AP_CONSO,
          },
        },
      },
    },
  },
  {
    ////////////////////////////////////////////////////////
    testCaseName: "Exemple5: deletion without complete types",
    completeTypes: { "1901": ["apconso"] },
    batchKey: "1901",
    types: "",
    batches: ["1812", "1901"],
    reduce_key: "123",
    reduce_values: [
      {
        key: "123",
        batch: {
          "1812": {
            apconso: {
              deleteme: AP_CONSO,
            },
          },
          "1901": {
            compact: {
              delete: {
                apconso: ["deleteme"],
              },
            },
          },
        },
        scope: "etablissement",
      },
      {
        key: "123",
        batch: {
          "1901": {
            apconso: {
              deleteme: AP_CONSO,
            },
          },
        },
        scope: "etablissement",
      },
    ],
    expected: {
      key: "123",
      scope: "etablissement",
      batch: {
        "1812": {
          apconso: {
            deleteme: AP_CONSO,
          },
        },
      },
    },
  },
  {
    ////////////////////////////////////////////////////////
    testCaseName: "Exemple6: only one batch",
    completeTypes: { "1901": [] },
    batchKey: "1901",
    types: "",
    batches: ["1901"],
    reduce_key: "123",
    reduce_values: [
      {
        key: "123",
        batch: {
          "1901": {
            apconso: {
              uneconso: AP_CONSO,
            },
          },
        },
        scope: "etablissement",
      },
      {
        key: "123",
        batch: {
          "1901": {
            apdemande: {
              unedemande: AP_DEMANDE,
            },
          },
        },
        scope: "etablissement",
      },
    ],
    expected: {
      key: "123",
      scope: "etablissement",
      batch: {
        "1901": {
          apconso: {
            uneconso: AP_CONSO,
          },
          apdemande: {
            unedemande: AP_DEMANDE,
          },
        },
      },
    },
  },
]

testCases.forEach(({ testCaseName, expected, ...testCase }) => {
  test.serial(`reduce: ${testCaseName}`, (t: ExecutionContext) => {
    jsParams.completeTypes = testCase.completeTypes
    jsParams.batchKey = testCase.batchKey
    jsParams.types = testCase.types
    jsParams.batches = testCase.batches

    const actualResults = reduce(testCase.reduce_key, testCase.reduce_values)
    t.deepEqual(actualResults, expected) // au lieu d'appel à compare()
  })
})
