import { SortieEffectif } from "./effectifs"
import { SortieDebit } from "./debits"

// Paramètres globaux utilisés par "public"
declare const serie_periode: Date[]

export type SortieJoinUrssaf = {
  effectif: (number | null)[]
  part_patronale: number[]
  part_ouvriere: number[]
  montant_majorations: number[]
}

export function joinUrssaf(
  effectif: SortieEffectif[],
  debit: SortieDebit[]
): SortieJoinUrssaf {
  const result: SortieJoinUrssaf = {
    effectif: [],
    part_patronale: [],
    part_ouvriere: [],
    montant_majorations: [],
  }

  debit.forEach((d, i) => {
    const e = effectif.filter(
      (e) => serie_periode[i].getTime() === e.periode.getTime()
    )
    if (e.length > 0) {
      result.effectif.push(e[0].effectif)
    } else {
      result.effectif.push(null)
    }
    result.part_patronale.push(d.part_patronale)
    result.part_ouvriere.push(d.part_ouvriere)
    result.montant_majorations.push(d.montant_majorations)
  })

  return result
}
