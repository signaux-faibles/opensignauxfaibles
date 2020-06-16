type V = {
  key: unknown
  scope: Scope
  batch: BatchValues
}

type Flattened = Partial<BatchValue> & {
  key: unknown
  scope: Scope
}

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
          if (
            v.batch[batch] &&
            v.batch[batch].compact &&
            v.batch[batch].compact.delete &&
            v.batch[batch].compact.delete[type]
          ) {
            v.batch[batch].compact.delete[type].forEach((hash) => {
              if (typeof m[type] === "object" && (m[type] as any)[hash])
                delete (m[type] as any)[hash]
            })
          }
          Object.assign(m[type], v.batch[batch][type])
        })
        return m
      },
      { key: v.key, scope: v.scope } as Flattened
    )

  return res
}
