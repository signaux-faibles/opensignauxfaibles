import { f } from "./functions"
import { ParPériode } from "../common/makePeriodeMap"
import { EntréeEffectif, EntréeEffectifEnt } from "../GeneratedTypes"
import { ParHash, Timestamp } from "../RawDataTypes"

// Paramètres globaux utilisés par "reduce.algo2"
declare const offset_effectif: number

type Clé = "effectif_ent" | "effectif" // effectif entreprise ou établissement
type ValeurEffectif = number

type EffectifReporté<K extends Clé> = K extends "effectif_ent"
  ? {
      /** Vaut 1 si cette valeur est absente pour la période donnée et que le dernier effectif connu a été utilisé à la place. */
      effectif_ent_reporte: 1 | 0
    }
  : {
      /** Vaut 1 si cette valeur est absente pour la période donnée et que le dernier effectif connu a été utilisé à la place. */
      effectif_reporte: 1 | 0
    }

type MonthOffset = 6 | 12 | 18 | 24
type EffectifPassé<K extends Clé> = {
  /** Valeur d'effectif il y a N mois. */
  [k in `${K}_past_${MonthOffset}`]: ValeurEffectif
}

type ValeursTransmisesEtab = {
  /** Nombre de personnes employées par l'établissement. */
  effectif: ValeurEffectif | null
}

type ValeursTransmisesEntr = {
  /** Nombre de personnes employées par l'entreprise. */
  effectif_ent: ValeurEffectif | null
}

export type ValeursTransmises<K extends Clé> = K extends "effectif_ent"
  ? ValeursTransmisesEntr
  : ValeursTransmisesEtab

type ValeursCalculées<K extends Clé> = EffectifReporté<K> & EffectifPassé<K>

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "effectifs"
  computed: ValeursCalculées<"effectif"> | ValeursCalculées<"effectif_ent">
  transmitted: ValeursTransmises<"effectif"> | ValeursTransmises<"effectif_ent">
}

export type SortieEffectifs<K extends Clé> = ValeursTransmises<K> &
  ValeursCalculées<K>

export function effectifs<K extends Clé>(
  entréeEffectif: ParHash<
    K extends "effectif_ent" ? EntréeEffectifEnt : EntréeEffectif
  >,
  periodes: Timestamp[],
  clé: K
): ParPériode<SortieEffectifs<K>> {
  "use strict"

  const sortieEffectif = f.makePeriodeMap<SortieEffectifs<K>>()

  // Construction d'une map[time] = effectif à cette periode
  const mapEffectif = f.makePeriodeMap<ValeurEffectif>()

  Object.keys(entréeEffectif).forEach((hash) => {
    const effectif = entréeEffectif[hash]
    if (effectif !== null && effectif !== undefined) {
      mapEffectif.set(effectif.periode, effectif.effectif)
    }
  })

  // On reporte dans les dernières périodes le dernier effectif connu
  // Ne reporter que si le dernier effectif est disponible
  const dernièrePériodeAvecEffectifConnu = f.dateAddMonth(
    new Date(periodes[periodes.length - 1] as number),
    offset_effectif + 1
  )
  const effectifÀReporter =
    mapEffectif.get(dernièrePériodeAvecEffectifConnu) ?? null

  const makeReporteProp = (clé: K) => `${clé}_reporte`

  periodes.forEach((time) => {
    sortieEffectif.set(time, {
      ...(sortieEffectif.get(time) ?? ({} as SortieEffectifs<K>)),
      [clé]: mapEffectif.get(time) || effectifÀReporter,
      [makeReporteProp(clé)]: mapEffectif.get(time) ? 0 : 1,
    })
  })

  const makePastProp = (clé: K, offset: MonthOffset) => `${clé}_past_${offset}`

  mapEffectif.forEach((effectifAtTime, time) => {
    const futureOffsets: MonthOffset[] = [6, 12, 18, 24]
    const futureTimestamps = futureOffsets
      .map((offset) => ({
        offset,
        timestamp: f
          .dateAddMonth(new Date(time), offset - offset_effectif - 1)
          // TODO: réfléchir à si l'offset est nécessaire pour l'algo.
          // Ces valeurs permettent de calculer les dernières variations réelles
          // d'effectif sur la période donnée (par exemple: 6 mois),
          // en excluant les valeurs reportées qui
          // pourraient conduire à des variations = 0
          .getTime(),
      }))
      .filter(({ timestamp }) => periodes.includes(timestamp))

    futureTimestamps.forEach(({ offset, timestamp }) => {
      sortieEffectif.set(timestamp, {
        ...(sortieEffectif.get(timestamp) ?? ({} as SortieEffectifs<K>)),
        [makePastProp(clé, offset)]: effectifAtTime,
      })
    })
  })
  return sortieEffectif
}
