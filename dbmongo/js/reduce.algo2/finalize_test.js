"use strict";

// Context: this golden-file-based test runner was designed to prevent
// regressions on the JS functions (common + algo2) used to compute the
// "Features" collection from the "RawData" collection.
//
// It requires the JS functions from common + algo2 (notably: map()),
// and a makeTestData() function to generate a realistic test data set.
//
// Please execute ../test/test_finalize_algo2.sh
// to fill these requirements and run the tests.

// Allow f.*() function calls to resolve to globally-defined functions
const f = this;

// Define global parameters that are required by JS functions
const jsParams = this; // => all properties of this object will become global. TODO: remove this when merging namespace (https://github.com/signaux-faibles/opensignauxfaibles/pull/40)
jsParams.actual_batch = "2002_1";
jsParams.date_debut = new Date("2014-01-01");
jsParams.date_fin = new Date("2016-01-01");
jsParams.serie_periode = f.generatePeriodSerie(date_debut, date_fin);
jsParams.includes = { "all": true };
jsParams.offset_effectif = 2;

let emit; // global emit() function that mapFct() will call

// Run a map() function designed for MongoDB, i.e. that calls emit() an
// inderminate number of times, instead of returning one value per iteration.
function runMongoMap (testData, mapFct) {
  const results = []; // holds all the { _id, value } objects emitted from mapFct()
  // define a emit() function that mapFct() can call
  emit = (key, value) => results.push({"_id": key, value});
  testData.forEach(entrepriseOuEtablissement => mapFct.call(entrepriseOuEtablissement)); // entrepriseOuEtablissement will be accessible through `this`, in mapFct()
  return results;
};

// Generate a realistic test data set
const testData = makeTestData({
  ISODate: (dateString) => new Date(dateString.replace('+0000', '+00:00')), // make sure that timezone format complies with the spec
  NumberInt: (int) => int,
});

// Print the output of the global map() function
var map_result = runMongoMap(testData, map); // -> [ { _id, value } ]
const keys = map_result.map(entrepriseOuEtablissement => entrepriseOuEtablissement._id);
const values = map_result.map(entrepriseOuEtablissement => entrepriseOuEtablissement.value);
var reduce_result = f.reduce(keys, values); // -> { }
var finalizeResult = f.finalize(Object.keys(reduce_result), Object.values(reduce_result));
