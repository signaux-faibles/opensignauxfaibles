function map() {
  if (value.scope=='etablissement') {
    v = Object.keys((this.value.batch || {}))
    .filter(b => b <= actual_batch)
    .sort()
    .reduce((m, batch) => {
      Object.keys(this.value.batch[batch])
      .filter(type => ['sirene', 'effectif', 'debit', 'ccsf', 'procol'].includes(type))
      .forEach((type) => {
        m[type] = (m[type] || {})
        var  array_delete = (this.value.batch[batch].compact.delete[type]||[])
        if (array_delete != {}) {array_delete.forEach(hash => {
          delete m[type][hash]
        })
        }
        Object.assign(m[type], this.value.batch[batch][type])
      })
      return m
    }, {})

    var offset_effectif = (date_fin_effectif.getUTCFullYear() - date_fin.getUTCFullYear()) * 12 + date_fin_effectif.getUTCMonth() - date_fin.getUTCMonth()

    var sirene = Object.keys(v.sirene || {}).reduce((accu, k) => {
      accu.raisonsociale = v.sirene[k].raisonsociale
      accu.nicsiege = v.sirene[k].nicsiege
      accu.adresse = v.sirene[k].adresse
      accu.gps = [v.sirene[k].lattitude, v.sirene[k].longitude]
      accu.creation = v.sirene[k].creation
      accu.ape = v.sirene[k].ape
      return accu
    }, {})


    // Effectifs
    var map_effectif = Object.keys(v.effectif || {}).reduce(function (map_effectif, hash) {
      var effectif = v.effectif[hash];
      if (effectif == null) {
          return map_effectif
      }
      var effectifTime = effectif.periode.getTime()
      map_effectif[effectifTime] = (map_effectif[effectifTime] || 0) + effectif.effectif
      return map_effectif
    }, {})

    effectif = map_effectif[date_fin_effectif.getTime()]
    effectif_precedent = map_effectif[DateAddMonth(date_fin_effectif, -12).getTime()]

    date_fin_output = DateAddMonth(date_fin, -6)
    liste_periodes = generatePeriodSerie(date_fin_output, date_fin)
    
    var output_array = liste_periodes.map(function (e) {
      return {
          "periode": e,
          "apart_heures_consommees": 0,
          "apart_motif_recours": 0,
          "debit_array": [],
          "etat_proc_collective": "in_bonis"
      }
    });

    var output_indexed = output_array.reduce(function (periode, val) {
      periode[val.periode.getTime()] = val
      return periode
    }, {});

    v.debit = (v.debit || {})
    var value_dette = {}
    var value_cotisation = {}

    // Débits 
    
    // var ecn = Object.keys(v.debit).reduce((m, h) => {
    //   var d = [h, v.debit[h]]
    //   var start = d[1].periode.start
    //   var end = d[1].periode.end
    //   var num_ecn = d[1].numero_ecart_negatif
    //   var compte = d[1].numero_compte
    //   var key = start + "-" + end + "-" + num_ecn + "-" + compte
    //   m[key] = (m[key] || []).concat([{
    //       "hash": d[0],
    //       "numero_historique": d[1].numero_historique,
    //       "date_traitement": d[1].date_traitement
    //   }])
    //   return m
    // }, {})

    Object.keys(v.debit).forEach(function (h) {
      var debit = v.debit[h]
      if (debit.part_ouvriere + debit.part_patronale > 0) {

          var debit_suivant = (v.debit[debit.debit_suivant] || {"date_traitement" : date_fin})
          date_limite = date_fin //new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))
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

    var urssaf = output_array.reduce((accu, output) => {
      if (output.montant_part_patronale + output.montant_part_ouvriere > 0) {
        return accu || true
      } 
      return accu
    }, false )

    function deal_with_procols(data_source, altar_or_procol){

      var codes  =  Object.keys(data_source).reduce(function(events,hash) {
        var the_event = data_source[hash]

        if (altar_or_procol == "altares")
          var etat = altaresToHuman(the_event.code_evenement);
        else if (altar_or_procol == "procol")
          var etat = procolToHuman(the_event.action_procol, the_event.stade_procol);

        if (etat != null)
          events.push({"etat": etat, "date_proc_col": new Date(the_event.date_effet)})

        return(events)
      },[]).sort(
        function(a,b) {return(a.date_proc_col.getTime() > b.date_proc_col.getTime())}
      )

      codes.forEach(
        function (event) {
          let periode_effet = new Date(Date.UTC(event.date_proc_col.getFullYear(), event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0))
          var time_til_last = Object.keys(output_indexed).filter(val => {return (val >= periode_effet)})

          time_til_last.forEach(time => {
            if (time in output_indexed) {
              output_indexed[time].etat_proc_collective = event.etat
              output_indexed[time].date_proc_collective = event.date_proc_col
              if (event.etat != "in_bonis")
                output_indexed[time].tag_outcome = "failure"
            }
          })
        }
      )
    }

    deal_with_procols((v.procol || {}), "procol")

    
    r = {
      sirene,
      effectif,
      effectif_precedent,
      urssaf,
      procol: (output_array[output_array.length-1]||{'etat_proc_collective' : 'in_bonis'}).etat_proc_collective
    }

    emit({siret: this.value.siret, batch: actual_batch}, r) 
  } else if (value.scope=='entreprise') {
      v = Object.keys((this.value.batch || {})).sort().filter(batch => batch <= actual_batch).reduce((m, batch) => {
        Object.keys(this.value.batch[batch])
            .filter(type => ['bdf', 'diane']
            .includes(type))
            .forEach((type) => {
                m[type] = (m[type] || {})
                var  array_delete = (this.value.batch[batch].compact.delete[type]||[])
                if (array_delete != {}) {array_delete.forEach(hash => {
                    delete m[type][hash]
                })
                }
                Object.assign(m[type], this.value.batch[batch][type])
        })
        return m
      }, { "siren": this.value.siren })
  
      v.bdf = Object.keys(v.bdf || {}).reduce((accu, key) => {
        var bdf = {
            annee: v.bdf[key].annee,
            arrete_bilan: v.bdf[key].arrete_bilan,
            poids_frng: v.bdf[key].poids_frng,
            taux_marge: v.bdf[key].taux_marge,
            delai_fournisseur: v.bdf[key].delai_fournisseur,
            dette_fiscale: v.bdf[key].dette_fiscale,
            financier_court_terme: v.bdf[key].financier_court_terme,
            frais_financier: v.bdf[key].frais_financier
        }
        accu = accu.concat(bdf)
        return accu
      }, []).sort((a1, a2) => a1.annee < a2.annee)
  
      v.diane = Object.keys(v.diane || {}).reduce((accu, key) => {
        accu = accu.concat(v.diane[key])
        return accu
      }, []).sort((a1, a2) => a1.exercice_diane < a2.exercice_diane)
  
      emit({siren: this.value.siren, batch: actual_batch}, v) 
  }
}