![CI](https://github.com/signaux-faibles/opensignauxfaibles/workflows/CI/badge.svg) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/47a9094cf7bd4f7387a10151b90ed609)](https://www.codacy.com/gh/signaux-faibles/opensignauxfaibles/dashboard?utm_source=github.com&utm_medium=referral&utm_content=signaux-faibles/opensignauxfaibles&utm_campaign=Badge_Grade)

# Open Signaux Faibles

Le projet [Signaux Faibles](https://beta.gouv.fr/startups/signaux-faibles.html) fournit une plateforme technique de détection anticipée d'entreprises en difficulté, en s'appuyant sur l'exploitation des signaux faibles.

La commande `sfdata` fournie dans ce dépôt centralise toutes les 
fonctionnalités d'import des données.

Contact: [contact@signaux-faibles.beta.gouv.fr](mailto:contact@signaux-faibles.beta.gouv.fr)

## Architecture

- Golang
- MongoDB 4.2

## Installation

```bash
$ git clone https://github.com/signaux-faibles/opensignauxfaibles.git
$ cd opensignauxfaibles
$ make build
$ go test ./...
```

Dans l'arbre de sources de l'installation go, vous trouverez tous les répertoires nécessaires à l'exécution.

Avant de démarrer les test il est nécessaire de lancer le démon Docker.

## Configuration

Voir `config-sample.toml` dans les sources.

Par ordre de priorité, la configuration peut être définie:

- via des variables d'environnement
- dans /etc/opensignauxfaibles/config.toml
- dans ~/.opensignauxfaibles/config.toml
- dans ./config.toml

Pour les variables d'environnement, les variables sont en majuscules et les 
points `.` des options imbriquées sont remplacées par des `_` (par exemple, 
"log.level" est défini via "LOG_LEVEL").

## Usage

La commande `sfdata` s'inscrit dans un workflow d'intégration de données. Pour plus d'informations, consulter [signaux-faibles/documentation](https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#workflow-classique).

## Développement et tests automatisés

Afin de prévenir les régressions, plusieurs types de tests automatisés sont inclus dans le dépôt:

- tests unitaires de `sfdata`: `$ go test ./...`
- tests de bout en bout: `$ go test -tags=e2e ./...`

Tous ces tests sont exécutés en environnement d'Intégration Continue (CI) après chaque commit poussé sur GitHub, grâce à GitHub actions, tel que défini dans les fichiers `yaml` du répertoire `.github/workflows`.

Il est possible de tous les exécuter en local: `$ ./test-all.sh`.
