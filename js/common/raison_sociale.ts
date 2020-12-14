export function raison_sociale /*eslint-disable-line @typescript-eslint/no-unused-vars */(
  denomination_unite_legale?: string | null,
  nom_unite_legale?: string | null,
  nom_usage_unite_legale?: string | null,
  prenom1_unite_legale?: string | null,
  prenom2_unite_legale?: string | null,
  prenom3_unite_legale?: string | null,
  prenom4_unite_legale?: string | null
): string {
  "use strict"
  const nomUsageUniteLegale = nom_usage_unite_legale
    ? nom_usage_unite_legale + "/"
    : ""
  const raison_sociale =
    denomination_unite_legale ||
    (
      nom_unite_legale +
      "*" +
      nomUsageUniteLegale +
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
