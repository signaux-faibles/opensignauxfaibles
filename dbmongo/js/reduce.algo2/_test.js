"use strict";

const jsParams = this; // => all properties of this object will become global. TODO: remove this when merging namespace (https://github.com/signaux-faibles/opensignauxfaibles/pull/40)
jsParams.actual_batch = "1905"
jsParams.date_debut = new Date("2014-01-01")
jsParams.date_fin = new Date("2018-02-01")
jsParams.serie_periode = f.generatePeriodSerie(new Date("2014-01-01"), new Date("2018-02-01"))
jsParams.offset_effectif = 2
jsParams.includes = { all: true }

objects.forEach(object => {
  f.value = object.value
  f._id = object._id
  f.map()
})

var intermediateResult = Object.values(pool).map(array => reducer(array, f.reduce))

var invertedIntermediateResult = Object.values(pool).map(array => invertedReducer(array, f.reduce))

var result = intermediateResult.map(r => f.finalize(null, r))

var invertedResult = invertedIntermediateResult.map(r => f.finalize(null, r))

print(JSON.stringify(sortObject(result)) == JSON.stringify(sortObject(invertedResult)))

// from https://gist.github.com/ninapavlich/1697bcc107052f5b884a794d307845fe
function sortObject(object) {
  if (!object) {
    return object;
  }

  const isArray = object instanceof Array;
  var sortedObj = {};
  if (isArray) {
    sortedObj = object.map((item) => sortObject(item));
  } else {
    var keys = Object.keys(object);
    // console.log(keys);
    keys.sort(function(key1, key2) {
      (key1 = key1.toLowerCase()), (key2 = key2.toLowerCase());
      if (key1 < key2) return -1;
      if (key1 > key2) return 1;
      return 0;
    });

    for (var index in keys) {
      var key = keys[index];
      if (typeof object[key] == 'object') {
        sortedObj[key] = sortObject(object[key]);
      } else {
        sortedObj[key] = object[key];
      }
    }
  }

  return sortedObj;
}