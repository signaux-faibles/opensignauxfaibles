export function generatePeriodSerie(date_debut: Date, date_fin: Date): Date[] {
  "use strict"
  const date_next = new Date(date_debut.getTime())
  const serie = []
  while (date_next.getTime() < date_fin.getTime()) {
    serie.push(new Date(date_next.getTime()))
    date_next.setUTCMonth(date_next.getUTCMonth() + 1)
  }
  return serie
}
