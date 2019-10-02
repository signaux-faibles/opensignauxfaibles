function ISODate(date) {
  let d = new Date(date)
  return d
}

var pool = {}

function emit(key, value) {
  id = key.siren + key.batch + key.periode.getTime()
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

f = {
  add,
  altaresToHuman,
  apart,
  ccsf,
  cibleApprentissage,
  compareDebit,
  compte,
  cotisation,
  cotisationsdettes,
  dateAddMonth,
  dealWithProcols,
  defaillances,
  delais,
  detteFiscale,
  effectifs,
  finalize,
  financierCourtTerme,
  flatten,
  fraisFinancier,
  generatePeriodSerie,
  interim,
  lookAhead,
  map,
  naf,
  outputs,
  poidsFrng,
  procolToHuman,
  reduce,
  repeatable,
  sirene,
  sirene_ul,
  tauxMarge,
}
