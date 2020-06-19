import { altaresToHuman, AltaresToHumanRes } from "../common/altaresToHuman"
import { procolToHuman, ProcolToHumanRes } from "./procolToHuman"

export type InputEvent = Événement

type OutputEvent = {
  etat: AltaresToHumanRes
  date_proc_col: Date
}

export type Output = {
  etat_proc_collective: unknown
  date_proc_collective: unknown
  tag_failure: boolean
}

export function dealWithProcols(
  data_source: { [hash: string]: InputEvent },
  altar_or_procol: "altares" | "procol",
  output_indexed: {
    [time: string]: Output
  }
): void {
  "use strict"
  const f = { altaresToHuman, procolToHuman } // DO_NOT_INCLUDE_IN_JSFUNCTIONS_GO
  const codes = Object.keys(data_source)
    .reduce((events, hash) => {
      const the_event = data_source[hash]

      let etat: AltaresToHumanRes | ProcolToHumanRes = null
      if (altar_or_procol == "altares")
        etat = f.altaresToHuman(the_event.code_evenement)
      else if (altar_or_procol == "procol")
        etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol)

      if (etat != null)
        events.push({
          etat: etat,
          date_proc_col: new Date(the_event.date_effet),
        })

      return events
    }, [] as OutputEvent[])
    .sort((a, b) => {
      return a.date_proc_col.getTime() > b.date_proc_col.getTime() ? 1 : 0
    })

  codes.forEach((event) => {
    const periode_effet = new Date(
      Date.UTC(
        event.date_proc_col.getFullYear(),
        event.date_proc_col.getUTCMonth(),
        1,
        0,
        0,
        0,
        0
      )
    )
    const time_til_last = Object.keys(output_indexed).filter((val) => {
      return val >= periode_effet.toString()
    })

    time_til_last.forEach((time) => {
      if (time in output_indexed) {
        output_indexed[time].etat_proc_collective = event.etat
        output_indexed[time].date_proc_collective = event.date_proc_col
        if (event.etat != "in_bonis") output_indexed[time].tag_failure = true
      }
    })
  })
}
