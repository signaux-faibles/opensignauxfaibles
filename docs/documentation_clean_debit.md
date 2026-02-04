# Documentation de la table clean_debit

La table clean_debit représente le montant des *dettes* (agrégées) qu'un
établissement a sur ses cotisation sociales.

- **siret**
- **periode** cf ci-dessous
- **part_ouvriere** en €
- **part_patronale** en €
- **is_last** : s'agit-il de la dernière situation connue pour l'établissement ?

# Période de prise en compte

La **periode**, bien qu'elle soit systématiquement indiquée au premier du 
mois, concerne en réalité la dernière situation connue entre le 21 du mois 
précédent et le 20 du mois de la période. 

Si les données sont importées un 7 janvier, la dernière période ("2026-01-01") 
sera partielle et représentera la dernière situation connue au 7 janvier

La rationalité derrière cela :

- Les entreprises ont jusqu'au 15 du mois pour régler leurs cotisations 
  sociales.
- Une tolérence fait qu'un petit retard n'est pas systématiquement 
  immédiatement considéré comme une dette. 
- Une absence au 20 du mois sera néanmoins considéré comme une dette sur les 
  cotisations du mois en cours.

**Uniquement les périodes sur lesquels il y a eu des évènements sont présents 
dans la base.** Tant qu'il n'y a pas d'évènement qui vient actualiser le 
montant de la dette, la dette est considérée comme présente. Par exemple, si 
deux évènements sont présents, le premier au "2025-02-01" avec une dette de 
1000 et le "2025-06-01" avec une dette de 0, alors l'entreprise avait une 
dette de 1000 € sur les mois de février, mars, avril, mai 2025, qui a été 
remboursée entre le 21 mai et le 20 juin 2025.

Vu la taille de la base, il a été préféré de laisser au consommateur des 
données le soin d'effectuer le "forward fill" (au besoin).

# Fermeture d'établissements et affectation de la dette sociale

Lorsqu'un établissement est remplacé par un autre (par exemple, dans le cas 
d'un changement d'adresse), la dette va "suivre" ce changement 
d'établissement. Les dettes pour l'établissement fermé disparaîtront des 
données, et seront affectées au SIRET du nouvel établissement.
