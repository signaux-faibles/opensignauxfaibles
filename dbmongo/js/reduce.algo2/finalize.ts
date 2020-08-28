import { Siret, SortieMap, SortieMapEtablissement } from "./map"
import * as f from "../common/omit"

type EntrepriseEnSortie = {
  effectif_entreprise: number
  apart_heures_consommees: number
  apart_entreprise: number
  montant_part_patronale: number
  montant_part_ouvriere: number
  debit_entreprise: number
  nbr_etablissements_connus: number
  random_order?: number
  siret: SiretOrSiren
  periode: unknown
}

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
  const entreprise: Partial<EntrepriseEnSortie> = v.entreprise || {}
  const etab: Record<
    Siret,
    SortieMapEtablissement & Partial<EntrepriseEnSortie>
  > = f.omit(v, "entreprise")

  Object.keys(etab).forEach((siret) => {
    const { effectif } = etab[siret]
    if (effectif) {
      entreprise.effectif_entreprise =
        (entreprise.effectif_entreprise || 0) + effectif // initialized to null
    }
    const { apart_heures_consommees } = etab[siret]
    if (apart_heures_consommees) {
      entreprise.apart_entreprise =
        (entreprise.apart_entreprise || 0) + apart_heures_consommees // initialized to 0
    }
    if (
      etab[siret].montant_part_patronale ||
      etab[siret].montant_part_ouvriere
    ) {
      entreprise.debit_entreprise =
        (entreprise.debit_entreprise || 0) +
        (etab[siret].montant_part_patronale || 0) +
        (etab[siret].montant_part_ouvriere || 0)
    }

    Object.assign(etab[siret], entreprise)
  })

  // une fois que les comptes sont faits...
  const output: Partial<EntrepriseEnSortie>[] = []
  const nb_connus = Object.keys(etab).length
  Object.keys(etab).forEach((siret) => {
    etab[siret].nbr_etablissements_connus = nb_connus
    output.push(etab[siret])
  })

  // NON: Pour l'instant, filtrage a posteriori
  // output = output.filter(siret_data => {
  //   return(siret_data.effectif) // Only keep if there is known effectif
  // })

  if (output.length > 0 && nb_connus <= 1500) {
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
