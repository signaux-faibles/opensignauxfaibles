import { f } from "./functions"
import { ProcolToHumanRes } from "../common/procolToHuman"
import { EntréeDefaillances, ParHash } from "../RawDataTypes"

export type SortieProcols = {
  etat: ProcolToHumanRes
  date_procol: Date
}

export function dealWithProcols(
  data_source: ParHash<EntréeDefaillances> = {}
): SortieProcols[] {
  return Object.keys(data_source)
    .reduce((events, hash) => {
      const the_event = data_source[hash]

      let etat: ProcolToHumanRes = null
      etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol)

      if (etat !== null)
        events.push({ etat, date_procol: new Date(the_event.date_effet) })

      return events
    }, [] as SortieProcols[])
    .sort((a, b) => a.date_procol.getTime() - b.date_procol.getTime())
}
