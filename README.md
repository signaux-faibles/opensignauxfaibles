![CI](https://github.com/signaux-faibles/opensignauxfaibles/workflows/CI/badge.svg) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/47a9094cf7bd4f7387a10151b90ed609)](https://www.codacy.com/gh/signaux-faibles/opensignauxfaibles/dashboard?utm_source=github.com&utm_medium=referral&utm_content=signaux-faibles/opensignauxfaibles&utm_campaign=Badge_Grade)

# Open Signaux Faibles

Le projet [Signaux Faibles](https://beta.gouv.fr/startups/signaux-faibles.html) fournit une plateforme technique de détection anticipée d'entreprises en difficulté, en s'appuyant sur l'exploitation des signaux faibles.

En intéragissant avec une base de données MongoDB, la commande `sfdata` fournie dans ce dépôt centralise toutes les fonctionnalités du module de traitement de données:

- Gestion des batches d'intégration
- Exécution des traitements
- Export des données

Note: Précedemment, ces fonctionnalités étaient mises à disposition via un serveur HTTP nommé `dbmongo`.

Contact: [contact@signaux-faibles.beta.gouv.fr](mailto:contact@signaux-faibles.beta.gouv.fr)

## Architecture

- Golang
- MongoDB 4.2
- Fonctions map-reduce: TypeScript (TS) et JavaScript (JS)

## Dépendances / pré-requis

- [Node.js](https://nodejs.org/) (voir version spécifiée dans `js/.nvmrc`), pour transpiler les fonctions map-reduce de TypeScript vers JavaScript et exécuter les tests automatisés, après toute modification de ces fonctions.
- [git secret](https://git-secret.io/), pour (dé)chiffrer les fichiers de données privées utilisés dans certains tests automatisés.

## Installation

```bash
$ git clone https://github.com/signaux-faibles/opensignauxfaibles.git
$ cd opensignauxfaibles
$ go generate ./...
$ go build -o sfdata
$ go test ./...
```

Dans l'arbre de sources de l'installation go, vous trouverez tous les répertoires nécessaires à l'exécution.

## Configuration

Voir `config-sample.toml` dans les sources.

Par ordre de priorité, le fichier de configuration peut se trouver dans:

- /etc/opensignauxfaibles/config.toml
- ~/.opensignauxfaibles/config.toml
- ./config.toml

## Usage

La commande `sfdata` s'inscrit dans un workflow d'intégration de données. Pour plus d'informations, consulter [signaux-faibles/documentation](https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#workflow-classique).

## Développement et tests automatisés

Afin de prévenir les régressions, plusieurs types de tests automatisés sont inclus dans le dépôt:

- tests unitaires de `sfdata`: `$ go test ./...`
- tests unitaires et d'intégration des fonctions map-reduce: `$ cd js && npm test`
- tests de bout en bout: `$ tests/test-*.sh`

Tous ces tests sont exécutés en environnement d'Intégration Continue (CI) après chaque commit poussé sur GitHub, grâce à GitHub actions, tel que défini dans les fichiers `yaml` du répertoire `.github/workflows`.

Il est possible de tous les exécuter en local: `$ ./test-all.sh`.

### Prérequis: accès aux données chiffrées

#### Introduction

Certains tests (ex: `test-reduce-2.sh` et `algo2_golden_tests.ts`) s'appuient sur des fichiers de données sensibles qui doivent donc être déchiffrés avant l'exécution de ces tests.

Reconnaissables à leur extension `.secret`, ces fichiers sont stockés dans le répertoire `/tests` du dépôt. Ils sont chiffrés et déchiffrés à l'aide de [git secret](https://git-secret.io/), commande basée sur l'outil de chiffrage GPG. Chaque développeur désireux d'accéder et/ou de modifier ces fichiers doit y être autorisé par un autre développeur y ayant déjà acccès, en intégrant sa clé publique GPG au porte-clés du projet (incarné par le fichier `.gitsecret/keys/pubring.kbx`).

À noter qu'une clé GPG a été générée et intégrée afin de permettre le déchiffrage de ces données lors de l'exécution des tests automatisés en environnement d'Intégration Continue (CI).

#### Commandes usuelles

```sh
$ git secret reveal # pour déchiffrer les fichiers *.secret
$ git secret changes # pour afficher les modifications apportées aux fichiers privées *en clair*
$ git secret hide -d # pour chiffrer puis effacer les fichiers privés *en clair* qui ont été modifiés
```

#### Procédure d'ajout d'une clé GPG

Instructions à suivre par le développeur demandant l'accès aux données privées:

1. installer `git secret` (cf [instructions](https://git-secret.io/installation), ex: `sudo apt-get install git-secret`) et `gpg` (si ce n'est pas encore le cas) sur votre machine
2. créer une clé GPG avec `$ gpg --gen-key` (cf [using GPG](https://git-secret.io/#using-gpg))
3. exporter la clé publique (`$ gpg --export --armor your.email@address.com > my-public-key.gpg`) puis l'envoyer à une personne ayant déjà accès aux fichiers chiffrés, pour qu'elle puisse vous y donner droit également (cf [adding someone to a repository](https://git-secret.io/#usage-adding-someone-to-a-repository-using-git-secret), instructions à suivre par un des développeurs ayant déjà accès)
4. une fois que le fichier `.gitsecret/keys/pubring.kbx` a bien été mis à jour, récupérer la dernière version de cette branche (`$ git pull`)
5. exécuter `$ git secret reveal` => les fichiers listés dans `.gitsecret/paths/mapping.cfg` seront déchiffrés
6. pour vérifier que le chiffrage fonctionne également: `$ git secret hide` va modifier les fichiers avec l'extension `.secret`. (vous n'avez pas besoin de créer un commit si vous n'avez pas modifié les fichiers de données après l'étape précédente)

#### Notes pour les développeurs autorisés à accéder aux données chiffrées

- Penser à rechiffrer les fichiers après l'intégration de toute nouvelle clé GPG, afin de les rendre accessible au propriétaire de cette clé.
- Sachant que la version déchiffrée des fichiers privées est ignorée par `git`, il ne faut pas oublier de les chiffrer à nouveau après chaque modification, puis de créer un commit incluant les modifications des fichiers `*.secret`. En réponse à cela, la documentation de `git secret` suggère la mise en place d'un "pre-commit hook" dans `git`.

### Développement des fonctions map-reduce (TS/JS)

```sh
$ cd js
$ nvm use # pour utiliser la version de Node.js spécifiée dans .nvmrc
$ npm install # pour installer les dépendances
$ npm test # pour exécuter les tests unitaires et d'intégration, tel que décrit dans package.json
```

### Intégration et test de modifications des fonctions map-reduce (TS/JS)

```sh
$ ./test-all.sh # va regénérer jsFunctions.go, recompiler sfdata et exécuter tous les tests
```
