# Description des données brutes

- [Préambule](#pr%C3%A9ambule)
- [Périmètre des données](#p%C3%A9rim%C3%A8tre-des-donn%C3%A9es)
  - [Entreprises :house:](#entreprises-house)
  - [Périmètre temporel :clock1:](#p%C3%A9rim%C3%A8tre-temporel-clock1)
- [Données importées](#donn%C3%A9es-import%C3%A9es)
  - [Base Sirene](#base-sirene)
  - [Données d'activité partielle](#donn%C3%A9es-dactivit%C3%A9-partielle)
    - [Données de demandes d'activité partielle](#donn%C3%A9es-de-demandes-dactivit%C3%A9-partielle)
      - [Table des motifs de recours à l'activité partielle](#table-des-motifs-de-recours-%C3%A0-lactivit%C3%A9-partielle)
      - [Table des périmètres du chômage](#table-des-p%C3%A9rim%C3%A8tres-du-ch%C3%B4mage)
      - [Table des recours antérieurs au chômage](#table-des-recours-ant%C3%A9rieurs-au-ch%C3%B4mage)
      - [Table des avis du CE](#table-des-avis-du-ce)
    - [Données de consommations d'activité partielle](#donn%C3%A9es-de-consommations-dactivit%C3%A9-partielle)
  - [Table de correspondance entre compte administratif URSSAF et siret](#table-de-correspondance-entre-compte-administratif-urssaf-et-siret)
  - [Données sur l'effectif](#donn%C3%A9es-sur-leffectif)
  - [Données sur les cotisations sociales et les débits (URSSAF caisse nationale)](#donn%C3%A9es-sur-les-cotisations-sociales-et-les-d%C3%A9bits-urssaf-caisse-nationale)
    - [Fichier de cotisations](#fichier-de-cotisations)
    - [Fichier de débits](#fichier-de-d%C3%A9bits)
      - [Codes état du compte](#codes-%C3%A9tat-du-compte)
      - [Codes procédure collective](#codes-proc%C3%A9dure-collective)
      - [Codes opération historique](#codes-op%C3%A9ration-historique)
      - [Code motif de l'écart négatif](#code-motif-de-l%C3%A9cart-n%C3%A9gatif)
    - [Fichier de délais](#fichier-de-d%C3%A9lais)
  - [Données sur le procédures collectives](#donn%C3%A9es-sur-le-proc%C3%A9dures-collectives)
  - [Données sur les CCSF](#donn%C3%A9es-sur-les-ccsf)
  - [Données fiscales](#donn%C3%A9es-fiscales)
  - [Ratios financiers](#ratios-financiers)

## Préambule

Ce dossier technique décrit les données utilisées dans le projet Signaux Faibles.

Les données utilisées proviennent des sources suivantes:

- **Données Sirene** Raison sociale, adresse, code APE, date de création etc.
- **Données DARES** Autorisations et consommations d'activité partielle, recours à l'intérim, déclaration des mouvements de main d’œuvre.
- **Données URSSAF** Données de défaillance, montant des cotisations, montant des dettes (part patronale, part ouvrière), demandes de délais de paiement, demandes préalables à l'embauche.
- **Données DGFiP** Liasses fiscales : bilans et comptes de résultat, ratios financiers.
- **Données INPI** Bilans publics, issus du RNCS.

Les sections ci-dessous détaillent la nature des données. Elles peuvent être à la fois exploitées dans l'interface web afin d'être présentées sous différents formats (tableau, graphes) aux utilisateurs, et à la fois en tant que sources pour le [modèle](algorithme-evaluation.md) prédictif. La distinction n'est pas systématiquement précisée, et les variables utilisées par ce dernier sont régulièrement amenées à évoluer et peuvent être consultées à travers les [fichiers de configuration](https://github.com/signaux-faibles/sf-datalake-mirror/tree/develop/src/sf_datalake/configuration) associés.

## Périmètre des données

### Entreprises :house:

L'entité de travail de base est l'unité légale (SIREN), que nous nommerons « entreprise » par abus de langage — même si ces deux notions sont distinctes, voir la [documentation](https://www.insee.fr/fr/metadonnees/definition/c1044) INSEE. 

## Données importées

### Base Sirene

Les fichiers de la base Sirene sont utilisés comme référence pour les établissements actifs. On extrait la raison sociale, le code NAF, l'adresse (y compris région et département qui permettent d'ouvrir les droits de consultation sur le terrain), et les dates d'activité.

Les fichiers utilisés sont les suivants :

- `StockEtablissement_utf8_geo.csv` (informations sur les établissements)
- `StockUniteLegale.csv` (informations sur les "entreprises")
- `StockEtablissementHistorique_utf8.csv`

La description détaillée des variables des fichiers de la base de données est disponible aux liens suivants :

- https://www.data.gouv.fr/en/datasets/base-sirene-des-entreprises-et-de-leurs-etablissements-siren-siret/
- https://github.com/cquest/geocodage-spd/blob/master/insee-sirene/README.md

|                          |                                             |
| ------------------------ | ------------------------------------------- |
| Source                   | Base Sirene et géo-sirene                   |
| Unité                    | siret                                       |
| Couverture siret         | Tout établissement actif                    |
| Fréquence de mise-à-jour | Source mise à jour et intégration mensuelle |

### Données d'activité partielle

Deux fichiers sont exploités : demandes et consommations d'activité partielle.

|                              |                                            |
| ---------------------------- | ------------------------------------------ |
| Source                       | DARES                                      |
| Unité                        | siret                                      |
| Couverture                   | Toutes les demandes et consommations       |
| Fréquence de mise-à-jour     | Mensuelle, export base complète            |
| Délai des données (conso)    | Mois précédent, à quelques exceptions près |
| Délai des données (demandes) | Temps réel                                 |

Données communes :

- **ID_DA** Numéro de la demande (11 caractères principalement des chiffres)
- **ETAB_SIRET** Numéro de SIRET de l'établissement (14 chiffres)
- **CODE_NAF2_2** Division sur 88 postes (2 premiers caractères de la variable CODE_NAF2, par exemple : `25`)
- **EFF_ENT** Effectif de l'entreprise (nombre entier)
- **EFF_ETAB** Effectif de l'établissement (nombre entier)
- **SOURCE** Source de la donnée (`SINAPSE` ou `EXTRANET`)
- **CODENAF2** et/ou **CODE_NAF2** Code NAF sur 5 positions (en majuscules, par exemple : `2573A`)
- **CODE_NAF_21** Section A...U sur 21 postes
- **CODE_NAF_TP** Pseudo indicatrice entreprises travaux publics vaut : `TP` si travaux publics, `AUTRES` sinon
- **EFF_ENT_TR_TDB** Effectif de l'entreprise par tranche (découpage Dares) : `1. Moins de 20 salariés`, `2. Entre 20 et 49 salariés`, `3. Entre 50 et 249 salariés`, `4. Entre 250 et +`

#### Données de demandes d'activité partielle

Données spécifiques aux demandes d'activité partielle :

- **EFF_ENT_TR** Effectif de l'entreprise par tranche (découpage DARES) : `1. Moins de 20 salariés`, `2. Entre 20 et 49 salariés`, `3. Entre 50 et 249 salariés`, `4. Entre 250 et 499 salariés`, `5. Entre 500 et 999 salariés`, `6. Plus de 1000 salariés`
- **ETAB_RSS** Dénomination (raison sociale) de l'établissement
- **ETAB_CODE_INSEE** Code INSEE de la commune de l'établissement
- **ETAB_VILLE** Ville de l'établissement
- **DATE_STATUT** Date du statut - création ou mise à jour de la demande - au format JJ/MM/AA
- **TX_PC** Taux horaire d'indemnisation (nombre avec 2 chiffres après la virgule)
- **TX_PC_ETAT_DARES** Taux de prise en charge par l'Etat (nombre avec 2 chiffres après la virgule)
- **TX_PC_UNEDIC_DARES** Taux de prise en charge pas l'Unédic (nombre avec 2 chiffres après la virgule)
- **DATE_DEB** Date de début de la période de chômage déterminée au format JJ/MM/AA
- **DATE_FIN** Date de fin de la période de chômage déterminée au format JJ/MM/AA
- **HTA** Nombre total d'heures autorisées (nombre décimal avec point)
- **MTA** Montant total autorisé (nombre décimal avec virgule)
- **EFF_AUTO** Effectifs autorisés (nombre entier)
- **PROD_HTA_EFF** ?
- **MOTIF_RECOURS_SE** Motif de recours à l'activité partielle (de `1` à `5`)
- **PERIMETRE_AP** Périmètre du chômage (de `1` à `4`)
- **RECOURS_ANTERIEUR** Activité partielle au cours des 36 derniers mois (3 ans) avant la demande (de `1` à `3`)
- **AVIS_CE** Avis du CE (de `0` à `3`)
- **S_HEURE_CONSOM_TOT** Nombre total d'heures consommées (nombre décimal avec point)
- **S_MONTANT_CONSOM_TOT** Montant total consommé (nombre décimal avec virgule)
- **S_EFF_CONSOM_TOT** (nombre entier)
- **CPT** ?
- **Date_Statut_Annee** Année du statut au format AAAA
- **Date_Statut_Mois** MOis du statut au format M
- **Date_Statut_Annee_Mois** Année et mois du statut au format AAAA_MM
- **Date_Statut_Annee_Trim** Année et trimestre du statut au format AAAA_T

##### Table des motifs de recours à l'activité partielle

| Code | Libellé                                                                             |
| ---- | ----------------------------------------------------------------------------------- |
| 1    | Conjoncture économique.                                                             |
| 2    | Difficultés d’approvisionnement en matières premières ou en énergie                 |
| 3    | Sinistre ou intempéries de caractère exceptionnel                                   |
| 4    | Transformation, restructuration ou modernisation des installations et des bâtiments |
| 5    | Autres circonstances exceptionnelles                                                |

##### Table des périmètres du chômage

| Code | Libellé                      |
| ---- | ---------------------------- |
| 1    | Réduction horaire tout Ets   |
| 2    | Réduction horaire partie Ets |
| 3    | Fermeture tempor. Tout Ets   |
| 4    | Fermeture tempor. Partie Ets |

##### Table des recours antérieurs au chômage

| Code | Libellé                                                      |
| ---- | ------------------------------------------------------------ |
| 1    | Aucun recours depuis 3 ans                                   |
| 2    | Recours au chômage partiel au cours des 3 années précédentes |
| 3    | Non renseigné                                                |

##### Table des avis du CE

| Code | Libellé       |
| ---- | ------------- |
| 0    | Non renseigné |
| 1    | Favorable     |
| 2    | Défavorable   |
| 3    | Sans objet    |

#### Données de consommations d'activité partielle

- **DEP** Code département sur 3 chiffres (par exemple : `073`)
- **REG** Code région sur 2 chiffres (par exemple : `23`)
- **MOIS** Mois au format MM/AAAA
- **HEURES** Heures consommées (chômées) dans le mois (nombre décimal avec point)
- **MONTANTS** Montants consommés dans le mois (nombre décimal avec virgule)
- **EFFECTIFS** Nombre de salariés en activité partielle dans le mois (nombre entier)
- **Date_Payement_Annee** Année de paiement au format AAAA
- **Date_Payement_Mois** Mois de paiement au format M
- **Date_Payement_Annee_Mois** Année et mois de paiement au format AAAA_MM
- **Date_Payement_Annee_Trim** Année et trimestre de paiement au format AAAA_T
- **Naf5_code** Code grand secteur
- **Naf5_libelle** Libellé du grand secteur
- **Naf38_code** Code NACE 38
- **Naf38_libelle** Libellé NACE 38

Les consommations affichées sur le frontend sont, depuis la mise en place de la nouvelle version de datapi, les consommations mensuelles issues des fichiers spécifiques livrés par la DARES.

Le valeur indiquée est le nombre d'équivalents temps plein correspondant à l'activité partielle (calculé à partir du volume d'heure mensuel **HEURES** sur la base d'une durée légale de 151,67 h).

Dans le cas où figurent différentes consommations mensuelles rattachées à différentes demandes concomitantes alors la valeur correspond à la somme des consommations.

### Table de correspondance entre compte administratif URSSAF et siret

|                          |                                           |
| ------------------------ | ----------------------------------------- |
| Source                   | URSSAF                                    |
| Couverture               | TODO                                      |
| Fréquence de mise-à-jour | Mensuellement                             |
| Délai des données        | En temps réel (parfois un mois de retard) |

- **Numéro de compte externe** Compte administratif URSSAF

- **Etat du compte** Compte ouvert (1) ou fermé (3) ?
- **Numéro siren** Numéro Siren de l'entreprise
- **Numéro d'établissement** Numéro Siret de l'établissement. Les numéros avec des lettres sont des sirets provisoires. Certains Siret sont des pseudo-siret de l'ACOSS, cf [détail](#pseudo-siret) ci-dessous.
- **Date de création de l'établissement** Date de création de l'établissement au format (A)AAMMJJ. (A)AA = AAAA - 1900.
- **Date de disparition de l'établissement** Date de disparition de l'établissement au format (A)AAMMJJ (cf date de création)

#### Pseudo-siret

Les pseudo-sirets sont des SIRET forgés par l'ACOSS pour le traitement spécifique de certaines activités avec travail temporaire, notamment : 
- Agences d'intérim,
- Agences de mannequinat,
- Agences de production de films.

Ces pseudo-sirets sont reconnaissables par le fait que le premier chiffre du NIC (10ème chiffre du SIRET) est 5,6,7 ou 9. 

### Données sur l'effectif

|                          |                               |
| ------------------------ | ----------------------------- |
| Source                   | URSSAF                        |
| Couverture               | ? (TODO)                      |
| Fréquence de mise-à-jour | Variable, tous les 3 à 6 mois |
| Délai des données        | 0 à 6 mois selon mise-à-jour  |

- **Siret** Siret de l'établissement. Certains Siret sont des pseudo-sirets, cf [détail](#pseudo-siret)
- **Compte** Compte administratif URSSAF
- **Rais_soc** Raison sociale
- **Ur_emet** Urssaf en charge de la gestion du compte
- **Dep** Département
- **effAAAAXY** effectif du mois AAAAXY. AAAA = année. X = trimestre. Y = N° du mois dans le trimestre (ex: 201631 vaut juillet 2016)

### Données sur les cotisations sociales et les débits (URSSAF caisse nationale)

Trois fichiers sont exploités : cotisations sociales, débits sur les cotisations sociales, et délais accordés sur les débits.

|                                 |                                         |
| ------------------------------- | --------------------------------------- |
| Source                          | URSSAF                                  |
| Couverture                      | ? (TODO)                                |
| Fréquence de mise-à-jour        | mensuelle (autour du 20 de chaque mois) |
| Délai des données (cotisations) | cotisations du mois précédent           |
| Délai des données (débits)      | débits sur les cot. du mois précédent   |

#### Fichier de cotisations

- **Compte** Compte administratif URSSAF
- **Periode_debit** Période en débit. _A ne pas prendre en compte_
- **ecn** Numéro écart négatif.
- **periode** Période. Format AAXY, cf ci-dessus l'effectif pour l'explication du format.
- **mer** Cotisation mise en recouvrement, en euros.
- **enc_direct** Cotisation encaissée directement, en euros.
- **cotis_due** Cotisation due, I euros. À utiliser pour calculer le montant moyen mensuel du : Somme cotisations dues / nb périodes

#### Fichier de débits

- **num_cpte** Compte administratif URSSAF
- **Siren** Siren de l'entreprise
- **Dt_immat** Date d'immatriculation du compte à l'Urssaf
- **Etat_cpte** Code état du compte. Cf la table ci-dessous.
- **Cd_pro_col** Code qui indique si le compte fait l'objet d'une procédure collective. Cf la table ci-dessous.
- **Periode** Période au format AAAAXY. Cf effectif pour l'explication.
- **Num_Ecn** L'écart négatif (ecn) correspond à une période en débit. Pour une même période, plusieurs débits peuvent être créés. On leur attribue un numéro d'ordre. Par exemple, 101, 201, 301 etc.; ou 101, 102, 201 etc. correspondent respectivement au 1er, 2ème et 3ème ecn de la période considérée.
- **Num_Hist_Ecn** Ordre des opérations pour un écart négatif donné.
- **Dt_trt_ecn** Date de comptabilisation de l'évènement (mise en recouvrement, paiement etc..). Format (A)AAMMJJ, ou (A)AA correspond à l'année à laquelle a été soustrait 1900. Exemple: 1160318 vaut 18 mars 2016, 990612 vaut 12 juin 1999.
- **Mt_PO** Montant des débits sur la part ouvrières **en centimes**. Sont exclues les pénalités et les majorations de retard.
- **Mt_PP** Montant des débits sur la part patronale **en centimes**. Sont exclues les pénalités et les majorations de retard.
- **Cd_op_ecn** Code opération historique de l'écart négatif. Cf table ci-dessous.
- **ecn** Code motif de l'écart négatif. Cf table ci-dessous

##### Codes état du compte

| Code | Description |
| ---- | ----------- |
| 1    | Actif       |
| 2    | Suspendu    |
| 3    | Radié       |

##### Codes procédure collective

| Code           | Description                             |
| -------------- | --------------------------------------- |
| 0, blanc, null | Pas de pro col                          |
| 1              | Pro col en cours                        |
| 2              | Pro col - plan de redressement en cours |
| 9              | Pro col sans dette à l'Urssaf           |

##### Codes opération historique

| Code | Description                               | Commentaire |
| ---- | ----------------------------------------- | ----------- |
| 1    | Mise en recouvrement                      |             |
| 2    | Paiement                                  |             |
| 3    | Admission en non valeur                   | Peu d'espoir de recouvrement |
| 4    | Remise de majoration de retard            |             |
| 5    | Abandon de solde debiteur                 | Montant trop faible pour justifier d'un recouvrement |
| 11   | Annulation de mise en recouvrement        |             |
| 12   | Annulation paiement                       |             |
| 13   | Annulation a-n-v                          |             |
| 14   | Annulation de remise de majoration retard |             |
| 15   | Annulation abandon solde debiteur         |             |


##### Code motif de l'écart négatif

| Code | Description                                                                                |
| ---- | ------------------------------------------------------------------------------------------ |
| 0    | Cde motif inconnu                                                                          |
| 1    | Retard dans le versement                                                                   |
| 2    | Absence ou insuffisance de versement                                                       |
| 3    | Taxation provisionelle. Déclarations non fournies                                          |
| 4    | Majorations de retard complémentaires Article R243-18 du code de la sécurité sociale       |
| 5    | Contrôle,chefs de redressement notifiés le JJ/MM/AA Article R243-59 de la Securité Sociale |
| 6    | Fourniture tardive des déclarations                                                        |
| 7    | Bases déclarées supérieures à Taxation provisionnelle                                      |
| 8    | Retard dans le versement et fourniture tardive des déclarations                            |
| 9    | Absence ou insuffisance de versement et fourniture tardive des déclarations                |
| 10   | Rappel sur contrôle et fourniture tardive des déclarations                                 |
| 11   | Régularisation d'une taxation provisionnelle                                               |
| 12   | Régularisation annuelle                                                                    |
| 13   | Rejet du titre de paiement par la banque .                                                 |
| 14   | Modification d'affectation d'un crédit                                                     |
| 15   | Annulation d'un crédit                                                                     |
| 16   | Régularisation suite à modification du Taux Accident du Travail                            |
| 17   | Régularisation suite à assujettissement au transport (origine débit sur PJ=4)              |
| 18   | Majorations pour non respect de paiement par moyen dématérialisé Article L243-14           |
| 19   | Rapprochement TR/BRC sous réserve de vérification ultérieure                               |
| 20   | Cotisations complémentaires suite modification des revenus déclarés                        |
| 21   | Cotisations complémentaires suite à non fourniture du contrat d'exonération                |
| 22   | Contrôle. Chefs de redressement notifiés le JJ/MM/AA. Article L324.9 du code du travail    |
| 23   | Cotisations complémentaires suite conditions d'exonération non remplies                    |
| 24   | Absence de versement                                                                       |
| 25   | Insuffisance de versement                                                                  |
| 26   | Absence de versement et fourniture tardive des déclarations                                |
| 27   | Insuffisance de versement et fourniture tardive des déclarations                           |

#### Fichier de délais

|                          |                         |
| ------------------------ | ----------------------- |
| Source                   | URSSAF                  |
| Couverture               | Tous les délais         |
| Fréquence de mise-à-jour | Mensuellement           |
| Délai des données        | Créations en temps réel |

- **Numero de compte externe** Compte administratif URSSAF
- **Numéro de structure** Le numéro de structure est l'identifiant d'un dossier contentieux
- **Date de création** Date de création du délai. Format aaaa-mm-jj
- **Date d'échéance** Date d'échéance du délai. Format aaaa-mm-jj
- **Durée délai** Durée du délai en jours.
- **Dénomination premiére ligne** Raison sociale de l'établissement.
- **Indic 6M** Délai inférieur ou supérieur à 6 mois? Modalités INF et SUP.
- **année ( Date de création )** Année de création du délai.
- **Montant global de l'échéancier** Montant global de l'échéancier, en euros.
- **Numéro de structure** Champs en double, cf plus haut.
- **Code externe du stade**

### Données sur le procédures collectives

|                          |                                   |
| ------------------------ | --------------------------------- |
| Source                   | URSSAF                            |
| Couverture               | Toutes les procédures collectives |
| Fréquence de mise-à-jour | Mensuellement                     |
| Délai des données        | Créations en temps réel ?         |

- **Siret** Siret de l'établissement
- **Siren** Siren de l'entreprise
- **Dt_effet** Date effet de la procédure collective au format JJMMMAAAA, par exemple 24FEB2014
- **Cat_v2** ? TODO
- **Lib_actx_stdx** Champs double: Nature procédure + évènement.

### Données sur les CCSF

Les délais gérés par les Commissions des Chefs des Services Financiers (CCSF) sont des actions contentieuses dont l'objectif est de négocier un plan de paiement de dettes avec des entreprises débitrices. Il y a une commission par département et y siègent des représentants des Urssaf, du Fisc, des Douanes, de Pole emploi et de la CCMSA.

C'est donc le même type de dispositif que les délais classiques URSSAF appelés aussi sursis à poursuite. Même si les CCSF traitent tout dossier, ce sont en général les entreprises d'assez grande taille avec des enjeux financiers supérieurs à la moyenne des dossiers de recouvrement qui sont concernées.

- **Compte** Compte administratif URSSAF
- **Date de traitement** Date de début de la procédure CCSF au format (A)AAMMJJ, ou (A)AA = AAAA -1900
- **Code externe du stade** Stade de la demande de délai: `DEBUT`, `REFUS`, `APPROB`, `FIN` ou `REJ`
- **Code externe de l'action** Code externe de l'action: `CCSF` (valeur systématique)

### Données fiscales

|                          |                           |
| ------------------------ | ------------------------- |
| Source                   | DGFiP                     |
| Couverture               | Ensemble du périmètre     |
| Fréquence de mise-à-jour | Trimestre                 |
| Délai des données        | Dépend du type de données |

Un nombre important de liasses fiscales sont mises à disposition par la DGFiP et il serait trop lourd de détailler l'ensemble des sources et schémas ici. Voici une liste non-exhaustive des données fournies, les références des formulaires permettant de retrouver les différentes informations renseignées par les entreprises et fournies à l'administration fiscale.

| Formulaire CERFA      | Contenu                                                                                 |
| --------------------- | --------------------------------------------------------------------------------------- |
| 2031-SD               | Résultats BIC - Impôt sur le revenu                                                     |
| 2033-A-SD             | RSI - Bilan simplifié                                                                   |
| 2033-B-SD             | RSI - Compte de résultat simplifié de l'exercice                                        |
| 2033-C-SD             | RSI - Immobilisations, amortissements, plus-values, moins-values                        |
| 2033-D-SD             | RSI - Relevé des provisions, amortissements dérogatoires, déficits                      |
| 2033-E-SD             | RSI - Détermination des effectifs et de la valeur ajoutée                               |
| 2033-F-SD             | RSI - Composition du capital social                                                     |
| 2033-G-SD             | RSI - Filiales et participations                                                        |
| 2050-SD               | RN - Bilan - Actif                                                                      |
| 2051-SD               | RN - Bilan - Passif                                                                     |
| 2052-SD et 2053-SD    | RN - Compte de résultat de l'exercice                                                   |
| 2054-SD               | RN - Immobilisations                                                                    |
| 2054bis-SD            | RN - Ecarts de réévaluation sur                                                         |
| 2055-SD               | RN - Amortissements                                                                     |
| 2056-SD               | RN - Provisions inscrites au bilan                                                      |
| 2057-SD               | RN - Etat des échéances des créances et des dettes à la clôture de l'exercice           |
| 2058-A-SD             | RN - Détermination du résultat fiscal                                                   |
| 2058-B-SD             | RN - Déficits, indemnités pour congés à payer et provisions non déductibles             |
| 2058-C-SD             | RN - Tableau d'affectation des résultats et renseignements divers                       |
| 2059-A-SD             | RN - Détermination des plus et moins values                                             |
| 2059-B-SD             | RN - Affectation des plus-values à court terme et des plus-values de fusion ou d'apport |
| 2059-C-SD             | RN - Suivi des moins-values à long terme                                                |
| 2059-D-SD             | RN - Réserve spéciale des plus-values à long terme                                      |
| 2059-E-SD             | RN - Détermination des effectifs et de la valeur ajoutée                                |
| 2059-F-SD             | RN - Composition du capital social                                                      |
| 2059-G-SD             | RN - Filiales et participations                                                         |
| 2065-SD et 2065bis-SD | Impôt sur les sociétés (IS)                                                             |
| 2066-SD               | Impôt sur les sociétés (IS) - Déclaration complémentaire                                |
| 3310-CA3-SD           | Taxe sur la valeur ajoutée (TVA) et taxes assimilées - Régime réel normal (RN)          |
| 3517-S-SD             | Taxe sur la valeur ajoutée (TVA) et taxes assimilées - Régime réel simplifié (RSI)      |

Ces données sont hébergées sur le lac de données de la DTNUM et ne sont pas [importées](procedure-import-donnees.md) en vue d'une mise en valeur pour les utilisateurs de l'application, mais sont exclusivement exploitées pour la [prédiction](algorithme-evaluation.md) de la défaillance.

### Ratios financiers

La consultation des données confidentielles ou semi-confidentielles de bilan n'a pas pu être étendue à l'ensemble des utilisateurs de l'application sur décision de la DGFiP et est donc réservée aux agents de cette administration. Afin de partiellement palier cette limitation, des données ouvertes — correspondant aux bilans publics — issues du registre national du commerce et des sociétés et publiées par l'INPI ont été traitées en collaboration avec l'équipe « Base commune entreprise » du ministère du travail, et sont désormais régulièrement diffusées de manière ouverte sur la plateforme open data du ministère de l'économie : https://data.economie.gouv.fr/explore/dataset/ratios_inpi_bce/information/

Ces ratios sont présentés au sein de l'application dans les [fiches entreprises](https://signaux-faibles.gitbook.io/guide-dutilisation-et-f.a.q.-de-signaux-faibles/consulter-des-informations-dentreprises/fiche-entreprise/les-informations-financieres), mais ont également été réutilisés immédiatement par d'autres projets publics tels que [fiche commune entreprise](https://fce.fabrique.social.gouv.fr/home), et [l'annuaire des entreprises](https://annuaire-entreprises.data.gouv.fr/).

Les ratios financiers présentés dans l'application sont détaillés [ici](https://signaux-faibles.gitbook.io/guide-dutilisation-et-f.a.q.-de-signaux-faibles/consulter-des-informations-dentreprises/analyse-financiere).
