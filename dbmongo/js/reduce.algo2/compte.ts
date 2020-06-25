import "../globals.ts"

type Periode = number // TODO: réutiliser le type Periode global ?

type SortieCompte = Record<
  Periode,
  {
    compte_urssaf: unknown
  }
>

export function compte(v: DonnéesCompte): SortieCompte {
  "use strict"
  const output_compte: SortieCompte = {}

  //  var offset_compte = 3
  Object.keys(v.compte).forEach((hash) => {
    const periode: Periode = v.compte[hash].periode.getTime()

    output_compte[periode] = output_compte[periode] || {}
    output_compte[periode].compte_urssaf = v.compte[hash].numero_compte
  })

  return output_compte
}
