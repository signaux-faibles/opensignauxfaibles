function sirene (v, output_array) {
  var sireneHashes = Object.keys(v.sirene || {})

  output_array.forEach(val => {
    // geolocalisation

    if (sireneHashes.length != 0) {
      var sirene = v.sirene[sireneHashes[0]]
      val.siren = val.siret.substring(0, 9)
      val.lattitude = sirene.lattitude || null 
      val.longitude = sirene.longitude || null 
      val.region = sirene.region || null 
      val.departement = sirene.departement || null 
      val.code_ape  = sirene.ape || null 
      val.raison_sociale = sirene.raison_sociale || null 
      val.activite_saisonniere = sirene.activite_saisoniere || null 
      val.productif = sirene.productif || null 
      val.tranche_ca = sirene.tranche_ca || null
      val.indice_monoactivite = sirene.indice_monoactivite || null
      val.date_creation = sirene.creation ? sirene.creation.getFullYear() : null
      val.age = val.periode.getFullYear() - val.date_creation
    }
  })
}
