import "../globals.ts"

// currentState() agrège un ensemble de batch, en tenant compte des suppressions
// pour renvoyer le dernier état connu des données.
// Note: similaire à flatten() de reduce.algo2.
export function currentState(batches: BatchValue[]): CurrentDataState {
  "use strict"

  // Retourne les clés de obj, en respectant le type défini dans le type de obj.
  // Contrat: obj ne doit contenir que les clés définies dans son type.
  const typedObjectKeys = <T>(obj: T): Array<keyof T> =>
    Object.keys(obj) as Array<keyof T>

  // Appelle fct() pour chaque propriété définie (non undefined) de obj.
  // Contrat: obj ne doit contenir que les clés définies dans son type.
  const forEachPopulatedProp = <T>(
    obj: T,
    fct: (key: keyof T, val: Required<T>[keyof T]) => unknown
  ) =>
    (Object.keys(obj) as Array<keyof T>).forEach((key) => {
      if (obj[key] !== undefined) fct(key, obj[key])
    })

  const currentState: CurrentDataState = batches.reduce(
    (m: CurrentDataState, batch: BatchValue) => {
      //1. On supprime les clés de la mémoire
      if (batch.compact) {
        forEachPopulatedProp(batch.compact.delete, (type, keysToDelete) => {
          keysToDelete.forEach((key) => {
            m[type].delete(key) // Should never fail or collection is corrupted
          })
        })
      }

      //2. On ajoute les nouvelles clés
      for (const type of typedObjectKeys(batch)) {
        if (type === "compact") continue
        m[type] = m[type] || new Set()
        for (const key in batch[type]) {
          m[type].add(key)
        }
      }
      return m
    },
    {}
  )

  return currentState
}
