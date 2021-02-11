/* eslint-disable no-use-before-define */

import {
  EntréeApConso,
  EntréeApDemande,
  EntréeDelai,
  EntréeCcsf,
  EntréeCompte,
  EntréeCotisation,
  EntréeDebit,
  EntréeDéfaillances,
  EntréeDiane,
  EntréeEffectif,
  EntréeEffectifEnt,
  EntréeEllisphere,
  EntréePaydex,
  EntréeSirene,
  EntréeSireneEntreprise,
} from "./GeneratedTypes"

// Types de données de base

export type Periode = string // Date.getTime().toString()
export type Timestamp = number // Date.getTime()
export type ParPériode<T> = Record<Periode, T>

export type Departement = string

export type Siren = string
export type Siret = string
export type SiretOrSiren = Siret | Siren // TODO: supprimer ce type, une fois que tous les champs auront été séparés
export type CodeAPE = string

export type DataHash = string
export type ParHash<T> = Record<DataHash, T>

// Données importées pour une entreprise ou établissement

export type EntrepriseDataValues = {
  key: Siren
  scope: "entreprise"
  batch: Record<BatchKey, Partial<BatchValueProps>> // TODO: remplacer par `Partial<EntrepriseBatchProps>>`, une fois que tous les champs auront été séparés
}

export type EtablissementDataValues = {
  key: Siret
  scope: "etablissement"
  batch: Record<BatchKey, Partial<BatchValueProps>> // TODO: remplacer par Partial<EtablissementBatchProps>, une fois que tous les champs auront été séparés
}

export type Scope = (EntrepriseDataValues | EtablissementDataValues)["scope"]

export type CompanyDataValues = EntrepriseDataValues | EtablissementDataValues

export type CompanyDataValuesWithFlags = CompanyDataValues & IndexFlags

export type IndexFlags = {
  index: {
    algo2: boolean // pour spécifier quelles données seront à calculer puis inclure dans Features, par Reduce.algo2
  }
}

// Données importées par les parseurs, pour chaque source de données

export type BatchKey = string

export type BatchValues = Record<BatchKey, BatchValue>

export type DataType = keyof BatchValueProps // => 'reporder' | 'effectif' | 'apconso' | ...

export type BatchValue = Partial<BatchValueProps>

type CommonBatchProps = {
  reporder: ParPériode<EntréeRepOrder> // RepOrder est généré par "compact", et non importé => Usage de Periode en guise de hash d'indexation
}

export type EntrepriseBatchProps = CommonBatchProps & {
  paydex: ParHash<EntréePaydex>
}

export type EtablissementBatchProps = CommonBatchProps & {
  apconso: ParHash<EntréeApConso>
}

// TODO: continuer d'extraire les propriétés vers EntrepriseBatchProps et EtablissementBatchProps, puis supprimer BatchValueProps et les types qui en dépendent
export type BatchValueProps = CommonBatchProps &
  EntrepriseBatchProps &
  EtablissementBatchProps & {
    effectif: ParHash<EntréeEffectif>
    apdemande: ParHash<EntréeApDemande>
    compte: ParHash<EntréeCompte>
    delai: ParHash<EntréeDelai>
    procol: ParHash<EntréeDéfaillances>
    cotisation: ParHash<EntréeCotisation>
    debit: ParHash<EntréeDebit>
    ccsf: ParHash<EntréeCcsf>
    sirene: ParHash<EntréeSirene>
    sirene_ul: ParHash<EntréeSireneEntreprise>
    effectif_ent: ParHash<EntréeEffectifEnt>
    bdf: ParHash<EntréeBdf>
    diane: ParHash<EntréeDiane>
    ellisphere: ParHash<EntréeEllisphere>
  }

// Détail des types de données

export type EntréeRepOrder = {
  random_order: number
  periode: Date
  siret: SiretOrSiren
}

export type EntréeBdf = {
  /** Date de clôture de l'exercice. */
  arrete_bilan_bdf: Date
  /** Année de l'exercice. */
  annee_bdf: number
  /** Raison sociale de l'entreprise. */
  raison_sociale: string
  /** Secteur d'activité. */
  secteur: string
  /** Siren de l'entreprise. */
  siren: SiretOrSiren
  /** Poids du fonds de roulement net global sur le chiffre d'affaire. Exprimé en %. */
  poids_frng: number
  /** Taux de marge, rapport de l'excédent brut d'exploitation (EBE) sur la valeur ajoutée (exprimé en %): 100*EBE / valeur ajoutee */
  taux_marge: number
  /** Délai estimé de paiement des fournisseurs (exprimé en jours): 360 * dettes fournisseurs / achats HT */
  delai_fournisseur: number
  /** Poids des dettes fiscales et sociales, par rapport à la valeur ajoutée (exprimé en %): 100 * dettes fiscales et sociales / Valeur ajoutee */
  dette_fiscale: number
  /** Poids du financement court terme (exprimé en %): 100 * concours bancaires courants / chiffre d'affaires HT */
  financier_court_terme: number
  /** Poids des frais financiers, sur l'excedent brut d'exploitation corrigé des produits et charges hors exploitation (exprimé en %): 100 * frais financiers / (EBE + Produits hors expl. - charges hors expl.) */
  frais_financier: number
}
