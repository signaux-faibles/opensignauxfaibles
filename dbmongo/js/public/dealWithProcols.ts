import { altaresToHuman, AltaresToHumanRes } from "../common/altaresToHuman"
import { procolToHuman, ProcolToHumanRes } from "../common/procolToHuman"

export type SortieProcols = {
  etat: AltaresToHumanRes | ProcolToHumanRes
  date_procol: Date
}

export function dealWithProcols(
  data_source: Record<DataHash, EntréeDefaillances>,
  altar_or_procol: "altares" | "procol"
): SortieProcols[] {
  const f = { altaresToHuman, procolToHuman } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO

  return Object.keys(data_source || {})
    .reduce((events, hash) => {
      const the_event = data_source[hash]

      let etat: AltaresToHumanRes | ProcolToHumanRes = null
      if (altar_or_procol === "altares")
        etat = f.altaresToHuman(the_event.code_evenement)
      else if (altar_or_procol === "procol")
        etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol)

      if (etat !== null)
        events.push({ etat, date_procol: new Date(the_event.date_effet) })

      return events
    }, [] as SortieProcols[])
    .sort((a, b) => a.date_procol.getTime() - b.date_procol.getTime())
}
