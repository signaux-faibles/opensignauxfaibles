import "../globals"

type DataType = string // TODO: use BatchDataType instead

export function applyPatchesToBatch(
  hashToAdd: Record<DataType, Set<DataHash>>,
  hashToDelete: Record<DataType, Set<DataHash>>,
  stockTypes: DataType[],
  currentBatch: BatchValue
): void {
  // Application des suppressions
  stockTypes.forEach((type) => {
    if (hashToDelete[type]) {
      currentBatch.compact = currentBatch.compact || { delete: {} }
      currentBatch.compact.delete = currentBatch.compact.delete || {}
      currentBatch.compact.delete[type] = [...hashToDelete[type]]
    }
  })

  // Application des ajouts
  type AllValueTypesButCompact = Exclude<keyof BatchValue, "compact">
  const typesToAdd = Object.keys(hashToAdd) as AllValueTypesButCompact[]
  typesToAdd.forEach((type) => {
    currentBatch[type] = [...hashToAdd[type]].reduce(
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
    Object.keys(compactDelete).forEach((type) => {
      if (compactDelete[type].length === 0) {
        delete compactDelete[type]
      }
    })
    if (Object.keys(compactDelete).length === 0) {
      delete currentBatch.compact
    }
  }
  // - types vides
  Object.keys(currentBatch).forEach((strType) => {
    const type = strType as keyof BatchValue
    if (Object.keys(currentBatch[type] || {}).length === 0) {
      delete currentBatch[type]
    }
  })
  // TODO: nettoyer le batch ?
}
