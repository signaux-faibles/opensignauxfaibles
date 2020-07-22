import * as f from "../common/raison_sociale"

type Input = {
  periode: Date
}

export type SortieSireneEntreprise = {
  raison_sociale: string // nom de l'entreprise
  statut_juridique: string | null // code numérique sérialisé en chaine de caractères
  date_creation_entreprise: number | null // année
  age_entreprise?: number // en années
}

export function entr_sirene(
  v: DonnéesSireneEntreprise,
  output_array: (Input & Partial<SortieSireneEntreprise>)[]
): void {
  "use strict"
  const sireneHashes = Object.keys(v.sirene_ul || {})
  output_array.forEach((val) => {
    if (sireneHashes.length !== 0) {
      const sirene = v.sirene_ul[sireneHashes[sireneHashes.length - 1]]
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
      if (
        val.date_creation_entreprise &&
        sirene.date_creation &&
        sirene.date_creation >= new Date("1901/01/01")
      ) {
        val.age_entreprise =
          val.periode.getFullYear() - val.date_creation_entreprise
      }
    }
  })
}
