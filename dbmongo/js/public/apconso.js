function apconso(apconso) {
  "use strict";
  return f.iterable(apconso).sort((p1, p2) => p1.periode < p2.periode)
}