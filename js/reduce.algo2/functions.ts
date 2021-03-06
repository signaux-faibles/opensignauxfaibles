import { compareDebit } from "../common/compareDebit"
import { dateAddMonth } from "../common/dateAddMonth"
import { flatten } from "../common/flatten"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { omit } from "../common/omit"
import { procolToHuman } from "../common/procolToHuman"
import { raison_sociale } from "../common/raison_sociale"
import { region } from "../common/region"
import { add } from "./add"
import { apart } from "./apart"
import { ccsf } from "./ccsf"
import { cibleApprentissage } from "./cibleApprentissage"
import { compte } from "./compte"
import { cotisation } from "./cotisation"
import { cotisationsdettes } from "./cotisationsdettes"
import { defaillances } from "./defaillances"
import { delais } from "./delais"
import { detteFiscale } from "./detteFiscale"
import { effectifs } from "./effectifs"
import { entr_bdf } from "./entr_bdf"
import { entr_diane } from "./entr_diane"
import { entr_paydex } from "./entr_paydex"
import { entr_sirene } from "./entr_sirene"
import { fraisFinancier } from "./fraisFinancier"
import { lookAhead } from "./lookAhead"
import { nbDays } from "./nbDays"
import { outputs } from "./outputs"
import { poidsFrng } from "./poidsFrng"
import { populateNafAndApe } from "./populateNafAndApe"
import { repeatable } from "./repeatable"
import { sirene } from "./sirene"
import { makePeriodeMap } from "../common/makePeriodeMap"

export const f = {
  makePeriodeMap,
  flatten,
  outputs,
  apart,
  compte,
  effectifs,
  add,
  repeatable,
  delais,
  defaillances,
  cotisationsdettes,
  ccsf,
  sirene,
  populateNafAndApe,
  cotisation,
  cibleApprentissage,
  entr_sirene,
  dateAddMonth,
  generatePeriodSerie,
  poidsFrng,
  detteFiscale,
  fraisFinancier,
  entr_bdf,
  omit,
  entr_diane,
  entr_paydex,
  lookAhead,
  compareDebit,
  procolToHuman,
  raison_sociale,
  region,
  nbDays,
}
