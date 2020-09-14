import { f } from "./functions"
import "./js_params"
import {
  CompanyDataValues,
  BatchKey,
  Siret,
  SiretOrSiren,
  ParPériode,
} from "../RawDataTypes"
import { SortieBdf } from "./entr_bdf"
import { SortieDiane } from "./entr_diane"
import { SortieSireneEntreprise } from "./entr_sirene"
import { DonnéesAgrégées } from "./outputs"

type SortieMapEntreprise = {
  periode: Date
} & Partial<SortieSireneEntreprise> &
  Partial<SortieBdf> &
  Partial<SortieDiane>

export type SortieMapEtablissement = Partial<DonnéesAgrégées>

type SortieMapEtablissements = Record<Siret, SortieMapEtablissement>

export type SortieMap = {
  entreprise?: SortieMapEntreprise
} & SortieMapEtablissements

export type CléSortieMap = {
  batch: BatchKey
  siren: SiretOrSiren
  periode: Date
  type: "apart" | "other"
}

export type EntréeMap = {
  _id: SiretOrSiren
  value: CompanyDataValues
}

declare function emit(key: CléSortieMap, value: SortieMap): void

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
    const periodes = Object.keys(output_indexed)
      .sort()
      .map((timestamp) => parseInt(timestamp))

    if (includes["apart"] || includes["all"]) {
      if (v.apconso && v.apdemande) {
        const output_apart = f.apart(v.apconso, v.apdemande)
        Object.keys(output_apart).forEach((periode) => {
          const data: SortieMapEtablissements = {
            [this._id]: {
              ...output_apart[periode],
              siret: this._id,
            },
          }
          emit(
            {
              batch: actual_batch,
              siren: this._id.substring(0, 9),
              periode: new Date(Number(periode)),
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

      if (v.interim) {
        const output_interim = f.interim(v.interim, output_indexed)
        f.add(output_interim, output_indexed)
      }

      if (v.reporder) {
        const output_repeatable = f.repeatable(v.reporder)
        f.add(output_repeatable, output_indexed)
      }

      let output_cotisationsdettes
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
        const output_delai = f.delais(v.delai, output_cotisationsdettes || {}, {
          premièreDate: serie_periode[0],
          dernièreDate: serie_periode[serie_periode.length - 1],
        })
        f.add(output_delai, output_indexed)
      }

      v.altares = v.altares || {}
      v.procol = v.procol || {}

      if (v.altares) {
        f.defaillances(v.altares, v.procol, output_indexed)
      }

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
      const output_indexed: ParPériode<SortieMapEntreprise> = {}

      for (const periode of serie_periode) {
        output_indexed[periode.getTime()] = {
          siren: v.key,
          periode,
          exercice_bdf: 0,
        }
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

      v.bdf = v.bdf || {}
      v.diane = v.diane || {}

      if (v.bdf) {
        const outputBdf = f.entr_bdf(v.bdf, periodes)
        f.add(outputBdf, output_indexed)
      }

      if (v.diane) {
        /*const outputDiane =*/ f.entr_diane(v.diane, output_indexed, periodes)
        // f.add(outputDiane, output_indexed)
        // TODO: rendre f.entr_diane() pure, c.a.d. faire en sorte qu'elle ne modifie plus output_indexed directement
      }

      serie_periode.forEach((date) => {
        const periode = output_indexed[date.getTime()]
        if (
          typeof periode.arrete_bilan_bdf !== "undefined" ||
          typeof periode.arrete_bilan_diane !== "undefined"
        ) {
          emit(
            {
              batch: actual_batch,
              siren: this._id.substring(0, 9),
              periode: periode.periode,
              type: "other",
            },
            {
              entreprise: periode,
            }
          )
        }
      })
    }
  }
}
