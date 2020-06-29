import "../globals.ts"

type DataType = string // TODO: use BatchDataType instead

/**
 * On recupère les clés ajoutées et les clés supprimées depuis currentBatch.
 * On ajoute aux clés supprimées les types stocks de la memoire.
 */
export function listHashesToAddAndDelete(
  currentBatch: BatchValue,
  stockTypes: DataType[],
  memory: CurrentDataState
): {
  hashToAdd: Record<DataType, Set<DataHash>>
  hashToDelete: Record<DataType, Set<DataHash>>
} {
  const hashToDelete: Record<DataType, Set<DataHash>> = {}
  const hashToAdd: Record<DataType, Set<DataHash>> = {}

  // Itération sur les types qui ont potentiellement subi des modifications
  // pour compléter hashToDelete et hashToAdd.
  // Les suppressions de types complets / stock sont gérés dans le bloc suivant.
  for (const type in currentBatch) {
    // Le type compact gère les clés supprimées
    // Ce type compact existe si le batch en cours a déjà été compacté.
    if (type === "compact") {
      const compactDelete = currentBatch.compact?.delete
      if (compactDelete) {
        Object.keys(compactDelete).forEach((deleteType) => {
          compactDelete[deleteType].forEach((hash) => {
            hashToDelete[deleteType] = hashToDelete[deleteType] || new Set()
            hashToDelete[deleteType].add(hash)
          })
        })
      }
    } else {
      for (const hash in currentBatch[type as keyof BatchValue]) {
        ;(hashToAdd[type] = hashToAdd[type] || new Set()).add(hash)
      }
    }
  }

  stockTypes.forEach((type) => {
    hashToDelete[type] = new Set([
      ...(hashToDelete[type] || new Set()),
      ...memory[type],
    ])
  })

  return {
    hashToAdd,
    hashToDelete,
  }
}
