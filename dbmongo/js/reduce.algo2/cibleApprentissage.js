function cibleApprentissage(output_indexed) {
  // Préparation des transformations de fonctions
  //
  let output_cotisation = output_indexed
  let output_procol = output_indexed

  let all_keys = [...new Set([...Object.keys(output_cotisation), ...Object.keys(output_procol)])]
  
  let merged_info = {}
  all_keys.forEach(k => {
    merged_info[k] = merged_info[k] || {outcome: false}
    if (output_procol[k].tag_failure || output_cotisation[k].tag_default)
      merged_info[k]["outcome"] = true
  })

  let output_cible = {}

  let output_cible1 = f.lookAhead(merged_info, "outcome", 18, true)
  let output_cible2 = f.lookAhead(output_cotisation, "tag_debit", 18, true)

  all_keys.forEach(k => {
    output_cible[k] = {}
    if (output_cible1[k]){
      output_cible[k].outcome = output_cible1[k].outcome
      output_cible[k].time_til_outcome = output_cible1[k].time_til_outcome
    }
    if (output_cible2[k])
      output_cible[k].debit_18m = output_cible2[k].outcome
  })

  return (output_cible)
  //let counter = -1
  //all_keys.sort((a,b)=> a<=b).forEach( k => {
  //  if (counter >= 0) counter = counter + 1 
  //  if (output_procol[k].tag_failure || output_cotisation[k].tag_default){
  //    counter = 0 
  //  }
  //  if (counter >= 0){
  //    output_cible[k] = output_cible[k] || {}
  //    output_cible[k].time_til_outcome = counter
  //  }
  //})

  //Object.keys(output_cible).forEach( k => {
  //  if ("time_til_outcome" in output_cible[k] && 
  //    output_cible[k].time_til_outcome <= 18){
  //    //||
  //    //(("arrete_bilan_diane" in v[siret] || "arrete_bilan_bdf" in v[siret]) && 
  //    //  v[siret].time_til_outcome <= 30)) &&
  //    //  !("arrete_bilan_diane" in v[siret] && v[siret].arrete_bilan_diane < key.periode &&  
  //    //  generatePeriodSerie(key.periode, v[siret].arrete_bilan_diane).length >= 18)  &&
  //    //!("arrete_bilan_bdf" in v[siret] && v[siret].arrete_bilan_bdf < key.periode && 
  //    //  generatePeriodSerie(key.periode, v[siret].arrete_bilan_bdf).length >= 18)) {
  //    output_cible[k].outcome = true
  //  } else output_cible[k].outcome = false
  //})
}

