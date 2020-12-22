import { f } from "./functions"
import { EntréeInterim, ParPériode, ParHash } from "../RawDataTypes"

type Input = {
  effectif: number | null
}

type MonthOffset = 6 | 12 | 18 | 24
type CléSortieInterim = `interim_ratio_past_${MonthOffset}`
type SortieInterim = {
  interim_proportion: number
} & {
  [K in CléSortieInterim]: number
}

export function interim(
  interim: ParHash<EntréeInterim>,
  output_indexed: ParPériode<Input>
): ParPériode<SortieInterim> {
  "use strict"

  const output_effectif = output_indexed
  // let periodes = Object.keys(output_indexed)
  // output_indexed devra être remplacé par output_effectif, et ne contenir que les données d'effectif.
  // periodes sera passé en argument.

  const output_interim: ParPériode<SortieInterim> = {}

  //  var offset_interim = 3

  for (const one_interim of Object.values(interim)) {
    const periode = one_interim.periode.getTime()
    // var periode_d = new Date(parseInt(interimTime))
    // var time_offset = f.dateAddMonth(time_d, -offset_interim)
    if (periode in output_effectif) {
      const out = output_interim[periode] ?? ({} as SortieInterim)
      const { effectif } = output_effectif[periode] ?? {}
      if (effectif) {
        out.interim_proportion = one_interim.etp / effectif
      }
      output_interim[periode] = out
    }

    const makePastProp = (offset: MonthOffset) =>
      `interim_ratio_past_${offset}` as CléSortieInterim

    const past_month_offsets: MonthOffset[] = [6, 12, 18, 24]
    past_month_offsets.forEach((offset) => {
      const pastPropName = makePastProp(offset)
      const time_past_offset = f.dateAddMonth(one_interim.periode, offset)
      if (
        periode in output_effectif &&
        time_past_offset.getTime() in output_effectif
      ) {
        const out =
          output_interim[time_past_offset.getTime()] ?? ({} as SortieInterim)
        const val_offset = output_interim[time_past_offset.getTime()]
        const { effectif } = output_effectif[periode] ?? {}
        if (effectif) {
          Object.assign(val_offset, {
            [pastPropName]: one_interim.etp / effectif,
          })
        }
        output_interim[time_past_offset.getTime()] = out
      }
    })
  }

  return output_interim
}
