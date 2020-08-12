import { iterable } from "./iterable"
import { debits } from "./debits"
import { apconso } from "./apconso"
import { apdemande } from "./apdemande"
import { flatten } from "./flatten"
import { compte } from "./compte"
import { effectifs, SortieEffectif } from "./effectifs"
import { delai } from "./delai"
import { bdf } from "./bdf"
import { diane } from "./diane"
import { sirene } from "./sirene"
import { cotisations } from "./cotisations"
import { dateAddDay } from "./dateAddDay"
import { dealWithProcols } from "./dealWithProcols"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { omit } from "../common/omit"

export type SortieMap = {
  effectif: SortieEffectif[]
  procol: Record<DataHash, EntréeDefaillances> // ou SortieProcols ?
  crp: unknown // TODO: à définir
} & Record<string, unknown> // TODO: à expliciter, cf reduce.algo2/map.ts

// Paramètres globaux utilisés par "public"
declare let actual_batch: BatchKey

declare function emit(key: string, value: Partial<SortieMap>): void

export function map(this: { value: CompanyDataValues }): void {
  /* DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO */ const f = {
    ...{ iterable, debits, apconso, apdemande, flatten, compte, effectifs }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ delai, dealWithProcols, sirene, cotisations, dateAddDay, omit }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ generatePeriodSerie, diane, bdf }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  const value = f.flatten(this.value, actual_batch)

  if (this.value.scope === "etablissement") {
    const vcmde: Partial<SortieMap> = {}
    vcmde.key = this.value.key
    vcmde.batch = actual_batch
    vcmde.effectif = f.effectifs(value.effectif)
    vcmde.dernier_effectif = vcmde.effectif[vcmde.effectif.length - 1]
    vcmde.sirene = f.sirene(f.iterable(value.sirene))
    vcmde.cotisation = f.cotisations(value.cotisation)
    vcmde.debit = f.debits(value.debit)
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
    const v: Partial<SortieMap> = {}
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
