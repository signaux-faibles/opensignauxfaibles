function reducer(array, reduce) {
  if (array.length === 1) {
    return array[0]
  } else {
    const newVal = reduce(array[0].key, [array[0].value, array[1].value])
    return reducer([newVal].concat(array.slice(2, array.length)), reduce)
  }
}

function invertedReducer(array, reduce) {
  if (array.length === 1) {
    return array[0]
  } else {
    const newVal = reduce(array[0].key, [
      array[array.length - 1].value,
      array[array.length - 2].value,
    ])
    return reducer([newVal].concat(array.slice(0, array.length - 2)), reduce)
  }
}

exports.reducer = reducer
exports.invertedReducer = invertedReducer
