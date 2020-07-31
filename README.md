# Open Signaux Faibles

Solution logicielle pour la détection anticipée d'entreprises en difficulté

## Architecture

- Back-end: golang
- Front-end: vuetify
- MongoDB 4.2
- Fonctions map-reduce: TypeScript (TS) et JavaScript (JS)

## Dépendances / pré-requis

- [Node.js](https://nodejs.org/) (voir version spécifiée dans `dbmongo/js/.nvmrc`), pour transpiler les fonctions map-reduce de TypeScript vers JavaScript et exécuter les tests automatisés, après toute modification de ces fonctions.
- [git-secret](https://git-secret.io/), pour (dé)chiffrer les fichiers de données privées utilisés dans certains tests automatisés. (voir description de https://github.com/signaux-faibles/opensignauxfaibles/pull/113)

## Installation

```bash
$ git clone https://github.com/signaux-faibles/opensignauxfaibles.git
$ cd opensignauxfaibles
$ cd dbmongo
$ go generate ./...
$ go build
$ go test ./...
```

Dans l'arbre de sources de l'installation go, vous trouverez tous les répertoires nécessaires à l'exécution.

TODO:

- linker correctement les procédures R avec le core golang
- provoquer l'installation des modules npm et la compilation webpack pour compiler l'exécutable golang tout compris.
- intégrer toutes les dépendances fichier dans l'exécutable golang pour le rendre plus portable et faciliter l'installation

## Configuration

Voir `config-sample.toml` dans les sources.

Par ordre de priorité, le fichier de configuration peut se trouver dans:

- /etc/opensignauxfaibles/config.toml
- ~/.opensignauxfaibles/config.toml
- ./config.toml

## Usage

Le serveur `dbmongo` s'inscrit dans un workflow d'intégration de données. Pour plus d'informations, consulter [signaux-faibles/documentation](https://github.com/signaux-faibles/documentation/blob/master/processus-traitement-donnees.md#workflow-classique).

## Développement des fonctions map-reduce (TS/JS)

```sh
$ cd dbmongo/js
$ nvm use # pour utiliser la version de Node.js spécifiée dans .nvmrc
$ npm install # pour installer les dépendances
$ npm test # pour exécuter les tests unitaires et d'intégration, tel que décrit dans package.json
```

## Intégration et test de modifications des fonctions map-reduce (TS/JS)

```sh
$ ./test-all.sh # va regénérer jsFunctions.go, recompiler dbmongo et exécuter tous les tests
```
