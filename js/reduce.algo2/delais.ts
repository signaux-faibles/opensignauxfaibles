import { f } from "./functions"
import { EntréeDelai, ParHash, ParPériode } from "../RawDataTypes"
import { SortieCotisationsDettes } from "./cotisationsdettes"

type DeepReadonly<T> = Readonly<T> // pas vraiment immutable pout l'instant, mais espoir que TS le permette prochainement

export type ChampsEntréeDelai = Pick<
  EntréeDelai,
  "date_creation" | "date_echeance" | "duree_delai" | "montant_echeancier"
>

// valeurs fournies, reportées par delais() dans chaque période
type ValeursTransmises = {
  /** Nombre de jours entre date_creation et date_echeance. */
  delai_nb_jours_total: EntréeDelai["duree_delai"]
  /** Montant global de l'échéancier, en euros. */
  delai_montant_echeancier: EntréeDelai["montant_echeancier"]
}

// valeurs calculées par delais()
type ValeursCalculuées = {
  /** Nombre de jours restants du délai. */
  delai_nb_jours_restants: number
  /** Ratio entre remboursement linéaire et effectif, à condition d'avoir le montant des parts ouvrière et patronale. */
  delai_deviation_remboursement?: number
}

// Valeurs retournées par delais(), pour chaque période
export type SortieDelais = ValeursTransmises & ValeursCalculuées

// Variables est inspecté pour générer docs/variables.json (cf generate-docs.ts)
export type Variables = {
  source: "delais"
  computed: ValeursCalculuées
  transmitted: ValeursTransmises
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
  vDelai: ParHash<ChampsEntréeDelai>,
  debitParPériode: DeepReadonly<ParPériode<SortieCotisationsDettes>>,
  intervalleTraitement: { premièreDate: Date; dernièreDate: Date }
): ParPériode<SortieDelais> {
  "use strict"
  const donnéesDélaiParPériode: ParPériode<SortieDelais> = {}
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
        const outputAtTime: SortieDelais = {
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
