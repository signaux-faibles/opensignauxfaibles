import { DebitComputedValues } from "./delais"
import { Output as DefaillancesOutput } from "./defaillances"
import { Output as CcsfOutput } from "./ccsf"
import { Output as SireneOutput } from "./sirene"
import { Output as NafOutput } from "./populateNafAndApe"
import { Output as CotisationOutput } from "./cotisation"
// import { Output as DealWithProcolsOutput } from "./dealWithProcols"

type DonnéesAgrégées = {
  siret: SiretOrSiren
  periode: Date
  effectif: number | null
  etat_proc_collective: "in_bonis" // ou ProcolToHumanRes ?
  interessante_urssaf: true
  outcome: false
} & DebitComputedValues &
  Partial<DefaillancesOutput> &
  Partial<CcsfOutput> &
  Partial<SireneOutput> &
  Partial<NafOutput> &
  Partial<CotisationOutput>

type IndexedOutput = Record<Periode, DonnéesAgrégées>

/**
 * Appelé par `map()` pour chaque entreprise/établissement, `outputs()` retourne
 * un tableau contenant un objet de base par période, ainsi qu'une version
 * indexée par période de ce tableau, afin de faciliter l'agrégation progressive
 * de données dans ces structures par `map()`.
 */
export function outputs(
  v: { key: SiretOrSiren },
  serie_periode: Date[]
): [DonnéesAgrégées[], IndexedOutput] {
  "use strict"
  const output_array: DonnéesAgrégées[] = serie_periode.map(function (e) {
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
