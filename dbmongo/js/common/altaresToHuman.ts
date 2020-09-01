export type AltaresToHumanRes =
  | "liquidation"
  | "in_bonis"
  | "continuation"
  | "sauvegarde"
  | "plan_sauvegarde"
  | "plan_redressement"
  | "cession"
  | null

type AltaresCode = EntréeDefaillances["code_evenement"]

export function altaresToHuman(code: AltaresCode): AltaresToHumanRes {
  "use strict"
  const codeLiquidation = [
    "PCL0108",
    "PCL010801",
    "PCL010802",
    "PCL030107",
    "PCL030307",
    "PCL030311",
    "PCL05010103",
    "PCL05010204",
    "PCL05010303",
    "PCL05010403",
    "PCL05010503",
    "PCL05010703",
    "PCL05011004",
    "PCL05011102",
    "PCL05011204",
    "PCL05011206",
    "PCL05011304",
    "PCL05011404",
    "PCL05011504",
    "PCL05011604",
    "PCL05011903",
    "PCL05012004",
    "PCL050204",
    "PCL0109",
    "PCL010901",
    "PCL030108",
    "PCL030308",
    "PCL05010104",
    "PCL05010205",
    "PCL05010304",
    "PCL05010404",
    "PCL05010504",
    "PCL05010803",
    "PCL05011005",
    "PCL05011103",
    "PCL05011205",
    "PCL05011207",
    "PCL05011305",
    "PCL05011405",
    "PCL05011505",
    "PCL05011904",
    "PCL05011605",
    "PCL05012005",
  ]
  const codePlanSauvegarde = [
    "PCL010601",
    "PCL0106",
    "PCL010602",
    "PCL030103",
    "PCL030303",
    "PCL03030301",
    "PCL05010101",
    "PCL05010202",
    "PCL05010301",
    "PCL05010401",
    "PCL05010501",
    "PCL05010506",
    "PCL05010701",
    "PCL05010705",
    "PCL05010801",
    "PCL05010805",
    "PCL05011002",
    "PCL05011202",
    "PCL05011302",
    "PCL05011402",
    "PCL05011502",
    "PCL05011602",
    "PCL05011901",
    "PCL0114",
    "PCL030110",
    "PCL030310",
  ]
  const codeRedressement = [
    "PCL0105",
    "PCL010501",
    "PCL010502",
    "PCL010503",
    "PCL030105",
    "PCL030305",
    "PCL05010102",
    "PCL05010203",
    "PCL05010302",
    "PCL05010402",
    "PCL05010502",
    "PCL05010702",
    "PCL05010706",
    "PCL05010802",
    "PCL05010806",
    "PCL05010901",
    "PCL05011003",
    "PCL05011101",
    "PCL05011203",
    "PCL05011303",
    "PCL05011403",
    "PCL05011503",
    "PCL05011603",
    "PCL05011902",
    "PCL05012003",
  ]
  const codeInBonis = [
    "PCL05",
    "PCL0501",
    "PCL050101",
    "PCL050102",
    "PCL050103",
    "PCL050104",
    "PCL050105",
    "PCL050106",
    "PCL050107",
    "PCL050108",
    "PCL050109",
    "PCL050110",
    "PCL050111",
    "PCL050112",
    "PCL050113",
    "PCL050114",
    "PCL050115",
    "PCL050116",
    "PCL050119",
    "PCL050120",
    "PCL050121",
    "PCL0503",
    "PCL050301",
    "PCL050302",
    "PCL0508",
    "PCL010504",
    "PCL010803",
    "PCL010902",
    "PCL050901",
    "PCL050902",
    "PCL050903",
    "PCL050904",
    "PCL0504",
    "PCL050303",
    "PCL050401",
    "PCL050402",
    "PCL050403",
    "PCL050404",
    "PCL050405",
    "PCL050406",
  ]
  const codeContinuation = ["PCL0202"]
  const codeSauvegarde = ["PCL0203", "PCL020301", "PCL0205", "PCL040408"]
  const codeCession = ["PCL0204", "PCL020401", "PCL020402", "PCL020403"]
  let res: AltaresToHumanRes = null
  if (codeLiquidation.includes(code)) res = "liquidation"
  else if (codePlanSauvegarde.includes(code)) res = "plan_sauvegarde"
  else if (codeRedressement.includes(code)) res = "plan_redressement"
  else if (codeInBonis.includes(code)) res = "in_bonis"
  else if (codeContinuation.includes(code)) res = "continuation"
  else if (codeSauvegarde.includes(code)) res = "sauvegarde"
  else if (codeCession.includes(code)) res = "cession"
  return res
}
