function defaillances (v, output_indexed) {
  // On filtre altares pour ne garder que les codes qui nous intéressents
  var altares_codes  =  Object.keys(v.altares).reduce((events,hash) => {
      var altares_event = v.altares[hash]
      var etat = altaresToHuman(altares_event.code_evenement)
      if (etat != null)
      events.push({"etat": etat, "date_proc_col": new Date(altares_event.date_effet)})
      return(events)
  },[{"etat" : "in_bonis", "date_proc_col" : new Date(0)}]).sort((a, b) => {
      return(a.date_proc_col.getTime() > b.date_proc_col.getTime())
    }
  )

  altares_codes.forEach(event => {
    var periode_effet = new Date(Date.UTC(event.date_proc_col.getFullYear(), event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0))
    var time_til_last = Object.keys(output_indexed).filter(val => {return (val >= periode_effet)})
    time_til_last.forEach(time => {
      if (time in output_indexed) {
        output_indexed[time].etat_proc_collective = event.etat
        output_indexed[time].date_proc_collective = event.date_proc_col
      }
    })
  })
}
  
  
  