import "../globals"
import test, { ExecutionContext } from "ava"
import { cotisationsdettes, SortieCotisationsDettes } from "./cotisationsdettes"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

// Supprime les propriétés de obj dont la valeur est indéfinie.
const deleteUndefinedProps = <T>(obj: T): void =>
  (Object.keys(obj) as Array<keyof T>).forEach((prop) =>
    typeof obj[prop] === "undefined" ? delete obj[prop] : {}
  )

test("La variable cotisation représente les cotisations sociales dues à une période donnée", (t: ExecutionContext) => {
  const date = new Date("2018-01-01")
  const datePlusUnMois = new Date("2018-02-01")

  const v: DonnéesCotisation & DonnéesDebit = {
    cotisation: {
      hash1: {
        periode: { start: date, end: datePlusUnMois },
        du: 100,
      },
    },
    debit: {},
  }

  const actual = cotisationsdettes(v, [date.getTime()])

  const expected = {
    [date.getTime()]: {
      montant_part_ouvriere: 0,
      montant_part_patronale: 0,
      cotisation: 100,
    },
  }

  t.deepEqual(actual, expected)
})

test("Le montant de dette d'une période est reporté dans les périodes suivantes", (t: ExecutionContext) => {
  const dureeEnMois = 13
  const dateDebut = new Date("2018-01-01")
  const dateFin = dateAddMonth(dateDebut, dureeEnMois)
  const periode = generatePeriodSerie(dateDebut, dateFin).map((date) =>
    date.getTime()
  )

  const moisRemboursement = 4
  const partOuvrière = 100
  const partPatronale = 200
  const v: DonnéesCotisation & DonnéesDebit = {
    cotisation: {},
    debit: {
      hash1: {
        periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
        part_ouvriere: partOuvrière,
        part_patronale: partPatronale,
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
  }

  const output = cotisationsdettes(v, periode, dateFin)

  const expPartOuvrière = Array(moisRemboursement)
    .fill(partOuvrière)
    .concat(Array(dureeEnMois - moisRemboursement).fill(0))

  const expPartPatronale = Array(moisRemboursement)
    .fill(partPatronale)
    .concat(Array(dureeEnMois - moisRemboursement).fill(0))

  for (let mois = 0; mois < 13; ++mois) {
    const expected = {
      montant_part_ouvriere: expPartOuvrière[mois],
      montant_part_patronale: expPartPatronale[mois],
      montant_part_ouvriere_past_1: expPartOuvrière[mois - 1],
      montant_part_patronale_past_1: expPartPatronale[mois - 1],
      montant_part_ouvriere_past_2: expPartOuvrière[mois - 2],
      montant_part_patronale_past_2: expPartPatronale[mois - 2],
      montant_part_ouvriere_past_3: expPartOuvrière[mois - 3],
      montant_part_patronale_past_3: expPartPatronale[mois - 3],
      montant_part_ouvriere_past_6: expPartOuvrière[mois - 6],
      montant_part_patronale_past_6: expPartPatronale[mois - 6],
      montant_part_ouvriere_past_12: expPartOuvrière[mois - 12],
      montant_part_patronale_past_12: expPartPatronale[mois - 12],
    }
    deleteUndefinedProps(expected)
    const actual = output[dateAddMonth(dateDebut, mois).getTime()]
    delete actual.interessante_urssaf // exclure interessante_urssaf car cette prop est considérée par un autre test
    t.deepEqual(actual, expected)
  }
})

const setupPeriodes = () => {
  const dureeEnMois = 13
  const dateDebut = new Date("2018-01-01")
  const dateFin = dateAddMonth(dateDebut, dureeEnMois)
  const periode = generatePeriodSerie(dateDebut, dateFin).map((date) =>
    date.getTime()
  )
  return { dateDebut, dateFin, periode }
}
const setupCompanyValuesForMontant = (dateDebut: Date) => ({
  cotisation: {},
  debit: {
    hash1: {
      periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
      part_ouvriere: partOuvrière,
      part_patronale: partPatronale,
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

const dureeEnMois = 13
const moisRemboursement = 4
const partOuvrière = 100
const partPatronale = 200

const expPartOuvrière = (
  partOuvrière: number,
  moisRemboursement: number,
  dureeEnMois: number
) =>
  Array(moisRemboursement)
    .fill(partOuvrière)
    .concat(Array(dureeEnMois - moisRemboursement).fill(0))

const { dateDebut, dateFin, periode } = setupPeriodes()

const generatePastTestCase = (
  ouvrièreOuPatronale: "ouvriere" | "patronale",
  décalageEnMois: number
) => ({
  assertion: `Le montant de part ${ouvrièreOuPatronale} d'une période est reporté dans montant_part_${ouvrièreOuPatronale}_past_${décalageEnMois}`,
  name: `montant_part_${ouvrièreOuPatronale}_past_${décalageEnMois}`,
  input: setupCompanyValuesForMontant(dateDebut),
  expected: [
    undefined,
    ...expPartOuvrière(
      partOuvrière,
      moisRemboursement,
      dureeEnMois - décalageEnMois
    ),
  ],
})

const generatePastTestCases = (
  ouvrièreOuPatronale: "ouvriere" | "patronale",
  décalagesEnMois: number[]
) =>
  décalagesEnMois.map((décalageEnMois) =>
    generatePastTestCase(ouvrièreOuPatronale, décalageEnMois)
  )

const testedProps = [
  ...generatePastTestCases("ouvriere", [1, 2]),
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
    const actual = cotisationsdettes(testedProp.input, periode, dateFin)

    testedProp.expected.forEach((expectedPropValue, indiceMois) => {
      const actualValue = actual[dateAddMonth(dateDebut, indiceMois).getTime()]
      t.is(
        actualValue[testedProp.name as keyof SortieCotisationsDettes],
        expectedPropValue,
        `mois: #${indiceMois}`
      )
    })
  })
})
