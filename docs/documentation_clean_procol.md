# Documentation de la table clean_procol

La table clean_procol stocke les données sur les évènements de procédures collectives.

- **siren**
- **date_effet**
- **action_procol**
- **stade_procol**
- **libelle_procol**

# Fonction auxiliaire

La fonction auxiliaire `procol_at_date` est stockée dans la base pour aider à 
déterminer les procédures en cours à une date donnée.

`
SELECT *
FROM procol_at_date('2026-01-01'::date);
`

/!\ Si théoriquement, la transformation d'une procol en une autre devrait 
cloturer la premières, les données ne sont pas parfaites et il arrive que 
plusieurs procédures collectives apparaissent comme actives à la même date.
La fonction ne gère pas encore bien cette situation et peut retourner 
plusieurs lignes pour une même entreprise. 

# Notes et limites connues

Ne figurent dans la base uniquement les entreprises qui ont eu des évènements 
de procédure collective. Par défaut, une entreprise absente des données doit 
donc être considérée comme "In bonis". 

Attention toutefois, certaines entreprises "In bonis" figurent dans la base, 
car une entreprise peut retourner au statut "In bonis" à la fin d'un plan de 
sauvegarde ou de redressement. Dans cette situation, `libelle_procol='In 
bonis'`.

Un plan de redressement peut durer 10 ans. Pour connaître précisément les 
entreprises en plan de redressement a une date précise, il faut donc 10 ans de 
profondeur historique antérieure à cette date.


