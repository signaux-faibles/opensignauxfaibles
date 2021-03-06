import { f } from "./functions"
import { EntréeDebit } from "../GeneratedTypes"
import { ParHash, Timestamp } from "../RawDataTypes"

type AccuItem = {
  hash: string
  numero_historique: number
  date_traitement: Date
}

type DetteItem = {
  periode: Date
  part_ouvriere: number
  part_patronale: number
  montant_majorations: number
}

export type SortieDebit = {
  part_ouvriere: number
  part_patronale: number
  montant_majorations: number
  periode: Date
}

// Paramètres globaux utilisés par "public"
declare const date_fin: Date
declare const serie_periode: Date[]

export function debits(vdebit: ParHash<EntréeDebit> = {}): SortieDebit[] {
  const last_treatment_day = 20
  const ecn = {} as Record<string, AccuItem[]>
  for (const [h, debit] of Object.entries(vdebit)) {
    const start = debit.periode.start
    const end = debit.periode.end
    const num_ecn = debit.numero_ecart_negatif
    const compte = debit.numero_compte
    const key = start + "-" + end + "-" + num_ecn + "-" + compte
    ecn[key] = (ecn[key] || []).concat([
      {
        hash: h,
        numero_historique: debit.numero_historique,
        date_traitement: debit.date_traitement,
      },
    ])
  }

  for (const ecnItem of Object.values(ecn)) {
    ecnItem.sort(f.compareDebit)
    const l = ecnItem.length
    ecnItem.forEach((e, idx) => {
      if (idx <= l - 2) {
        const hashedDataInVDebit = vdebit[e?.hash]
        const next = ecnItem[idx + 1]
        if (hashedDataInVDebit !== undefined && next !== undefined) {
          hashedDataInVDebit.debit_suivant = next.hash
        }
      }
    })
  }

  const value_dette: Record<Timestamp, DetteItem[]> = {}

  for (const debit of Object.values(vdebit)) {
    const nextDate =
      (debit.debit_suivant && vdebit[debit.debit_suivant]?.date_traitement) ||
      date_fin

    //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
    const jour_traitement = debit.date_traitement.getUTCDate()
    const jour_traitement_suivant = nextDate.getUTCDate()
    let date_traitement_debut
    if (jour_traitement <= last_treatment_day) {
      date_traitement_debut = new Date(
        Date.UTC(
          debit.date_traitement.getFullYear(),
          debit.date_traitement.getUTCMonth()
        )
      )
    } else {
      date_traitement_debut = new Date(
        Date.UTC(
          debit.date_traitement.getFullYear(),
          debit.date_traitement.getUTCMonth() + 1
        )
      )
    }

    let date_traitement_fin
    if (jour_traitement_suivant <= last_treatment_day) {
      date_traitement_fin = new Date(
        Date.UTC(nextDate.getFullYear(), nextDate.getUTCMonth())
      )
    } else {
      date_traitement_fin = new Date(
        Date.UTC(nextDate.getFullYear(), nextDate.getUTCMonth() + 1)
      )
    }

    const periode_debut = date_traitement_debut
    const periode_fin = date_traitement_fin

    //generatePeriodSerie exlue la dernière période
    f.generatePeriodSerie(periode_debut, periode_fin).map((date) => {
      const time = date.getTime()
      value_dette[time] = (value_dette[time] || []).concat([
        {
          periode: debit.periode.start,
          part_ouvriere: debit.part_ouvriere,
          part_patronale: debit.part_patronale,
          montant_majorations: /*debit.montant_majorations ||*/ 0, // TODO: montant_majorations n'est pas fourni par les fichiers debit de l'urssaf pour l'instant, mais on aimerait y avoir accès un jour.
        },
      ])
    })
  }

  return serie_periode.map((p) =>
    (value_dette[p.getTime()] || []).reduce(
      (m, c) => {
        m.part_ouvriere += c.part_ouvriere
        m.part_patronale += c.part_patronale
        m.montant_majorations += c.montant_majorations
        return m
      },
      {
        part_ouvriere: 0,
        part_patronale: 0,
        montant_majorations: 0,
        periode: f.dateAddDay(f.dateAddMonth(p, 1), -1),
      } as SortieDebit
    )
  )
}
