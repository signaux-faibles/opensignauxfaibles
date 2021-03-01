import { f } from "./functions"
import { EntréeApConso, EntréeApDemande } from "../GeneratedTypes"
import { ParPériode } from "../common/makePeriodeMap"

type ApConsoHash = string

type ApDemandeHash = string

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "apart"
  computed: {
    /** Nombre d'heures d'activité partielle consommées sur la période considérée. */
    apart_heures_consommees: EntréeApConso["heure_consomme"]
    /** Cumul du nombre d'heures d'activité partielle consommées depuis date_debut. */
    apart_heures_consommees_cumulees: EntréeApConso["heure_consomme"]
  }
  transmitted: {
    /** Nombre total d'heures d'activité partielle autorisées (nombre décimal). */
    apart_heures_autorisees: EntréeApDemande["hta"]
    /** Motif de recours à l'activité partielle:
     * 1	Conjoncture économique.
     * 2	Difficultés d’approvisionnement en matières premières ou en énergie
     * 3	Sinistre ou intempéries de caractère exceptionnel
     * 4	Transformation, restructuration ou modernisation des installations et des bâtiments
     * 5	Autres circonstances exceptionnelles
     */
    apart_motif_recours: EntréeApDemande["motif_recours_se"]
  }
}

export type SortieAPart = Variables["computed"] & Variables["transmitted"]

export function apart(
  apconso: Record<ApConsoHash, EntréeApConso>,
  apdemande: Record<ApDemandeHash, EntréeApDemande>
): ParPériode<SortieAPart> {
  "use strict"

  const output_apart = f.newParPériode<SortieAPart>()

  // Mapping (pour l'instant vide) du hash de la demande avec les hash des consos correspondantes
  const apart: Record<
    string,
    {
      demande: ApDemandeHash
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
    series.forEach((période) => {
      output_apart.set(période, {
        ...(output_apart.get(période) ?? ({} as SortieAPart)),
        apart_heures_autorisees: apdemandeEntry.hta,
      })
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
          const période = apconso[h]?.periode
          if (période === undefined) {
            return
          }
          const current = output_apart.get(période) ?? ({} as SortieAPart)
          const heureConso = apconso[h]?.heure_consomme
          if (heureConso !== undefined) {
            current.apart_heures_consommees =
              (current.apart_heures_consommees ?? 0) + heureConso
          }
          const motifRecours = apdemande[apartEntry.demande]?.motif_recours_se
          if (motifRecours !== undefined) {
            current.apart_motif_recours = motifRecours
          }
          output_apart.set(période, current)
        })

      // Heures consommees cumulees sur la demande
      const series = f.generatePeriodSerie(
        apartEntry.periode_debut,
        apartEntry.periode_fin
      )
      series.reduce((accu, période) => {
        //output_apart est déjà défini pour les heures autorisées
        const current = output_apart.get(période) ?? ({} as SortieAPart)
        accu = accu + (current.apart_heures_consommees || 0)
        output_apart.set(période, {
          ...current,
          apart_heures_consommees_cumulees: accu,
        })
        // TODO: on pourrait ajouter une méthode append (ou upsert) à ParPériode() pour alléger la logique ci-dessus
        return accu
      }, 0)
    }
  }

  // Note: à la fin de l'opération map-reduce, sfdata va calculer la propriété
  // ratio_apart depuis apart.crossComputation.json.

  return output_apart
}
