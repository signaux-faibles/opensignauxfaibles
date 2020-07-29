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

const testedProps = [
  {
    assertion:
      "La variable cotisation_moy12m de chaque mois est égale au montant de cotisation, quand celui-ci est stable sur tout la période",
    name: "cotisation_moy12m",
    input: forEachMonth(() => ({ cotisation: 10 })),
    expected: new Array(periodeSerie.length).fill(10),
  },
  {
    assertion:
      "La variable cotisation_moy12m décroit quand une cotisation est passée",
    name: "cotisation_moy12m",
    input: forEachMonth(({ month }) => ({ cotisation: month === 0 ? 10 : 0 })),
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
    name: "cotisation_moy12m",
    input: forEachMonth(({ month }) => ({
      cotisation: month === 12 ? 10 : 0,
    })),
    expected: new Array(12).fill(0).concat([10 / 12]),
  },
]

test("cotisation retourne les mêmes périodes que fournies en entrée", (t) => {
  const input = forEachMonth(() => ({}))
  const actual = cotisation(input)
  t.deepEqual(Object.keys(actual), Object.keys(input))
})

testedProps.forEach((testedProp) => {
  test(testedProp.assertion, (t) => {
    const actual = cotisation(testedProp.input)
    testedProp.expected.forEach((expectedPropValue, indiceMois) => {
      const actualValue = actual[dateAddMonth(dateDebut, indiceMois).getTime()]
      t.is(
        actualValue[testedProp.name as keyof SortieCotisation],
        expectedPropValue,
        `mois: #${indiceMois}, expected: ${expectedPropValue}`
      )
    })
  })
})
