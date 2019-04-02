function reduce(key, values) {
  if (key.scope="entreprise") {
    values = values.reduce((m, v) => {
      if (v.sirets) {
        m.sirets = (m.sirets || []).concat(v.sirets)
        delete v.sirets
      }
      Object.assign(m, v)
      return m
    }, {})
  }
  return values
}