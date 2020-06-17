type Entreprise = {
  effectif: number
  effectif_entreprise: number
  apart_heures_consommees: number
  apart_entreprise: number
  montant_part_patronale: number
  montant_part_ouvriere: number
  debit_entreprise: number
  nbr_etablissements_connus: number
}

type Clé = {
  batch: unknown
  siren: SiretOrSiren
  periode: unknown
  type: unknown
}

type V = Record<SiretOrSiren, Entreprise> & {
  // _id: Clé
  // value: Record<SiretOrSiren, Etablissement>
  entreprise: Entreprise
}

type Output = unknown[] | { incomplete: true } | undefined

declare function print(str: string): void

export function finalize(
  key: Clé,
  v: V // { _[key: string]: T }
): Output {
  "use strict"
  const bsonsize =
    // @ts-expect-error: Object.bsonsize is not known by TypeScript, but it exists when the function in run by MongoDB
    Object.bsonsize || ((obj: unknown): number => JSON.stringify(obj).length) // DO // _NOT_INCLUDE_IN_JSFUNCTIONS_GO
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

  const etablissements_connus: Record<SiretOrSiren, boolean> = {}
  const entreprise: Entreprise = v.entreprise || {}

  Object.keys(v).forEach((siret) => {
    if (siret != "entreprise") {
      etablissements_connus[siret] = true
      if (v[siret].effectif) {
        entreprise.effectif_entreprise =
          (entreprise.effectif_entreprise || 0) + v[siret].effectif // initialized to null
      }
      if (v[siret].apart_heures_consommees) {
        entreprise.apart_entreprise =
          (entreprise.apart_entreprise || 0) + v[siret].apart_heures_consommees // initialized to 0
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
    if (siret != "entreprise") {
      Object.assign(v[siret], entreprise)
    }
  })

  // une fois que les comptes sont faits...
  const output: Entreprise[] = []
  const nb_connus = Object.keys(etablissements_connus).length
  Object.keys(v).forEach((siret) => {
    if (siret != "entreprise" && v[siret]) {
      v[siret].nbr_etablissements_connus = nb_connus
      output.push(v[siret])
    }
  })

  // NON: Pour l'instant, filtrage a posteriori
  // output = output.filter(siret_data => {
  //   return(siret_data.effectif) // Only keep if there is known effectif
  // })

  if (output.length > 0 && nb_connus <= 1500) {
    if (bsonsize(output) + bsonsize({ _id: key }) < maxBsonSize) {
      return output
    } else {
      print(
        "Warning: my name is " +
          JSON.stringify(key, null, 2) +
          " and I died in reduce.algo2/finalize.js"
      )
      return { incomplete: true }
    }
  }
}
