import { f } from "./functions"
import {
  EntréeApConso,
  EntréeApDemande,
  EntréeCompte,
  EntréeDelai,
  EntréeDéfaillances,
  EntréePaydex,
} from "../GeneratedTypes"
import {
  CompanyDataValues,
  BatchKey,
  EntréeDiane,
  EntréeEllisphere,
  EntréeSirene,
  EntréeSireneEntreprise,
} from "../RawDataTypes"
import { SortieDebit } from "./debits"
import { Bdf } from "./bdf"

type SortieMapCommon = {
  key: string
  batch: string
}

type SortieMapEtablissement = SortieMapCommon & {
  sirene: Partial<EntréeSirene>
  debit: SortieDebit[]
  apconso: EntréeApConso[]
  apdemande: EntréeApDemande[]
  delai: EntréeDelai[]
  compte?: EntréeCompte
  procol: EntréeDéfaillances[]
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
  sirene_ul: EntréeSireneEntreprise
  ellisphere: EntréeEllisphere
  paydex: EntréePaydex[]
}

export type SortieMap = SortieMapEtablissement | SortieMapEntreprise

// Paramètres globaux utilisés par "public"
declare const actual_batch: BatchKey
declare const serie_periode: Date[]

// Types de données en entrée et sortie
export type Input = { _id: unknown; value: CompanyDataValues }
export type OutKey = string
export type OutValue = Partial<SortieMap>
declare function emit(key: string, value: OutValue): void

export function map(this: Input): void {
  const value = f.flatten(this.value, actual_batch)

  if (this.value.scope === "etablissement") {
    const vcmde: Partial<SortieMapEtablissement> = {}
    vcmde.key = this.value.key
    vcmde.batch = actual_batch
    vcmde.sirene = f.sirene(Object.values(value.sirene ?? {}))
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
    vcmde.procol = Object.values(value.procol ?? {})

    emit("etablissement_" + this.value.key, vcmde)
  } else if (this.value.scope === "entreprise") {
    const v: Partial<SortieMapEntreprise> = {}
    const diane = f.diane(value.diane)
    const bdf = f.bdf(value.bdf)
    const sirene_ul = Object.values(value.sirene_ul ?? {})[0] ?? null
    const ellisphere = Object.values(value.ellisphere ?? {})[0] ?? null

    if (sirene_ul) {
      sirene_ul.raison_sociale = f.raison_sociale(
        sirene_ul.raison_sociale,
        sirene_ul.nom_unite_legale,
        sirene_ul.nom_usage_unite_legale,
        sirene_ul.prenom1_unite_legale,
        sirene_ul.prenom2_unite_legale,
        sirene_ul.prenom3_unite_legale,
        sirene_ul.prenom4_unite_legale
      )
    }
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
    if (ellisphere) {
      v.ellisphere = ellisphere
    }

    if (value.paydex) {
      v.paydex = Object.values(value.paydex).sort(
        (p1, p2) => p1.date_valeur.getTime() - p2.date_valeur.getTime()
      )
    }

    if (Object.keys(v) !== []) {
      emit("entreprise_" + this.value.key, v)
    }
  }
}
