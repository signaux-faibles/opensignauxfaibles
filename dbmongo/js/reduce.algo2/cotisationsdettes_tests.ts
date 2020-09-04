import test from "ava"
import { cotisationsdettes, SortieCotisationsDettes } from "./cotisationsdettes"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "../common/dateAddMonth"

const dureeEnMois = 13
const moisRemboursement = 4
const dateDebut = new Date("2018-01-01")
const dateFin = dateAddMonth(dateDebut, dureeEnMois)
const periode = generatePeriodSerie(dateDebut, dateFin).map((date) =>
  date.getTime()
)

const montantCotisation = 100
const montantPartOuvrière = 100
const montantPartPatronale = 200

const expectedDette = (montant: number, moisRemboursement: number): number[] =>
  Array(moisRemboursement)
    .fill(montant)
    .concat(Array(dureeEnMois - moisRemboursement).fill(0))

const setupCompanyValuesWithCotisation = (
  dateDebut: Date,
  montantCotisation: number
) => ({
  cotisation: {
    hash1: {
      periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
      du: montantCotisation,
    },
  },
  debit: {},
})

const setupCompanyValuesForMontant = (dateDebut: Date) => ({
  cotisation: {},
  debit: {
    hash1: {
      periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
      part_ouvriere: montantPartOuvrière,
      part_patronale: montantPartPatronale,
      date_traitement: dateDebut,
      debit_suivant: "",
      numero_compte: "",
      numero_ecart_negatif: 1,
      numero_historique: 2,
    },
    hash2: {
      periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
      part_ouvriere: 0,
      part_patronale: 0,
      date_traitement: dateAddMonth(dateDebut, moisRemboursement),
      debit_suivant: "",
      numero_compte: "",
      numero_ecart_negatif: 1,
      numero_historique: 3,
    },
  },
})

const setupCompanyValues = (dateDebut: Date) => ({
  cotisation: {
    hash1: {
      periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
      du: 60,
    },
  },
  debit: {
    // dette initiale
    hashDetteInitiale: {
      periode: {
        start: dateDebut,
        end: dateAddMonth(dateDebut, 1),
      },
      numero_ecart_negatif: 1,
      numero_historique: 2,
      numero_compte: "",
      date_traitement: dateDebut,
      debit_suivant: "hashRemboursement",
      part_ouvriere: 60,
      part_patronale: 0,
    },
    // remboursement la dette
    hashRemboursement: {
      periode: {
        start: dateDebut,
        end: dateAddMonth(dateDebut, 1),
      },
      numero_ecart_negatif: 1, // même valeur que pour le débit précédent
      numero_historique: 3, // incrémentation depuis le débit précédent
      numero_compte: "",
      date_traitement: dateAddMonth(dateDebut, 1),
      debit_suivant: "",
      part_ouvriere: 0,
      part_patronale: 0,
    },
  },
})

// TODO: refactoriser les générateurs de données de tests pour éviter la
// redondance et faciliter les tests futurs.

const generatePastTestCase = (
  ouvrièreOuPatronale: "ouvriere" | "patronale",
  montantDette: number,
  décalageEnMois: number
) => ({
  assertion: `Le montant de part ${ouvrièreOuPatronale} d'une période est reporté dans montant_part_${ouvrièreOuPatronale}_past_${décalageEnMois}`,
  name: `montant_part_${ouvrièreOuPatronale}_past_${décalageEnMois}`,
  input: setupCompanyValuesForMontant(dateDebut),
  expected: [
    ...Array(décalageEnMois).fill(undefined),
    ...expectedDette(montantDette, moisRemboursement),
  ].slice(0, dureeEnMois),
})

const testedProps = [
  {
    assertion:
      "La variable cotisation représente les cotisations sociales dues à une période donnée",
    name: "cotisation",
    input: setupCompanyValuesWithCotisation(dateDebut, montantCotisation),
    expected: [montantCotisation],
  },
  {
    assertion:
      "montant_part_patronale est annulé après le remboursement de la dette",
    name: "montant_part_patronale",
    input: setupCompanyValuesForMontant(dateDebut),
    expected: expectedDette(montantPartPatronale, moisRemboursement),
  },
  {
    assertion:
      "montant_part_ouvriere est annulé après le remboursement de la dette",
    name: "montant_part_ouvriere",
    input: setupCompanyValuesForMontant(dateDebut),
    expected: expectedDette(montantPartOuvrière, moisRemboursement),
  },
  ...[1, 2, 3, 6, 12].map((décalageEnMois) =>
    generatePastTestCase("ouvriere", montantPartOuvrière, décalageEnMois)
  ),
  ...[1, 2, 3, 6, 12].map((décalageEnMois) =>
    generatePastTestCase("patronale", montantPartPatronale, décalageEnMois)
  ),
  {
    assertion:
      "interessante_urssaf est vrai quand l'entreprise n'a pas eu de débit (dette) sur les 6 derniers mois",
    name: "interessante_urssaf",
    input: setupCompanyValues(dateDebut),
    expected: Array(6).fill(false).concat(Array(2).fill(undefined)),
  },
]

testedProps.forEach((testedProp) => {
  test(testedProp.assertion, (t) => {
    const { cotisation, debit } = testedProp.input
    const actual = cotisationsdettes(cotisation, debit, periode, dateFin)

    testedProp.expected.forEach((expectedPropValue, indiceMois) => {
      const actualValue = actual[dateAddMonth(dateDebut, indiceMois).getTime()]
      t.is(
        actualValue[testedProp.name as keyof SortieCotisationsDettes],
        expectedPropValue,
        `mois: #${indiceMois}, expected: ${expectedPropValue}`
      )
    })
  })
})
