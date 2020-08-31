import "../globals"

export function poidsFrng(diane: EntréeDiane): number | null {
  "use strict"

  return typeof diane["couverture_ca_fdr"] === "number"
    ? (diane["couverture_ca_fdr"] / 360) * 100
    : null
}
