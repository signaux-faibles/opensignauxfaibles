import test from "ava"
import { defaillances, SortieDefaillances } from "./defaillances"
import { EntréeDéfaillances } from "../GeneratedTypes"
import { ParHash } from "../RawDataTypes"
import { parPériode } from "../test/helpers/parPeriode"
import { ParPériode } from "../common/makePeriodeMap"

type OutputIndexed = Parameters<typeof defaillances>[1]

test("Une ouverture de liquidation est prise en compte dans la période courante et les suivantes", (t: test) => {
  const output_indexed: OutputIndexed = parPériode({
    ["2018-01-01"]: {},
    ["2018-02-01"]: {},
    ["2018-03-01"]: {},
  })

  const date_ouverture = new Date("2018-02-12")
  const data_source = {
    ["123"]: {
      action_procol: "liquidation",
      stade_procol: "ouverture",
      date_effet: date_ouverture,
    },
  } as ParHash<EntréeDéfaillances>

  defaillances(data_source, output_indexed)
  const expected: ParPériode<Partial<SortieDefaillances>> = parPériode({
    ["2018-01-01"]: {},
    ["2018-02-01"]: {
      date_proc_collective: date_ouverture,
      etat_proc_collective: "liquidation",
      tag_failure: true,
    },
    ["2018-03-01"]: {
      date_proc_collective: date_ouverture,
      etat_proc_collective: "liquidation",
      tag_failure: true,
    },
  })

  t.deepEqual(output_indexed, expected)
})

test("Une ouverture puis cloture d'un redressement sont pris en compte, tag_failure reste à TRUE", (t: test) => {
  const output_indexed: OutputIndexed = parPériode({
    ["2018-01-01"]: {},
    ["2018-02-01"]: {},
    ["2018-03-01"]: {},
  })

  const date_ouverture = new Date("2018-02-12")
  const date_cloture = new Date("2018-03-05")

  const data_source = {
    ["123"]: {
      action_procol: "redressement",
      stade_procol: "ouverture",
      date_effet: date_ouverture,
    },
    ["456"]: {
      action_procol: "redressement",
      stade_procol: "fin_procedure",
      date_effet: date_cloture,
    },
  } as ParHash<EntréeDéfaillances>

  defaillances(data_source, output_indexed)
  const expected: ParPériode<Partial<SortieDefaillances>> = parPériode({
    ["2018-01-01"]: {},
    ["2018-02-01"]: {
      date_proc_collective: date_ouverture,
      etat_proc_collective: "plan_redressement",
      tag_failure: true,
    },
    ["2018-03-01"]: {
      date_proc_collective: date_cloture,
      etat_proc_collective: "in_bonis",
      tag_failure: true,
    },
  })

  t.deepEqual(output_indexed, expected)
})
