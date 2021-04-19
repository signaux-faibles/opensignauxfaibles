import test from "ava"
import { finalize } from "./finalize"
import { CléSortieMap } from "./map"

const siret1 = "12345678901234"
const siret2 = "12345678901235"

const clé = {
  batch: "dummy",
  siren: "012345678",
  periode: new Date("2014-01-01"),
  type: "other" as CléSortieMap["type"],
}

test(`finalize() fait la somme des effectifs des établissement rattachés à l'entreprise`, (t) => {
  const results = finalize(clé, {
    siret1: { siret: siret1, effectif: 3 },
    siret2: { siret: siret2, effectif: 4 },
  })
  t.deepEqual(results, [
    {
      effectif: 3,
      effectif_entreprise: 7,
      nbr_etablissements_connus: 2,
      siret: siret1,
    },
    {
      effectif: 4,
      effectif_entreprise: 7,
      nbr_etablissements_connus: 2,
      siret: siret2,
    },
  ] as unknown)
})

test(`finalize() fait la somme des heures d'activité partielle des établissement rattachés à l'entreprise`, (t) => {
  const results = finalize(clé, {
    siret1: { siret: siret1, apart_heures_consommees: 3 },
    siret2: { siret: siret2, apart_heures_consommees: 4 },
  })
  t.deepEqual(results, [
    {
      apart_heures_consommees: 3,
      apart_entreprise: 7,
      nbr_etablissements_connus: 2,
      siret: siret1,
    },
    {
      apart_heures_consommees: 4,
      apart_entreprise: 7,
      nbr_etablissements_connus: 2,
      siret: siret2,
    },
  ] as unknown)
})

test(`finalize() calcule la dette totale de l'entreprise à partir de celle des établissement`, (t) => {
  const results = finalize(clé, {
    siret1: { siret: siret1, montant_part_patronale: 3 },
    siret2: { siret: siret2, montant_part_ouvriere: 4 },
  })
  t.deepEqual(results, [
    {
      montant_part_patronale: 3,
      debit_entreprise: 7,
      nbr_etablissements_connus: 2,
      siret: siret1,
    },
    {
      montant_part_ouvriere: 4,
      debit_entreprise: 7,
      nbr_etablissements_connus: 2,
      siret: siret2,
    },
  ] as unknown)
})
