import test from "ava"
import { finalize } from "./finalize"
import { CléSortieMap, SortieMap } from "./map"
import { setGlobals } from "../test/helpers/setGlobals"

const clé = {
  batch: "dummy",
  siren: "012345678",
  periode: new Date("2014-01-01"),
  type: "other" as CléSortieMap["type"],
}

test(`finalize() fait la somme des effectifs des établissement rattachés à l'entreprise`, (t) => {
  const results = finalize(clé, {
    siret1: { effectif: 3 },
    siret2: { effectif: 4 },
  })
  t.deepEqual(results, [
    {
      effectif: 3,
      effectif_entreprise: 7,
      nbr_etablissements_connus: 2,
    },
    {
      effectif: 4,
      effectif_entreprise: 7,
      nbr_etablissements_connus: 2,
    },
  ] as unknown)
})

test(`finalize() fait la somme des heures d'activité partielle des établissement rattachés à l'entreprise`, (t) => {
  const results = finalize(clé, {
    siret1: { apart_heures_consommees: 3 },
    siret2: { apart_heures_consommees: 4 },
  })
  t.deepEqual(results, [
    {
      apart_heures_consommees: 3,
      apart_entreprise: 7,
      nbr_etablissements_connus: 2,
    },
    {
      apart_heures_consommees: 4,
      apart_entreprise: 7,
      nbr_etablissements_connus: 2,
    },
  ] as unknown)
})

test(`finalize() calcule la dette totale de l'entreprise à partir de celle des établissement`, (t) => {
  const results = finalize(clé, {
    siret1: { montant_part_patronale: 3 },
    siret2: { montant_part_ouvriere: 4 },
  })
  t.deepEqual(results, [
    {
      montant_part_patronale: 3,
      debit_entreprise: 7,
      nbr_etablissements_connus: 2,
    },
    {
      montant_part_ouvriere: 4,
      debit_entreprise: 7,
      nbr_etablissements_connus: 2,
    },
  ] as unknown)
})

test("finalize() retourne un tableau vide au dela de 1500 établissements pour une même entreprise", (t) => {
  const clé = {
    batch: "dummy",
    siren: "012345678",
    periode: new Date("2014-01-01"),
    type: "other" as CléSortieMap["type"],
  }
  const etablissements: SortieMap = {}
  for (let i = 0; i <= 1500; ++i) {
    // 1500 = cf maxEtabParEntr de finalize.ts
    const siret = `${i}`
    etablissements[siret] = { siret }
  }
  const results = finalize(clé, etablissements)
  t.deepEqual(results, [])
})

test("finalize() retourne un objet incomplet en cas de dépassement de taille autorisée", (t) => {
  setGlobals({ print: () => {} }) // eslint-disable-line @typescript-eslint/no-empty-function
  const clé = {
    batch: "x".repeat(16777216), // cf maxBsonSize de finalize.ts
    siren: "012345678",
    periode: new Date("2014-01-01"),
    type: "other" as CléSortieMap["type"],
  }
  const results = finalize(clé, {
    "12345678901234": { siret: "12345678901234" },
  })
  t.deepEqual(results, { incomplete: true })
})
