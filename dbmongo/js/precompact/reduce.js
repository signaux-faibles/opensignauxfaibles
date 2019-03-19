function reduce(key, values) {
  //fusion des objets dans values
  let reduced_value = values.reduce((m, value) => {
    Object.keys((value.batch||{})).forEach(batch => {
      m.batch = (m.batch||{})
      m.batch[batch] = (m.batch[batch] || {})
      Object.keys(value.batch[batch]).forEach(type => {
        m.batch[batch][type] = (m.batch[batch][type] || {})
        Object.assign(m.batch[batch][type],value.batch[batch][type])
      })
    })
    return m
  }, {"key": key, "scope": values[0].scope  })

  return(reduced_value)
}
