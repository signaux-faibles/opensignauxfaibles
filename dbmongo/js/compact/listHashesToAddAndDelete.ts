import { f } from "./functions"
import { DataType, DataHash } from "../RawDataTypes"
import { CurrentDataState } from "./currentState"
import { BatchValueWithCompact } from "./applyPatchesToBatch"

/**
 * On recupère les clés ajoutées et les clés supprimées depuis currentBatch.
 * On ajoute aux clés supprimées les types stocks de la memoire.
 */
export function listHashesToAddAndDelete(
  currentBatch: BatchValueWithCompact,
  stockTypes: DataType[],
  memory: CurrentDataState
): {
  hashToAdd: Partial<Record<DataType, Set<DataHash>>>
  hashToDelete: Partial<Record<DataType, Set<DataHash>>>
} {
  const hashToDelete: Partial<Record<DataType, Set<DataHash>>> = {}
  const hashToAdd: Partial<Record<DataType, Set<DataHash>>> = {}

  // Itération sur les types qui ont potentiellement subi des modifications
  // pour compléter hashToDelete et hashToAdd.
  // Les suppressions de types complets / stock sont gérés dans le bloc suivant.
  f.forEachPopulatedProp(currentBatch, (type) => {
    // Le type compact gère les clés supprimées
    // Ce type compact existe si le batch en cours a déjà été compacté.
    if (type === "compact") {
      const compactDelete = currentBatch.compact?.delete
      if (compactDelete) {
        f.forEachPopulatedProp(compactDelete, (deleteType, keysToDelete) => {
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
