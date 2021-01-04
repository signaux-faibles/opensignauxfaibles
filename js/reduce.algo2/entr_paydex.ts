import { ParHash, ParPériode, EntréePaydex } from "../RawDataTypes"
import { f } from "./functions"

export type SortiePaydex = {
  [K in `paydex_nb_jours${"" | "_past_1" | "_past_12"}`]: number | null
}

export function entr_paydex(
  vPaydex: ParHash<EntréePaydex>,
  sériePériode: Date[]
): ParPériode<SortiePaydex> {
  "use strict"
  const paydexParPériode: ParPériode<SortiePaydex> = {}
  // initialisation (avec valeurs N/A par défaut)
  for (const période of sériePériode) {
    paydexParPériode[période.getTime()] = {
      paydex_nb_jours: null,
      paydex_nb_jours_past_1: null,
      paydex_nb_jours_past_12: null,
    }
  }
  // population des valeurs
  for (const entréePaydex of Object.values(vPaydex)) {
    const période = Date.UTC(
      entréePaydex.date_valeur.getUTCFullYear(),
      entréePaydex.date_valeur.getUTCMonth(),
      1
    )
    const moisSuivant = f.dateAddMonth(new Date(période), 1).getTime()
    const annéeSuivante = f.dateAddMonth(new Date(période), 12).getTime()
    f.add(
      {
        [période]: { paydex_nb_jours: entréePaydex.nb_jours },
        [moisSuivant]: { paydex_nb_jours_past_1: entréePaydex.nb_jours },
        [annéeSuivante]: { paydex_nb_jours_past_12: entréePaydex.nb_jours },
      },
      paydexParPériode
    )
  }
  return paydexParPériode
}
