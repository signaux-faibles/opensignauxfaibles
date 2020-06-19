import { flatten, V as FlattenedEntreprise } from "./flatten"
import { outputs } from "./outputs"
import { apart } from "./apart"
import { compte, V as CompteInput } from "./compte"
import { effectifs } from "./effectifs"
import { interim } from "./interim"
import { add } from "./add"
import { repeatable } from "./repeatable"
import { delais, V as DelaisV } from "./delais"
import { defaillances } from "./defaillances"
import { cotisationsdettes } from "./cotisationsdettes"
import { ccsf } from "./ccsf"
import { sirene } from "./sirene"
import { populateNafAndApe, NAF } from "./populateNafAndApe"
import { cotisation } from "./cotisation"
import { cibleApprentissage } from "./cibleApprentissage"

declare const naf: NAF
declare const actual_batch: BatchKey
declare const includes: Record<"all" | "apart", boolean>

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
  value: FlattenedEntreprise
}): void {
  "use strict"
  /* DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO */ const f = {
    ...{ flatten, outputs, apart, compte, effectifs, interim, add }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ repeatable, delais, defaillances, cotisationsdettes, ccsf }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
    ...{ sirene, populateNafAndApe, cotisation, cibleApprentissage }, // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  const v = f.flatten(this.value, actual_batch)

  if (v.scope === "etablissement") {
    const [
      output_array, // [ OutputValue ], in chronological order
      output_indexed, // { Periode -> OutputValue }
    ] = f.outputs(v, serie_periode)

    // Les periodes qui nous interessent, triées
    const periodes = Object.keys(output_indexed).sort((a, b) =>
      a >= b ? 1 : 0
    )

    if (includes["apart"] || includes["all"]) {
      if (v.apconso && v.apdemande) {
        const output_apart = f.apart(v.apconso, v.apdemande)
        Object.keys(output_apart).forEach((periode) => {
          const data: Record<SiretOrSiren, { siret?: SiretOrSiren }> = {}
          data[this._id] = {
            ...output_apart[periode],
            siret: this._id,
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
        const output_compte = f.compte(v as CompteInput)
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

      if (v.delai) {
        const output_delai = f.delais(v as DelaisV, output_indexed)
        f.add(output_delai, output_indexed)
      }

      v.altares = v.altares || {}
      v.procol = v.procol || {}

      if (v.altares) {
        f.defaillances(v as DonnéesDefaillances, output_indexed)
      }

      if (v.cotisation && v.debit) {
        const output_cotisationsdettes = f.cotisationsdettes(
          v as DonnéesCotisationsDettes,
          periodes
        )
        f.add(output_cotisationsdettes, output_indexed)
      }

      if (v.ccsf) {
        f.ccsf(v as DonnéesCcsf, output_array)
      }
      if (v.sirene) {
        f.sirene(v as DonnéesSirene, output_array)
      }

      f.populateNafAndApe(output_indexed, naf)

      f.cotisation(output_indexed, output_array)

      const output_cible = f.cibleApprentissage(output_indexed, 18)
      f.add(output_cible, output_indexed)
      output_array.forEach((val) => {
        const data: Record<SiretOrSiren, typeof val> = {}
        data[this._id] = val
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
      const output_array = serie_periode.map(function (e) {
        return {
          siren: v.key,
          periode: e,
          exercice_bdf: 0,
          arrete_bilan_bdf: new Date(0),
          exercice_diane: 0,
          arrete_bilan_diane: new Date(0),
        }
      })

      var output_indexed = output_array.reduce(function (periode, val) {
        periode[val.periode.getTime()] = val
        return periode
      }, {})

      if (v.sirene_ul) {
        f.sirene_ul(v, output_array)
      }

      const periodes = Object.keys(output_indexed).sort((a, b) => a >= b)
      if (v.effectif_ent) {
        const output_effectif_ent = f.effectifs(
          v.effectif_ent,
          periodes,
          "effectif_ent"
        )
        f.add(output_effectif_ent, output_indexed)
      }

      var output_indexed = output_array.reduce(function (periode, val) {
        periode[val.periode.getTime()] = val
        return periode
      }, {})

      v.bdf = v.bdf || {}
      v.diane = v.diane || {}

      Object.keys(v.bdf).forEach((hash) => {}, {})

      v.bdf = v.bdf || {}
      v.diane = v.diane || {}

      Object.keys(v.bdf).forEach((hash) => {
        const periode_arrete_bilan = new Date(
          Date.UTC(
            v.bdf[hash].arrete_bilan_bdf.getUTCFullYear(),
            v.bdf[hash].arrete_bilan_bdf.getUTCMonth() + 1,
            1,
            0,
            0,
            0,
            0
          )
        )
        const periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7)
        const series = f.generatePeriodSerie(
          periode_dispo,
          f.dateAddMonth(periode_dispo, 13)
        )

        series.forEach((periode) => {
          Object.keys(v.bdf[hash])
            .filter((k) => {
              const omit = ["raison_sociale", "secteur", "siren"]
              return v.bdf[hash][k] != null && !omit.includes(k)
            })
            .forEach((k) => {
              if (periode.getTime() in output_indexed) {
                output_indexed[periode.getTime()][k] = v.bdf[hash][k]
                output_indexed[periode.getTime()].exercice_bdf =
                  output_indexed[periode.getTime()].annee_bdf - 1
              }

              const past_year_offset = [1, 2]
              past_year_offset.forEach((offset) => {
                const periode_offset = f.dateAddMonth(periode, 12 * offset)
                const variable_name = k + "_past_" + offset
                if (
                  periode_offset.getTime() in output_indexed &&
                  k != "arrete_bilan_bdf" &&
                  k != "exercice_bdf"
                ) {
                  output_indexed[periode_offset.getTime()][variable_name] =
                    v.bdf[hash][k]
                }
              })
            })
        })
      })

      Object.keys(v.diane)
        .filter((hash) => v.diane[hash].arrete_bilan_diane)
        .forEach((hash) => {
          //v.diane[hash].arrete_bilan_diane = new Date(Date.UTC(v.diane[hash].exercice_diane, 11, 31, 0, 0, 0, 0))
          const periode_arrete_bilan = new Date(
            Date.UTC(
              v.diane[hash].arrete_bilan_diane.getUTCFullYear(),
              v.diane[hash].arrete_bilan_diane.getUTCMonth() + 1,
              1,
              0,
              0,
              0,
              0
            )
          )
          const periode_dispo = f.dateAddMonth(periode_arrete_bilan, 7) // 01/08 pour un bilan le 31/12, donc algo qui tourne en 01/09
          const series = f.generatePeriodSerie(
            periode_dispo,
            f.dateAddMonth(periode_dispo, 14) // periode de validité d'un bilan auprès de la Banque de France: 21 mois (14+7)
          )

          series.forEach((periode) => {
            Object.keys(v.diane[hash])
              .filter((k) => {
                const omit = [
                  "marquee",
                  "nom_entreprise",
                  "numero_siren",
                  "statut_juridique",
                  "procedure_collective",
                ]
                return v.diane[hash][k] != null && !omit.includes(k)
              })
              .forEach((k) => {
                if (periode.getTime() in output_indexed) {
                  output_indexed[periode.getTime()][k] = v.diane[hash][k]
                }

                // Passé

                const past_year_offset = [1, 2]
                past_year_offset.forEach((offset) => {
                  const periode_offset = f.dateAddMonth(periode, 12 * offset)
                  const variable_name = k + "_past_" + offset

                  if (
                    periode_offset.getTime() in output_indexed &&
                    k != "arrete_bilan_diane" &&
                    k != "exercice_diane"
                  ) {
                    output_indexed[periode_offset.getTime()][variable_name] =
                      v.diane[hash][k]
                  }
                })
              })
          })

          series.forEach((periode) => {
            if (periode.getTime() in output_indexed) {
              // Recalcul BdF si ratios bdf sont absents
              if (
                !("poids_frng" in output_indexed[periode.getTime()]) &&
                f.poidsFrng(v.diane[hash]) !== null
              ) {
                output_indexed[periode.getTime()].poids_frng = f.poidsFrng(
                  v.diane[hash]
                )
              }
              if (
                !("dette_fiscale" in output_indexed[periode.getTime()]) &&
                f.detteFiscale(v.diane[hash]) !== null
              ) {
                output_indexed[
                  periode.getTime()
                ].dette_fiscale = f.detteFiscale(v.diane[hash])
              }
              if (
                !("frais_financier" in output_indexed[periode.getTime()]) &&
                f.fraisFinancier(v.diane[hash]) !== null
              ) {
                output_indexed[
                  periode.getTime()
                ].frais_financier = f.fraisFinancier(v.diane[hash])
              }

              const bdf_vars = [
                "taux_marge",
                "poids_frng",
                "dette_fiscale",
                "financier_court_terme",
                "frais_financier",
              ]
              const past_year_offset = [1, 2]
              bdf_vars.forEach((k) => {
                if (k in output_indexed[periode.getTime()]) {
                  past_year_offset.forEach((offset) => {
                    const periode_offset = f.dateAddMonth(periode, 12 * offset)
                    const variable_name = k + "_past_" + offset

                    if (periode_offset.getTime() in output_indexed) {
                      output_indexed[periode_offset.getTime()][variable_name] =
                        output_indexed[periode.getTime()][k]
                    }
                  })
                }
              })
            }
          })
        })

      output_array.forEach((periode, index) => {
        if (
          (periode.arrete_bilan_bdf || new Date(0)).getTime() == 0 &&
          (periode.arrete_bilan_diane || new Date(0)).getTime() == 0
        ) {
          delete output_array[index]
        }
        if ((periode.arrete_bilan_bdf || new Date(0)).getTime() == 0) {
          delete periode.arrete_bilan_bdf
        }
        if ((periode.arrete_bilan_diane || new Date(0)).getTime() == 0) {
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
