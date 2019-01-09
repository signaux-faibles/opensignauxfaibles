function cotisationsdettes(v, output_array, output_indexed) {
  var value_cotisation = {}
  Object.keys(v.cotisation).forEach(function (h) {
      var cotisation = v.cotisation[h]
      var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
      periode_cotisation.forEach(function (date_cotisation) {
          value_cotisation[date_cotisation.getTime()] = (value_cotisation[date_cotisation.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
      })
  })
  
  var value_dette = {}
  
  Object.keys(v.debit).forEach(function (h) {
      var debit = v.debit[h]
      if (debit.part_ouvriere + debit.part_patronale > 0) {

          var debit_suivant = (v.debit[debit.debit_suivant] || {"date_traitement" : date_fin})
          date_limite = date_fin//new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))
          date_traitement_debut = new Date(
              Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth())
          )
          
          date_traitement_fin = new Date(
              Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth())
          )
          
          periode_debut = (date_traitement_debut.getTime() >= date_limite.getTime() ? date_limite : date_traitement_debut)
          periode_fin = (date_traitement_fin.getTime() >= date_limite.getTime() ? date_limite : date_traitement_fin)
          
          generatePeriodSerie(periode_debut, periode_fin).map(function (date) {
              time = date.getTime()
              value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale }])
          })
      }
  })    

  Object.keys(output_indexed).forEach(function (time) {
      if (time in value_cotisation){
          output_indexed[time].cotisation = value_cotisation[time].reduce((a,cot) => a + cot,0)
      }
      
      if (time in value_dette) {
          output_indexed[time].debit_array = value_dette[time]
      }
  })

  output_array.forEach(function (val) {
      
      val.montant_dette = val.debit_array.reduce(function (m, dette) {
          m.part_ouvriere += dette.part_ouvriere
          m.part_patronale += dette.part_patronale
          return m
      }, { "part_ouvriere": 0, "part_patronale": 0 })
      
      val.montant_part_ouvriere = val.montant_dette.part_ouvriere
      val.montant_part_patronale = val.montant_dette.part_patronale

      delete val.montant_dette
      delete val.debit_array
      
  })
}