import { Diane } from "./fraisFinancier"

export function poidsFrng(diane: Diane): number | null {
  "use strict"
  if ("couverture_ca_fdr" in diane && diane["couverture_ca_fdr"] !== null) {
    return (diane["couverture_ca_fdr"] / 360) * 100
  } else {
    return null
  }
}
