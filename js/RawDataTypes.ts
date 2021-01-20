/* eslint-disable no-use-before-define */

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
    effectif_ent: ParHash<EntréeEffectif>
    bdf: ParHash<EntréeBdf>
    diane: ParHash<EntréeDiane>
    ellisphere: ParHash<EntréeEllisphere>
  }

// Détail des types de données

export type EntréeCcsf = {
  /** Date de début de la procédure CCSF */
  date_traitement: Date
  stade: string // TODO: choisir un type plus précis
  action: string // TODO: choisir un type plus précis
}

export type EntréeDéfaillances = {
  /** Nature de la procédure de défaillance. */
  action_procol: "liquidation" | "redressement" | "sauvegarde"
  /** Evénement survenu dans le cadre de cette procédure. */
  stade_procol:
    | "abandon_procedure"
    | "solde_procedure"
    | "fin_procedure"
    | "plan_continuation"
    | "ouverture"
    | "inclusion_autre_procedure"
    | "cloture_insuffisance_actif"
  /** Date effet de la procédure collective. */
  date_effet: Date
}

export type EntréeApConso = {
  id_conso: string
  periode: Date
  heure_consomme: number
}

export type EntréeApDemande = {
  id_demande: string
  periode: { start: Date; end: Date }
  hta: number /* Nombre total d'heures autorisées */
  motif_recours_se: number /* Cause d'activité partielle */
  effectif_entreprise?: number
  effectif?: number
  date_statut?: Date
  mta?: number
  effectif_autorise?: number
  heure_consomme?: number
  montant_consomme?: number
  effectif_consomme?: number
}

export type EntréeCompte = {
  /** Date à laquelle cet établissement est associé à ce numéro de compte URSSAF. */
  periode: Date
  /** Numéro SIRET de l'établissement. Les numéros avec des Lettres sont des sirets provisoires. */
  siret: string
  /** Compte administratif URSSAF. */
  numero_compte: number
}

export type EntréeInterim = {
  periode: Date
  etp: number
}

export type EntréeRepOrder = {
  random_order: number
  periode: Date
  siret: SiretOrSiren
}

export type EntréeEffectif = {
  numero_compte: string
  periode: Date
  /** Nombre de personnes employées par l'entreprise ou l'établissement. */
  effectif: number
}

// Valeurs attendues par delais(), pour chaque période. (cf lib/urssaf/delai.go)
export type EntréeDelai = {
  date_creation: Date
  date_echeance: Date
  /** Nombre de jours entre date_creation et date_echeance. */
  duree_delai: number
  /** Montant global de l'échéancier, en euros. */
  montant_echeancier: number
}

export type EntréeCotisation = {
  periode: { start: Date; end: Date }
  du: number
}

/**
 * Représente un reste à payer (dette) sur cotisation sociale ou autre.
 */
export type EntréeDebit = {
  periode: { start: Date; end: Date } // Periode pour laquelle la cotisation était attendue
  numero_ecart_negatif: number // identifiant du débit pour une période donnée (comme une sorte de numéro de facture)
  numero_historique: number // identifiant d'un remboursement (partiel ou complet) d'un débit
  numero_compte: string // identifiant URSSAF d'établissement (équivalent du SIRET)
  date_traitement: Date // Date de constatation du débit (exemple: remboursement, majoration ou autre modification du montant)
  debit_suivant: string // Hash d'un autre débit
  // le montant est ventilé entre ces deux valeurs, exprimées en euros (€):
  part_ouvriere: number
  part_patronale: number
  montant_majorations?: number // exploité par map-reduce "public", mais pas par "reduce.algo2"
}

export type EntréeSirene = {
  ape: CodeAPE
  latitude: number
  longitude: number
  departement: Departement
  raison_sociale: string
  date_creation: Date
}

export type EntréeSireneEntreprise = {
  raison_sociale: string
  nom_unite_legale: string
  nom_usage_unite_legale: string
  prenom1_unite_legale: string
  prenom2_unite_legale: string
  prenom3_unite_legale: string
  prenom4_unite_legale: string
  statut_juridique: string | null // code numérique sérialisé en chaine de caractères
  date_creation: Date
}

export type EntréeBdf = {
  arrete_bilan_bdf: Date
  annee_bdf: number
  exercice_bdf: number
  raison_sociale: string
  secteur: string
  siren: SiretOrSiren
} & EntréeBdfRatios

export type EntréeBdfRatios = {
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

/**
 * Champs récupérés lors de l'import depuis des fichiers Diane.
 * Le commentaire de chaque champ permet de générer sa documentation.
 * Cf https://github.com/signaux-faibles/opensignauxfaibles/pull/291.
 */
export type EntréeDiane = {
  /** Année de l'exercice */
  exercice_diane: number
  /** Date d'arrêté du bilan */
  arrete_bilan_diane: Date
  /** Couverture du chiffre d'affaire par le fonds de roulement (exprimé en jours): Fonds de roulement net global / Chiffre d'affaires net * 360 */
  couverture_ca_fdr?: number
  /** Intérêts et charges assimilées. */
  interets?: number
  /** Excédent brut d'exploitation. */
  excedent_brut_d_exploitation?: number
  /** Produits financiers. */
  produits_financiers?: number
  /** Produits exceptionnels. */
  produit_exceptionnel?: number
  /** Charges exceptionnelles. */
  charge_exceptionnelle?: number
  /** Charges financières. */
  charges_financieres?: number
  /** Chiffre d'affaires */
  ca?: number
  /** Concours bancaires courants. (Pour recalculer les frais financiers court terme de la Banque de France) */
  concours_bancaire_courant?: number
  /** Valeur ajoutée. */
  valeur_ajoutee?: number
  /** Dette fiscale et sociale */
  dette_fiscale_et_sociale?: number
  marquee: unknown // TODO: propriété non trouvée en sortie du parseur Diane => à supprimer ?
  /** Raison sociale */
  nom_entreprise: string
  /** Numéro siren */
  numero_siren: SiretOrSiren
  /** Statut juridique */
  statut_juridique: string
  /** Présence d'une procédure collective en cours */
  procedure_collective: boolean
  /** Effectif consolidé à l'entreprise */
  effectif_consolide: number
  /** Frais de Recherche et Développement */
  frais_de_RetD: number
  /** Concessions, brevets, et droits similaires */
  conces_brev_et_droits_sim: number
  /** Nombre d'établissements secondaires de l'entreprise, en plus du siège. */
  nombre_etab_secondaire: number
  /** Nombre de filiales de l'entreprise. Dans la base de données des liens capitalistiques, le concept de filiale ne fait aucune référence au pourcentage d’appartenance entre le parent et la fille. Dans ce sens, si l'entreprise A est enregistrée comme ayant des intérêts dans l'entreprise B avec un très petit, ou même un pourcentage de participation inconnu, l'entreprise B sera considérée filiale de l'entreprise A. */
  nombre_filiale: number
  /** Nombre d'entreprises dans le groupe (groupe défini par les liens capitalistique d'au moins 50,01%) */
  taille_compo_groupe: number
  /** Durée de l'exercice en mois. */
  nombre_mois: number
  /** Équilibre financier: Ressources durables / Emplois stables */
  equilibre_financier: number
  /** Indépendance financière (exprimé en %): Fonds propres / Ressources durables * 100 */
  independance_financiere: number
  /** Endettement (exprimé en %): Dettes de caractère financier / Ressources durables * 100 */
  endettement: number
  /** Autonomie financière Fonds propres / Total bilan * 100 */
  autonomie_financiere: number
  /** Degré d'amortissement des immobilisations corporelles (exprimé en %): Amortissements des immobilisations corporelles / Immobilisation corporelles brutes * 100 */
  degre_immo_corporelle: number
  /** Financement de l'actif circulant net: Fonds de roulement net global / Actif circulant net */
  financement_actif_circulant: number
  /** Liquidité générale: Actif circulant net / Dettes à court terme */
  liquidite_generale: number
  /** Liquidité réduite: Actif circulant net hors stocks / Dettes à court terme */
  liquidite_reduite: number
  /** Rotation des stocks (exprimé en jours): Stock / Chiffre d'affaires net * 360. Selon la nomenclature NAF Rév. 2 pour les secteurs d'activité 45, 46, 47, 95 (sauf 9511Z) ainsi que pour les codes d'activités 2319Z, 3831Z et 3832Z : Marchandises / (Achats de marchandises + Variation de stock) * 360 */
  rotation_stocks: number
  /** Crédit clients (exprimé en jours): (Clients + Effets portés à l'escompte et non échus) / Chiffre d'affaires TTC * 360 */
  credit_client: number
  /** Crédit fournisseurs (exprimé en jours): Fournisseurs / Achats TTC * 360 */
  credit_fournisseur: number
  /** Chiffre d'affaire par effectif (exprimé en k€/emploi): Chiffre d'affaires net / Effectif * 1000 */
  ca_par_effectif: number
  /** Taux d'intérêt financier (exprimé en %): Intérêts / Chiffre d'affaires net * 100 */
  taux_interet_financier: number
  /** Intérêts sur chiffre d'affaire (exprimé en %): Total des charges financières / Chiffre d'affaires net * 100 */
  taux_interet_sur_ca: number
  /** Endettement global (exprimé en jours): (Dettes + Effets portés à l'escompte et non échus) / Chiffre d'affaires net * 360 */
  endettement_global: number
  /** Taux d'endettement (exprimé en %): Dettes de caractère financier / (Capitaux propres + autres fonds propres) * 100 */
  taux_endettement: number
  /** Capacité de remboursement: Dettes de caractère financier / Capacité d'autofinancement avant répartition */
  capacite_remboursement: number
  /** Capacité d'autofinancement (exprimé en %): Capacité d'autofinancement avant répartition / (Chiffre d'affaires net + Subvention d'exploitation) * 100 */
  capacite_autofinancement: number
  /** Couverture du chiffre d'affaire par le besoin en fonds de roulement (exprimé en jours): Besoins en fonds de roulement / Chiffre d'affaires net * 360 */
  couverture_ca_besoin_fdr: number
  /** PoidsBFRExploitation Poids des besoins en fonds de roulement d'exploitation (exprimé en %): Besoins en fonds de roulement d'exploitation / Chiffre d'affaires net * 100 */
  poids_bfr_exploitation: number
  /** Exportation Exportation (exprimé en %): (Chiffre d'affaires net - Chiffre d'affaires net en France) / Chiffre d'affaires net * 100 */
  exportation: number
  /** Efficacité économique (exprimé en k€/emploi): Valeur ajoutée / Effectif * 1000 */
  efficacite_economique: number
  /** Productivité du potentiel de production: Valeur ajoutée / Immobilisations corporelles et incorporelles brutes */
  productivite_potentiel_production: number
  /** Productivtié du capital financier: Valeur ajoutée / Actif circulant net + Effets portés à l'escompte et non échus */
  productivite_capital_financier: number
  /** Productivité du capital investi: Valeur ajoutée / Total de l'actif + Effets portés à l'escompte et non échus */
  productivite_capital_investi: number
  /** Taux d'investissement productif (exprimé en %): Immobilisations à valeur d'acquisition / Valeur ajoutée * 100 */
  taux_d_investissement_productif: number
  /** Rentabilité économique (exprimé en %): Excédent brut d'exploitation / Chiffre d'affaires net + Subventions d'exploitation * 100 */
  rentabilite_economique: number
  /** Performance (exprimé en %): Résultat courant avant impôt / Chiffre d'affaires net + Subventions d'exploitation * 100 */
  performance: number
  /** Rendement brut des fonds propres (exprimé en %): Résultat courant avant impôt / Fonds propres nets * 100 */
  rendement_brut_fonds_propres: number
  /** Rentabilité nette (exprimé en %): Bénéfice ou perte / Chiffre d'affaires net + Subventions d'exploitation * 100 */
  rentabilite_nette: number
  /** Rendement des capitaux propres (exprimé en %): Bénéfice ou perte / Capitaux propres nets * 100 */
  rendement_capitaux_propres: number
  /** RendementRessourcesDurables Rendement des ressources durables (exprimé en %): Résultat courant avant impôts + Intérêts et charges assimilées / Ressources durables nettes * 100 */
  rendement_ressources_durables: number
  /** Taux de marge commerciale (exprimé en %): Marge commerciale / Vente de marchandises * 100 */
  taux_marge_commerciale: number
  /** Taux de valeur ajoutée (exprimé en %): Valeur ajoutée / Chiffre d'affaires net * 100 */
  taux_valeur_ajoutee: number
  /** Part des salariés (exprimé en %): (Charges de personnel + Participation des salariés aux résultats) / Valeur ajoutée * 100 */
  part_salaries: number
  /** Part de l'État (exprimé en %): Impôts et taxes / Valeur ajoutée * 100 */
  part_etat: number
  /** Part des prêteurs (exprimé en %): Intérêts / Valeur ajoutée * 100 */
  part_preteur: number
  /** Part de l'autofinancement (exprimé en %): Capacité d'autofinancement avant répartition / Valeur ajoutée * 100 */
  part_autofinancement: number
  /** Chiffre d'affaires à l'exportation */
  ca_exportation: number
  /** Achats de marchandises */
  achat_marchandises: number
  /** Achats de matières premières et autres approvisionnement. */
  achat_matieres_premieres: number
  /** Production de l'exercice. */
  production: number
  /** Marge commerciale. */
  marge_commerciale: number
  /** Consommation de l'exercice. */
  consommation: number
  /** Autres achats et charges externes. */
  autres_achats_charges_externes: number
  /** Charges de personnel. */
  charge_personnel: number
  /** Impôts, taxes et versements assimilés. */
  impots_taxes: number
  /** Subventions d'exploitation. */
  subventions_d_exploitation: number
  /** Autres produits, charges et reprises. */
  autres_produits_charges_reprises: number
  /** Dotation d'exploitation aux amortissements et aux provisions. */
  dotation_amortissement: number
  /** Résultat d'exploitation. */
  resultat_expl: number
  /** Opérations en commun. */
  operations_commun: number
  /** Résultat courant avant impôts. */
  resultat_avant_impot: number
  /** Participation des salariés aux résultats. */
  participation_salaries: number
  /** Impôts sur les bénéfices et impôts différés. */
  impot_benefice: number
  /** Bénéfice ou perte. */
  benefice_ou_perte: number
}

export type EntréeEllisphere = {
  siren: string
  code_groupe: string
  siren_groupe: string
  refid_groupe: string
  raison_sociale_groupe: string
  adresse_groupe: string
  personne_pou_m_groupe: string
  niveau_detention: number
  part_financiere: number
  code_filiere: string
  refid_filiere: string
  personne_pou_m_filiere: string
}

export type EntréePaydex = {
  date_valeur: Date
  nb_jours: number
}
