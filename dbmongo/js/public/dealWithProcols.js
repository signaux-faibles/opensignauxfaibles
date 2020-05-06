function dealWithProcols(data_source, altar_or_procol, output_indexed){
  "use strict";
  return Object.keys(data_source || {}).reduce((events,hash) => {
    var the_event = data_source[hash]

    let etat = {}
    if (altar_or_procol == "altares")
      etat = f.altaresToHuman(the_event.code_evenement);
    else if (altar_or_procol == "procol")
      etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol);

    if (etat != null)
      events.push({"etat": etat, "date_procol": new Date(the_event.date_effet)})

    return(events)
  },[]).sort(
    (a,b) => {return(a.date_procol.getTime() > b.date_procol.getTime())}
  )
}
