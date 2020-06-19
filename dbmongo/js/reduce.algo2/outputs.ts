type Output = {
  siret: SiretOrSiren
  periode: Date
  effectif: null
  etat_proc_collective: "in_bonis" // ou ProcolToHumanRes ?
  interessante_urssaf: true
  outcome: false
}

type IndexedOutput = Record<string, Output>

export function outputs(
  v: { key: SiretOrSiren },
  serie_periode: Date[]
): [Output[], IndexedOutput] {
  "use strict"
  const output_array: Output[] = serie_periode.map(function (e) {
    return {
      siret: v.key,
      periode: e,
      effectif: null,
      etat_proc_collective: "in_bonis",
      interessante_urssaf: true,
      outcome: false,
    }
  })

  const output_indexed = output_array.reduce(function (periodes, val) {
    periodes[val.periode.getTime()] = val
    return periodes
  }, {} as IndexedOutput)

  return [output_array, output_indexed]
}
