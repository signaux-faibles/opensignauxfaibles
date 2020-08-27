import "../globals"

export function detteFiscale(diane: EntréeDiane): number | null {
  "use strict"

  const ratio =
    (diane["dette_fiscale_et_sociale"] ?? NaN) /
    (diane["valeur_ajoutee"] ?? NaN)
  return isNaN(ratio) ? null : ratio * 100
}
