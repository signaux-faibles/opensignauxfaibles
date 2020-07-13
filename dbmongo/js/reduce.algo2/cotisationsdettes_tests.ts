import "../globals"
import test, { ExecutionContext } from "ava"
import { cotisationsdettes } from "./cotisationsdettes"
import { generatePeriodSerie } from "../common/generatePeriodSerie"
import { dateAddMonth } from "./dateAddMonth"

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

  const expPartOuvrière = Array(moisRemboursement)
    .fill(100)
    .concat(Array(dureeEnMois - moisRemboursement).fill(0))

  const expPartPatronale = Array(moisRemboursement)
    .fill(200)
    .concat(Array(dureeEnMois - moisRemboursement).fill(0))

  const expectedInteressanteUrssaf = Array(9)
    .fill(false)
    .concat(Array(4).fill(undefined))

  for (let mois = 0; mois < 13; ++mois) {
    t.log({ période: mois }, actual[dateAddMonth(dateDebut, mois).getTime()])
    const expected = {
      interessante_urssaf: expectedInteressanteUrssaf[mois],
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
    Object.keys(expected).forEach((p) => {
      const prop = p as keyof typeof expected
      if (typeof expected[prop] === "undefined") {
        delete expected[prop]
      }
    })
    t.deepEqual(actual[dateAddMonth(dateDebut, mois).getTime()], expected)
  }
})

// TODO: comportement surprenant avec debut_suivant, le test passe a moitié avec valeur non nulle mais echoue avec chaine vide
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
      // dette initiale
      hash1: {
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
      // remboursement la dette
      hash2: {
        periode: {
          start: dateDebut,
          end: dateAddMonth(dateDebut, 1),
        },
        numero_ecart_negatif: 1,
        numero_historique: 3,
        numero_compte: "3",
        date_traitement: dateAddMonth(dateDebut, 1),
        debit_suivant: "",
        part_ouvriere: 0,
        part_patronale: 0,
      },
    },
  }

  const actual = cotisationsdettes(v, periode)

  t.log(actual)

  for (const month of [0, 1, 2, 3, 4, 5, 6]) {
    t.false(
      actual[dateAddMonth(dateDebut, month).getTime()].interessante_urssaf
    )
  }

  t.is(
    actual[dateAddMonth(dateDebut, 7).getTime()].interessante_urssaf,
    undefined
  )
})
