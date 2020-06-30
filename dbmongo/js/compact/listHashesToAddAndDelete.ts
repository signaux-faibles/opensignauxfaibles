import "../globals.ts"
import { forEachPopulatedProp } from "../common/forEachPopulatedProp"

/**
 * On recupère les clés ajoutées et les clés supprimées depuis currentBatch.
 * On ajoute aux clés supprimées les types stocks de la memoire.
 */
export function listHashesToAddAndDelete(
  currentBatch: BatchValue,
  stockTypes: BatchDataType[],
  memory: CurrentDataState
): {
  hashToAdd: Partial<Record<BatchDataType, Set<DataHash>>>
  hashToDelete: Partial<Record<BatchDataType, Set<DataHash>>>
} {
  const hashToDelete: Partial<Record<BatchDataType, Set<DataHash>>> = {}
  const hashToAdd: Partial<Record<BatchDataType, Set<DataHash>>> = {}

  // Itération sur les types qui ont potentiellement subi des modifications
  // pour compléter hashToDelete et hashToAdd.
  // Les suppressions de types complets / stock sont gérés dans le bloc suivant.
  forEachPopulatedProp(currentBatch, (type) => {
    // Le type compact gère les clés supprimées
    // Ce type compact existe si le batch en cours a déjà été compacté.
    if (type === "compact") {
      const compactDelete = currentBatch.compact?.delete
      if (compactDelete) {
        forEachPopulatedProp(compactDelete, (deleteType, keysToDelete) => {
          keysToDelete.forEach((hash) => {
            ;(hashToDelete[deleteType] =
              hashToDelete[deleteType] || new Set()).add(hash)
          })
        })
      }
    } else {
      for (const hash in currentBatch[type]) {
        ;(hashToAdd[type] = hashToAdd[type] || new Set()).add(hash)
      }
    }
  })

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
