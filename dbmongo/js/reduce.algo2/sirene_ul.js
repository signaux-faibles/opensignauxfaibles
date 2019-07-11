function sirene_ul(v, output_array) {
  var sireneHashes = Object.keys(v.sirene_ul || {})
  output_array.forEach(val => {
    if (sireneHashes.length != 0) {
      var sirene = v.sirene_ul[sireneHashes[0]]
      val.siren = val.siren
      val.raison_sociale = sirene.raison_sociale || null
      val.statut_juridique = sirene.statut_juridique || null
    }
  })
}
