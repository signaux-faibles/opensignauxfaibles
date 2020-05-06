function cotisations(vcotisation) {
  "use strict";
  var offset_cotisation = 0 
  var value_cotisation = {}
  
  // Répartition des cotisations sur toute la période qu'elle concerne
  vcotisation = vcotisation || {}
  Object.keys(vcotisation).forEach(function (h) {
    var cotisation = vcotisation[h]
    var periode_cotisation = f.generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
    periode_cotisation.forEach(date_cotisation => {
      let date_offset = f.dateAddMonth(date_cotisation, offset_cotisation)
      value_cotisation[date_offset.getTime()] = (value_cotisation[date_offset.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
    })
  })

  var output_cotisation = []

  serie_periode.forEach(p => {
    output_cotisation.push(
      (value_cotisation[p.getTime()] || []) 
        .reduce((m,c) => m+c, 0)
    )
  })

  return(output_cotisation)
}
