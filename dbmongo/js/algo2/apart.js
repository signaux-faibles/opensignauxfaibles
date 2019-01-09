function apart (v, output_indexed) {
  var apart = Object.keys(v.apdemande).reduce((apart, hash) => {
    apart[v.apdemande[hash].id_demande.substring(0, 10)] = {
        "demande": hash,
        "consommation": []
    }
    return apart
  }, {})



  Object.keys(v.apconso).forEach(hash => {
    var valueap = v.apconso[hash]
    if (valueap.id_conso.substring(0, 10) in apart) {
        apart[valueap.id_conso.substring(0, 10)].consommation.push(hash)
    }
  })

  Object.keys(apart).forEach(k => {
    v.apdemande[apart[k].demande].hash_consommation = apart[k].consommation
    for (j in apart[k].consommation) {
        v.apconso[apart[k].consommation[j]].hash_demande = apart[k].demande;
    }
  })

  Object.keys(v.apconso).forEach(h => {
    var conso = v.apconso[h]
    if (conso.hash_demande) {
      var time = conso.periode.getTime()
        if (time in output_indexed){
          output_indexed[time].apart_heures_consommees = output_indexed[time].apart_heures_consommees + conso.heure_consomme;
          output_indexed[time].apart_motif_recours = v.apdemande[conso.hash_demande].motif_recours_se;
        }
      }
  })
}