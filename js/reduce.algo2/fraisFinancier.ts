import { EntréeDiane } from "../RawDataTypes"

export type ChampsDiane = Pick<
  EntréeDiane,
  | "interets"
  | "excedent_brut_d_exploitation"
  | "produits_financiers"
  | "produit_exceptionnel"
  | "charge_exceptionnelle"
  | "charges_financieres"
>

export function fraisFinancier(diane: ChampsDiane): number | null {
  "use strict"
  const ratio =
    (diane["interets"] ?? NaN) /
    ((diane["excedent_brut_d_exploitation"] ?? NaN) +
      (diane["produits_financiers"] ?? NaN) +
      (diane["produit_exceptionnel"] ?? NaN) -
      (diane["charge_exceptionnelle"] ?? NaN) -
      (diane["charges_financieres"] ?? NaN))
  return isNaN(ratio) ? null : ratio * 100
}
