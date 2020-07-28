/*global globalThis*/

// Combinaison des tests de map_test.js et finalize_test.js.
//
// Golden-file-based tests to prevent regressions on the JS functions
// (common + algo2) used to compute the "Features" collection from the
// "RawData" collection.
//
// To update golden files: `$ npx ava algo2_golden_tests.ts -- --update`
//                      or `$ npm run test:update-golden-files`

import test, { before, after } from "ava"
import * as fs from "fs"
import * as util from "util"
import * as childProcess from "child_process"
import { naf } from "../test/data/naf"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { map } from "./map"
import { finalize } from "./finalize"
import { reduce } from "./reduce"
import { runMongoMap, parseMongoObject } from "../test/helpers/mongodb"

const INPUT_FILE = "reduce_test_data.json"
const MAP_GOLDEN_FILE = "map_golden.log"
const FINALIZE_GOLDEN_FILE = "finalize_golden.log"

// En Intégration Continue, certains tests seront ignorés.
const serialOrSkip = process.env.CI ? "skip" : "serial"

const updateGoldenFiles = process.argv.slice(2).includes("--update")

const exec = (command: string): Promise<{ stdout: string; stderr: string }> =>
  util.promisify(childProcess.exec)(command)

const context = (() => {
  const remotePath = "stockage:/home/centos/opensignauxfaibles_tests"
  const localPath = "./test_data_algo2"

  return {
    setup: async () => {
      await exec(`mkdir ${localPath} | true`)
      const command = `scp ${remotePath}/* ${localPath}`
      console.warn(`$ ${command}`) // eslint-disable-line no-console
      const { stderr } = await exec(command)
      if (stderr) throw new Error(stderr)
    },
    tearDown: () => exec(`rm -r ${localPath}`),
    readFile: async (filename: string): Promise<string> =>
      util.promisify(fs.readFile)(`${localPath}/${filename}`, "utf8"),
    writeFile: async (filename: string, data: string): Promise<void> => {
      await util.promisify(fs.writeFile)(`${localPath}/${filename}`, data)
      await exec(`scp ${localPath}/${filename} ${remotePath}/`)
    },
  }
})()

before("récupération des données", async () => {
  await context.setup() // step will fail in case of error while downloading golden files
})

after("suppression des données temporaires", async () => {
  await context.tearDown()
})

test[serialOrSkip](
  "l'application de reduce.algo2 sur reduce_test_data.json donne le même résultat que d'habitude",
  async (t) => {
    type TestDataItem = { _id: string; value: CompanyDataValuesWithFlags }
    const testData = parseMongoObject(
      await context.readFile(INPUT_FILE)
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
      await context.writeFile(MAP_GOLDEN_FILE, mapOutput)
    }

    const mapExpected = await context.readFile(MAP_GOLDEN_FILE)
    t.deepEqual(mapOutput, mapExpected)

    const valuesPerKey: Record<string, unknown[]> = {}
    mapResult.forEach(({ _id, value }) => {
      const idString = JSON.stringify(_id)
      valuesPerKey[idString] = valuesPerKey[idString] || []
      valuesPerKey[idString].push(value)
    })

    const finalizeResult = Object.keys(valuesPerKey).map((key) =>
      f.finalize(JSON.parse(key), f.reduce(key, valuesPerKey[key]))
    )
    const finalizeOutput = JSON.stringify(finalizeResult, null, 2)

    if (updateGoldenFiles) {
      await context.writeFile(FINALIZE_GOLDEN_FILE, finalizeOutput)
    }

    const finalizeExpected = await context.readFile(FINALIZE_GOLDEN_FILE)
    t.deepEqual(finalizeOutput, finalizeExpected)
  }
)
