import { EntréeDéfaillances } from "../RawDataTypes"

export type ProcolToHumanRes =
  | "liquidation"
  | "in_bonis"
  | "continuation"
  | "sauvegarde"
  | "plan_sauvegarde"
  | "plan_redressement"
  | null

type Action = EntréeDéfaillances["action_procol"]

type Stade = EntréeDéfaillances["stade_procol"]

export function procolToHuman(action: Action, stade: Stade): ProcolToHumanRes {
  "use strict"
  let res: ProcolToHumanRes = null
  if (action === "liquidation" && stade !== "abandon_procedure")
    res = "liquidation"
  else if (stade === "abandon_procedure" || stade === "fin_procedure")
    res = "in_bonis"
  else if (action === "redressement" && stade === "plan_continuation")
    res = "continuation"
  else if (action === "sauvegarde" && stade === "plan_continuation")
    res = "sauvegarde"
  else if (action === "sauvegarde") res = "plan_sauvegarde"
  else if (action === "redressement") res = "plan_redressement"
  return res
}
