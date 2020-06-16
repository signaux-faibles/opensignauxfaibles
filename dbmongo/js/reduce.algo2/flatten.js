function flatten(v, actual_batch) {
  "use strict"
  var res = Object.keys(v.batch || {})
    .sort()
    .filter((batch) => batch <= actual_batch)
    .reduce(
      (m, batch) => {
        // Types intéressants = nouveaux types, ou types avec suppressions
        var delete_types = Object.keys(
          (v.batch[batch].compact || {}).delete || {}
        )
        var new_types = Object.keys(v.batch[batch])
        var all_interesting_types = [
          ...new Set([...delete_types, ...new_types]),
        ]

        all_interesting_types.forEach((type) => {
          m[type] = m[type] || {}
          // On supprime les clés qu'il faut
          if (
            v.batch[batch] &&
            v.batch[batch].compact &&
            v.batch[batch].compact.delete &&
            v.batch[batch].compact.delete[type] &&
            v.batch[batch].compact.delete[type] != {}
          ) {
            v.batch[batch].compact.delete[type].forEach((hash) => {
              delete m[type][hash]
            })
          }
          Object.assign(m[type], v.batch[batch][type])
        })
        return m
      },
      { key: v.key, scope: v.scope }
    )

  return res
}

exports.flatten = flatten
