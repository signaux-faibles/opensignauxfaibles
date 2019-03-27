function bdf(hs) {
  return f.iterable(hs).sort((a, b) => a.annee_bdf < b.annee_bdf)
}