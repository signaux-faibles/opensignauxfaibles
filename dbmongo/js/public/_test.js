jsParams = {}
jsParams.actual_batch = "1905"
date_debut = new Date("2014-01-01")
date_fin = new Date("2018-02-01")
jsParams.serie_periode = f.generatePeriodSerie(new Date("2014-01-01"), new Date("2018-02-01"))
offset_effectif = 2

objects.forEach(object => {
  f.value = object.value
  f._id = object._id
  f.map()
})

var intermediateResult = Object.values(pool).map(array => ({
  key: array[0].key,
  value: reducer(array, f.reduce)
}))

var invertedIntermediateResult = Object.values(pool).map(array => ({
  key: array[0].key,
  value: invertedReducer(array, f.reduce)
}))

var result = intermediateResult.map(r => ({
  _id: r.key,
  value: f.finalize(r.key, r.value)
}))

var invertedResult = invertedIntermediateResult.map(r => ({
  _id: r.key,
  value: f.finalize(r.key, r.value)
}))

print(JSON.stringify(result) == JSON.stringify(invertedResult))
