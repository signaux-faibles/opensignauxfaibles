export type V = {
  // TODO: donner un nom plus explicite au type
  key: SiretOrSiren
  scope: Scope
  batch: BatchValues
}

export type Flattened = {
  key: SiretOrSiren
  scope: Scope
} & Partial<BatchValue>

/**
 * Appelé par `map()`, `flatten()` transforme les données importées (*Batches*)
 * d'une entreprise ou établissement afin de retourner un unique objet *plat*
 * contenant les valeurs finales de chaque type de données.
 *
 * Pour cela:
 * - il supprime les clés `compact.delete` des *Batches* en entrées;
 * - il agrège les propriétés apportées par chaque *Batch*, dans l'ordre chrono.
 */
export function flatten(v: V, actual_batch: string): Flattened {
  "use strict"
  const res = Object.keys(v.batch || {})
    .sort()
    .filter((batch) => batch <= actual_batch)
    .reduce(
      (m, batch) => {
        // Types intéressants = nouveaux types, ou types avec suppressions
        const delete_types = Object.keys(
          (v.batch[batch].compact || {}).delete || {}
        )
        const new_types = Object.keys(v.batch[batch])
        const all_interesting_types = [
          ...new Set([...delete_types, ...new_types]),
        ] as BatchDataType[]

        all_interesting_types.forEach((type) => {
          m[type] = m[type] || ({} as any)
          // On supprime les clés qu'il faut
          const batchData = v.batch[batch]
          const keysToDelete = batchData?.compact?.delete?.[type] || []
          for (const hash of keysToDelete) {
            if (typeof m[type] === "object" && (m[type] as any)[hash])
              delete (m[type] as any)[hash]
          }
          Object.assign(m[type], v.batch[batch][type])
        })
        return m
      },
      { key: v.key, scope: v.scope } as Flattened
    )

  return res
}
