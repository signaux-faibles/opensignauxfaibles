// Ces tests visent à couvrir toutes les fonctions invoquées par map(), en lui
// fournissant un jeu de données minimal mais incluant tous les types d'entrée
// inclus dans RawData.
//
// => Penser à ajouter les nouveaux types dans rawEtabData et rawEntrData.
//
import test, { ExecutionContext } from "ava"
import { map } from "../public/map"
import { reduce } from "../public/reduce"
import { finalize } from "../public/finalize"
import { runMongoMap } from "../test/helpers/mongodb"
import { setGlobals } from "../test/helpers/setGlobals"
import {
  EntrepriseBatchProps,
  EntrepriseDataValues,
  EtablissementBatchProps,
  EtablissementDataValues,
  Siret,
} from "../RawDataTypes"
import { EntréeApConso, EntréeApDemande } from "../GeneratedTypes"
import { apdemande } from "./apdemande"

// test data inspired by test.sh
const siret: Siret = "01234567891011"
const batchKey = "1910"
const dates = [
  new Date("2015-12-01T00:00:00.000+0000"),
  new Date("2016-01-01T00:00:00.000+0000"),
]
setGlobals({
  actual_batch: batchKey, // used by map()
  serie_periode: dates, // used by effectifs(), which is called by map()
})

const rawEtabData: EtablissementBatchProps = {
  reporder: {},
  apconso: {
    somehash: {
      id_conso: "",
      periode: new Date(),
      heure_consomme: 1,
    } as EntréeApConso,
  },
}

const rawData: EtablissementDataValues = {
  batch: {
    [batchKey]: rawEtabData,
  },
  scope: "etablissement",
  key: siret,
}

const etablissementKey = rawData.scope + "_" + siret

const expectedMapResult = {
  _id: etablissementKey,
  value: {
    apconso: [rawEtabData.apconso.somehash],
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

const expectedReduceResults = expectedMapResult.value

const expectedFinalizeResultValue = expectedMapResult.value

test(`apdemande() classe les entrées par ordre antichronologique`, (t) => {
  const entrée = {
    hash1: { periode: { start: new Date(0) } } as EntréeApDemande,
    hash2: { periode: { start: new Date(1) } } as EntréeApDemande,
  }
  const sortie = apdemande(entrée)
  t.deepEqual(sortie, [entrée.hash2, entrée.hash1])
})

// exécution complète de la chaine "public"

test.serial(
  `public.map() retourne toutes les propriétés d'entreprise attendues sur le frontal`,
  (t: ExecutionContext) => {
    const rawEntrData: EntrepriseBatchProps = {
      reporder: {},
      paydex: { somehash: { date_valeur: new Date(), nb_jours: 1 } },
    }
    const rawData: EntrepriseDataValues = {
      scope: "entreprise",
      key: siret.substr(0, 9), // siren
      batch: { [batchKey]: rawEntrData },
    }
    const expectedMapResult = {
      _id: rawData.scope + "_" + rawData.key,
      value: {
        key: rawData.key,
        batch: batchKey,
        paydex: [rawEntrData.paydex.somehash],
      },
    }
    const mapResults = runMongoMap(map, [{ _id: null, value: rawData }])
    t.deepEqual(mapResults, [expectedMapResult])
  }
)

test.serial(
  `public.map() retourne les propriétés d'établissement présentées sur le frontal`,
  (t: ExecutionContext) => {
    const mapResults = runMongoMap(map, [{ _id: null, value: rawData }])
    t.deepEqual(mapResults, [expectedMapResult])
  }
)

test.serial(
  `public.reduce() retourne les propriétés d'établissement, telles quelles`,
  (t: ExecutionContext) => {
    const reduceValues = [expectedMapResult.value]
    const reduceResults = reduce({ scope: rawData.scope }, reduceValues)
    t.deepEqual(reduceResults, expectedReduceResults)
  }
)

test.serial(
  `public.finalize() retourne les propriétés d'établissement, telles quelles`,
  (t: ExecutionContext) => {
    const finalizeResultValue = finalize(
      { scope: rawData.scope },
      expectedReduceResults
    )
    t.deepEqual(finalizeResultValue, expectedFinalizeResultValue)
  }
)
