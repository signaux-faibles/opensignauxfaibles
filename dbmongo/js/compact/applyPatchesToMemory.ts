import "../globals"

type DataType = string // TODO: use BatchDataType instead

export function applyPatchesToMemory(
  hashToAdd: Record<DataType, Set<DataHash>>,
  hashToDelete: Record<DataType, Set<DataHash>>,
  memory: CurrentDataState
): void {
  // Prise en compte des suppressions de clés dans la mémoire
  Object.keys(hashToDelete).forEach((type) => {
    hashToDelete[type].forEach((hash) => {
      memory[type].delete(hash)
    })
  })

  // Prise en compte des ajouts de clés dans la mémoire
  Object.keys(hashToAdd).forEach((type) => {
    hashToAdd[type].forEach((hash) => {
      memory[type] = memory[type] || new Set()
      memory[type].add(hash)
    })
  })
}
