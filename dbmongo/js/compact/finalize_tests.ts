import test, { ExecutionContext } from "ava"
import "../globals"
import { finalize } from "./finalize"
import { setGlobals } from "../test/helpers/setGlobals"
import {
  CompanyDataValuesWithCompact,
  BatchValueWithCompact,
} from "./applyPatchesToBatch"
import { EntréeRepOrder } from "../RawDataTypes"

const SIRET = "123"
const DATE_DEBUT = new Date("2014-01-01")
const DATE_FIN = new Date("2015-01-01")

const AP_CONSO = {
  periode: DATE_DEBUT,
  id_conso: "",
  heure_consomme: 0,
}

const EFFECTIF = {
  periode: DATE_DEBUT,
  effectif: 1,
  numero_compte: "456",
}

type TestCase = {
  testCaseName: string
  finalizeObject: CompanyDataValuesWithCompact
  expected: unknown
}

const testCases: TestCase[] = [
  // example 1
  {
    testCaseName: "add random_order",
    finalizeObject: {
      key: SIRET,
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
    expected: {
      key: SIRET,
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
          reporder: {
            [DATE_DEBUT.toString()]: {
              periode: DATE_DEBUT,
              siret: SIRET,
            },
            [DATE_FIN.toString()]: {
              periode: DATE_FIN,
              siret: SIRET,
            },
          },
        },
      },
      index: {
        algo1: false,
        algo2: false,
      },
    },
  },
  // example 2
  {
    testCaseName: "random_order already present",
    finalizeObject: {
      key: SIRET,
      scope: "etablissement",
      batch: {
        "1901": {
          effectif: {
            a: EFFECTIF,
          },
          reporder: {
            [DATE_DEBUT.toString()]: {
              periode: DATE_DEBUT,
            } as EntréeRepOrder,
            [DATE_FIN.toString()]: {
              periode: DATE_FIN,
            } as EntréeRepOrder,
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
      index: {
        algo1: false,
        algo2: false,
      },
    },
    expected: {
      key: SIRET,
      scope: "etablissement",
      batch: {
        "1901": {
          effectif: {
            a: EFFECTIF,
          },
          reporder: {
            [DATE_DEBUT.toString()]: {
              periode: DATE_DEBUT,
            },
            [DATE_FIN.toString()]: {
              periode: DATE_FIN,
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
      index: {
        algo1: true,
        algo2: true,
      },
    },
  },
  // example 3
  {
    testCaseName: "partial random_order",
    finalizeObject: {
      key: SIRET,
      scope: "etablissement",
      batch: {
        "1901": {
          effectif: {
            a: EFFECTIF,
          },
          reporder: {
            [DATE_DEBUT.toString()]: {
              periode: DATE_DEBUT,
            } as EntréeRepOrder,
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
      index: {
        algo1: false,
        algo2: false,
      },
    },
    expected: {
      key: SIRET,
      scope: "etablissement",
      batch: {
        "1901": {
          effectif: {
            a: EFFECTIF,
          },
          reporder: {
            [DATE_DEBUT.toString()]: {
              periode: DATE_DEBUT,
            },
          },
        },
        "1902": {
          apconso: {
            c: AP_CONSO,
          },
          reporder: {
            [DATE_FIN.toString()]: {
              periode: DATE_FIN,
              siret: SIRET,
            },
          },
          compact: {
            delete: {
              apconso: ["b"],
            },
          },
        },
      },
      index: {
        algo1: true,
        algo2: true,
      },
    },
  },
  // example 4
  {
    testCaseName:
      "partial random_order with batch in between (names with text)",
    finalizeObject: {
      key: SIRET,
      scope: "etablissement",
      batch: {
        "1901_1repeatable": {
          reporder: {
            "4124ad3ec7264743785e6a0b107cbc41": {
              siret: "00578004400011",
              periode: DATE_FIN,
              random_order: 0.5391696081492233,
            },
          },
        },
        "1901_2other": {
          other_stuff: {},
        } as BatchValueWithCompact,
        "1902": {
          apconso: {
            c: AP_CONSO,
          },
        },
      },
      index: {
        algo1: false,
        algo2: false,
      },
    },
    expected: {
      key: SIRET,
      scope: "etablissement",
      batch: {
        "1902": {
          apconso: {
            c: AP_CONSO,
          },
          reporder: {
            [DATE_DEBUT.toString()]: {
              periode: DATE_DEBUT,
              siret: SIRET,
            },
          },
        },
        "1901_1repeatable": {
          reporder: {
            "4124ad3ec7264743785e6a0b107cbc41": {
              siret: "00578004400011",
              periode: DATE_FIN,
            },
          },
        },
        "1901_2other": {
          other_stuff: {},
        },
      },
      index: {
        algo1: false,
        algo2: false,
      },
    },
  },
  // example 5
  {
    testCaseName: "Always keep only first reporder",
    finalizeObject: {
      key: SIRET,
      scope: "etablissement",
      batch: {
        "1901": {
          reporder: {
            "4124ad3ec7264743785e6a0b107cbc41": {
              siret: "00578004400011",
              periode: DATE_DEBUT,
              random_order: 0.5391696081492233,
            },
          },
        },
        "1902": {
          apconso: {
            c: AP_CONSO,
          },
          reporder: {
            "4124ad3ec7264743785e6a0b107cbc42": {
              siret: "00578004400011",
              periode: DATE_DEBUT,
              random_order: 0.12,
            },
            "4124ad3ec7264743785e6a0b107cbc43": {
              siret: "00578004400011",
              periode: DATE_FIN,
              random_order: 0.83,
            },
          },
        },
      },
      index: {
        algo1: false,
        algo2: false,
      },
    },
    expected: {
      key: SIRET,
      scope: "etablissement",
      batch: {
        "1901": {
          reporder: {
            "4124ad3ec7264743785e6a0b107cbc41": {
              siret: "00578004400011",
              periode: DATE_DEBUT,
            },
          },
        },
        "1902": {
          apconso: {
            c: AP_CONSO,
          },
          reporder: {
            "4124ad3ec7264743785e6a0b107cbc43": {
              siret: "00578004400011",
              periode: DATE_FIN,
            },
          },
        },
      },
      index: {
        algo1: false,
        algo2: false,
      },
    },
  },
]

const excludeRandomOrder = (obj: unknown): unknown =>
  Object.entries(obj as Record<string, unknown>).reduce(
    (acc, [prop, val]) =>
      prop === "random_order"
        ? acc
        : {
            ...acc,
            [prop]:
              typeof val === "object" && val?.constructor?.name === "Object" // to make sure it's an object, but not an array, nor a Date instance
                ? excludeRandomOrder(val)
                : val,
          },
    {}
  )

testCases.forEach(({ testCaseName, expected, finalizeObject }) => {
  test.serial(`finalize: ${testCaseName}`, (t: ExecutionContext) => {
    setGlobals({ serie_periode: [DATE_DEBUT, DATE_FIN] })
    // exécution du test
    const actual = finalize(SIRET, finalizeObject)
    t.deepEqual(excludeRandomOrder(actual), expected)
  })
})
