import "../globals"
import { forEachPopulatedProp } from "../common/forEachPopulatedProp"

export function applyPatchesToMemory(
  hashToAdd: Partial<Record<DataType, Set<DataHash>>>,
  hashToDelete: Partial<Record<DataType, Set<DataHash>>>,
  memory: CurrentDataState
): void {
  const f = { forEachPopulatedProp } /*DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO*/
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
