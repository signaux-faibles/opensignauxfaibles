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
    expected: new Array(12).fill(0).concat([10 / 12]),
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
