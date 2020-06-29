import "../globals.ts"

import { listHashesToAddAndDelete } from "./listHashesToAddAndDelete"

/**
 * Appelée par reduce(), compactBatch() va générer un diff entre les
 * données de batch et les données précédentes fournies par memory.
 * Pré-requis: les batches précédents doivent avoir été compactés.
 */
export function compactBatch(
  currentBatch: BatchValue,
  memory: CurrentDataState,
  batchKey: string
): BatchValue {
  // Les types où il y a potentiellement des suppressions
  const stock_types = completeTypes[batchKey].filter(
    (type) => (memory[type] || new Set()).size > 0
  )

  const { hashToAdd, hashToDelete } = listHashesToAddAndDelete(currentBatch)

  //
  // 2. On ajoute aux cles supprimees les types stocks de la memoire.
  // ----------------------------------------------------------------

  stock_types.forEach((type) => {
    hashToDelete[type] = new Set([
      ...(hashToDelete[type] || new Set()),
      ...memory[type],
    ])
  })

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

  Object.keys(hashToAdd).forEach((type) => {
    // 4.a Pour chaque cle ajoutee: est-ce qu'elle est dans la memoire ? Si oui on filtre cette cle
    // i.e. on herite de la memoire. (pas de maj de la memoire)
    // ---------------------------------------------------------------------------------------------
    hashToAdd[type] = new Set(
      [...hashToAdd[type]].filter((hash) => {
        return !(memory[type] || new Set()).has(hash)
      })
    )

    // 4.b Pour chaque cle ajoutee restante: on ajoute à la memoire.
    // -------------------------------------------------------------

    hashToAdd[type].forEach((hash) => {
      memory[type] = memory[type] || new Set()
      memory[type].add(hash)
    })
  })

  // 5. On met à jour reduced_value
  // -------------------------------
  stock_types.forEach((type) => {
    if (hashToDelete[type]) {
      currentBatch.compact = currentBatch.compact || { delete: {} }
      currentBatch.compact.delete = currentBatch.compact.delete || {}
      currentBatch.compact.delete[type] = [...hashToDelete[type]]
    }
  })

  // filtrage des données en fonction de new_types et hashToAdd
  type AllValueTypesButCompact = Exclude<keyof BatchValue, "compact">
  const typesToAdd = Object.keys(hashToAdd) as AllValueTypesButCompact[]
  typesToAdd.forEach((type) => {
    currentBatch[type] = [...hashToAdd[type]].reduce(
      (typedBatchValues, hash) => ({
        ...typedBatchValues,
        [hash]: currentBatch[type]?.[hash],
      }),
      {}
    )
  })

  // 6. nettoyage
  // ------------

  if (currentBatch) {
    //types vides
    Object.keys(currentBatch).forEach((strType) => {
      const type = strType as keyof BatchValue
      if (Object.keys(currentBatch[type] || {}).length === 0) {
        delete currentBatch[type]
      }
    })
    //hash à supprimer vides (compact.delete)
    const compactDelete = currentBatch.compact?.delete
    if (compactDelete) {
      Object.keys(compactDelete).forEach((type) => {
        if (compactDelete[type].length === 0) {
          delete compactDelete[type]
        }
      })
      if (Object.keys(compactDelete).length === 0) {
        delete currentBatch.compact
      }
    }
  }
  return currentBatch
}
