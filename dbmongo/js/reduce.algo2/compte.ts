import "../globals.ts"
import { EntréeCompte } from "../RawDataTypes"

type SortieCompte = ParPériode<{ compte_urssaf: unknown }>

export function compte(compte: Record<string, EntréeCompte>): SortieCompte {
  "use strict"
  const output_compte: SortieCompte = {}

  //  var offset_compte = 3
  Object.keys(compte).forEach((hash) => {
    const periode: Periode = compte[hash].periode.getTime().toString()
    output_compte[periode] = output_compte[periode] || {}
    output_compte[periode].compte_urssaf = compte[hash].numero_compte
  })

  return output_compte
}
