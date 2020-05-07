// Context: this golden-file-based test runner was designed to prevent
// regressions on the JS functions (common + algo2) used to compute the
// "Features" collection from the "RawData" collection.
// 
// It requires the JS functions from common + algo2 (notably: map()),
// and a makeTestData() function to generate a realistic test data set.
//
// => Please execute ../test/test_map_reduce_algo2.sh to fill these
// requirements and run the tests.

// Allow f.*() function calls to resolve to globally-defined functions 
f = this;

// Define global parameters that are required by JS functions
jsParams = {}
jsParams.naf = naf; // loaded from data/naf.js
jsParams.actual_batch = "2002_1";
jsParams.date_debut = new Date("2014-01-01");
jsParams.date_fin = new Date("2016-01-01");
jsParams.serie_periode = f.generatePeriodSerie(jsParams.date_debut, jsParams.date_fin);
jsParams.includes = { "all": true };
jsParams.offset_effectif = 2;
Object.freeze(jsParams);

// Run a map() function designed for MongoDB, i.e. that calls emit() an
// inderminate number of times, instead of returning one value per iteration.
function runMongoMap (testData, mapFct) {
  const results = []; // holds all the { _id, value } objects emitted from mapFct()
  emit = (key, value) => results.push({"_id": key, value}); // define a emit() function that mapFct() can call
  testData.forEach(entrepriseOuEtablissement => mapFct.call(entrepriseOuEtablissement)); // entrepriseOuEtablissement will be accessible through `this`, in mapFct()
  return results;
};

// Generate a realistic test data set
const testData = makeTestData({
  ISODate: (dateString) => new Date(dateString.replace('+0000', '+00:00')), // make sure that timezone format complies with the spec
  NumberInt: (int) => int,
});

// Print the output of the global map() function
print(JSON.stringify(runMongoMap(testData, map), null, 2));
