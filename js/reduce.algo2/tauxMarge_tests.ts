import test from "ava"
import { tauxMarge } from "./tauxMarge"
import { EntréeDiane } from "../GeneratedTypes"

test(`tauxMarge est calculé en divisant l'excedent brut d'exploitation par la valeur ajoutée`, (t) => {
  const entréeDiane: EntréeDiane = {
    excedent_brut_d_exploitation: 1528,
    valeur_ajoutee: 3419,
  }
  t.is(tauxMarge(entréeDiane), (100 * 1528) / 3419)
})

test(`tauxMarge est null si l'excedent brut d'exploitation n'est pas défini`, (t) => {
  const entréeDiane: EntréeDiane = {
    valeur_ajoutee: 3419,
  }
  t.is(tauxMarge(entréeDiane), null)
})

test(`tauxMarge est null si la valeur ajoutée n'est pas définie`, (t) => {
  const entréeDiane: EntréeDiane = {
    excedent_brut_d_exploitation: 1528,
  }
  t.is(tauxMarge(entréeDiane), null)
})

// $ cd js && $(npm bin)/ava ./reduce.algo2/tauxMarge_tests.ts
