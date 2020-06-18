type Output = {
  code_ape
  code_naf
  libelle_naf
  code_ape_niveau2: string
  code_ape_niveau3: string
  code_ape_niveau4: string
  libelle_ape2
  libelle_ape3
  libelle_ape4
  libelle_ape5
}

export function populateNafAndApe(
  output_indexed: { [k: string]: Output },
  naf: { n5to1; n1; n2; n3; n4; n5 }
): void
