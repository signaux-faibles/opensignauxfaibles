

type DataType = string // TODO: use BatchDataType instead

/**
 * Modification de hashToAdd et hashToDelete pour retirer les redondances.
**/
export function fixRedundantPatches(
  hashToAdd: Record<DataType, Set<DataHash>>,
  hashToDelete: Record<DataType, Set<DataHash>>,
  memory: CurrentDataState
) {
  Object.keys(hashToDelete).forEach((type) => {
    // 3.a Pour chaque cle supprimee: est-ce qu'elle est bien dans la
    // memoire ? sinon on la retire de la liste des clés supprimées (pas de
    // maj memoire)
    // -----------------------------------------------------------------------------------------------------------------
    hashToDelete[type] = new Set(
      [...hashToDelete[type]].filter((hash) => {
        return (memory[type] || new Set()).has(hash)
      })
    )

    // 3.b Est-ce qu'elle a ete egalement ajoutee en même temps que
    // supprimée ? (par exemple remplacement d'un stock complet à
    // l'identique) Dans ce cas là, on retire cette clé des valeurs ajoutées
    // et supprimées
    // i.e. on herite de la memoire. (pas de maj de la memoire)
    // ------------------------------------------------------------------------------

    hashToDelete[type] = new Set(
      [...hashToDelete[type]].filter((hash) => {
        const also_added = (hashToAdd[type] || new Set()).has(hash)
        if (also_added) {
          hashToAdd[type].delete(hash)
        }
        return !also_added
      })
    )

    // 3.c On retire les cles restantes de la memoire.
    // --------------------------------------------------
    hashToDelete[type].forEach((hash) => {
      memory[type].delete(hash)
    })
  })
}
