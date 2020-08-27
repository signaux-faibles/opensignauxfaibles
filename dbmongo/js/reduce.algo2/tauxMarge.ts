import "../globals"

export function tauxMarge(diane: EntréeDiane): number | null {
  "use strict"

  const ratio =
    (diane["excedent_brut_d_exploitation"] ?? NaN) /
    (diane["valeur_ajoutee"] ?? NaN)
  return isNaN(ratio) ? null : ratio * 100
}
