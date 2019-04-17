function cibleApprentissage(output_indexed) {
  // Préparation des transformations de fonctions
  //
  let output_cotisation = output_indexed
  let output_procol = output_indexed
  let output_cible = {}

  let counter = -1
  let all_keys = [...new Set([...output_cotisation, ...output_procol])]
  all_keys.sort((a,b)=> a<=b).forEach( k => {
    if (counter >= 0) counter = counter + 1 
    if (output_procol[k].tag_failure || output_cotisation[k].tag_default){
      counter = 0 
    }
    if (counter >= 0){
      output_cible[k] = output_cible[k] || {}
      output_cible[k].time_til_outcome = counter
    }
  })

  Object.keys(output_cible).forEach( k => {
    if ("time_til_outcome" in output_cible[k] && 
      output_cible[k].time_til_outcome <= 18){
      //||
      //(("arrete_bilan_diane" in v[siret] || "arrete_bilan_bdf" in v[siret]) && 
      //  v[siret].time_til_outcome <= 30)) &&
      //  !("arrete_bilan_diane" in v[siret] && v[siret].arrete_bilan_diane < key.periode &&  
      //  generatePeriodSerie(key.periode, v[siret].arrete_bilan_diane).length >= 18)  &&
      //!("arrete_bilan_bdf" in v[siret] && v[siret].arrete_bilan_bdf < key.periode && 
      //  generatePeriodSerie(key.periode, v[siret].arrete_bilan_bdf).length >= 18)) {
      output_cible[k].outcome = true
    } else output_cible[k].outcome = false
  })
  return (output_cible)
}

