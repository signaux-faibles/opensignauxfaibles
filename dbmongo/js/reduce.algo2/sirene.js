function sirene(v, output_array) {
    "use strict"
    var sireneHashes = Object.keys(v.sirene || {})

    output_array.forEach((val) => {
        // geolocalisation

        if (sireneHashes.length != 0) {
            var sirene = v.sirene[sireneHashes[sireneHashes.length - 1]]
            val.siren = val.siret.substring(0, 9)
            val.latitude = sirene.lattitude || null
            val.longitude = sirene.longitude || null
            val.departement = sirene.departement || null
            if (val.departement) {
                val.region = f.region(val.departement)
            }
            var regexp_naf = /^[0-9]{4}[A-Z]$/
            if (sirene.ape && sirene.ape.match(regexp_naf)) {
                val.code_ape = sirene.ape
            }
            val.raison_sociale = sirene.raison_sociale || null
            // val.activite_saisonniere = sirene.activite_saisoniere || null
            // val.productif = sirene.productif || null
            // val.tranche_ca = sirene.tranche_ca || null
            // val.indice_monoactivite = sirene.indice_monoactivite || null
            val.date_creation_etablissement = sirene.date_creation
                ? sirene.date_creation.getFullYear()
                : null
            val.age =
                sirene.date_creation &&
                sirene.date_creation >= new Date("1901/01/01")
                    ? val.periode.getFullYear() -
                      val.date_creation_etablissement
                    : null
        }
    })
}
