// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import "../globals"
import { map } from "./map"
import { flatten } from "../common/flatten"
import { outputs } from "./outputs"
import { repeatable } from "./repeatable"
import { add } from "./add"
import { defaillances } from "./defaillances"
import { dealWithProcols } from "./dealWithProcols"
import { populateNafAndApe } from "./populateNafAndApe"
import { cotisation } from "./cotisation"
import { dateAddMonth } from "../common/dateAddMonth"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { cibleApprentissage } from "./cibleApprentissage"
import { lookAhead } from "./lookAhead"
import { reduce } from "./reduce"
import { finalize, EntréeFinalize, EntrepriseEnEntrée } from "./finalize"
import { runMongoMap } from "../test/helpers/mongodb"

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any
global.f = {
  flatten,
  outputs,
  repeatable,
  add,
  defaillances,
  dealWithProcols,
  populateNafAndApe,
  cotisation,
  dateAddMonth,
  generatePeriodSerie,
  cibleApprentissage,
  lookAhead,
}

// test data inspired by test-api.sh
const siret: SiretOrSiren = "01234567891011"
const siren = siret.substr(0, 9)
const scope: Scope = "etablissement"
const batchKey = "1910"
const dates = [
  new Date("2015-12-01T00:00:00.000+0000"),
  new Date("2016-01-01T00:00:00.000+0000"),
]
global.actual_batch = batchKey // used by map()
global.serie_periode = dates // used by map()
global.includes = { all: true } // used by map()
global.naf = {} // used by map()

// même valeur en entrée que pour ../compact/ava_tests.ts
const rawData = {
  batch: {
    [batchKey]: {},
  },
  scope,
  index: { algo1: false, algo2: false },
  key: siret,
}

const makeValue = (periode: Date) => ({
  effectif: null,
  etat_proc_collective: "in_bonis",
  interessante_urssaf: true,
  outcome: false,
  periode,
  siret,
})

const expectedMapResults = dates.map((periode) => ({
  _id: {
    batch: batchKey,
    siren,
    periode,
    type: "other",
  },
  value: {
    [siret]: makeValue(periode),
  },
}))

const expectedReduceResults = expectedMapResults

// Structure légèrement différente de celle de test-api.golden.txt car
// l'API effectue une passe d'agrégation en plus: "cross-computation", qui est
// en cours de développement en Go (cf `reduceFinalAggregation`).
const expectedFinalizeResults = expectedMapResults.map(({ _id }) => ({
  _id,
  value: [
    // Un élément par établissement, alors que cross-computation retourne un document par établissement.
    {
      ...makeValue(_id.periode),
      nbr_etablissements_connus: 1,
    },
  ],
}))

// exécution complète de la chaine "reduce.algo2"

test.serial(
  `reduce.algo2.map() émet un objet par période`,
  (t: ExecutionContext) => {
    type MapResult = { _id: unknown; value: Record<string, EntrepriseEnEntrée> }
    const mapResults = runMongoMap(map, [{ _id: siret, value: rawData }])
    t.deepEqual(mapResults as MapResult[], expectedMapResults)
  }
)

test.serial(
  `reduce.algo2.reduce() émet un objet par période`,
  (t: ExecutionContext) => {
    const reduceResults = expectedMapResults.map(({ _id, value }) => {
      // Note: on suppose qu'il n'y a qu'une valeur par clé
      return { _id, value: reduce(_id, [value]) }
    })
    t.deepEqual(reduceResults, expectedReduceResults)
  }
)

test.serial(
  `reduce.algo2.finalize() émet un objet par période`,
  (t: ExecutionContext) => {
    const finalizeResult = expectedReduceResults.map(({ _id, value }) => {
      // Note: on suppose qu'il n'y a qu'une valeur par clé
      return {
        _id,
        value: finalize(_id, value as EntréeFinalize),
      }
    })
    t.deepEqual(finalizeResult, expectedFinalizeResults)
  }
)
