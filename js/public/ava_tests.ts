// Ces tests visent à couvrir toutes les fonctions invoquées par map(), en lui
// fournissant un jeu de données minimal mais incluant tous les types d'entrée
// inclus dans RawData.
//
// => Penser à ajouter les nouveaux types dans rawEtabData et rawEntrData.
//
// TODO: renommer ce fichier --> `map_tests.ts`.

import test, { ExecutionContext } from "ava"
import { map } from "./map"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { runMongoMap } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"
import { Scope, SiretOrSiren } from "../RawDataTypes"

// test data inspired by test.sh
const siret: SiretOrSiren = "01234567891011"
const scope: Scope = "etablissement"
const batchKey = "1910"
const dates = [
  new Date("2015-12-01T00:00:00.000+0000"),
  new Date("2016-01-01T00:00:00.000+0000"),
]
setGlobals({
  actual_batch: batchKey, // used by map()
  serie_periode: dates, // used by effectifs(), which is called by map()
})

const rawData = {
  batch: {
    [batchKey]: {},
  },
  scope,
  index: { algo2: false }, // car il n'y a pas de données justifiant que l'établissement compte 10 employés ou pas
  key: siret,
}

const etablissementKey = scope + "_" + siret

const expectedMapResults = {
  [etablissementKey]: {
    apconso: [],
    apdemande: [],
    batch: batchKey,
    compte: undefined,
    cotisation: dates.map(() => 0),
    debit_montant_majorations: dates.map(() => 0),
    debit_part_ouvriere: dates.map(() => 0),
    debit_part_patronale: dates.map(() => 0),
    delai: [],
    effectif: [null, null],
    key: siret,
    periodes: dates,
    procol: [],
    sirene: {},
  },
}

const expectedReduceResults = expectedMapResults[etablissementKey]

const expectedFinalizeResultValue = expectedMapResults[etablissementKey]

// exécution complète de la chaine "public"

test.serial(
  `public.map() retourne les propriétés d'établissement présentées sur le frontal`,
  (t: ExecutionContext) => {
    const mapResults: Record<string, unknown> = {}
    runMongoMap(map, [{ _id: null, value: rawData }]).map(
      ({ _id, value }) => (mapResults[_id as string] = value)
    )
    t.deepEqual(mapResults, expectedMapResults)
  }
)

test.serial(
  `public.reduce() retourne les propriétés d'établissement, telles quelles`,
  (t: ExecutionContext) => {
    const reduceValues = [expectedMapResults[etablissementKey] ?? {}]
    const reduceResults = reduce({ scope }, reduceValues)
    t.deepEqual(reduceResults, expectedReduceResults)
  }
)

test.serial(
  `public.finalize() retourne les propriétés d'établissement, telles quelles`,
  (t: ExecutionContext) => {
    const finalizeResultValue = finalize({ scope }, expectedReduceResults)
    t.deepEqual(finalizeResultValue, expectedFinalizeResultValue)
  }
)
