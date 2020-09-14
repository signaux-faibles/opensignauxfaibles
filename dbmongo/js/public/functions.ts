import { iterable } from "./iterable"
import { debits } from "./debits"
import { apconso } from "./apconso"
import { apdemande } from "./apdemande"
import { flatten } from "../common/flatten"
import { compte } from "./compte"
import { effectifs } from "./effectifs"
import { joinUrssaf } from "./joinUrssaf"
import { delai } from "./delai"
import { bdf } from "./bdf"
import { diane } from "./diane"
import { sirene } from "./sirene"
import { cotisations } from "./cotisations"
import { dateAddDay } from "./dateAddDay"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { omit } from "../common/omit"

export const f = {
  iterable,
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
}
