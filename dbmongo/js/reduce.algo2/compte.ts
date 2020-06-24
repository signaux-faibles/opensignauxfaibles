import "../globals.ts"

type Hash = string

export type V = {
  compte: Record<Hash, Compte>
}

type Periode = number

type TypeEnSortie = Record<
  Periode,
  {
    compte_urssaf: unknown
  }
>

export function compte(v: V): TypeEnSortie {
  "use strict"
  const output_compte: TypeEnSortie = {}

  //  var offset_compte = 3
  Object.keys(v.compte).forEach((hash) => {
    const periode: Periode = v.compte[hash].periode.getTime()

    output_compte[periode] = output_compte[periode] || {}
    output_compte[periode].compte_urssaf = v.compte[hash].numero_compte
  })

  return output_compte
}
