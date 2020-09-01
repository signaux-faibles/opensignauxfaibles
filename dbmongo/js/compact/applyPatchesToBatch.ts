import "../globals"
import { forEachPopulatedProp } from "../common/forEachPopulatedProp"
import {
  DataType,
  BatchKey,
  BatchValue,
  Scope,
  IndexFlags,
} from "../RawDataTypes"

export type BatchValueWithCompact = BatchValue & {
  compact?: { delete: Partial<Record<DataType, DataHash[]>> }
}

export type CompanyDataValuesWithCompact = {
  key: SiretOrSiren
  scope: Scope
  batch: Record<BatchKey, BatchValueWithCompact>
} & Partial<IndexFlags>

export function applyPatchesToBatch(
  hashToAdd: Partial<Record<DataType, Set<DataHash>>>,
  hashToDelete: Partial<Record<DataType, Set<DataHash>>>,
  stockTypes: DataType[],
  currentBatch: BatchValueWithCompact
): void {
  // Application des suppressions
  stockTypes.forEach((type) => {
    const hashesToDelete = hashToDelete[type]
    if (hashesToDelete) {
      currentBatch.compact = currentBatch.compact || { delete: {} }
      currentBatch.compact.delete = currentBatch.compact.delete || {}
      currentBatch.compact.delete[type] = [...hashesToDelete]
    }
  })

  // Application des ajouts
  forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
    currentBatch[type] = [...hashesToAdd].reduce(
      (typedBatchValues, hash) => ({
        ...typedBatchValues,
        [hash]: currentBatch[type]?.[hash],
      }),
      {}
    )
  })

  // Retrait des propriété vides
  // - compact.delete vides
  const compactDelete = currentBatch.compact?.delete
  if (compactDelete) {
    forEachPopulatedProp(compactDelete, (type, keysToDelete) => {
      if (keysToDelete.length === 0) {
        delete compactDelete[type]
      }
    })
    if (Object.keys(compactDelete).length === 0) {
      delete currentBatch.compact
    }
  }
  // - types vides
  forEachPopulatedProp(currentBatch, (type, typedBatchData) => {
    if (Object.keys(typedBatchData).length === 0) {
      delete currentBatch[type]
    }
  })
}
