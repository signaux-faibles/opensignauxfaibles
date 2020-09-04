import { forEachPopulatedProp } from "../common/forEachPopulatedProp"
import { DataType, DataHash } from "../RawDataTypes"
import { CurrentDataState } from "./currentState"

/**
 * Modification de hashToAdd et hashToDelete pour retirer les redondances.
 **/
export function fixRedundantPatches(
  hashToAdd: Partial<Record<DataType, Set<DataHash>>>,
  hashToDelete: Partial<Record<DataType, Set<DataHash>>>,
  memory: CurrentDataState
): void {
  forEachPopulatedProp(hashToDelete, (type, hashesToDelete) => {
    // Pour chaque cle supprimee: est-ce qu'elle est bien dans la
    // memoire ? sinon on la retire de la liste des clés supprimées (pas de
    // maj memoire)
    // -----------------------------------------------------------------------------------------------------------------
    hashToDelete[type] = new Set(
      [...hashesToDelete].filter((hash) => {
        return (memory[type] || new Set()).has(hash)
      })
    )

    // Est-ce qu'elle a ete egalement ajoutee en même temps que
    // supprimée ? (par exemple remplacement d'un stock complet à
    // l'identique) Dans ce cas là, on retire cette clé des valeurs ajoutées
    // et supprimées
    // i.e. on herite de la memoire. (pas de maj de la memoire)
    // ------------------------------------------------------------------------------
    hashToDelete[type] = new Set(
      [...(hashToDelete[type] || new Set())].filter((hash) => {
        const hashesToAdd = hashToAdd[type] || new Set()
        const also_added = hashesToAdd.has(hash)
        if (also_added) {
          hashesToAdd.delete(hash)
        }
        return !also_added
      })
    )
  })

  forEachPopulatedProp(hashToAdd, (type, hashesToAdd) => {
    // Pour chaque cle ajoutee: est-ce qu'elle est dans la memoire ? Si oui on filtre cette cle
    // i.e. on herite de la memoire. (pas de maj de la memoire)
    // ---------------------------------------------------------------------------------------------
    hashToAdd[type] = new Set(
      [...hashesToAdd].filter((hash) => {
        return !(memory[type] || new Set()).has(hash)
      })
    )
  })
}
