import * as f from "../common/generatePeriodSerie"

type ApConsoHash = string

type Hash = string

type Timestamp = string

type SortieAPart = {
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
  const apart = Object.keys(apdemande).reduce((apart, hash) => {
    apart[apdemande[hash].id_demande.substring(0, 9)] = {
      demande: hash,
      consommation: [],
      periode_debut: new Date(0),
      periode_fin: new Date(0),
    }
    return apart
  }, {} as Record<string, { demande: Hash; consommation: ApConsoHash[]; periode_debut: Date; periode_fin: Date }>)

  // on note le nombre d'heures demandées dans output_apart
  Object.keys(apdemande).forEach((hash) => {
    const periode_deb = apdemande[hash].periode.start
    const periode_fin = apdemande[hash].periode.end

    // Des periodes arrondies aux débuts de périodes
    // TODO meilleur arrondi
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
    apart[
      apdemande[hash].id_demande.substring(0, 9)
    ].periode_debut = periode_deb_floor
    apart[
      apdemande[hash].id_demande.substring(0, 9)
    ].periode_fin = periode_fin_ceil

    const series = f.generatePeriodSerie(periode_deb_floor, periode_fin_ceil)
    series.forEach((date) => {
      const time = date.getTime()
      output_apart[time] = output_apart[time] || {}
      output_apart[time].apart_heures_autorisees = apdemande[hash].hta
    })
  })

  // relier les consos faites aux demandes (hashs) dans apart
  Object.keys(apconso).forEach((hash) => {
    const valueap = apconso[hash]
    if (valueap.id_conso.substring(0, 9) in apart) {
      apart[valueap.id_conso.substring(0, 9)].consommation.push(hash)
    }
  })

  Object.keys(apart).forEach((k) => {
    if (apart[k].consommation.length > 0) {
      apart[k].consommation
        .sort(
          (a, b) => apconso[a].periode.getTime() - apconso[b].periode.getTime()
        )
        .forEach((h) => {
          const time = apconso[h].periode.getTime()
          output_apart[time] = output_apart[time] || {}
          output_apart[time].apart_heures_consommees =
            (output_apart[time].apart_heures_consommees || 0) +
            apconso[h].heure_consomme
          output_apart[time].apart_motif_recours =
            apdemande[apart[k].demande].motif_recours_se
        })

      // Heures consommees cumulees sur la demande
      const series = f.generatePeriodSerie(
        apart[k].periode_debut,
        apart[k].periode_fin
      )
      series.reduce((accu, date) => {
        const time = date.getTime()

        //output_apart est déjà défini pour les heures autorisées
        accu = accu + (output_apart[time].apart_heures_consommees || 0)
        output_apart[time].apart_heures_consommees_cumulees = accu
        return accu
      }, 0)
    }
  })

  //Object.keys(output_apart).forEach(time => {
  //  if (output_effectif && time in output_effectif){
  //    output_apart[time].ratio_apart = (output_apart[time].apart_heures_consommees || 0) / (output_effectif[time].effectif * 157.67)
  //    //nbr approximatif d'heures ouvrées par mois
  //  }
  //})
  return output_apart
}
