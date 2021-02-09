import test, { ExecutionContext } from "ava"
import { currentState } from "./currentState"
import { EntréeApConso, EntréeApDemande } from "../GeneratedTypes"

const makeApDemande = () =>
  ({
    id_demande: "",
    periode: { start: new Date(), end: new Date() },
    hta: 0,
    motif_recours_se: 0,
  } as EntréeApDemande)

const makeApConso = () =>
  ({
    id_conso: "",
    periode: new Date(),
    heure_consomme: 0,
  } as EntréeApConso)

test("currentState", (t: ExecutionContext) => {
  const actualRes = currentState([
    {
      apdemande: { a: makeApDemande(), b: makeApDemande() },
      apconso: { c: makeApConso() },
    },
    {
      apdemande: { d: makeApDemande() },
      compact: { delete: { apdemande: ["a", "b"] } },
    },
  ])
  const expectedRes = {
    apdemande: new Set(["d"]),
    apconso: new Set(["c"]),
  }
  t.deepEqual(actualRes, expectedRes)
})
