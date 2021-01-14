import { EntréeDelai, ParHash } from "../RawDataTypes"

export function delai(delai?: ParHash<EntréeDelai>): EntréeDelai[] {
  return Object.values(delai ?? {})
}
