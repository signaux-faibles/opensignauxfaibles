type Region =
  | "Auvergne-Rhône-Alpes"
  | "Hauts-de-France"
  | "Provence-Alpes-Côte d'Azur"
  | "Grand Est"
  | "Occitanie"
  | "Normandie"
  | "Centre-Val de Loire"
  | "Nouvelle-Aquitaine"
  | "Corse"
  | "Bourgogne-Franche-Comté"
  | "Bretagne"
  | "Pays de la Loire"
  | "Île-de-France"

export function region(departement: Departement): Region | "" {
  "use strict"
  let reg: Region | "" = ""
  switch (departement) {
    case "01" ||
      "03" ||
      "07" ||
      "15" ||
      "26" ||
      "38" ||
      "42" ||
      "43" ||
      "63" ||
      "69" ||
      "73" ||
      "74":
      reg = "Auvergne-Rhône-Alpes"
      break
    case "02" || "59" || "60" || "62" || "80":
      reg = "Hauts-de-France"
      break
    case "04" || "05" || "06" || "13" || "83" || "84":
      reg = "Provence-Alpes-Côte d'Azur"
      break
    case "08" ||
      "10" ||
      "51" ||
      "52" ||
      "54" ||
      "55" ||
      "57" ||
      "67" ||
      "68" ||
      "88":
      reg = "Grand Est"
      break
    case "09" ||
      "11" ||
      "12" ||
      "30" ||
      "31" ||
      "32" ||
      "34" ||
      "46" ||
      "48" ||
      "65" ||
      "66" ||
      "81" ||
      "82":
      reg = "Occitanie"
      break
    case "14" || "27" || "50" || "61" || "76":
      reg = "Normandie"
      break
    case "18" || "28" || "36" || "37" || "41" || "45":
      reg = "Centre-Val de Loire"
      break
    case "16" ||
      "17" ||
      "19" ||
      "23" ||
      "24" ||
      "33" ||
      "40" ||
      "47" ||
      "64" ||
      "79" ||
      "86" ||
      "87":
      reg = "Nouvelle-Aquitaine"
      break
    case "20":
      reg = "Corse"
      break
    case "21" || "25" || "39" || "58" || "70" || "71" || "89" || "90":
      reg = "Bourgogne-Franche-Comté"
      break
    case "22" || "29" || "35" || "56":
      reg = "Bretagne"
      break
    case "44" || "49" || "53" || "72" || "85":
      reg = "Pays de la Loire"
      break
    case "75" || "77" || "78" || "91" || "92" || "93" || "94" || "95":
      reg = "Île-de-France"
      break
  }
  return reg
}
