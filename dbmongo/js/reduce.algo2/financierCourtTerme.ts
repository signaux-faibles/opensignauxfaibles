function financierCourtTerme(diane) {
  "use strict"
  if (
    "concours_bancaire_courant" in diane &&
    diane["concours_bancaire_courant"] !== null &&
    "ca" in diane &&
    diane["ca"] !== null &&
    diane["ca"] != 0
  ) {
    return (diane["concours_bancaire_courant"] / diane["ca"]) * 100
  } else {
    return null
  }
}
