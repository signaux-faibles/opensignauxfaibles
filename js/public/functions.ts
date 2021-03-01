import { makePeriodeMap } from "../common/makePeriodeMap"
import { compareDebit } from "../common/compareDebit"
import { dateAddMonth } from "../common/dateAddMonth"
import { flatten } from "../common/flatten"
import { raison_sociale } from "../common/raison_sociale"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { omit } from "../common/omit"
import { procolToHuman } from "../common/procolToHuman"
import { apconso } from "./apconso"
import { apdemande } from "./apdemande"
import { bdf } from "./bdf"
import { compte } from "./compte"
import { cotisations } from "./cotisations"
import { dateAddDay } from "./dateAddDay"
import { debits } from "./debits"
import { delai } from "./delai"
import { diane } from "./diane"
import { effectifs } from "./effectifs"
import { joinUrssaf } from "./joinUrssaf"
import { sirene } from "./sirene"

export const f = {
  newParPériode: makePeriodeMap,
  dateAddMonth,
  debits,
  apconso,
  apdemande,
  flatten,
  compte,
  effectifs,
  delai,
  sirene,
  cotisations,
  dateAddDay,
  omit,
  generatePeriodSerie,
  diane,
  bdf,
  joinUrssaf,
  procolToHuman,
  compareDebit,
  raison_sociale,
}
