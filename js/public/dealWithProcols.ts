import { f } from "./functions"
import { ProcolToHumanRes } from "../common/procolToHuman"
import { EntréeDéfaillances, ParHash } from "../RawDataTypes"

export type SortieProcols = {
  etat: ProcolToHumanRes
  date_procol: Date
}

export function dealWithProcols(
  data_source: ParHash<EntréeDéfaillances> = {}
): SortieProcols[] {
  const events: SortieProcols[] = []
  for (const the_event of Object.values(data_source)) {
    let etat: ProcolToHumanRes = null
    etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol)

    if (etat !== null)
      events.push({ etat, date_procol: new Date(the_event.date_effet) })
  }

  return events.sort(
    (a, b) => a.date_procol.getTime() - b.date_procol.getTime()
  )
}
