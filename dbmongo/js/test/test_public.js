
actual_batch = "1905"
date_debut = new Date("2014-01-01")
date_fin = new Date("2018-02-01")
serie_periode = f.generatePeriodSerie(new Date("2014-01-01"), new Date("2018-02-01"))

f.value = object.value
f._id = object._id

f.map()

print(JSON.stringify(pool, null, 2))