
f = this;
actual_batch = "2002_1";
date_debut = new Date("2014-01-01")
date_fin = new Date("2016-01-01")
serie_periode = f.generatePeriodSerie(date_debut, date_fin);
includes = {"all": true}
offset_effectif = 2

const notreMap = (testData) => {
  const results = [];
  emit = (key, value) => results.push({"_id": key, value});
  testData.forEach(entrepriseOuEtablissement => map.call(entrepriseOuEtablissement)); // will call emit an inderminate number of times
  // testData contains _id and value properties. testData is passed as this
  return results;
};

print(JSON.stringify(notreMap(testData), null, 2));
