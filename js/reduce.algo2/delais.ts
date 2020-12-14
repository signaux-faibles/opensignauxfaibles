import { f } from "./functions"
import { EntréeDelai, ParHash, ParPériode } from "../RawDataTypes"

type DeepReadonly<T> = Readonly<T> // pas vraiment, mais espoire que TS le supporte prochainement

// Valeurs attendues pour chaque période, lors de l'appel à delais()
export type DebitComputedValues = {
  montant_part_patronale: number
  montant_part_ouvriere: number
}

// Valeurs retournées par delais(), pour chaque période
export type DelaiComputedValues = {
  // valeurs fournies, reportées par delais() dans chaque période:
  delai_nb_jours_total: number // nombre de jours entre date_creation et date_echeance
  delai_montant_echeancier: number // exprimé en euros
  // valeurs calculées par delais():
  delai_nb_jours_restants: number
  delai_deviation_remboursement?: number // ratio entre remboursement linéaire et effectif, à condition d'avoir le montant des parts ouvrière et patronale
}

/**
 * Calcule pour chaque période le nombre de jours restants du délai accordé et
 * un indicateur de la déviation par rapport à un remboursement linéaire du
 * montant couvert par le délai. Un "délai" étant une demande accordée de délai
 * de paiement des cotisations sociales, pour un certain montant
 * (delai_montant_echeancier) et pendant une certaine période
 * (delai_nb_jours_total).
 * Contrat: cette fonction ne devrait être appelée que s'il y a eu au moins une
 * demande de délai.
 */
export function delais(
  vDelai: ParHash<EntréeDelai>,
  debitParPériode: DeepReadonly<ParPériode<DebitComputedValues>>,
  intervalleTraitement: { premièreDate: Date; dernièreDate: Date }
): ParPériode<DelaiComputedValues> {
  "use strict"
  const donnéesDélaiParPériode: ParPériode<DelaiComputedValues> = {}
  Object.values(vDelai).forEach((delai) => {
    if (delai.duree_delai <= 0) {
      return
    }

    // On arrondit les dates au premier jour du mois.
    const date_creation = new Date(
      Date.UTC(
        delai.date_creation.getUTCFullYear(),
        delai.date_creation.getUTCMonth(),
        1,
        0,
        0,
        0,
        0
      )
    )
    const date_echeance = new Date(
      Date.UTC(
        delai.date_echeance.getUTCFullYear(),
        delai.date_echeance.getUTCMonth(),
        1,
        0,
        0,
        0,
        0
      )
    )
    // Création d'un tableau de timestamps à raison de 1 par mois.
    f.generatePeriodSerie(date_creation, date_echeance)
      .filter(
        (date) =>
          date >= intervalleTraitement.premièreDate &&
          date <= intervalleTraitement.dernièreDate
      )
      .map(function (debutDeMois) {
        const time = debutDeMois.getTime()
        const remainingDays = f.nbDays(debutDeMois, delai.date_echeance)
        const inputAtTime = debitParPériode[time]
        const outputAtTime: DelaiComputedValues = {
          delai_nb_jours_restants: remainingDays,
          delai_nb_jours_total: delai.duree_delai,
          delai_montant_echeancier: delai.montant_echeancier,
        }
        if (
          typeof inputAtTime?.montant_part_patronale !== "undefined" &&
          typeof inputAtTime?.montant_part_ouvriere !== "undefined"
        ) {
          const detteActuelle =
            inputAtTime.montant_part_patronale +
            inputAtTime.montant_part_ouvriere
          const detteHypothétiqueRemboursementLinéaire =
            (delai.montant_echeancier * remainingDays) / delai.duree_delai
          outputAtTime.delai_deviation_remboursement =
            (detteActuelle - detteHypothétiqueRemboursementLinéaire) /
            delai.montant_echeancier
        }
        donnéesDélaiParPériode[time] = outputAtTime
      })
  })
  return donnéesDélaiParPériode
}