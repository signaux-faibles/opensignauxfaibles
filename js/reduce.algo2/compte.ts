import { f } from "./functions"
import { EntréeCompte } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"
import { ParPériode } from "../common/makePeriodeMap"

export type SortieCompte = {
  /** Compte administratif URSSAF */
  compte_urssaf: string
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "compte"
  computed: unknown // unknown ~= aucune variable n'est calculée
  transmitted: SortieCompte
}

export function compte(
  compte: ParHash<EntréeCompte>
): ParPériode<SortieCompte> {
  "use strict"
  const output_compte = f.makePeriodeMap<SortieCompte>()

  //  var offset_compte = 3
  for (const { periode, numero_compte } of Object.values(compte)) {
    output_compte.assign(periode, { compte_urssaf: numero_compte })
  }

  return output_compte
}
