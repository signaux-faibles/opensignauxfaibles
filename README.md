
# Open Signaux Faibles

Solution logicielle pour la détection anticipée d'entreprises en difficulté

## Architecture

- backend golang
- frontend vuetify
- mongodb

## Dépendances / pré-requis

- `npx` (installé avec [Node.js](https://nodejs.org/)), pour lancer la transpilation des fichiers TypeScript vers JavaScript

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
