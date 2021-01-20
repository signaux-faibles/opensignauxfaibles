import { EntréeCompte, ParPériode, Periode, ParHash } from "../RawDataTypes"

export type SortieCompte = {
  /** Compte administratif URSSAF */
  compte_urssaf: string
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "compte"
  computed: SortieCompte
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

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
