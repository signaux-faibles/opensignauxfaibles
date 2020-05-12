function raison_sociale /*eslint-disable-line @typescript-eslint/no-unused-vars */(denomination_unite_legale, nom_unite_legale, nom_usage_unite_legale, prenom1_unite_legale, prenom2_unite_legale, prenom3_unite_legale, prenom4_unite_legale) {
    "use strict";
    if (nom_usage_unite_legale === void 0) { nom_usage_unite_legale = ""; }
    nom_usage_unite_legale = nom_usage_unite_legale + "/";
    var raison_sociale = denomination_unite_legale ||
        (nom_unite_legale +
            "*" +
            nom_usage_unite_legale +
            prenom1_unite_legale +
            " " +
            (prenom2_unite_legale || "") +
            " " +
            (prenom3_unite_legale || "") +
            " " +
            (prenom4_unite_legale || "") +
            " ").trim() + "/";
    return raison_sociale;
}
