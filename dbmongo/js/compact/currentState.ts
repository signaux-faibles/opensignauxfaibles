export function currentState(batches) {
  "use strict"
  type Keys = Set<string>
  type State = { [key: string]: Keys }
  const currentState: State = batches.reduce((m: State, batch) => {
    //1. On supprime les clés de la mémoire
    Object.keys((batch.compact || { delete: [] }).delete).forEach((type) => {
      batch.compact.delete[type].forEach((key) => {
        m[type].delete(key) // Should never fail or collection is corrupted
      })
    })

    //2. On ajoute les nouvelles clés
    Object.keys(batch)
      .filter((type) => type !== "compact")
      .forEach((type) => {
        m[type] = m[type] || new Set()

        Object.keys(batch[type]).forEach((key) => {
          m[type].add(key)
        })
      })
    return m
  }, {})

  return currentState
}
