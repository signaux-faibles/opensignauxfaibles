function finalize(k, v) {
  const maxBsonSize = 16777216;

  //v de la forme
  // .id: {batch / siren / periode}
  // value: {siret1: {}, siret2: {}, "siren": {}}
  //
  ///
  ///////////////////////////////////////////////
  // consolidation a l'echelle de l'entreprise //
  ///////////////////////////////////////////////
  ///
  //

  let etablissements_connus = []
  let entreprise = (v.entreprise || {})

  Object.keys(v).forEach(siret =>{
    if (siret != "entreprise" && siret != "siren" ) {
      etablissements_connus[siret] = true
      //if (v[siret]){  // always TRUE
      //    var time = v[siret].periode.getTime()
      if (v[siret].effectif){
        entreprise.effectif_entreprise = (entreprise.effectif_entreprise || 0) + v[siret].effectif // initialized to null
      }
      if (v[siret].apart_heures_consommees){
        entreprise.apart_entreprise = (entreprise.apart_entreprise || 0) + v[siret].apart_heures_consommees // initialized to 0
      }
      if (v[siret].montant_part_patronale || v[siret].montant_part_ouvriere){
        entreprise.debit_entreprise = (entreprise.debit_entreprise || 0) +
          (v[siret].montant_part_patronale || 0) +
          (v[siret].montant_part_ouvriere || 0)
      }
      // not initialized
      //}
    }
  })

  Object.keys(v).forEach(siret =>{
    Object.assign(v[siret], entreprise)
  })

  //une fois que les comptes sont faits...
  let output = []
  let nb_connus = Object.keys(etablissements_connus).length
  Object.keys(v).forEach(siret => {
    if (siret != "entreprise" && siret != "siren" && v[siret]) {
      v[siret].nbr_etablissements_connus = nb_connus
      output.push(v[siret])
    }
  })

  if (output.length > 0 && nb_connus <= 1500){
    if ((Object.bsonsize(output)  + Object.bsonsize({"_id": k})) < maxBsonSize){
      return output
    } else {
      print("Warning: my name is " + JSON.stringify(key, null, 2) + " and I died in reduce.algo2/finalize.js")
      return {"incomplete": true}
    }
  }
}
