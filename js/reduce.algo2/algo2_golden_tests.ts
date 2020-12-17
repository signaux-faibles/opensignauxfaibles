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
import { map, EntréeMap, CléSortieMap, SortieMap } from "./map"
import { finalize } from "./finalize"
import { reduce } from "./reduce"
import { TestDataItem } from "../test/data/objects"
import { setGlobals } from "../test/helpers/setGlobals"
import {
  runMongoMap,
  runMongoReduce,
  parseMongoObject,
  serializeAsMongoObject,
} from "../test/helpers/mongodb"
import { compare } from "concordance"

const INPUT_FILE = "../tests/input-data/RawData.sample.json"
const MAP_GOLDEN_FILE =
  "../tests/output-snapshots/reduce-map-output.golden.json"
const FINALIZE_GOLDEN_FILE =
  "../tests/output-snapshots/reduce-Features.golden.json"

// En Intégration Continue, certains tests seront ignorés.
const serialOrSkip = process.env.SKIP_PRIVATE ? "skip" : "serial"

const updateGoldenFiles = process.argv.slice(2).includes("--update")

const readFile = async (filename: string): Promise<string> =>
  util.promisify(fs.readFile)(filename, "utf8")

const writeFile = async (filename: string, data: string): Promise<void> => {
  await util.promisify(fs.writeFile)(filename, data)
  console.warn(`ℹ️ Updated ${filename} => run: $ git secret hide`) // eslint-disable-line no-console
}

// N'affichera le diff complet que si les tests ne tournent pas en CI.
// (pour éviter une fuite de données privée des fichiers golden master)
const safeDeepEqual = (t: ExecCtx, actual: unknown, expected: unknown) => {
  if (process.env.CI) {
    const { pass } = compare(actual, expected)
    if (!pass) {
      t.fail(
        "Results don't match the snapshot. Diff of private data forbidden on CI."
      )
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

    // Define global parameters that are required by JS functions
    const date_debut = new Date("2014-01-01")
    const date_fin = new Date("2016-01-01")
    const jsParams = {
      actual_batch: "2002_1",
      date_debut,
      date_fin,
      serie_periode: generatePeriodSerie(date_debut, date_fin),
      includes: { all: true },
      offset_effectif: 2,
      naf,
    }
    setGlobals(jsParams)

    const mapResult = runMongoMap<EntréeMap, CléSortieMap, SortieMap>(
      map,
      testData
    )
    const mapOutput = JSON.stringify(mapResult, null, 2)

    if (updateGoldenFiles) {
      await writeFile(MAP_GOLDEN_FILE, mapOutput)
    }

    const mapExpected = await readFile(MAP_GOLDEN_FILE)
    safeDeepEqual(t, parseMongoObject(mapOutput), parseMongoObject(mapExpected))

    const reduceResult = runMongoReduce(reduce, mapResult)

    const finalizeResult = reduceResult
      .map(({ _id, value }) => finalize(_id, value))
      .map((finalizedEntry) => {
        if (
          typeof finalizedEntry === "undefined" ||
          "incomplete" in finalizedEntry ||
          finalizedEntry[0] === undefined
        ) {
          return {}
        }
        const value = finalizedEntry[0]
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

    const finalizeOutput = serializeAsMongoObject(finalizeResult) // finalizeOutput doit être parfaitement identique au golden master qui serait mis à jour depuis test-reduce-2.sh => d'où l'appel à serializeAsMongoObject()
    if (updateGoldenFiles) {
      await writeFile(FINALIZE_GOLDEN_FILE, finalizeOutput)
    }

    const finalizeExpected = await readFile(FINALIZE_GOLDEN_FILE)
    safeDeepEqual(
      t,
      parseMongoObject(finalizeOutput),
      parseMongoObject(finalizeExpected)
    )
  }
)
