function flatten(v, actual_batch) {
  return Object.keys((v.batch || {})).sort().filter(batch => batch <= actual_batch).reduce((m, batch) => {
      Object.keys(v.batch[batch]).forEach((type) => {
          m[type] = (m[type] || {})
          var  array_delete = (v.batch[batch].compact.delete[type]||[])
          if (array_delete != {}) {array_delete.forEach(hash => {
              delete m[type][hash]
          })
          }
          Object.assign(m[type], v.batch[batch][type])
      })
      return m
  }, { "key": v.key, scope: v.scope })
}