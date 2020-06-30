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
  output_indexed: { [k: string]: Partial<Input> & Partial<SortieNAF> },
  naf: NAF
): void {
  "use strict"
  Object.keys(output_indexed).forEach((k) => {
    const code_ape = output_indexed[k].code_ape
    if (code_ape) {
      const code_naf = naf.n5to1[code_ape]
      output_indexed[k].code_naf = code_naf
      output_indexed[k].libelle_naf = naf.n1[code_naf]
      const code_ape_niveau2 = code_ape.substring(0, 2)
      output_indexed[k].code_ape_niveau2 = code_ape_niveau2
      const code_ape_niveau3 = code_ape.substring(0, 3)
      output_indexed[k].code_ape_niveau3 = code_ape_niveau3
      const code_ape_niveau4 = code_ape.substring(0, 4)
      output_indexed[k].code_ape_niveau4 = code_ape_niveau4
      output_indexed[k].libelle_ape2 = naf.n2[code_ape_niveau2]
      output_indexed[k].libelle_ape3 = naf.n3[code_ape_niveau3]
      output_indexed[k].libelle_ape4 = naf.n4[code_ape_niveau4]
      output_indexed[k].libelle_ape5 = naf.n5[code_ape]
    }
  })
}
