function compareDebit (a,b) {
  if (a.numero_historique < b.numero_historique) return -1
  if (a.numero_historique > b.numero_historique) return 1
  return 0
}