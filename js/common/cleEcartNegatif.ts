import { EntréeDebit } from "../GeneratedTypes"

export function cleEcartNegatif(debit: EntréeDebit): string {
  const start = debit.periode.start
  const end = debit.periode.end
  const num_ecn = debit.numero_ecart_negatif
  const compte = debit.numero_compte
  return start + "-" + end + "-" + num_ecn + "-" + compte
}
