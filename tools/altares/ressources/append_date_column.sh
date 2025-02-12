#!/bin/bash

# Comment utiliser ce script
# ./append_date_column.sh E_202403144975-0_202408010201.csv 2024-08-15

date=$2
header="REFERENCE_CLIENT;SIREN;SIRET;ETAT_ACTIVITE_ETP;PAYDEX;RETARD_MOYEN_PAIEMENTS;NOMBRE_FOURNISSEURS_ANALYSES;MONTANT_TOTAL_ENCOURS_ETUDIES;MONTANT_TOTAL_ENCOURS_ECHUS_NON_REGLES;FPI_30;FPI_90;ETAT_PROCEDURE_COLLECTIVE;DIFFUSIBLE;DATE_VALEUR"

# Ajouter l'en-tête au nouveau fichier
echo $header > "$date".csv

# Lire le fichier ligne par ligne et ajouter la date à la fin de chaque ligne
while IFS= read -r line
do
  echo "${line%%[[:space:]]}$date" >> "$date".csv
done < "$1"

sed -i '2d' "$date".csv