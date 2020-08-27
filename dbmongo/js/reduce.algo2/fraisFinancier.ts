export type DianeProperty =
  | "interets"
  | "excedent_brut_d_exploitation"
  | "produits_financiers"
  | "produit_exceptionnel"
  | "charge_exceptionnelle"
  | "charges_financieres"

export type Diane = {
  [prop in DianeProperty]: number
}

export type DianePartial = {
  [prop in DianeProperty]: number | null
}

export function fraisFinancier(diane: DianePartial): number | null {
  "use strict"
  if (
    typeof diane["interets"] === "number" &&
    typeof diane["excedent_brut_d_exploitation"] === "number" &&
    typeof diane["produits_financiers"] === "number" &&
    typeof diane["charges_financieres"] === "number" &&
    typeof diane["charge_exceptionnelle"] === "number" &&
    typeof diane["produit_exceptionnel"] === "number" &&
    diane["excedent_brut_d_exploitation"] +
      diane["produits_financiers"] +
      diane["produit_exceptionnel"] -
      diane["charge_exceptionnelle"] -
      diane["charges_financieres"] !==
      0
  ) {
    return (
      (diane["interets"] /
        (diane["excedent_brut_d_exploitation"] +
          diane["produits_financiers"] +
          diane["produit_exceptionnel"] -
          diane["charge_exceptionnelle"] -
          diane["charges_financieres"])) *
      100
    )
  } else {
    return null
  }
}
