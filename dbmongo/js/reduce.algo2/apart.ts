import { f } from "./functions"
import { EntréeApDemande, EntréeApConso } from "../RawDataTypes"

type ApConsoHash = string

type Hash = string

type Timestamp = string

export type SortieAPart = {
  apart_heures_autorisees: unknown
  apart_heures_consommees: number
  apart_motif_recours: EntréeApDemande["motif_recours_se"]
  apart_heures_consommees_cumulees: number
}

export function apart(
  apconso: Record<ApConsoHash, EntréeApConso>,
  apdemande: Record<Hash, EntréeApDemande>
): Record<Timestamp, SortieAPart> {
  "use strict"

  const output_apart: Record<Timestamp, SortieAPart> = {}

  // Mapping (pour l'instant vide) du hash de la demande avec les hash des consos correspondantes
  const apart: Record<
    string,
    {
      demande: Hash
      consommation: ApConsoHash[]
      periode_debut: Date
      periode_fin: Date
    }
  > = {}
  for (const [hash, apdemandeEntry] of Object.entries(apdemande)) {
    apart[apdemandeEntry.id_demande.substring(0, 9)] = {
      demande: hash,
      consommation: [],
      periode_debut: new Date(0),
      periode_fin: new Date(0),
    }
  }

  // on note le nombre d'heures demandées dans output_apart
  for (const apdemandeEntry of Object.values(apdemande)) {
    const periode_deb = apdemandeEntry.periode.start
    const periode_fin = apdemandeEntry.periode.end

    // Des periodes arrondies aux débuts de périodes
    // TODO: arrondir au debut du mois le plus proche, au lieu de tronquer la date. (ex: cas du dernier jour d'un mois)
    const periode_deb_floor = new Date(
      Date.UTC(
        periode_deb.getUTCFullYear(),
        periode_deb.getUTCMonth(),
        1,
        0,
        0,
        0,
        0
      )
    )
    const periode_fin_ceil = new Date(
      Date.UTC(
        periode_fin.getUTCFullYear(),
        periode_fin.getUTCMonth() + 1,
        1,
        0,
        0,
        0,
        0
      )
    )
    const apartForSiren = apart[apdemandeEntry.id_demande.substring(0, 9)]
    if (apartForSiren === undefined) {
      const error = (message: string): never => {
        throw new Error(message)
      }
      error("siren should be included in apart")
    } else {
      apartForSiren.periode_debut = periode_deb_floor
      apartForSiren.periode_fin = periode_fin_ceil
    }

    const series = f.generatePeriodSerie(periode_deb_floor, periode_fin_ceil)
    series.forEach((date) => {
      const time = date.getTime()
      output_apart[time] = {
        ...(output_apart[time] ?? ({} as SortieAPart)),
        apart_heures_autorisees: apdemandeEntry.hta,
      }
    })
  }

  // relier les consos faites aux demandes (hashs) dans apart
  for (const [hash, valueap] of Object.entries(apconso)) {
    const apartForSiren = apart[valueap.id_conso.substring(0, 9)]
    if (apartForSiren !== undefined) {
      apartForSiren.consommation.push(hash)
    }
  }

  for (const apartEntry of Object.values(apart)) {
    if (apartEntry.consommation.length > 0) {
      apartEntry.consommation
        .sort(
          (a, b) =>
            (apconso[a]?.periode ?? new Date()).getTime() -
            (apconso[b]?.periode ?? new Date()).getTime() // TODO: use `never` type assertion here?
        )
        .forEach((h) => {
          const time = apconso[h]?.periode.getTime()
          if (time === undefined) {
            return
          }
          const current = output_apart[time] ?? ({} as SortieAPart)
          const heureConso = apconso[h]?.heure_consomme
          if (heureConso !== undefined) {
            current.apart_heures_consommees =
              (current.apart_heures_consommees ?? 0) + heureConso
          }
          const motifRecours = apdemande[apartEntry.demande]?.motif_recours_se
          if (motifRecours !== undefined) {
            current.apart_motif_recours = motifRecours
          }
          output_apart[time] = current
        })

      // Heures consommees cumulees sur la demande
      const series = f.generatePeriodSerie(
        apartEntry.periode_debut,
        apartEntry.periode_fin
      )
      series.reduce((accu, date) => {
        const time = date.getTime()

        //output_apart est déjà défini pour les heures autorisées
        const current = output_apart[time] ?? ({} as SortieAPart)
        accu = accu + (current.apart_heures_consommees || 0)
        output_apart[time] = {
          ...current,
          apart_heures_consommees_cumulees: accu,
        }

        return accu
      }, 0)
    }
  }

  // Note: à la fin de l'opération map-reduce, dbmongo va calculer la propriété
  // ratio_apart depuis apart.crossComputation.json.

  return output_apart
}
