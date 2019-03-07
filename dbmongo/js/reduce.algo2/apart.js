function apart (v, output_effectif) {

  var output_apart = {}

  // Mapping (pour l'instant vide) du hash de la demande avec les hash des consos correspondantes
  var apart = Object.keys(v.apdemande).reduce((apart, hash) => {
    apart[v.apdemande[hash].id_demande.substring(0, 9)] = {
      "demande": hash,
      "consommation": [],
      "periode_debut": 0,
      "periode_fin": 0
    }
    return apart
  }, {})

  // on note le nombre d'heures demandées dans output_apart
  Object.keys(v.apdemande).forEach(hash => {
    
    var periode_deb = v.apdemande[hash].periode.start
    var periode_fin = v.apdemande[hash].periode.end
    // Des periodes arrondies aux débuts de périodes
    // TODO meilleur arrondi
    var periode_deb_floor = new Date(Date.UTC(periode_deb.getUTCFullYear(), periode_deb.getUTCMonth(), 1, 0, 0, 0, 0))
    var periode_fin_ceil = new Date(Date.UTC(periode_fin.getUTCFullYear(), periode_fin.getUTCMonth() + 1, 1, 0, 0, 0, 0))
    apart[v.apdemande[hash].id_demande.substring(0, 9)].periode_debut = periode_deb_floor
    apart[v.apdemande[hash].id_demande.substring(0, 9)].periode_fin = periode_fin_ceil
    
    var series = generatePeriodSerie(periode_deb_floor, periode_fin_ceil)
    series.forEach( date => {
      let time = date.getTime()
      output_apart[time] = output_apart[time] || {}
      output_apart[time].apart_heures_autorisees = v.apdemande[hash].hta
    })
  })

  // relier les consos faites aux demandes (hashs) dans apart
  Object.keys(v.apconso).forEach(hash => {
    var valueap = v.apconso[hash]
    if (valueap.id_conso.substring(0, 9) in apart) {
      apart[valueap.id_conso.substring(0, 9)].consommation.push(hash)
    }
  })

  //Object.keys(apart).forEach(k => {
  //  v.apdemande[apart[k].demande].hash_consommation = apart[k].consommation
  //  for (let j in apart[k].consommation) {
  //    v.apconso[apart[k].consommation[j]].hash_demande = apart[k].demande
  //  }
  //})

  Object.keys(apart).forEach(k => {
    if (apart[k].consommation.length > 0) {
      apart[k].consommation.sort(
        (a,b) => (v.apconso[a].periode.getTime() >= v.apconso[b].periode.getTime())
      ).forEach( (h) => {
        var time = v.apconso[h].periode.getTime()
        output_apart[time] = output_apart[time] || {}
        output_apart[time].apart_heures_consommees = (output_apart[time].apart_heures_consommees || 0) + v.apconso[h].heure_consomme
        output_apart[time].apart_motif_recours = v.apdemande[apart[k].demande].motif_recours_se
      })

      // Heures consommees cumulees sur la demande
      let series = generatePeriodSerie(apart[k].periode_debut, apart[k].periode_fin)
      series.reduce( (accu, date) => {
        let time = date.getTime()
        //output_apart est déjà défini pour les heures autorisées
        accu = accu + (output_apart[time].apart_heures_consommees || 0)
        output_apart[time].apart_heures_consommees_cumulees = accu
        return(accu)
      }, 0)
    }
  })
  //Object.keys(v.apconso).forEach(h => {
  //  // Pour toutes les consos
  //  var conso = v.apconso[h]
  //  // on regard s'il y a une demande correspondante
  //  if (conso.hash_demande) {
  //    var time = conso.periode.getTime()
  //    if (time in periodes){
  //
  //      output_apart[time].apart_heures_consommees = output_apart[time].apart_heures_consommees + conso.heure_consomme
  //      output_apart[time].apart_motif_recours = v.apdemande[conso.hash_demande].motif_recours_se
  //    }
  //  }
  //})


  Object.keys(output_apart).forEach(time => {
    if (time in output_effectif){
      output_apart[time].ratio_apart = (output_apart[time].apart_heures_consommees || 0) / (output_effectif[time].effectif * 157.67) 
      print(output_apart[time].apart_heures_consommees)
      print(output_effectif[time].effectif)
      //nbr approximatif d'heures ouvrées par mois
    }
  })
  return(output_apart)
}
