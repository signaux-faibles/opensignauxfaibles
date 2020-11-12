import test from "ava"
import { dealWithProcols, InputEvent } from "./dealWithProcols"
import { ParHash } from "../RawDataTypes"



test("Une entrée en liquidation est prise en compte dans la période courante et les suivantes", (t) => {
  const output_indexed =  {
    ["2018-01-01"]: {},
    ["2018-02-01"]: {},
    ["2018-03-01"]: {}
  }

  const date_proc_collective = new Date("2018-02-12")
  const data_source = {
    ["123"]: {
      action_procol: "liquidation",
      stade_procol: "ouverture",
      date_effet: date_proc_collective
    }
  } as ParHash<InputEvent>

  dealWithProcols(data_source, output_indexed)
  const expected = {
    ["2018-01-01"]: {},
    ["2018-02-01"]: {
      date_proc_collective: date_proc_collective,
      etat_proc_collective: "liquidation",
      tag_failure: true
    },
    ["2018-03-01"]: {
      date_proc_collective: date_proc_collective,
      etat_proc_collective: "liquidation",
      tag_failure: true
    }
  }

  t.deepEqual(output_indexed, expected)
})
