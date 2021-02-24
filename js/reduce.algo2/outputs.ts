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
import { SiretOrSiren, ParPériode } from "../RawDataTypes"

export type DonnéesAgrégées = {
  siret: SiretOrSiren
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
  v: { key: SiretOrSiren },
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

  const rawOutputIndexed = output_array.reduce(function (periodes, val) {
    periodes[val.periode.getTime()] = val
    return periodes
  }, {} as ParPériode<DonnéesAgrégées>)

  type OutputIndexed = typeof rawOutputIndexed

  const validator = {
    set(
      obj: OutputIndexed,
      prop: keyof OutputIndexed,
      value: OutputIndexed[keyof OutputIndexed]
    ): boolean {
      const timestamp = parseInt(prop, 10)
      if (isNaN(timestamp) || new Date(timestamp).getTime() !== timestamp) {
        throw new RangeError("output_indexed only accepts timestamps as keys")
      }
      obj[prop] = value // The default behavior to store the value
      return true // Indicate success
    },
  }

  const output_indexed = new Proxy(rawOutputIndexed, validator)

  // output_indexed["abd"] = {} as DonnéesAgrégées // => npm test fails with "output_indexed only accepts timestamps as keys" (at runtime)

  return [output_array, output_indexed]
}
