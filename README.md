![CI](https://github.com/signaux-faibles/opensignauxfaibles/workflows/CI/badge.svg)

# Open Signaux Faibles

Le projet [Signaux Faibles](https://beta.gouv.fr/startups/signaux-faibles.html) fournit une plateforme technique de détection anticipée d'entreprises en difficulté, en s'appuyant sur l'exploitation des signaux faibles.

La commande `sfdata` fournie dans ce dépôt centralise toutes les fonctionnalités d'import et de préparation des données.

Les données préparées s'adressent à deux consommateurs aval :

- La data science (qui a besoin d'une grande profondeur historique)
- Le front-end de l'application Signaux Faibles

Contact: [contact@signaux-faibles.beta.gouv.fr](mailto:contact@signaux-faibles.beta.gouv.fr)

## Architecture

- Golang pour l'import des données
- PostgreSQL 17 pour le stockage et la préparation des données

## Build et tests

```bash
# Cloner le code en local
git clone https://github.com/signaux-faibles/opensignauxfaibles.git
cd opensignauxfaibles

# Compiler le binaire
make build

# Compiler pour la production (Linux AMD64)
make build-prod

# Exécuter tous les tests (unit + e2e)
./test-all.sh

# Exécuter tous les tests et mettre à jour les snapshots/golden files
./test-all.sh --update

# Exécuter uniquement les tests unitaires
go test ./...

# Exécuter uniquement les tests e2e
go test -tags=e2e ./...

# Exécuter un seul test
go test ./lib/engine -run TestSpecificTest

# Démarrer le conteneur PostgreSQL pour les tests manuels locaux
make start-postgres

# Arrêter le conteneur PostgreSQL
make stop-postgres
```

## Usage 

```bash
# Afficher l'aide
./sfdata --help

# Importer des données depuis un répertoire de batch
./sfdata import --batch 1802 --path /path/to/data

# Importer avec une configuration de batch explicite
./sfdata import --batch 1802 --batch-config ./batch.json

# Importer sans filtrage
./sfdata import --batch 1802 --no-filter

# Importer uniquement des parsers spécifiques (plusieurs possibles)
./sfdata import --batch 1802 --parsers apconso --parsers cotisation

# Dry run (parser sans écrire en DB/CSV)
./sfdata import --batch 1802 --dry-run

# Parser un seul fichier vers stdout
./sfdata parseFile --parsers cotisation --file /path/to/file.csv

# Exporter les données des vues postgres pour la data science
./sfdata export --path /path/to/dir/
```

## Configuration

Voir `config-sample.toml` dans les sources.

Par ordre de priorité, la configuration peut être définie:

- via des variables d'environnement (majuscules, points remplacés par des `_`, 
  ex. `LOG_LEVEL`)
- dans "/etc/opensignauxfaibles/config.toml"
- dans "~/.opensignauxfaibles/config.toml"
- dans "./config.toml"

Pour les variables d'environnement, les variables sont en majuscules et les 
points `.` des options imbriquées sont remplacées par des `_` (par exemple, 
"log.level" est défini via "LOG_LEVEL").

### Configuration de batch spécifique

Voici un modèle de configuration de batch :

```json
 {
    "key": "1802",
    "files": {
      "apconso": [
        "1802/apconso_20180201.csv",
        "1802/apconso_20180215.csv.gz"
      ],
      "effectif": [
        "1802/effectif_202402.csv"
      ],
      "effectif_ent": [
        "1802/effectif_entreprise_202402.csv"
      ],
      "filter": [
        "1802/filter_custom_siren.csv"
      ]
  }
}
```

Les types de parser disponibles peuvent être consultés dans 
"lib/engine/parser_types.go".

Le type "filter" permet de donner une liste de siren (csv avec une colonne 
"siren") pour filtrer l'import.

## Ajouter un nouveau parser

Lors de l'ajout d'un nouveau parser de sources de données :

1. Définir la constante de type de parser dans `lib/engine/parser_types.go`
2. Créer le répertoire du parser : `lib/parsing/[parser_name]/`
3. Implémenter les interfaces `Parser` et `ParserInst` (peut s'appuyer sur des 
   implémentations existantes, par ex. CsvParserInst)
4. Définir la struct tuple implémentant l'interface `Tuple` (Key(), Scope(), 
   Type()). Les tags `csv` et `sql` définissent le nom des colonnes dans les 
   sorties .csv et postgreSQL respectivement. 
5. Enregistrer le parser dans `lib/registry/main.go`
6. Créer la table `stg_[parser_name]` via une migration 
   (`lib/db/migrations/`), puis l'ajouter dans `lib/db/tables.go`
7. Ajouter une vue ou une vue matérialisée `clean_[parser_name]` dans une 
   migration puis pour les vues matérialisées, l'ajouter dans 
   `lib/db/tables.go`.
8. Ajouter le support du sink dans `lib/sinks/postgresSink.go`. Si la vue est 
   matérialisée, sélectionner les conditions de mise à jour de la vue dans la 
   fonction `CreateSink`. 
9. Ajouter la reconnaissance du pattern de nom de fichier dans 
   `lib/prepare-import/parsertypes.go` pour que les fichiers types pour ce 
   parser soient reconnus automatiquement.

# Notes d'architecture

```
Fichiers de données → Parser → Filtre → Sink (PostgreSQL + CSV)
```

Le pipeline d'importation se compose de :

1. **Préparation du batch** (`lib/prepare-import/`) : Découvre les fichiers de 
   données et infère leurs types de parser depuis les noms de fichiers, ou 
   charge une configuration de batch explicite. Les heuristiques peuvent être 
   consultées dans `lib/prepare-import/parsertypes.go`
2. **Parsing** (`lib/parsing/` pour l'implémentation de chaque parser, 
   `lib/engine/` pour la mécanique générale) : Lit les fichiers de données 
   brutes et extrait des tuples structurés (parallélisé)
3. **Filtrage** (`lib/filter/`) : Applique un filtrage basé sur le SIREN pour limiter le volume
4. **Sinks** (`lib/sinks/`) : destination des données. Écrit les données 
   nettoyées dans les tables PostgreSQL et les fichiers CSV (les données des  
   fichiers CSV sont avant préparation, car la préparation des données se fait 
   via des vues Postgres)

## Base de données

- Architecture à deux couches :
  - tables `stg_*` : Données brutes/staging importées 
  - tables/vues `clean_*` : Données enrichies et nettoyées. Ce sont ces tables qui doivent être utilisées par les consommateurs des données downstream.

Voir [la documentation des tables](./docs/documentation_tables.md) pour plus 
d'informations.

Les vues matérialisées des données enrichies sont mis à jour automatiquement 
lorsque les données dont elles dépendent sont importées, via une dépendance 
exlpicite dans le code (cf `viewsToRefresh` dans 
[./lib/sinks/postgresSink.go](./lib/sinks/postgresSink.go))

## Migrations de base de données

Les migrations sont définies dans `lib/db/migrations.go`.
Elles sont automatiquement effectuées au début de l'import. 

Le fonctionnement est simple : les migrations sont numérotées dans l'ordre, la 
table `migrations` stocke la dernière migration appliquée, et golang applique 
les migrations suivantes au besoin (via l'utilitaire 
[`tern`](https://github.com/JackC/tern))

Il n'y a pour l'instant pas de commande pour n'effectuer que les migrations, 
cependant, l'import d'un petit fichier (e.g. `--parsers delai`) est un 
contournement pour migrer la base de données sans import complet.

Le test de bout en bout exécute toutes les migrations et vérifie donc qu'elles 
sont valides.

## Parsers

- Chaque source de données a un parser dédié (ex. `apconso`, `urssaf`, 
  `effectif`)
- Les types de parser sont définis dans `lib/engine/parser_types.go`

## Filtrage

Les données sont filtrées à l'import via leur SIREN, afin de n'importer que 
des données d'intérêt et restreindre le volume de données stockées. Le 
filtrage s'applique sur tous les fichiers sauf les fichiers `sirene` et 
`sirene_ul`, car le front-end a besoin de l'intégralité des données sur ces 
bases.

Par sécurité, l'absence de filtre est par défaut une erreur. Pour importer 
l'intégralité des données sans filtrage, il est nécessaire d'utiliser le flag 
`--no-filter`.

Si aucun filtre explicite n'est fourni (fichier commençant par "filter", ou 
explicitement défini dans un batch JSON), le filtre va être lu de la base de 
données. Si aucun filtre n'est stocké en base, il faut que le batch importé 
possède un fichier de type "effectif_ent" (effectif de l'entreprise), afin de 
générer le filtre.

Le système de filtrage se fait en trois étapes :

- Si pas de filtre explicite fourni, et en présence d'un fichier 
  "effectif_ent", on met à jour le filtre (entreprises qui ont atteint ou 
  dépassé 10 salariés dans les 120 derniers mois)
- Le filtre explicite ou la table `stg_filter_import` est utilisé pour le 
  filtrage à l'import (à l'exception des données "sirene" et "sirene_ul" qui 
  sont importées intégralement)
- Un filtrage supplémentaire pour arriver au périmètre définitif (table 
  `clean_filter`, cf. [documentation 
  spécifique](./docs/documentation_clean_filter.md)). Ce périmètre définitif 
  est utilisé pour la construction des vues `clean_[parser_name]`.

Ces trois étapes permettent d'écarter le plus gros volume de données qui ne 
nous intéressent pas à l'import, en laissant la possibilité d'affiner le 
filtrage dans un second temps.

### Faire évoluer la logique de filtrage

Pour faire évoluer la logique de filtrage, utiliser `CREATE OR REPLACE VIEW` 
sur la vue `siren_blacklist_logic`, puis `REFRESH MATERIALIZED VIEW 
siren_blacklist;` qui stocke une copie matérialisée pour des raisons de 
performance (cette construction en deux étape vient du fait que les vues 
matérialisées ne permettent pas de mise-à-jour sur place et nécessitent un 
`DROP ... CASCADE` qu'on souhaite éviter).

# Performance 

Pour des questions de performances, les indexes sont supprimés au début de 
l'import des données, et reconstruits a posteriori. 

Pour éviter qu'ils ne soient définitivement perdus si une anomalie survient 
pendant l'import (perte de connexion à la base de données par exemple), ils 
sont stockés dans une table `tmp_saved_indexes`. Ainsi, si la reconstruction 
échoue, il suffit de relancer la pipeline pour que les indexes stockés en base 
soient reconstruits. 

Par ailleurs opensignauxfaibles augmente, le temps de l'import, les paramètres 
`work_mem` (pour le rafraîchissement des vues matérialisées) et 
`maintenance_work_mem` (pour la reconstruction des indexes), pour une 
meilleure utilisation de la mémoire vive.

# Consommation des données 

Les consommateur de données doivent exclusivement consommer les données 
préparées des vues préfixées par `clean_xxx`. 

**Note concernant l'utilisation des indexes** : pour rechercher par SIREN dans 
une table indexée par SIRET, utiliser une requête du type `siret LIKE 
'123456789%'` pour profiter de l'index. Le fait de ne pas cumuler les indexes 
réduit sensiblement le temps de reconstruction des indexes à chaque import. 
