exports.__esModule = true;
function fraisFinancier(diane) {
    "use strict";
    if (("interets" in diane) && (diane["interets"] !== null) &&
        ("excedent_brut_d_exploitation" in diane) && (diane["excedent_brut_d_exploitation"] !== null) &&
        ("produits_financiers" in diane) && (diane["produits_financiers"] !== null) &&
        ("charges_financieres" in diane) && (diane["charges_financieres"] !== null) &&
        ("charge_exceptionnelle" in diane) && (diane["charge_exceptionnelle"] !== null) &&
        ("produit_exceptionnel" in diane) && (diane["produit_exceptionnel"] !== null) &&
        diane["excedent_brut_d_exploitation"] + diane["produits_financiers"] + diane["produit_exceptionnel"] - diane["charge_exceptionnelle"] - diane["charges_financieres"] != 0) {
        return diane["interets"] / (diane["excedent_brut_d_exploitation"] + diane["produits_financiers"] + diane["produit_exceptionnel"] -
            diane["charge_exceptionnelle"] - diane["charges_financieres"]) * 100;
    }
    else {
        return null;
    }
}
exports.fraisFinancier = fraisFinancier;
