// initially available as /dbmongo/js/test/public/lib_public.js

function ISODate(date) {
  let d = new Date(date)
  return d
}

f = {
  altaresToHuman,
  apconso,
  apdemande,
  bdf,
  compareDebit,
  cotisations,
  dateAddDay,
  dateAddMonth,
  dealWithProcols,
  debits,
  delai,
  diane,
  effectifs,
  finalize,
  flatten,
  generatePeriodSerie,
  idEntreprise,
  iterable,
  map,
  procolToHuman,
  reduce,
  sirene,
}


var pool = {}

function emit(key, value) {
  id = JSON.stringify(key)
  pool[id] = (pool[id] || []).concat([{key, value}])
}

function reducer(array, reduce) {
  if (array.length == 1) {
    return array[0]
  } else {
    let newVal = reduce(array[0].key, [array[0].value, array[1].value])
    return reducer([newVal].concat(array.slice(2, array.length)), reduce)
  }
}

function invertedReducer(array, reduce) {
  if (array.length == 1) {
    return array[0]
  } else {
    let newVal = reduce(array[0].key, [array[array.length-1].value, array[array.length-2].value])
    return reducer([newVal].concat(array.slice(0, array.length-2)), reduce)
  }
}
