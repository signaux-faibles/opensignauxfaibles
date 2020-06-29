import "../globals"

type DataType = string // TODO: use BatchDataType instead

export function applyPatchesToMemory(
  hashToAdd: Record<DataType, Set<DataHash>>,
  hashToDelete: Record<DataType, Set<DataHash>>,
  memory: CurrentDataState
): void {

  // On retire les cles restantes de la memoire.
  Object.keys(hashToDelete).forEach((type) => {
    hashToDelete[type].forEach((hash) => {
      memory[type].delete(hash)
    })
  })

  // Pour chaque cle ajoutee restante: on ajoute à la memoire.
  Object.keys(hashToAdd).forEach((type) => {
    hashToAdd[type].forEach((hash) => {
      memory[type] = memory[type] || new Set()
      memory[type].add(hash)
    })
  })
}
