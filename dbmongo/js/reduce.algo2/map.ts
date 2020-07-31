import { flatten } from "./flatten"
import { outputs, DonnéesAgrégées } from "./outputs"
import { apart } from "./apart"
import { compte } from "./compte"
import { effectifs } from "./effectifs"
import { interim } from "./interim"
import { add } from "./add"
import { repeatable } from "./repeatable"
import { delais } from "./delais"
import { defaillances } from "./defaillances"
import { cotisationsdettes } from "./cotisationsdettes"
import { ccsf } from "./ccsf"
import { sirene } from "./sirene"
import { populateNafAndApe } from "./populateNafAndApe"
import { cotisation } from "./cotisation"
import { cibleApprentissage } from "./cibleApprentissage"
import { entr_sirene, SortieSireneEntreprise } from "./entr_sirene"
import { dateAddMonth } from "./dateAddMonth"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { poidsFrng } from "./poidsFrng"
import { detteFiscale } from "./detteFiscale"
import { fraisFinancier } from "./fraisFinancier"
import { entr_bdf, SortieBdf } from "./entr_bdf"
import { omit } from "../common/omit"
import { entr_diane, SortieDiane } from "./entr_diane"

type Siret = string

type SortieMapEntreprise = {
  periode: Date
} & Partial<SortieSireneEntreprise> &
  Partial<EntréeBdf> & // TODO: est-ce nécéssaire d'inclure les types d'entrée ?
  Partial<EntréeDiane> &
  Partial<EntréeBdf> &
  Partial<SortieBdf> &
  Partial<SortieDiane>

type SortieMapEtablissement = Partial<DonnéesAgrégées>

type SortieMap =
  | { entreprise: SortieMapEntreprise }
  | Record<Siret, SortieMapEtablissement>

type CléSortieMap = {
  batch: BatchKey
  siren: SiretOrSiren
  periode: Date
  type: "apart" | "other"
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
export function map(this: {
  _id: SiretOrSiren
  value: CompanyDataValues
}): void {
  "use strict"
  /* DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO */ const f = {
    ...{ flatten, outputs, apart, compte, effectifs, interim, add }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ repeatable, delais, defaillances, cotisationsdettes, ccsf }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ sirene, populateNafAndApe, cotisation, cibleApprentissage }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ entr_sirene, dateAddMonth, generatePeriodSerie, poidsFrng }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ detteFiscale, fraisFinancier, entr_bdf, omit, entr_diane }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

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
          const data: SortieMap = {
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
        const output_compte = f.compte(v as DonnéesCompte)
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
          v as DonnéesCotisation & DonnéesDebit,
          periodes,
          date_fin
        )
        f.add(output_cotisationsdettes, output_indexed)
      }

      if (v.delai) {
        const output_delai = f.delais(
          v as DonnéesDelai,
          output_cotisationsdettes || {},
          {
            premièreDate: serie_periode[0],
            dernièreDate: serie_periode[serie_periode.length - 1],
          }
        )
        f.add(output_delai, output_indexed)
      }

      v.altares = v.altares || {}
      v.procol = v.procol || {}

      if (v.altares) {
        f.defaillances(v as DonnéesDefaillances, output_indexed)
      }

      if (v.ccsf) {
        f.ccsf(v as DonnéesCcsf, output_array)
      }
      if (v.sirene) {
        f.sirene(v as DonnéesSirene, output_array)
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
          arrete_bilan_bdf: new Date(0),
          exercice_diane: 0,
          arrete_bilan_diane: new Date(0),
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
        const outputDiane = f.entr_diane(v.diane, output_indexed, periodes)
        f.add(outputDiane, output_indexed)
      }

      serie_periode.forEach((date) => {
        const periode = output_indexed[date.getTime()]
        if (
          (periode.arrete_bilan_bdf || new Date(0)).getTime() === 0 &&
          (periode.arrete_bilan_diane || new Date(0)).getTime() === 0
        ) {
          delete output_indexed[date.getTime()]
        }
        if ((periode.arrete_bilan_bdf || new Date(0)).getTime() === 0) {
          delete periode.arrete_bilan_bdf
        }
        if ((periode.arrete_bilan_diane || new Date(0)).getTime() === 0) {
          delete periode.arrete_bilan_diane
        }

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
      })
    }
  }
}
