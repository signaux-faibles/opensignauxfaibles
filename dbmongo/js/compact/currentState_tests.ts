import test, { ExecutionContext } from "ava"
import "../globals"
import { currentState } from "./currentState"

const makeApDemande = (): EntréeApDemande => ({
  id_demande: "",
  periode: { start: new Date(), end: new Date() },
  hta: 0,
  motif_recours_se: 0,
})

const makeApConso = (): EntréeApConso => ({
  id_conso: "",
  periode: new Date(),
  heure_consomme: 0,
})

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
