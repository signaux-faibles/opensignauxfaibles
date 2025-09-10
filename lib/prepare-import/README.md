# prepare-import

Scripts de préparation à l'importation de données dans le processus d'intégration de Signaux Faibles.

Lors de la constitution d'un batch, la commande `prepare-import` génère un document JSON destiné à être inséré dans la collection `Admin` de la base de données, à partir de fichiers de données mis à disposition dans un répertoire.

Elle vise à supporter tous les types de fichiers décrits dans le tableau fourni dans la [section "Spécificités de l'import" de la documentation](https://documentation/blob/master/processus-traitement-donnees.md#sp%C3%A9cificit%C3%A9s-de-limport) de Signaux Faibles.

La rencontre de fichiers non supportés n'empêchera pas la génération d'un batch, mais ceux-ci seront listés dans la sortie d'erreurs. (`stderr`)

## Dépendances

- Go
- `awk`

## Usage

```sh
make # Installe les dépendances, y compris de test (-t), et compile le binaire
make test # Exécute les tests
./prepare-import . # Retourne la définition du batch au format JSON, depuis le répertoire courant
```

Après toute modification du rendu de prepare-import, penser à mettre à jour le
golden file avec la commande:

```sh
go test --update
```

## Contribution

Nous suivons la specification [Conventional Commits](https://www.conventionalcommits.org/) pour le nommage des commits intégrés à la branche `master`. Ceci nous permet d'automatiser la génération de numéros de version avec [hekike/unchain: Tooling for conventional commit messages](https://github.com/hekike/unchain). (alternative à [semantic-release](https://github.com/semantic-release/semantic-release))
