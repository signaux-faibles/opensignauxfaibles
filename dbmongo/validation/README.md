Chaque fichier de ce répertoire décrit les règles de validation qui doivent s'appliquer au type de données correspondant au moment de l'importation dans les collections `ImportedData` et `RawData`, avant que ces données ne soient traitées par les chaînes `public` et `reduce.algo2`.

Par exemple, le fichier `delai.schema.json` décrit les champs de l'entrée `delai` générée par le parseur des données URSSAF.

Ces règles sont exprimées dans la [version étendue par MongoDB de JSON Schema](https://docs.mongodb.com/manual/reference/operator/query/jsonSchema).

Après chaque modification d'un schema, penser à mettre à jour les définitions de types correspondants dans les fichiers `.go` et `.ts`.
