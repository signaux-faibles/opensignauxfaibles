function tauxMarge(diane) {
    "use strict"
    if (
        "excedent_brut_d_exploitation" in diane &&
        diane["excedent_brut_d_exploitation"] !== null &&
        "valeur_ajoutee" in diane &&
        diane["valeur_ajoutee"] !== null &&
        diane["excedent_brut_d_exploitation"] != 0
    ) {
        return (
            (diane["excedent_brut_d_exploitation"] / diane["valeur_ajoutee"]) *
            100
        )
    } else {
        return null
    }
}
