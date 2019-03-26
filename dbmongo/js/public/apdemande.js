function apdemande(apdemande) {
  return f.iterable(apdemande).sort((p1, p2) => p1.periode < p2.periode)
}