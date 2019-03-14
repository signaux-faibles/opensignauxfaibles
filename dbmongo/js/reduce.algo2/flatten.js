function flatten(v, actual_batch) {
  return Object.keys(v.batch || {})
    .sort()
    .filter(batch => batch <= actual_batch)
    .reduce((m, batch) => {
      Object.keys(v.batch[batch]).forEach((type) => {
        m[type] = (m[type] || {})
        // On supprime les clés qu'il faut
        if (v.batch[batch] && v.batch[batch].compact && v.batch[batch].compact.delete &&
          v.batch[batch].compact.delete[type] && v.batch[batch].compact.delete[type] != {}) {

          v.batch[batch].compact.delete[type].forEach(hash => {
            delete m[type][hash]
          })
        }
        Object.assign(m[type], v.batch[batch][type])
      })
      return m
    }, { "key": v.key, scope: v.scope })
}
