# Documentation de la table clean_filter

La table clean_filter représente le périmètre (liste de SIREN) des données
préparées, la plupart des vues `clean_...` sont définies sur ce périmètre.

Il ne concerne pas les tables `clean_sirene` et `clean_sirene_ul` qui ne sont
pas filtrées pour les besoins du front-end.

Elle est construite sur la base d'un filtre initial sur la base de l'effectif
de l'entreprise `stg_import_filter`, auquel sont retirés un certain nombre de
cas de figures agrégés dans `siren_blacklist`.

# Calcul du périmètre

Le calcul du périmètre est **découplé de l'import**. Il se fait via la
commande `computePerimeter` :

```bash
./sfdata computePerimeter --batch <batch_key> --path <chemin_données>
```

Cette commande :
1. Lit le fichier `effectif_ent` présent dans le batch
2. Calcule le périmètre SIREN selon les critères d'effectif (voir ci-dessous)
3. Écrit le résultat dans la table `stg_filter_import` en base de données

Le périmètre doit être calculé **avant** de lancer un import. L'import ne
recalcule jamais le périmètre automatiquement : il utilise le filtre existant
en base ou un fichier filtre explicite.

**Workflow typique :**

```bash
# 1. Calculer ou mettre à jour le périmètre
./sfdata computePerimeter --batch 1902 --path ./data/

# 2. Importer les données (utilise le périmètre en base)
./sfdata import --batch 1902 --path ./data/
```

Options de `computePerimeter` :
- `--batch` : identifiant du batch (obligatoire)
- `--path` : répertoire contenant les données brutes
- `--batch-config` : chemin vers un fichier de configuration batch explicite

# Filtrages

Les filtrages effectués portent :

- sur l'effectif, ne sont conservées que les entreprises ayant atteint ou
  dépassé 10 salariés dans les 120 derniers mois.
- sur la nature juridique : sont filtrées les catégories juridiques associées
  à des organismes publiques
- sur l'activité principale de l'entreprise (/!\ ce n'est pas la même chose
  que l'activité du siège de l'entreprise)
- sur les entreprises ayant leur siège à l'étranger.

  Voir la définition de la vue `siren_blacklist_logic` pour consulter le
  détail des catégories juridiques et activités filtrées (actuellement
  consultable dans la [migration
  042](../lib/db/migrations/042_change_perimeter.sql), mais qui peut évoluer
  par la suite)
