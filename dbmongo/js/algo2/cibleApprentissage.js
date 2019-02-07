function cibleApprentissage(output_indexed) {
  counter = -1
  Object.keys(output_indexed).sort((a,b)=> a<=b).forEach( k => {
    if (counter >=0) counter = counter + 1 
    if (output_indexed[k].tag_outcome == "default" || output_indexed[k].tag_outcome == "failure"){
      counter = 0 
    }
    if (counter >= 0){
      output_indexed[k].time_til_outcome = counter
    }
  })

  Object.keys(output_indexed).forEach( k => {

      if ("time_til_outcome" in output_indexed[k] && 
        output_indexed[k].time_til_outcome <= 18){
        //||
        //(("arrete_bilan_diane" in v[siret] || "arrete_bilan_bdf" in v[siret]) && 
        //  v[siret].time_til_outcome <= 30)) &&
        //  !("arrete_bilan_diane" in v[siret] && v[siret].arrete_bilan_diane < key.periode &&  
        //  generatePeriodSerie(key.periode, v[siret].arrete_bilan_diane).length >= 18)  &&
        //!("arrete_bilan_bdf" in v[siret] && v[siret].arrete_bilan_bdf < key.periode && 
        //  generatePeriodSerie(key.periode, v[siret].arrete_bilan_bdf).length >= 18)) {
        output_indexed[k].outcome = true
      } else output_indexed[k].outcome = false
  })
}
