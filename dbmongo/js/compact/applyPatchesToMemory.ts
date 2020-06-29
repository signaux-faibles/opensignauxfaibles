import "../globals"

type DataType = string // TODO: use BatchDataType instead

export function applyPatchesToMemory(
  hashToAdd: Record<DataType, Set<DataHash>>,
  hashToDelete: Record<DataType, Set<DataHash>>,
  memory: CurrentDataState
): void {

  Object.keys(hashToDelete).forEach((type) => {
    // 3.c On retire les cles restantes de la memoire.
    // --------------------------------------------------
    hashToDelete[type].forEach((hash) => {
      memory[type].delete(hash)
    })
  })

  Object.keys(hashToAdd).forEach((type) => {
    //
    // 4.b Pour chaque cle ajoutee restante: on ajoute à la memoire.
    // -------------------------------------------------------------

    hashToAdd[type].forEach((hash) => {
      memory[type] = memory[type] || new Set()
      memory[type].add(hash)
    })
  })
}
