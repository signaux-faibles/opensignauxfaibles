/* Generated by bundle-schemas.go - DO NOT EDIT */

package engine

var validationSchemas = map[string]string{
"apconso.schema.json": `{
  "title": "EntréeApConso",
  "description": "Champs importés par le parseur lib/apconso/main.go de sfdata.",
  "bsonType": "object",
  "required": ["id_conso", "periode", "heure_consomme", "montant", "effectif"],
  "properties": {
    "id_conso": {
      "description": "Numéro de la demande (11 caractères principalement des chiffres)",
      "bsonType": "string"
    },
    "heure_consomme": {
      "description": "Heures consommées (chômées) dans le mois",
      "bsonType": "number"
    },
    "montant": {
      "description": "Montants consommés dans le mois",
      "bsonType": "number"
    },
    "effectif": {
      "description": "Nombre de salariés en activité partielle dans le mois",
      "bsonType": "number"
    },
    "periode": {
      "description": "Mois considéré",
      "bsonType": "date"
    }
  },
  "additionalProperties": false
}
`,
"apdemande.schema.json": `{
  "title": "EntréeApDemande",
  "description": "Champs importés par le parseur lib/apdemande/main.go de sfdata.",
  "bsonType": "object",
  "required": [
    "id_demande",
    "periode",
    "hta",
    "motif_recours_se",
    "effectif_entreprise",
    "effectif",
    "date_statut",
    "mta",
    "effectif_autorise",
    "heure_consommee",
    "montant_consommee",
    "effectif_consomme",
    "perimetre"
  ],
  "properties": {
    "id_demande": {
      "description": "Numéro de la demande (11 caractères principalement des chiffres)",
      "bsonType": "string"
    },
    "periode": {
      "description": "Période de chômage",
      "bsonType": "object",
      "required": ["start", "end"],
      "properties": {
        "start": { "bsonType": "date" },
        "end": { "bsonType": "date" }
      },
      "additionalProperties": false
    },
    "hta": {
      "description": "Nombre total d'heures autorisées",
      "bsonType": "number"
    },
    "motif_recours_se": {
      "description": "Motif de recours à l'activité partielle: \n 1 - Conjoncture économique. \n 2 - Difficultés d’approvisionnement en matières premières ou en énergie \n 3 - Sinistre ou intempéries de caractère exceptionnel \n 4 - Transformation, restructuration ou modernisation des installations et des bâtiments \n 5 - Autres circonstances exceptionnelles",
      "bsonType": "number"
    },
    "effectif_entreprise": {
      "description": "Effectif de l'entreprise",
      "bsonType": "number"
    },
    "effectif": {
      "description": "Effectif de l'établissement",
      "bsonType": "number"
    },
    "date_statut": {
      "description": "Date du statut - création ou mise à jour de la demande",
      "bsonType": "date"
    },
    "mta": {
      "description": "Montant total autorisé",
      "bsonType": ["number", "null"]
    },
    "effectif_autorise": {
      "description": "Effectifs autorisés",
      "bsonType": "number"
    },
    "heure_consommee": {
      "description": "Nombre total d'heures consommées",
      "bsonType": ["number", "null"]
    },
    "montant_consommee": {
      "description": "Montant total consommé",
      "bsonType": ["number", "null"]
    },
    "effectif_consomme": {
      "description": "Effectifs consommés",
      "bsonType":  ["number", "null"]
    },
    "perimetre": {
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"bdf.schema.json": `{
  "title": "EntréeBdf",
  "description": "Champs importés par le parseur lib/bdf/main.go de sfdata.",
  "bsonType": "object",
  "required": [
    "siren",
    "annee_bdf",
    "arrete_bilan_bdf",
    "raison_sociale",
    "secteur",
    "poids_frng",
    "taux_marge",
    "delai_fournisseur",
    "dette_fiscale",
    "financier_court_terme",
    "frais_financier"
  ],
  "properties": {
    "arrete_bilan_bdf": {
      "description": "Date de clôture de l'exercice.",
      "bsonType": "date"
    },
    "annee_bdf": {
      "description": "Année de l'exercice.",
      "bsonType": "number"
    },
    "raison_sociale": {
      "description": "Raison sociale de l'entreprise.",
      "bsonType": "string"
    },
    "secteur": {
      "description": "Secteur d'activité.",
      "bsonType": "string"
    },
    "siren": {
      "description": "Siren de l'entreprise.",
      "bsonType": "string",
      "pattern": "^[0-9]{9}$"
    },
    "poids_frng": {
      "description": "Poids du fonds de roulement net global sur le chiffre d'affaire. Exprimé en %.",
      "bsonType": "number"
    },
    "taux_marge": {
      "description": "Taux de marge, rapport de l'excédent brut d'exploitation (EBE) sur la valeur ajoutée (exprimé en %): 100*EBE / valeur ajoutee",
      "bsonType": "number"
    },
    "delai_fournisseur": {
      "description": "Délai estimé de paiement des fournisseurs (exprimé en jours): 360 * dettes fournisseurs / achats HT",
      "bsonType": "number"
    },
    "dette_fiscale": {
      "description": "Poids des dettes fiscales et sociales, par rapport à la valeur ajoutée (exprimé en %): 100 * dettes fiscales et sociales / Valeur ajoutee",
      "bsonType": "number"
    },
    "financier_court_terme": {
      "description": "Poids du financement court terme (exprimé en %): 100 * concours bancaires courants / chiffre d'affaires HT",
      "bsonType": "number"
    },
    "frais_financier": {
      "description": "Poids des frais financiers, sur l'excedent brut d'exploitation corrigé des produits et charges hors exploitation (exprimé en %): 100 * frais financiers / (EBE + Produits hors expl. - charges hors expl.)",
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"ccsf.schema.json": `{
  "title": "EntréeCcsf",
  "description": "Champs importés par le parseur lib/urssaf/ccsf.go de sfdata.",
  "bsonType": "object",
  "required": ["date_traitement", "stade", "action"],
  "properties": {
    "date_traitement": {
      "bsonType": "date",
      "description": "Date de début de la procédure CCSF"
    },
    "stade": {
      "bsonType": "string",
      "description": "Stade de la demande de délai",
      "enum": ["DEBUT", "REFUS", "APPROB", "FIN", "REJ"]
    },
    "action": {
      "bsonType": "string",
      "description": "Code externe de l'action",
      "enum": ["CCSF"]
    }
  },
  "additionalProperties": false
}
`,
"compte.schema.json": `{
  "title": "EntréeCompte",
  "description": "Champs importés par le parseur lib/urssaf/compte.go de sfdata.",
  "bsonType": "object",
  "required": ["periode", "siret", "numero_compte"],
  "properties": {
    "periode": {
      "bsonType": "date",
      "description": "Date à laquelle cet établissement est associé à ce numéro de compte URSSAF."
    },
    "siret": {
      "bsonType": "string",
      "description": "Numéro SIRET de l'établissement. Les numéros avec des Lettres sont des sirets provisoires."
    },
    "numero_compte": {
      "bsonType": "string",
      "description": "Compte administratif URSSAF."
    }
  },
  "additionalProperties": false
}
`,
"cotisation.schema.json": `{
  "title": "EntréeCotisation",
  "description": "Champs importés par le parseur lib/urssaf/cotisation.go de sfdata.",
  "bsonType": "object",
  "required": ["numero_compte", "periode", "encaisse", "du"],
  "properties": {
    "numero_compte": {
      "description": "Compte administratif URSSAF.",
      "bsonType": "string"
    },
    "periode": {
      "description": "Période sur laquelle le montants s'appliquent.",
      "bsonType": "object",
      "required": ["start", "end"],
      "properties": {
        "start": { "bsonType": "date" },
        "end": { "bsonType": "date" }
      },
      "additionalProperties": false
    },
    "encaisse": {
      "description": "Cotisation encaissée directement, en euros.",
      "bsonType": "number"
    },
    "du": {
      "description": "Cotisation due, en euros. À utiliser pour calculer le montant moyen mensuel du: Somme cotisations dues / nb périodes.",
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"debit.schema.json": `{
  "title": "EntréeDebit",
  "description": "Champs importés par le parseur lib/sirene_ul/main.go de sfdata. Représente un reste à payer (dette) sur cotisation sociale ou autre.",
  "bsonType": "object",
  "required": [
    "numero_compte",
    "numero_ecart_negatif",
    "date_traitement",
    "part_ouvriere",
    "part_patronale",
    "numero_historique",
    "etat_compte",
    "code_procedure_collective",
    "periode",
    "code_operation_ecart_negatif",
    "code_motif_ecart_negatif",
    "recours_en_cours"
  ],
  "properties": {
    "numero_compte": {
      "description": "Identifiant URSSAF d'établissement (équivalent du SIRET).",
      "bsonType": "string"
    },
    "periode": {
      "description": "Période sur laquelle le montants s'appliquent.",
      "bsonType": "object",
      "required": ["start", "end"],
      "properties": {
        "start": { "bsonType": "date" },
        "end": { "bsonType": "date" }
      },
      "additionalProperties": false
    },
    "numero_ecart_negatif": {
      "description": "L'écart négatif (ecn) correspond à une période en débit. Pour une même période, plusieurs débits peuvent être créés. On leur attribue un numéro d'ordre. Par exemple, 101, 201, 301 etc.; ou 101, 102, 201 etc. correspondent respectivement au 1er, 2ème et 3ème ecn de la période considérée.",
      "bsonType": "string"
    },
    "numero_historique": {
      "description": "Ordre des opérations pour un écart négatif donné.",
      "bsonType": "number"
    },
    "date_traitement": {
      "description": "Date de constatation du débit (exemple: remboursement, majoration ou autre modification du montant)",
      "bsonType": "date"
    },
    "debit_suivant": {
      "description": "Hash du débit suivant. (généré par sfdata)",
      "bsonType": "string"
    },
    "part_ouvriere": {
      "description": "Montant des débits sur la part ouvrières, exprimées en euros (€). Sont exclues les pénalités et les majorations de retard.",
      "bsonType": "number"
    },
    "part_patronale": {
      "description": "Montant des débits sur la part patronale, exprimées en euros (€). Sont exclues les pénalités et les majorations de retard.",
      "bsonType": "number"
    },
    "etat_compte": {
      "description": "Code état du compte: 1 (Actif), 2 (Suspendu) ou 3 (Radié).",
      "bsonType": "number"
    },
    "code_procedure_collective": {
      "description": "Code qui indique si le compte fait l'objet d'une procédure collective: 1 (en cours), 2 (plan de redressement en cours), 9 (procedure collective sans dette à l'Urssaf) ou valeur nulle en cas d'absence de procédure collective.",
      "bsonType": "string"
    },
    "code_operation_ecart_negatif": {
      "description": "Code opération historique de l'écart négatif: \n 1 Mise en recouvrement \n 2 Paiement \n 3 Admission en non valeur \n 4 Remise de majoration de retard \n 5 Abandon de solde debiteur \n 11 Annulation de mise en recouvrement \n 12 Annulation paiement \n 13 Annulation a-n-v \n 14 Annulation de remise de majoration retard \n 15 Annulation abandon solde debiteur",
      "bsonType": "string"
    },
    "code_motif_ecart_negatif": {
      "description": "Code motif de l'écart négatif: \n 0 Cde motif inconnu \n 1 Retard dans le versement \n 2 Absence ou insuffisance de versement \n 3 Taxation provisionelle. Déclarations non fournies \n 4 Majorations de retard complémentaires Article R243-18 du code de la sécurité sociale \n 5 Contrôle,chefs de redressement notifiés le JJ/MM/AA Article R243-59 de la Securité Sociale \n 6 Fourniture tardive des déclarations \n 7 Bases déclarées supérieures à Taxation provisionnelle \n 8 Retard dans le versement et fourniture tardive des déclarations \n 9 Absence ou insuffisance de versement et fourniture tardive des déclarations \n 10 Rappel sur contrôle et fourniture tardive des déclarations \n 11 Régularisation d'une taxation provisionnelle \n 12 Régularisation annuelle \n 13 Rejet du titre de paiement par la banque . \n 14 Modification d'affectation d'un crédit \n 15 Annulation d'un crédit \n 16 Régularisation suite à modification du Taux Accident du Travail \n 17 Régularisation suite à assujettissement au transport (origine débit sur PJ=4) \n 18 Majorations pour non respect de paiement par moyen dématérialisé Article L243-14 \n 19 Rapprochement TR/BRC sous réserve de vérification ultérieure \n 20 Cotisations complémentaires suite modification des revenus déclarés \n 21 Cotisations complémentaires suite à non fourniture du contrat d'exonération \n 22 Contrôle. Chefs de redressement notifiés le JJ/MM/AA. Article L324.9 du code du travail \n 23 Cotisations complémentaires suite conditions d'exonération non remplies \n 24 Absence de versement \n 25 Insuffisance de versement \n 26 Absence de versement et fourniture tardive des déclarations \n 27 Insuffisance de versement et fourniture tardive des déclarations",
      "bsonType": "string"
    },
    "recours_en_cours": {
      "description": "Recours en cours.",
      "bsonType": "bool"
    }
  },
  "additionalProperties": false
}
`,
"delai.schema.json": `{
  "title": "EntréeDelai",
  "description": "Champs importés par le parseur lib/urssaf/delai.go de sfdata.",
  "bsonType": "object",
  "required": [
    "numero_compte",
    "numero_contentieux",
    "date_creation",
    "date_echeance",
    "duree_delai",
    "denomination",
    "indic_6m",
    "annee_creation",
    "montant_echeancier",
    "stade",
    "action"
  ],
  "properties": {
    "numero_compte": {
      "bsonType": "string",
      "description": "Compte administratif URSSAF."
    },
    "numero_contentieux": {
      "bsonType": "string",
      "description": "Le numéro de structure est l'identifiant d'un dossier contentieux."
    },
    "date_creation": {
      "bsonType": "date",
      "description": "Date de création du délai."
    },
    "date_echeance": {
      "bsonType": "date",
      "description": "Date d'échéance du délai."
    },
    "duree_delai": {
      "bsonType": "number",
      "description": "Durée du délai en jours: nombre de jours entre date_creation et date_echeance.",
      "minimum": 1
    },
    "denomination": {
      "bsonType": "string",
      "description": "Raison sociale de l'établissement."
    },
    "indic_6m": {
      "bsonType": "string",
      "description": "Délai inférieur ou supérieur à 6 mois ? Modalités INF et SUP."
    },
    "annee_creation": {
      "bsonType": "number",
      "description": "Année de création du délai."
    },
    "montant_echeancier": {
      "bsonType": "number",
      "description": "Montant global de l'échéancier, en euros.",
      "minimum": 0.01
    },
    "stade": {
      "bsonType": "string",
      "description": "Code externe du stade."
    },
    "action": {
      "bsonType": "string",
      "description": "Code externe de l'action."
    }
  },
  "additionalProperties": false
}
`,
"detect_invalid_entries.pipeline.json": `[
  { "$project": { "_id": 1, "batches": { "$objectToArray": "$value.batch" } } },
  { "$unwind": { "path": "$batches", "preserveNullAndEmptyArrays": false } },
  {
    "$project": {
      "_id": 1,
      "batchKey": "$batches.k",
      "dataPerHash": { "$objectToArray": "$batches.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": "$dataPerHash.k",
      "dataPerHash": "$dataPerHash.v"
    }
  },
  { "$match": { "dataPerHash": { "$not": { "$type": 3 } } } }
]
`,
"diane.schema.json": `{
  "title": "EntréeDiane",
  "description": "Champs importés par le parseur lib/diane/main.go de sfdata.",
  "bsonType": "object",
  "properties": {
    "exercice_diane": {
      "description": "Année de l'exercice",
      "bsonType": "number"
    },
    "arrete_bilan_diane": {
      "description": "Date d'arrêté du bilan",
      "bsonType": "date"
    },
    "couverture_ca_fdr": {
      "description": "Couverture du chiffre d'affaire par le fonds de roulement (exprimé en jours): Fonds de roulement net global / Chiffre d'affaires net * 360",
      "bsonType": "number"
    },
    "interets": {
      "description": "Intérêts et charges assimilées.",
      "bsonType": "number"
    },
    "excedent_brut_d_exploitation": {
      "description": "Excédent brut d'exploitation.",
      "bsonType": "number"
    },
    "produits_financiers": {
      "description": "Produits financiers.",
      "bsonType": "number"
    },
    "produit_exceptionnel": {
      "description": "Produits exceptionnels.",
      "bsonType": "number"
    },
    "charge_exceptionnelle": {
      "description": "Charges exceptionnelles.",
      "bsonType": "number"
    },
    "charges_financieres": {
      "description": "Charges financières.",
      "bsonType": "number"
    },
    "ca": {
      "description": "Chiffre d'affaires",
      "bsonType": "number"
    },
    "concours_bancaire_courant": {
      "description": "Concours bancaires courants. (Pour recalculer les frais financiers court terme de la Banque de France)",
      "bsonType": "number"
    },
    "valeur_ajoutee": {
      "description": "Valeur ajoutée.",
      "bsonType": "number"
    },
    "dette_fiscale_et_sociale": {
      "description": "Dette fiscale et sociale",
      "bsonType": "number"
    },
    "nom_entreprise": {
      "description": "Raison sociale",
      "bsonType": "string"
    },
    "numero_siren": {
      "description": "Numéro siren",
      "bsonType": "string"
    },
    "statut_juridique": {
      "description": "Statut juridique",
      "bsonType": "string"
    },
    "procedure_collective": {
      "description": "Présence d'une procédure collective en cours",
      "bsonType": "bool"
    },
    "effectif_consolide": {
      "description": "Effectif consolidé à l'entreprise",
      "bsonType": "number"
    },
    "frais_de_RetD": {
      "description": "Frais de Recherche et Développement",
      "bsonType": "number"
    },
    "conces_brev_et_droits_sim": {
      "description": "Concessions, brevets, et droits similaires",
      "bsonType": "number"
    },
    "nombre_etab_secondaire": {
      "description": "Nombre d'établissements secondaires de l'entreprise, en plus du siège.",
      "bsonType": "number"
    },
    "nombre_filiale": {
      "description": "Nombre de filiales de l'entreprise. Dans la base de données des liens capitalistiques, le concept de filiale ne fait aucune référence au pourcentage d’appartenance entre le parent et la fille. Dans ce sens, si l'entreprise A est enregistrée comme ayant des intérêts dans l'entreprise B avec un très petit, ou même un pourcentage de participation inconnu, l'entreprise B sera considérée filiale de l'entreprise A.",
      "bsonType": "number"
    },
    "taille_compo_groupe": {
      "description": "Nombre d'entreprises dans le groupe (groupe défini par les liens capitalistique d'au moins 50,01%)",
      "bsonType": "number"
    },
    "nombre_mois": {
      "description": "Durée de l'exercice en mois.",
      "bsonType": "number"
    },
    "equilibre_financier": {
      "description": "Équilibre financier: Ressources durables / Emplois stables",
      "bsonType": "number"
    },
    "independance_financiere": {
      "description": "Indépendance financière (exprimé en %): Fonds propres / Ressources durables * 100",
      "bsonType": "number"
    },
    "endettement": {
      "description": "Endettement (exprimé en %): Dettes de caractère financier / Ressources durables * 100",
      "bsonType": "number"
    },
    "autonomie_financiere": {
      "description": "Autonomie financière Fonds propres / Total bilan * 100",
      "bsonType": "number"
    },
    "degre_immo_corporelle": {
      "description": "Degré d'amortissement des immobilisations corporelles (exprimé en %): Amortissements des immobilisations corporelles / Immobilisation corporelles brutes * 100",
      "bsonType": "number"
    },
    "financement_actif_circulant": {
      "description": "Financement de l'actif circulant net: Fonds de roulement net global / Actif circulant net",
      "bsonType": "number"
    },
    "liquidite_generale": {
      "description": "Liquidité générale: Actif circulant net / Dettes à court terme",
      "bsonType": "number"
    },
    "liquidite_reduite": {
      "description": "Liquidité réduite: Actif circulant net hors stocks / Dettes à court terme",
      "bsonType": "number"
    },
    "rotation_stocks": {
      "description": "Rotation des stocks (exprimé en jours): Stock / Chiffre d'affaires net * 360. Selon la nomenclature NAF Rév. 2 pour les secteurs d'activité 45, 46, 47, 95 (sauf 9511Z) ainsi que pour les codes d'activités 2319Z, 3831Z et 3832Z : Marchandises / (Achats de marchandises + Variation de stock) * 360",
      "bsonType": "number"
    },
    "credit_client": {
      "description": "Crédit clients (exprimé en jours): (Clients + Effets portés à l'escompte et non échus) / Chiffre d'affaires TTC * 360",
      "bsonType": "number"
    },
    "credit_fournisseur": {
      "description": "Crédit fournisseurs (exprimé en jours): Fournisseurs / Achats TTC * 360",
      "bsonType": "number"
    },
    "ca_par_effectif": {
      "description": "Chiffre d'affaire par effectif (exprimé en k€/emploi): Chiffre d'affaires net / Effectif * 1000",
      "bsonType": "number"
    },
    "taux_interet_financier": {
      "description": "Taux d'intérêt financier (exprimé en %): Intérêts / Chiffre d'affaires net * 100",
      "bsonType": "number"
    },
    "taux_interet_sur_ca": {
      "description": "Intérêts sur chiffre d'affaire (exprimé en %): Total des charges financières / Chiffre d'affaires net * 100",
      "bsonType": "number"
    },
    "endettement_global": {
      "description": "Endettement global (exprimé en jours): (Dettes + Effets portés à l'escompte et non échus) / Chiffre d'affaires net * 360",
      "bsonType": "number"
    },
    "taux_endettement": {
      "description": "Taux d'endettement (exprimé en %): Dettes de caractère financier / (Capitaux propres + autres fonds propres) * 100",
      "bsonType": "number"
    },
    "capacite_remboursement": {
      "description": "Capacité de remboursement: Dettes de caractère financier / Capacité d'autofinancement avant répartition",
      "bsonType": "number"
    },
    "capacite_autofinancement": {
      "description": "Capacité d'autofinancement (exprimé en %): Capacité d'autofinancement avant répartition / (Chiffre d'affaires net + Subvention d'exploitation) * 100",
      "bsonType": "number"
    },
    "couverture_ca_besoin_fdr": {
      "description": "Couverture du chiffre d'affaire par le besoin en fonds de roulement (exprimé en jours): Besoins en fonds de roulement / Chiffre d'affaires net * 360",
      "bsonType": "number"
    },
    "poids_bfr_exploitation": {
      "description": "PoidsBFRExploitation Poids des besoins en fonds de roulement d'exploitation (exprimé en %): Besoins en fonds de roulement d'exploitation / Chiffre d'affaires net * 100",
      "bsonType": "number"
    },
    "exportation": {
      "description": "Exportation Exportation (exprimé en %): (Chiffre d'affaires net - Chiffre d'affaires net en France) / Chiffre d'affaires net * 100",
      "bsonType": "number"
    },
    "efficacite_economique": {
      "description": "Efficacité économique (exprimé en k€/emploi): Valeur ajoutée / Effectif * 1000",
      "bsonType": "number"
    },
    "productivite_potentiel_production": {
      "description": "Productivité du potentiel de production: Valeur ajoutée / Immobilisations corporelles et incorporelles brutes",
      "bsonType": "number"
    },
    "productivite_capital_financier": {
      "description": "Productivtié du capital financier: Valeur ajoutée / Actif circulant net + Effets portés à l'escompte et non échus",
      "bsonType": "number"
    },
    "productivite_capital_investi": {
      "description": "Productivité du capital investi: Valeur ajoutée / Total de l'actif + Effets portés à l'escompte et non échus",
      "bsonType": "number"
    },
    "taux_d_investissement_productif": {
      "description": "Taux d'investissement productif (exprimé en %): Immobilisations à valeur d'acquisition / Valeur ajoutée * 100",
      "bsonType": "number"
    },
    "rentabilite_economique": {
      "description": "Rentabilité économique (exprimé en %): Excédent brut d'exploitation / Chiffre d'affaires net + Subventions d'exploitation * 100",
      "bsonType": "number"
    },
    "performance": {
      "description": "Performance (exprimé en %): Résultat courant avant impôt / Chiffre d'affaires net + Subventions d'exploitation * 100",
      "bsonType": "number"
    },
    "rendement_brut_fonds_propres": {
      "description": "Rendement brut des fonds propres (exprimé en %): Résultat courant avant impôt / Fonds propres nets * 100",
      "bsonType": "number"
    },
    "rentabilite_nette": {
      "description": "Rentabilité nette (exprimé en %): Bénéfice ou perte / Chiffre d'affaires net + Subventions d'exploitation * 100",
      "bsonType": "number"
    },
    "rendement_capitaux_propres": {
      "description": "Rendement des capitaux propres (exprimé en %): Bénéfice ou perte / Capitaux propres nets * 100",
      "bsonType": "number"
    },
    "rendement_ressources_durables": {
      "description": "RendementRessourcesDurables Rendement des ressources durables (exprimé en %): Résultat courant avant impôts + Intérêts et charges assimilées / Ressources durables nettes * 100",
      "bsonType": "number"
    },
    "taux_marge_commerciale": {
      "description": "Taux de marge commerciale (exprimé en %): Marge commerciale / Vente de marchandises * 100",
      "bsonType": "number"
    },
    "taux_valeur_ajoutee": {
      "description": "Taux de valeur ajoutée (exprimé en %): Valeur ajoutée / Chiffre d'affaires net * 100",
      "bsonType": "number"
    },
    "part_salaries": {
      "description": "Part des salariés (exprimé en %): (Charges de personnel + Participation des salariés aux résultats) / Valeur ajoutée * 100",
      "bsonType": "number"
    },
    "part_etat": {
      "description": "Part de l'État (exprimé en %): Impôts et taxes / Valeur ajoutée * 100",
      "bsonType": "number"
    },
    "part_preteur": {
      "description": "Part des prêteurs (exprimé en %): Intérêts / Valeur ajoutée * 100",
      "bsonType": "number"
    },
    "part_autofinancement": {
      "description": "Part de l'autofinancement (exprimé en %): Capacité d'autofinancement avant répartition / Valeur ajoutée * 100",
      "bsonType": "number"
    },
    "ca_exportation": {
      "description": "Chiffre d'affaires à l'exportation",
      "bsonType": "number"
    },
    "achat_marchandises": {
      "description": "Achats de marchandises",
      "bsonType": "number"
    },
    "achat_matieres_premieres": {
      "description": "Achats de matières premières et autres approvisionnement.",
      "bsonType": "number"
    },
    "production": {
      "description": "Production de l'exercice.",
      "bsonType": "number"
    },
    "marge_commerciale": {
      "description": "Marge commerciale.",
      "bsonType": "number"
    },
    "consommation": {
      "description": "Consommation de l'exercice.",
      "bsonType": "number"
    },
    "autres_achats_charges_externes": {
      "description": "Autres achats et charges externes.",
      "bsonType": "number"
    },
    "charge_personnel": {
      "description": "Charges de personnel.",
      "bsonType": "number"
    },
    "impots_taxes": {
      "description": "Impôts, taxes et versements assimilés.",
      "bsonType": "number"
    },
    "subventions_d_exploitation": {
      "description": "Subventions d'exploitation.",
      "bsonType": "number"
    },
    "autres_produits_charges_reprises": {
      "description": "Autres produits, charges et reprises.",
      "bsonType": "number"
    },
    "dotation_amortissement": {
      "description": "Dotation d'exploitation aux amortissements et aux provisions.",
      "bsonType": "number"
    },
    "resultat_expl": {
      "description": "Résultat d'exploitation.",
      "bsonType": "number"
    },
    "operations_commun": {
      "description": "Opérations en commun.",
      "bsonType": "number"
    },
    "resultat_avant_impot": {
      "description": "Résultat courant avant impôts.",
      "bsonType": "number"
    },
    "participation_salaries": {
      "description": "Participation des salariés aux résultats.",
      "bsonType": "number"
    },
    "impot_benefice": {
      "description": "Impôts sur les bénéfices et impôts différés.",
      "bsonType": "number"
    },
    "benefice_ou_perte": {
      "description": "Bénéfice ou perte.",
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"effectif.schema.json": `{
  "title": "EntréeEffectif",
  "description": "Champs importés par le parseur lib/urssaf/effectif.go de sfdata.",
  "bsonType": "object",
  "required": ["numero_compte", "periode", "effectif"],
  "properties": {
    "numero_compte": {
      "description": "Compte administratif URSSAF.",
      "bsonType": "string"
    },
    "periode": {
      "bsonType": "date"
    },
    "effectif": {
      "description": "Nombre de personnes employées par l'établissement.",
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"effectif_ent.schema.json": `{
  "title": "EntréeEffectifEnt",
  "description": "Champs importés par le parseur lib/urssaf/effectif_ent.go de sfdata.",
  "bsonType": "object",
  "required": ["periode", "effectif"],
  "properties": {
    "periode": {
      "bsonType": "date"
    },
    "effectif": {
      "description": "Nombre de personnes employées par l'entreprise.",
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"ellisphere.schema.json": `{
  "title": "EntréeEllisphere",
  "description": "Champs importés par le parseur lib/ellisphere/main.go de sfdata.",
  "bsonType": "object",
  "properties": {
    "code_groupe": {
      "description": "Code du groupe/actionnaire.",
      "bsonType": "string"
    },
    "siren_groupe": {
      "description": "Siren du groupe/actionnaire.",
      "bsonType": "string",
      "pattern": "^[0-9]{9}$"
    },
    "refid_groupe": {
      "description": "Référence du groupe/actionnaire.",
      "bsonType": "string"
    },
    "raison_sociale_groupe": {
      "description": "Raison sociale du groupe/actionnaire.",
      "bsonType": "string"
    },
    "adresse_groupe": {
      "description": "Adresse du groupe/actionnaire.",
      "bsonType": "string"
    },
    "personne_pou_m_groupe": {
      "description": "Groupe/actionnaire: personne physique (P) ou morale (M)",
      "bsonType": "string",
      "enum": ["P", "M"]
    },
    "niveau_detention": {
      "description": "Le Rang exprime le nombre d’intermédiaires entre 2 entités. (c.a.d. entre l'actionnaire et la filiale)",
      "bsonType": "number"
    },
    "part_financiere": {
      "description": "Le Pourcentage d’intérêt exprime la part mathématique du capital de la société détenue directement ou indirectement par l’entité mère.",
      "bsonType": "number"
    },
    "code_filiere": {
      "description": "Code de la filiale.",
      "bsonType": "string"
    },
    "refid_filiere": {
      "description": "Référence de la filiale.",
      "bsonType": "string"
    },
    "personne_pou_m_filiere": {
      "description": "Filiale: personne physique (P) ou morale (M)",
      "bsonType": "string",
      "enum": ["P", "M"]
    }
  },
  "additionalProperties": false
}
`,
"flatten_data_entries.pipeline.json": `[
  { "$project": { "_id": 1, "batches": { "$objectToArray": "$value.batch" } } },
  { "$unwind": { "path": "$batches", "preserveNullAndEmptyArrays": false } },
  {
    "$project": {
      "_id": 1,
      "batchKey": "$batches.k",
      "dataPerHash": { "$objectToArray": "$batches.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": "$dataPerHash.k",
      "dataPerHash": { "$objectToArray": "$dataPerHash.v" }
    }
  },
  {
    "$unwind": { "path": "$dataPerHash", "preserveNullAndEmptyArrays": false }
  },
  {
    "$project": {
      "_id": 1,
      "batchKey": 1,
      "dataType": 1,
      "dataHash": "$dataPerHash.k",
      "dataObject": "$dataPerHash.v"
    }
  }
]
`,
"paydex.schema.json": `{
  "title": "EntréePaydex",
  "description": "Champs importés par le parseur lib/paydex/main.go de sfdata.",
  "bsonType": "object",
  "required": ["date_valeur", "nb_jours"],
  "properties": {
    "date_valeur": {
      "bsonType": "date"
    },
    "nb_jours": {
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"procol.schema.json": `{
  "title": "EntréeDéfaillances",
  "description": "Champs importés par le parseur lib/urssaf/procol.go de sfdata.",
  "bsonType": "object",
  "required": ["action_procol", "stade_procol", "date_effet"],
  "properties": {
    "action_procol": {
      "description": "Nature de la procédure de défaillance.",
      "bsonType": "string",
      "enum": ["liquidation", "redressement", "sauvegarde"]
    },
    "stade_procol": {
      "description": "Evénement survenu dans le cadre de cette procédure.",
      "bsonType": "string",
      "enum": [
        "abandon_procedure",
        "solde_procedure",
        "fin_procedure",
        "plan_continuation",
        "ouverture",
        "inclusion_autre_procedure",
        "cloture_insuffisance_actif"
      ]
    },
    "date_effet": {
      "bsonType": "date",
      "description": "Date effet de la procédure collective."
    }
  },
  "additionalProperties": false
}
`,
"sirene.schema.json": `{
  "title": "EntréeSirene",
  "description": "Champs importés par le parseur lib/sirene/main.go de sfdata.",
  "bsonType": "object",
  "properties": {
    "siren": {
      "description": "Numéro Siren de l'entreprise",
      "bsonType": "string",
      "pattern": "^[0-9]{9}$"
    },
    "nic": {
      "description": "Numéro interne de classement de l'établissement",
      "bsonType": "string"
    },
    "siege": {
      "description": "Qualité de siège ou non de l’établissement",
      "bsonType": "bool"
    },
    "complement_adresse": {
      "description": "Complément d’adresse",
      "bsonType": "string"
    },
    "numero_voie": {
      "description": "Numéro de la voie de l’adresse",
      "bsonType": "string"
    },
    "indrep": {
      "description": "Indice de répétition dans la voie",
      "bsonType": "string"
    },
    "type_voie": {
      "description": "Type de voie",
      "bsonType": "string"
    },
    "voie": {
      "description": "Libellé de voie",
      "bsonType": "string"
    },
    "commune": {
      "description": "Libellé de la commune",
      "bsonType": "string"
    },
    "commune_etranger": {
      "description": "Libellé de la commune pour un établissement situé à l’étranger",
      "bsonType": "string"
    },
    "distribution_speciale": {
      "description": "Distribution spéciale de l’établissement",
      "bsonType": "string"
    },
    "code_commune": {
      "description": "Code commune de l’établissement",
      "bsonType": "string"
    },
    "code_cedex": {
      "description": "Code cedex",
      "bsonType": "string"
    },
    "cedex": {
      "description": "Libellé du code cedex",
      "bsonType": "string"
    },
    "code_pays_etranger": {
      "description": "Code pays pour un établissement situé à l’étranger",
      "bsonType": "string"
    },
    "pays_etranger": {
      "description": "Libellé du pays pour un établissement situé à l’étranger",
      "bsonType": "string"
    },
    "code_postal": {
      "description": "Code postal",
      "bsonType": "string"
    },
    "departement": {
      "description": "Code de département généré à partir du code postal (ex: 2A et 2B pour la Corse)",
      "bsonType": "string"
    },
    "ape": {
      "description": "Activité principale de l'établissement pendant la période, dans le cas où celui-ci est renseigné selon la deuxième version de nomenclature NAF",
      "bsonType": "string"
    },
    "code_activite": {
      "description": "Activité principale de l'établissement pendant la période, dans le cas où celui-ci est renseigné dans un format différent de la deuxième version de nomenclature NAF",
      "bsonType": "string"
    },
    "nomen_activite": {
      "description": "Nomenclature NAF employée pour renseigner le code d'activité/APE de l'établissement (cf https://www.insee.fr/fr/information/2416409), si autre que deuxième révision",
      "bsonType": "string"
    },
    "date_creation": {
      "description": "Date de création de l’établissement",
      "bsonType": "date"
    },
    "longitude": {
      "description": "Géolocalisation des locaux: longitude",
      "bsonType": "number"
    },
    "latitude": {
      "description": "Géolocalisation des locaux: latitude",
      "bsonType": "number"
    }
  },
  "additionalProperties": false
}
`,
"sirene_ul.schema.json": `{
  "title": "EntréeSireneEntreprise",
  "description": "Champs importés par le parseur lib/sirene_ul/main.go de sfdata.",
  "bsonType": "object",
  "required": ["raison_sociale", "statut_juridique"],
  "properties": {
    "siren": {
      "description": "Numéro Siren de l'entreprise",
      "bsonType": "string",
      "pattern": "^[0-9]{9}$"
    },
    "nic": {
      "description": "Numéro interne de classement (Nic) de l’unité légale",
      "bsonType": "string",
      "pattern": "^[0-9]{9}$"
    },
    "raison_sociale": {
      "description": "Dénomination de l’unité légale",
      "bsonType": "string"
    },
    "nom_unite_legale": {
      "description": "Nom de naissance de la personne physique",
      "bsonType": "string"
    },
    "nom_usage_unite_legale": {
      "description": "Nom d’usage de la personne physique",
      "bsonType": "string"
    },
    "prenom1_unite_legale": {
      "description": "Premier prénom déclaré pour une personne physique",
      "bsonType": "string"
    },
    "prenom2_unite_legale": {
      "description": "Deuxième prénom déclaré pour une personne physique",
      "bsonType": "string"
    },
    "prenom3_unite_legale": {
      "description": "Troisième prénom déclaré pour une personne physique",
      "bsonType": "string"
    },
    "prenom4_unite_legale": {
      "description": "Quatrième prénom déclaré pour une personne physique",
      "bsonType": "string"
    },
    "statut_juridique": {
      "description": "Catégorie juridique de l'unité légale. Cf https://www.insee.fr/fr/information/2028129",
      "bsonType": "string",
      "pattern": "^[0-9]{4}$"
    },
    "date_creation": {
      "description": "Date de création de l'unité légale",
      "bsonType": "date"
    }
  },
  "additionalProperties": false
}
`,
}
