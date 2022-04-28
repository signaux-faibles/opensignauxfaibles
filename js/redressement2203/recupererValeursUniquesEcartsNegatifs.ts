import { EntréeDebit } from "../GeneratedTypes"

export function recupererValeursUniquesEcartsNegatifs(
  debits: EntréeDebit[]
): string[] {
  const ecartsNegatifs = debits.map((debit) => debit.numero_ecart_negatif)
  return [...new Set(ecartsNegatifs)]
}
