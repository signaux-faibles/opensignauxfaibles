"use strict"

// common fakes from /dbmongo/js/test/public/lib_public.js
//               and /dbmongo/js/test/algo2/lib_algo2.js

Object.bsonsize = function (obj) {
  return JSON.stringify(obj).length
}

function ISODate(date) {
  return new Date(date)
}

const f = this /* = {
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
  ...
}*/

const exports = { f }
const globalThis = { f } // permet d'importer la variable globale f depuis les tests TypeScript

function reducer(array, reduce) {
  if (array.length == 1) {
    return array[0]
  } else {
    const newVal = reduce(array[0].key, [array[0].value, array[1].value])
    return reducer([newVal].concat(array.slice(2, array.length)), reduce)
  }
}

function invertedReducer(array, reduce) {
  if (array.length == 1) {
    return array[0]
  } else {
    const newVal = reduce(array[0].key, [
      array[array.length - 1].value,
      array[array.length - 2].value,
    ])
    return reducer([newVal].concat(array.slice(0, array.length - 2)), reduce)
  }
}
