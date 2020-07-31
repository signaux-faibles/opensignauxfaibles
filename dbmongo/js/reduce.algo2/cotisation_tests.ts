import "../globals"
import test from "ava"
import { cotisation, Input, SortieCotisation } from "./cotisation"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "../common/dateAddMonth"

const dureeEnMois = 13
const dateDebut = new Date("2018-01-01")
const dateFin = dateAddMonth(dateDebut, dureeEnMois)
const periodeSerie = generatePeriodSerie(dateDebut, dateFin) // => length === dureeEnMois

const forEachMonth = (
  fct: ({ periode, month }: { periode: Date; month: number }) => Partial<Input>
) =>
  periodeSerie.reduce(
    (acc, periode, month) => ({
      ...acc,
      [periode.getTime()]: { periode, ...fct({ periode, month }) },
    }),
    {}
  )

const testCases = [
  {
    assertion:
      "La variable cotisation_moy12m est calculée sur la base de 12 mois de données, pas moins",
    input: forEachMonth(() => ({ cotisation: 10 })),
    propName: "cotisation_moy12m",
    expected: [...new Array(11).fill(undefined), 10, 10],
  },
  {
    assertion:
      "La variable cotisation_moy12m est nulle jusqu'à la présence d'une cotisation non nulle",
    input: forEachMonth(({ month }) => ({
      cotisation: month === 12 ? 10 : 0,
    })),
    propName: "cotisation_moy12m",
    expected: [...new Array(11).fill(undefined), 0, 10 / 12],
  },
  {
    assertion:
      "La variable cotisation_moy12m n'est pas calculée s'il manque un montant de cotisation au sein de la période",
    input: forEachMonth(({ month }) => ({
      cotisation: month === 0 ? undefined : 10,
    })),
    propName: "cotisation_moy12m",
    expected: [...new Array(12).fill(undefined), 10],
  },
  {
    assertion:
      "La variable ratio_dette divise montant_part_ouvriere et montant_part_patronale par cotisation_moy12m",
    input: forEachMonth(({ month }) => ({
      cotisation: [10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10][month],
      montant_part_ouvriere: 5,
      montant_part_patronale: [...new Array(12).fill(0), 5][month],
    })),
    propName: "ratio_dette",
    expected: [...new Array(11).fill(undefined), 1 / 2, 1],
  },
  {
    assertion:
      "La variable ratio_dette n'est pas calculée si on n'a pas 12 mois d'historique de cotisations",
    input: forEachMonth(({ month }) => ({
      cotisation: month === 0 ? undefined : 10,
      montant_part_ouvriere: 5,
      montant_part_patronale: 5,
    })),
    propName: "ratio_dette",
    expected: [...new Array(12).fill(undefined), (5 + 5) / 10],
  },
  {
    assertion:
      "La variable ratio_dette considère tout montant de part ouvrière ou patronale manquant comme nul",
    input: forEachMonth(({ month }) => ({
      cotisation: 10,
      montant_part_ouvriere: [undefined, ...new Array(11).fill(10), undefined][
        month
      ],
      montant_part_patronale: [undefined, undefined, ...new Array(11).fill(10)][
        month
      ],
    })),
    propName: "ratio_dette",
    expected: [...new Array(11).fill(undefined), 2, 1],
  },
  {
    assertion:
      "La variable ratio_dette_moy12m considère tout montant de part ouvrière ou patronale manquant comme nul",
    input: forEachMonth(({ month }) => ({
      cotisation: 10,
      montant_part_ouvriere: undefined,
      montant_part_patronale: month < 12 ? 10 : undefined,
    })),
    propName: "ratio_dette_moy12m",
    expected: [...new Array(11).fill(undefined), 1, 0.9166666666666666],
  },
]

testCases.forEach(({ assertion, input, propName, expected }) => {
  test(assertion, (t) => {
    const actual = cotisation(input)
    expected.forEach((expectedPropValue, indiceMois) => {
      const actualValue = actual[dateAddMonth(dateDebut, indiceMois).getTime()]
      t.is(
        actualValue[propName as keyof SortieCotisation],
        expectedPropValue,
        `mois: #${indiceMois}, expected: ${expectedPropValue}`
      )
    })
  })
})

test("cotisation retourne les mêmes périodes que fournies en entrée", (t) => {
  const input = forEachMonth(() => ({}))
  const actual = cotisation(input)
  t.deepEqual(Object.keys(actual), Object.keys(input))
})
