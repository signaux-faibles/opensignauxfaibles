import "../globals.ts"

type DataType = string // TODO: use BatchDataType instead

export function listHashesToAddAndDelete(
  currentBatch: BatchValue
): {
  hashToAdd: Record<DataType, Set<DataHash>>
  hashToDelete: Record<DataType, Set<DataHash>>
} {
  // 1. On recupère les cles ajoutes et les cles supprimes
  // -----------------------------------------------------

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
        Object.keys(compactDelete).forEach((delete_type) => {
          compactDelete[delete_type].forEach((hash) => {
            hashToDelete[delete_type] = hashToDelete[delete_type] || new Set()
            hashToDelete[delete_type].add(hash)
          })
        })
      }
    } else {
      for (const hash in currentBatch[type as keyof BatchValue]) {
        ;(hashToAdd[type] = hashToAdd[type] || new Set()).add(hash)
      }
    }
  }
  return {
    hashToAdd,
    hashToDelete,
  }
}
