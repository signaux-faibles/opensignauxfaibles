import { BatchValue, DataType, Scope, SiretOrSiren } from "../RawDataTypes"
import { CompanyDataValuesWithCompact } from "../compact/applyPatchesToBatch"

export type FlattenedImportedData = {
  key: SiretOrSiren
  scope: Scope
} & BatchValue

/**
 * Appelé par `map()`, `flatten()` transforme les données importées (*Batches*)
 * d'une entreprise ou établissement afin de retourner un unique objet *plat*
 * contenant les valeurs finales de chaque type de données.
 *
 * Pour cela:
 * - il supprime les clés `compact.delete` des *Batches* en entrées;
 * - il agrège les propriétés apportées par chaque *Batch*, dans l'ordre chrono.
 */
export function flatten(
  v: CompanyDataValuesWithCompact,
  actual_batch: string
): FlattenedImportedData {
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
        ] as DataType[]

        all_interesting_types.forEach((type) => {
          const typedData = m[type]
          if (typeof typedData === "object") {
            // On supprime les clés qu'il faut
            const keysToDelete = v.batch[batch]?.compact?.delete?.[type] || []
            for (const hash of keysToDelete) {
              delete typedData[hash]
            }
          } else {
            m[type] = {}
          }
          Object.assign(m[type], v.batch[batch][type])
        })
        return m
      },
      { key: v.key, scope: v.scope } as FlattenedImportedData
    )

  return res
}
