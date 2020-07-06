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
