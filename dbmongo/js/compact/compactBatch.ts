import "../globals.ts"

import { listHashesToAddAndDelete } from "./listHashesToAddAndDelete"
import { fixRedundantPatches } from "./fixRedundantPatches"
import { applyPatchesToMemory } from "./applyPatchesToMemory"
import { applyPatchesToBatch } from "./applyPatchesToBatch"

/**
 * Appelée par reduce(), compactBatch() va générer un diff entre les
 * données de batch et les données précédentes fournies par memory.
 * Pré-requis: les batches précédents doivent avoir été compactés.
 */
export function compactBatch(
  currentBatch: BatchValue,
  memory: CurrentDataState,
  batchKey: string
): BatchValue {
  // Les types où il y a potentiellement des suppressions
  const stockTypes = completeTypes[batchKey].filter(
    (type) => (memory[type] || new Set()).size > 0
  )

  const { hashToAdd, hashToDelete } = listHashesToAddAndDelete(
    currentBatch,
    stockTypes,
    memory
  )

  fixRedundantPatches(hashToAdd, hashToDelete, memory)
  applyPatchesToMemory(hashToAdd, hashToDelete, memory)
  applyPatchesToBatch(hashToAdd, hashToDelete, stockTypes, currentBatch)

  // 6. nettoyage
  // ------------

  if (currentBatch) {
    //types vides
    Object.keys(currentBatch).forEach((strType) => {
      const type = strType as keyof BatchValue
      if (Object.keys(currentBatch[type] || {}).length === 0) {
        delete currentBatch[type]
      }
    })
    //hash à supprimer vides (compact.delete)
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
  }
  return currentBatch
}
