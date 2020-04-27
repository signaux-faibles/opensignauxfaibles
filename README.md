
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
$ go get -v -d -u github.com/signaux-faibles/opensignauxfaibles/dbmongo
$ cd $(go env GOPATH)/src/github.com/signaux-faibles/opensignauxfaibles/dbmongo
$ cd lib/engine
$ go generate # pour générer jsFunctions
$ cd -
$ go build
```

Dans l'arbre de sources de l'installation go, vous trouverez tous les répertoires nécessaires à l'exécution.

TODO:

- linker correctement les procédures R avec le core golang
- provoquer l'installation des modules npm et la compilation webpack pour compiler l'exécutable golang tout compris.
- intégrer toutes les dépendances fichier dans l'exécutable golang pour le rendre plus portable et faciliter l'installation

## Configuration

Voir config.toml.example dans les sources
Par ordre de priorité, le fichier de configuration peut se trouver dans:

- /etc/opensignauxfaibles/config.toml
- ~/.opensignauxfaibles/config.toml
- ./config.toml
