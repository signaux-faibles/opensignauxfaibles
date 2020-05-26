function generatePeriodSerie (date_debut, date_fin) {
  "use strict";
  var date_next = new Date(date_debut.getTime())
  var serie = []
  while (date_next.getTime() < date_fin.getTime()) {
    serie.push(new Date(date_next.getTime()))
    date_next.setUTCMonth(date_next.getUTCMonth() + 1)
  }
  return serie
}
exports.generatePeriodSerie = generatePeriodSerie
