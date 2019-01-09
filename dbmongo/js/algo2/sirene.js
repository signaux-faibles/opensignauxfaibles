function sirene (v, output_array) {
  var sireneHashes = Object.keys(v.sirene || {})

  output_array.forEach(val => {
      if (sireneHashes.length != 0) {
          sirene = v.sirene[sireneHashes[0]]
      }
    
      val.lattitude = (sirene || { "lattitude": null }).lattitude
      val.longitude = (sirene || { "longitude": null }).longitude
      val.region = (sirene || {"region": null}).region
      val.departement = (sirene || {"departement": null}).departement
      val.code_ape  = (sirene || { "ape": null}).ape
      val.raison_sociale = (sirene || {"raisonsociale": null}).raisonsociale
      val.activite_saisonniere = (sirene || {"activitesaisoniere": null}).activitesaisoniere
      val.productif = (sirene || {"productif": null}).productif
      val.debut_activite = (sirene || {"debut_activite":null}).debut_activite.getFullYear()
      val.tranche_ca = (sirene || {"trancheca":null}).trancheca
      val.indice_monoactivite = (sirene || {"indicemonoactivite": null}).indicemonoactivite  
  })
}