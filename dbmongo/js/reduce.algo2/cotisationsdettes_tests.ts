import "../globals"
import test, { ExecutionContext } from "ava"
import { cotisationsdettes } from "./cotisationsdettes"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

function décaler(tableau: number[], décalage: number): number[] {
  return Array(décalage).fill(undefined).concat(tableau).slice(0, -décalage)
}

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

test.only("Le montant de dette d'une période est reporté dans les périodes suivantes", (t: ExecutionContext) => {
  const dureeEnMois = 13
  const dateDebut = new Date("2018-01-01")
  const dateFin = dateAddMonth(dateDebut, dureeEnMois)
  const periode = generatePeriodSerie(dateDebut, dateFin).map((date) =>
    date.getTime()
  )

  ;(globalThis as any).date_fin = dateFin // TODO: transformer ce parametre en parametre local de fonction

  const moisRemboursement = 4
  const v: DonnéesCotisation & DonnéesDebit = {
    cotisation: {},
    debit: {
      hash1: {
        periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
        part_ouvriere: 100,
        part_patronale: 200,
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

  const actual = cotisationsdettes(v, periode)

  const expectedMontantPartOuvrière = Array(moisRemboursement).fill(100).concat(
    Array(dureeEnMois - moisRemboursement).fill(0))

  const expectedMontantPartPatronale = Array(moisRemboursement).fill(200).concat(
    Array(dureeEnMois - moisRemboursement).fill(0))

  const expectedMontantPartOuvrièrePast1 = décaler(
    expectedMontantPartOuvrière,
    1
  )
  const expectedMontantPartOuvrièrePast2 = décaler(
    expectedMontantPartOuvrière,
    2
  )
  const expectedMontantPartOuvrièrePast3 = décaler(
    expectedMontantPartOuvrière,
    3
  )
  const expectedMontantPartOuvrièrePast6 = décaler(
    expectedMontantPartOuvrière,
    6
  )
  const expectedMontantPartOuvrièrePast12 = décaler(
    expectedMontantPartOuvrière,
    12
  )
  const expectedMontantPartPatronalePast1 = décaler(
    expectedMontantPartPatronale,
    1
  )
  const expectedMontantPartPatronalePast2 = décaler(
    expectedMontantPartPatronale,
    2
  )
  const expectedMontantPartPatronalePast3 = décaler(
    expectedMontantPartPatronale,
    3
  )
  const expectedMontantPartPatronalePast6 = décaler(
    expectedMontantPartPatronale,
    6
  )
  const expectedMontantPartPatronalePast12 = décaler(
    expectedMontantPartPatronale,
    12
  )

  const expectedInteressanteUrssaf = Array(9)
    .fill(false)
    .concat(Array(4).fill(undefined))

  for (let période = 0; période < 13; ++période) {
    t.log({ période }, actual[dateAddMonth(dateDebut, période).getTime()])
    const expected = {
      interessante_urssaf: expectedInteressanteUrssaf[période],
      montant_part_ouvriere: expectedMontantPartOuvrière[période],
      montant_part_patronale: expectedMontantPartPatronale[période],
      montant_part_ouvriere_past_1: expectedMontantPartOuvrièrePast1[période],
      montant_part_patronale_past_1: expectedMontantPartPatronalePast1[période],
      montant_part_ouvriere_past_2: expectedMontantPartOuvrièrePast2[période],
      montant_part_patronale_past_2: expectedMontantPartPatronalePast2[période],
      montant_part_ouvriere_past_3: expectedMontantPartOuvrièrePast3[période],
      montant_part_patronale_past_3: expectedMontantPartPatronalePast3[période],
      montant_part_ouvriere_past_6: expectedMontantPartOuvrièrePast6[période],
      montant_part_patronale_past_6: expectedMontantPartPatronalePast6[période],
      montant_part_ouvriere_past_12: expectedMontantPartOuvrièrePast12[période],
      montant_part_patronale_past_12:
        expectedMontantPartPatronalePast12[période],
    }
    Object.keys(expected).forEach((p) => {
      const prop = p as keyof typeof expected
      if (typeof expected[prop] === "undefined") {
        delete expected[prop]
      }
    })
    t.deepEqual(actual[dateAddMonth(dateDebut, période).getTime()], expected)
  }
})

test("interessante_urssaf est vrai quand l'entreprise n'a pas eu de débit (dette) sur les 6 derniers mois", (t: ExecutionContext) => {
  const dateDebut = new Date("2018-01-01")
  const periode = generatePeriodSerie(
    dateDebut,
    dateAddMonth(dateDebut, 8)
  ).map((date) => date.getTime())

  ;(globalThis as any).date_fin = dateAddMonth(dateDebut, 8) // utilisé par cotisationsdettes lors du traitement des débits

  const v = {
    cotisation: {
      hash1: {
        periode: { start: dateDebut, end: dateAddMonth(dateDebut, 1) },
        du: 60,
      },
    },
    debit: {
      // tentative de répartition du montant de la dette (part ouvrière: 100%)
      [dateDebut.getTime()]: {
        periode: {
          start: dateDebut,
          end: dateAddMonth(dateDebut, 1),
        },
        numero_ecart_negatif: 1,
        numero_historique: 2,
        numero_compte: "3",
        date_traitement: dateDebut,
        debit_suivant: dateAddMonth(dateDebut, 1).toString(),
        part_ouvriere: 60,
        part_patronale: 0,
      },
      // tentative de remboursement la dette
      [dateAddMonth(dateDebut, 1).getTime()]: {
        periode: {
          start: dateAddMonth(dateDebut, 1),
          end: dateAddMonth(dateDebut, 2),
        },
        numero_ecart_negatif: 1,
        numero_historique: 2,
        numero_compte: "3",
        date_traitement: dateAddMonth(dateDebut, 1),
        debit_suivant: "",
        part_ouvriere: 0,
        part_patronale: 0,
      },
    },
    /*
    dettes: {
      hash1: {
        periode: dateDebut,
        part_ouvriere: 30,
        part_patronale: 0,
      },
    },
    */
  }

  const actual = cotisationsdettes(v, periode)

  t.log(actual)

  t.true(actual[dateAddMonth(dateDebut, 7).getTime()].interessante_urssaf)

  for (const month of [0, 1, 2, 3, 4, 5, 6]) {
    t.false(
      actual[dateAddMonth(dateDebut, month).getTime()].interessante_urssaf
    )
  }
})
