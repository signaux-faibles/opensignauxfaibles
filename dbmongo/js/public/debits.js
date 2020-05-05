"use strict";

function debits(vdebit) {
  const  { date_fin, serie_periode } = this; // parameters passed from Go Pipeline as global variables  

  const last_treatment_day = 20
  vdebit = vdebit || {}
  var ecn = Object.keys(vdebit).reduce((accu, h) => {
      let debit = vdebit[h]
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

  Object.keys(ecn).forEach(i => {
      ecn[i].sort(f.compareDebit)
      var l = ecn[i].length
      ecn[i].forEach((e, idx) => {
          if (idx <= l - 2) {
              vdebit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
          }
      })
  })

  var value_dette = {}

  Object.keys(vdebit).forEach(function (h) {
    var debit = vdebit[h]

    var debit_suivant = (vdebit[debit.debit_suivant] || {"date_traitement" : date_fin})
    
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

    //generatePeriodSerie exlue la dernière période
    f.generatePeriodSerie(periode_debut, periode_fin).map(date => {
      let time = date.getTime()
      value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale, "montant_majorations": debit.montant_majorations}])
    })
  })    

  const output_dette = []
  serie_periode.forEach(p => {
    output_dette.push(
      (value_dette[p.getTime()] || [])
        .reduce((m,c) => {
          return {
            part_ouvriere: m.part_ouvriere + c.part_ouvriere,
            part_patronale: m.part_patronale + c.part_patronale,
            periode: f.dateAddDay(f.dateAddMonth(p,1),-1) }
          }, {part_ouvriere: 0, part_patronale: 0})
    )
  })

  return(output_dette)
}
