/*global globalThis*/

// Combinaison des tests de map_test.js et finalize.js.

// Context: this golden-file-based test runner was designed to prevent
// regressions on the JS functions (common + algo2) used to compute the
// "Features" collection from the "RawData" collection.
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

const serialOrSkip = process.env.CI ? "skip" : "serial"

const exec = (command: string) =>
  new Promise((resolve, reject) =>
    childProcess.exec(
      command,
      (err: Error | null, stdout: string, stderr: string) =>
        err ? reject(err) : resolve({ stdout, stderr })
    )
  )

const context = (() => {
  const remotePath = "stockage:/home/centos/opensignauxfaibles_tests"
  const goldenPath = "./test_data_algo2"
  const outFile = `${goldenPath}/algo2_stdout.log`
  // const goldenFileContent: Record<string, string> = {}
  // const promisedDownload: Promise<void> | null = null

  const getGoldenFile = async (filename: string): Promise<string> => {
    /*
    if (goldenFileContent[filename]) return goldenFileContent[filename]
    if (!promisedDownload) {
      promisedDownload = null
    }
    await promiseToDownload
    */
    return util.promisify(fs.readFile)(`${goldenPath}/${filename}`, "utf8")
  }

  const print = async (content: string) =>
    util.promisify(fs.appendFile)(outFile, content + "\n")

  return {
    setup: async (t) => {
      await exec(`mkdir ${goldenPath} | true`)
      const command = `scp ${remotePath}/* ${goldenPath}`
      console.warn(`$ ${command}`) // eslint-disable-line no-console
      const { stderr } = (await exec(command)) as { stderr: string }
      t.error(new Error(stderr))
      // prepare the outFile
      util.promisify(fs.writeFile)(outFile, "")
    },
    tearDown: () => exec(`rm -r ${goldenPath}`),
    getGoldenFile,
    print,
  }
})()

const loadTestData = async (filename: string) => {
  const content = await context.getGoldenFile(filename)
  return JSON.parse(
    content
      .replace(/ISODate\("([^"]+)"\)/g, `{ "_ISODate": "$1" }`)
      .replace(/NumberInt\(([^)]+)\)/g, "$1"),
    (_key, value: unknown) =>
      value && typeof value === "object" && (value as any)._ISODate
        ? new Date((value as any)._ISODate)
        : value
  )
  // TODO: ISODate: (date: string) => new Date(date.replace("+0000", "+00:00")), // make sure that timezone format complies with the spec
}

before("préparation des golden files", async () => {
  await context.setup()
})

after("libération des golden files", async () => {
  // await context.tearDown() // TODO
})

test[serialOrSkip](
  "l'application de reduce.algo2 sur reduce_test_data.json donne le même résultat que d'habitude",
  async (t) => {
    const testData = await loadTestData("reduce_test_data.json")
    // console.log(util.inspect(testData, { depth: Infinity, colors: true }))

    const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any

    const f = {
      generatePeriodSerie,
      map,
      finalize,
      reduce,
    }

    // Define global parameters that are required by JS functions
    const jsParams = global
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

    // Print the output of the f.map() function
    await context.print(JSON.stringify(mapResult, null, 2))

    const valuesPerKey: Record<string, unknown[]> = {}
    mapResult.forEach(({ _id, value }) => {
      const idString = JSON.stringify(_id)
      valuesPerKey[idString] = valuesPerKey[idString] || []
      valuesPerKey[idString].push(value)
    })

    const finalizeResult = Object.keys(valuesPerKey).map((key) =>
      f.finalize(JSON.parse(key), f.reduce(key, valuesPerKey[key]))
    )

    // Print the output of the f.finalize() function
    await context.print(JSON.stringify(finalizeResult, null, 2))

    // await new Promise((resolve) => setTimeout(resolve, 2000))

    t.pass()
  }
)
