/**
 * This file was automatically generated by generate-types.sh.
 *
 * DO NOT MODIFY IT BY HAND.
 *
 * Instead:
 * - modify the validation/*.schema.json files;
 * - then, run generate-types.sh to regenerate this file.
 */

export interface EntreeBdf {
  siren: string
  [k: string]: unknown
}
export interface EntreeDelai {
  date_creation: Date
  date_echeance: Date
  /**
   * doit valoir 1 ou plus
   */
  duree_delai: number
  /**
   * doit valoir plus que 0 euros
   */
  montant_echeancier: number
  [k: string]: unknown
}
