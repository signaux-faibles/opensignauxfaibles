function region(departement){
  "use strict";
  var reg = ""
  switch (departement){
    case "01":
    case "03":
    case "07":
    case "15":
    case "26":
    case "38":
    case "42":
    case "43":
    case "63":
    case "69":
    case "69":
    case "73":
    case "74":
      reg = "Auvergne-Rhône-Alpes"
      break
    case "02":
    case "59":
    case "60":
    case "62":
    case "80":
      reg = "Hauts-de-France"
      break
    case "04":
    case "05":
    case "06":
    case "13":
    case "83":
    case "84":
      reg = "Provence-Alpes-Côte d'Azur"
      break
    case "08":
    case "10":
    case "51":
    case "52":
    case "54":
    case "55":
    case "57":
    case "67":
    case "68":
    case "88":
      reg = "Grand Est"
      break
    case "09":
    case "11":
    case "12":
    case "30":
    case "31":
    case "32":
    case "34":
    case "46":
    case "48":
    case "65":
    case "66":
    case "81":
    case "82":
      reg = "Occitanie"
      break
    case "14":
    case "27":
    case "50":
    case "61":
    case "76":
      reg = "Normandie"
      break
    case "18":
    case "28":
    case "36":
    case "37":
    case "41":
    case "45":
      reg = "Centre-Val de Loire"
      break
    case "16":
    case "17":
    case "19":
    case "23":
    case "24":
    case "33":
    case "40":
    case "47":
    case "64":
    case "79":
    case "86":
    case "87":
      reg = "Nouvelle-Aquitaine"
      break
    case "20":
      reg = "Corse"
      break
    case "21":
    case "25":
    case "39":
    case "58":
    case "70":
    case "71":
    case "89":
    case "90":
      reg = "Bourgogne-Franche-Comté"
      break
    case "22":
    case "29":
    case "35":
    case "56":
      reg = "Bretagne"
      break
    case "44":
    case "49":
    case "53":
    case "72":
    case "85":
      reg = "Pays de la Loire"
      break
    case "75":
    case "77":
    case "78":
    case "91":
    case "92":
    case "93":
    case "94":
    case "95":
      reg = "Île-de-France"
      break
  }
  return(reg)
}
