import { SortieAPart } from "./apart"
import { SortieCotisationsDettes } from "./cotisationsdettes"
import { SortieDefaillances } from "./defaillances"
import { SortieCcsf } from "./ccsf"
import { SortieSirene } from "./sirene"
import { SortieNAF } from "./populateNafAndApe"
import { SortieCotisation } from "./cotisation"

export type DonnéesAgrégées = {
  siret: SiretOrSiren
  periode: Date
  effectif: number | null
  etat_proc_collective: "in_bonis" // ou ProcolToHumanRes ?
  interessante_urssaf: true
  outcome: false
} & Partial<SortieCotisationsDettes> &
  Partial<SortieDefaillances> &
  Partial<SortieCcsf> &
  Partial<SortieSirene> &
  Partial<SortieNAF> &
  Partial<SortieAPart> &
  Partial<SortieCotisation>

type IndexDonnéesAgrégées = Record<Timestamp, DonnéesAgrégées>

/**
 * Appelé par `map()` pour chaque entreprise/établissement, `outputs()` retourne
 * un tableau contenant un objet de base par période, ainsi qu'une version
 * indexée par période de ce tableau, afin de faciliter l'agrégation progressive
 * de données dans ces structures par `map()`.
 */
export function outputs(
  v: { key: SiretOrSiren },
  serie_periode: Date[]
): [DonnéesAgrégées[], IndexDonnéesAgrégées] {
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
  }, {} as IndexDonnéesAgrégées)

  return [output_array, output_indexed]
}
