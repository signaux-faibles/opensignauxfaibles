import { f } from "./functions"
import { ParPériode } from "../common/makePeriodeMap"
import { EntréePaydex } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"

export type SortiePaydex = {
  /** Nombre de jours de retard de paiement moyen, basé sur trois expériences de paiement minimum (provenant de trois fournisseurs distincts). */
  paydex_nb_jours: number | null
  paydex_nb_jours_past_3: number | null
  paydex_nb_jours_past_6: number | null
  paydex_nb_jours_past_12: number | null
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "entr_paydex"
  computed: Omit<SortiePaydex, "paydex_nb_jours">
  transmitted: Pick<SortiePaydex, "paydex_nb_jours">
}

export function entr_paydex(
  vPaydex: ParHash<EntréePaydex>,
  sériePériode: Date[]
): ParPériode<SortiePaydex> {
  "use strict"
  const paydexParPériode = f.makePeriodeMap<SortiePaydex>()
  // initialisation (avec valeurs N/A par défaut)
  for (const période of sériePériode) {
    paydexParPériode.set(période, {
      paydex_nb_jours: null,
      paydex_nb_jours_past_3: null,
      paydex_nb_jours_past_6: null,
      paydex_nb_jours_past_12: null,
    })
  }
  // population des valeurs
  for (const entréePaydex of Object.values(vPaydex)) {
    const période = Date.UTC(
      entréePaydex.date_valeur.getUTCFullYear(),
      entréePaydex.date_valeur.getUTCMonth(),
      1
    )
    const mois3Suivant = f.dateAddMonth(new Date(période), 3).getTime()
    const mois6Suivant = f.dateAddMonth(new Date(période), 6).getTime()
    const annéeSuivante = f.dateAddMonth(new Date(période), 12).getTime()
    const donnéesAdditionnelles = f.makePeriodeMap<Partial<SortiePaydex>>([
      [période, { paydex_nb_jours: entréePaydex.nb_jours }],
      [mois3Suivant, { paydex_nb_jours_past_3: entréePaydex.nb_jours }],
      [mois6Suivant, { paydex_nb_jours_past_6: entréePaydex.nb_jours }],
      [annéeSuivante, { paydex_nb_jours_past_12: entréePaydex.nb_jours }],
    ])
    f.add(donnéesAdditionnelles, paydexParPériode)
  }
  return paydexParPériode
}
