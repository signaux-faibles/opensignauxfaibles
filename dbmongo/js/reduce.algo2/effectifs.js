function effectifs (effobj, periodes, effectif_name) {
  "use strict";

  let output_effectif = {}

  // Construction d'une map[time] = effectif à cette periode
  let map_effectif = Object.keys(effobj).reduce((m, hash) => {
    var effectif = effobj[hash]
    if (effectif == null) {
      return m
    }
    var effectifTime = effectif.periode.getTime()
    m[effectifTime] = (m[effectifTime] || 0) + effectif.effectif
    return m
  }, {})

  //ne reporter que si le dernier est disponible
  // 1- quelle periode doit être disponible
  var last_period = new Date(parseInt(periodes[periodes.length - 1]))
  var last_period_offset = f.dateAddMonth(last_period, jsParams.offset_effectif + 1)
  // 2- Cette période est-elle disponible ?

  var available = map_effectif[last_period_offset.getTime()] ? 1 : 0


  //pour chaque periode (elles sont triees dans l'ordre croissant)
  periodes.reduce((accu, time) => {
    var periode = new Date(parseInt(time))
    // si disponible on reporte l'effectif tel quel, sinon, on recupère l'accu
    output_effectif[time] = output_effectif[time] || {}
    output_effectif[time][effectif_name] = map_effectif[time] || (available ? accu : null)


    // le cas échéant, on met à jour l'accu avec le dernier effectif disponible
    accu = map_effectif[time] || accu

    output_effectif[time][effectif_name + "_reporte"] = map_effectif[time] ? 0 : 1
    return(accu)
  }, null)

  Object.keys(map_effectif).forEach(time => {
    var periode = new Date(parseInt(time))
    var past_month_offsets = [6,12,18,24]
    past_month_offsets.forEach(lookback => {
      // On ajoute un offset pour partir de la dernière période où l'effectif est connu
      var time_past_lookback = f.dateAddMonth(periode, lookback - jsParams.offset_effectif - 1)

      var variable_name_effectif = effectif_name + "_past_" + lookback
      output_effectif[time_past_lookback.getTime()] = output_effectif[time_past_lookback.getTime()] || {}
      output_effectif[time_past_lookback.getTime()][variable_name_effectif] = map_effectif[time]
    })
  })

  // On supprime les effectifs 'null'
  Object.keys(output_effectif).forEach(k => {
    if (output_effectif[k].effectif == null && output_effectif[k].effectif_ent == null) {
      delete output_effectif[k]
    }
  })
  return(output_effectif)
}
