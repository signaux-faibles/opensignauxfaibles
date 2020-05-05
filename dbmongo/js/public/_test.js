actual_batch = "1905"
date_debut = new Date("2014-01-01")
date_fin = new Date("2018-02-01")
serie_periode = f.generatePeriodSerie(new Date("2014-01-01"), new Date("2018-02-01"))
offset_effectif = 2

objects.forEach(object => {
  f.value = object.value
  f._id = object._id
  f.map()
})

var intermediateResult = []
Object.keys(pool).forEach(k => {
  array = pool[k]
  intermediateResult.push({key: pool[k][0].key, value: reducer(array, f.reduce)})
})

var invertedIntermediateResult = []
Object.keys(pool).forEach(k => {
  array = pool[k]
  invertedIntermediateResult.push({key: pool[k][0].key, value: invertedReducer(array, f.reduce)})
})

var result = []
intermediateResult.forEach(r => {
  result.push({_id: r.key, value: f.finalize(r.key, r.value)})
})

var invertedResult = []
invertedIntermediateResult.forEach(r => {
  invertedResult.push({_id: r.key, value: f.finalize(r.key, r.value)})
})

print(JSON.stringify(result) == JSON.stringify(invertedResult))
