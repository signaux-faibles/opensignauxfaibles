/*global globalThis*/

// Combinaison des tests de map_test.js et finalize_test.js.
//
// Golden-file-based tests to prevent regressions on the JS functions
// (common + algo2) used to compute the "Features" collection from the
// "RawData" collection.
//
// Please execute ./test_algo2.sh to run this test suite.

import test, { before, after } from "ava"
import * as fs from "fs"
import * as util from "util"
import * as childProcess from "child_process"
import { naf } from "../test/data/naf"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { map } from "../reduce.algo2/map"
import { finalize } from "../reduce.algo2/finalize"
import { reduce } from "../reduce.algo2/reduce"
import { runMongoMap } from "../test/helpers/mongodb"

// En Intégration Continue, certains tests seront ignorés.
const serialOrSkip = process.env.CI ? "skip" : "serial"

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
  }
})()

const loadTestData = async (filename: string) =>
  JSON.parse(
    (await context.readFile(filename))
      .replace(/ISODate\("([^"]+)"\)/g, `{ "_ISODate": "$1" }`)
      .replace(/NumberInt\(([^)]+)\)/g, "$1"),
    (_key, value: unknown) =>
      value && typeof value === "object" && (value as any)._ISODate
        ? new Date((value as any)._ISODate)
        : value
  )

before("récupération des données", async () => {
  await context.setup() // step will fail in case of error while downloading golden files
})

after("suppression des données temporaires", async () => {
  await context.tearDown()
})

test[serialOrSkip](
  "l'application de reduce.algo2 sur reduce_test_data.json donne le même résultat que d'habitude",
  async (t) => {
    const testData = await loadTestData("reduce_test_data.json")

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

    const mapResult = runMongoMap(
      f.map,
      testData as any[] // TODO: as { _id: string; value: CompanyDataValuesWithFlags }[]
    ) // -> [ { _id, value } ]

    const mapOutput = JSON.stringify(mapResult, null, 2)
    const mapExpected = await context.readFile("map_golden.log")
    t.deepEqual(mapOutput, mapExpected.trim())

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
    const finalizeExpected = await context.readFile("finalize_golden.log")
    t.deepEqual(finalizeOutput, finalizeExpected.trim())
  }
)
