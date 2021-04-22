// $ cd js && $(npm bin)/ava ./reduce.algo2/ccsf_tests.ts

import test from "ava"
import { ccsf } from "./ccsf"
import { ParHash } from "../RawDataTypes"
import { EntréeCcsf } from "../GeneratedTypes"

const makeUTCDate = (year: number, month: number, day?: number): Date =>
  new Date(Date.UTC(year, month, day || 1))

test(`ccsf retourne la date de début de la procédure CCSF pour chaque période`, (t) => {
  const vCcsf: ParHash<EntréeCcsf> = {
    hash1: { date_traitement: makeUTCDate(2021, 2, 15) } as EntréeCcsf,
    hash2: { date_traitement: makeUTCDate(2021, 0, 15) } as EntréeCcsf,
    hash3: { date_traitement: makeUTCDate(2021, 1, 15) } as EntréeCcsf,
  }
  const output: { periode: Date; date_ccsf?: Date }[] = [
    { periode: makeUTCDate(2021, 0) },
    { periode: makeUTCDate(2021, 1) },
    { periode: makeUTCDate(2021, 2) },
  ]
  ccsf(vCcsf, output)
  t.deepEqual(output, [
    { periode: makeUTCDate(2021, 0) },
    { periode: makeUTCDate(2021, 1), date_ccsf: makeUTCDate(2021, 0, 15) },
    { periode: makeUTCDate(2021, 2), date_ccsf: makeUTCDate(2021, 1, 15) },
  ])
})
