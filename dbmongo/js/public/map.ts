import { iterable } from "./iterable"
import { debits, SortieDebit } from "./debits"
import { apconso } from "./apconso"
import { apdemande } from "./apdemande"
import { flatten } from "../common/flatten"
import { compte } from "./compte"
import { effectifs } from "./effectifs"
import { joinUrssaf } from "./joinUrssaf"
import { delai } from "./delai"
import { bdf, Bdf } from "./bdf"
import { diane } from "./diane"
import { sirene } from "./sirene"
import { cotisations } from "./cotisations"
import { dateAddDay } from "./dateAddDay"
import { dealWithProcols, SortieProcols } from "./dealWithProcols"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { omit } from "../common/omit"
import {
  CompanyDataValues,
  BatchKey,
  EntréeApConso,
  EntréeApDemande,
  EntréeDelai,
  EntréeDefaillances,
  EntréeDiane,
  ParHash,
} from "../RawDataTypes"

type SortieMapCommon = {
  key: string
  batch: string
}

type SortieMapEtablissement = SortieMapCommon & {
  sirene: unknown
  debit: SortieDebit[]
  apconso: EntréeApConso[]
  apdemande: EntréeApDemande[]
  delai: EntréeDelai[]
  compte: unknown
  procol: ParHash<EntréeDefaillances>
  last_procol: SortieProcols
  periodes: Date[]
  effectif: (number | null)[]
  cotisation: number[]
  debit_part_patronale: number[]
  debit_part_ouvriere: number[]
  debit_montant_majorations: number[]
  idEntreprise: string
}

type SortieMapEntreprise = SortieMapCommon & {
  diane: EntréeDiane[]
  bdf: Bdf[]
  sirene_ul: unknown
  crp: unknown
}

export type SortieMap = SortieMapEtablissement | SortieMapEntreprise

// Paramètres globaux utilisés par "public"
declare let actual_batch: BatchKey
declare let serie_periode: Date[]

// Types de données en entrée et sortie
export type Input = { _id: unknown; value: CompanyDataValues }
export type OutKey = string
export type OutValue = Partial<SortieMap>
declare function emit(key: string, value: OutValue): void

export function map(this: Input): void {
  /* DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO */ const f = {
    ...{ iterable, debits, apconso, apdemande, flatten, compte, effectifs }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ delai, dealWithProcols, sirene, cotisations, dateAddDay, omit }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ generatePeriodSerie, diane, bdf, joinUrssaf }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  const value = f.flatten(this.value, actual_batch)

  if (this.value.scope === "etablissement") {
    const vcmde: Partial<SortieMapEtablissement> = {}
    vcmde.key = this.value.key
    vcmde.batch = actual_batch
    vcmde.sirene = f.sirene(f.iterable(value.sirene))
    vcmde.periodes = serie_periode
    const effectif = f.effectifs(value.effectif)
    const debit = f.debits(value.debit)
    const join = f.joinUrssaf(effectif, debit)
    vcmde.debit_part_patronale = join.part_patronale
    vcmde.debit_part_ouvriere = join.part_ouvriere
    vcmde.debit_montant_majorations = join.montant_majorations
    vcmde.effectif = join.effectif
    vcmde.cotisation = f.cotisations(value.cotisation)
    vcmde.apconso = f.apconso(value.apconso)
    vcmde.apdemande = f.apdemande(value.apdemande)
    vcmde.delai = f.delai(value.delai)
    vcmde.compte = f.compte(value.compte)
    vcmde.procol = undefined // Note: initialement, l'expression ci-dessous était affectée à vcmde.procol, puis écrasée plus bas. J'initialise quand même vcmde.procol ici pour ne pas faire échouer test-api.sh sur l'ordre des propriétés.
    const procol = [
      ...f.dealWithProcols(value.altares, "altares"),
      ...f.dealWithProcols(value.procol, "procol"),
    ]
    vcmde.last_procol = procol[procol.length - 1] || { etat: "in_bonis" }
    vcmde.idEntreprise = "entreprise_" + this.value.key.slice(0, 9)
    vcmde.procol = value.procol

    emit("etablissement_" + this.value.key, vcmde)
  } else if (this.value.scope === "entreprise") {
    const v: Partial<SortieMapEntreprise> = {}
    const diane = f.diane(value.diane)
    const bdf = f.bdf(value.bdf)
    const sirene_ul = (value.sirene_ul || {})[
      Object.keys(value.sirene_ul || {})[0] || ""
    ]
    const crp = value.crp
    v.key = this.value.key
    v.batch = actual_batch

    if (diane.length > 0) {
      v.diane = diane
    }
    if (bdf.length > 0) {
      v.bdf = bdf
    }
    if (sirene_ul) {
      v.sirene_ul = sirene_ul
    }
    if (crp) {
      v.crp = crp
    }
    if (Object.keys(v) !== []) {
      emit("entreprise_" + this.value.key, v)
    }
  }
}
