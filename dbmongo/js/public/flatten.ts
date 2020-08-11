export type FlattenedImportedData = {
  key: SiretOrSiren
  scope: Scope
} & Partial<BatchValue>

// Note: cette fonction a été copiée depuis reduce.algo2/flatten.ts
// TODO: déplacer reduce.algo2/flatten.ts dans common/ puis la réutiliser aux deux endroits
export function flatten(
  v: CompanyDataValues,
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
      { key: v.key, scope: v.scope } as FlattenedImportedData
    )

  return res
}
