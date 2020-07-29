import "../globals"
import test from "ava"
import { cotisation, SortieCotisation } from "./cotisation"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

const dureeEnMois = 13
const dateDebut = new Date("2018-01-01")
const dateFin = dateAddMonth(dateDebut, dureeEnMois)
const periode = generatePeriodSerie(dateDebut, dateFin)

const testedProps = [
  {
    assertion:
      "La variable cotisation_moy12m de chaque mois est égale au montant de cotisation, quand celui-ci est stable sur tout la période",
    name: "cotisation_moy12m",
    input: periode.reduce(
      (acc, periode) => ({
        ...acc,
        [periode.getTime()]: {
          periode,
          cotisation: 10,
        },
      }),
      {}
    ),
    expected: new Array(periode.length).fill(10),
  },
]

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
