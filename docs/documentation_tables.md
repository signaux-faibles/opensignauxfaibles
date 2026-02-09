# Documentation des tables dans PostgreSQL

- Architecture à deux couches :
  - tables `stg_*` : Données brutes/staging importées 
  - tables/vues `clean_*` : Données préparées et nettoyées. Ce sont ces tables qui doivent être utilisées par les consommateurs des données downstream.

## Données nettoyées

Ces vues (parfois matéralisées, parfois non selon leur complexité) sont les données filtrées et préparées pour consommation avale.

|           name            |       type        |                                             description                                              |
|---------------------------|-------------------|------------------------------------------------------------------------------------------------------|
| [clean_filter](./documentation_clean_filter.md)              | view              | Périmètre des données enrichies = stg_filter_import - siren_blacklist                                |
| [clean_ap](./documentation_clean_ap.md)                  | materialized view | Données enrichie et agrégée d'activité partielle                                                     |
| clean_sirene              | view              | Données enrichies sur les établissements (non filtrées sur le périmètre SF)                          |
| clean_sirene_ul           | view              | Données enrichies sur les entreprises (non filtrées sur le périmètre SF)                             |
| clean_sirene_histo        | view              | Données historiques enrichies sur les établissements (on ne conserve que les évènements qui impliquent un changement d'état administratif) |
| clean_cotisation          | view              | Données enrichies sur les cotisations                                                                |
| [clean_debit](./documentation_clean_debit.md)               | materialized view | Données enrichies sur les débits                                                                     |
| [clean_procol](./documentation_clean_procol.md)              | view              | Données enrichies de procédures collectives                                                          |
| clean_delai               | view              | Données enrichies de délais de paiement des cotisations sociales                                     |
| clean_effectif            | view              | Données enrichies des effectifs d'établissements                                                     |
| clean_effectif_ent        | view              | Données enrichies des effectifs d'entreprises                                                        |


## Données techniques

Ces données permettent le bon fonctionnement de la pipeline : données de migration, de filtrage,
logs, labels pour préparér les données (les labels sont injectés directement via les migrations).

Le périmètre (liste de SIREN) de l'import pour les consommateurs aval est défini dans `clean_filter`.

### Périmètre et filtrage

|           name            |       type        |                                             description                                              |
|---------------------------|-------------------|------------------------------------------------------------------------------------------------------|
| stg_filter_import         | table             | Périmètre d'import des données brutes, filtré sur l'effectif uniquement                               |
| siren_blacklist           | materialized view | Siren à exclure du périmètre d'import (à privilégier sur clean_filter lorsque la performance compte) |
| clean_filter              | view              | Périmètre des données préparées = stg_filter_import - siren_blacklist                                |

### Migrations 

|           name            |       type        |                                             description                                              |
|---------------------------|-------------------|------------------------------------------------------------------------------------------------------|
| migrations                | table             | Dernière migration de base de donnée appliquée                                                       |

### Labels

|           name            |       type        |                                             description                                              |
|---------------------------|-------------------|------------------------------------------------------------------------------------------------------|
| labels_motif_recours      | table             | Libellés pour les recours à l'activité partielle                                                     |
| categories_juridiques     | table             | Libellés pour les catégories juridiques                                                              |
| naf_codes                 | table             | Libellés pour la nomenclature d'activité                                                             |

### Logs

|           name            |       type        |                                             description                                              |
|---------------------------|-------------------|------------------------------------------------------------------------------------------------------|
| import_logs               | table             | Logs des données importées via OpenSignauxFaibles, préparés après chaque import                      |

### Données intermédiaires

Pour simplifier ou améliorer la performance de la préparation de certaines requêtes, des vues intermédiaires sont crées.

|           name            |       type        |                                             description                                              |
|---------------------------|-------------------|------------------------------------------------------------------------------------------------------|
| stg_apconso_by_period     | view              | Données intermédiaires                                                                               |
| stg_apdemande_by_period   | materialized view | Données intermédiaires                                                                               |
| stg_tmp_debits_simplified | materialized view | Données intermédiaires                                                                               |

## Données brutes

Ce sont les données importées, après une première phase de nettoyage 
(renommage, filtrage, standardisation des données). 

|           name            |       type        |                                             description                                              |
|---------------------------|-------------------|------------------------------------------------------------------------------------------------------|
| stg_apdemande             | table             | Données brutes d'autorisation d'activité partielle                                                   |
| stg_apconso               | table             | Données brutes de consommation d'activité partielle                                                  |
| stg_sirene                | table             | Données brutes sur les établissements (non filtrées sur le périmètre SF)                             |
| stg_sirene_ul             | table             | Données brutes sur les entreprises (non filtrées sur le périmètre SF)                                |
| stg_sirene_histo          | table             | Données historiques brutes sur les établissements                                                    |
| stg_cotisation            | table             | Données brutes sur les cotisation                                                                    |
| stg_debit                 | table             | Données brutes sur les débits                                                                        |
| stg_procol                | table             | Données brutes de procédures collectives                                                             |
| stg_delai                 | table             | Données brutes de délais de paiement des cotisations sociales                                        |
| stg_effectif              | table             | Données brutes des effectifs d'établissements                                                        |
| stg_effectif_ent          | table             | Données brutes des effectifs d'entreprises                                                           |
