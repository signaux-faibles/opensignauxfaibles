function apart (v, output_indexed) {
  var apart = Object.keys(v.apdemande).reduce((apart, hash) => {
    apart[v.apdemande[hash].id_demande.substring(0, 9)] = {
      "demande": hash,
      "consommation": []
    }
    return apart
  }, {})

  Object.keys(v.apdemande).forEach(hash => {
    var periode_deb = v.apdemande[hash].periode.start
    var periode_fin = v.apdemande[hash].periode.end
    var periode_deb_floor = new Date(Date.UTC(periode_deb.getUTCFullYear(), periode_deb.getUTCMonth(), 1, 0, 0, 0, 0))
    var periode_fin_ceil = new Date(Date.UTC(periode_fin.getUTCFullYear(), periode_fin.getUTCMonth() + 1, 1, 0, 0, 0, 0))
    var series = generatePeriodSerie(periode_deb_floor, periode_fin_ceil)
    series.forEach( date => {
      let time = date.getTime()
      if (time in output_indexed){
        output_indexed[time].apart_heures_autorisees = v.apdemande[hash].hta
      }  
    })
  })

  Object.keys(v.apconso).forEach(hash => {
    var valueap = v.apconso[hash]
    if (valueap.id_conso.substring(0, 9) in apart) {
      apart[valueap.id_conso.substring(0, 9)].consommation.push(hash)
    }
  })

  // relier apdemande et apconso
  Object.keys(apart).forEach(k => {
    v.apdemande[apart[k].demande].hash_consommation = apart[k].consommation
    for (let j in apart[k].consommation) {
      v.apconso[apart[k].consommation[j]].hash_demande = apart[k].demande
    }
  })

  Object.keys(v.apconso).forEach(h => {
    var conso = v.apconso[h]
    if (conso.hash_demande) {
      var time = conso.periode.getTime()
      if (time in output_indexed){
        output_indexed[time].apart_heures_consommees = output_indexed[time].apart_heures_consommees + conso.heure_consomme
        output_indexed[time].apart_motif_recours = v.apdemande[conso.hash_demande].motif_recours_se
      }
    }
  })

  output_array.forEach(val => {
    if (val !== null && val.effectif > 0){
      val.ratio_apart = val.apart_heures_consommees / (val.effectif * 157.67)
    }
  })
}