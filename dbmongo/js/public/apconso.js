function apconso(apconso) {
  return f.iterable(apconso).sort((p1, p2) => p1.periode < p2.periode)
}