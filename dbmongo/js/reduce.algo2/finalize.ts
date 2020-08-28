import { SortieMap } from "./map"

type EntrepriseEnEntrée = {
  effectif: number | null
} & Partial<EntrepriseEnSortie>

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

type EntréeFinalize = Record<SiretOrSiren | "entreprise", EntrepriseEnEntrée>

type SortieFinalize =
  | Partial<EntrepriseEnSortie>[]
  | { incomplete: true }
  | undefined

declare function print(str: string): void

export function finalize(k: Clé, mapV: SortieMap): SortieFinalize {
  "use strict"
  const v = mapV as EntréeFinalize // TODO: améliorer l'alignement avec type SortieMap

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

  const etablissements_connus: Record<SiretOrSiren, boolean> = {}
  const entreprise: Partial<EntrepriseEnSortie> = v.entreprise || {}

  Object.keys(v).forEach((siret) => {
    if (siret !== "entreprise") {
      etablissements_connus[siret] = true
      const { effectif } = v[siret]
      if (effectif) {
        entreprise.effectif_entreprise =
          (entreprise.effectif_entreprise || 0) + effectif // initialized to null
      }
      const { apart_heures_consommees } = v[siret]
      if (apart_heures_consommees) {
        entreprise.apart_entreprise =
          (entreprise.apart_entreprise || 0) + apart_heures_consommees // initialized to 0
      }
      if (v[siret].montant_part_patronale || v[siret].montant_part_ouvriere) {
        entreprise.debit_entreprise =
          (entreprise.debit_entreprise || 0) +
          (v[siret].montant_part_patronale || 0) +
          (v[siret].montant_part_ouvriere || 0)
      }
    }
  })

  Object.keys(v).forEach((siret) => {
    if (siret !== "entreprise") {
      Object.assign(v[siret], entreprise)
    }
  })

  // une fois que les comptes sont faits...
  const output: Partial<EntrepriseEnSortie>[] = []
  const nb_connus = Object.keys(etablissements_connus).length
  Object.keys(v).forEach((siret) => {
    if (siret !== "entreprise" && v[siret]) {
      v[siret].nbr_etablissements_connus = nb_connus
      output.push(v[siret])
    }
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
