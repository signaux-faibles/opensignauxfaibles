import { Siret, SortieMap, SortieMapEtablissement } from "./map"
import * as f from "../common/omit"

type Accumulateurs = {
  effectif_entreprise: number
  apart_entreprise: number
  debit_entreprise: number
  nbr_etablissements_connus: number
}

export type EntrepriseEnSortie = SortieMapEtablissement & Accumulateurs

export type Clé = {
  batch: unknown
  siren: SiretOrSiren
  periode: unknown
  type: unknown
}

type SortieFinalize = Partial<EntrepriseEnSortie>[] | { incomplete: true }

declare function print(str: string): void

export function finalize(k: Clé, v: SortieMap): SortieFinalize {
  "use strict"

  const maxBsonSize = 16777216
  const bsonsize = (obj: unknown): number => JSON.stringify(obj).length // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  // v de la forme
  // _id: {batch / siren / periode / type}
  // value: {siret1: {}, siret2: {}, "siren": {}}
  //
  ///
  ///////////////////////////////////////////////
  // consolidation a l'echelle de l'entreprise //
  ///////////////////////////////////////////////
  ///
  //

  // extraction de l'entreprise et des établissements depuis v
  const etab: Record<Siret, SortieMapEtablissement> = f.omit(v, "entreprise")
  const entr: Partial<EntrepriseEnSortie> = { ...v.entreprise }

  const output: Partial<EntrepriseEnSortie>[] = Object.keys(etab).map(
    (siret) => {
      const { effectif } = etab[siret]
      if (effectif) {
        entr.effectif_entreprise = entr.effectif_entreprise || 0 + effectif
      }
      const { apart_heures_consommees } = etab[siret]
      if (apart_heures_consommees) {
        entr.apart_entreprise =
          (entr.apart_entreprise || 0) + apart_heures_consommees
      }
      if (
        etab[siret].montant_part_patronale ||
        etab[siret].montant_part_ouvriere
      ) {
        entr.debit_entreprise =
          (entr.debit_entreprise || 0) +
          (etab[siret].montant_part_patronale || 0) +
          (etab[siret].montant_part_ouvriere || 0)
      }

      return {
        ...etab[siret],
        ...entr,
        nbr_etablissements_connus: Object.keys(etab).length,
      }
    }
  )

  // NON: Pour l'instant, filtrage a posteriori
  // output = output.filter(siret_data => {
  //   return(siret_data.effectif) // Only keep if there is known effectif
  // })

  if (output.length > 0 && output.length <= 1500) {
    if (bsonsize(output) + bsonsize({ _id: k }) < maxBsonSize) {
      return output
    } else {
      print(
        "Warning: my name is " +
          JSON.stringify(k, null, 2) +
          " and I died in reduce.algo2/finalize.js"
      )
      return { incomplete: true }
    }
  } else {
    return [] // ajouté pour résoudre erreur TS7030 (Not all code paths return a value)
  }
}
