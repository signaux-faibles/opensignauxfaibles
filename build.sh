#!/bin/bash
# build.sh compile les éléments et dispose le nécessaire dans le répertoire dist.
# ce script ne se charge pas d'installer les éléments nécessaires à la compilation.
mkdir -p dist/static

if ! [ -f config.toml ]; then
  echo config.toml non trouvé, abandon
  exit 1
fi

# compilation et recopie des éléments nécessaires
go build -ldflags="-s -w" -o sfdata
if command -v upx &> /dev/null
then
    upx sfdata
else
    echo "upx n'a pas été trouvé, le binaire ne sera pas compressé"
fi


cp sfdata ../dist
cp config.toml ../dist

# création de l'archive
tar cvzf dist.tar.gz dist/ &> /dev/null

# cleanup
rm sfdata
rm -r dist/

echo dist.tar.gz généré