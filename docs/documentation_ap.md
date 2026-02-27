# Documentation de la table clean_ap

La table clean_ap est une vue déjà agrégée des consommations et des demandes.

`etp_consomme=0` signifie que pas de consommation n'est connue, mais pourra 
potentiellement être modifié avec la mise à jour des données.


`is_last=true` lorsqu'il y a une consommation connue à la dernière date pour 
laquelle il y a des données de consommation dans les données source. 
