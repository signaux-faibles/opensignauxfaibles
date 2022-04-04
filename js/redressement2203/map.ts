/* eslint-disable @typescript-eslint/no-unused-vars */
import {
  CompanyDataValues,
  BatchKey,
  SortieRedressementUrssaf2203,
} from "../RawDataTypes"

export type SortieMap = SortieRedressementUrssaf2203

// Paramètres globaux utilisés par "public"
declare const actual_batch: BatchKey

// Types de données en entrée et sortie
export type Input = { _id: unknown; value: CompanyDataValues }
export type OutKey = string
export type OutValue = Partial<SortieMap>
declare function emit(key: string, value: OutValue): void

export function map(this: Input): void {
  const testDate = new Date("2021-09-01")
  const batches = this.value.batch
  const beforeBatches = []
  const afterBatches = []
  for (const [_, value] of Object.entries(batches)) {
    if (value.debit) {
      const debit = value.debit
      for (const [_, value2] of Object.entries(debit)) {
        value2.periode.start > testDate
          ? afterBatches.push(value2)
          : beforeBatches.push(value2)
      }
    }
  }
  const latestBatchBeforeDate = beforeBatches.reduce((a, b) =>
    a.periode.start > b.periode.start ? a : b
  )
  const latestBatchAfterDate = afterBatches.reduce((a, b) =>
    a.periode.start > b.periode.start ? a : b
  )

  emit(this.value.key, {
    partPatronaleAncienne: latestBatchBeforeDate.part_patronale,
    partOuvriereAncienne: latestBatchBeforeDate.part_ouvriere,
    partPatronaleRecente: latestBatchAfterDate.part_patronale,
    partOuvriereRecente: latestBatchAfterDate.part_ouvriere,
  })
}
