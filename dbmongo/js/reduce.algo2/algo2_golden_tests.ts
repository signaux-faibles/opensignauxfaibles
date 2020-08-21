/*global globalThis*/

// Combinaison des tests de map_test.js et finalize_test.js.
//
// Tests to prevent regressions on the JS functions (reduce.algo2 + common)
// used to compute the "Features" collection from the "RawData" collection.
//
// To update golden files: `$ npx ava algo2_golden_tests.ts -- --update`
//                      or `$ npm run test:update-golden-files`
//
// These tests require the presence of private files specified in the constants
// below. => Make sure to:
// - run `$ git secret reveal` before running these tests;
// - run `$ git secret hide` (to encrypt changes) after updating.

import test, { ExecutionContext as ExecCtx } from "ava"
import * as fs from "fs"
import * as util from "util"
import { naf } from "../test/data/naf"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { map } from "./map"
import { finalize, Clé, EntrepriseEnEntrée } from "./finalize"
import { reduce } from "./reduce"
import { TestDataItem } from "../test/data/objects"
import {
  runMongoMap,
  runMongoReduce,
  parseMongoObject,
  serializeAsMongoObject,
} from "../test/helpers/mongodb"

const INPUT_FILE = "../../tests/input-data/RawData.sample.json"
const MAP_GOLDEN_FILE =
  "../../tests/output-snapshots/reduce-map-output.golden.json"
const FINALIZE_GOLDEN_FILE =
  "../../tests/output-snapshots/reduce-Features.golden.json"

const PRIVATE_LINE_DIFF_THRESHOLD = 30

// En Intégration Continue, certains tests seront ignorés.
const serialOrSkip = process.env.SKIP_PRIVATE ? "skip" : "serial"

const updateGoldenFiles = process.argv.slice(2).includes("--update")

const readFile = async (filename: string): Promise<string> =>
  util.promisify(fs.readFile)(filename, "utf8")

const writeFile = async (filename: string, data: string): Promise<void> => {
  await util.promisify(fs.writeFile)(filename, data)
  console.warn(`ℹ️ Updated ${filename} => run: $ git secret hide`) // eslint-disable-line no-console
}

const countLines = (str: string) => str.split(/[\r\n]+/).length

// N'affichera le diff complet que si les tests ne tournent pas en CI.
// (pour éviter une fuite de données privée des fichiers golden master)
const safeDeepEqual = (t: ExecCtx, actual: string, expected: string) => {
  if (process.env.CI) {
    const [expectedLines, actualLines] = [expected, actual].map(countLines)
    if (Math.abs(expectedLines - actualLines) > PRIVATE_LINE_DIFF_THRESHOLD) {
      t.fail("the diff is too large => not displaying on CI")
      return
    }
  }
  t.deepEqual(actual, expected)
}

test[serialOrSkip](
  "l'application de reduce.algo2 sur reduce_test_data.json donne le même résultat que d'habitude",
  async (t) => {
    const testData = parseMongoObject(
      await readFile(INPUT_FILE)
    ) as TestDataItem[]

    const f = {
      generatePeriodSerie,
      map,
      finalize,
      reduce,
    }

    // Define global parameters that are required by JS functions
    const jsParams = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any
    jsParams.actual_batch = "2002_1"
    jsParams.date_debut = new Date("2014-01-01")
    jsParams.date_fin = new Date("2016-01-01")
    jsParams.serie_periode = f.generatePeriodSerie(
      jsParams.date_debut,
      jsParams.date_fin
    )
    jsParams.includes = { all: true }
    jsParams.offset_effectif = 2
    jsParams.naf = naf

    const mapResult = runMongoMap(f.map, testData)
    const mapOutput = JSON.stringify(mapResult, null, 2)

    if (updateGoldenFiles) {
      await writeFile(MAP_GOLDEN_FILE, mapOutput)
    }

    const mapExpected = await readFile(MAP_GOLDEN_FILE)
    safeDeepEqual(t, mapOutput, mapExpected)

    const reduceResult = runMongoReduce(
      f.reduce,
      mapResult as { _id: Clé; value: Record<string, EntrepriseEnEntrée> }[]
    )

    const finalizeResult = reduceResult
      .map(({ _id, value }) => f.finalize(_id, value))
      .map((finalizedEntry) => {
        const value = (finalizedEntry as any[])[0]
        delete value.random_order
        return {
          _id: {
            batch: jsParams.actual_batch,
            siret: value.siret,
            periode: value.periode,
          },
          value,
        }
      })

    const finalizeOutput = serializeAsMongoObject(finalizeResult) // finalizeOutput doit être parfaitement identique au golden master qui serait mis à jour depuis test-api-reduce-2.sh => d'où l'appel à serializeAsMongoObject()

    if (updateGoldenFiles) {
      await writeFile(FINALIZE_GOLDEN_FILE, finalizeOutput)
    }

    const finalizeExpected = await readFile(FINALIZE_GOLDEN_FILE)

    safeDeepEqual(t, finalizeOutput, finalizeExpected)
  }
)
