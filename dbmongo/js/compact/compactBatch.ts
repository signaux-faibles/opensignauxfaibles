import "../globals.ts"

import { listHashesToAddAndDelete } from "./listHashesToAddAndDelete"
import { fixRedundantPatches } from "./fixRedundantPatches"
import { applyPatchesToMemory } from "./applyPatchesToMemory"
import { applyPatchesToBatch } from "./applyPatchesToBatch"

// Paramètres globaux utilisés par "compact"
declare const completeTypes: Record<BatchKey, DataType[]>

/**
 * Appelée par reduce(), compactBatch() va générer un diff entre les
 * données de batch et les données précédentes fournies par memory.
 * Paramètres modifiés: currentBatch et memory.
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

  return currentBatch
}
