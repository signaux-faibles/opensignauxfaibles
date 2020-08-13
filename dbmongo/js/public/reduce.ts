import { SortieMap } from "./map"

type V = Partial<SortieMap> & { sirets?: unknown[] }

export function reduce(_key: { scope: Scope }, values: V[]): V {
  return values.reduce((m, v) => {
    if (v.sirets) {
      // TODO: je n'ai pas trouvé d'affectation de valeur dans la propriété "sirets" => est-elle toujours d'actualité ?
      m.sirets = (m.sirets || []).concat(v.sirets)
      delete v.sirets
    }
    Object.assign(m, v)
    return m
  }, {} as V)
}
