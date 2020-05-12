declare function debug(string) // supported by jsc, to print in stdout

//eslint-disable-next-line @typescript-eslint/no-unused-vars
function raison_sociale(
  denomination_unite_legale: string,
  nom_unite_legale: string,
  nom_usage_unite_legale = "",
  prenom1_unite_legale: string,
  prenom2_unite_legale: string,
  prenom3_unite_legale: string,
  prenom4_unite_legale: string
): string {
  "use strict"
  nom_usage_unite_legale = nom_usage_unite_legale + "/"
  const raison_sociale =
    denomination_unite_legale ||
    (
      nom_unite_legale +
      "*" +
      nom_usage_unite_legale +
      prenom1_unite_legale +
      " " +
      (prenom2_unite_legale || "") +
      " " +
      (prenom3_unite_legale || "") +
      " " +
      (prenom4_unite_legale || "") +
      " "
    ).trim() + "/"

  return raison_sociale
}
