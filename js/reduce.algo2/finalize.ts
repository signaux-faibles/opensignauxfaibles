import { f } from "./functions"
import { SiretOrSiren, Siret, BatchKey } from "../RawDataTypes"
import {
  CléSortieMap,
  SortieMap,
  SortieMapEntreprise,
  SortieMapEtablissement,
} from "./map"

type Accumulateurs = {
  /** Cumul du nombre de personnes employées par tous les établissements de l'entreprise. */
  effectif_entreprise?: number
  /** Cumul du nombre d'heures d'activité partielle consommées par tous les établissements de l'entreprise. */
  apart_entreprise?: number
  /** Cumul du montant des débits de tous les établissements de l'entreprise. */
  debit_entreprise?: number
  /** Nombre d'établissements rattachés à l'entreprise. */
  nbr_etablissements_connus: number
}

type DonnéesEntreprise = SortieMapEntreprise & Accumulateurs

type SortieEtabAvecEntreprise = SortieMapEtablissement & DonnéesEntreprise

export type Clé = {
  batch: BatchKey
  siren: SiretOrSiren
  periode: Date
  type: CléSortieMap["type"]
}

export type SortieFinalize = SortieEtabAvecEntreprise[] | { incomplete: true }

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "finalize"
  computed: Accumulateurs
  transmitted: unknown // Note: les autres champs ont été documentés dans les autres fichiers constituants les types inclus dans SortieEtabAvecEntreprise
}

declare function print(str: string): void

const bsonsize = (obj: unknown): number => JSON.stringify(obj).length // will not be included in jsFunctions.go

export function finalize(k: Clé, v: SortieMap): SortieFinalize {
  "use strict"

  const maxBsonSize = 16777216

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
  const établissements: Record<Siret, SortieMapEtablissement> = f.omit(
    v,
    "entreprise"
  )
  const entr: DonnéesEntreprise = { ...v.entreprise } as DonnéesEntreprise // on suppose que v.entreprise est défini

  const output: SortieEtabAvecEntreprise[] = Object.keys(établissements)
    .map((siret) => {
      const etab: SortieMapEtablissement = établissements[siret] ?? {}
      if (etab.effectif) {
        entr.effectif_entreprise =
          (entr.effectif_entreprise || 0) + etab.effectif
      }
      if (etab.apart_heures_consommees) {
        entr.apart_entreprise =
          (entr.apart_entreprise || 0) + etab.apart_heures_consommees
      }
      if (etab.montant_part_patronale || etab.montant_part_ouvriere) {
        entr.debit_entreprise =
          (entr.debit_entreprise ?? 0) +
          (etab.montant_part_patronale ?? 0) +
          (etab.montant_part_ouvriere ?? 0)
      }
      return etab
    })
    .map((etab) => ({
      ...etab, // TODO: s'assurer que certains champs de données d'établissement ne sont pas écrasés par des données d'entreprise portant le même nom
      ...entr,
      nbr_etablissements_connus: Object.keys(établissements).length,
    }))

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
