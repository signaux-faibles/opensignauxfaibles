type CodeAPENiveau2 = string
type CodeAPENiveau3 = string
type CodeAPENiveau4 = string

export type Output = {
  code_ape: CodeAPE
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

export type NAF = {
  n5to1: Record<CodeAPE, CodeNAF>
  n1: Record<CodeNAF, string>
  n2: Record<CodeAPENiveau2, string>
  n3: Record<CodeAPENiveau3, string>
  n4: Record<CodeAPENiveau4, string>
  n5: Record<CodeAPE, string>
}

export function populateNafAndApe(
  output_indexed: { [k: string]: Output },
  naf: NAF
): void {
  "use strict"
  Object.keys(output_indexed).forEach((k) => {
    if (
      "code_ape" in output_indexed[k] &&
      output_indexed[k].code_ape !== null
    ) {
      const code_ape = output_indexed[k].code_ape
      output_indexed[k].code_naf = naf.n5to1[code_ape]
      output_indexed[k].libelle_naf = naf.n1[output_indexed[k].code_naf]
      output_indexed[k].code_ape_niveau2 = code_ape.substring(0, 2)
      output_indexed[k].code_ape_niveau3 = code_ape.substring(0, 3)
      output_indexed[k].code_ape_niveau4 = code_ape.substring(0, 4)
      output_indexed[k].libelle_ape2 =
        naf.n2[output_indexed[k].code_ape_niveau2]
      output_indexed[k].libelle_ape3 =
        naf.n3[output_indexed[k].code_ape_niveau3]
      output_indexed[k].libelle_ape4 =
        naf.n4[output_indexed[k].code_ape_niveau4]
      output_indexed[k].libelle_ape5 = naf.n5[code_ape]
    }
  })
}
