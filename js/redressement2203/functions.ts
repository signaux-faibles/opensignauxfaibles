import { flatten } from "../common/flatten"
import { recupererDetteTotale } from "./recupererDetteTotale"
import { recupererValeursUniquesEcartsNegatifs } from "./recupererValeursUniquesEcartsNegatifs"
import { cleEcartNegatif } from "../common/cleEcartNegatif"
import { cotisation } from "./cotisation"
import { outputs } from "./outputs"
import { makePeriodeMap } from "../common/makePeriodeMap"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "../common/dateAddMonth"

export const f = {
  flatten,
  recupererDetteTotale,
  recupererValeursUniquesEcartsNegatifs,
  cleEcartNegatif,
  makePeriodeMap,
  outputs,
  cotisation,
  generatePeriodSerie,
  dateAddMonth,
}
