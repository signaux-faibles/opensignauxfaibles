// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import "../globals"
import { map } from "./map"
import { flatten } from "../common/flatten"
import { effectifs } from "./effectifs"
import { iterable } from "./iterable"
import { sirene } from "./sirene"
import { cotisations } from "./cotisations"
import { debits } from "./debits"
import { apconso } from "./apconso"
import { delai } from "./delai"
import { compte } from "./compte"
import { dealWithProcols } from "./dealWithProcols"
import { reduce } from "./reduce"
import { finalize } from "./finalize"
import { runMongoMap } from "../test/helpers/mongodb"

const global = globalThis as any // eslint-disable-line @typescript-eslint/no-explicit-any
global.f = {
  flatten,
  effectifs,
  iterable,
  sirene,
  cotisations,
  debits,
  apconso,
  delai,
  compte,
  dealWithProcols,
}

const ISODate = (date: string): Date => new Date(date)

// test data inspired by test-api.sh
const SIREN_LENGTH = 9
const siret: SiretOrSiren = "01234567891011"
const scope: Scope = "etablissement"
const batchKey = "1910"
const dates = [
  ISODate("2015-12-01T00:00:00.000+0000"),
  ISODate("2016-01-01T00:00:00.000+0000"),
]
global.actual_batch = batchKey // used by map()
global.serie_periode = dates // used by effectifs(), which is called by map()

const rawData = {
  batch: {
    [batchKey]: {},
  },
  scope,
  index: { algo1: false, algo2: false }, // car il n'y a pas de données justifiant que l'établissement compte 10 employés ou pas
  key: siret,
}

const etablissementKey = scope + "_" + siret

const expectedMapResults = {
  // TODO: structure et valeurs à confirmer
  [etablissementKey]: {
    apconso: [],
    apdemande: [],
    batch: batchKey,
    compte: undefined,
    cotisation: dates.map(() => 0),
    debit: dates.map(() => ({ part_ouvriere: 0, part_patronale: 0 })),
    delai: [],
    dernier_effectif: undefined,
    effectif: [],
    idEntreprise: "entreprise_" + siret.substr(0, SIREN_LENGTH),
    key: siret,
    last_procol: {
      etat: "in_bonis",
    },
    procol: undefined,
    sirene: {},
  },
}

const expectedReduceResults = expectedMapResults[etablissementKey] // TODO: à confirmer

const expectedFinalizeResultValue = expectedMapResults[etablissementKey] // TODO: à confirmer

// exécution complète de la chaine "public"

test.serial(
  `public.map() retourne les propriétés d'établissement présentées sur le frontal`,
  (t: ExecutionContext) => {
    const mapResults: Record<string, unknown> = {}
    runMongoMap(map, [{ value: rawData }]).map(
      ({ _id, value }) => (mapResults[_id as string] = value)
    )
    t.deepEqual(mapResults, expectedMapResults)
  }
)

test.serial(
  `public.reduce() retourne les propriétés d'établissement, telles quelles`,
  (t: ExecutionContext) => {
    const reduceValues = [expectedMapResults[etablissementKey]]
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
