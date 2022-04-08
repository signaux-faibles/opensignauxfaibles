/* eslint-disable @typescript-eslint/no-unused-vars */
import { f } from "./functions"
import {
  CompanyDataValues,
  SortieRedressementUrssaf2203,
} from "../RawDataTypes"

export type SortieMap = SortieRedressementUrssaf2203

declare const dateStr: string

// Types de données en entrée et sortie
export type Input = { _id: unknown; value: CompanyDataValues }
export type OutKey = string
export type OutValue = Partial<SortieMap>
declare function emit(key: string, value: OutValue): void

export function map(this: Input): void {
  const testDate = new Date(dateStr)

  const values = f.flatten(this.value, "2203")
  const beforeBatches = []
  const afterBatches = []
  if (values.debit) {
    for (const debit of Object.values(values.debit)) {
      debit.periode.start > testDate
        ? afterBatches.push(debit)
        : beforeBatches.push(debit)
    }
  }

  const latestBatchBeforeDate =
    beforeBatches.length > 0
      ? beforeBatches.reduce((a, b) =>
          a.periode.start > b.periode.start ? a : b
        )
      : null
  const latestBatchAfterDate =
    afterBatches.length > 0
      ? afterBatches.reduce((a, b) =>
          a.periode.start > b.periode.start ? a : b
        )
      : null

  emit(this.value.key, {
    partPatronaleAncienne: latestBatchBeforeDate
      ? latestBatchBeforeDate.part_patronale
      : 0,
    partOuvriereAncienne: latestBatchBeforeDate
      ? latestBatchBeforeDate.part_ouvriere
      : 0,
    partPatronaleRecente: latestBatchAfterDate
      ? latestBatchAfterDate.part_patronale
      : 0,
    partOuvriereRecente: latestBatchAfterDate
      ? latestBatchAfterDate.part_ouvriere
      : 0,
  })
}
