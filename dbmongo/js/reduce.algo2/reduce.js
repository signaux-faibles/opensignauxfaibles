function reduce(key, values) {
  try {
    return values.reduce((val, accu) => {
      return Object.assign(accu, val)
    }, {})
  } catch {
        print("My name is " + key + " and I died in reduce.algo2/reduce.js")
  }
}
