import { EntréeCompte, ParPériode, Periode, ParHash } from "../RawDataTypes"

type SortieCompte = ParPériode<{ compte_urssaf: unknown }> // TODO: choisir un type plus précis

export function compte(compte: ParHash<EntréeCompte>): SortieCompte {
  "use strict"
  const output_compte: SortieCompte = {}

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
