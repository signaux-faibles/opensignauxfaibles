function bdf(hs) {
  "use strict";
  return f.iterable(hs).sort((a, b) => a.annee_bdf < b.annee_bdf)
}