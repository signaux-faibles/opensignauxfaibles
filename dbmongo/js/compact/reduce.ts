import "../globals.ts"
import * as f from "./currentState"

// Entrée: données d'entreprises venant de ImportedData, regroupées par entreprise ou établissement.
// Sortie: un objet fusionné par entreprise ou établissement, contenant les données historiques et les données importées, à destination de la collection RawData.
// Opérations: retrait des données doublons et application des corrections de données éventuelles.
export function reduce(
  key: SiretOrSiren,
  values: CompanyDataValues[]
): CompanyDataValues {
  "use strict"

  // Appelle fct() pour chaque propriété définie (non undefined) de obj.
  // Contrat: obj ne doit contenir que les clés définies dans son type.
  const forEachPopulatedProp = <T>(
    obj: T,
    fct: (key: keyof T, val: MakeRequired<T>[keyof T]) => unknown
  ) =>
    (Object.keys(obj) as Array<keyof T>).forEach((key) => {
      if (obj[key] !== undefined) fct(key, obj[key])
    })

  // Tester si plusieurs batchs. Reduce complet uniquement si plusieurs
  // batchs. Sinon, juste fusion des attributs
  const auxBatchSet = new Set()

  const severalBatches = values.some((value) => {
    auxBatchSet.add(Object.keys(value.batch || {}))
    return auxBatchSet.size > 1
  })

  //fusion des attributs dans values
  const reduced_value: CompanyDataValues = values.reduce(
    (m, value: CompanyDataValues) => {
      Object.keys(value.batch).forEach((batch) => {
        type DataType = keyof BatchValue
        m.batch[batch] = (Object.keys(value.batch[batch]) as DataType[]).reduce(
          (batchValues: BatchValue, type: DataType) => ({
            ...batchValues,
            [type]: value.batch[batch][type],
          }),
          m.batch[batch] || {}
        )
      })
      return m
    },
    { key: key, scope: values[0].scope, batch: {} }
  )

  // Cette fonction reduce() est appelée à deux moments:
  // 1. agregation par établissement d'objets ImportedData. Dans cet étape, on
  // ne travaille généralement que sur un seul batch.
  // 2. agregation de ces résultats au sein de RawData, en fusionnant avec les
  // données potentiellement présentes. Dans cette étape, on fusionne
  // généralement les données de plusieurs batches. (données historiques)

  if (!severalBatches) return reduced_value

  //////////////////////////////////////////////////
  // ETAPES DE LA FUSION AVEC DONNÉES HISTORIQUES //
  //////////////////////////////////////////////////

  // 0. On calcule la memoire au moment du batch à modifier
  const memory_batches: BatchValue[] = Object.keys(reduced_value.batch)
    .filter((batch) => batch < batchKey)
    .sort()
    .reduce((m: BatchValue[], batch: string) => {
      m.push(reduced_value.batch[batch])
      return m
    }, [])

  const memory = f.currentState(memory_batches)

  // Pour tous les batchs à modifier, c'est-à-dire le batch ajouté et tous les
  // suivants.
  const modified_batches = batches.filter((batch) => batch >= batchKey)

  modified_batches.forEach((batch: string) => {
    reduced_value.batch[batch] = reduced_value.batch[batch] || {}

    // Les types où il y a potentiellement des suppressions
    const stock_types = completeTypes[batch].filter(
      (type) => (memory[type] || new Set()).size > 0
    )

    // 1. On recupère les cles ajoutes et les cles supprimes
    // -----------------------------------------------------

    const hashToDelete: { [dataType: string]: Set<DataHash> } = {}
    type DataType = string
    const hashToAdd: Record<DataType, Set<DataHash>> = {}

    // Itération sur les types qui ont potentiellement subi des modifications
    // pour compléter hashToDelete et hashToAdd.
    // Les suppressions de types complets / stock sont gérés dans le bloc suivant.
    const currentBatch = reduced_value.batch[batch]
    for (const type in currentBatch) {
      // Le type compact gère les clés supprimées
      // Ce type compact existe si le batch en cours a déjà été compacté.
      if (type === "compact") {
        const compactDelete = reduced_value.batch[batch].compact?.delete
        if (compactDelete) {
          forEachPopulatedProp(compactDelete, (delete_type, keysToDelete) => {
            keysToDelete.forEach((hash) => {
              hashToDelete[delete_type] = hashToDelete[delete_type] || new Set()
              hashToDelete[delete_type].add(hash)
            })
          })
        }
      } else {
        for (const hash in currentBatch[type as keyof BatchValue]) {
          hashToAdd[type] = hashToAdd[type] || new Set()
          hashToAdd[type].add(hash)
        }
      }
    }

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
        const reducedBatch = reduced_value.batch[batch]
        reducedBatch.compact = reducedBatch.compact || { delete: {} }
        reducedBatch.compact.delete = reducedBatch.compact.delete || {}
        reducedBatch.compact.delete[type] = [...hashToDelete[type]]
      }
    })

    // filtrage des données en fonction de new_types et hashToAdd
    type AllValueTypesButCompact = Exclude<keyof BatchValue, "compact">
    const typesToAdd = Object.keys(hashToAdd) as AllValueTypesButCompact[]
    typesToAdd.forEach((type) => {
      reduced_value.batch[batch][type] = [...hashToAdd[type]].reduce(
        (typedBatchValues, hash) => ({
          ...typedBatchValues,
          [hash]: reduced_value.batch[batch][type]?.[hash],
        }),
        {}
      )
    })

    // 6. nettoyage
    // ------------

    if (reduced_value.batch[batch]) {
      //types vides
      Object.keys(reduced_value.batch[batch]).forEach((strType) => {
        const type = strType as keyof BatchValue
        if (Object.keys(reduced_value.batch[batch][type] || {}).length === 0) {
          delete reduced_value.batch[batch][type]
        }
      })
      //hash à supprimer vides (compact.delete)
      const compactDelete = reduced_value.batch[batch].compact?.delete
      if (compactDelete) {
        forEachPopulatedProp(compactDelete, (type, keysToDelete) => {
          if (keysToDelete.length === 0) {
            delete compactDelete[type]
          }
        })
        if (Object.keys(compactDelete).length === 0) {
          delete reduced_value.batch[batch].compact
        }
      }
      //batchs vides
      if (Object.keys(reduced_value.batch[batch]).length === 0) {
        delete reduced_value.batch[batch]
      }
    }
  })

  return reduced_value
}
