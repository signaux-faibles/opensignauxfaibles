function sirene_ul(v, output_array) {
  "use strict"
  var sireneHashes = Object.keys(v.sirene_ul || {})
  output_array.forEach((val) => {
    if (sireneHashes.length != 0) {
      var sirene = v.sirene_ul[sireneHashes[sireneHashes.length - 1]]
      val.siren = val.siren
      val.raison_sociale = f.raison_sociale(
        sirene.raison_sociale,
        sirene.nom_unite_legale,
        sirene.nom_usage_unite_legale,
        sirene.prenom1_unite_legale,
        sirene.prenom2_unite_legale,
        sirene.prenom3_unite_legale,
        sirene.prenom4_unite_legale
      )
      val.statut_juridique = sirene.statut_juridique || null
      val.date_creation_entreprise = sirene.date_creation
        ? sirene.date_creation.getFullYear()
        : null
      val.age_entreprise =
        sirene.date_creation && sirene.date_creation >= new Date("1901/01/01")
          ? val.periode.getFullYear() - val.date_creation_entreprise
          : null
    }
  })
}
