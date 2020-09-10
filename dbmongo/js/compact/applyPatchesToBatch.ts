import "../globals"
import { forEachPopulatedProp } from "../common/forEachPopulatedProp"

export function applyPatchesToBatch(
  hashToAdd: Partial<Record<DataType, Set<DataHash>>>,
  hashToDelete: Partial<Record<DataType, Set<DataHash>>>,
  stockTypes: DataType[],
  currentBatch: BatchValue
): void {
  const f = { forEachPopulatedProp } /*DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO*/
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
  f.forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
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
    f.forEachPopulatedProp(compactDelete, (type, keysToDelete) => {
      if (keysToDelete.length === 0) {
        delete compactDelete[type]
      }
    })
    if (Object.keys(compactDelete).length === 0) {
      delete currentBatch.compact
    }
  }
  // - types vides
  f.forEachPopulatedProp(currentBatch, (type, typedBatchData) => {
    if (Object.keys(typedBatchData).length === 0) {
      delete currentBatch[type]
    }
  })
}
