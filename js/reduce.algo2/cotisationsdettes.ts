import { f } from "./functions"
import { ParPériode } from "../common/makePeriodeMap"
import { EntréeCotisation, EntréeDebit } from "../GeneratedTypes"
import { Timestamp, ParHash } from "../RawDataTypes"

// Champs de EntréeCotisation nécéssaires à cotisationsdettes
type ChampsEntréeCotisation = Pick<EntréeCotisation, "periode" | "du">

// Champs de EntréeDebit nécéssaires à cotisationsdettes
type ChampsEntréeDebit = Pick<
  EntréeDebit,
  | "numero_compte"
  | "periode"
  | "part_ouvriere"
  | "part_patronale"
  | "numero_ecart_negatif"
  | "numero_historique"
  | "date_traitement"
  | "debit_suivant"
>

type EcartNegatif = {
  hash: string
  numero_historique: EntréeDebit["numero_historique"]
  date_traitement: EntréeDebit["date_traitement"]
}

type Dette = {
  periode: EntréeDebit["periode"]["start"]
  part_ouvriere: EntréeDebit["part_ouvriere"]
  part_patronale: EntréeDebit["part_patronale"]
}

type MonthOffset = 1 | 2 | 3 | 6 | 12
type CotisationsDettesPassees = {
  [K in
    | `montant_part_ouvriere_past_${MonthOffset}`
    | `montant_part_patronale_past_${MonthOffset}`]: number
}

export type SortieCotisationsDettes = {
  /** Règle métier URSSAF. true: si l'entreprise n'a pas eu de débit (dette) sur les 6 derniers mois. Pas utile dans les travaux de data science. */
  interessante_urssaf: boolean
  /** montant (€) des mensualités de règlement des cotisations sociales */
  cotisation: number
  /** montant (€) de la dette imputable au réglement des cotisatisations sociales des employés */
  montant_part_ouvriere: number
  /** montant (€) de la dette imputable au réglement des cotisatisations sociales des dirigeants */
  montant_part_patronale: number
} & CotisationsDettesPassees

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "cotisationsdettes"
  computed: SortieCotisationsDettes
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

/**
 * Calcule les variables liées aux cotisations sociales et dettes sur ces
 * cotisations.
 */
export function cotisationsdettes(
  vCotisation: ParHash<ChampsEntréeCotisation>,
  vDebit: ParHash<ChampsEntréeDebit>,
  periodes: Timestamp[],
  finPériode: Date // correspond à la variable globale date_fin
): ParPériode<SortieCotisationsDettes> {
  "use strict"

  // Tous les débits traitées après ce jour du mois sont reportées à la période suivante
  // Permet de s'aligner avec le calendrier de fourniture des données
  const lastAccountedDay = 20

  const sortieCotisationsDettes = f.makePeriodeMap<SortieCotisationsDettes>()

  const value_cotisation = f.makePeriodeMap<number[]>()

  // Répartition des cotisations sur toute la période qu'elle concerne
  for (const cotisation of Object.values(vCotisation)) {
    const periode_cotisation = f.generatePeriodSerie(
      cotisation.periode.start,
      cotisation.periode.end
    )
    periode_cotisation.forEach((date_cotisation) => {
      value_cotisation.set(
        date_cotisation,
        (value_cotisation.get(date_cotisation) || []).concat([
          cotisation.du / periode_cotisation.length,
        ])
      )
    })
  }

  // relier les débits
  // ecn: ecart negatif
  // map les débits: clé fabriquée maison => [{hash, numero_historique, date_traitement}, ...]
  // Pour un même compte, les débits avec le même num_ecn (chaque émission de facture) sont donc regroupés
  const ecn: ParHash<EcartNegatif[]> = {}
  for (const [h, debit] of Object.entries(vDebit)) {
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

  // Pour chaque numero_ecn, on trie et on chaîne les débits avec debit_suivant
  for (const ecnEntry of Object.values(ecn)) {
    ecnEntry.sort(f.compareDebit)
    const l = ecnEntry.length
    ecnEntry
      .filter((_, idx) => idx <= l - 2)
      .forEach((e, idx) => {
        const vDebitForHash = vDebit[e.hash]
        const next = (ecnEntry?.[idx + 1] || {}).hash
        if (vDebitForHash && next !== undefined)
          vDebitForHash.debit_suivant = next
      })
  }

  const value_dette = f.makePeriodeMap<Dette[]>()
  // Pour chaque objet debit:
  // debit_traitement_debut => periode de traitement du débit
  // debit_traitement_fin => periode de traitement du debit suivant, ou bien finPériode
  // Entre ces deux dates, c'est cet objet qui est le plus à jour.
  for (const debit of Object.values(vDebit)) {
    const nextDate =
      (debit.debit_suivant && vDebit[debit.debit_suivant]?.date_traitement) ||
      finPériode

    //Selon le jour du traitement, cela passe sur la période en cours ou sur la suivante.
    const jour_traitement = debit.date_traitement.getUTCDate()
    const jour_traitement_suivant = nextDate.getUTCDate()
    let date_traitement_debut
    if (jour_traitement <= lastAccountedDay) {
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
    if (jour_traitement_suivant <= lastAccountedDay) {
      date_traitement_fin = new Date(
        Date.UTC(nextDate.getFullYear(), nextDate.getUTCMonth())
      )
    } else {
      date_traitement_fin = new Date(
        Date.UTC(nextDate.getFullYear(), nextDate.getUTCMonth() + 1)
      )
    }

    //f.generatePeriodSerie exlue la dernière période
    f.generatePeriodSerie(date_traitement_debut, date_traitement_fin).forEach(
      (date) => {
        value_dette.set(date, [
          ...(value_dette.get(date) ?? []),
          {
            periode: debit.periode.start,
            part_ouvriere: debit.part_ouvriere,
            part_patronale: debit.part_patronale,
          },
        ])
      }
    )
  }

  // TODO faire numero de compte ailleurs
  // Array des numeros de compte
  //var numeros_compte = Array.from(new Set(
  //  Object.keys(vCotisation).map(function (h) {
  //    return(vCotisation[h].numero_compte)
  //  })
  //))

  periodes.forEach(function (time) {
    const val =
      sortieCotisationsDettes.get(time) ?? ({} as SortieCotisationsDettes)
    //val.numero_compte_urssaf = numeros_compte
    const valueCotis = value_cotisation.get(time)
    if (valueCotis !== undefined) {
      // somme de toutes les cotisations dues pour une periode donnée
      val.cotisation = valueCotis.reduce((a, cot) => a + cot, 0)
    }

    // somme de tous les débits (part ouvriere, part patronale)
    val.montant_part_ouvriere = (value_dette.get(time) || []).reduce(
      (acc, { part_ouvriere }) => acc + part_ouvriere,
      0
    )
    val.montant_part_patronale = (value_dette.get(time) || []).reduce(
      (acc, { part_patronale }) => acc + part_patronale,
      0
    )
    sortieCotisationsDettes.set(time, val)

    const monthOffsets: MonthOffset[] = [1, 2, 3, 6, 12]
    const futureTimestamps = monthOffsets
      .map((offset) => ({
        offset,
        timestamp: f.dateAddMonth(new Date(time), offset).getTime(),
      }))
      .filter(({ timestamp }) => periodes.includes(timestamp))

    futureTimestamps.forEach(({ offset, timestamp }) => {
      sortieCotisationsDettes.set(timestamp, {
        ...(sortieCotisationsDettes.get(timestamp) ??
          ({} as SortieCotisationsDettes)),
        [`montant_part_ouvriere_past_${offset}`]: val.montant_part_ouvriere,
        [`montant_part_patronale_past_${offset}`]: val.montant_part_patronale,
      })
    })

    if (val.montant_part_ouvriere + val.montant_part_patronale > 0) {
      const futureTimestamps = [0, 1, 2, 3, 4, 5]
        .map((offset) => ({
          timestamp: f.dateAddMonth(new Date(time), offset).getTime(),
        }))
        .filter(({ timestamp }) => periodes.includes(timestamp))

      futureTimestamps.forEach(({ timestamp }) => {
        sortieCotisationsDettes.set(timestamp, {
          ...(sortieCotisationsDettes.get(timestamp) ??
            ({} as SortieCotisationsDettes)),
          interessante_urssaf: false,
        })
      })
    }
  })

  return sortieCotisationsDettes
}
