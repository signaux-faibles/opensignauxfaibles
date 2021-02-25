import { f } from "./functions"
import { EntréeSireneEntreprise } from "../GeneratedTypes"
import { ParHash, ParPériode } from "../RawDataTypes"

type VariablesTransmises = {
  /** Catégorie juridique de l’unité légale. Nomenclature: https://www.insee.fr/fr/information/2028129. */
  statut_juridique: string | null
}

export type SortieSireneEntreprise = VariablesTransmises & {
  /** Nom de l'entreprise. Composé à partir des champs raison_sociale, nom_unite_legale, nom_usage_unite_legale, prenom1_unite_legale, prenom2_unite_legale, prenom3_unite_legale et prenom4_unite_legal. */
  raison_sociale: string
  /** Année de création de l'entreprise. */
  date_creation_entreprise: number | null
  /** Age de l'entreprise, en nombre d'années. */
  age_entreprise?: number
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "entr_sirene"
  computed: Omit<SortieSireneEntreprise, keyof VariablesTransmises>
  transmitted: VariablesTransmises
}

export function entr_sirene(
  sirene_ul: ParHash<EntréeSireneEntreprise>,
  sériePériode: Date[]
): ParPériode<Partial<SortieSireneEntreprise>> {
  "use strict"
  const retourEntrSirene = new ParPériode<Partial<SortieSireneEntreprise>>()
  const sireneHashes = Object.keys(sirene_ul || {})
  sériePériode.forEach((période) => {
    if (sireneHashes.length !== 0) {
      const sirene = sirene_ul[
        sireneHashes[sireneHashes.length - 1] as string
      ] as EntréeSireneEntreprise
      const val: Partial<SortieSireneEntreprise> = {}
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
          période.getFullYear() - val.date_creation_entreprise
      }
      retourEntrSirene.set(période, val)
    }
  })
  return retourEntrSirene
}
