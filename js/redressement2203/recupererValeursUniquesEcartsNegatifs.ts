import { EntréeDebit } from "../GeneratedTypes"
import { f } from "./functions"

export function recupererValeursUniquesEcartsNegatifs(
  debits: EntréeDebit[]
): string[] {
  const ecartsNegatifs = debits.map((debit) => f.cleEcartNegatif(debit))
  return [...new Set(ecartsNegatifs)]
}
