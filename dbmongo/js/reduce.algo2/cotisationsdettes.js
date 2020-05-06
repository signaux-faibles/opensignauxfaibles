function cotisationsdettes(v, periodes) {
  "use strict";

  // Tous les débits traitées après ce jour du mois sont reportées à la période suivante
  // Permet de s'aligner avec le calendrier de fourniture des données
  const last_treatment_day = 20

  var output_cotisationsdettes = {}

  // TODO Cotisations avec un mois de retard ? Bizarre, plus maintenant que l'export se fait le 20
  // var offset_cotisation = 1
  const offset_cotisation = 0
  var value_cotisation = {}

  // Répartition des cotisations sur toute la période qu'elle concerne
  Object.keys(v.cotisation).forEach(function (h) {
    var cotisation = v.cotisation[h]
    var periode_cotisation = f.generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
    periode_cotisation.forEach(date_cotisation => {
      let date_offset = f.dateAddMonth(date_cotisation, offset_cotisation)
      value_cotisation[date_offset.getTime()] = (value_cotisation[date_offset.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
    })
  })



  // relier les débits
  // ecn: ecart negatif
  // map les débits: clé fabriquée maison => [{hash, numero_historique, date_traitement}, ...]
  // Pour un même compte, les débits avec le même num_ecn (chaque émission de facture) sont donc regroupés
  var ecn = Object.keys(v.debit).reduce((accu, h) => {
      //pour chaque debit
      let debit = v.debit[h]

      var start = debit.periode.start
      var end = debit.periode.end
      var num_ecn = debit.numero_ecart_negatif
      var compte = debit.numero_compte
      var key = start + "-" + end + "-" + num_ecn + "-" + compte
      accu[key] = (accu[key] || []).concat([{
          "hash": h,
          "numero_historique": debit.numero_historique,
          "date_traitement": debit.date_traitement
      }])
      return accu
  }, {})

  // Pour chaque numero_ecn, on trie et on chaîne les débits avec debit_suivant
  Object.keys(ecn).forEach(i => {
      ecn[i].sort(f.compareDebit)
      var l = ecn[i].length
      ecn[i].forEach((e, idx) => {
          if (idx <= l - 2) {
              v.debit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
          }
      })
  })

  var value_dette = {}
  // Pour chaque objet debit:
  // debit_traitement_debut => periode de traitement du débit
  // debit_traitement_fin => periode de traitement du debit suivant, ou bien date_fin
  // Entre ces deux dates, c'est cet objet qui est le plus à jour.
  Object.keys(v.debit).forEach(function (h) {
    var debit = v.debit[h]

    var debit_suivant = (v.debit[debit.debit_suivant] || {"date_traitement" : date_fin})

    //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
    let jour_traitement = debit.date_traitement.getUTCDate()
    let jour_traitement_suivant = debit_suivant.date_traitement.getUTCDate()
    let date_traitement_debut
    if (jour_traitement <= last_treatment_day){
      date_traitement_debut = new Date(
        Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth())
      )
    } else {
      date_traitement_debut = new Date(
        Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth() + 1)
      )
    }

    let date_traitement_fin
    if (jour_traitement_suivant <= last_treatment_day) {
      date_traitement_fin = new Date(
        Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth())
      )
    } else {
      date_traitement_fin = new Date(
        Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth() + 1)
      )
    }

    let periode_debut = date_traitement_debut
    let periode_fin = date_traitement_fin

    //f.generatePeriodSerie exlue la dernière période
    f.generatePeriodSerie(periode_debut, periode_fin).map(date => {
      let time = date.getTime()
      value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale, "montant_majorations": debit.montant_majorations}])
    })
  })

  // TODO faire numero de compte ailleurs
  // Array des numeros de compte
  //var numeros_compte = Array.from(new Set(
  //  Object.keys(v.cotisation).map(function (h) {
  //    return(v.cotisation[h].numero_compte)
  //  })
  //))

  periodes.forEach(function (time) {
    output_cotisationsdettes[time] = output_cotisationsdettes[time] || {}
    var val = output_cotisationsdettes[time]
  //output_cotisationsdettes[time].numero_compte_urssaf = numeros_compte
    if (time in value_cotisation){
      // somme de toutes les cotisations dues pour une periode donnée
      val.cotisation = value_cotisation[time].reduce((a,cot) => a + cot,0)
    }

    // somme de tous les débits (part ouvriere, part patronale, montant_majorations)
    let montant_dette = (value_dette[time] || []).reduce(function (m, dette) {
      m.montant_part_ouvriere += dette.part_ouvriere
      m.montant_part_patronale += dette.part_patronale
      m.montant_majorations += dette.montant_majorations
      return m
    }, {"montant_part_ouvriere": 0, "montant_part_patronale": 0, "montant_majorations": 0})
    val = Object.assign(val, montant_dette)


    let past_month_offsets = [1,2,3,6,12]
    let time_d = new Date(parseInt(time))

    past_month_offsets.forEach(offset => {
      let time_offset = f.dateAddMonth(time_d, offset)
      let variable_name_part_ouvriere = "montant_part_ouvriere_past_" + offset
      let variable_name_part_patronale = "montant_part_patronale_past_" + offset
      output_cotisationsdettes[time_offset.getTime()] = output_cotisationsdettes[time_offset.getTime()] || {}
      let val_offset = output_cotisationsdettes[time_offset.getTime()]
      val_offset[variable_name_part_ouvriere] = val.montant_part_ouvriere
      val_offset[variable_name_part_patronale] = val.montant_part_patronale
    })

    let future_month_offsets = [0, 1, 2, 3, 4, 5]
    if (val.montant_part_ouvriere + val.montant_part_patronale > 0){
      future_month_offsets.forEach(offset => {
        let time_offset = f.dateAddMonth(time_d, offset)
        output_cotisationsdettes[time_offset.getTime()] = output_cotisationsdettes[time_offset.getTime()] || {}
        output_cotisationsdettes[time_offset.getTime()].interessante_urssaf = false
      })
    }
  })

  return(output_cotisationsdettes)
}
