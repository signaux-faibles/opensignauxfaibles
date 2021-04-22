// Ces tests visent à couvrir les fonctions du traitement map-reduce purgeBatch

// $ cd js && $(npm bin)/ava ./purgeBatch/ava_tests.ts

import test, { ExecutionContext } from "ava"
import { map } from "../purgeBatch/map"
import { reduce } from "../purgeBatch/reduce"
import { finalize } from "../purgeBatch/finalize"
import { runMongoMap } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"
import { CompanyDataValues } from "../RawDataTypes"

test.serial(
  `purgeBatch() ne conserve que les données des batches antérieurs au batch fourni`,
  (t: ExecutionContext) => {
    setGlobals({ fromBatchKey: "1911_2" }) // important: utiliser test.serial() pour que ces paramètres ne soient utilisés que pour ce test
    const rawData: CompanyDataValues = {
      batch: {
        "1910": {},
        "1911": {},
        "1911_1": {},
        "1911_2": {}, // <-- on va purger à partir de ce batch
        "1911_3": {},
        "1912": {},
        "1912_1": {},
      },
      scope: "etablissement",
      key: "01234567891011",
    }
    const mapResults = runMongoMap(map, [{ _id: null, value: rawData }]) as {
      value: CompanyDataValues
    }[]
    const reduceResults = reduce(
      {},
      mapResults.map(({ value }) => value)
    )
    const finalizeResultValue = reduceResults.map((value) =>
      finalize({ scope: rawData.scope }, value)
    )
    t.is(finalizeResultValue.length, 1)
    finalizeResultValue.forEach((finalizedValue) =>
      t.deepEqual(finalizedValue.batch, {
        "1910": {},
        "1911": {},
        "1911_1": {},
      })
    )
  }
)
