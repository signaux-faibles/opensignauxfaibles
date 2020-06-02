type Output = {
  siret: SiretOrSiren
  periode: Periode
  effectif: null
  etat_proc_collective: "in_bonis"
  interessante_urssaf: true
  outcome: false
}

export function outputs(
  v: { key: SiretOrSiren },
  serie_periode: Periode[]
): [Output[], { [k: string]: Output }]
