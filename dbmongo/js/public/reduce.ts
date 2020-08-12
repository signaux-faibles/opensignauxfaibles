import { PublicMapResult } from "./map"

type V = PublicMapResult & { sirets?: unknown[] }

export function reduce(_key: { scope: Scope }, values: V[]): V {
  // if (key.scope = "entreprise") { // TODO: cette expression est toujours vraie var elle affecte "entreprise" à key.scope => à corriger
  return values.reduce((m, v) => {
    if (v.sirets) {
      // TODO: je n'ai pas trouvé d'affectation de valeur dans la propriété "sirets" => est-elle toujours d'actualité ?
      m.sirets = (m.sirets || []).concat(v.sirets)
      delete v.sirets
    }
    Object.assign(m, v)
    return m
  }, {} as V)
  // }
  // return (values as unknown) as V // TODO: veut-on vraiment retourner values tel quel ? (tableau de type V[] au lieu d'objet de type V)
}
