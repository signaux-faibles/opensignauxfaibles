import { f } from "./functions"
import { DataType, DataHash } from "../RawDataTypes"
import { CurrentDataState } from "./currentState"

export function applyPatchesToMemory(
  hashToAdd: Partial<Record<DataType, Set<DataHash>>>,
  hashToDelete: Partial<Record<DataType, Set<DataHash>>>,
  memory: CurrentDataState
): void {
  // Prise en compte des suppressions de clés dans la mémoire
  f.forEachPopulatedProp(hashToDelete, (type, hashesToDelete) => {
    hashesToDelete.forEach((hash) => {
      memory[type].delete(hash)
    })
  })

  // Prise en compte des ajouts de clés dans la mémoire
  f.forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
    hashesToAdd.forEach((hash) => {
      memory[type] = memory[type] || new Set()
      memory[type].add(hash)
    })
  })
}
