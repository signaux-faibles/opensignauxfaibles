
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
  intermediateResult.push(reducer(array, f.reduce))
})

var result = []
intermediateResult.forEach(r => {
  result.push(f.finalize(null, r))
})

print(JSON.stringify(result, null, 2))