import { EntréeCompte, ParPériode, Periode, ParHash } from "../RawDataTypes"

export type SortieCompte = { compte_urssaf: number }

export function compte(
  compte: ParHash<EntréeCompte>
): ParPériode<SortieCompte> {
  "use strict"
  const output_compte: ParPériode<SortieCompte> = {}

  //  var offset_compte = 3
  for (const compteEntry of Object.values(compte)) {
    const periode: Periode = compteEntry.periode.getTime().toString()
    output_compte[periode] = {
      ...(output_compte[periode] ?? {}),
      compte_urssaf: compteEntry.numero_compte,
    }
  }

  return output_compte
}
