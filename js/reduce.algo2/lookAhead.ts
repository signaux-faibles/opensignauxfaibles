import { f } from "./functions"
import { ParPériode } from "../common/makePeriodeMap"
import { SortieDefaillances } from "./defaillances"
import { SortieCotisation } from "./cotisation"
import { Timestamp } from "../RawDataTypes"

export type Outcome = {
  /** Distance de l'évènement, exprimé en nombre de périodes. */
  time_til_outcome: number
  outcome: boolean
}

type EntréeLookAhead = {
  outcome?: Outcome["outcome"]
  tag_default?: SortieCotisation["tag_default"]
  tag_failure?: SortieDefaillances["tag_failure"]
}

export function lookAhead(
  data: ParPériode<EntréeLookAhead>,
  attr_name: keyof EntréeLookAhead, // "outcome" | "tag_default" | "tag_failure",
  n_months: number,
  past: boolean
): ParPériode<Outcome> {
  "use strict"
  // Est-ce que l'évènement se répercute dans le passé (past = true on pourra se
  // demander: que va-t-il se passer) ou dans le future (past = false on
  // pourra se demander que s'est-il passé
  const chronologic = (pérA: Timestamp, pérB: Timestamp) => pérA - pérB
  const reverse = (pérA: Timestamp, pérB: Timestamp) => pérB - pérA

  let counter = -1
  const output = [...data.keys()]
    .sort(past ? reverse : chronologic)
    .reduce((m, période) => {
      // Si on a déjà détecté quelque chose, on compte le nombre de périodes
      if (counter >= 0) counter = counter + 1

      if (data.get(période)?.[attr_name]) {
        // si l'évènement se produit on retombe à 0
        counter = 0
      }

      if (counter >= 0) {
        // l'évènement s'est produit
        m.set(période, {
          time_til_outcome: counter,
          outcome: counter <= n_months ? true : false,
        })
      }
      return m
    }, f.makePeriodeMap<Outcome>())

  return output
}
