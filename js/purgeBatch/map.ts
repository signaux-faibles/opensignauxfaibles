import { BatchKey, CompanyDataValues } from "../RawDataTypes"

// Paramètres globaux utilisés par "public"
declare const fromBatchKey: BatchKey

declare function emit(key: unknown, value: CompanyDataValues): void

export type Input = { _id: unknown; value: CompanyDataValues }

export function map(this: Input): void {
  "use strict"
  const batches = Object.keys(this.value.batch)
  batches
    .filter((key) => key >= fromBatchKey)
    .forEach((key) => {
      delete this.value.batch[key]
    })
  // With a merge output, sending a new object, even empty, is compulsory
  emit(this._id, this.value)
}
