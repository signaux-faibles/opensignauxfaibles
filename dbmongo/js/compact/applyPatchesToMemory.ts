import "../globals"
import { forEachPopulatedProp } from "../common/forEachPopulatedProp"

export function applyPatchesToMemory(
  hashToAdd: Partial<Record<BatchDataType, Set<DataHash>>>,
  hashToDelete: Partial<Record<BatchDataType, Set<DataHash>>>,
  memory: CurrentDataState
): void {
  // Prise en compte des suppressions de clés dans la mémoire
  forEachPopulatedProp(hashToDelete, (type, hashesToDelete) => {
    hashesToDelete.forEach((hash) => {
      memory[type].delete(hash)
    })
  })

  // Prise en compte des ajouts de clés dans la mémoire
  forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
    hashesToAdd.forEach((hash) => {
      memory[type] = memory[type] || new Set()
      memory[type].add(hash)
    })
  })
}
