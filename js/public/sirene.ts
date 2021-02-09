import { EntréeSirene } from "../GeneratedTypes"

// Cette fonction retourne les données sirene les plus récentes
export function sirene(sireneArray: EntréeSirene[]): Partial<EntréeSirene> {
  return sireneArray[sireneArray.length - 1] || {} // TODO: vérifier que sireneArray est bien classé dans l'ordre chronologique -> c'est sûr qu'il ne l'est pas, vérifier que pour toute la base on a bien un objet sirene unique !
}
