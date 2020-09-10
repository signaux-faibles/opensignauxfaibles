// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import { map, EntréeMap, CléSortieMap, SortieMap } from "./map"
import { reduce } from "./reduce"
import { finalize, EntrepriseEnSortie } from "./finalize"
import { setGlobals } from "../test/helpers/setGlobals"
import { runMongoMap } from "../test/helpers/mongodb"
import { Scope, SiretOrSiren } from "../RawDataTypes"

// test data inspired by test-api.sh
const siret: SiretOrSiren = "01234567891011"
const siren = siret.substr(0, 9)
const scope: Scope = "etablissement"
const batchKey = "1910"
const dates = [
  new Date("2015-12-01T00:00:00.000+0000"),
  new Date("2016-01-01T00:00:00.000+0000"),
]

setGlobals({
  // used by map()
  actual_batch: batchKey,
  serie_periode: dates,
  includes: { all: true },
  naf: {},
})

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
  } as SortieMap,
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
    } as EntrepriseEnSortie,
  ],
}))

// exécution complète de la chaine "reduce.algo2"

test.serial(
  `reduce.algo2.map() émet un objet par période`,
  (t: ExecutionContext) => {
    const mapResults = runMongoMap<EntréeMap, CléSortieMap, SortieMap>(map, [
      { _id: siret, value: rawData },
    ])
    t.deepEqual(mapResults, expectedMapResults)
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
        value: finalize(_id, value),
      }
    })
    t.deepEqual(finalizeResult, expectedFinalizeResults)
  }
)
