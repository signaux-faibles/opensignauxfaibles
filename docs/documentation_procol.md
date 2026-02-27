# Documentation de la table clean_procol

La table clean_procol stocke les données sur les évènements de procédures collectives.

- **siren**
- **date_effet** Date d'effet de l'évènement
- **action_procol** cf données brutes
- **stade_procol** cf données brutes
- **libelle_procol** Nom d'affichage de l'état de procédure collective après l'évènement.

Une ligne = un évènement.

Toute entreprise absente de cette base, ou avant le premier évènement, est par 
défaut In bonis.

La profondeur historique des données est depuis 2006, car un évènement de plan 
de continuation peut durer jusqu'à 10 ans. Cela permet donc de connaître avec 
précision la situation d'une entreprise en 2016, au début de la fenêtre de 
traitement des données. 

Attention, il ne suffit pas de prendre la dernière date d'effet pour connaître 
la situation de l'entreprise, on peut trouver des situations avec une 
chronologie du type "redressement : ouverture" 
> "liquidation : ouverture" > "redressement : inclusion autre procédure". Dans ce cas, 
la liquidation a toujours cours alors que le dernier évènement choronologique 
concerne la cloture du redressement. 

La fin automatique d'un plan de continuation au bout de 10 ans ne fait pas 
l'objet d'une annonce BODACC, donc ne constitue pas un évènement dans cette 
table. 

Contrairement à la table `stg_procol`, cette table est indexée à l'entreprise. 
La seule raison d'avoir un détail à l'établissement dans `stg_procol` est le 
stade `solde procédure`, dont la date d'effet peut varier d'un établissement à 
l'autre. Ce stade ne nous intéresse pas dans le cadre de Signaux Faibles. 

La fonction auxiliaire `procol_at_date` et la vue `clean_procol_at_date` 
décrites ci-dessous prémachent le travail de connaître la situation d'une 
entreprise à une date donnée.

Attention, certains évènements constituent un retour à l'état "In bonis". La 
présence d'un évènement avec `action_procol='redressement'` ne signifie par 
exemple pas que l'entreprise fait forcément l'objet d'un redressement 
*ultèrieurement* à l'évènement.

# Fonction auxiliaire

La fonction auxiliaire `procol_at_date` est stockée dans la base pour aider à 
déterminer les procédures en cours à une date donnée.

Exemple d'utilisation (renvoie toutes les entreprises).

`
SELECT *
FROM procol_at_date('2026-01-01'::date);
`

La logique derrière : 

- pour chaque action de procédure collective, regarder le dernier état connu
- filtrer les actions terminées, y compris le plan de continuation de plus de 
  10 ans. 
- en présence de plusieurs procédures, ne conserver que la plus critique 
  (liquidation > redressement > sauvegarde)

L'absence de retour de cette fonction pour une entreprise donnée signifie que 
l'entreprise est (par défaut) "In bonis" à la date donnée.

# Documentation de la table clean_procol_at_date

La vue matérialisée `clean_procol_at_date` facilite l'utilisation des données 
de procol en décomposant toutes les paires "siren" (qui a au moins un 
évènement) x "période" entre 2016 et aujourd'hui, et pour chaque paire, l'état 
de l'entreprise à cette période. 

L'absence de valeur pour une période donnée, ou l'absence de l'entreprise 
signifie que l'entreprise est par défaut "In bonis".

# Notes et limites connues


Un plan de redressement peut durer 10 ans. Pour connaître précisément les 
entreprises en plan de redressement a une date précise, il faut donc 10 ans de 
profondeur historique antérieure à cette date.

