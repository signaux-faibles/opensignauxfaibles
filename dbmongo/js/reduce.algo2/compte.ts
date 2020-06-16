import "../globals.ts"

type TypeEnSortie = Record<number, { compte_urssaf: unknown }>

export function compte(v: any): TypeEnSortie {
  "use strict"
  const output_compte: TypeEnSortie = {}

  //  var offset_compte = 3
  Object.keys(v.compte).forEach((hash) => {
    const periode: number = v.compte[hash].periode.getTime()

    output_compte[periode] = output_compte[periode] || {}
    output_compte[periode].compte_urssaf = v.compte[hash].numero_compte
  })

  return output_compte
}
