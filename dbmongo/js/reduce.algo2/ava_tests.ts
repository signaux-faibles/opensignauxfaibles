// Objectif de cette suite de tests d'intégration:
// Vérifier la compatibilité des types et mesurer la couverture lors du passage
// de données entre les fonctions map(), reduce() et finalize(), en s'appuyant
// sur le jeu de données minimal utilisé dans notre suite de bout en bout
// définie dans test-api.sh.

import test, { ExecutionContext } from "ava"
import "../globals"
import { map } from "./map"
import { flatten } from "./flatten.js"
import { outputs } from "./outputs.js"
import { repeatable } from "./repeatable.js"
import { add } from "./add.js"
import { defaillances } from "./defaillances.js"
import { dealWithProcols } from "./dealWithProcols.js"
import { populateNafAndApe } from "./populateNafAndApe.js"
import { cotisation } from "./cotisation.js"
import { dateAddMonth } from "./dateAddMonth.js"
import { generatePeriodSerie } from "../common/generatePeriodSerie.js"
import { cibleApprentissage } from "./cibleApprentissage.js"
import { lookAhead } from "./lookAhead.js"
import { reduce } from "./reduce.js"
import { finalize } from "./finalize.js"

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

const ISODate = (date: string): Date => new Date(date)

;(Object as any).bsonsize = (obj: any): number => JSON.stringify(obj).length // used by finalize()

const runMongoMap = (mapFct: () => void, keyVal: any): any => {
  const results: { _id: any; value: any }[] = []
  globalThis.emit = (key: any, value: any): void => {
    results.push({ _id: key, value })
  }
  mapFct.call(keyVal)
  return results
}

// test data inspired by test-api.sh
const siret: SiretOrSiren = "01234567891011"
const scope: Scope = "etablissement"
const batchKey = "1910"
const dates = [
  ISODate("2015-12-01T00:00:00.000+0000"),
  ISODate("2016-01-01T00:00:00.000+0000"),
]
global.actual_batch = batchKey // used by map()
global.serie_periode = dates // used by map()
global.includes = { all: true } // used by map()
global.naf = {} // used by map()

// même valeur en entrée que pour ../compact/ava_tests.ts
const rawData = {
  batch: {
    [batchKey]: {
      reporder: dates.reduce(
        (reporder, date) => ({
          ...reporder,
          [date.toString()]: { periode: date, siret },
        }),
        {}
      ),
    } as BatchValue, // TODO: rendre optionnelles les props de BatchValues, pour retirer ce cast
  },
  scope,
  index: { algo1: false, algo2: false }, // car il n'y a pas de données justifiant que l'établissement compte 10 employés ou pas
  key: siret,
}

// valeurs résultantes de l'exécution de map() => à vérifier et à ré-écrire de manière plus concise
const expectedMapResults = [
  {
    _id: {
      batch: "1910",
      siren: "012345678",
      periode: new Date("2015-12-01 00:00:00 UTC"),
      type: "other",
    },
    value: {
      "01234567891011": {
        cotisation_moy12m: 0,
        effectif: null,
        etat_proc_collective: "in_bonis",
        interessante_urssaf: true,
        outcome: false,
        periode: new Date("2015-12-01 00:00:00 UTC"),
        random_order: undefined,
        siret: "01234567891011",
      },
    },
  },
  {
    _id: {
      batch: "1910",
      siren: "012345678",
      periode: new Date("2016-01-01 00:00:00 UTC"),
      type: "other",
    },
    value: {
      "01234567891011": {
        cotisation_moy12m: 0,
        effectif: null,
        etat_proc_collective: "in_bonis",
        interessante_urssaf: true,
        outcome: false,
        periode: new Date("2016-01-01 00:00:00 UTC"),
        random_order: undefined,
        siret: "01234567891011",
      },
    },
  },
]

// valeurs résultantes de l'exécution de reduce() => à vérifier et à ré-écrire de manière plus concise
const expectedReduceResults = [
  {
    _id: {
      batch: "1910",
      periode: new Date("2015-12-01 00:00:00 UTC"),
      siren: "012345678",
      type: "other",
    },
    value: {
      "01234567891011": {
        cotisation_moy12m: 0,
        effectif: null,
        etat_proc_collective: "in_bonis",
        interessante_urssaf: true,
        outcome: false,
        periode: new Date("2015-12-01 00:00:00 UTC"),
        random_order: undefined,
        siret: "01234567891011",
      },
    },
  },
  {
    _id: {
      batch: "1910",
      periode: new Date("2016-01-01 00:00:00 UTC"),
      siren: "012345678",
      type: "other",
    },
    value: {
      "01234567891011": {
        cotisation_moy12m: 0,
        effectif: null,
        etat_proc_collective: "in_bonis",
        interessante_urssaf: true,
        outcome: false,
        periode: new Date("2016-01-01 00:00:00 UTC"),
        random_order: undefined,
        siret: "01234567891011",
      },
    },
  },
]

// extrait de test-api.golden-master.txt, pour les dates spécifiées plus haut
const expectedFinalizeResults = [
  {
    _id: {
      batch: "1910",
      siret: "01234567891011",
      periode: ISODate("2015-12-01T00:00:00Z"),
    },
    value: {
      siret: "01234567891011",
      periode: ISODate("2015-12-01T00:00:00Z"),
      effectif: null,
      etat_proc_collective: "in_bonis",
      interessante_urssaf: true,
      outcome: false,
      cotisation_moy12m: 0,
      nbr_etablissements_connus: 1,
    },
  },
  {
    _id: {
      batch: "1910",
      siret: "01234567891011",
      periode: ISODate("2016-01-01T00:00:00Z"),
    },
    value: {
      siret: "01234567891011",
      periode: ISODate("2016-01-01T00:00:00Z"),
      effectif: null,
      etat_proc_collective: "in_bonis",
      interessante_urssaf: true,
      outcome: false,
      cotisation_moy12m: 0,
      nbr_etablissements_connus: 1,
    },
  },
]

// exécution complète de la chaine "reduce.algo2"

test.serial(
  `reduce.algo2.map() émet un objet par période`,
  (t: ExecutionContext) => {
    const mapResults = runMongoMap(map, { _id: siret, value: rawData })
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

test.serial(`reduce.algo2.finalize()`, (t: ExecutionContext) => {
  const finalizeResult = expectedReduceResults.map(({ _id, value }) => {
    // Note: on suppose qu'il n'y a qu'une valeur par clé
    const result = finalize(_id, value)
    return { _id, value: "incomplete" in result ? result : result[0] } // TODO: pourquoi préciser [0] ici ? 🤔 => c'était attendu: un élément par établissment. l'agregation cross-computation va exploser ces etablissements au niveau le plus haut des données en sorties.
  })
  t.log(JSON.stringify(finalizeResult, null, 2))
  t.deepEqual(finalizeResult, expectedFinalizeResults as any) // ⚠️ Les types sont incompatibles => réparer la déclaration TS de finalize ?
})

// il manque une agregation qui sort dans la base Features_debug (_debug car on a spécifié la clé)
// => changer la valeur attendue: récupérer la valeur intermédiaire de l'appel à /reduce au lieu de la sortie finale
// (cette agrégation ne s'appuie pas sur des scripts JS, cf `reduceFinalAggregation()`)
// (a.k.a. cross-computation) l'objectif était de reduce type par type, mais c'est pas fini.
// pierre est serein sur la sortie actuelle de finalize().
