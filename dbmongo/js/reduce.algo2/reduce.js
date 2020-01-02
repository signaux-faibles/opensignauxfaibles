function reduce(key, values) {
  return values.reduce((val, accu) => {
    return Object.assign(accu, val)
  }, {})
}
