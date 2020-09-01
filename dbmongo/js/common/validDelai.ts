// Valeurs attendues par delais(), pour chaque période. (cf dbmongo/lib/urssaf/delai.go)
export type EntréeDelai = {
  date_creation: Date
  date_echeance: Date
  duree_delai: number // nombre de jours entre date_creation et date_echeance
  montant_echeancier: number // exprimé en euros
}

export function validDelai(delai: EntréeDelai): void {
  const règles: Array<(delai: EntréeDelai) => boolean> = [
    ({ duree_delai }) => duree_delai > 0,
  ]
  if (!règles[0](delai)) {
    throw new Error(`delai invalide, règle: ${règles[0].toString()}`)
  }
}
