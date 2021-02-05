import { f } from "./functions"
import { EntréeSirene } from "../GeneratedTypes"
import { Siret, Siren, ParHash } from "../RawDataTypes"

type Input = {
  periode: Date
  siret: Siret
}

type VariablesTransmises = {
  /** Raison sociale de l'établissement. */
  raison_sociale: string | null
  /** Latitude (géolocalisation) de l'établissement. */
  latitude: EntréeSirene["latitude"] | null
  /** Longitude (géolocalisation) de l'établissement. */
  longitude: EntréeSirene["longitude"] | null
  /** Département où est localisé l'établissement. */
  departement: EntréeSirene["departement"] | null
  /** Code APE/NAF de l'établissement. */
  code_ape: EntréeSirene["ape"]
}

export type SortieSirene = VariablesTransmises & {
  /** Numéro SIREN de l'établissement. */
  siren: Siren
  /** Région où est localisé l'établissement. */
  region: ReturnType<typeof f.region>
  /** Année de création de l'établissement. */
  date_creation_etablissement: number | null
  /** Âge de l'établissement, exprimé en années. */
  age: number | null
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "sirene"
  computed: Omit<SortieSirene, keyof VariablesTransmises>
  transmitted: VariablesTransmises
}

export function sirene(
  vSirene: ParHash<EntréeSirene>,
  output_array: (Input & Partial<SortieSirene>)[]
): void {
  "use strict"

  const sireneHashes = Object.keys(vSirene || {})

  output_array.forEach((val) => {
    // geolocalisation

    if (sireneHashes.length !== 0) {
      const sirene = vSirene[
        sireneHashes[sireneHashes.length - 1] as string
      ] as EntréeSirene
      val.siren = val.siret.substring(0, 9)
      val.latitude = sirene.latitude || null
      val.longitude = sirene.longitude || null
      val.departement = sirene.departement || null
      if (val.departement) {
        val.region = f.region(val.departement)
      }
      const regexp_naf = /^[0-9]{4}[A-Z]$/
      if (sirene.ape && sirene.ape.match(regexp_naf)) {
        val.code_ape = sirene.ape
      }
      val.raison_sociale = /*sirene.raison_sociale ||*/ null
      // val.activite_saisonniere = sirene.activite_saisoniere || null
      // val.productif = sirene.productif || null
      // val.tranche_ca = sirene.tranche_ca || null
      // val.indice_monoactivite = sirene.indice_monoactivite || null
      val.date_creation_etablissement = sirene.date_creation
        ? sirene.date_creation.getFullYear()
        : null
      if (val.date_creation_etablissement) {
        val.age =
          sirene.date_creation && sirene.date_creation >= new Date("1901/01/01")
            ? val.periode.getFullYear() - val.date_creation_etablissement
            : null
      }
    }
  })
}
