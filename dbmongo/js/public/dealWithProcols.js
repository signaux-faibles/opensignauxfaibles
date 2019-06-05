function dealWithProcols(data_source, altar_or_procol, output_indexed){
  data_source = data_source || {}
  var codes  =  Object.keys(data_source).reduce((events,hash) => {
    var the_event = data_source[hash]

    let etat = {}
    if (altar_or_procol == "altares")
      etat = f.altaresToHuman(the_event.code_evenement);
    else if (altar_or_procol == "procol")
      etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol);

    if (etat != null)
      events.push({"etat": etat, "date_proc_col": new Date(the_event.date_effet)})

    return(events)
  },[]).sort(
    (a,b) => {return(a.date_proc_col.getTime() > b.date_proc_col.getTime())}
  )
  return(codes)
}
