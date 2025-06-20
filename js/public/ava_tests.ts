// Ces tests visent à couvrir toutes les fonctions invoquées par map(), en lui
// fournissant un jeu de données minimal mais incluant tous les types d'entrée
// inclus dans RawData.
//
// => Penser à ajouter les nouveaux types dans rawEtabData et rawEntrData.
//
import test from "ava"
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
import { EntréeApDemande } from "../GeneratedTypes"
import { apdemande } from "./apdemande"

// test data inspired by test.sh
const siret: Siret = "01234567891011"
const batchKey = "1910"
const dateDébut = new Date("2015-12-01T00:00:00.000+0000")
const dateFin = new Date("2016-01-01T00:00:00.000+0000")
const dates = [dateDébut, dateFin]
setGlobals({
  actual_batch: batchKey, // used by map()
  serie_periode: dates, // used by effectifs(), which is called by map()
})

const rawEtabData: EtablissementBatchProps = {
  reporder: {},
  apdemande: {
    b88032c02d1724d92279a47599b112dd: {
      id_demande: "S044130482",
      effectif_entreprise: 10,
      effectif: 10,
      date_statut: dateDébut,
      periode: {
        start: dateDébut,
        end: dateFin,
      },
      hta: 1500,
      mta: 10000,
      effectif_autorise: 30,
      motif_recours_se: 2,
      heure_consommee: 0,
      montant_consommee: 0,
      effectif_consomme: 0,
      perimetre: 5,
    },
  },
  apconso: {
    e3af9bbf7d37e62fbbf88efdae464746: {
      id_conso: "S044130482",
      heure_consomme: 326.0,
      montant: 2530.13,
      effectif: 20,
      periode: dateDébut,
    },
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
    apconso: [rawEtabData.apconso.e3af9bbf7d37e62fbbf88efdae464746],
    apdemande: [rawEtabData.apdemande.b88032c02d1724d92279a47599b112dd],
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

test(`apdemande() classe les entrées par ordre antichronologique`, (t: test) => {
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
  (t: test) => {
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
  (t: test) => {
    const mapResults = runMongoMap(map, [{ _id: null, value: rawData }])
    t.deepEqual(mapResults, [expectedMapResult])
  }
)

test.serial(
  `public.reduce() retourne les propriétés d'établissement, telles quelles`,
  (t: test) => {
    const reduceValues = [expectedMapResult.value]
    const reduceResults = reduce({ scope: rawData.scope }, reduceValues)
    t.deepEqual(reduceResults, expectedReduceResults)
  }
)

test.serial(
  `public.finalize() retourne les propriétés d'établissement, telles quelles`,
  (t: test) => {
    const finalizeResultValue = finalize(
      { scope: rawData.scope },
      expectedReduceResults
    )
    t.deepEqual(finalizeResultValue, expectedFinalizeResultValue)
  }
)
