import { f } from "./functions"
import { CompanyDataValues, BatchKey, Siret, Siren } from "../RawDataTypes"
import { SortiePaydex } from "./entr_paydex"
import { SortieSireneEntreprise } from "./entr_sirene"
import { DonnéesAgrégées } from "./outputs"
import { NAF } from "./populateNafAndApe"
import { SortieEffectifs } from "./effectifs"
import { SortieCotisationsDettes } from "./cotisationsdettes"

export type SortieMapEntreprise = {
  siren: Siren
  periode: Date
} & Partial<SortieSireneEntreprise> &
  Partial<SortiePaydex> &
  Partial<SortieEffectifs<"effectif_ent">>

export type SortieMapEtablissement = Partial<DonnéesAgrégées>

type SortieMapEtablissements = Record<Siret, SortieMapEtablissement>

export type SortieMap = {
  entreprise?: SortieMapEntreprise
} & SortieMapEtablissements

export type CléSortieMap = {
  batch: BatchKey
  siren: Siren
  periode: Date
  type: "apart" | "other"
}

export type EntréeMap = {
  _id: Siret | Siren
  value: CompanyDataValues
}

declare function emit(key: CléSortieMap, value: SortieMap): void

// Paramètres globaux utilisés par "reduce.algo2"
declare const naf: NAF
declare const actual_batch: BatchKey
declare const includes: Record<"all" | "apart", boolean>
declare const serie_periode: Date[]
declare const date_fin: Date

/**
 * `map()` est appelée pour chaque entreprise/établissement.
 *
 * Une entreprise/établissement est rattachée à des données de plusieurs types,
 * groupées par *Batch* (groupements de fichiers de données importés).
 *
 * Pour chaque période d'un entreprise/établissement, un objet contenant toutes
 * les données agrégées est émis (par appel à `emit()`), à destination de
 * `reduce()`, puis de `finalize()`.
 */
export function map(this: EntréeMap): void {
  "use strict"

  const v = f.flatten(this.value, actual_batch)

  if (v.scope === "etablissement") {
    const [
      output_array, // DonnéesAgrégées[] dans l'ordre chronologique
      output_indexed, // { Periode -> DonnéesAgrégées }
    ] = f.outputs(v, serie_periode)

    // Les periodes qui nous interessent, triées
    const periodes = serie_periode.map((date) => date.getTime())

    if (includes["apart"] || includes["all"]) {
      if (v.apconso && v.apdemande) {
        const output_apart = f.apart(v.apconso, v.apdemande)
        output_apart.forEach((current, periode) => {
          if (!output_indexed.has(periode)) return // limiter dans le scope temporel du batch.
          const data: SortieMapEtablissements = {
            [this._id]: {
              ...current,
              siret: this._id,
            },
          }
          emit(
            {
              batch: actual_batch,
              siren: this._id.substring(0, 9),
              periode: new Date(periode),
              type: "apart",
            },
            data
          )
        })
      }
    }

    if (includes["all"]) {
      if (v.compte) {
        const output_compte = f.compte(v.compte)
        f.add(output_compte, output_indexed)
      }

      if (v.effectif) {
        const output_effectif = f.effectifs(v.effectif, periodes, "effectif")
        f.add(output_effectif, output_indexed)
      }

      if (v.reporder) {
        const output_repeatable = f.repeatable(v.reporder)
        f.add(output_repeatable, output_indexed)
      }

      let output_cotisationsdettes = f.makePeriodeMap<SortieCotisationsDettes>()
      if (v.cotisation && v.debit) {
        output_cotisationsdettes = f.cotisationsdettes(
          v.cotisation,
          v.debit,
          periodes,
          date_fin
        )
        f.add(output_cotisationsdettes, output_indexed)
      }

      if (v.delai) {
        const premièreDate = serie_periode[0]
        const dernièreDate = serie_periode[serie_periode.length - 1]
        if (premièreDate === undefined || dernièreDate === undefined) {
          const error = (message: string): never => {
            throw new Error(message)
          }
          error("serie_periode should not contain undefined values")
        } else {
          const output_delai = f.delais(v.delai, output_cotisationsdettes, {
            premièreDate,
            dernièreDate,
          })
          f.add(output_delai, output_indexed)
        }
      }

      v.procol = v.procol || {}

      f.defaillances(v.procol, output_indexed)

      if (v.ccsf) {
        f.ccsf(v.ccsf, output_array)
      }
      if (v.sirene) {
        f.sirene(v.sirene, output_array)
      }

      f.populateNafAndApe(output_indexed, naf)

      const output_cotisation = f.cotisation(output_indexed)
      f.add(output_cotisation, output_indexed)

      const output_cible = f.cibleApprentissage(output_indexed, 18)
      f.add(output_cible, output_indexed)

      output_array.forEach((val) => {
        const data: SortieMap = {
          [this._id]: val,
        }
        emit(
          {
            batch: actual_batch,
            siren: this._id.substring(0, 9),
            periode: val.periode,
            type: "other",
          },
          data
        )
      })
    }
  }

  if (v.scope === "entreprise") {
    if (includes["all"]) {
      const output_indexed = f.makePeriodeMap<SortieMapEntreprise>()
      for (const periode of serie_periode) {
        output_indexed.set(periode, {
          siren: v.key,
          periode,
        })
      }

      if (v.sirene_ul) {
        const outputEntrSirene = f.entr_sirene(v.sirene_ul, serie_periode)
        f.add(outputEntrSirene, output_indexed)
      }

      const periodes = serie_periode.map((date) => date.getTime())

      if (v.effectif_ent) {
        const output_effectif_ent = f.effectifs(
          v.effectif_ent,
          periodes,
          "effectif_ent"
        )
        f.add(output_effectif_ent, output_indexed)
      }

      if (v.paydex) {
        const paydexParPériode = f.entr_paydex(v.paydex, serie_periode)
        f.add(paydexParPériode, output_indexed)
      }
    }
  }
}
