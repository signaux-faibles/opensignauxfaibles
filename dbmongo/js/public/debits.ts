import { dateAddDay } from "./dateAddDay"
import { compareDebit } from "../common/compareDebit"
import { dateAddMonth } from "../common/dateAddMonth"
import { generatePeriodSerie } from "../common/generatePeriodSerie"

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
  periode?: Date
}

// Paramètres globaux utilisés par "public"
declare let date_fin: Date
declare let serie_periode: Date[]

export function debits(
  vdebit: Record<DataHash, EntréeDebit> = {}
): SortieDebit[] {
  const f = { compareDebit, generatePeriodSerie, dateAddMonth, dateAddDay } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  const last_treatment_day = 20
  const ecn = Object.keys(vdebit).reduce((accu, h) => {
    const debit = vdebit[h]
    const start = debit.periode.start
    const end = debit.periode.end
    const num_ecn = debit.numero_ecart_negatif
    const compte = debit.numero_compte
    const key = start + "-" + end + "-" + num_ecn + "-" + compte
    accu[key] = (accu[key] || []).concat([
      {
        hash: h,
        numero_historique: debit.numero_historique,
        date_traitement: debit.date_traitement,
      },
    ])
    return accu
  }, {} as Record<string, AccuItem[]>)

  Object.keys(ecn).forEach((i) => {
    ecn[i].sort(f.compareDebit)
    const l = ecn[i].length
    ecn[i].forEach((e, idx) => {
      if (idx <= l - 2) {
        vdebit[e.hash].debit_suivant = ecn[i][idx + 1].hash
      }
    })
  })

  const value_dette: Record<number, DetteItem[]> = {}

  Object.keys(vdebit).forEach(function (h) {
    const debit = vdebit[h]

    const debit_suivant = vdebit[debit.debit_suivant] || {
      date_traitement: date_fin,
    }

    //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
    const jour_traitement = debit.date_traitement.getUTCDate()
    const jour_traitement_suivant = debit_suivant.date_traitement.getUTCDate()
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
        Date.UTC(
          debit_suivant.date_traitement.getFullYear(),
          debit_suivant.date_traitement.getUTCMonth()
        )
      )
    } else {
      date_traitement_fin = new Date(
        Date.UTC(
          debit_suivant.date_traitement.getFullYear(),
          debit_suivant.date_traitement.getUTCMonth() + 1
        )
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
          montant_majorations: debit.montant_majorations || 0,
        },
      ])
    })
  })

  return serie_periode.map((p) =>
    (value_dette[p.getTime()] || []).reduce(
      (m, c) => ({
        part_ouvriere: m.part_ouvriere + c.part_ouvriere,
        part_patronale: m.part_patronale + c.part_patronale,
        periode: f.dateAddDay(f.dateAddMonth(p, 1), -1),
      }),
      { part_ouvriere: 0, part_patronale: 0 } as SortieDebit
    )
  )
}
