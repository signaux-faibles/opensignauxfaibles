import { f } from "./functions"
import { EntréeCotisation } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"

export type Input = {
  periode: Date
  cotisation?: number
  montant_part_patronale?: number
  montant_part_ouvriere?: number
}

export type SortieCotisation = {
  /** Montant moyen de cotisations calculé sur 12 mois consécutifs. */
  cotisation_moy12m?: number
  /** ratio_dette = (montant_part_ouvriere + montant_part_patronale) / cotisation_moy12m */
  ratio_dette: number
  /** Moyenne de ratio_dette sur 12 mois. */
  ratio_dette_moy12m?: number
  /** Survenance d'un débit d'au moins 1% des cotisations */
  tag_debit: boolean
  /** Survenance de trois débits de 100% (ou plus) des cotisations */
  tag_default: boolean
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "cotisation"
  computed: SortieCotisation
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

export type SortieCotisationsDettes = {
  interessante_urssaf: boolean
  cotisation: number
  montant_part_ouvriere: number
  montant_part_patronale: number
}

export function cotisation(
  vCotisation: ParHash<EntréeCotisation>,
  dateFin: Date // correspond à la variable globale date_fin
): number {
  const dateDebutObservation = new Date(dateFin)
  const dateFinObservation = new Date(dateFin.getTime())
  dateDebutObservation.setFullYear(dateDebutObservation.getFullYear() - 1)

  const value_cotisation = f.makePeriodeMap<number>()

  // Répartition des cotisations sur toute la période qu'elle concerne
  for (const cotisation of Object.values(vCotisation)) {
    const periode_cotisation = f.generatePeriodSerie(
      cotisation.periode.start,
      cotisation.periode.end
    )

    periode_cotisation.forEach((date_cotisation) => {
      if (
        date_cotisation.getTime() <= dateFinObservation.getTime() &&
        date_cotisation.getTime() >= dateDebutObservation.getTime()
      ) {
        value_cotisation.set(
          date_cotisation,
          (value_cotisation.get(date_cotisation) || 0) +
            cotisation.du / periode_cotisation.length
        )
      }
    })
  }

  let somme = 0
  let taille = 0

  for (const t of value_cotisation.values()) {
    somme += t
    taille++
  }
  const result = taille > 0 ? somme / taille : 0

  return Math.round(result * 100) / 100
}
