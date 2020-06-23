export function compareDebit(
  a: { numero_historique: number },
  b: { numero_historique: number }
): number {
  "use strict"
  if (a.numero_historique < b.numero_historique) return -1
  if (a.numero_historique > b.numero_historique) return 1
  return 0
}
