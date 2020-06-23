import "../globals.ts"

// currentState() agrège un ensemble de batch, en tenant compte des suppressions
// pour renvoyer le dernier état connu des données.
// Note: similaire à flatten() de reduce.algo2.
export function currentState(batches: BatchValue[]): CurrentDataState {
  "use strict"
  const currentState: CurrentDataState = batches.reduce(
    (m: CurrentDataState, batch: BatchValue) => {
      //1. On supprime les clés de la mémoire
      if (batch.compact) {
        for (const type of Object.keys(batch.compact.delete)) {
          batch.compact.delete[type].forEach((key) => {
            m[type].delete(key) // Should never fail or collection is corrupted
          })
        }
      }

      //2. On ajoute les nouvelles clés
      for (const type in batch) {
        if (type === "compact") continue
        m[type] = m[type] || new Set()
        for (const key in batch[type as keyof Exclude<BatchValue, "compact">]) {
          // note: nous ne serions pas prévenus si `compact` était défini dans `batch`
          m[type].add(key)
        }
      }
      return m
    },
    {}
  )

  return currentState
}
