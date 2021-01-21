import { CodeAPE, ParPériode } from "../RawDataTypes"

type CodeAPENiveau2 = string
type CodeAPENiveau3 = string
type CodeAPENiveau4 = string

type CodeNAF = string

export type NAF = {
  n2to1: Record<CodeAPENiveau2, CodeNAF>
  n3to1: Record<CodeAPENiveau3, CodeNAF>
  n4to1: Record<CodeAPENiveau4, CodeNAF>
  n5to1: Record<CodeAPE, CodeNAF>
  n1: Record<CodeNAF, string>
  n2: Record<CodeAPENiveau2, string>
  n3: Record<CodeAPENiveau3, string>
  n4: Record<CodeAPENiveau4, string>
  n5: Record<CodeAPE, string>
}

type Input = {
  code_ape: CodeAPE
}

export type SortieNAF = {
  /** Code APE (code d'activité principale), aussi appelé code NAF (nomenclature d’activité française). */
  code_naf: CodeNAF
  /** Libellé de code NAF/APE. */
  libelle_naf: string
  /** Deuxième niveau du code NAF/APE: Section et Division. */
  code_ape_niveau2: CodeAPENiveau2
  /** Troisième niveau du code NAF/APE: Section, Division et Groupe. */
  code_ape_niveau3: CodeAPENiveau3
  /** Quatrième niveau du code NAF/APE: Section, Division, Groupe et Classe. (sans la Sous-classe) */
  code_ape_niveau4: CodeAPENiveau4
  /** Libellé de code NAF/APE de deuxième niveau: Division. */
  libelle_ape2: string
  /** Libellé de code NAF/APE de troisième niveau: Groupe. */
  libelle_ape3: string
  /** Libellé de code NAF/APE de quatrième niveau: Classe. */
  libelle_ape4: string
  /** Libellé de code NAF/APE de cinquième niveau: Sous-classe. */
  libelle_ape5: string
}

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "populateNafAndApe"
  computed: SortieNAF
  transmitted: unknown // unknown ~= aucune variable n'est transmise directement depuis RawData
}

export function populateNafAndApe(
  output_indexed: ParPériode<Partial<Input> & Partial<SortieNAF>>,
  naf: NAF
): void {
  "use strict"
  for (const outputForKey of Object.values(output_indexed)) {
    const code_ape = outputForKey.code_ape
    if (code_ape) {
      const code_naf = naf.n5to1[code_ape]
      outputForKey.code_naf = code_naf
      outputForKey.libelle_naf = code_naf ? naf.n1[code_naf] : undefined
      const code_ape_niveau2 = code_ape.substring(0, 2)
      outputForKey.code_ape_niveau2 = code_ape_niveau2
      const code_ape_niveau3 = code_ape.substring(0, 3)
      outputForKey.code_ape_niveau3 = code_ape_niveau3
      const code_ape_niveau4 = code_ape.substring(0, 4)
      outputForKey.code_ape_niveau4 = code_ape_niveau4
      outputForKey.libelle_ape2 = naf.n2[code_ape_niveau2]
      outputForKey.libelle_ape3 = naf.n3[code_ape_niveau3]
      outputForKey.libelle_ape4 = naf.n4[code_ape_niveau4]
      outputForKey.libelle_ape5 = naf.n5[code_ape]
    }
  }
}
