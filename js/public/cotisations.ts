import { f } from "./functions"
import { EntréeCotisation } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"

// Paramètres globaux utilisés par "public"
declare const serie_periode: Date[]

export function cotisations(
  vcotisation: ParHash<EntréeCotisation> = {}
): number[] {
  const offset_cotisation = 0
  const value_cotisation: Record<number, number[]> = {}

  // Répartition des cotisations sur toute la période qu'elle concerne
  for (const cotisation of Object.values(vcotisation)) {
    const periode_cotisation = f.generatePeriodSerie(
      cotisation.periode.start,
      cotisation.periode.end
    )
    periode_cotisation.forEach((date_cotisation) => {
      const date_offset = f.dateAddMonth(date_cotisation, offset_cotisation)
      value_cotisation[date_offset.getTime()] = (
        value_cotisation[date_offset.getTime()] || []
      ).concat([cotisation.du / periode_cotisation.length])
    })
  }

  const output_cotisation: number[] = []

  serie_periode.forEach((p) => {
    output_cotisation.push(
      (value_cotisation[p.getTime()] || []).reduce((m, c) => m + c, 0)
    )
  })

  return output_cotisation
}
