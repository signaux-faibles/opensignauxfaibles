import "../globals"
import test, { ExecutionContext } from "ava"
import { cotisationsdettes, SortieCotisationsDettes } from "./cotisationsdettes"
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

test.only("Le montant de dette d'une période est reporté dans les périodes suivantes", (t: ExecutionContext) => {
  const dateDebut = new Date("2018-01-01")
  const dateFin = dateAddMonth(dateDebut, 13)
  const periode = generatePeriodSerie(dateDebut, dateFin).map((date) =>
    date.getTime()
  )

  ;(globalThis as any).date_fin = dateFin // TODO: transformer ce parametre en parametre local de fonction

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
        date_traitement: dateAddMonth(dateDebut, 4),
        debit_suivant: "",
        numero_compte: "",
        numero_ecart_negatif: 1,
        numero_historique: 3,
      },
    },
  }

  const actual = cotisationsdettes(v, periode)

  const montants = {
    montant_part_ouvriere: 100,
    montant_part_patronale: 200,
    interessante_urssaf: false,
  }

  const montantsUnMois = {
    ...montants,
    montant_part_ouvriere_past_1: 100,
    montant_part_patronale_past_1: 200,
  }

  const montantsDeuxMois = {
    ...montantsUnMois,
    montant_part_ouvriere_past_2: 100,
    montant_part_patronale_past_2: 200,
  }

  const montantsTroisMois = {
    ...montantsDeuxMois,
    montant_part_ouvriere_past_3: 100,
    montant_part_patronale_past_3: 200,
  }

  const montantsSixMois = {
    ...montantsTroisMois,
    montant_part_ouvriere_past_6: 100,
    montant_part_patronale_past_6: 200,
  }

  const montantsDouzeMois = {
    ...montantsSixMois,
    montant_part_ouvriere_past_12: 100,
    montant_part_patronale_past_12: 200,
  }

  t.deepEqual(actual[dateAddMonth(dateDebut, 1).getTime()], montantsUnMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 2).getTime()], montantsDeuxMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 3).getTime()], montantsTroisMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 4).getTime()], {...montantsTroisMois, montant_part_ouvriere: 0, montant_part_patronale: 0} as SortieCotisationsDettes)
  t.deepEqual(actual[dateAddMonth(dateDebut, 5).getTime()], {...montantsTroisMois, montant_part_ouvriere: 0, montant_part_patronale: 0} as SortieCotisationsDettes)
  t.deepEqual(actual[dateAddMonth(dateDebut, 6).getTime()], {...montantsSixMois, montant_part_ouvriere: 0, montant_part_patronale: 0} as SortieCotisationsDettes)
  t.deepEqual(actual[dateAddMonth(dateDebut, 7).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 8).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 9).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 10).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 11).getTime()], montantsSixMois)
  t.deepEqual(actual[dateAddMonth(dateDebut, 12).getTime()], montantsDouzeMois)
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
