import "../globals"
import test from "ava"
import { cotisation, Input, SortieCotisation } from "./cotisation"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

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
      "La variable cotisation_moy12m de chaque mois est égale au montant de cotisation, quand celui-ci est stable sur tout la période",
    input: forEachMonth(() => ({ cotisation: 10 })),
    propName: "cotisation_moy12m",
    expected: new Array(periodeSerie.length).fill(10),
  },
  {
    assertion:
      "La variable cotisation_moy12m décroit quand une cotisation est passée",
    input: forEachMonth(({ month }) => ({ cotisation: month === 0 ? 10 : 0 })),
    propName: "cotisation_moy12m",
    expected: [
      10 / 1,
      10 / 2,
      10 / 3,
      10 / 4,
      10 / 5,
      10 / 6,
      10 / 7,
      10 / 8,
      10 / 9,
      10 / 10,
      10 / 11,
      10 / 12,
      0,
    ],
  },
  {
    assertion:
      "La variable cotisation_moy12m est nulle jusqu'à ce qu'une cotisation soit présente",
    input: forEachMonth(({ month }) => ({
      cotisation: month === 12 ? 10 : 0,
    })),
    propName: "cotisation_moy12m",
    expected: [...new Array(12).fill(0), 10 / 12],
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
      montant_part_ouvriere: [5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5][month],
      montant_part_patronale: [5, 5, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5][month],
    })),
    propName: "ratio_dette",
    expected: [1, 1, 1, 1, 1, 1 / 2, 1 / 2, 1 / 2, 1 / 2, 1 / 2, 1, 1, 1],
  },
  {
    assertion:
      "La variable ratio_dette n'est pas calculée s'il manque un montant de cotisation au sein de la période",
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
    expected: [0, 1, ...new Array(10).fill(2), 1],
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
