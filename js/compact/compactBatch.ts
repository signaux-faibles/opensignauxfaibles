import { f } from "./functions"
import { BatchValueWithCompact } from "./applyPatchesToBatch"
import { DataType, BatchValue, BatchKey } from "../RawDataTypes"
import { CurrentDataState } from "./currentState"

// Paramètres globaux utilisés par "compact"
declare const completeTypes: Record<BatchKey, DataType[]>

/**
 * Appelée par reduce(), compactBatch() va générer un diff entre les
 * données de batch et les données précédentes fournies par memory.
 * Paramètres modifiés: currentBatch et memory.
 * Pré-requis: les batches précédents doivent avoir été compactés.
 */
export function compactBatch(
  currentBatch: BatchValueWithCompact,
  memory: CurrentDataState,
  fromBatchKey: string
): BatchValue {
  // Les types où il y a potentiellement des suppressions
  const stockTypes = (completeTypes[fromBatchKey] ?? []).filter(
    (type) => (memory[type] || new Set()).size > 0
  )

  const { hashToAdd, hashToDelete } = f.listHashesToAddAndDelete(
    currentBatch,
    stockTypes,
    memory
  )

  f.fixRedundantPatches(hashToAdd, hashToDelete, memory)
  f.applyPatchesToMemory(hashToAdd, hashToDelete, memory)
  f.applyPatchesToBatch(hashToAdd, hashToDelete, stockTypes, currentBatch)

  return currentBatch
}
