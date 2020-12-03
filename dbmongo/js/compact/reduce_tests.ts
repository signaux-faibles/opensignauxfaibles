import test, { ExecutionContext } from "ava"
import { reduce } from "./reduce"
import { setGlobals } from "../test/helpers/setGlobals"
import { DataType, BatchKey, CompanyDataValues } from "../RawDataTypes"
import { CompanyDataValuesWithCompact } from "./applyPatchesToBatch"

const REDUCE_KEY = "123"

const AP_CONSO = {
  id_conso: "",
  periode: new Date(),
  heure_consomme: 0,
}

const AP_DEMANDE = {
  id_demande: "",
  periode: { start: new Date(), end: new Date() },
  hta: 0,
  motif_recours_se: 0,
}

type TestCase = {
  testCaseName: string
  completeTypes: Record<BatchKey, DataType[]>
  fromBatchKey: string
  batches: string[]
  reduce_values: CompanyDataValuesWithCompact[]
  expected: unknown
}

const testCases: TestCase[] = [
  {
    ////////////////////////////////////////////////////////
    testCaseName: "Exemple1: complete type deletion",
    completeTypes: { "1902": ["apconso"] },
    fromBatchKey: "1902",
    batches: ["1901", "1902"],
    reduce_values: [
      {
        key: REDUCE_KEY,
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
        key: REDUCE_KEY,
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
      key: REDUCE_KEY,
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
    fromBatchKey: "1902",
    batches: ["1901", "1902"],
    reduce_values: [
      {
        key: REDUCE_KEY,
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
        key: REDUCE_KEY,
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
      key: REDUCE_KEY,
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
    fromBatchKey: "1901",
    batches: ["1812", "1901", "1902"],
    reduce_values: [
      {
        key: REDUCE_KEY,
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
        key: REDUCE_KEY,
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
      key: REDUCE_KEY,
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
    fromBatchKey: "1901",
    batches: ["1812", "1901"],
    reduce_values: [
      {
        key: REDUCE_KEY,
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
        key: REDUCE_KEY,
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
      key: REDUCE_KEY,
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
    fromBatchKey: "1901",
    batches: ["1812", "1901"],
    reduce_values: [
      {
        key: REDUCE_KEY,
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
        key: REDUCE_KEY,
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
      key: REDUCE_KEY,
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
    fromBatchKey: "1901",
    batches: ["1901"],
    reduce_values: [
      {
        key: REDUCE_KEY,
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
        key: REDUCE_KEY,
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
      key: REDUCE_KEY,
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
    // définition des valeurs de paramètres globaux utilisés par les fonctions de "compact"
    setGlobals({
      completeTypes: testCase.completeTypes,
      fromBatchKey: testCase.fromBatchKey,
      batches: testCase.batches,
    })
    // exécution du test
    const actualResults = reduce(REDUCE_KEY, testCase.reduce_values)
    t.deepEqual(actualResults, expected)
  })
})

test.serial(
  `reduce retourne 2 cotisations depuis deux objets importés couvrant le même batch`,
  (t: ExecutionContext) => {
    // initialisation des paramètres de compact
    const batchId = "1910"
    const hashCotisation = ["hash1", "hash2"]
    setGlobals({
      fromBatchKey: batchId,
      batches: [batchId],
      completeTypes: { [batchId]: [] },
    })
    // execution de compact sur les données importées
    const siret = ""
    const entréeCotisation = {
      periode: {
        start: new Date(),
        end: new Date(),
      },
      du: 64012.0,
    }
    const reduceResults = reduce(siret, [
      {
        scope: "etablissement",
        key: siret,
        batch: {
          [batchId]: {
            cotisation: {
              [hashCotisation[0]]: entréeCotisation,
            },
          },
        },
      },
      {
        scope: "etablissement",
        key: siret,
        batch: {
          [batchId]: {
            cotisation: {
              [hashCotisation[1]]: entréeCotisation,
            },
          },
        },
      },
    ])
    // test sur les données compactées de cotisation
    const cotisations = reduceResults.batch[batchId].cotisation || {}
    t.deepEqual(Object.keys(cotisations), hashCotisation)
  }
)

test.serial(
  `reduce crashe quand une entrée de donnée est indéfinie`,
  (t: ExecutionContext) => {
    // définition des valeurs de paramètres globaux utilisés par les fonctions de "compact"
    const batchKeyWithInvalidData = "2009"
    const fromBatchKey = "2008"
    const oldBatchKey = "2007"
    setGlobals({
      fromBatchKey,
      batches: [fromBatchKey],
      completeTypes: { [fromBatchKey]: [] },
    })
    const key = "01234567891011"
    const previousRawDataValue: CompanyDataValues = {
      key,
      scope: "etablissement",
      batch: { [oldBatchKey]: {} },
    }
    const importedDataValue: CompanyDataValues = {
      key,
      scope: "etablissement",
      batch: { [batchKeyWithInvalidData]: { cotisation: undefined } },
    }
    // exécution du test
    const reducedData = reduce(key, [importedDataValue])
    const error = t.throws(() => {
      /*const mergeOutput =*/ reduce(key, [previousRawDataValue, reducedData])
    })
    t.regex(error.message, /Cannot convert undefined or null to object/)
  }
)

test.serial(
  `https://github.com/signaux-faibles/opensignauxfaibles/issues/248`, // TODO
  (t: ExecutionContext) => {
    // définition des valeurs de paramètres globaux utilisés par les fonctions de "compact"
    const oldBatchKey = "1910_6"
    const fromBatchKey = "2011_0_urssaf"
    const nextBatchKey = "2011_1_sirene"
    setGlobals({
      fromBatchKey,
      batches: [fromBatchKey, nextBatchKey],
      completeTypes: { [fromBatchKey]: [], [nextBatchKey]: [] },
    })
    const key = "000000000"
    const scope = "entreprise"
    const previousRawDataValue: CompanyDataValues = {
      key,
      scope,
      batch: { [oldBatchKey]: {} },
    }
    const importedDataValue: CompanyDataValues = {
      key,
      scope,
      batch: { [fromBatchKey]: {} },
    }
    // exécution du test
    const reducedData = reduce(key, [importedDataValue])
    reduce(key, [previousRawDataValue, reducedData])
    t.pass()
    // TODO: enrichir le test pour vérifier que les données de nextBatchKey sont bien intégrées lors du compactage
  }
)
