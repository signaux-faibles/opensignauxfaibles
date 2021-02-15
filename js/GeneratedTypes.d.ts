/**
 * This file was automatically generated by generate-types.sh.
 *
 * DO NOT MODIFY IT BY HAND.
 *
 * Instead:
 * - modify the validation/*.schema.json files;
 * - then, run generate-types.sh to regenerate this file.
 */

/**
 * Champs importés par le parseur lib/apconso/main.go de sfdata.
 */
export interface EntréeApConso {
  /**
   * Numéro de la demande (11 caractères principalement des chiffres)
   */
  id_conso: string
  /**
   * Heures consommées (chômées) dans le mois
   */
  heure_consomme: number
  /**
   * Montants consommés dans le mois
   */
  montant: number
  /**
   * Nombre de salariés en activité partielle dans le mois
   */
  effectif: number
  /**
   * Mois considéré
   */
  periode: Date
}
/**
 * Champs importés par le parseur lib/apdemande/main.go de sfdata.
 */
export interface EntréeApDemande {
  /**
   * Numéro de la demande (11 caractères principalement des chiffres)
   */
  id_demande: string
  /**
   * Période de chômage
   */
  periode: {
    start: Date
    end: Date
  }
  /**
   * Nombre total d'heures autorisées
   */
  hta: number
  /**
   * Motif de recours à l'activité partielle:
   *  1 - Conjoncture économique.
   *  2 - Difficultés d’approvisionnement en matières premières ou en énergie
   *  3 - Sinistre ou intempéries de caractère exceptionnel
   *  4 - Transformation, restructuration ou modernisation des installations et des bâtiments
   *  5 - Autres circonstances exceptionnelles
   */
  motif_recours_se: number
  /**
   * Effectif de l'entreprise
   */
  effectif_entreprise: number
  /**
   * Effectif de l'établissement
   */
  effectif: number
  /**
   * Date du statut - création ou mise à jour de la demande
   */
  date_statut: Date
  /**
   * Montant total autorisé
   */
  mta: number
  /**
   * Effectifs autorisés
   */
  effectif_autorise: number
  /**
   * Nombre total d'heures consommées
   */
  heure_consommee: number
  /**
   * Montant total consommé
   */
  montant_consommee: number
  effectif_consomme: number
}
/**
 * Note: CE SCHEMA EST INCOMPLET POUR L'INSTANT. Cf https://github.com/signaux-faibles/opensignauxfaibles/pull/143
 */
export interface EntréeBdf {
  siren: string
}
/**
 * Champs importés par le parseur lib/urssaf/ccsf.go de sfdata.
 */
export interface EntréeCcsf {
  /**
   * Date de début de la procédure CCSF
   */
  date_traitement: Date
  /**
   * TODO: choisir un type plus précis
   */
  stade: string
  /**
   * TODO: choisir un type plus précis
   */
  action: string
}
/**
 * Champs importés par le parseur lib/urssaf/compte.go de sfdata.
 */
export interface EntréeCompte {
  /**
   * Date à laquelle cet établissement est associé à ce numéro de compte URSSAF.
   */
  periode: Date
  /**
   * Numéro SIRET de l'établissement. Les numéros avec des Lettres sont des sirets provisoires.
   */
  siret: string
  /**
   * Compte administratif URSSAF.
   */
  numero_compte: string
}
/**
 * Champs importés par le parseur lib/urssaf/cotisation.go de sfdata.
 */
export interface EntréeCotisation {
  /**
   * Compte administratif URSSAF.
   */
  numero_compte: string
  /**
   * Période sur laquelle le montants s'appliquent.
   */
  periode: {
    start: Date
    end: Date
  }
  /**
   * Cotisation encaissée directement, en euros.
   */
  encaisse: number
  /**
   * Cotisation due, en euros. À utiliser pour calculer le montant moyen mensuel du: Somme cotisations dues / nb périodes.
   */
  du: number
}
/**
 * Champs importés par le parseur lib/sirene_ul/main.go de sfdata. Représente un reste à payer (dette) sur cotisation sociale ou autre.
 */
export interface EntréeDebit {
  /**
   * Identifiant URSSAF d'établissement (équivalent du SIRET).
   */
  numero_compte: string
  /**
   * Période sur laquelle le montants s'appliquent.
   */
  periode: {
    start: Date
    end: Date
    [k: string]: unknown
  }
  /**
   * L'écart négatif (ecn) correspond à une période en débit. Pour une même période, plusieurs débits peuvent être créés. On leur attribue un numéro d'ordre. Par exemple, 101, 201, 301 etc.; ou 101, 102, 201 etc. correspondent respectivement au 1er, 2ème et 3ème ecn de la période considérée.
   */
  numero_ecart_negatif: string
  /**
   * Ordre des opérations pour un écart négatif donné.
   */
  numero_historique: number
  /**
   * Date de constatation du débit (exemple: remboursement, majoration ou autre modification du montant)
   */
  date_traitement: Date
  /**
   * Hash d'un autre débit // TODO: non fourni par le parseur, ce champ devrait être défini dans un type de sortie.
   */
  debit_suivant?: string
  /**
   * Montant des débits sur la part ouvrières, exprimées en euros (€). Sont exclues les pénalités et les majorations de retard.
   */
  part_ouvriere: number
  /**
   * Montant des débits sur la part patronale, exprimées en euros (€). Sont exclues les pénalités et les majorations de retard.
   */
  part_patronale: number
  /**
   * Code état du compte: 1 (Actif), 2 (Suspendu) ou 3 (Radié).
   */
  etat_compte: number
  /**
   * Code qui indique si le compte fait l'objet d'une procédure collective: 1 (en cours), 2 (plan de redressement en cours), 9 (procedure collective sans dette à l'Urssaf) ou valeur nulle en cas d'absence de procédure collective.
   */
  code_procedure_collective: string
  /**
   * Code opération historique de l'écart négatif:
   *  1 Mise en recouvrement
   *  2 Paiement
   *  3 Admission en non valeur
   *  4 Remise de majoration de retard
   *  5 Abandon de solde debiteur
   *  11 Annulation de mise en recouvrement
   *  12 Annulation paiement
   *  13 Annulation a-n-v
   *  14 Annulation de remise de majoration retard
   *  15 Annulation abandon solde debiteur
   */
  code_operation_ecart_negatif: string
  /**
   * Code motif de l'écart négatif:
   *  0 Cde motif inconnu
   *  1 Retard dans le versement
   *  2 Absence ou insuffisance de versement
   *  3 Taxation provisionelle. Déclarations non fournies
   *  4 Majorations de retard complémentaires Article R243-18 du code de la sécurité sociale
   *  5 Contrôle,chefs de redressement notifiés le JJ/MM/AA Article R243-59 de la Securité Sociale
   *  6 Fourniture tardive des déclarations
   *  7 Bases déclarées supérieures à Taxation provisionnelle
   *  8 Retard dans le versement et fourniture tardive des déclarations
   *  9 Absence ou insuffisance de versement et fourniture tardive des déclarations
   *  10 Rappel sur contrôle et fourniture tardive des déclarations
   *  11 Régularisation d'une taxation provisionnelle
   *  12 Régularisation annuelle
   *  13 Rejet du titre de paiement par la banque .
   *  14 Modification d'affectation d'un crédit
   *  15 Annulation d'un crédit
   *  16 Régularisation suite à modification du Taux Accident du Travail
   *  17 Régularisation suite à assujettissement au transport (origine débit sur PJ=4)
   *  18 Majorations pour non respect de paiement par moyen dématérialisé Article L243-14
   *  19 Rapprochement TR/BRC sous réserve de vérification ultérieure
   *  20 Cotisations complémentaires suite modification des revenus déclarés
   *  21 Cotisations complémentaires suite à non fourniture du contrat d'exonération
   *  22 Contrôle. Chefs de redressement notifiés le JJ/MM/AA. Article L324.9 du code du travail
   *  23 Cotisations complémentaires suite conditions d'exonération non remplies
   *  24 Absence de versement
   *  25 Insuffisance de versement
   *  26 Absence de versement et fourniture tardive des déclarations
   *  27 Insuffisance de versement et fourniture tardive des déclarations
   */
  code_motif_ecart_negatif: string
  /**
   * Recours en cours.
   */
  recours_en_cours: {
    [k: string]: unknown
  }
}
/**
 * Champs importés par le parseur lib/urssaf/delai.go de sfdata.
 */
export interface EntréeDelai {
  /**
   * Compte administratif URSSAF.
   */
  numero_compte: string
  /**
   * Le numéro de structure est l'identifiant d'un dossier contentieux.
   */
  numero_contentieux: string
  /**
   * Date de création du délai.
   */
  date_creation: Date
  /**
   * Date d'échéance du délai.
   */
  date_echeance: Date
  /**
   * Durée du délai en jours: nombre de jours entre date_creation et date_echeance.
   */
  duree_delai: number
  /**
   * Raison sociale de l'établissement.
   */
  denomination: string
  /**
   * Délai inférieur ou supérieur à 6 mois ? Modalités INF et SUP.
   */
  indic_6m: string
  /**
   * Année de création du délai.
   */
  annee_creation: number
  /**
   * Montant global de l'échéancier, en euros.
   */
  montant_echeancier: number
  /**
   * Code externe du stade.
   */
  stade: string
  /**
   * Code externe de l'action.
   */
  action: string
}
/**
 * Champs importés par le parseur lib/diane/main.go de sfdata.
 */
export interface EntréeDiane {
  /**
   * Année de l'exercice
   */
  exercice_diane?: number
  /**
   * Date d'arrêté du bilan
   */
  arrete_bilan_diane?: Date
  /**
   * Couverture du chiffre d'affaire par le fonds de roulement (exprimé en jours): Fonds de roulement net global / Chiffre d'affaires net * 360
   */
  couverture_ca_fdr?: number
  /**
   * Intérêts et charges assimilées.
   */
  interets?: number
  /**
   * Excédent brut d'exploitation.
   */
  excedent_brut_d_exploitation?: number
  /**
   * Produits financiers.
   */
  produits_financiers?: number
  /**
   * Produits exceptionnels.
   */
  produit_exceptionnel?: number
  /**
   * Charges exceptionnelles.
   */
  charge_exceptionnelle?: number
  /**
   * Charges financières.
   */
  charges_financieres?: number
  /**
   * Chiffre d'affaires
   */
  ca?: number
  /**
   * Concours bancaires courants. (Pour recalculer les frais financiers court terme de la Banque de France)
   */
  concours_bancaire_courant?: number
  /**
   * Valeur ajoutée.
   */
  valeur_ajoutee?: number
  /**
   * Dette fiscale et sociale
   */
  dette_fiscale_et_sociale?: number
  /**
   * Raison sociale
   */
  nom_entreprise?: string
  /**
   * Numéro siren
   */
  numero_siren?: string
  /**
   * Statut juridique
   */
  statut_juridique?: string
  /**
   * Présence d'une procédure collective en cours
   */
  procedure_collective?: {
    [k: string]: unknown
  }
  /**
   * Effectif consolidé à l'entreprise
   */
  effectif_consolide?: number
  /**
   * Frais de Recherche et Développement
   */
  frais_de_RetD?: number
  /**
   * Concessions, brevets, et droits similaires
   */
  conces_brev_et_droits_sim?: number
  /**
   * Nombre d'établissements secondaires de l'entreprise, en plus du siège.
   */
  nombre_etab_secondaire?: number
  /**
   * Nombre de filiales de l'entreprise. Dans la base de données des liens capitalistiques, le concept de filiale ne fait aucune référence au pourcentage d’appartenance entre le parent et la fille. Dans ce sens, si l'entreprise A est enregistrée comme ayant des intérêts dans l'entreprise B avec un très petit, ou même un pourcentage de participation inconnu, l'entreprise B sera considérée filiale de l'entreprise A.
   */
  nombre_filiale?: number
  /**
   * Nombre d'entreprises dans le groupe (groupe défini par les liens capitalistique d'au moins 50,01%)
   */
  taille_compo_groupe?: number
  /**
   * Durée de l'exercice en mois.
   */
  nombre_mois?: number
  /**
   * Équilibre financier: Ressources durables / Emplois stables
   */
  equilibre_financier?: number
  /**
   * Indépendance financière (exprimé en %): Fonds propres / Ressources durables * 100
   */
  independance_financiere?: number
  /**
   * Endettement (exprimé en %): Dettes de caractère financier / Ressources durables * 100
   */
  endettement?: number
  /**
   * Autonomie financière Fonds propres / Total bilan * 100
   */
  autonomie_financiere?: number
  /**
   * Degré d'amortissement des immobilisations corporelles (exprimé en %): Amortissements des immobilisations corporelles / Immobilisation corporelles brutes * 100
   */
  degre_immo_corporelle?: number
  /**
   * Financement de l'actif circulant net: Fonds de roulement net global / Actif circulant net
   */
  financement_actif_circulant?: number
  /**
   * Liquidité générale: Actif circulant net / Dettes à court terme
   */
  liquidite_generale?: number
  /**
   * Liquidité réduite: Actif circulant net hors stocks / Dettes à court terme
   */
  liquidite_reduite?: number
  /**
   * Rotation des stocks (exprimé en jours): Stock / Chiffre d'affaires net * 360. Selon la nomenclature NAF Rév. 2 pour les secteurs d'activité 45, 46, 47, 95 (sauf 9511Z) ainsi que pour les codes d'activités 2319Z, 3831Z et 3832Z : Marchandises / (Achats de marchandises + Variation de stock) * 360
   */
  rotation_stocks?: number
  /**
   * Crédit clients (exprimé en jours): (Clients + Effets portés à l'escompte et non échus) / Chiffre d'affaires TTC * 360
   */
  credit_client?: number
  /**
   * Crédit fournisseurs (exprimé en jours): Fournisseurs / Achats TTC * 360
   */
  credit_fournisseur?: number
  /**
   * Chiffre d'affaire par effectif (exprimé en k€/emploi): Chiffre d'affaires net / Effectif * 1000
   */
  ca_par_effectif?: number
  /**
   * Taux d'intérêt financier (exprimé en %): Intérêts / Chiffre d'affaires net * 100
   */
  taux_interet_financier?: number
  /**
   * Intérêts sur chiffre d'affaire (exprimé en %): Total des charges financières / Chiffre d'affaires net * 100
   */
  taux_interet_sur_ca?: number
  /**
   * Endettement global (exprimé en jours): (Dettes + Effets portés à l'escompte et non échus) / Chiffre d'affaires net * 360
   */
  endettement_global?: number
  /**
   * Taux d'endettement (exprimé en %): Dettes de caractère financier / (Capitaux propres + autres fonds propres) * 100
   */
  taux_endettement?: number
  /**
   * Capacité de remboursement: Dettes de caractère financier / Capacité d'autofinancement avant répartition
   */
  capacite_remboursement?: number
  /**
   * Capacité d'autofinancement (exprimé en %): Capacité d'autofinancement avant répartition / (Chiffre d'affaires net + Subvention d'exploitation) * 100
   */
  capacite_autofinancement?: number
  /**
   * Couverture du chiffre d'affaire par le besoin en fonds de roulement (exprimé en jours): Besoins en fonds de roulement / Chiffre d'affaires net * 360
   */
  couverture_ca_besoin_fdr?: number
  /**
   * PoidsBFRExploitation Poids des besoins en fonds de roulement d'exploitation (exprimé en %): Besoins en fonds de roulement d'exploitation / Chiffre d'affaires net * 100
   */
  poids_bfr_exploitation?: number
  /**
   * Exportation Exportation (exprimé en %): (Chiffre d'affaires net - Chiffre d'affaires net en France) / Chiffre d'affaires net * 100
   */
  exportation?: number
  /**
   * Efficacité économique (exprimé en k€/emploi): Valeur ajoutée / Effectif * 1000
   */
  efficacite_economique?: number
  /**
   * Productivité du potentiel de production: Valeur ajoutée / Immobilisations corporelles et incorporelles brutes
   */
  productivite_potentiel_production?: number
  /**
   * Productivtié du capital financier: Valeur ajoutée / Actif circulant net + Effets portés à l'escompte et non échus
   */
  productivite_capital_financier?: number
  /**
   * Productivité du capital investi: Valeur ajoutée / Total de l'actif + Effets portés à l'escompte et non échus
   */
  productivite_capital_investi?: number
  /**
   * Taux d'investissement productif (exprimé en %): Immobilisations à valeur d'acquisition / Valeur ajoutée * 100
   */
  taux_d_investissement_productif?: number
  /**
   * Rentabilité économique (exprimé en %): Excédent brut d'exploitation / Chiffre d'affaires net + Subventions d'exploitation * 100
   */
  rentabilite_economique?: number
  /**
   * Performance (exprimé en %): Résultat courant avant impôt / Chiffre d'affaires net + Subventions d'exploitation * 100
   */
  performance?: number
  /**
   * Rendement brut des fonds propres (exprimé en %): Résultat courant avant impôt / Fonds propres nets * 100
   */
  rendement_brut_fonds_propres?: number
  /**
   * Rentabilité nette (exprimé en %): Bénéfice ou perte / Chiffre d'affaires net + Subventions d'exploitation * 100
   */
  rentabilite_nette?: number
  /**
   * Rendement des capitaux propres (exprimé en %): Bénéfice ou perte / Capitaux propres nets * 100
   */
  rendement_capitaux_propres?: number
  /**
   * RendementRessourcesDurables Rendement des ressources durables (exprimé en %): Résultat courant avant impôts + Intérêts et charges assimilées / Ressources durables nettes * 100
   */
  rendement_ressources_durables?: number
  /**
   * Taux de marge commerciale (exprimé en %): Marge commerciale / Vente de marchandises * 100
   */
  taux_marge_commerciale?: number
  /**
   * Taux de valeur ajoutée (exprimé en %): Valeur ajoutée / Chiffre d'affaires net * 100
   */
  taux_valeur_ajoutee?: number
  /**
   * Part des salariés (exprimé en %): (Charges de personnel + Participation des salariés aux résultats) / Valeur ajoutée * 100
   */
  part_salaries?: number
  /**
   * Part de l'État (exprimé en %): Impôts et taxes / Valeur ajoutée * 100
   */
  part_etat?: number
  /**
   * Part des prêteurs (exprimé en %): Intérêts / Valeur ajoutée * 100
   */
  part_preteur?: number
  /**
   * Part de l'autofinancement (exprimé en %): Capacité d'autofinancement avant répartition / Valeur ajoutée * 100
   */
  part_autofinancement?: number
  /**
   * Chiffre d'affaires à l'exportation
   */
  ca_exportation?: number
  /**
   * Achats de marchandises
   */
  achat_marchandises?: number
  /**
   * Achats de matières premières et autres approvisionnement.
   */
  achat_matieres_premieres?: number
  /**
   * Production de l'exercice.
   */
  production?: number
  /**
   * Marge commerciale.
   */
  marge_commerciale?: number
  /**
   * Consommation de l'exercice.
   */
  consommation?: number
  /**
   * Autres achats et charges externes.
   */
  autres_achats_charges_externes?: number
  /**
   * Charges de personnel.
   */
  charge_personnel?: number
  /**
   * Impôts, taxes et versements assimilés.
   */
  impots_taxes?: number
  /**
   * Subventions d'exploitation.
   */
  subventions_d_exploitation?: number
  /**
   * Autres produits, charges et reprises.
   */
  autres_produits_charges_reprises?: number
  /**
   * Dotation d'exploitation aux amortissements et aux provisions.
   */
  dotation_amortissement?: number
  /**
   * Résultat d'exploitation.
   */
  resultat_expl?: number
  /**
   * Opérations en commun.
   */
  operations_commun?: number
  /**
   * Résultat courant avant impôts.
   */
  resultat_avant_impot?: number
  /**
   * Participation des salariés aux résultats.
   */
  participation_salaries?: number
  /**
   * Impôts sur les bénéfices et impôts différés.
   */
  impot_benefice?: number
  /**
   * Bénéfice ou perte.
   */
  benefice_ou_perte?: number
}
/**
 * Champs importés par le parseur lib/urssaf/effectif_ent.go de sfdata.
 */
export interface EntréeEffectifEnt {
  periode: Date
  /**
   * Nombre de personnes employées par l'entreprise.
   */
  effectif: number
}
/**
 * Champs importés par le parseur lib/urssaf/effectif.go de sfdata.
 */
export interface EntréeEffectif {
  /**
   * Compte administratif URSSAF.
   */
  numero_compte: string
  periode: Date
  /**
   * Nombre de personnes employées par l'établissement.
   */
  effectif: number
}
/**
 * Champs importés par le parseur lib/ellisphere/main.go de sfdata.
 */
export interface EntréeEllisphere {
  code_groupe?: string
  siren_groupe?: string
  refid_groupe?: string
  raison_sociale_groupe?: string
  adresse_groupe?: string
  personne_pou_m_groupe?: string
  niveau_detention?: number
  part_financiere?: number
  code_filiere?: string
  refid_filiere?: string
  personne_pou_m_filiere?: string
}
/**
 * Champs importés par le parseur lib/paydex/main.go de sfdata.
 */
export interface EntréePaydex {
  date_valeur: Date
  nb_jours: number
}
/**
 * Champs importés par le parseur lib/urssaf/procol.go de sfdata.
 */
export interface EntréeDéfaillances {
  /**
   * Nature de la procédure de défaillance.
   */
  action_procol: "liquidation" | "redressement" | "sauvegarde"
  /**
   * Evénement survenu dans le cadre de cette procédure.
   */
  stade_procol:
    | "abandon_procedure"
    | "solde_procedure"
    | "fin_procedure"
    | "plan_continuation"
    | "ouverture"
    | "inclusion_autre_procedure"
    | "cloture_insuffisance_actif"
  /**
   * Date effet de la procédure collective.
   */
  date_effet: Date
}
/**
 * Champs importés par le parseur lib/sirene_ul/main.go de sfdata.
 */
export interface EntréeSireneEntreprise {
  /**
   * Numéro Siren de l'entreprise
   */
  siren?: string
  /**
   * Numéro interne de classement (Nic) de l’unité légale
   */
  nic?: string
  /**
   * Dénomination de l’unité légale
   */
  raison_sociale: string
  /**
   * Nom de naissance de la personne physique
   */
  nom_unite_legale?: string
  /**
   * Nom d’usage de la personne physique
   */
  nom_usage_unite_legale?: string
  /**
   * Premier prénom déclaré pour une personne physique
   */
  prenom1_unite_legale?: string
  /**
   * Deuxième prénom déclaré pour une personne physique
   */
  prenom2_unite_legale?: string
  /**
   * Troisième prénom déclaré pour une personne physique
   */
  prenom3_unite_legale?: string
  /**
   * Quatrième prénom déclaré pour une personne physique
   */
  prenom4_unite_legale?: string
  /**
   * Catégorie juridique de l'unité légale. Cf https://www.insee.fr/fr/information/2028129
   */
  statut_juridique: string
  /**
   * Date de création de l'unité légale
   */
  date_creation?: Date
}
/**
 * Champs importés par le parseur lib/sirene/main.go de sfdata.
 */
export interface EntréeSirene {
  /**
   * Numéro Siren de l'entreprise
   */
  siren?: string
  /**
   * Numéro interne de classement de l'établissement
   */
  nic?: string
  /**
   * Qualité de siège ou non de l’établissement
   */
  siege?: {
    [k: string]: unknown
  }
  /**
   * Complément d’adresse
   */
  complement_adresse?: string
  /**
   * Numéro de la voie de l’adresse
   */
  numero_voie?: string
  /**
   * Indice de répétition dans la voie
   */
  indrep?: string
  /**
   * Type de voie
   */
  type_voie?: string
  /**
   * Libellé de voie
   */
  voie?: string
  /**
   * Libellé de la commune
   */
  commune?: string
  /**
   * Libellé de la commune pour un établissement situé à l’étranger
   */
  commune_etranger?: string
  /**
   * Distribution spéciale de l’établissement
   */
  distribution_speciale?: string
  /**
   * Code commune de l’établissement
   */
  code_commune?: string
  /**
   * Code cedex
   */
  code_cedex?: string
  /**
   * Libellé du code cedex
   */
  cedex?: string
  /**
   * Code pays pour un établissement situé à l’étranger
   */
  code_pays_etranger?: string
  /**
   * Libellé du pays pour un établissement situé à l’étranger
   */
  pays_etranger?: string
  /**
   * Code postal
   */
  code_postal?: string
  /**
   * Code de département généré à partir du code postal (ex: 2A et 2B pour la Corse)
   */
  departement?: string
  /**
   * Activité principale de l'établissement pendant la période, dans le cas où celui-ci est renseigné selon la deuxième version de nomenclature NAF
   */
  ape?: string
  /**
   * Activité principale de l'établissement pendant la période, dans le cas où celui-ci est renseigné dans un format différent de la deuxième version de nomenclature NAF
   */
  code_activite?: string
  /**
   * Nomenclature NAF employée pour renseigner le code d'activité/APE de l'établissement (cf https://www.insee.fr/fr/information/2416409), si autre que deuxième révision
   */
  nomen_activite?: string
  /**
   * Date de création de l’établissement
   */
  date_creation?: Date
  /**
   * Géolocalisation des locaux: longitude
   */
  longitude?: number
  /**
   * Géolocalisation des locaux: latitude
   */
  latitude?: number
}
