#!/usr/bin/env bash

# Se mettre dans le répertoire d'output
# Lancer: "bash /chemin/du/script/dl_sirene.sh 71 25 21 etc.
# Pour télécharger et fusionner les fichiers des départements concernés

echo "$*"

for dep in "$@"
do
 wget "http://data.cquest.org/geo_sirene/v2019/2019-06/dep/geo_siret_${dep}.csv.gz"
 gzip -d "geo_siret_${dep}.csv.gz"
 echo "$(tail -n +2 geo_siret_${dep}.csv)" > "geo_siret_${dep}.csv"
done

cat geo_siret_*.csv >> "sirene_$(echo $* | tr ' ' '_').csv"




