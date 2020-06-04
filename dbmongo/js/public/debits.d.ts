export function debits(vdebit: {
  [h: string]: {
    periode: { start: Date; end: Date }
    numero_ecart_negatif: any
    numero_compte: any
    numero_historique: any
    date_traitement: any
  }
}): {
  part_ouvriere: number
  part_patronale: number
  periode: Periode
}[]
