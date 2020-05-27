import "../globals.ts"

// currentState() agrège un ensemble de batch, en tenant compte des suppressions
// pour renvoyer le dernier état connu des données.
// Note: similaire à flatten() de reduce.algo2.
export function currentState(batches: BatchValue[]): CurrentDataState {
  "use strict"
  const currentState: CurrentDataState = batches.reduce(
    (m: CurrentDataState, batch: BatchValue) => {
      //1. On supprime les clés de la mémoire
      Object.keys((batch.compact || { delete: [] }).delete).forEach((type) => {
        batch.compact.delete[type].forEach((key) => {
          m[type].delete(key) // Should never fail or collection is corrupted
        })
      })

      //2. On ajoute les nouvelles clés
      Object.keys(batch)
        .filter((type) => type !== "compact")
        .forEach((type: keyof BatchValue) => {
          m[type] = m[type] || new Set()

          Object.keys(batch[type]).forEach((key) => {
            m[type].add(key)
          })
        })
      return m
    },
    {}
  )

  return currentState
}
