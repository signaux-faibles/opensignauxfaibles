import { EntréeCompte } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"

export function compte(
  compte?: ParHash<EntréeCompte>
): EntréeCompte | undefined {
  const c = Object.values(compte ?? {})
  return c.length > 0 ? c[c.length - 1] : undefined
}
