// Cette fonction retourne les données sirene les plus récentes
export function sirene(sireneArray: unknown[]): unknown {
  return sireneArray[sireneArray.length - 1] || {} // TODO: vérifier que sireneArray est bien classé dans l'ordre chronologique
}
