#!/bin/bash
# build.sh compile les éléments et dispose le nécessaire dans le répertoire dist.
# ce script ne se charge pas d'installer les éléments nécessaires à la compilation.
mkdir -p dist/dbmongo 
mkdir dist/frontend 

# dbmongo: compilation et recopie des éléments nécessaires
cd dbmongo 
go build -ldflags="-s -w"
upx dbmongo
cp ./dbmongo ../dist/dbmongo/ 
cp -r js/ ../dist/dbmongo 
cp config.toml ../dist/dbmongo

# frontend: compilation et recopie
echo $PWD
cd ../frontend &&
yarn build
cp -r dist/* ../dist/frontend

# création de l'archive
cd ..
tar cvJf dist.tar.xz dist/

# cleanup

rm dbmongo/dbmongo
rm -r frontend/dist 
rm -r dist/