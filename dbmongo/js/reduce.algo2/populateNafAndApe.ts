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
  code_naf: CodeNAF
  libelle_naf: string
  code_ape_niveau2: CodeAPENiveau2
  code_ape_niveau3: CodeAPENiveau3
  code_ape_niveau4: CodeAPENiveau4
  libelle_ape2: string
  libelle_ape3: string
  libelle_ape4: string
  libelle_ape5: string
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
      if (code_naf) {
        outputForKey.libelle_naf = naf.n1[code_naf]
      }
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
