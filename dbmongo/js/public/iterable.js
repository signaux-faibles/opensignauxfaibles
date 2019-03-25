function iterable(dict) {
  try {
    return Object.keys(dict).map(h => {
    return dict[h]
  })
  } catch(error) {
    return []
  }
}