import { f } from "./functions"
import { SortieAPart } from "./apart"
import { SortieRepeatable } from "./repeatable"
import { SortieCotisationsDettes } from "./cotisationsdettes"
import { SortieEffectifs } from "./effectifs"
import { SortieDefaillances } from "./defaillances"
import { SortieCcsf } from "./ccsf"
import { SortieSirene } from "./sirene"
import { SortieNAF } from "./populateNafAndApe"
import { SortieDelais } from "./delais"
import { SortieCibleApprentissage } from "./cibleApprentissage"
import { SortieCotisation } from "./cotisation"
import { SortieCompte } from "./compte"
import { Siret, Siren } from "../RawDataTypes"
import { ParPériode } from "../common/makePeriodeMap"

export type DonnéesAgrégées = {
  siret: Siret
  periode: Date
} & Partial<SortieCotisationsDettes> &
  Partial<SortieEffectifs<"effectif">> &
  Partial<SortieDefaillances> &
  Partial<SortieCcsf> &
  Partial<SortieSirene> &
  Partial<SortieNAF> &
  Partial<SortieAPart> &
  Partial<SortieRepeatable> &
  Partial<SortieDelais> &
  Partial<SortieCotisation> &
  Partial<SortieCompte> &
  Partial<SortieCibleApprentissage> &
  Partial<SortieCotisationsDettes>

/**
 * Appelé par `map()` pour chaque entreprise/établissement, `outputs()` retourne
 * un tableau contenant un objet de base par période, ainsi qu'une version
 * indexée par période de ce tableau, afin de faciliter l'agrégation progressive
 * de données dans ces structures par `map()`.
 */
export function outputs(
  v: { key: Siret | Siren },
  serie_periode: Date[]
): [DonnéesAgrégées[], ParPériode<DonnéesAgrégées>] {
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

  const output_indexed: ParPériode<DonnéesAgrégées> = f.makePeriodeMap<DonnéesAgrégées>()
  for (const val of output_array) {
    output_indexed.set(val.periode, val)
  }

  return [output_array, output_indexed]
}
